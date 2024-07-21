package engine

import (
	"log"

	"github.com/google/uuid"
)

type Projector[model any] interface {
	Project(identifier uuid.UUID) model
	IsDestroyed(identifier uuid.UUID) bool
	ListIdentifiers(filters ...func(model) bool) []uuid.UUID
	ListDeletedIdentifiers(filters ...func(model) bool) []uuid.UUID
}

type effectMapping[model any] struct {
	effects       []string
	aggregateFunc func(currentData model, nextEffectData any) model
}

func NewFieldEffectMapping[model, effectModel any](effects []Effect[effectModel], aggregateFunc func(currentData model, nextEffectData effectModel) model) effectMapping[model] {
	es := []string{}
	for _, effect := range effects {
		es = append(es, string(effect))
	}
	af := func(currentData model, nextEffectData any) model {
		return aggregateFunc(currentData, nextEffectData.(effectModel))
	}

	return effectMapping[model]{
		effects:       es,
		aggregateFunc: af,
	}
}

type StoreProjector[model any] struct {
	store          *Store
	effectMappings []effectMapping[model]
}

func NewStoreProjector[model any](store *Store, effectMappings ...effectMapping[model]) Projector[model] {
	return StoreProjector[model]{
		store:          store,
		effectMappings: effectMappings,
	}
}

func (p StoreProjector[model]) Project(identifier uuid.UUID) model {
	var result model

	events, err := p.store.GetEventsByEntityID(identifier)
	if err != nil {
		log.Fatalln(err)
	}

	if len(p.effectMappings) <= 0 {
		return result
	}

	for i := 0; i < len(events); i++ {
		event := events[i]

		for _, mapping := range p.effectMappings {
			for _, effect := range mapping.effects {
				if effect == event.Effect {
					result = mapping.aggregateFunc(result, event.Data)
					break
				}
			}
		}
	}

	return result
}

func (p StoreProjector[model]) ListIdentifiers(filters ...func(m model) bool) []uuid.UUID {
	eventsMap := p.store.GetEvents()

	identifiers := []uuid.UUID{}
	for entityID := range eventsMap {
		projection := p.Project(entityID)

		isMatchedFilters := true
		for _, filter := range filters {
			if !filter(projection) {
				isMatchedFilters = false
				break
			}
		}

		if isMatchedFilters {
			identifiers = append(identifiers, entityID)
		}
	}

	return identifiers
}

func (p StoreProjector[model]) ListDeletedIdentifiers(filters ...func(m model) bool) []uuid.UUID {
	eventsMap := p.store.destroyedEventsSet

	identifiers := []uuid.UUID{}
	for entityID := range eventsMap {
		projection := p.Project(entityID)

		isMatchedFilters := true
		for _, filter := range filters {
			if !filter(projection) {
				isMatchedFilters = false
				break
			}
		}

		if isMatchedFilters {
			identifiers = append(identifiers, entityID)
		}
	}

	return identifiers
}

func (p StoreProjector[model]) IsDestroyed(identifier uuid.UUID) bool {
	_, isExist := p.store.destroyedEventsSet[identifier]
	return isExist
}
