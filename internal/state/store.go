package state

import (
	"github.com/google/uuid"
)

type Store struct {
	counter *int

	eventsSet map[uuid.UUID][]Event
}

func NewStore() Store {
	defaultCounter := 0
	return Store{
		eventsSet: make(map[uuid.UUID][]Event),
		counter:   &defaultCounter,
	}
}

func (i Store) GetEventsByEntityID(id uuid.UUID) ([]Event, error) {
	return i.eventsSet[id], nil
}

func (i *Store) AppendEvent(e Event) error {
	events := i.eventsSet[e.EntityID]
	if events == nil {
		events = []Event{}
	}

	i.eventsSet[e.EntityID] = append(events, e)
	(*i.counter)++
	return nil
}

func (i Store) GetEvents() map[uuid.UUID][]Event {
	events := i.eventsSet // Shallow clone
	return events
}

func (i *Store) GetCounter() int {
	return *i.counter
}
