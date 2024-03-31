package example

import (
	"thief/internal/base"
	"thief/internal/model"

	"github.com/google/uuid"
)

const (
	EntityTypeThief string = "THIEF"
)

const (
	effectThiefMove base.Effect = "THIEF_MOVE_EFFECT"
)

const (
	stateThiefActive base.State = "THIEF_ACTIVE"
)

const (
	FieldThiefPosition  = "Position"
	fieldThiefMoveInput = "MoveInput"
)

type thiefMoveOutput struct {
	inputID  uuid.UUID
	position model.Position
}

var thiefStates = map[base.State]map[base.Effect]base.Gate{
	stateThiefActive: {
		effectThiefMove: base.NewGate(stateThiefActive, func(pm base.ProjectorManager, selfID uuid.UUID) interface{} {
			lastMoveInput := pm.GetEntityProjector(EntityTypeThief).Project(selfID, fieldThiefMoveInput).(uuid.UUID)

			input := pm.GetEntityProjector(EntityTypeController).Project(selfID, fieldControllerThiefInput).(ControllerMoveInput)

			if input.ID == uuid.Nil || input.ID == lastMoveInput {
				return nil
			}

			position := pm.GetEntityProjector(EntityTypeThief).Project(selfID, FieldThiefPosition).(model.Position)
			var newPosition model.Position

			switch input.Value {
			case MoveInputJump:
				// TODO: logic here
			case MoveInputLeft:
				newPosition.X = position.X - 1
			case MoveInputRight:
				newPosition.X = position.X + 1
			}

			return thiefMoveOutput{
				inputID:  input.ID,
				position: newPosition,
			}
		}),
	},
}

var ThiefStateMachine = base.NewStateMachine(EntityTypeThief, stateThiefActive, thiefStates)
