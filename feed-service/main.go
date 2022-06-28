package main

import (
	"fmt"
	"jose827corrza/go-cqrs/database"
	"jose827corrza/go-cqrs/events"
	"jose827corrza/go-cqrs/repository"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PostgresDB       string `json:"POSTGRES_DB"`
	PostgresUser     string `json:"POSTGRES_USER"`
	PostgresPassword string `json:"POSTGRES_PASSWORD"`
	NatsAddress      string `json:"NATS_ADDRESS"`
}

func main() {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatalf("%v", err)
	}

	addr := fmt.Sprintf("postgress://%s:%s@postgress/%s?sslmode=disable",
		config.PostgresUser,
		config.PostgresPassword,
		config.PostgresDB)
	repo, err := database.NewPostgresRepository(addr)
	if err != nil {
		log.Fatal(err)
	}
	repository.SetRepository(repo)

	//Ahora el NATS
	natAddr := fmt.Sprintf("nats://%s", config.NatsAddress)
	n, err := events.NewNats(natAddr)
	if err != nil {
		log.Fatal(err)
	}
	events.SetEventStore(n)

	defer events.Close()

	//Router
	router := newRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

func newRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/feeds", CreateFeedHandler).Methods(http.MethodPost)
	return
}
