package knight

import (
	"thief/base/engine"
	"thief/scared/model"
	"thief/scared/projectors"

	"github.com/google/uuid"
)

const (
	stateActive    engine.State = "Active"
	stateDestroyed engine.State = engine.StateDestroyed
)

const (
	EffectInit    engine.Effect[model.SpawnKnightData] = "Init"
	EffectDestroy engine.Effect[any]                   = "Destroy"
)

var StateMachine = engine.NewStateMachine(EffectInit.ToState(stateActive), engine.Nodes{
	stateActive: {
		EffectDestroy.ToStateWhen(
			stateDestroyed,
			func(selfID uuid.UUID) (any, bool) {
				return nil, projectors.ProjectorEntityType.IsDestroyed(selfID)
			},
		),
	},
})
