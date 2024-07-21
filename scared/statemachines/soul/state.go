package soul

import (
	"thief/base/engine"
	"thief/scared/model"
	"thief/scared/projectors"

	"github.com/google/uuid"
)

const (
	stateUnarmed   engine.State = "Unarmed"
	stateArmed     engine.State = "Armed"
	stateDestroyed engine.State = engine.StateDestroyed
)

const (
	EffectInit    engine.Effect[model.Soul]      = "Init"
	EffectArm     engine.Effect[model.ArmWeapon] = "Arm"
	EffectDisarm  engine.Effect[any]             = "Disarm"
	EffectDestroy engine.Effect[any]             = "Destroy"
)

var StateMachine = engine.NewStateMachine(EffectInit.ToState(stateUnarmed), engine.Nodes{
	stateUnarmed: {
		EffectArm.ToStateWhen(
			stateArmed,
			func(selfID uuid.UUID) (model.ArmWeapon, bool) {
				return model.ArmWeapon{
					OwnerID:    selfID,
					TemplateID: model.ClawID,
				}, true
			},
		),
		EffectDestroy.ToStateWhen(
			stateDestroyed,
			func(selfID uuid.UUID) (interface{}, bool) {
				return nil, projectors.ProjectorEntityType.IsDestroyed(selfID)
			},
		),
	},
	stateArmed: {
		EffectDestroy.ToStateWhen(
			stateDestroyed,
			func(selfID uuid.UUID) (interface{}, bool) {
				return nil, projectors.ProjectorEntityType.IsDestroyed(selfID)
			},
		),
	},
})
