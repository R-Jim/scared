package ship

import (
	"thief/base/engine"

	"github.com/google/uuid"
)

const (
	stateUnarmed   engine.State = "Unarmed"
	stateArmed     engine.State = "Armed"
	stateDestroyed engine.State = "Destroyed"
)

const (
	EffectArm     engine.Effect = "Arm"
	EffectDisarm  engine.Effect = "Disarm"
	EffectDestroy engine.Effect = "Destroy"
)

var StateMachine = engine.NewStateMachine(stateUnarmed, engine.Nodes{
	stateUnarmed: {
		EffectArm: engine.NewGate(
			stateArmed,
			func(selfID uuid.UUID) (interface{}, bool) {
				return selfID, true
			},
		),
	},
})
