package state

import (
	"fmt"

	"github.com/google/uuid"
	pkgerrors "github.com/pkg/errors"
)

type EventDataProducerFunc func(selfID uuid.UUID, pm ProjectorManager, inputData interface{}) interface{}

type Gate struct {
	unlockByEventEffect string                // Only input effect of this is allowed to unlock gate
	eventProducerFunc   EventDataProducerFunc // produces an event to unlock the gate
	nextState           string                // Result state after gate unlocked
}

type Node struct {
	State string // State identifier
	Gates []Gate // gates will run unlock condition check every life cycle
	// PassiveGates []Gate // Passive gates only run when invoke by external caller
}

// transitionSet is used to map out the entire state machine and it traversable nodes
type transitionSet struct {
	Effect      string
	ResultState string
}

type StateMachine struct {
	entityType string
	nodes      []Node

	transitionMapping map[string][]transitionSet
}

func NewStateMachine(entityType string, nodes []Node) StateMachine {
	stateMachine := StateMachine{
		entityType: entityType,
		nodes:      nodes,
	}

	transitionMapping := map[string][]transitionSet{}
	for _, node := range nodes {
		resultTransitionSet := transitionMapping[node.State]
		if resultTransitionSet == nil {
			resultTransitionSet = []transitionSet{}
		}

		for _, gate := range node.Gates {
			var isAlreadyExist bool

			// validate if effect transition already in transition set
			for _, transitionSet := range resultTransitionSet {
				if gate.unlockByEventEffect == transitionSet.Effect && gate.nextState == transitionSet.ResultState {
					isAlreadyExist = true
					break
				}
			}

			if !isAlreadyExist {
				resultTransitionSet = append(resultTransitionSet, transitionSet{
					Effect:      gate.unlockByEventEffect,
					ResultState: gate.nextState,
				})
			}
		}

		transitionMapping[node.State] = resultTransitionSet
	}

	stateMachine.transitionMapping = transitionMapping

	return stateMachine
}

// Returns current node of the instance
func (s StateMachine) GetNode(instanceID uuid.UUID, events []Event) (Node, error) {
	if len(s.nodes) <= 0 {
		return Node{}, pkgerrors.WithStack(fmt.Errorf("no node config for state machine[%s]", s.entityType))
	}

	var currentNode Node
	for _, event := range events {
		if nextNode := getNextNode(currentNode.State, event.Effect, s.nodes, s.transitionMapping); nextNode != nil {
			currentNode = *nextNode
		}
	}

	return currentNode, nil
}

func getNextNode(currentState string, effect string, nodes []Node, transitionMapping map[string][]transitionSet) *Node {
	if len(nodes) <= 0 {
		return nil
	}

	if effect == "INIT" {
		return &nodes[0]
	}

	transitionSet := transitionMapping[currentState]

	for _, transition := range transitionSet {
		if effect == transition.Effect {
			for _, node := range nodes {
				if transition.ResultState == node.State {
					return &node
				}
			}
		}
	}

	return nil
}
