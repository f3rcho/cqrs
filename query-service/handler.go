package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/f3rcho/cqrs/events"
	"github.com/f3rcho/cqrs/models"
	"github.com/f3rcho/cqrs/repository"
	"github.com/f3rcho/cqrs/search"
)

func onCreatedFeed(m events.CreatedFeedMessage) {
	feed := models.Feed{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
	}
	if err := search.IndexFeed(context.Background(), feed); err != nil {
		log.Println(err)
	}
}

func listFeedsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	feeds, err := repository.ListFeeds(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feeds)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	query := r.URL.Query().Get("q")
	if len(query) == 0 {
		http.Error(w, "query is required", http.StatusBadRequest)
		return
	}
	feeds, err := search.SearchFeed(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feeds)
}
