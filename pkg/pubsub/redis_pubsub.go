// Package pubsub provides a Redis client for publishing and subscribing messages.
package pubsub

import (
	"context"
	"fmt"
	"log"

	"github.com/davidado/go-one-way-broadcast/pb"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
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
func (ps *RedisPubSub) Publish(ctx context.Context, topic string, m *pb.Message) error {
	b, err := proto.Marshal(m)
	if err != nil {
		log.Println("error marshalling message:", err)
		return err
	}
	err = ps.conn.Publish(ctx, topic, b).Err()
	if err != nil {
		return err
	}
	return nil
}

// Subscribe listens to a topic for incoming messages
// then adds them to the eventStream channel.
func (ps *RedisPubSub) Subscribe(ctx context.Context, topic string, eventStream chan *pb.Message) {
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
			msg := &pb.Message{}
			err := proto.Unmarshal([]byte(b.Payload), msg)
			if err != nil {
				log.Println("error unmarshalling message:", err)
				continue
			}

			eventStream <- msg
		}
	}
}
