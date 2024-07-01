package spawnersoul

import (
	"thief/base/engine"
	"thief/scared"

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
			stateActive,
			func(selfID uuid.UUID) (interface{}, bool) {
				soulIDs := scared.ProjectorEntityType.ListIdentifiers(func(et scared.EntityType) bool {
					return scared.EntityTypeSoul == et
				})

				if len(soulIDs) > 0 {
					return nil, false
				}

				return scared.PositionSoulSpawn, true
			},
		),
	},
})
