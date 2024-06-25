package engine

import (
	"github.com/google/uuid"
)

type ComposerExternalInput struct {
	store        *Store
	stateMachine stateMachine
}

func NewComposerExternalInput(store *Store, sm stateMachine) *ComposerExternalInput {
	return &ComposerExternalInput{
		store:        store,
		stateMachine: sm,
	}
}

type InputData struct {
	Key            string
	TransitionData interface{}
}

func (c ComposerExternalInput) TransitionByInput(entityID uuid.UUID, effect Effect, key string) (bool, error) {
	events, err := c.store.GetEventsByEntityID(entityID)
	if err != nil {
		return false, err
	}

	state := c.stateMachine.GetState(events)
	if err != nil {
		return false, err
	}

	for unlockByEventEffect, gate := range c.stateMachine.nodes[state] {
		if effect != unlockByEventEffect {
			continue
		}

		if data, isUnlocked := gate.outputUnlockFunc(entityID); isUnlocked {
			if err := c.store.AppendEvent(initEvent(unlockByEventEffect, entityID, InputData{
				Key:            key,
				TransitionData: data,
			})); err != nil {
				return false, err
			}
			return true, nil
		}
	}

	return false, nil
}
