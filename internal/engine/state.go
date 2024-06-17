package engine

import (
	"log"

	"github.com/google/uuid"
)

type outputProducerFunc func(pm ProjectorManager, selfID uuid.UUID) (interface{}, bool)

// gate represents a traversable path from the current node
type gate struct {
	outputState      State              // Result state after gate unlocked
	outputUnlockFunc outputProducerFunc // produces a data to unlock the gate
}

// NewGate returns a new Gate with:
//   - The expected output State.
//   - The required unlock func. The func will retrieve/validate necessary data to produce the output data. If outputUnlockFunc's isUnlocked == true, the caller of the func(usually a composer) will receive the output data and append the corresponding Event to unlock the Gate and traverse to the output State
func NewGate(outputState State, o outputProducerFunc) gate {
	return gate{
		outputState:      outputState,
		outputUnlockFunc: o,
	}
}

// State represents State node of the State machine
type State string

type Nodes map[State]map[Effect]gate

type stateMachine struct {
	entityType EntityType // state machine identifier
	nodes      Nodes      // mapping of traversable nodes
}

// NewStateMachine returns a new state machine. A state machine must have list of "nodes" that traversable from a "defaultState"
func NewStateMachine(entityType EntityType, defaultState State, nodes Nodes) stateMachine {
	if len(nodes) <= 0 {
		log.Fatalf("missing nodes config for state machine[%s]\n", entityType)
	}

	// Init state
	nodes[""] = map[Effect]gate{
		EffectInit: {
			outputState: defaultState, // to init a new state machine to the default state
		},
	}

	// TODO: add validation to make sure all nodes is traversable, beginning from the default State

	return stateMachine{
		entityType: entityType,
		nodes:      nodes,
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
func (s stateMachine) GetState(events []Event) State {
	if len(s.nodes) <= 0 {
		log.Fatalf("no node config for state machine[%s]", s.entityType)
	}

	currentState := State("")
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
