package engine

import "github.com/google/uuid"

type ComposerDestroyer struct {
	store        *Store
	stateMachine stateMachine
	plannedIDs   []uuid.UUID
}

func NewComposerDestroyer(store *Store, sm stateMachine) *ComposerDestroyer {
	return &ComposerDestroyer{
		store:        store,
		stateMachine: sm,
	}
}

const (
	StateDestroyed State = "DESTROYED"
)

func (c *ComposerDestroyer) PlanDestroyIDs() {
	for entityID, events := range c.store.GetEvents() {
		currentState := c.stateMachine.GetState(events)

		if currentState == StateDestroyed {
			c.plannedIDs = append(c.plannedIDs, entityID)
		}
	}
}

func (c *ComposerDestroyer) Commit() {
	for _, entityID := range c.plannedIDs {
		c.store.destroySet(entityID)
	}
}
