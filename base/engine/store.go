package engine

import (
	"log"

	"github.com/google/uuid"
)

type Store struct {
	name    string
	counter *int // used to check if a store has new events

	eventsSet          map[uuid.UUID][]Event
	destroyedEventsSet map[uuid.UUID][]Event

	hooks []Hook
}

func NewStore(name string) *Store {
	defaultCounter := 0
	return &Store{
		name:               name,
		eventsSet:          make(map[uuid.UUID][]Event),
		destroyedEventsSet: make(map[uuid.UUID][]Event),
		counter:            &defaultCounter,
	}
}

func (i Store) GetEventsByEntityID(id uuid.UUID) ([]Event, error) {
	return i.eventsSet[id], nil
}

func (i *Store) AppendEvent(e Event) error {
	if e.ID == uuid.Nil {
		log.Fatalf("fail to append, no event ID, event: %v", e)
	}

	events := i.eventsSet[e.EntityID]
	if events == nil {
		events = []Event{}
	}

	i.eventsSet[e.EntityID] = append(events, e)
	(*i.counter)++

	for _, hook := range i.hooks {
		hook(e)
	}

	return nil
}

func (i Store) GetEvents() map[uuid.UUID][]Event {
	events := i.eventsSet // Shallow clone
	return events
}

func (i *Store) GetCounter() int {
	return *i.counter
}

func (i *Store) destroySet(id uuid.UUID) {
	i.destroyedEventsSet[id] = i.eventsSet[id]
	delete(i.eventsSet, id)
}

func (i *Store) IsDestroyed(id uuid.UUID) bool {
	_, isExist := i.destroyedEventsSet[id]
	return isExist
}

func (i *Store) AddHook(hook Hook) {
	i.hooks = append(i.hooks, hook)
}
