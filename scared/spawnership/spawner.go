package spawnership

import (
	"log"
	"thief/base/engine"
	"thief/scared"

	"github.com/google/uuid"
)

const (
	storePosition   = "StorePosition"
	storeEntityType = "StoreEntityType"
	storeHealth     = "StoreHealth"
	storeShipGuard  = "StoreShipGuard"
)

type spawnerShip struct {
	storePosition   *engine.Store
	storeEntityType *engine.Store
	storeHealth     *engine.Store
	storeShipGuard  *engine.Store
}

func NewSpawner(storePosition, storeEntityType, storeHealth, storeShipGuard *engine.Store) engine.Spawner {
	return spawnerShip{
		storePosition:   storePosition,
		storeEntityType: storeEntityType,
		storeHealth:     storeHealth,
		storeShipGuard:  storeShipGuard,
	}
}

func (s spawnerShip) GetInitEvents(eventData interface{}) map[string]engine.Event {
	inputData := eventData.(engine.InputData)
	templateID, err := uuid.Parse(inputData.Key)
	if err != nil {
		log.Fatalln(err)
	}

	shipID := uuid.New()
	template := scared.ShipTemplates[templateID]

	log.Println("spawned ship")

	return map[string]engine.Event{
		storePosition: {
			EntityID: shipID,
			Effect:   engine.EffectInit,
			Data:     scared.PositionShipSpawn,
		},
		storeEntityType: {
			EntityID: shipID,
			Effect:   engine.EffectInit,
			Data:     scared.EntityTypeShip,
		},
		storeHealth: {
			EntityID: shipID,
			Effect:   engine.EffectInit,
			Data:     100,
		},
		storeShipGuard: {
			EntityID: shipID,
			Effect:   engine.EffectInit,
			Data: scared.ShipGuard{
				Quantity: template.GuardQuantity,
			},
		},
	}
}

func (s spawnerShip) GetStore(store string) *engine.Store {
	switch store {
	case storePosition:
		return s.storePosition
	case storeEntityType:
		return s.storeEntityType
	case storeHealth:
		return s.storeHealth
	case storeShipGuard:
		return s.storeShipGuard
	default:
		return nil
	}
}
