package church

import (
	"thief/base/engine"
	"thief/scared/model"
	"thief/scared/statemachines/acolyte"
	"thief/scared/statemachines/entitytype"
	"thief/scared/statemachines/health"
	"thief/scared/statemachines/position"

	"github.com/google/uuid"
)

type spawnerChurch struct {
	storePosition   *engine.Store
	storeEntityType *engine.Store
	storeHealth     *engine.Store
	storeAcolyte    *engine.Store
}

func NewChurch(storePosition, storeEntityType, storeHealth, storeAcolyte *engine.Store) engine.Receiver[model.Position] {
	return spawnerChurch{
		storePosition:   storePosition,
		storeEntityType: storeEntityType,
		storeHealth:     storeHealth,
		storeAcolyte:    storeAcolyte,
	}
}

func (s spawnerChurch) GetEvents(entityID uuid.UUID, data model.Position) map[*engine.Store][]engine.Event {
	return map[*engine.Store][]engine.Event{
		s.storePosition: {
			position.EffectInit.NewEvent(entityID, data),
		},
		s.storeEntityType: {
			entitytype.EffectInit.NewEvent(entityID, model.EntityTypeChurch),
		},
		s.storeHealth: {
			health.EffectInit.NewEvent(entityID, 1),
		},
		s.storeAcolyte: {
			acolyte.EffectInit.NewEvent(entityID, 2),
		},
	}
}
