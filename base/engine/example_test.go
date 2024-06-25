package engine

// TODO: update outdated example test
// type PositionProjector struct {
// 	playerID       uuid.UUID
// 	playerPosition model.Position
// }

// func NewPlayerProjector() Projector {
// 	return PositionProjector{
// 		playerID: uuid.New(),
// 		playerPosition: model.Position{
// 			X: 5,
// 		},
// 	}
// }

// func (p PositionProjector) Project(identifier uuid.UUID) interface{} {
// 	return p.playerPosition
// }

// func (p PositionProjector) ListIdentifiers() []uuid.UUID {
// 	return []uuid.UUID{p.playerID}
// }

// type EnemyProjector struct {
// 	enemyStore *Store
// }

// func NewEnemyProjector(enemyStore *Store) Projector {
// 	return EnemyProjector{
// 		enemyStore: enemyStore,
// 	}
// }

// func (p EnemyProjector) Project(identifier uuid.UUID, field Field) interface{} {
// 	events, err := p.enemyStore.GetEventsByEntityID(identifier)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	for i := len(events) - 1; i >= 0; i-- {
// 		event := events[i]
// 		switch field {
// 		case "Position":
// 			switch event.Effect {
// 			case EnemyEventMove.Effect:
// 				return ParseData[model.Position](event)
// 			case EnemyEventInit.Effect:
// 				return model.Position{}
// 			}
// 		case "Target":
// 			switch event.Effect {
// 			case EnemyEventTargetRelease.Effect, EnemyEventTargetAcquired.Effect:
// 				return ParseData[targetData](event)
// 			case EnemyEventInit.Effect, EnemyEventForceTargetRelease.Effect:
// 				return targetData{}
// 			}
// 		case "TargetReleaseLastInput":
// 			switch event.Effect {
// 			case EnemyEventForceTargetRelease.Effect:
// 				return ParseData[forceTargetReleaseData](event)
// 			default:
// 				return forceTargetReleaseData{}
// 			}
// 		}
// 	}

// 	return nil
// }

// func (p EnemyProjector) ListIdentifiers() []uuid.UUID {
// 	identifiers := []uuid.UUID{}
// 	for id, _ := range p.enemyStore.GetEvents() {
// 		identifiers = append(identifiers, id)
// 	}

// 	return identifiers
// }

// type ControllerProjector struct {
// 	controllerStore *Store
// }

// func NewControllerProjector(controllerStore *Store) Projector {
// 	return ControllerProjector{
// 		controllerStore: controllerStore,
// 	}
// }

// func (p ControllerProjector) Project(identifier uuid.UUID, field Field) interface{} {
// 	events, err := p.controllerStore.GetEventsByEntityID(identifier)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	for i := len(events) - 1; i >= 0; i-- {
// 		event := events[i]
// 		switch field {
// 		case "EnemyTargetReleaseInput":
// 			switch event.Effect {
// 			case ControllerEventEnemyTargetRelease.Effect:
// 				return ParseData[ControllerInput](event)
// 			case ControllerEventInit.Effect:
// 				return ControllerInput{}
// 			}
// 		}
// 	}

// 	return nil
// }

// func (p ControllerProjector) ListIdentifiers() []uuid.UUID {
// 	identifiers := []uuid.UUID{}
// 	for id, _ := range p.controllerStore.GetEvents() {
// 		identifiers = append(identifiers, id)
// 	}

// 	return identifiers
// }

// func Test_EnemyStateMachine(t *testing.T) {
// 	enemyStore := NewStore()
// 	controllerStore := NewStore()

// 	enemyStateMachine := enemyPatrolStateMachine
// 	enemyID := uuid.New()

// 	enemyStore.AppendEvent(initEvent(EnemyEventInit.Effect, enemyID, nil))
// 	controllerStore.AppendEvent(initEvent(ControllerEventInit.Effect, enemyID, nil))

// 	pm := ProjectorManager{
// 		mapping: map[EntityType]Projector{
// 			"PLAYER":     NewPlayerProjector(),
// 			"ENEMY":      NewEnemyProjector(&enemyStore),
// 			"CONTROLLER": NewControllerProjector(&controllerStore),
// 		},
// 	}

// 	enemyComposer := LifeCycleComposer{
// 		store:            &enemyStore,
// 		projectorManager: pm,
// 		stateMachine:     &enemyStateMachine,
// 	}

// 	controllerComposer := SystemInputComposer{
// 		store:            &controllerStore,
// 		projectorManager: pm,
// 		stateMachine:     &controllerStateMachine,
// 	}

// 	for i := 0; i < 4; i++ {
// 		enemyComposer.Operate()

// 		position := pm.Get("ENEMY").Project(enemyID, "Position").(model.Position)
// 		if position.X != i {
// 			t.Fatal("fail move enemy to target")
// 		}
// 	}

// 	target := pm.Get("ENEMY").Project(enemyID, "Target").(targetData)
// 	if target.id == uuid.Nil {
// 		t.Fatal("no target for release")
// 	}

// 	result, err := controllerComposer.TransitionByInput(enemyID, ControllerEventEnemyTargetRelease.Effect, nil)
// 	if err != nil || !result {
// 		t.Fatal("force release target failed")
// 	}

// 	input := pm.Get("CONTROLLER").Project(enemyID, "EnemyTargetReleaseInput").(ControllerInput)
// 	if input.ID == uuid.Nil {
// 		t.Fatal("should have controller input")
// 	}

// 	enemyComposer.Operate()

// 	target = pm.Get("ENEMY").Project(enemyID, "Target").(targetData)
// 	if target.id != uuid.Nil {
// 		t.Fatal("should force release target")
// 	}

// 	// double check, auto re-assign target
// 	enemyComposer.Operate()

// 	target = pm.Get("ENEMY").Project(enemyID, "Target").(targetData)
// 	if target.id == uuid.Nil {
// 		t.Fatal("double check, should auto re-assign target")
// 	}

// 	position := pm.Get("ENEMY").Project(enemyID, "Position").(model.Position)
// 	if position.X != 3 {
// 		t.Fatal("double check, should remain position")
// 	}
// }

// type affectedTargetCondition struct {
// 	entityTypes []string
// 	position    model.Position
// 	startAt     time.Time
// 	endAt       time.Time
// }

// type damageData struct {
// 	value int
// }

// type strikeData struct {
// 	id                      uuid.UUID
// 	sourceID                uuid.UUIDs
// 	affectedTargetCondition affectedTargetCondition
// 	damageData              damageData
// }

/*
	logic implementation for player strike enemy/enemy strike player
	1) ActiveGate or user input request gate activation. Attack source has Gate that produce strike event
	2) Target has ActiveGate to check for Strike attempt from possible source. The ActiveGate produce hit event

	QA:
	1) How target know the strike already performed and has corresponding hit?
	-> Each strike has an id, the ActiveGate of the target will check possible attack source's strike data if these are any strike without a corresponding hit.
	2) How to reduce the number of strike to corresponding hit checks?
	-> Strike data can has start_at datetime. Target check the strike at until now - buff_time <= start_at
*/
