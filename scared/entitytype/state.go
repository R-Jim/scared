package entitytype

import (
	"thief/base/engine"
	"thief/scared"

	"github.com/google/uuid"
)

const (
	stateActive    engine.State = "Active"
	stateDestroyed engine.State = engine.StateDestroyed
)

const (
	EffectDestroy engine.Effect = "Destroy"
)

const (
	MaxMoveRange = 1
)

var StateMachine = engine.NewStateMachine(stateActive, engine.Nodes{
	stateActive: {
		EffectDestroy: engine.NewGate(
			stateDestroyed,
			func(selfID uuid.UUID) (interface{}, bool) {
				return nil, scared.ProjectorHealth.IsDestroyed(selfID)
			},
		),
	},
})
