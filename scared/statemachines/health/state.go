package health

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
	EffectInit        engine.Effect[int] = "Init"
	EffectHit         engine.Effect[int] = "HitHealth"
	EffectHitDevotion engine.Effect[int] = "HitDevotion"
	EffectDestroy     engine.Effect[any] = "Destroy"
)

const (
	MaxMoveRange = 1
)

var StateMachine = engine.NewStateMachine(EffectInit.ToState(stateActive), engine.Nodes{
	stateActive: {
		EffectHit.ToStateWhen(
			stateActive,
			func(selfID uuid.UUID) (int, bool) {
				return 0, false // consumer only
			},
		),
		EffectDestroy.ToStateWhen(
			stateDestroyed,
			func(selfID uuid.UUID) (interface{}, bool) {
				if projectors.ProjectorEntityType.Project(selfID) == model.EntityTypeChurch {
					return nil, projectors.ProjectorAcolyte.Project(selfID) <= 0
				}

				return nil, projectors.ProjectorHealth.Project(selfID) <= 0
			},
		),
	},
})
