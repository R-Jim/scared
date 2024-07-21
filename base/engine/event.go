package engine

import (
	"time"

	"github.com/google/uuid"
)

type Effect[model any] string

func (e Effect[model]) ToState(outputState State) transition {
	return transition{
		effect:      string(e),
		outputState: outputState,
	}
}

func (e Effect[model]) ToStateWhen(outputState State, transitionFunc func(selfID uuid.UUID) (model, bool)) transition {
	tf := func(selfID uuid.UUID) (any, bool) {
		return transitionFunc(selfID)
	}

	return transition{
		effect:         string(e),
		outputState:    outputState,
		transitionFunc: tf,
	}
}

func (e Effect[model]) NewEvent(entityID uuid.UUID, data model) Event {
	return Event{
		ID:       uuid.New(),
		Effect:   string(e),
		EntityID: entityID,
		Data:     data,
	}
}

type Event struct {
	ID       uuid.UUID
	EntityID uuid.UUID

	Effect string
	Data   interface{}

	CreatedAt time.Time
}
