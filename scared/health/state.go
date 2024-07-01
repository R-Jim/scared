package health

import (
	"thief/base/engine"
	"thief/scared"

	"github.com/google/uuid"
)

const (
	stateActive    engine.State = "Active"
	stateDestroyed engine.State = engine.StateDestroyed
)

const (
	EffectHit     engine.Effect = "Hit"
	EffectDestroy engine.Effect = "Destroy"
)

const (
	MaxMoveRange = 1
)

var StateMachine = engine.NewStateMachine(stateActive, engine.Nodes{
	stateActive: {
		EffectHit: engine.NewGate(
			stateActive,
			func(selfID uuid.UUID) (interface{}, bool) {
				selfTargetedEquippedWeaponIDs := scared.ProjectorEquippedWeapon.ListIdentifiers()
				if len(selfTargetedEquippedWeaponIDs) == 0 {
					return nil, false
				}

				consumedWeaponActiveLogs := scared.ProjectorConsumedWeaponHitLogs.Project(selfID)
				unConsumedLogs := []scared.Log[int]{}

				for _, selfTargetedEquippedWeaponID := range selfTargetedEquippedWeaponIDs {
					weaponHitLogs := scared.ProjectorWeaponHitLogs.Project(selfTargetedEquippedWeaponID)

					for _, log := range weaponHitLogs {
						isTargeted := false
						for _, targetID := range log.Targets {
							if targetID == selfID {
								isTargeted = true
								break
							}
						}

						if !isTargeted {
							continue
						}

						isLogConsumed := false
						for _, consumedLogID := range consumedWeaponActiveLogs {
							if log.ID == consumedLogID {
								isLogConsumed = true
								break
							}
						}

						if !isLogConsumed {
							unConsumedLogs = append(unConsumedLogs, log)
						}
					}
				}

				if len(unConsumedLogs) <= 0 {
					return nil, false
				}

				return unConsumedLogs, true
			},
		),
		EffectDestroy: engine.NewGate(
			stateDestroyed,
			func(selfID uuid.UUID) (interface{}, bool) {
				health := scared.ProjectorHealth.Project(selfID)
				if health > 0 {
					return nil, false
				}

				return nil, true
			},
		),
	},
})
