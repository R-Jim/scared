package ship

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
	EffectInit    engine.Effect[model.ShipGuard] = "Init"
	EffectArm     engine.Effect[model.ArmWeapon] = "Arm"
	EffectDisarm  engine.Effect[any]             = "Disarm"
	EffectDestroy engine.Effect[any]             = "Destroy"
)

var StateMachineShipGuard = engine.NewStateMachine(EffectInit.ToState(stateUnarmed), engine.Nodes{
	stateUnarmed: {
		EffectArm.ToStateWhen(
			stateArmed,
			func(selfID uuid.UUID) (model.ArmWeapon, bool) {
				weaponTemplateID := projectors.ProjectorEntityAssignedWeapon.Project(selfID)
				if weaponTemplateID == uuid.Nil {
					return model.ArmWeapon{}, false
				}

				return model.ArmWeapon{
					OwnerID:    selfID,
					TemplateID: weaponTemplateID,
				}, true
			},
		),
		EffectDestroy.ToStateWhen(
			stateDestroyed,
			func(selfID uuid.UUID) (any, bool) {
				return nil, projectors.ProjectorEntityType.IsDestroyed(selfID)
			},
		),
	},
	stateArmed: {
		EffectDisarm.ToStateWhen(
			stateArmed,
			func(selfID uuid.UUID) (any, bool) {
				weaponTemplateID := projectors.ProjectorEntityAssignedWeapon.Project(selfID)
				if weaponTemplateID != uuid.Nil {
					return nil, false
				}

				return nil, true
			},
		),
		EffectDestroy.ToStateWhen(
			stateDestroyed,
			func(selfID uuid.UUID) (any, bool) {
				return nil, projectors.ProjectorEntityType.IsDestroyed(selfID)
			},
		),
	},
})

const (
	stateBlessingAltarActive    engine.State = "Active"
	stateBlessingAltarDestroyed engine.State = engine.StateDestroyed

	EffectBlessingAltarInit    engine.Effect[uuid.UUID]             = "Init"
	EffectBlessAcolyteToKnight engine.Effect[model.SpawnKnightData] = "EffectBlessAcolyteToKnight"
)

var StateMachineShipBlessingAltar = engine.NewStateMachine(EffectBlessingAltarInit.ToState(stateBlessingAltarActive), engine.Nodes{
	stateBlessingAltarActive: {
		EffectBlessAcolyteToKnight.ToStateWhen(
			stateBlessingAltarActive,
			func(selfID uuid.UUID) (model.SpawnKnightData, bool) {
				altarData := projectors.ProjectorBlessingAltar.Project(selfID)

				numberOfAssignedAcolyte := projectors.ProjectorAcolyte.Project(selfID)

				if altarData.NumberOfBlessedAcolyte >= numberOfAssignedAcolyte {
					return model.SpawnKnightData{}, false
				}

				return model.SpawnKnightData{
					Position: projectors.ProjectorPosition.Project(altarData.OwnerID),
				}, true
			},
		),
	},
})
