package search

import (
	"context"

	"github.com/f3rcho/cqrs/models"
)

type SearchRepository interface {
	IndexFeed(ctx context.Context, feed models.Feed) error
	SearchFeed(ctx context.Context, query string) ([]models.Feed, error)
}

var repository SearchRepository

func SetSearchRepository(repo SearchRepository) {
	repository = repo
}

func IndexFeed(ctx context.Context, feed models.Feed) error {
	return repository.IndexFeed(ctx, feed)
}

func SearchFeed(ctx context.Context, query string) ([]models.Feed, error) {
	return repository.SearchFeed(ctx, query)
}
