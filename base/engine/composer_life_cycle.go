package engine

import (
	"log"
)

func NewComposerLifeCycle(store *Store, sm stateMachine) *ComposerLifeCycle {
	return &ComposerLifeCycle{
		store:        store,
		stateMachine: sm,
	}
}

type ComposerLifeCycle struct {
	store        *Store
	stateMachine stateMachine
}

func (c ComposerLifeCycle) Operate() {
	resultEvents := []Event{}

	for id, events := range c.store.GetEvents() {
		currentState := c.stateMachine.GetState(events)

		for effect, gate := range c.stateMachine.nodes[currentState] {
			if data, isUnlocked := gate.outputUnlockFunc(id); isUnlocked {
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
