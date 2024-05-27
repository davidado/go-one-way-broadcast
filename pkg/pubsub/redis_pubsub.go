// Package pubsub provides a Redis client for publishing and subscribing messages.
package pubsub

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// RedisPubSub is a Redis client for publishing and subscribing messages.
type RedisPubSub struct {
	conn *redis.Client
}

// NewRedisPubSub creates a new RedisPubSub client.
func NewRedisPubSub(conn *redis.Client) *RedisPubSub {
	return &RedisPubSub{
		conn: conn,
	}
}

// Publish publishes a message to a topic.
func (ps *RedisPubSub) Publish(ctx context.Context, topic string, b []byte) error {
	err := ps.conn.Publish(ctx, topic, b).Err()
	if err != nil {
		return err
	}
	return nil
}

// Subscribe listens to a topic for incoming messages
// then adds them to the eventStream channel.
func (ps *RedisPubSub) Subscribe(ctx context.Context, topic string, eventStream chan []byte) {
	subscriber := ps.conn.Subscribe(ctx, topic)

	defer func() {
		err := subscriber.Unsubscribe(ctx, topic)
		if err != nil {
			log.Println("error cancelling consumer:", err)
		}
		fmt.Println("unsubscribed from topic", topic)
	}()

	fmt.Printf(" [*] %s waiting for messages.\n", topic)

	for {
		select {
		case <-ctx.Done():
			return
		case b := <-subscriber.Channel():
			eventStream <- []byte(b.Payload)
		}
	}
}
