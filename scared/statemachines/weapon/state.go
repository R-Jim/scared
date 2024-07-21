package weapon

import (
	"thief/base/engine"
	"thief/scared/model"
	"thief/scared/projectors"

	"github.com/google/uuid"
)

const (
	stateIdle      engine.State = "Idle"
	stateActive    engine.State = "Active"
	stateCoolDown  engine.State = "CoolDown"
	stateDestroyed engine.State = engine.StateDestroyed
)

const (
	EffectInit     engine.Effect[model.EquippedWeapon] = "Init"
	EffectActive   engine.Effect[any]                  = "Active"
	EffectHitEnemy engine.Effect[model.WeaponLog]      = "HitEnemy"
	EffectCoolDown engine.Effect[int]                  = "CoolDown"
	EffectReActive engine.Effect[any]                  = "ReActive"
	EffectDestroy  engine.Effect[any]                  = "Destroy"
)

func addEffectDestroyTransitionToNodes(nodes engine.Nodes) engine.Nodes {
	for key, transitions := range nodes {
		nodes[key] = append(transitions, EffectDestroy.ToStateWhen(
			stateDestroyed,
			func(selfID uuid.UUID) (interface{}, bool) {
				equippedWeapon := projectors.ProjectorEquippedWeapon.Project(selfID)
				if projectors.ProjectorEntityType.IsDestroyed(equippedWeapon.OwnerID) {
					return nil, true
				}

				return nil, false
			},
		))
	}

	return nodes
}

var StateMachineEquippedWeapon = engine.NewStateMachine(EffectInit.ToState(stateIdle), addEffectDestroyTransitionToNodes(engine.Nodes{
	stateIdle: {
		EffectActive.ToStateWhen(
			stateActive,
			func(selfID uuid.UUID) (any, bool) {
				equippedWeapon := projectors.ProjectorEquippedWeapon.Project(selfID)
				ownerEntityType := projectors.ProjectorEntityType.Project(equippedWeapon.OwnerID)

				if ownerEntityType == model.EntityTypeSoul {
					return nil, true
				}

				return nil, projectors.ProjectorActiveWeapon.Project(selfID)
			},
		),
	},
	stateActive: {
		EffectHitEnemy.ToStateWhen(
			stateCoolDown,
			func(selfID uuid.UUID) (model.WeaponLog, bool) {
				equippedWeapon := projectors.ProjectorEquippedWeapon.Project(selfID)
				ownerPosition := projectors.ProjectorPosition.Project(equippedWeapon.OwnerID)
				ownerEntityType := projectors.ProjectorEntityType.Project(equippedWeapon.OwnerID)

				targetIDs := projectors.ProjectorEntityType.ListIdentifiers(
					func(t model.EntityType) bool {
						for _, targetType := range model.EntityTypeAttackTargetMapping[ownerEntityType] {
							if targetType == t {
								return true
							}
						}
						return false
					},
				)

				var nearestTargetID uuid.UUID
				var nearestTargetDistance float64

				weaponRange := projectors.ProjectorEquippedWeaponRange.Project(selfID)

				for _, targetID := range targetIDs {
					targetPosition := projectors.ProjectorPosition.Project(targetID)

					distance := ownerPosition.DistanceOf(targetPosition)
					if (nearestTargetID == uuid.Nil || distance < nearestTargetDistance) && distance <= float64(weaponRange) {
						nearestTargetID = targetID
						nearestTargetDistance = distance
					}
				}

				if nearestTargetID == uuid.Nil {
					return model.WeaponLog{}, false
				}

				weaponDamage := projectors.ProjectorEquippedWeaponDamage.Project(selfID)
				coolDown := projectors.ProjectorEquippedWeaponCoolDown.Project(selfID)

				// log.Printf("fire range: %d, damage: %d\n", weaponRange, weaponDamage.Plus-weaponDamage.Minus)
				return model.WeaponLog{
					CoolDown:     coolDown,
					DevotionCost: projectors.ProjectorEquippedWeaponDevotionCost.Project(selfID),
					Log: model.Log[model.Stat]{
						ID: uuid.New(),
						Value: model.Stat{
							Minus: weaponDamage.Minus,
							Plus:  weaponDamage.Plus,
						},
						Targets: []uuid.UUID{nearestTargetID},
					},
				}, true
			},
		),
	},
	stateCoolDown: {
		EffectCoolDown.ToStateWhen(
			stateCoolDown,
			func(selfID uuid.UUID) (int, bool) {
				coolDownCount := projectors.ProjectorEquippedWeaponCoolDownCount.Project(selfID)
				if coolDownCount <= 0 {
					return 0, false
				}

				return 1, true
			},
		),
		EffectReActive.ToStateWhen(
			stateActive,
			func(selfID uuid.UUID) (interface{}, bool) {
				coolDownCount := projectors.ProjectorEquippedWeaponCoolDownCount.Project(selfID)
				if coolDownCount > 0 {
					return nil, false
				}

				return nil, true
			},
		),
	},
}))

const (
	stateRuneSlotEmpty     engine.State = "Empty"
	stateRuneSlotActive    engine.State = "Active"
	stateRuneSlotDestroyed engine.State = engine.StateDestroyed
)

const (
	EffectRuneSlotInit    engine.Effect[model.WeaponRuneSlotType] = "Init"
	EffectRuneSlotActive  engine.Effect[uuid.UUID]                = "Active"
	EffectRuneSlotDestroy engine.Effect[any]                      = "Destroy"
)

func addEffectDestroyTransitionToEquippedRuneSlotNodes(nodes engine.Nodes) engine.Nodes {
	for key, transitions := range nodes {
		nodes[key] = append(transitions, EffectRuneSlotDestroy.ToStateWhen(
			stateRuneSlotDestroyed,
			func(selfID uuid.UUID) (interface{}, bool) {
				equippedWeaponIDs := projectors.ProjectorEquippedWeapon.ListIdentifiers(func(ew model.EquippedWeapon) bool {
					for _, runeSlotID := range ew.RuneSlotIDs {
						if selfID == runeSlotID {
							return true
						}
					}
					return false
				})

				return nil, len(equippedWeaponIDs) <= 0
			},
		))
	}

	return nodes
}

var StateMachineEquippedWeaponRuneSlot = engine.NewStateMachine(EffectRuneSlotInit.ToState(stateRuneSlotEmpty), addEffectDestroyTransitionToEquippedRuneSlotNodes(engine.Nodes{
	stateRuneSlotEmpty: {
		EffectRuneSlotActive.ToStateWhen(
			stateRuneSlotActive,
			func(selfID uuid.UUID) (uuid.UUID, bool) {
				runeTemplateID := projectors.ProjectorRuneTemplateForEquippedWeaponRuneSlot.Project(selfID)
				if runeTemplateID != uuid.Nil {
					return runeTemplateID, true
				}

				return uuid.Nil, false
			},
		),
	},
	stateRuneSlotActive: {},
}))
