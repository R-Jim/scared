package example

import (
	"thief/internal/engine"
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
	effectThiefAddEnergy engine.Effect = "THIEF_ADD_ENERGY_EFFECT"
	effectThiefMove      engine.Effect = "THIEF_MOVE_EFFECT"
)

const (
	stateThiefActive engine.State = "THIEF_ACTIVE"
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

var thiefStates = map[engine.State]map[engine.Effect]engine.Gate{
	stateThiefActive: {
		effectThiefAddEnergy: engine.NewGate(stateThiefActive, func(pm engine.ProjectorManager, selfID uuid.UUID) interface{} {
			thiefProjector := pm.Get(EntityTypeThief)

			lastMoveInput := thiefProjector.Project(selfID, fieldThiefMoveInput).(uuid.UUID)
			input := pm.Get(EntityTypeController).Project(selfID, fieldControllerThiefInput).(ControllerMoveInput)

			if input.ID == uuid.Nil || input.ID == lastMoveInput {
				return nil
			}

			energy := thiefProjector.Project(selfID, FieldThiefEnergy).(thiefEnergy)
			position := thiefProjector.Project(selfID, FieldThiefPosition).(model.Position)

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
		effectThiefMove: engine.NewGate(stateThiefActive, func(pm engine.ProjectorManager, selfID uuid.UUID) interface{} {
			thiefProjector := pm.Get(EntityTypeThief)

			energy := thiefProjector.Project(selfID, FieldThiefEnergy).(thiefEnergy)
			position := thiefProjector.Project(selfID, FieldThiefPosition).(model.Position)
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

var ThiefStateMachine = engine.NewStateMachine(EntityTypeThief, stateThiefActive, thiefStates)
