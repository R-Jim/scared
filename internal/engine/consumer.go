package engine

import (
	"log"

	"github.com/google/uuid"
)

type composer struct {
	store            *Store
	projectorManager ProjectorManager
	stateMachine     *stateMachine
}

func NewComposer(store *Store, pm ProjectorManager, sm *stateMachine) composer {
	return composer{
		store:            store,
		projectorManager: pm,
		stateMachine:     sm,
	}
}

type LifeCycleComposer composer

func (c LifeCycleComposer) Operate() {
	resultEvents := []Event{}

	for id, events := range c.store.GetEvents() {
		currentState := c.stateMachine.getState(events)

		for effect, gate := range c.stateMachine.nodes[currentState] {
			if data := gate.outputUnlockFunc(c.projectorManager, id); data != nil {
				resultEvents = append(resultEvents, initEvent(effect, id, data))
				if gate.outputState != currentState {
					break
				}
			}
		}
	}

	for _, resultEvent := range resultEvents {
		if err := c.store.AppendEvent(resultEvent); err != nil {
			log.Fatal(err)
		}
	}
}

type SystemInputComposer composer

func (c SystemInputComposer) TransitionByInput(entityID uuid.UUID, effect Effect, inputData interface{}) (bool, error) {
	events, err := c.store.GetEventsByEntityID(entityID)
	if err != nil {
		return false, err
	}

	state := c.stateMachine.getState(events)
	if err != nil {
		return false, err
	}

	for unlockByEventEffect, gate := range c.stateMachine.nodes[state] {
		if effect != unlockByEventEffect {
			continue
		}

		if data := gate.outputUnlockFunc(c.projectorManager, entityID); data != nil {
			if err := c.store.AppendEvent(initEventWithSystemData(unlockByEventEffect, entityID, data, inputData)); err != nil {
				return false, err
			}
			return true, nil
		}
	}

	return false, nil
}
