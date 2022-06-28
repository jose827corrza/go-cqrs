package events

import (
	"context"
	"jose827corrza/go-cqrs/models"
)

type EventStore interface {
	Close()
	PublishCreatedFeed(ctx context.Context, feed *models.Feed) error             //Para crear
	SubscribeCreatedFeed(ctx context.Context) (<-chan CreatedFeedMessage, error) //Para suscribirse a un canal
	OnCreateFeed(f func(CreatedFeedMessage)) error                               //callback para que reacciones cuando un feed ha sido creado
}

//Muy parecido al repository
var eventStore EventStore

func SetEventStore(store EventStore) {
	eventStore = store
}
func Close() {
	eventStore.Close()
}

func PublishCreatedFeed(ctx context.Context, feed *models.Feed) error {
	return eventStore.PublishCreatedFeed(ctx, feed)
}

func SubscribeCreatedFeed(ctx context.Context) (<-chan CreatedFeedMessage, error) {
	return eventStore.SubscribeCreatedFeed(ctx)
}

func OnCreateFeed(ctx context.Context, f func(CreatedFeedMessage)) error {
	return eventStore.OnCreateFeed(f)
}
