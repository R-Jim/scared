package engine

// TODO: update example
// var (
// 	EnemyEventInit = Event{
// 		Effect: "INIT",
// 	}
// 	EnemyEventTargetAcquired = Event{
// 		Effect: "TARGET_ACQUIRE",
// 	}
// 	EnemyEventTargetRelease = Event{
// 		Effect: "TARGET_RELEASE",
// 	}
// 	EnemyEventForceTargetRelease = Event{
// 		Effect: "FORCE_TARGET_RELEASE",
// 	}
// 	EnemyEventMove = Event{
// 		Effect: "MOVE",
// 	}
// )

// var (
// 	ControllerEventInit = Event{
// 		Effect: "INIT",
// 	}
// 	ControllerEventEnemyTargetRelease = Event{
// 		Effect: "TARGET_RELEASE",
// 	}
// )

// type targetData struct {
// 	targetType string
// 	id         uuid.UUID
// }

// type forceTargetReleaseData struct {
// 	inputID uuid.UUID
// }

// var enemyPatrolStateMachine = NewStateMachine("IDLE", map[State]map[Effect]gate{
// 	"IDLE": {
// 		EnemyEventTargetAcquired.Effect: {
// 			outputState: "TARGET_ACQUIRED",
// 			outputUnlockFunc: func(pm ProjectorManager, selfID uuid.UUID) (interface{}, bool) {
// 				playerIDs := pm.Get("PLAYER").ListIdentifiers()
// 				playerPositions := make([]model.Position, len(playerIDs))

// 				for index, playerID := range playerIDs {
// 					playerPositions[index] = pm.Get("Position").Project(playerID).(model.Position)
// 				}

// 				enemyPosition := pm.Get("Position").Project(selfID).(model.Position)

// 				if len(playerIDs) < 1 {
// 					return nil, false
// 				}

// 				isInRange := math.Sqrt(math.Pow(float64(enemyPosition.X-playerPositions[0].X), 2)) <= 5
// 				if !isInRange {
// 					return nil, false
// 				}

// 				return targetData{
// 					targetType: "PLAYER",
// 					id:         playerIDs[0],
// 				}, true
// 			},
// 		},
// 	},
// 	"TARGET_ACQUIRED": {
// 		EnemyEventForceTargetRelease.Effect: {
// 			outputState: "IDLE",
// 			outputUnlockFunc: func(pm ProjectorManager, selfID uuid.UUID) (interface{}, bool) {
// 				targetReleaseData := pm.Get("TargetReleaseLastInput").Project(selfID).(forceTargetReleaseData)

// 				input := pm.Get("EnemyTargetReleaseInput").Project(selfID).(ControllerInput)

// 				if input.ID != uuid.Nil && input.ID != targetReleaseData.inputID {
// 					return forceTargetReleaseData{
// 						inputID: input.ID,
// 					}, true
// 				}

// 				return nil, false
// 			},
// 		},
// 		EnemyEventMove.Effect: gate{
// 			outputState: "TARGET_ACQUIRED",
// 			outputUnlockFunc: func(pm ProjectorManager, selfID uuid.UUID) (interface{}, bool) {
// 				target := pm.Get("Target").Project(selfID).(targetData)

// 				if target.id == uuid.Nil {
// 					return nil, false
// 				}

// 				playerIDs := pm.Get("PLAYER").ListIdentifiers()
// 				playerPositions := make([]model.Position, len(playerIDs))

// 				for index, playerID := range playerIDs {
// 					playerPositions[index] = pm.Get("Position").Project(playerID).(model.Position)
// 				}

// 				position := pm.Get("Position").Project(selfID).(model.Position)

// 				for index, playerID := range playerIDs {
// 					if playerID == target.id {
// 						playerPosition := playerPositions[index]
// 						if playerPosition.X == position.X {
// 							return nil, false
// 						}

// 						if playerPosition.X > position.X {
// 							return model.Position{X: position.X + 1}, true
// 						} else if playerPosition.X < position.X {
// 							return model.Position{X: position.X - 1}, true
// 						}
// 					}
// 				}

// 				return nil, false
// 			},
// 		},
// 	},
// })

// type ControllerInput struct {
// 	ID    uuid.UUID
// 	Value interface{}
// }

// var controllerStateMachine = NewStateMachine("ACTIVE", map[State]map[Effect]gate{
// 	"ACTIVE": {
// 		ControllerEventEnemyTargetRelease.Effect: {
// 			outputState: "ACTIVE",
// 			outputUnlockFunc: func(pm ProjectorManager, selfID uuid.UUID) (interface{}, bool) {
// 				return ControllerInput{
// 					ID: uuid.New(),
// 				}, true
// 			},
// 		},
// 	},
// })
