package state

import (
	"log"
	"reflect"
	"time"

	"github.com/google/uuid"
)

type Effect string

type Event struct {
	ID       uuid.UUID
	EntityID uuid.UUID

	Effect     Effect
	Data       interface{}
	SystemData interface{}

	FromState string
	ToState   string

	CreatedAt time.Time
}

func ParseData[T interface{}](e Event) T {
	data, ok := e.Data.(T)
	if !ok {
		log.Fatalf("failed to parse data for effect[%s]: %s", e.Effect, reflect.TypeOf(data))
	}
	return data
}
