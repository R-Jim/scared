package base

import (
	"log"

	"github.com/google/uuid"
)

type outputProducerFunc func(pm ProjectorManager, selfID uuid.UUID) interface{}

// Gate represents a traversable path from the current node
type Gate struct {
	outputState      State              // Result state after gate unlocked
	outputUnlockFunc outputProducerFunc // produces a data to unlock the gate
}

// NewGate returns a new Gate with:
//   - The expected output State.
//   - The required unlock func. The func will retrieve/validate necessary data to produce the output data. If output data != nil, the caller of the func(usually a composer) will receive the output data and append the corresponding Event to unlock the Gate and traverse to the next State node
func NewGate(outputState State, o outputProducerFunc) Gate {
	return Gate{
		outputState:      outputState,
		outputUnlockFunc: o,
	}
}

// State represents State node of the State machine
type State string

type stateMachine struct {
	entityType   string                    // state machine identifier
	defaultState State                     // default beginning state
	nodes        map[State]map[Effect]Gate // mapping of traversable nodes
}

// NewStateMachine returns a new state machine. A state machine must have list of "nodes" that traversable from a "defaultState"
func NewStateMachine(entityType string, defaultState State, nodes map[State]map[Effect]Gate) stateMachine {
	if len(nodes) <= 0 {
		log.Fatalf("missing nodes config for state machine[%s]\n", entityType)
	}

	// TODO: add validation to make sure all nodes is traversable, beginning from the default State

	return stateMachine{
		entityType:   entityType,
		defaultState: defaultState,
		nodes:        nodes,
	}

	// For render only
	// transitionMapping := map[State]map[Effect]State{}
	// for state, gates := range nodes {
	// 	resultTransitionSet := transitionMapping[state]
	// 	if resultTransitionSet == nil {
	// 		resultTransitionSet = map[Effect]State{}
	// 	}

	// 	for effect, gate := range gates {
	// 		var isAlreadyExist bool

	// 		// validate if effect transition already in transition set
	// 		for transitionEffect, transitionState := range resultTransitionSet {
	// 			if effect == transitionEffect && gate.nextState == transitionState {
	// 				isAlreadyExist = true
	// 				break
	// 			}
	// 		}

	// 		if !isAlreadyExist {
	// 			resultTransitionSet[effect] = gate.nextState
	// 		}
	// 	}

	// 	transitionMapping[state] = resultTransitionSet
	// }

	// stateMachine.transitionMapping = transitionMapping

	// return stateMachine
}

// Returns current state of the instance
func (s stateMachine) getState(events []Event) State {
	if len(s.nodes) <= 0 {
		log.Fatalf("no node config for state machine[%s]", s.entityType)
	}

	currentState := s.defaultState
	for _, event := range events {
		node, isExist := s.nodes[currentState]
		if !isExist {
			return currentState
		}

		gate, isExist := node[event.Effect]
		if isExist {
			currentState = gate.outputState
		}
	}

	return currentState
}
