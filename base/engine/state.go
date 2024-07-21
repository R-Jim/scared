package engine

import (
	"log"

	"github.com/google/uuid"
)

// State represents State node of the State machine
type State string

type transition struct {
	effect         string
	outputState    State
	transitionFunc func(selfID uuid.UUID) (any, bool)
}

type Nodes map[State][]transition

type stateMachine struct {
	entryTransition transition
	nodes           Nodes // traversable nodes transitions
}

// NewStateMachine returns a new state machine. A state machine must have list of "nodes". An entity must much the first node's transition effect to begin traverse
func NewStateMachine(entryTransition transition, nodes Nodes) stateMachine {
	// TODO: add validation to make sure all nodes is traversable, beginning from the default State
	return stateMachine{
		entryTransition: entryTransition,
		nodes:           nodes,
	}
}

// Returns current state of the instance
func (s stateMachine) getState(events []Event) State {
	if len(events) <= 0 {
		log.Fatalln("No events to get state")
	}

	if events[0].Effect != s.entryTransition.effect {
		log.Fatalln("Invalid entry effect")
	}

	currentState := s.entryTransition.outputState
	for _, event := range events[1:] {
		for _, transition := range s.nodes[currentState] {
			if transition.effect == event.Effect {
				currentState = transition.outputState
				break
			}
		}
	}

	return currentState
}
