package engine

import (
	"log"

	"github.com/google/uuid"
)

type Projector interface {
	Project(identifier uuid.UUID, field Field) interface{}
	ListIdentifiers() []uuid.UUID
}

type fieldEffectMapping[model any] struct {
	field         Field
	effects       []Effect
	aggregateFunc func(a, b model) model
	getFunc       func(m model) interface{}
}

func NewFieldEffectMapping[model any](field Field, effects []Effect, aggregateFunc func(a, b model) model, getFunc func(m model) interface{}) fieldEffectMapping[model] {
	return fieldEffectMapping[model]{
		field:         field,
		effects:       effects,
		aggregateFunc: aggregateFunc,
		getFunc:       getFunc,
	}
}

type storeProjector[model any] struct {
	store         *Store
	fieldMappings []fieldEffectMapping[model]
}

func NewStoreProjector[model any](store *Store, fieldMappings ...fieldEffectMapping[model]) Projector {
	return storeProjector[model]{
		store:         store,
		fieldMappings: fieldMappings,
	}
}

func (p storeProjector[model]) Project(identifier uuid.UUID, field Field) interface{} {
	events, err := p.store.GetEventsByEntityID(identifier)
	if err != nil {
		log.Fatalln(err)
	}

	var mapping *fieldEffectMapping[model]

	for _, fm := range p.fieldMappings {
		if fm.field == field {
			mapping = &fm
			break
		}
	}
	if mapping == nil {
		return nil
	}

	var result model

	for i := 0; i < len(events); i++ {
		var eventData model
		event := events[i]

		for _, effect := range mapping.effects {
			if effect == event.Effect {
				if event.Data != nil {
					eventData = event.Data.(model)
				}

				result = mapping.aggregateFunc(result, eventData)
				break
			}
		}
	}

	return mapping.getFunc(result)
}

func (p storeProjector[model]) ListIdentifiers() []uuid.UUID {
	eventsMap := p.store.GetEvents()

	identifiers := []uuid.UUID{}
	for entityID := range eventsMap {
		identifiers = append(identifiers, entityID)
	}

	return identifiers
}

type ProjectorTypeMapping map[EntityType]Projector

type ProjectorManager struct {
	mapping ProjectorTypeMapping
}

func NewProjectorManager(mapping ProjectorTypeMapping) ProjectorManager {
	return ProjectorManager{
		mapping: mapping,
	}
}

func (m ProjectorManager) Get(entityType EntityType) Projector {
	return m.mapping[entityType]
}
