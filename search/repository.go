package search

import (
	"context"

	"github.com/f3rcho/cqrs/models"
)

type Repository interface {
	InsertFeed(ctx context.Context, feed *models.Feed) error
	ListFeeds(ctx context.Context) ([]*models.Feed, error)
}

var repository Repository

func SetRepository(repo Repository) {
	repository = repo
}

func InsertFeed(ctx context.Context, feed *models.Feed) error {
	return repository.InsertFeed(ctx, feed)
}

func ListFeeds(ctx context.Context) ([]*models.Feed, error) {
	return repository.ListFeeds(ctx)
}
