package example

import (
	"thief/internal/engine"

	"github.com/google/uuid"
)

const (
	EntityTypeController = "CONTROLLER"
)

const (
	EffectControllerMove engine.Effect = "CONTROLLER_MOVE_EFFECT"
)

const (
	stateControllerActive engine.State = "CONTROLLER_ACTIVE"
)

const (
	fieldControllerThiefInput = "ThiefInput"
)

type MoveInput string

const (
	MoveInputLeft  MoveInput = "LEFT"
	MoveInputRight MoveInput = "RIGHT"
	MoveInputJump  MoveInput = "JUMP"
)

type ControllerMoveInput struct {
	ID     uuid.UUID
	Inputs []MoveInput
}

var controllerStates = map[engine.State]map[engine.Effect]engine.Gate{
	stateControllerActive: {
		EffectControllerMove: engine.NewGate(stateControllerActive, func(pm engine.ProjectorManager, selfID uuid.UUID) interface{} {
			return ControllerMoveInput{
				ID: uuid.New(),
			}
		}),
	},
}

var ControllerStateMachine = engine.NewStateMachine(EntityTypeController, stateControllerActive, controllerStates)
