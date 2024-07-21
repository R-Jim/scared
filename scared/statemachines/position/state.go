package position

import (
	"math"
	"thief/base/engine"
	"thief/scared/model"
	"thief/scared/projectors"

	"github.com/google/uuid"
)

const (
	StateActive    engine.State = "Active"
	StateDestroyed engine.State = engine.StateDestroyed
)

const (
	EffectInit    engine.Effect[model.Position] = "Init"
	EffectMove    engine.Effect[model.Position] = "Move"
	EffectDestroy engine.Effect[any]            = "Destroy"
)

const (
	MaxMoveRange = 1
)

var StateMachine = engine.NewStateMachine(EffectInit.ToState(StateActive), engine.Nodes{
	StateActive: {
		EffectMove.ToStateWhen(
			StateActive,
			func(selfID uuid.UUID) (model.Position, bool) {
				var targetPosition *model.Position

				targetID := projectors.ProjectorTarget.Project(selfID)
				if targetID != uuid.Nil {
					p := projectors.ProjectorPosition.Project(targetID)
					targetPosition = &p
				} else {
					waypointIDs := projectors.ProjectorWaypoint.ListIdentifiers(func(w model.Waypoint) bool {
						return w.OwnerID == selfID
					})

					if len(waypointIDs) > 0 {
						waypoint := projectors.ProjectorWaypoint.Project(waypointIDs[0])
						targetPosition = waypoint.Position
					}
				}

				if targetPosition == nil {
					return model.Position{}, false
				}

				selfPosition := projectors.ProjectorPosition.Project(selfID)

				distance := selfPosition.DistanceOf(*targetPosition)

				if distance < MaxMoveRange {
					return model.Position{
						X: targetPosition.X - selfPosition.X,
						Y: targetPosition.Y - selfPosition.Y,
					}, true
				}

				steps := int(math.Ceil(distance) / MaxMoveRange)

				var nextX, nextY int

				if steps%2 == 0 {
					nextX = 1
					nextY = 1
				} else if steps%3 == 0 || steps%5 == 0 {
					nextX = 1
				} else { // mod 1, 7
					nextY = 1
				}

				if selfPosition.X > targetPosition.X {
					nextX = -nextX
				}
				if selfPosition.Y > targetPosition.Y {
					nextY = -nextY
				}

				return model.Position{
					X: nextX,
					Y: nextY,
				}, true
			},
		),
		EffectDestroy.ToStateWhen(
			StateDestroyed,
			func(selfID uuid.UUID) (interface{}, bool) {
				return nil, projectors.ProjectorHealth.IsDestroyed(selfID)
			},
		),
	},
})
