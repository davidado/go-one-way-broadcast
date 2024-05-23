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
)

// PubSubber is an interface for subscribing to a topic and
// publishing messages to it.
type PubSubber interface {
	Publish(ctx context.Context, topic string, m *pb.Message) error
	Subscribe(ctx context.Context, topic string, eventStream chan *pb.Message)
}

// HandleTopic serves the index.html file.
func HandleTopic(w http.ResponseWriter, r *http.Request) {
	dir, _ := os.Getwd()
	http.ServeFile(w, r, dir+"/public/index.html")
}

// HandleEvents handles server sent events (SSE) for a topic.
func HandleEvents(pubsub PubSubber) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventStream := make(chan *pb.Message)
		defer close(eventStream)

		topic := r.PathValue("topic")
		go pubsub.Subscribe(r.Context(), topic, eventStream)

		sendServerEvents(w, r, eventStream)
	}
}

func sendServerEvents(w http.ResponseWriter, r *http.Request, eventStream chan *pb.Message) {
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
		case m := <-eventStream:
			// Flush the message to the client.
			j, err := json.Marshal(struct {
				Topic string `json:"topic"`
				Title string `json:"title"`
				Body  string `json:"body"`
			}{
				Topic: m.Topic,
				Title: m.Title,
				Body:  m.Body,
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
			m := &pb.Message{
				Topic: topic,
				Title: "New Message!",
				Body:  strconv.Itoa(i),
			}
			pubsub.Publish(r.Context(), topic, m)
			time.Sleep(1 * time.Second)
			fmt.Println(" [*] Published message to topic:", string(m.Body))
		}
	}
}
