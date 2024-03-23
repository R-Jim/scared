package state

import (
	"log"
	"testing"
	"thief/internal/model"
	"time"

	"github.com/google/uuid"
)

type PlayerProjector struct {
	playerID       uuid.UUID
	playerPosition model.Position
}

func NewPlayerProjector() Projector {
	return PlayerProjector{
		playerID: uuid.New(),
		playerPosition: model.Position{
			X: 5,
		},
	}
}

func (p PlayerProjector) Project(identifier uuid.UUID, field string) interface{} {
	if field == "Position" {
		return p.playerPosition
	}
	return nil
}

func (p PlayerProjector) ListIdentifiers() []uuid.UUID {
	return []uuid.UUID{p.playerID}
}

type EnemyProjector struct {
	enemyStore *Store
}

func NewEnemyProjector(enemyStore *Store) Projector {
	return EnemyProjector{
		enemyStore: enemyStore,
	}
}

func (p EnemyProjector) Project(identifier uuid.UUID, field string) interface{} {
	events, err := p.enemyStore.GetEventsByEntityID(identifier)
	if err != nil {
		log.Fatalln(err)
	}

	for i := len(events) - 1; i >= 0; i-- {
		event := events[i]
		switch field {
		case "Position":
			switch event.Effect {
			case EnemyEventMove.Effect:
				position, err := ParseData[model.Position](event)
				if err != nil {
					log.Fatalln(err)
				}
				return position
			case EnemyEventInit.Effect:
				return model.Position{}
			}
		case "Target":
			switch event.Effect {
			case EnemyEventTargetRelease.Effect, EnemyEventTargetAcquired.Effect:
				targetData, err := ParseData[targetData](event)
				if err != nil {
					log.Fatalln(err)
				}
				return targetData
			case EnemyEventInit.Effect:
				return targetData{}
			}
		}
	}

	return nil
}

func (p EnemyProjector) ListIdentifiers() []uuid.UUID {
	identifiers := []uuid.UUID{}
	for id, _ := range p.enemyStore.GetEvents() {
		identifiers = append(identifiers, id)
	}

	return identifiers
}

func Test_EnemyStateMachine(t *testing.T) {
	enemyStore := NewStore()

	enemyStateMachine := NewStateMachine(EnemyPatrolStateMachine.entityType, EnemyPatrolStateMachine.nodes)
	enemyID := uuid.New()

	enemyStore.AppendEvent(InitEvent(EnemyEventInit.Effect, enemyID, nil))

	pm := ProjectorManager{
		EntityProjectorMap: map[string]Projector{
			"PLAYER": NewPlayerProjector(),
			"ENEMY":  NewEnemyProjector(&enemyStore),
		},
	}

	enemyComposer := Composer{
		store:            &enemyStore,
		projectorManager: pm,
		stateMachine:     &enemyStateMachine,
	}

	for i := 0; i < 4; i++ {
		enemyComposer.OperateStateLifeCycle()

		position := pm.GetEntityProjector("ENEMY").Project(enemyID, "Position").(model.Position)
		if position.X != i {
			t.Fatal("fail move enemy to target")
		}
	}

	target := pm.GetEntityProjector("ENEMY").Project(enemyID, "Target").(targetData)
	if target.id == uuid.Nil {
		t.Fatal("no target for release")
	}

	result, err := enemyComposer.RequestStateTransition(enemyID, EnemyEventTargetRelease.Effect, nil)
	if err != nil || !result {
		t.Fatal("force release target failed")
	}

	target = pm.GetEntityProjector("ENEMY").Project(enemyID, "Target").(targetData)
	if target.id != uuid.Nil {
		t.Fatal("should force release target")
	}
}

type affectedTargetCondition struct {
	entityTypes []string
	position    model.Position
	startAt     time.Time
	endAt       time.Time
}

type damageData struct {
	value int
}

type strikeData struct {
	id                      uuid.UUID
	sourceID                uuid.UUIDs
	affectedTargetCondition affectedTargetCondition
	damageData              damageData
}

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
