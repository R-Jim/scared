package state

import (
	"log"

	"github.com/google/uuid"
)

type Composer struct {
	store            *Store
	projectorManager ProjectorManager
	stateMachine     *StateMachine
}

func NewComposer(store *Store, projectorManager ProjectorManager, stateMachine *StateMachine) Composer {
	return Composer{
		store:            store,
		projectorManager: projectorManager,
		stateMachine:     stateMachine,
	}
}

func (c Composer) OperateStateLifeCycle() {
	resultEvents := []Event{}

	for id, events := range c.store.GetEvents() {
		node, err := c.stateMachine.GetNode(id, events)
		if err != nil {
			log.Fatal(err)
		}

		for _, gate := range node.ActiveGates {
			if data := gate.eventProducerFunc(id, c.projectorManager, nil); data != nil {
				resultEvents = append(resultEvents, InitEvent(gate.unlockByEventEffect, id, data))
				break
			}
		}
	}

	for _, resultEvent := range resultEvents {
		if err := c.store.AppendEvent(resultEvent); err != nil {
			log.Fatal(err)
		}
	}
}

func (c Composer) RequestStateTransition(entityID uuid.UUID, effect string, inputData interface{}) (bool, error) {
	events, err := c.store.GetEventsByEntityID(entityID)
	if err != nil {
		return false, err
	}

	node, err := c.stateMachine.GetNode(entityID, events)
	if err != nil {
		return false, err
	}

	for _, gate := range node.PassiveGates {
		if effect != gate.unlockByEventEffect {
			continue
		}

		if data := gate.eventProducerFunc(entityID, c.projectorManager, inputData); data != nil {
			if err := c.store.AppendEvent(InitEvent(gate.unlockByEventEffect, entityID, data)); err != nil {
				return false, err
			}
			return true, nil
		}
	}

	return false, nil
}
