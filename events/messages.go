package events

import "time"

type Message interface {
	Type() string
}

//Esta es la que va a ser transportada por NATS y que sea procesada por los dif servicios
type CreatedFeedMessage struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

//Implementar la interface
func (m CreatedFeedMessage) Type() string {
	return "created_feed"
}
