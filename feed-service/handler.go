package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/f3rcho/cqrs/events"
	"github.com/f3rcho/cqrs/models"
	"github.com/f3rcho/cqrs/repository"
	"github.com/segmentio/ksuid"
)

type CreateFeedRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func handlerCreateFeed(w http.ResponseWriter, r *http.Request) {
	var req CreateFeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdAt := time.Now().UTC()
	id, err := ksuid.NewRandom()
	if err != nil {
		http.Error(w, "failed to create feed", http.StatusInternalServerError)
		return
	}

	feed := models.Feed{
		ID:          id.String(),
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   createdAt,
	}

	if err := repository.InsertFeed(r.Context(), &feed); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := events.PublishCreatedFeed(r.Context(), &feed); err != nil {
		log.Printf("failed to publish created feed event: %v\n", err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feed)
}
