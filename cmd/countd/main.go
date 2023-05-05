package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/diptanw/countd-zmq/internal/platform/async"
	"github.com/diptanw/countd-zmq/internal/platform/logger"
	"github.com/diptanw/countd-zmq/internal/platform/storage"
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

	log.Infof("connecting ZMQ socket receiver to %q", config.ZmqReceiverAddr)

	sockWriter, err := zeromq.Connect(config.ZmqReceiverAddr)
	if err != nil {
		return err
	}

	defer sockWriter.Close()

	hot, err := readHotList(config.HotList)
	if err != nil {
		return err
	}

	store := storage.NewInMemory[*word.Aggregate]()
	processor := word.NewProcessor(
		word.NewRepository(store),
		json.NewEncoder(&async.ConcurrentWriter{Writer: sockWriter}),
		async.NewScheduler(log), hot)

	listener := server.NewStreamReader(sockReader, debugMessages(log), func(m server.Message) error {
		if len(m.Data) == 0 {
			return errors.New("empty message")
		}

		return processor.Aggregate(ctx, string(m.Data))
	})

	return server.New(listener, log).Serve(ctx)
}

// debugMessages logs raw socket messages when debug logging is enabled.
func debugMessages(log logger.Logger) server.Handler {
	return func(message server.Message) error {
		log.Debugf("sock: %s", message.Data)

		return nil
	}
}

func readHotList(path string) ([]string, error) {
	if path == "" {
		return nil, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening hot-list file: %w", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var words []string

	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := file.Close(); err != nil {
		return nil, fmt.Errorf("closing hot-list file: %w", err)
	}

	return words, nil
}
