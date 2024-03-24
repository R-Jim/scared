package state

import (
	"log"

	"github.com/google/uuid"
)

type outputProducerFunc func(selfID uuid.UUID, pm ProjectorManager) interface{}

type gate struct {
	nextState          State              // Result state after gate unlocked
	outputProducerFunc outputProducerFunc // produces a data to unlock the gate
}

func NewGate(nextState State, o outputProducerFunc) gate {
	return gate{
		nextState:          nextState,
		outputProducerFunc: o,
	}
}

type State string

type stateMachine struct {
	entityType   string
	defaultState State
	nodes        map[State]map[Effect]gate
}

func NewStateMachine(entityType string, defaultState State, nodes map[State]map[Effect]gate) stateMachine {
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
			currentState = gate.nextState
		}
	}

	return currentState
}
