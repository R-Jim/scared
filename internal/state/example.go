package state

import (
	"log"
	"math"
	"thief/internal/model"

	"github.com/google/uuid"
)

const EntityTypeEnemy = "ENEMY"

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
	EnemyEventMove = Event{
		Effect: "MOVE",
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

var EnemyPatrolStates = []Node{
	{
		State: "IDLE",
		ActiveGates: []Gate{
			{
				nextState:           "TARGET_ACQUIRED",
				unlockByEventEffect: EnemyEventTargetAcquired.Effect,
				eventProducerFunc: func(selfID uuid.UUID, pm ProjectorManager, _ interface{}) interface{} {
					playerPositions, playerIDs, err := FieldValues[model.Position]("PLAYER", "Position", pm)
					if err != nil {
						log.Fatalln(err)
					}

					enemyPosition, err := FieldValue[model.Position](selfID, EntityTypeEnemy, "Position", pm)
					if err != nil {
						log.Fatalln(err)
					}

					if len(playerIDs) < 1 {
						return Event{}
					}

					isInRange := math.Sqrt(math.Pow(float64(enemyPosition.X-playerPositions[0].X), 2)) <= 5
					if !isInRange {
						return Event{}
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
		ActiveGates: []Gate{
			{
				nextState:           "IDLE",
				unlockByEventEffect: EnemyEventTargetRelease.Effect,
				eventProducerFunc: func(selfID uuid.UUID, pm ProjectorManager, _ interface{}) interface{} {
					target, err := FieldValue[targetData](selfID, EntityTypeEnemy, "Target", pm)
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
					target, err := FieldValue[targetData](selfID, EntityTypeEnemy, "Target", pm)
					if err != nil {
						log.Fatalln(err)
					}

					if target.id == uuid.Nil {
						log.Println("no target to move")
						return Event{}
					}

					playerPositions, playerIDs, err := FieldValues[model.Position]("PLAYER", "Position", pm)
					if err != nil {
						log.Fatalln(err)
					}

					position, err := FieldValue[model.Position](selfID, EntityTypeEnemy, "Position", pm)
					if err != nil {
						log.Fatalln(err)
					}

					for index, playerID := range playerIDs {
						if playerID == target.id {
							playerPosition := playerPositions[index]
							if playerPosition.X == position.X {
								return Event{}
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
		PassiveGates: []Gate{
			{
				nextState:           "IDLE",
				unlockByEventEffect: EnemyEventTargetRelease.Effect,
				eventProducerFunc: func(selfID uuid.UUID, pm ProjectorManager, _ interface{}) interface{} {
					return targetData{}
				},
			},
		},
	},
}

var EnemyPatrolStateMachine = StateMachine{
	entityType: "ENEMY",
	nodes:      EnemyPatrolStates,
}
