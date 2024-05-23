// Sample code for sending server sent events (SSE) to a client.
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/davidado/go-one-way-broadcast/config"
	"github.com/davidado/go-one-way-broadcast/pkg/api"
	"github.com/davidado/go-one-way-broadcast/pkg/pubsub"

	"github.com/redis/go-redis/v9"
)

func initRedis(ctx context.Context, conn *redis.Client) {
	if err := conn.Ping(ctx).Err(); err != nil {
		log.Fatal("unable to ping Redis:", err)
	}
	log.Println("successfully connected to Redis")
}

func main() {
	ctx := context.Background()
	conn := redis.NewClient(&redis.Options{
		Addr:     config.Envs.RedisHost,
		Password: config.Envs.RedisPassword, // no password set
		DB:       0,                         // use default DB
	})
	initRedis(ctx, conn)
	defer conn.Close()

	pubsub := pubsub.NewRedisPubSub(conn)

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
