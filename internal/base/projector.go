package base

import (
	"log"

	"github.com/google/uuid"
)

type Projector interface {
	Project(identifier uuid.UUID, field string) interface{}
	ListIdentifiers() []uuid.UUID
}

type storeProjector struct {
	store        *Store
	fieldMapping map[string]map[Effect]func(eventData interface{}) interface{}
}

func NewStoreProjector(store *Store, fieldMapping map[string]func() interface{}) Projector {
	return storeProjector{
		store: store,
	}
}

func (p storeProjector) Project(identifier uuid.UUID, field string) interface{} {
	events, err := p.store.GetEventsByEntityID(identifier)
	if err != nil {
		log.Fatalln(err)
	}

	for i := len(events) - 1; i >= 0; i-- {
		supportedEffects, ok := p.fieldMapping[field]
		if !ok {
			continue
		}

		event := events[i]

		mappingFunc, ok := supportedEffects[event.Effect]
		if !ok {
			continue
		}

		return mappingFunc(event.Data)
	}
	return nil
}

func (p storeProjector) ListIdentifiers() []uuid.UUID {
	eventsMap := p.store.GetEvents()

	identifiers := []uuid.UUID{}
	for entityID := range eventsMap {
		identifiers = append(identifiers, entityID)
	}

	return identifiers
}

type ProjectorManager struct {
	entityProjectorMap map[string]Projector
}

func NewProjectorManager(projectorMapping map[string]Projector) ProjectorManager {
	return ProjectorManager{
		entityProjectorMap: projectorMapping,
	}
}

func (m ProjectorManager) GetEntityProjector(entityType string) Projector {
	return m.entityProjectorMap[entityType]
}
