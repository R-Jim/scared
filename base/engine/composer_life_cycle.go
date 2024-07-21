package engine

import (
	"log"

	"github.com/google/uuid"
)

const (
	StateDestroyed State = "Destroyed"
)

func NewComposerLifeCycle(store *Store, sm stateMachine) *ComposerLifeCycle {
	return &ComposerLifeCycle{
		store:        store,
		stateMachine: sm,
	}
}

type ComposerLifeCycle struct {
	store               *Store
	stateMachine        stateMachine
	plannedDestroyedIDs []uuid.UUID
}

func (c *ComposerLifeCycle) Operate() {
	resultEvents := []Event{}

	for id, events := range c.store.GetEvents() {
		currentState := c.stateMachine.getState(events)

		for _, transition := range c.stateMachine.nodes[currentState] {
			if transition.transitionFunc == nil {
				continue
			}

			if data, isUnlocked := transition.transitionFunc(id); isUnlocked {
				resultEvents = append(resultEvents, Event{
					ID:       uuid.New(),
					EntityID: id,
					Effect:   transition.effect,
					Data:     data,
				})
				currentState = transition.outputState
				break
			}
		}

		if currentState == StateDestroyed {
			c.plannedDestroyedIDs = append(c.plannedDestroyedIDs, id)
		}
	}

	for _, resultEvent := range resultEvents {
		if err := c.store.AppendEvent(resultEvent); err != nil {
			log.Fatal(err)
		}
	}
}

func (c *ComposerLifeCycle) CommitDestroyedIDs() {
	if len(c.plannedDestroyedIDs) > 0 {
		log.Printf("[%s] remove IDs: %v\n", c.store.name, c.plannedDestroyedIDs)
	}

	for _, entityID := range c.plannedDestroyedIDs {
		c.store.destroySet(entityID)
	}

	c.plannedDestroyedIDs = []uuid.UUID{}
}
