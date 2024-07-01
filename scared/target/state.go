package target

import (
	"thief/base/engine"
	"thief/scared"

	"github.com/google/uuid"
)

const (
	stateIdle      engine.State = "Idle"
	stateHasTarget engine.State = "HasTarget"
	stateDestroyed engine.State = engine.StateDestroyed
)

const (
	EffectSelectTarget engine.Effect = "SelectTarget"
	EffectDestroy      engine.Effect = "Destroy"
)

var StateMachine = engine.NewStateMachine(stateIdle, engine.Nodes{
	stateIdle: {
		EffectSelectTarget: engine.NewGate(
			stateHasTarget,
			func(selfID uuid.UUID) (interface{}, bool) {
				selfPosition := scared.ProjectorPosition.Project(selfID)
				entityType := scared.ProjectorEntityType.Project(selfID)

				targetTypes := scared.HunterTargetMapping[entityType]

				targetIDs := scared.ProjectorEntityType.ListIdentifiers(
					func(t scared.EntityType) bool {
						for _, targetType := range targetTypes {
							if t == targetType {
								return true
							}
						}
						return false
					},
				)

				var nearestTargetID uuid.UUID
				var nearestTargetDistance float64

				for _, shipID := range targetIDs {
					shipPosition := scared.ProjectorPosition.Project(shipID)

					distance := selfPosition.DistanceOf(shipPosition)
					if nearestTargetID == uuid.Nil || distance < nearestTargetDistance {
						nearestTargetID = shipID
						nearestTargetDistance = distance
					}
				}

				if nearestTargetID == uuid.Nil {
					return nil, false
				}

				return nearestTargetID, true
			},
		),
		EffectDestroy: engine.NewGate(
			stateDestroyed,
			func(selfID uuid.UUID) (interface{}, bool) {
				return nil, scared.ProjectorHealth.IsDestroyed(selfID)
			},
		),
	},
	stateHasTarget: {
		EffectDestroy: engine.NewGate(
			stateDestroyed,
			func(selfID uuid.UUID) (interface{}, bool) {
				return nil, scared.ProjectorHealth.IsDestroyed(selfID)
			},
		),
	},
})
