package repository

import (
	"context"
	"jose827corrza/go-cqrs/models"
)

type Repository interface {
	Close()
	InsertFeed(ctx context.Context, feed *models.Feed) error
	ListFeed(ctx context.Context) ([]*models.Feed, error)
}

//Implementacion abstracta
var implementation Repository

func SetRepository(r Repository) {
	implementation = r
}
func Close() {
	implementation.Close()
}
func InsertFeed(ctx context.Context, feed *models.Feed) error {
	return implementation.InsertFeed(ctx, feed)
}
func ListFeed(ctx context.Context) ([]*models.Feed, error) {
	return implementation.ListFeed(ctx)
}
