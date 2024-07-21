package spawnership

import (
	"thief/base/engine"
	"thief/scared/model"
	"thief/scared/statemachines/acolyte"
	"thief/scared/statemachines/entitytype"
	"thief/scared/statemachines/health"
	"thief/scared/statemachines/position"
	"thief/scared/statemachines/ship"

	"github.com/google/uuid"
)

type spawnerShip struct {
	storePosition          *engine.Store
	storeEntityType        *engine.Store
	storeHealth            *engine.Store
	storeShipGuard         *engine.Store
	storeShipBlessingAltar *engine.Store
	storeAcolyte           *engine.Store
}

func NewShip(storePosition, storeEntityType, storeHealth, storeShipGuard, storeShipBlessingAltar, storeAcolyte *engine.Store) engine.Receiver[model.SpawnShipData] {
	return spawnerShip{
		storePosition:          storePosition,
		storeEntityType:        storeEntityType,
		storeHealth:            storeHealth,
		storeShipGuard:         storeShipGuard,
		storeShipBlessingAltar: storeShipBlessingAltar,
		storeAcolyte:           storeAcolyte,
	}
}

func (s spawnerShip) GetEvents(entityID uuid.UUID, data model.SpawnShipData) map[*engine.Store][]engine.Event {
	shipID := uuid.New()
	template := model.TemplateShips[data.TemplateID]

	return map[*engine.Store][]engine.Event{
		s.storePosition: {
			position.EffectInit.NewEvent(shipID, data.Position),
		},
		s.storeEntityType: {
			entitytype.EffectInit.NewEvent(shipID, model.EntityTypeShip),
		},
		s.storeHealth: {
			health.EffectInit.NewEvent(shipID, 20),
		},
		s.storeShipGuard: {
			ship.EffectInit.NewEvent(shipID, model.ShipGuard{
				Quantity: template.GuardQuantity,
			}),
		},
		s.storeShipBlessingAltar: {
			ship.EffectBlessingAltarInit.NewEvent(uuid.New(), shipID),
		},
		s.storeAcolyte: {
			acolyte.EffectInit.NewEvent(shipID, template.AcolyteQuantity),
		},
	}
}
