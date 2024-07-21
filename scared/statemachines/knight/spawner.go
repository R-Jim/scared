package knight

import (
	"thief/base/engine"
	"thief/scared/model"
	"thief/scared/projectors"
	"thief/scared/statemachines/entitytype"
	"thief/scared/statemachines/health"
	"thief/scared/statemachines/position"

	"github.com/google/uuid"
)

type spawnerKnight struct {
	storePosition   *engine.Store
	storeEntityType *engine.Store
	storeTarget     *engine.Store
	storeHealth     *engine.Store
	storeKnight     *engine.Store
}

func NewKnight(storePosition, storeEntityType, storeTarget, storeHealth, storeKnight *engine.Store) engine.Receiver[model.SpawnKnightData] {
	return spawnerKnight{
		storePosition:   storePosition,
		storeEntityType: storeEntityType,
		storeTarget:     storeTarget,
		storeHealth:     storeHealth,
		storeKnight:     storeKnight,
	}
}

func (s spawnerKnight) GetEvents(entityID uuid.UUID, data model.SpawnKnightData) map[*engine.Store][]engine.Event {
	KnightID := uuid.New()

	return map[*engine.Store][]engine.Event{
		s.storePosition: {
			position.EffectInit.NewEvent(KnightID, data.Position),
		},
		s.storeEntityType: {
			entitytype.EffectInit.NewEvent(KnightID, model.EntityTypeKnight),
		},
		s.storeHealth: {
			health.EffectInit.NewEvent(KnightID, 10),
		},
		s.storeKnight: {
			EffectInit.NewEvent(KnightID, model.SpawnKnightData{
				Position: projectors.ProjectorPosition.Project(entityID),
			}),
		},
	}
}
