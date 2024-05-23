# Overview

A one-way message broadcasting application implemented in Go that handles server sent events (SSE) to connected clients over a swappable message broker.

## Features

The primary use case for this application is for sending one-way real-time notifications to connected clients through server sent events. However, the program can easily be expanded to include other types of live data changes including stock ticker or progress updates by modifying the `message.proto` file in the `pb` directory (refer to the protobuf section for more details).

Websockets are a viable alternative to server sent events for accomplishing the same use cases but unlike websockets, SSE has the main advantage of being built on top of HTTP. As such, SSE retains advantages over websockets including http/2 multiplexing, compression, security, and implementation simplicity.

This application is deployable to multiple servers behind a load balancer. Messages are published to brokers such as Redis, RabbitMQ, or one of your own implementation.

## Prerequisites

- Go 1.22+ since the project uses its routing features but it should be fairly straightforward to implement your own router if using earlier versions of Go.

- A running instance of Redis or RabbitMQ. If you do not already have a running instance, you can run `docker run -d --name redis -p 6379:6379 redis` to run a local Redis instance in Docker.

## Installation 

1. Download or clone the project to your workspace: `git clone https://github.com/davidado/go-one-way-broadcast`

2. Edit the Makefile and set the MAIN_PACKAGE_PATH to either `./cmd/redis` or `./cmd/rabbitmq` (default is `./cmd/redis`).

## Usage

1. Execute `make run` in the root path (same directory as the Makefile).

2. Go to http://localhost:8080/[topic] where [topic] can be anything you choose.

3. Go to http://localhost:8080/push/[topic] to push data into the message broker which is then displayed into the page you entered above.

4. Listen to http://localhost:8080/events/[topic] in your own application for data updates (refer to `/public/index.html`).

## Protobuf

This application uses protocol buffers to allow for multiple services to publish messages with the message broker regardless of programming language. Modify `/pb/message.proto` if custom data is required, install protocol buffers in your system, and run `protoc -I=$SRC_DIR --go_out=$DST_DIR $SRC_DIR/message.proto` (refer to https://protobuf.dev/getting-started/gotutorial for installation and usage details).

## Contributing

Contributions are welcome.
