package state

import (
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	pkgerrors "github.com/pkg/errors"
)

type Event struct {
	ID       uuid.UUID
	EntityID uuid.UUID

	Effect string
	Data   interface{}

	FromState string
	ToState   string

	CreatedAt time.Time
}

func ParseData[T interface{}](e Event) (T, error) {
	data, ok := e.Data.(T)
	if !ok {
		return data, pkgerrors.WithStack(fmt.Errorf("failed to parse data for effect[%s]: %s", e.Effect, reflect.TypeOf(data)))
	}
	return data, nil
}
