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
