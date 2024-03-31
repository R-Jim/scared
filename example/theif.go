package example

import (
	"thief/internal/base"
	"thief/internal/model"

	"github.com/google/uuid"
)

const (
	XEnergyPerMoveInput    = 2
	YEnergyPerMoveInput    = 10
	XEnergyConsumePerCycle = 1
	YEnergyConsumePerCycle = 1
)

const (
	EntityTypeThief string = "THIEF"
)

const (
	effectThiefAddEnergy base.Effect = "THIEF_ADD_ENERGY_EFFECT"
	effectThiefMove      base.Effect = "THIEF_MOVE_EFFECT"
)

const (
	stateThiefActive base.State = "THIEF_ACTIVE"
)

const (
	FieldThiefPosition  = "Position"
	fieldThiefMoveInput = "MoveInput"
	FieldThiefEnergy    = "Energy"
)

type thiefMoveOutput struct {
	inputID uuid.UUID
	energy  thiefEnergy
}

type thiefEnergy struct {
	X int
	Y int
}

var thiefStates = map[base.State]map[base.Effect]base.Gate{
	stateThiefActive: {
		effectThiefAddEnergy: base.NewGate(stateThiefActive, func(pm base.ProjectorManager, selfID uuid.UUID) interface{} {
			lastMoveInput := pm.GetEntityProjector(EntityTypeThief).Project(selfID, fieldThiefMoveInput).(uuid.UUID)

			input := pm.GetEntityProjector(EntityTypeController).Project(selfID, fieldControllerThiefInput).(ControllerMoveInput)

			if input.ID == uuid.Nil || input.ID == lastMoveInput {
				return nil
			}

			energy := pm.GetEntityProjector(EntityTypeThief).Project(selfID, FieldThiefEnergy).(thiefEnergy)
			position := pm.GetEntityProjector(EntityTypeThief).Project(selfID, FieldThiefPosition).(model.Position)

			for _, input := range input.Inputs {
				switch input {
				case MoveInputJump:
					if energy.Y <= 0 && position.Y == 0 {
						energy.Y = YEnergyPerMoveInput
					}
				case MoveInputLeft:
					energy.X = -XEnergyPerMoveInput
				case MoveInputRight:
					energy.X = XEnergyPerMoveInput
				}
			}

			return thiefMoveOutput{
				inputID: input.ID,
				energy:  energy,
			}
		}),
		effectThiefMove: base.NewGate(stateThiefActive, func(pm base.ProjectorManager, selfID uuid.UUID) interface{} {
			energy := pm.GetEntityProjector(EntityTypeThief).Project(selfID, FieldThiefEnergy).(thiefEnergy)
			position := pm.GetEntityProjector(EntityTypeThief).Project(selfID, FieldThiefPosition).(model.Position)
			if energy.X == 0 && energy.Y == 0 && position.Y == 0 {
				return nil
			}

			if energy.X > 0 {
				position.X++
			} else if energy.X < 0 {
				position.X--
			}

			if energy.Y > 0 {
				position.Y++
			} else if energy.Y == 0 && position.Y > 0 {
				position.Y--
			}

			return position
		}),
	},
}

var ThiefStateMachine = base.NewStateMachine(EntityTypeThief, stateThiefActive, thiefStates)
