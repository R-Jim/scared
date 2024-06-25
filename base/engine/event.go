package engine

import (
	"time"

	"github.com/google/uuid"
)

type Effect string

const (
	EffectInit Effect = "INIT"
)

type Event struct {
	ID       uuid.UUID
	EntityID uuid.UUID

	Effect Effect
	Data   interface{}

	CreatedAt time.Time
}

func initEvent(effect Effect, entityID uuid.UUID, data interface{}) Event {
	return Event{
		ID:       uuid.New(),
		Effect:   effect,
		EntityID: entityID,
		Data:     data,
	}
}
