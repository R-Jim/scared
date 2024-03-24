package state

import (
	"log"
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

type targetData struct {
	targetType string
	id         uuid.UUID
}

type forceTargetReleaseData struct {
	inputID uuid.UUID
}

var EnemyPatrolStates = map[State]map[Effect]gate{
	"IDLE": {
		EnemyEventTargetAcquired.Effect: {
			nextState: "TARGET_ACQUIRED",
			outputProducerFunc: func(selfID uuid.UUID, pm ProjectorManager) interface{} {
				playerPositions, playerIDs := FieldValues[model.Position](pm, "PLAYER", "Position")

				enemyPosition := FieldValue[model.Position](pm, selfID, EntityTypeEnemy, "Position")

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
			nextState: "IDLE",
			outputProducerFunc: func(selfID uuid.UUID, pm ProjectorManager) interface{} {
				targetReleaseData := FieldValue[forceTargetReleaseData](pm, selfID, EntityTypeEnemy, "TargetReleaseLastInput")

				input := FieldValue[ControllerInput](pm, selfID, EntityTypeController, "EnemyTargetReleaseInput")

				if input.ID != uuid.Nil && input.ID != targetReleaseData.inputID {
					return forceTargetReleaseData{
						inputID: input.ID,
					}
				}

				return nil
			},
		},
		EnemyEventTargetRelease.Effect: {
			nextState: "IDLE",
			outputProducerFunc: func(selfID uuid.UUID, pm ProjectorManager) interface{} {
				target := FieldValue[targetData](pm, selfID, EntityTypeEnemy, "Target")

				if target.id != uuid.Nil {
					return nil
				}

				return targetData{}
			},
		},
		EnemyEventMove.Effect: gate{
			nextState: "TARGET_ACQUIRED",
			outputProducerFunc: func(selfID uuid.UUID, pm ProjectorManager) interface{} {
				target := FieldValue[targetData](pm, selfID, EntityTypeEnemy, "Target")

				if target.id == uuid.Nil {
					log.Println("no target to move")
					return nil
				}

				playerPositions, playerIDs := FieldValues[model.Position](pm, "PLAYER", "Position")

				position := FieldValue[model.Position](pm, selfID, EntityTypeEnemy, "Position")

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

				log.Println("no match target to move")
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

var ControllerStates = map[State]map[Effect]gate{
	"ACTIVE": {
		ControllerEventEnemyTargetRelease.Effect: {
			nextState: "ACTIVE",
			outputProducerFunc: func(selfID uuid.UUID, pm ProjectorManager) interface{} {
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
