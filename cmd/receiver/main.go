package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/diptanw/countd-zmq/internal/platform/async"
	"github.com/diptanw/countd-zmq/internal/platform/logger"
	"github.com/diptanw/countd-zmq/internal/platform/zeromq"
	"github.com/diptanw/countd-zmq/internal/server"
	"github.com/diptanw/countd-zmq/internal/word"
)

func main() {
	config := readConfig()
	log := logger.New(os.Stdout, logger.Level(config.LogLevel))

	defer func() {
		if err := recover(); err != nil {
			log.Errorf("%s", err)
			os.Exit(1)
		}
	}()

	if err := serve(context.Background(), config, log); err != nil {
		panic(err)
	}
}

func serve(ctx context.Context, config Config, log logger.Logger) error {
	log.Infof("binding ZMQ socket listener to %q", config.ZmqBindAddr)

	sockReader, err := zeromq.Bind(config.ZmqBindAddr)
	if err != nil {
		return err
	}

	defer sockReader.Close()

	pool := async.NewPool()
	listener := server.NewStreamReader(sockReader, processMessages(log, pool))

	go func() {
		if err := pool.Run(ctx, 30); err != nil {
			log.Errorf("worker job error: %s", err)
		}
	}()

	return server.New(listener, log).Serve(ctx)
}

func processMessages(log logger.Logger, pool async.Pool) server.Handler {
	return func(message server.Message) error {
		pool.Enqueue(func(ctx context.Context) error {
			var aggr word.Aggregate

			if err := json.Unmarshal(message.Data, &aggr); err != nil {
				return err
			}

			text := fmt.Sprintf("Word: %q; Count: %d; First seen: %s; Last Seen: %s; Delta: %d",
				aggr.Word,
				aggr.TotalCount,
				aggr.FirstSeenAt.Format(time.RFC3339),
				aggr.LastSeenAt.Format(time.RFC3339),
				aggr.DeltaCount,
			)

			if aggr.DeltaCount == 0 {
				// Aggregates with zero delta with debug level for cleaner output.
				log.Debugf(text)
			} else {
				log.Infof(text)
			}

			return nil
		})

		return nil
	}
}
