package main

import (
	"flag"

	"github.com/diptanw/countd-zmq/internal/platform/logger"
)

// Config is a struct that contains service configuration.
type Config struct {
	LogLevel        int
	ZmqBindAddr     string
	ZmqReceiverAddr string
	HotList         string
}

func readConfig() Config {
	var config Config

	flag.StringVar(&config.ZmqBindAddr, "zmq-bind-addr", "tcp://*:6500", "an address of ZMQ socket listener")
	flag.StringVar(&config.ZmqReceiverAddr, "zmq-receiver-addr", "tcp://localhost:6501", "an address of ZMQ socket receiver")
	flag.IntVar(&config.LogLevel, "log-level", int(logger.Info), "a logging level")
	flag.StringVar(&config.HotList, "hot-list-path", "hotlist", "a hot-list file path")
	flag.Parse()

	return config
}
