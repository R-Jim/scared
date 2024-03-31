package example

import (
	"thief/internal/base"

	"github.com/google/uuid"
)

const (
	EntityTypeController = "CONTROLLER"
)

const (
	EffectControllerMove base.Effect = "CONTROLLER_MOVE_EFFECT"
)

const (
	stateControllerActive base.State = "CONTROLLER_ACTIVE"
)

const (
	fieldControllerThiefInput = "ThiefInput"
)

type moveInput string

const (
	MoveInputLeft  moveInput = "LEFT"
	MoveInputRight moveInput = "RIGHT"
	MoveInputJump  moveInput = "JUMP"
)

type ControllerMoveInput struct {
	ID    uuid.UUID
	Value moveInput
}

var controllerStates = map[base.State]map[base.Effect]base.Gate{
	stateControllerActive: {
		EffectControllerMove: base.NewGate(stateControllerActive, func(pm base.ProjectorManager, selfID uuid.UUID) interface{} {
			return ControllerMoveInput{
				ID: uuid.New(),
			}
		}),
	},
}

var ControllerStateMachine = base.NewStateMachine(EntityTypeController, stateControllerActive, controllerStates)
