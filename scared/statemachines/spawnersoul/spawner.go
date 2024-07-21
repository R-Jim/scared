package spawnersoul

import (
	"thief/base/engine"
	"thief/scared/model"
	"thief/scared/statemachines/entitytype"
	"thief/scared/statemachines/health"
	"thief/scared/statemachines/position"
	"thief/scared/statemachines/soul"
	"thief/scared/statemachines/target"

	"github.com/google/uuid"
)

type spawnerSoul struct {
	storePosition   *engine.Store
	storeEntityType *engine.Store
	storeTarget     *engine.Store
	storeHealth     *engine.Store
	storeSoul       *engine.Store
}

func NewSoul(storePosition, storeEntityType, storeTarget, storeHealth, storeSoul *engine.Store) engine.Receiver[model.SpawnSoulData] {
	return spawnerSoul{
		storePosition:   storePosition,
		storeEntityType: storeEntityType,
		storeTarget:     storeTarget,
		storeHealth:     storeHealth,
		storeSoul:       storeSoul,
	}
}

func (s spawnerSoul) GetEvents(entityID uuid.UUID, data model.SpawnSoulData) map[*engine.Store][]engine.Event {
	return map[*engine.Store][]engine.Event{
		s.storePosition: {
			position.EffectInit.NewEvent(data.ID, data.Position),
		},
		s.storeEntityType: {
			entitytype.EffectInit.NewEvent(data.ID, model.EntityTypeSoul),
		},
		s.storeTarget: {
			target.EffectInit.NewEvent(data.ID, nil),
		},
		s.storeHealth: {
			health.EffectInit.NewEvent(data.ID, 10),
		},
		s.storeSoul: {
			soul.EffectInit.NewEvent(data.ID, model.Soul{
				TemplateID: data.SoulTemplateID,
			}),
		},
	}
}
