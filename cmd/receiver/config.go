package main

import (
	"flag"

	"github.com/diptanw/countd-zmq/internal/platform/logger"
)

// Config is a struct that contains service configuration.
type Config struct {
	LogLevel    int
	ZmqBindAddr string
}

func readConfig() Config {
	var config Config

	flag.StringVar(&config.ZmqBindAddr, "zmq-bind-addr", "tcp://*:6501", "an address of ZMQ socket listener")
	flag.IntVar(&config.LogLevel, "log-level", int(logger.Info), "a logging level")
	flag.Parse()

	return config
}
