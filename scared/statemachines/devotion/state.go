package devotion

import (
	"thief/base/engine"
	"thief/scared/projectors"

	"github.com/google/uuid"
)

const (
	EffectInit    engine.Effect[int] = "Init"
	EffectAdd     engine.Effect[int] = "Add"
	EffectConsume engine.Effect[int] = "Consume"
)

const (
	stateGeneratorActive   engine.State = "Active"
	stateGeneratorCoolDown engine.State = "CoolDown"
)

const (
	EffectGeneratorInit     engine.Effect[int] = "Init"
	EffectGeneratorRun      engine.Effect[int] = "Run"
	EffectGeneratorCoolDown engine.Effect[int] = "CoolDown"
	EffectGeneratorReActive engine.Effect[any] = "ReActive"
)

var StateMachineDevotionGenerator = engine.NewStateMachine(EffectGeneratorInit.ToState(stateGeneratorActive), engine.Nodes{
	stateGeneratorActive: {
		EffectGeneratorRun.ToStateWhen(
			stateGeneratorCoolDown,
			func(selfID uuid.UUID) (int, bool) {
				return projectors.ProjectorDevotionGeneratorCoolDown.Project(selfID), true
			},
		),
	},
	stateGeneratorCoolDown: {
		EffectGeneratorCoolDown.ToStateWhen(
			stateGeneratorCoolDown,
			func(selfID uuid.UUID) (int, bool) {
				coolDownCount := projectors.ProjectorDevotionGeneratorCoolDownCount.Project(selfID)
				if coolDownCount <= 0 {
					return 0, false
				}

				return 1, true
			},
		),
		EffectGeneratorReActive.ToStateWhen(
			stateGeneratorActive,
			func(selfID uuid.UUID) (interface{}, bool) {
				coolDownCount := projectors.ProjectorDevotionGeneratorCoolDownCount.Project(selfID)
				if coolDownCount > 0 {
					return nil, false
				}

				return nil, true
			},
		),
	},
})
