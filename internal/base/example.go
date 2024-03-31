package base

import (
	"math"
	"thief/internal/model"

	"github.com/google/uuid"
)

const EntityTypeEnemy = "ENEMY"
const EntityTypeController = "CONTROLLER"

var (
	EnemyEventInit = Event{
		Effect: "INIT",
	}
	EnemyEventTargetAcquired = Event{
		Effect: "TARGET_ACQUIRE",
	}
	EnemyEventTargetRelease = Event{
		Effect: "TARGET_RELEASE",
	}
	EnemyEventForceTargetRelease = Event{
		Effect: "FORCE_TARGET_RELEASE",
	}
	EnemyEventMove = Event{
		Effect: "MOVE",
	}
)

var (
	ControllerEventInit = Event{
		Effect: "INIT",
	}
	ControllerEventEnemyTargetRelease = Event{
		Effect: "TARGET_RELEASE",
	}
)

func initEvent(effect Effect, entityID uuid.UUID, data interface{}) Event {
	return Event{
		ID:       uuid.New(),
		Effect:   effect,
		EntityID: entityID,
		Data:     data,
	}
}
func initEventWithSystemData(effect Effect, entityID uuid.UUID, data interface{}, systemData interface{}) Event {
	return Event{
		ID:         uuid.New(),
		Effect:     effect,
		EntityID:   entityID,
		Data:       data,
		SystemData: systemData,
	}
}

type targetData struct {
	targetType string
	id         uuid.UUID
}

type forceTargetReleaseData struct {
	inputID uuid.UUID
}

var EnemyPatrolStates = map[State]map[Effect]Gate{
	"IDLE": {
		EnemyEventTargetAcquired.Effect: {
			outputState: "TARGET_ACQUIRED",
			outputUnlockFunc: func(pm ProjectorManager, selfID uuid.UUID) interface{} {
				playerIDs := pm.GetEntityProjector("PLAYER").ListIdentifiers()
				playerPositions := make([]model.Position, len(playerIDs))

				for index, playerID := range playerIDs {
					playerPositions[index] = pm.GetEntityProjector("PLAYER").Project(playerID, "Position").(model.Position)
				}

				enemyPosition := pm.GetEntityProjector(EntityTypeEnemy).Project(selfID, "Position").(model.Position)

				if len(playerIDs) < 1 {
					return nil
				}

				isInRange := math.Sqrt(math.Pow(float64(enemyPosition.X-playerPositions[0].X), 2)) <= 5
				if !isInRange {
					return nil
				}

				return targetData{
					targetType: "PLAYER",
					id:         playerIDs[0],
				}
			},
		},
	},
	"TARGET_ACQUIRED": {
		EnemyEventForceTargetRelease.Effect: {
			outputState: "IDLE",
			outputUnlockFunc: func(pm ProjectorManager, selfID uuid.UUID) interface{} {
				targetReleaseData := pm.GetEntityProjector(EntityTypeEnemy).Project(selfID, "TargetReleaseLastInput").(forceTargetReleaseData)

				input := pm.GetEntityProjector(EntityTypeController).Project(selfID, "EnemyTargetReleaseInput").(ControllerInput)

				if input.ID != uuid.Nil && input.ID != targetReleaseData.inputID {
					return forceTargetReleaseData{
						inputID: input.ID,
					}
				}

				return nil
			},
		},
		EnemyEventMove.Effect: Gate{
			outputState: "TARGET_ACQUIRED",
			outputUnlockFunc: func(pm ProjectorManager, selfID uuid.UUID) interface{} {
				target := pm.GetEntityProjector(EntityTypeEnemy).Project(selfID, "Target").(targetData)

				if target.id == uuid.Nil {
					return nil
				}

				playerIDs := pm.GetEntityProjector("PLAYER").ListIdentifiers()
				playerPositions := make([]model.Position, len(playerIDs))

				for index, playerID := range playerIDs {
					playerPositions[index] = pm.GetEntityProjector("PLAYER").Project(playerID, "Position").(model.Position)
				}

				position := pm.GetEntityProjector(EntityTypeEnemy).Project(selfID, "Position").(model.Position)

				for index, playerID := range playerIDs {
					if playerID == target.id {
						playerPosition := playerPositions[index]
						if playerPosition.X == position.X {
							return nil
						}

						if playerPosition.X > position.X {
							return model.Position{X: position.X + 1}
						} else if playerPosition.X < position.X {
							return model.Position{X: position.X - 1}
						}
					}
				}

				return nil
			},
		},
	},
}

var EnemyPatrolStateMachine = stateMachine{
	entityType: "ENEMY",
	nodes:      EnemyPatrolStates,
}

type ControllerInput struct {
	ID    uuid.UUID
	Value interface{}
}

var ControllerStates = map[State]map[Effect]Gate{
	"ACTIVE": {
		ControllerEventEnemyTargetRelease.Effect: {
			outputState: "ACTIVE",
			outputUnlockFunc: func(pm ProjectorManager, selfID uuid.UUID) interface{} {
				return ControllerInput{
					ID: uuid.New(),
				}
			},
		},
	},
}

var controllerStateMachine = stateMachine{
	entityType:   "CONTROLLER",
	defaultState: "ACTIVE",
	nodes:        ControllerStates,
}
