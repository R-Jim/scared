package target

import (
	"thief/base/engine"
	"thief/scared/model"
	"thief/scared/projectors"

	"github.com/google/uuid"
)

const (
	stateIdle      engine.State = "Idle"
	stateHasTarget engine.State = "HasTarget"
	stateDestroyed engine.State = engine.StateDestroyed
)

const (
	EffectInit          engine.Effect[any]       = "Init"
	EffectSelectTarget  engine.Effect[uuid.UUID] = "SelectTarget"
	EffectReleaseTarget engine.Effect[uuid.UUID] = "ReleaseTarget"
	EffectDestroy       engine.Effect[any]       = "Destroy"
)

var StateMachine = engine.NewStateMachine(EffectInit.ToState(stateIdle), engine.Nodes{
	stateIdle: {
		EffectSelectTarget.ToStateWhen(
			stateHasTarget,
			func(selfID uuid.UUID) (uuid.UUID, bool) {
				selfPosition := projectors.ProjectorPosition.Project(selfID)
				entityType := projectors.ProjectorEntityType.Project(selfID)

				targetTypes := model.EntityTypeMoveTargetMapping[entityType]

				targetIDs := projectors.ProjectorEntityType.ListIdentifiers(
					func(t model.EntityType) bool {
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
					shipPosition := projectors.ProjectorPosition.Project(shipID)

					distance := selfPosition.DistanceOf(shipPosition)
					if nearestTargetID == uuid.Nil || distance < nearestTargetDistance {
						nearestTargetID = shipID
						nearestTargetDistance = distance
					}
				}

				if nearestTargetID == uuid.Nil {
					return uuid.Nil, false
				}

				return nearestTargetID, true
			},
		),
		EffectDestroy.ToStateWhen(
			stateDestroyed,
			func(selfID uuid.UUID) (interface{}, bool) {
				return nil, projectors.ProjectorHealth.IsDestroyed(selfID)
			},
		),
	},
	stateHasTarget: {
		EffectReleaseTarget.ToStateWhen(
			stateIdle,
			func(selfID uuid.UUID) (uuid.UUID, bool) {
				targetID := projectors.ProjectorTarget.Project(selfID)
				return uuid.Nil, projectors.ProjectorEntityType.IsDestroyed(targetID)
			},
		),
		EffectDestroy.ToStateWhen(
			stateDestroyed,
			func(selfID uuid.UUID) (interface{}, bool) {
				return nil, projectors.ProjectorHealth.IsDestroyed(selfID)
			},
		),
	},
})
