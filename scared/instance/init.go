package instance

import (
	"thief/base/engine"
	"thief/scared"
	"thief/scared/animator"
	"thief/scared/model"
	"thief/scared/projectors"
	"thief/scared/statemachines/acolyte"
	"thief/scared/statemachines/church"
	"thief/scared/statemachines/devotion"
	"thief/scared/statemachines/entitytype"
	"thief/scared/statemachines/health"
	"thief/scared/statemachines/position"
	"thief/scared/statemachines/runeplacement"
	"thief/scared/statemachines/ship"
	"thief/scared/statemachines/spawnership"
	"thief/scared/statemachines/target"
	"thief/scared/statemachines/weapon"

	"github.com/google/uuid"
)

var (
	storePosition               = engine.NewStore("Position")
	storeEntityType             = engine.NewStore("EntityType")
	storeHealth                 = engine.NewStore("Health")
	storeTarget                 = engine.NewStore("Target")
	storeEquippedWeapon         = engine.NewStore("EquippedWeapon")
	storeEquippedWeaponRuneSlot = engine.NewStore("EquippedWeaponRuneSlot")

	storeSpawnerShip = engine.NewStore("SpawnerShip")
	storeShipGuard   = engine.NewStore("ShipGuard")

	storeRunePlacement = engine.NewStore("RunePlacement")
	storeSoul          = engine.NewStore("Soul")

	storeDevotion          = engine.NewStore("Devotion")
	storeDevotionGenerator = engine.NewStore("DevotionGenerator")

	storeShipBlessingAltar = engine.NewStore("ShipBlessingAltar")

	storeKnight = engine.NewStore("Knight")

	storeAcolyte = engine.NewStore("Acolyte")

	storeChurch = engine.NewStore("Church")
)

func InitEntities() {
	scared.DevotionID = uuid.New()

	storeDevotion.AppendEvent(devotion.EffectInit.NewEvent(scared.DevotionID, 100))

	storeDevotionGenerator.AppendEvent(devotion.EffectGeneratorInit.NewEvent(scared.DevotionID, 30))

	storeSpawnerShip.AppendEvent(spawnership.EffectInit.NewEvent(uuid.New(), nil))

	storeRunePlacement.AppendEvent(runeplacement.EffectInit.NewEvent(uuid.New(), model.SpawnRunePlacementData{
		Position: model.PositionRunePlacement,
	}))

	storeChurch.AppendEvent(church.EffectInit.NewEvent(uuid.New(), model.Position{X: 100, Y: 50}))
}

func InitProjector() {
	projectors.ProjectorEntityType = engine.NewStoreProjector(storeEntityType,
		engine.NewFieldEffectMapping([]engine.Effect[model.EntityType]{entitytype.EffectInit}, func(_ model.EntityType, nextEventData model.EntityType) model.EntityType {
			return nextEventData
		}),
	)

	projectors.ProjectorPosition = engine.NewStoreProjector(storePosition,
		engine.NewFieldEffectMapping([]engine.Effect[model.Position]{position.EffectInit}, func(_ model.Position, nextEventData model.Position) model.Position {
			return model.Position{
				X: nextEventData.X,
				Y: nextEventData.Y,
			}
		}),
		engine.NewFieldEffectMapping([]engine.Effect[model.Position]{position.EffectMove}, func(currentData model.Position, nextEventData model.Position) model.Position {
			return model.Position{
				X: currentData.X + nextEventData.X,
				Y: currentData.Y + nextEventData.Y,
			}
		}),
	)

	projectors.ProjectorEquippedWeapon = engine.NewStoreProjector(storeEquippedWeapon,
		engine.NewFieldEffectMapping([]engine.Effect[model.EquippedWeapon]{weapon.EffectInit}, func(_ model.EquippedWeapon, nextEventData model.EquippedWeapon) model.EquippedWeapon {
			return nextEventData
		}),
	)

	projectors.ProjectorEquippedWeaponCoolDown = engine.NewStoreProjector(storeEquippedWeapon,
		engine.NewFieldEffectMapping([]engine.Effect[model.EquippedWeapon]{weapon.EffectInit}, func(currentData int, nextEventData model.EquippedWeapon) int {
			template := model.TemplateWeapons[nextEventData.TemplateID]

			weaponCoolDown := template.CoolDown

			for _, runeSlotID := range nextEventData.RuneSlotIDs {
				runeSlot := projectors.ProjectorEquippedWeaponRuneSlot.Project(runeSlotID)
				if runeSlot.Type == model.WeaponRuneSlotTypeCoolDown && runeSlot.RuneID != uuid.Nil {
					weaponCoolDown = weaponCoolDown / model.RuneTemplates[runeSlot.RuneID].Modifier.Minus
				}
			}

			return weaponCoolDown
		}),
	)

	projectors.ProjectorEquippedWeaponCoolDownCount = engine.NewStoreProjector(storeEquippedWeapon,
		engine.NewFieldEffectMapping([]engine.Effect[model.WeaponLog]{weapon.EffectHitEnemy}, func(currentData int, nextEventData model.WeaponLog) int {
			return currentData + nextEventData.CoolDown
		}),
		engine.NewFieldEffectMapping([]engine.Effect[int]{weapon.EffectCoolDown}, func(currentData int, nextEventData int) int {
			return currentData - nextEventData
		}),
	)

	projectors.ProjectorEquippedWeaponRange = engine.NewStoreProjector(storeEquippedWeapon,
		engine.NewFieldEffectMapping([]engine.Effect[model.EquippedWeapon]{weapon.EffectInit}, func(currentData int, nextEventData model.EquippedWeapon) int {
			template := model.TemplateWeapons[nextEventData.TemplateID]

			weaponRange := template.Range

			for _, runeSlotID := range nextEventData.RuneSlotIDs {
				runeSlot := projectors.ProjectorEquippedWeaponRuneSlot.Project(runeSlotID)
				if runeSlot.Type == model.WeaponRuneSlotTypeRange && runeSlot.RuneID != uuid.Nil {
					weaponRange = weaponRange * model.RuneTemplates[runeSlot.RuneID].Modifier.Plus
				}
			}

			return weaponRange
		}),
	)

	projectors.ProjectorEquippedWeaponDamage = engine.NewStoreProjector(storeEquippedWeapon,
		engine.NewFieldEffectMapping([]engine.Effect[model.EquippedWeapon]{weapon.EffectInit}, func(currentData model.Stat, nextEventData model.EquippedWeapon) model.Stat {
			template := model.TemplateWeapons[nextEventData.TemplateID]

			damage := template.Damage

			for _, runeSlotID := range nextEventData.RuneSlotIDs {
				runeSlot := projectors.ProjectorEquippedWeaponRuneSlot.Project(runeSlotID)
				if runeSlot.Type == model.WeaponRuneSlotTypeDamage && runeSlot.RuneID != uuid.Nil {
					damage.Minus = damage.Minus * model.RuneTemplates[runeSlot.RuneID].Modifier.Minus
					damage.Plus = damage.Plus * model.RuneTemplates[runeSlot.RuneID].Modifier.Plus
				}
			}

			return damage
		}),
	)

	projectors.ProjectorEquippedWeaponDevotionCost = engine.NewStoreProjector(storeEquippedWeapon,
		engine.NewFieldEffectMapping([]engine.Effect[model.EquippedWeapon]{weapon.EffectInit}, func(currentData int, nextEventData model.EquippedWeapon) int {
			var cost int

			for _, runeSlotID := range nextEventData.RuneSlotIDs {
				runeSlot := projectors.ProjectorEquippedWeaponRuneSlot.Project(runeSlotID)
				if runeSlot.Type == model.WeaponRuneSlotTypeDamage && runeSlot.RuneID != uuid.Nil {
					cost += model.RuneTemplates[runeSlot.RuneID].DevotionCost
				}
			}

			return cost
		}),
	)

	projectors.ProjectorEquippedWeaponRuneSlot = engine.NewStoreProjector(storeEquippedWeaponRuneSlot,
		engine.NewFieldEffectMapping([]engine.Effect[model.WeaponRuneSlotType]{weapon.EffectRuneSlotInit}, func(_ model.WeaponRuneSlot, nextEventData model.WeaponRuneSlotType) model.WeaponRuneSlot {
			return model.WeaponRuneSlot{
				Type: nextEventData,
			}
		}),
		engine.NewFieldEffectMapping([]engine.Effect[uuid.UUID]{weapon.EffectRuneSlotActive}, func(currentData model.WeaponRuneSlot, nextEventData uuid.UUID) model.WeaponRuneSlot {
			currentData.RuneID = nextEventData
			return currentData
		}),
	)

	projectorRuneToEquippedWeaponRuneSlot := newProjectorRuneToEquippedWeaponRuneSlot()
	SetRuneToEquippedWeaponRuneSlotFunc = projectorRuneToEquippedWeaponRuneSlot.AddRuneToRuneSlot
	projectors.ProjectorRuneTemplateForEquippedWeaponRuneSlot = projectorRuneToEquippedWeaponRuneSlot

	projectors.ProjectorWeaponHitLogs = engine.NewStoreProjector(storeEquippedWeapon,
		engine.NewFieldEffectMapping([]engine.Effect[model.WeaponLog]{weapon.EffectHitEnemy}, func(currentData []model.Log[model.Stat], nextEventData model.WeaponLog) []model.Log[model.Stat] {
			return append(currentData, nextEventData.Log)
		}),
	)

	projectors.ProjectorHealth = engine.NewStoreProjector(storeHealth,
		engine.NewFieldEffectMapping([]engine.Effect[int]{health.EffectHit}, func(currentData int, nextEventData int) int {
			return currentData + nextEventData
		}),
		engine.NewFieldEffectMapping([]engine.Effect[int]{health.EffectInit}, func(currentData int, nextEventData int) int {
			return nextEventData
		}),
	)

	projectors.ProjectorTarget = engine.NewStoreProjector(storeTarget,
		engine.NewFieldEffectMapping([]engine.Effect[uuid.UUID]{target.EffectSelectTarget, target.EffectReleaseTarget}, func(_ uuid.UUID, nextEventData uuid.UUID) uuid.UUID {
			return nextEventData
		}),
	)

	projectorWaypoint := newProjectorWaypoint()
	SetWaypointFunc = projectorWaypoint.SetWaypoint
	projectors.ProjectorWaypoint = projectorWaypoint

	projectorEntityAssignedWeapon := newProjectorEntityAssignedWeapon()
	AssignWeaponToEntityFunc = projectorEntityAssignedWeapon.AssignWeaponToEntity
	projectors.ProjectorEntityAssignedWeapon = projectorEntityAssignedWeapon

	projectors.ProjectorDevotion = engine.NewStoreProjector(storeDevotion,
		engine.NewFieldEffectMapping([]engine.Effect[int]{devotion.EffectInit}, func(currentData int, nextEffectData int) int {
			return nextEffectData
		}),
		engine.NewFieldEffectMapping([]engine.Effect[int]{devotion.EffectConsume}, func(currentData int, nextEffectData int) int {
			return currentData - nextEffectData
		}),
		engine.NewFieldEffectMapping([]engine.Effect[int]{devotion.EffectAdd}, func(currentData int, nextEffectData int) int {
			return currentData + nextEffectData
		}),
	)

	projectors.ProjectorDevotionGeneratorCoolDown = engine.NewStoreProjector(storeDevotionGenerator,
		engine.NewFieldEffectMapping([]engine.Effect[int]{devotion.EffectGeneratorInit}, func(currentData int, nextEffectData int) int {
			return nextEffectData
		}),
	)

	projectors.ProjectorDevotionGeneratorCoolDownCount = engine.NewStoreProjector(storeDevotionGenerator,
		engine.NewFieldEffectMapping([]engine.Effect[int]{devotion.EffectGeneratorRun}, func(currentData int, nextEventData int) int {
			return currentData + nextEventData
		}),
		engine.NewFieldEffectMapping([]engine.Effect[int]{devotion.EffectGeneratorCoolDown}, func(currentData int, nextEventData int) int {
			return currentData - 1
		}),
	)

	projectorActiveWeapon := newProjectorActiveWeapon()
	SetActiveWeapon = projectorActiveWeapon.SetActiveWeapon
	projectors.ProjectorActiveWeapon = projectorActiveWeapon

	// SetAcolyte = projectorAcolyte.SetAcolyte
	// TransferAcolyte = projectorAcolyte.TransformerAcolyte
	projectors.ProjectorAcolyte = engine.NewStoreProjector(storeAcolyte,
		engine.NewFieldEffectMapping([]engine.Effect[int]{acolyte.EffectInit}, func(currentData int, nextEventData int) int {
			return nextEventData
		}),
		engine.NewFieldEffectMapping([]engine.Effect[int]{acolyte.EffectDeposit}, func(currentData int, nextEventData int) int {
			return currentData + nextEventData
		}),
		engine.NewFieldEffectMapping([]engine.Effect[int]{acolyte.EffectWithdraw}, func(currentData int, nextEventData int) int {
			return currentData - nextEventData
		}),
	)

	projectors.ProjectorBlessingAltar = engine.NewStoreProjector(storeShipBlessingAltar,
		engine.NewFieldEffectMapping([]engine.Effect[uuid.UUID]{ship.EffectBlessingAltarInit}, func(currentData model.BlessingAltarData, nextEventData uuid.UUID) model.BlessingAltarData {
			return model.BlessingAltarData{
				OwnerID: nextEventData,
			}
		}),
		engine.NewFieldEffectMapping([]engine.Effect[model.SpawnKnightData]{ship.EffectBlessAcolyteToKnight}, func(currentData model.BlessingAltarData, nextEventData model.SpawnKnightData) model.BlessingAltarData {
			currentData.NumberOfBlessedAcolyte += 1
			return currentData
		}),
	)

	projectors.ProjectorRunePlacement = engine.NewStoreProjector(storeRunePlacement,
		engine.NewFieldEffectMapping([]engine.Effect[model.SpawnRunePlacementData]{runeplacement.EffectInit}, func(currentData model.RunePlacementData, nextEventData model.SpawnRunePlacementData) model.RunePlacementData {
			return model.RunePlacementData{
				Position: nextEventData.Position,
			}
		}),
		engine.NewFieldEffectMapping([]engine.Effect[model.SpawnSoulData]{runeplacement.EffectSpawnGuardian}, func(currentData model.RunePlacementData, nextEventData model.SpawnSoulData) model.RunePlacementData {
			currentData.SpawnedSoulIDs = append(currentData.SpawnedSoulIDs, nextEventData.ID)
			return currentData
		}),
	)
}

const (
	SpawnShip      = "SpawnShip"
	ChangeWayPoint = "ChangeWayPoint"
	EquipWeapon    = "EquipWeapon"
	EquipRune      = "EquipRune"
)

func InitAnimators() []engine.Animator {
	projectileAnimator := animator.NewProjectileAnimator()
	storeEquippedWeapon.AddHook(projectileAnimator.GetHook())

	soulAnimator := animator.NewSoulAnimator()
	shipAnimator := animator.NewShipAnimator()
	knightAnimator := animator.NewKnightAnimator()
	churchAnimator := animator.NewChurchAnimator()

	storeEntityType.AddHook(soulAnimator.GetHook())
	storeEntityType.AddHook(shipAnimator.GetHook())
	storeEntityType.AddHook(knightAnimator.GetHook())

	storeChurch.AddHook(churchAnimator.GetHook())

	runePlacementAnimator := animator.NewRunePlacementAnimator()
	storeRunePlacement.AddHook(runePlacementAnimator.GetHook())

	return []engine.Animator{
		projectileAnimator,
		soulAnimator,
		shipAnimator,
		knightAnimator,
		runePlacementAnimator,
		churchAnimator,
	}
}
