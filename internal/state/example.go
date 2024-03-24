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

func InitEvent(effect string, entityID uuid.UUID, data interface{}) Event {
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

var EnemyPatrolStates = []Node{
	{
		State: "IDLE",
		Gates: []Gate{
			{
				nextState:           "TARGET_ACQUIRED",
				unlockByEventEffect: EnemyEventTargetAcquired.Effect,
				eventProducerFunc: func(selfID uuid.UUID, pm ProjectorManager, _ interface{}) interface{} {
					playerPositions, playerIDs, err := FieldValues[model.Position](pm, "PLAYER", "Position")
					if err != nil {
						log.Fatalln(err)
					}

					enemyPosition, err := FieldValue[model.Position](pm, selfID, EntityTypeEnemy, "Position")
					if err != nil {
						log.Fatalln(err)
					}

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
	},
	{
		State: "TARGET_ACQUIRED",
		Gates: []Gate{
			{
				nextState:           "IDLE",
				unlockByEventEffect: EnemyEventForceTargetRelease.Effect,
				eventProducerFunc: func(selfID uuid.UUID, pm ProjectorManager, _ interface{}) interface{} {
					targetReleaseData, err := FieldValue[forceTargetReleaseData](pm, selfID, EntityTypeEnemy, "TargetReleaseLastInput")
					if err != nil {
						log.Fatalln(err)
					}

					input, err := FieldValue[ControllerInput](pm, selfID, EntityTypeController, "EnemyTargetReleaseInput")
					if err != nil {
						log.Fatalln(err)
					}

					if input.ID != uuid.Nil && input.ID != targetReleaseData.inputID {
						return forceTargetReleaseData{
							inputID: input.ID,
						}
					}

					return nil
				},
			},
			{
				nextState:           "IDLE",
				unlockByEventEffect: EnemyEventTargetRelease.Effect,
				eventProducerFunc: func(selfID uuid.UUID, pm ProjectorManager, _ interface{}) interface{} {
					target, err := FieldValue[targetData](pm, selfID, EntityTypeEnemy, "Target")
					if err != nil {
						log.Fatalln(err)
					}

					if target.id != uuid.Nil {
						return nil
					}

					return targetData{}
				},
			},
			{
				nextState:           "TARGET_ACQUIRED",
				unlockByEventEffect: EnemyEventMove.Effect,
				eventProducerFunc: func(selfID uuid.UUID, pm ProjectorManager, _ interface{}) interface{} {
					target, err := FieldValue[targetData](pm, selfID, EntityTypeEnemy, "Target")
					if err != nil {
						log.Fatalln(err)
					}

					if target.id == uuid.Nil {
						log.Println("no target to move")
						return nil
					}

					playerPositions, playerIDs, err := FieldValues[model.Position](pm, "PLAYER", "Position")
					if err != nil {
						log.Fatalln(err)
					}

					position, err := FieldValue[model.Position](pm, selfID, EntityTypeEnemy, "Position")
					if err != nil {
						log.Fatalln(err)
					}

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
	},
}

var EnemyPatrolStateMachine = StateMachine{
	entityType: "ENEMY",
	nodes:      EnemyPatrolStates,
}

type ControllerInput struct {
	ID    uuid.UUID
	Value interface{}
}

var ControllerStates = []Node{
	{
		State: "ACTIVE",
		Gates: []Gate{
			{
				nextState:           "ACTIVE",
				unlockByEventEffect: ControllerEventEnemyTargetRelease.Effect,
				eventProducerFunc: func(selfID uuid.UUID, pm ProjectorManager, _ interface{}) interface{} {
					return ControllerInput{
						ID: uuid.New(),
					}
				},
			},
		},
	},
}

var ControllerStateMachine = StateMachine{
	entityType: "CONTROLLER",
	nodes:      ControllerStates,
}
