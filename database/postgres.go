package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/f3rcho/cqrs/models"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgressRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db}, nil
}

func (repo *PostgresRepository) Close() {
	repo.db.Close()
}

func (repo *PostgresRepository) InsertFeed(ctx context.Context, feed *models.Feed) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO feeds (id, title, description) VALUES ($1, $2, $3)", feed.ID, feed.Title, feed.Description)
	return err
}

func (repo *PostgresRepository) ListFeeds(ctx context.Context) ([]*models.Feed, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT * FROM feeds")
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var feeds []*models.Feed

	for rows.Next() {
		var feed = models.Feed{}
		if err = rows.Scan(&feed.ID, &feed.Title, &feed.Description, &feed.CreatedAt); err == nil {
			feeds = append(feeds, &feed)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return feeds, nil
}
