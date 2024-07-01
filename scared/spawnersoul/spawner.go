package spawnersoul

import (
	"log"
	"thief/base/engine"
	"thief/scared"

	"github.com/google/uuid"
)

const (
	storePosition   = "StorePosition"
	storeEntityType = "StoreEntityType"
	storeTarget     = "StoreTarget"
	storeHealth     = "StoreHealth"
)

type spawnerShip struct {
	storePosition   *engine.Store
	storeEntityType *engine.Store
	storeTarget     *engine.Store
	storeHealth     *engine.Store
}

func NewSpawner(storePosition, storeEntityType, storeTarget, storeHealth *engine.Store) engine.Spawner {
	return spawnerShip{
		storePosition:   storePosition,
		storeEntityType: storeEntityType,
		storeTarget:     storeTarget,
		storeHealth:     storeHealth,
	}
}

func (s spawnerShip) GetInitEvents(eventData interface{}) map[string]engine.Event {
	position := eventData.(scared.Position)

	soulID := uuid.New()

	log.Println("spawned soul")

	return map[string]engine.Event{
		storePosition: {
			EntityID: soulID,
			Effect:   engine.EffectInit,
			Data:     position,
		},
		storeEntityType: {
			EntityID: soulID,
			Effect:   engine.EffectInit,
			Data:     scared.EntityTypeSoul,
		},
		storeTarget: {
			EntityID: soulID,
			Effect:   engine.EffectInit,
		},
		storeHealth: {
			EntityID: soulID,
			Effect:   engine.EffectInit,
			Data:     10,
		},
	}
}

func (s spawnerShip) GetStore(store string) *engine.Store {
	switch store {
	case storePosition:
		return s.storePosition
	case storeEntityType:
		return s.storeEntityType
	case storeTarget:
		return s.storeTarget
	case storeHealth:
		return s.storeHealth
	default:
		return nil
	}
}
