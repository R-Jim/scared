package scared

import (
	"thief/base/engine"

	"github.com/google/uuid"
)

var (
	ProjectorEntityType engine.Projector[EntityType]

	ProjectorPosition engine.Projector[Position]

	ProjectorEquippedWeapon engine.Projector[EquippedWeapon]

	ProjectorEquippedWeaponCoolDown      engine.Projector[int]
	ProjectorEquippedWeaponCoolDownCount engine.Projector[int]
	ProjectorEquippedWeaponRange         engine.Projector[int]
	ProjectorEquippedWeaponDamage        engine.Projector[Stat]

	ProjectorEquippedWeaponRuneSlot engine.Projector[WeaponRuneSlot]
	ProjectorEquippedRune           engine.Projector[EquippedRune]

	ProjectorWeaponHitLogs engine.Projector[[]Log[int]]

	ProjectorHealth engine.Projector[int]

	ProjectorConsumedWeaponHitLogs engine.Projector[[]uuid.UUID]

	ProjectorTarget engine.Projector[uuid.UUID]
)
