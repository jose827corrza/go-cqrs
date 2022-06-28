package main

import (
	"fmt"
	"jose827corrza/go-cqrs/events"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	NatsAddress string `json:"NATS_ADDRESS"`
}

func main() {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatalf("%v", err)
	}
	hub := NewHub()

	//NATS
	natAddr := fmt.Sprintf("nats://%s", config.NatsAddress)
	n, err := events.NewNats(natAddr)
	if err != nil {
		log.Fatal(err)
	}

	err = n.OnCreateFeed(func(m events.CreatedFeedMessage) {
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
