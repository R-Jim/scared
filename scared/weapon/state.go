package weapon

import (
	"log"
	"thief/base/engine"
	"thief/scared"

	"github.com/google/uuid"
)

const (
	stateActive    engine.State = "Active"
	stateCoolDown  engine.State = "CoolDown"
	stateDestroyed engine.State = "Destroyed"
)

const (
	EffectHitEnemy engine.Effect = "HitEnemy"
	EffectCoolDown engine.Effect = "CoolDown"
	EffectReActive engine.Effect = "ReActive"
)

var StateMachineEquippedWeapon = engine.NewStateMachine(stateActive, engine.Nodes{
	stateActive: {
		EffectHitEnemy: engine.NewGate(
			stateCoolDown,
			func(selfID uuid.UUID) (interface{}, bool) {
				equippedWeapon := scared.ProjectorEquippedWeapon.Project(selfID)
				ownerPosition := scared.ProjectorPosition.Project(equippedWeapon.OwnerID)
				ownerEntityType := scared.ProjectorEntityType.Project(equippedWeapon.OwnerID)

				targetIDs := scared.ProjectorEntityType.ListIdentifiers(
					func(t scared.EntityType) bool {
						for _, targetType := range scared.EntityTypeTargetMapping[ownerEntityType] {
							if targetType == t {
								return true
							}
						}
						return false
					},
				)

				var nearestTargetID uuid.UUID
				var nearestTargetDistance float64

				weaponRange := scared.ProjectorEquippedWeaponRange.Project(selfID)

				for _, targetID := range targetIDs {
					targetPosition := scared.ProjectorPosition.Project(targetID)

					distance := ownerPosition.DistanceOf(targetPosition)
					if (nearestTargetID == uuid.Nil || distance < nearestTargetDistance) && distance <= float64(weaponRange) {
						nearestTargetID = targetID
						nearestTargetDistance = distance
					}
				}

				if nearestTargetID == uuid.Nil {
					return nil, false
				}

				weaponDamage := scared.ProjectorEquippedWeaponDamage.Project(selfID)
				coolDown := scared.ProjectorEquippedWeaponCoolDown.Project(selfID)

				log.Printf("fire range: %d, damage: %d\n", weaponRange, weaponDamage.Plus-weaponDamage.Minus)
				return scared.WeaponLog{
					CoolDown: coolDown,
					Log: scared.Log[int]{
						ID:      uuid.New(),
						Minus:   weaponDamage.Minus,
						Plus:    weaponDamage.Plus,
						Targets: []uuid.UUID{nearestTargetID},
					},
				}, true
			},
		),
	},
	stateCoolDown: {
		EffectCoolDown: engine.NewGate(
			stateCoolDown,
			func(selfID uuid.UUID) (interface{}, bool) {
				coolDownCount := scared.ProjectorEquippedWeaponCoolDownCount.Project(selfID)
				if coolDownCount <= 0 {
					return nil, false
				}

				return nil, true
			},
		),
		EffectReActive: engine.NewGate(
			stateActive,
			func(selfID uuid.UUID) (interface{}, bool) {
				coolDownCount := scared.ProjectorEquippedWeaponCoolDownCount.Project(selfID)
				if coolDownCount > 0 {
					return nil, false
				}

				return nil, true
			},
		),
	},
})

const (
	stateRuneSlotEmpty     engine.State = "Empty"
	stateRuneSlotRequested engine.State = "Requested"
	stateRuneSlotActive    engine.State = "Active"
	stateRuneSlotDestroyed engine.State = "Destroyed"
)

const (
	EffectRuneSlotRequest engine.Effect = "Request"
	EffectRuneSlotActive  engine.Effect = "Active"
)

var LifeCycleStateMachineEquippedWeaponRuneSlot = engine.NewStateMachine(stateRuneSlotEmpty, engine.Nodes{
	stateRuneSlotEmpty: {
		EffectRuneSlotActive: engine.NewGate(
			stateRuneSlotActive,
			func(selfID uuid.UUID) (interface{}, bool) {
				equippedRuneIDs := scared.ProjectorEquippedRune.ListIdentifiers()
				for _, equippedRuneID := range equippedRuneIDs {
					equippedRune := scared.ProjectorEquippedRune.Project(equippedRuneID)
					if equippedRune.WeaponRuneSlotID == selfID {
						return equippedRuneID, true
					}
				}

				return nil, false
			},
		),
	},
})

var ExternalInputStateMachineEquippedWeaponRuneSlot = engine.NewStateMachine(stateRuneSlotEmpty, engine.Nodes{
	stateRuneSlotEmpty: {
		EffectRuneSlotRequest: engine.NewGate(
			stateRuneSlotRequested,
			func(selfID uuid.UUID) (interface{}, bool) {
				return selfID, true
			},
		),
	},
})
