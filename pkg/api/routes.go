// Package api provides the HTTP handlers for the API.
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/davidado/go-one-way-broadcast/pb"
	"google.golang.org/protobuf/proto"
)

// PubSubber is an interface for subscribing to a topic and
// publishing messages to it.
type PubSubber interface {
	Publish(ctx context.Context, topic string, b []byte) error
	Subscribe(ctx context.Context, topic string, eventStream chan []byte)
}

// HandleTopic serves the index.html file.
func HandleTopic(w http.ResponseWriter, r *http.Request) {
	dir, _ := os.Getwd()
	http.ServeFile(w, r, dir+"/public/index.html")
}

// HandleEvents handles server sent events (SSE) for a topic.
func HandleEvents(pubsub PubSubber) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventStream := make(chan []byte)
		defer close(eventStream)

		topic := r.PathValue("topic")
		go pubsub.Subscribe(r.Context(), topic, eventStream)

		sendServerEvents(w, r, eventStream)
	}
}

func sendServerEvents(w http.ResponseWriter, r *http.Request, eventStream chan []byte) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	// Set the headers.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		select {
		case <-r.Context().Done():
			return
		case b := <-eventStream:
			msg := &pb.Message{}
			err := proto.Unmarshal(b, msg)
			if err != nil {
				log.Println("error unmarshalling message:", err)
				continue
			}
			// Flush the message to the client.
			j, err := json.Marshal(struct {
				Topic string `json:"topic"`
				Title string `json:"title"`
				Body  string `json:"body"`
			}{
				Topic: msg.Topic,
				Title: msg.Title,
				Body:  msg.Body,
			})
			if err != nil {
				log.Println("error marshalling message:", err)
				continue
			}
			fmt.Fprintf(w, "data: %s\n\n", j)
			flusher.Flush()
		}
	}
}

// HandlePush publishes messages to a topic.
func HandlePush(pubsub PubSubber) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		topic := r.PathValue("topic")
		for i := 0; i < 10; i++ {
			msg := &pb.Message{
				Topic: topic,
				Title: "New Message!",
				Body:  strconv.Itoa(i),
			}
			b, err := proto.Marshal(msg)
			if err != nil {
				log.Println("error marshalling message:", err)
				return
			}
			pubsub.Publish(r.Context(), topic, b)
			time.Sleep(1 * time.Second)
			fmt.Println(" [*] Published message to topic:", string(msg.Body))
		}
	}
}
