package projectors

import (
	"thief/base/engine"
	"thief/scared/model"

	"github.com/google/uuid"
)

var (
	ProjectorEntityType engine.Projector[model.EntityType]

	ProjectorPosition engine.Projector[model.Position]

	ProjectorEquippedWeapon engine.Projector[model.EquippedWeapon]

	ProjectorEntityAssignedWeapon        engine.Projector[uuid.UUID] // weapon templateID
	ProjectorEquippedWeaponCoolDown      engine.Projector[int]
	ProjectorEquippedWeaponCoolDownCount engine.Projector[int]
	ProjectorEquippedWeaponRange         engine.Projector[int]
	ProjectorEquippedWeaponDamage        engine.Projector[model.Stat]
	ProjectorEquippedWeaponDevotionCost  engine.Projector[int]

	ProjectorEquippedWeaponRuneSlot                engine.Projector[model.WeaponRuneSlot]
	ProjectorRuneTemplateForEquippedWeaponRuneSlot engine.Projector[uuid.UUID]

	ProjectorWaypoint engine.Projector[model.Waypoint]

	ProjectorWeaponHitLogs engine.Projector[[]model.Log[model.Stat]]

	ProjectorHealth engine.Projector[int]

	ProjectorTarget engine.Projector[uuid.UUID]

	ProjectorDevotion                       engine.Projector[int]
	ProjectorDevotionGeneratorCoolDown      engine.Projector[int]
	ProjectorDevotionGeneratorCoolDownCount engine.Projector[int]

	ProjectorActiveWeapon engine.Projector[bool]

	ProjectorBlessingAltar engine.Projector[model.BlessingAltarData]
	ProjectorAcolyte       engine.Projector[int]

	ProjectorRunePlacement engine.Projector[model.RunePlacementData]
)
