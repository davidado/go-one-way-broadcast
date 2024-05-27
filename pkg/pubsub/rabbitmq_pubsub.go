// Package pubsub provides a publish/subscribe interface for sending messages.
package pubsub

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQPubSub is a Redis client for publishing and subscribing messages.
type RabbitMQPubSub struct {
	conn *amqp.Connection
}

// NewRabbitMQPubSub creates a new RabbitMQPubSub client.
func NewRabbitMQPubSub(conn *amqp.Connection) *RabbitMQPubSub {
	return &RabbitMQPubSub{
		conn: conn,
	}
}

func configureExchange(ch *amqp.Channel, topic string) error {
	err := ch.ExchangeDeclare(
		topic,    // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return err
	}
	return nil
}

// Publish publishes a message to a topic.
func (ps *RabbitMQPubSub) Publish(ctx context.Context, topic string, b []byte) error {
	ch, err := ps.conn.Channel()
	if err != nil {
		log.Println("error creating channel:", err)
	}
	defer ch.Close()

	err = configureExchange(ch, topic)
	if err != nil {
		log.Println("error declaring exchange:", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(
		ctx,
		topic, // exchange
		"",    // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        b,
		})
	if err != nil {
		log.Println("error publishing message:", err)
		return err
	}
	return nil
}

// Subscribe listens to a topic for incoming messages
// then adds them to the eventStream channel.
func (ps *RabbitMQPubSub) Subscribe(ctx context.Context, topic string, eventStream chan []byte) {
	ch, err := ps.conn.Channel()
	if err != nil {
		log.Println("error creating channel:", err)
		return
	}

	defer func() {
		err := ch.Cancel(topic, false)
		if err != nil {
			log.Println("error cancelling consumer:", err)
		}
		ch.Close()
		log.Println("unsubscribed from topic", topic)
	}()

	err = configureExchange(ch, topic)
	if err != nil {
		log.Println("failed to declare an exchange:", err)
		return
	}

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Println("failed to declare a queue:", err)
		return
	}

	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		topic,  // exchange
		false,
		nil,
	)
	if err != nil {
		log.Println("failed to bind a queue:", err)
		return
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	log.Printf(" [*] %s waiting for messages.", topic)

	for {
		select {
		case <-ctx.Done():
			return
		case d := <-msgs:
			eventStream <- d.Body
		}
	}
}
