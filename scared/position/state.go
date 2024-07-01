package position

import (
	"thief/base/engine"
	"thief/scared"

	"github.com/google/uuid"
)

const (
	StateActive    engine.State = "Active"
	StateDestroyed engine.State = engine.StateDestroyed
)

const (
	EffectMove    engine.Effect = "Move"
	EffectDestroy engine.Effect = "Destroy"
)

const (
	MaxMoveRange = 1
)

var StateMachine = engine.NewStateMachine(StateActive, engine.Nodes{
	StateActive: {
		EffectMove: engine.NewGate(
			StateActive,
			func(selfID uuid.UUID) (interface{}, bool) {
				targetID := scared.ProjectorTarget.Project(selfID)
				if targetID == uuid.Nil {
					return nil, false
				}

				selfPosition := scared.ProjectorPosition.Project(selfID)
				targetPosition := scared.ProjectorPosition.Project(targetID)

				distance := selfPosition.DistanceOf(targetPosition)

				if distance < MaxMoveRange {
					return scared.Position{
						X: targetPosition.X - selfPosition.X,
						Y: targetPosition.Y - selfPosition.Y,
					}, true
				}

				steps := int(distance / MaxMoveRange)

				return scared.Position{
					X: (targetPosition.X - selfPosition.X) / steps,
					Y: (targetPosition.Y - selfPosition.Y) / steps,
				}, true
			},
		),
		EffectDestroy: engine.NewGate(
			StateDestroyed,
			func(selfID uuid.UUID) (interface{}, bool) {
				return nil, scared.ProjectorHealth.IsDestroyed(selfID)
			},
		),
	},
})
