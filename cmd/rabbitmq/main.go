// Sample code for sending server sent events (SSE) to a client.
package main

import (
	"log"
	"net/http"

	"github.com/davidado/go-one-way-broadcast/config"
	"github.com/davidado/go-one-way-broadcast/pkg/api"
	"github.com/davidado/go-one-way-broadcast/pkg/pubsub"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial(config.Envs.RabbitMQHost)
	if err != nil {
		log.Fatal("unable to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	pubsub := pubsub.NewRabbitMQPubSub(conn)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{topic}", api.HandleTopic)
	mux.HandleFunc("GET /events/{topic}", api.HandleEvents(pubsub))
	mux.HandleFunc("GET /push/{topic}", api.HandlePush(pubsub))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	srv.ListenAndServe()
}
