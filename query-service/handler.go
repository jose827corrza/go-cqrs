package main

import (
	"context"
	"encoding/json"
	"jose827corrza/go-cqrs/events"
	"jose827corrza/go-cqrs/models"
	"jose827corrza/go-cqrs/repository"
	"jose827corrza/go-cqrs/search"
	"log"
	"net/http"
)

func onCreatedFeed(msg events.CreatedFeedMessage) {
	feed := models.Feed{
		ID:          msg.ID,
		Title:       msg.Title,
		Description: msg.Description,
		CreatedAt:   msg.CreatedAt,
	}
	if err := search.IndexFeed(context.Background(), feed); err != nil {
		log.Fatalf("Failed to index feed: %v", err)
	}
}

func listFeedHandler(w http.ResponseWriter, req *http.Request) {
	var err error
	feeds, err := repository.ListFeed(req.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feeds)
}

func searchHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var err error
	//Esto es el manejo de un quiery param
	// ... /search?q<query_value>
	query := req.URL.Query().Get("q")
	if len(query) == 0 {
		http.Error(w, "Query is required", http.StatusBadRequest)
		return
	}

	feeds, err := search.SearchFeed(ctx, query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feeds)
}
