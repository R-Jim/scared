package spawnership

import (
	"thief/base/engine"

	"github.com/google/uuid"
)

const (
	stateActive  engine.State = "Active"
	stateSpawned engine.State = "Spawned"
)

const (
	EffectSpawn engine.Effect = "Spawn"
)

var StateMachine = engine.NewStateMachine(stateActive, engine.Nodes{
	stateActive: {
		EffectSpawn: engine.NewGate(
			stateSpawned,
			func(selfID uuid.UUID) (interface{}, bool) {
				return nil, true
			},
		),
	},
})
