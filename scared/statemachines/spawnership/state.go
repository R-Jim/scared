package spawnership

import (
	"thief/base/engine"
	"thief/scared/model"

	"github.com/google/uuid"
)

const (
	stateActive  engine.State = "Active"
	stateSpawned engine.State = "Spawned"
)

const (
	EffectInit  engine.Effect[any]                 = "Init"
	EffectSpawn engine.Effect[model.SpawnShipData] = "Spawn"
)

var StateMachine = engine.NewStateMachine(EffectInit.ToState(stateActive), engine.Nodes{
	stateActive: {
		EffectSpawn.ToStateWhen(
			stateSpawned,
			func(selfID uuid.UUID) (model.SpawnShipData, bool) {
				return model.SpawnShipData{
					TemplateID: model.ShipDawnBreakID,
					Position:   model.PositionShipSpawn,
				}, true
			},
		),
	},
})
