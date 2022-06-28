package main

import (
	"encoding/json"
	"jose827corrza/go-cqrs/events"
	"jose827corrza/go-cqrs/models"
	"jose827corrza/go-cqrs/repository"
	"log"
	"net/http"
	"time"

	"github.com/segmentio/ksuid"
)

type CreateFeedRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func CreateFeedHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateFeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//Se supone el insert no necesita esto..
	createdAt := time.Now().UTC()
	id, err := ksuid.NewRandom()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	feed := models.Feed{
		ID:          id.String(),
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   createdAt,
	}
	//Envia a DB
	if err := repository.InsertFeed(r.Context(), &feed); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//Transmitir el feed que se acaba de crear a NATS
	if err := events.PublishCreatedFeed(r.Context(), &feed); err != nil {
		log.Printf("Failed to publish created feed event: %v", err)
	}
	//Ya esto es la respuesta al cliente
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&feed) //Revisar, aca puede ser sin puntero
}
