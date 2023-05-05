# Aggregated Word Count

- [Overview](#overview)
- [Architecture](#architecture)
- [Development](#development)
- [Setup and Testing](#setup-and-testing)

## Overview

The solution intended to demonstrate knowledge of Go Programing Language using Go's standard library. Some parts of the
production running service suppose to be replaced by more robust and powerful open-source libraries, e.g.:

- [Viper](https://github.com/spf13/viper) Powerful framework for building CLI applications.
- [Zap](https://github.com/uber-go/zap) Fast, structured, leveled logging.

All components are designed to be testable, but no tests are provided at this moment (I will add them later).

## Architecture

The solution implements [Pipeline](https://zeromq.org/socket-api/?language=go&library=zmq4#pipeline-pattern) pattern
with `PUSH` and `PULL` Socket Types for service communication.

```
[stdin]--> emitter --push--> [ZMQ] <--pull-- countd --push--> [ZMQ] <--pull-- receiver
```

## Development

This project follows design and development principles described in:

- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Package Oriented Design](https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html)

**Go 1.18 is required for type parameters support! (use previous commit for Go 1.17)**

To build the application (zeromq library must be installed in the target platform):

```go
go build -o ./build/countd ./cmd/countd
```

To run tests with coverage:
```go
go test -race -coverprofile=build/coverage.out ./... && go tool cover -html=build/coverage.out
```

## Setup and Testing

To run the application using docker-compose:

```sh
docker-compose up --build
```

Example output:

```bash
zmq-poc-countd-receiver-1  | INF: binding ZMQ socket listener to "tcp://*:6500"
zmq-poc-countd-receiver-1  | INF: server: starting...
zmq-poc-countd-countd-1    | INF: binding ZMQ socket listener to "tcp://*:6500"
zmq-poc-countd-countd-1    | INF: connecting ZMQ socket receiver to "tcp://receiver:6500"
zmq-poc-countd-countd-1    | INF: server: starting...
zmq-poc-countd-emitter-1   | 2022/03/14 11:38:28 Connecting to tcp://countd:6500
zmq-poc-countd-emitter-1   | 2022/03/14 11:38:28 Read 367 words
zmq-poc-countd-receiver-1  | INF: Word: "wildness--to"; Count: 1; First seen: 2022-03-14T11:38:28Z; Last Seen: 2022-03-14T11:38:28Z
zmq-poc-countd-receiver-1  | INF: Word: "is"; Count: 1; First seen: 2022-03-14T11:38:29Z; Last Seen: 2022-03-14T11:38:29Z
zmq-poc-countd-receiver-1  | INF: Word: "see"; Count: 1; First seen: 2022-03-14T11:38:29Z; Last Seen: 2022-03-14T11:38:29Z
zmq-poc-countd-receiver-1  | INF: Word: "the"; Count: 1; First seen: 2022-03-14T11:38:31Z; Last Seen: 2022-03-14T11:38:31Z
zmq-poc-countd-receiver-1  | INF: Word: "lurk"; Count: 1; First seen: 2022-03-14T11:38:33Z; Last Seen: 2022-03-14T11:38:33Z
zmq-poc-countd-receiver-1  | INF: Word: "freely"; Count: 1; First seen: 2022-03-14T11:38:35Z; Last Seen: 2022-03-14T11:38:35Z
zmq-poc-countd-receiver-1  | INF: Word: "serenely"; Count: 1; First seen: 2022-03-14T11:38:35Z; Last Seen: 2022-03-14T11:38:35Z
zmq-poc-countd-receiver-1  | INF: Word: "be"; Count: 1; First seen: 2022-03-14T11:38:35Z; Last Seen: 2022-03-14T11:38:35Z
zmq-poc-countd-receiver-1  | INF: Word: "we"; Count: 1; First seen: 2022-03-14T11:38:35Z; Last Seen: 2022-03-14T11:38:35Z
zmq-poc-countd-receiver-1  | INF: Word: "ground"; Count: 1; First seen: 2022-03-14T11:38:37Z; Last Seen: 2022-03-14T11:38:37Z
zmq-poc-countd-receiver-1  | INF: Word: "tortoises"; Count: 1; First seen: 2022-03-14T11:38:40Z; Last Seen: 2022-03-14T11:38:40Z
zmq-poc-countd-receiver-1  | INF: Word: "the"; Count: 2; First seen: 2022-03-14T11:38:31Z; Last Seen: 2022-03-14T11:38:39Z
zmq-poc-countd-receiver-1  | INF: Word: "be"; Count: 2; First seen: 2022-03-14T11:38:35Z; Last Seen: 2022-03-14T11:38:45Z
zmq-poc-countd-receiver-1  | INF: Word: "need"; Count: 1; First seen: 2022-03-14T11:38:46Z; Last Seen: 2022-03-14T11:38:46Z
zmq-poc-countd-receiver-1  | INF: Word: "lurk"; Count: 2; First seen: 2022-03-14T11:38:33Z; Last Seen: 2022-03-14T11:38:47Z
zmq-poc-countd-receiver-1  | INF: Word: "prey"; Count: 1; First seen: 2022-03-14T11:38:47Z; Last Seen: 2022-03-14T11:38:47Z
zmq-poc-countd-receiver-1  | INF: Word: "refreshed"; Count: 1; First seen: 2022-03-14T11:38:49Z; Last Seen: 2022-03-14T11:38:49Z
zmq-poc-countd-receiver-1  | INF: Word: "out"; Count: 1; First seen: 2022-03-14T11:38:49Z; Last Seen: 2022-03-14T11:38:49Z
zmq-poc-countd-receiver-1  | INF: Word: "compensation"; Count: 1; First seen: 2022-03-14T11:38:50Z; Last Seen: 2022-03-14T11:38:50Z
zmq-poc-countd-receiver-1  | INF: Word: "and"; Count: 1; First seen: 2022-03-14T11:38:50Z; Last Seen: 2022-03-14T11:38:50Z
zmq-poc-countd-receiver-1  | INF: Word: "feeding"; Count: 1; First seen: 2022-03-14T11:38:52Z; Last Seen: 2022-03-14T11:38:52Z
zmq-poc-countd-receiver-1  | INF: Word: "which"; Count: 1; First seen: 2022-03-14T11:38:52Z; Last Seen: 2022-03-14T11:38:52Z
zmq-poc-countd-receiver-1  | INF: Word: "untenable"; Count: 1; First seen: 2022-03-14T11:38:52Z; Last Seen: 2022-03-14T11:38:52Z
zmq-poc-countd-receiver-1  | INF: Word: "where"; Count: 1; First seen: 2022-03-14T11:38:53Z; Last Seen: 2022-03-14T11:38:53Z
zmq-poc-countd-receiver-1  | INF: Word: "must"; Count: 1; First seen: 2022-03-14T11:38:55Z; Last Seen: 2022-03-14T11:38:55Z
zmq-poc-countd-receiver-1  | INF: Word: "hear"; Count: 1; First seen: 2022-03-14T11:38:58Z; Last Seen: 2022-03-14T11:38:58Z
zmq-poc-countd-receiver-1  | INF: Word: "village"; Count: 1; First seen: 2022-03-14T11:39:00Z; Last Seen: 2022-03-14T11:39:00Z
zmq-poc-countd-receiver-1  | INF: Word: "the"; Count: 3; First seen: 2022-03-14T11:38:31Z; Last Seen: 2022-03-14T11:38:57Z
zmq-poc-countd-receiver-1  | INF: Word: "life"; Count: 1; First seen: 2022-03-14T11:39:03Z; Last Seen: 2022-03-14T11:39:03Z
zmq-poc-countd-receiver-1  | INF: Word: "enough"; Count: 1; First seen: 2022-03-14T11:39:06Z; Last Seen: 2022-03-14T11:39:06Z
zmq-poc-countd-receiver-1  | INF: Word: "with"; Count: 1; First seen: 2022-03-14T11:39:08Z; Last Seen: 2022-03-14T11:39:08Z
zmq-poc-countd-receiver-1  | INF: Word: "rife"; Count: 1; First seen: 2022-03-14T11:39:08Z; Last Seen: 2022-03-14T11:39:08Z
zmq-poc-countd-receiver-1  | INF: Word: "when"; Count: 1; First seen: 2022-03-14T11:39:13Z; Last Seen: 2022-03-14T11:39:13Z
zmq-poc-countd-receiver-1  | INF: Word: "to"; Count: 1; First seen: 2022-03-14T11:39:14Z; Last Seen: 2022-03-14T11:39:14Z
zmq-poc-countd-receiver-1  | INF: Word: "be"; Count: 3; First seen: 2022-03-14T11:38:35Z; Last Seen: 2022-03-14T11:39:06Z
zmq-poc-countd-receiver-1  | INF: Word: "bear"; Count: 1; First seen: 2022-03-14T11:39:18Z; Last Seen: 2022-03-14T11:39:18Z
zmq-poc-countd-receiver-1  | INF: Word: "on"; Count: 1; First seen: 2022-03-14T11:39:18Z; Last Seen: 2022-03-14T11:39:18Z
zmq-poc-countd-receiver-1  | INF: Word: "pulp--tadpoles"; Count: 1; First seen: 2022-03-14T11:39:19Z; Last Seen: 2022-03-14T11:39:19Z
zmq-poc-countd-receiver-1  | INF: Word: "that"; Count: 1; First seen: 2022-03-14T11:39:19Z; Last Seen: 2022-03-14T11:39:19Z
zmq-poc-countd-receiver-1  | INF: Word: "see"; Count: 2; First seen: 2022-03-14T11:38:29Z; Last Seen: 2022-03-14T11:39:13Z
zmq-poc-countd-receiver-1  | INF: Word: "and"; Count: 2; First seen: 2022-03-14T11:38:50Z; Last Seen: 2022-03-14T11:39:16Z
zmq-poc-countd-receiver-1  | INF: Word: "time"; Count: 1; First seen: 2022-03-14T11:39:23Z; Last Seen: 2022-03-14T11:39:23Z
zmq-poc-countd-receiver-1  | INF: Word: "will"; Count: 1; First seen: 2022-03-14T11:39:25Z; Last Seen: 2022-03-14T11:39:25Z
zmq-poc-countd-receiver-1  | INF: Word: "are"; Count: 1; First seen: 2022-03-14T11:39:27Z; Last Seen: 2022-03-14T11:39:27Z
zmq-poc-countd-receiver-1  | INF: Word: "belly"; Count: 1; First seen: 2022-03-14T11:39:28Z; Last Seen: 2022-03-14T11:39:28Z
zmq-poc-countd-receiver-1  | INF: Word: "liability"; Count: 1; First seen: 2022-03-14T11:39:28Z; Last Seen: 2022-03-14T11:39:28Z
zmq-poc-countd-receiver-1  | INF: Word: "the"; Count: 5; First seen: 2022-03-14T11:38:31Z; Last Seen: 2022-03-14T11:39:29Z
zmq-poc-countd-receiver-1  | INF: Word: "a"; Count: 1; First seen: 2022-03-14T11:39:32Z; Last Seen: 2022-03-14T11:39:32Z
zmq-poc-countd-receiver-1  | INF: Word: "impression"; Count: 1; First seen: 2022-03-14T11:39:34Z; Last Seen: 2022-03-14T11:39:34Z
zmq-poc-countd-receiver-1  | INF: Word: "toads"; Count: 1; First seen: 2022-03-14T11:39:34Z; Last Seen: 2022-03-14T11:39:34Z
zmq-poc-countd-receiver-1  | INF: Word: "to"; Count: 2; First seen: 2022-03-14T11:39:14Z; Last Seen: 2022-03-14T11:39:24Z
zmq-poc-countd-receiver-1  | INF: Word: "three"; Count: 1; First seen: 2022-03-14T11:39:35Z; Last Seen: 2022-03-14T11:39:35Z
zmq-poc-countd-receiver-1  | INF: Word: "we"; Count: 2; First seen: 2022-03-14T11:38:35Z; Last Seen: 2022-03-14T11:39:26Z
zmq-poc-countd-receiver-1  | INF: Word: "thunder-cloud"; Count: 1; First seen: 2022-03-14T11:39:37Z; Last S
```

Note, the default configuration doesn't output aggregates with 0 delta count, change the `receiver` logging level to
"debug=3" for verbose output.
