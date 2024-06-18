package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/f3rcho/cqrs/events"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	NatsAddress string `envconfig:"NATS_ADDRESS"`
}

func main() {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}
	hub := NewHub()

	n, err := events.NewNats(fmt.Sprintf("nats://%s", config.NatsAddress))
	if err != nil {
		log.Fatal(err)
	}

	err = n.OnCreatedFeed(func(m events.CreatedFeedMessage) {
		hub.Broadcast(newCreatedFeedMessage(m.ID, m.Title, m.Description, m.CreatedAt), nil)
	})
	if err != nil {
		log.Fatal(err)
	}
	events.SetEventStore(n)
	defer events.Close()

	go hub.Run()
	http.HandleFunc("/ws", hub.HandleWebSocket)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
