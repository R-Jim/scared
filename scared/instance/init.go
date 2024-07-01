package instance

import (
	"thief/base/engine"
	"thief/scared"
	"thief/scared/entitytype"
	"thief/scared/health"
	"thief/scared/position"
	"thief/scared/scaredrune"
	"thief/scared/ship"
	"thief/scared/spawnership"
	"thief/scared/spawnersoul"
	"thief/scared/target"
	"thief/scared/weapon"

	"github.com/google/uuid"
)

var (
	StorePosition               = engine.NewStore()
	StoreEntityType             = engine.NewStore()
	StoreHealth                 = engine.NewStore()
	StoreTarget                 = engine.NewStore()
	StoreEquippedWeapon         = engine.NewStore()
	StoreEquippedWeaponRuneSlot = engine.NewStore()
	StoreEquippedRune           = engine.NewStore()
	StoreSpawnerShip            = engine.NewStore()
	StoreSpawnerSoul            = engine.NewStore()
	StoreShipGuard              = engine.NewStore()
)

func InitSpawner() map[scared.EntityType]uuid.UUID {
	spawnerShipID := uuid.New()

	StoreSpawnerShip.AppendEvent(engine.Event{
		ID:       uuid.New(),
		EntityID: spawnerShipID,
		Effect:   engine.EffectInit,
	})

	spawnerSoulID := uuid.New()

	StoreSpawnerSoul.AppendEvent(engine.Event{
		ID:       uuid.New(),
		EntityID: spawnerSoulID,
		Effect:   engine.EffectInit,
	})

	return map[scared.EntityType]uuid.UUID{
		scared.EntityTypeShip: spawnerShipID,
	}
}

func InitProjector() {
	scared.ProjectorEntityType = engine.NewStoreProjector[scared.EntityType](StoreEntityType,
		engine.NewFieldEffectMapping[scared.EntityType]([]engine.Effect{engine.EffectInit}, func(_ scared.EntityType, nextEventData interface{}) scared.EntityType {
			return nextEventData.(scared.EntityType)
		}),
	)

	scared.ProjectorPosition = engine.NewStoreProjector[scared.Position](StorePosition,
		engine.NewFieldEffectMapping[scared.Position]([]engine.Effect{engine.EffectInit}, func(_ scared.Position, nextEventData interface{}) scared.Position {
			bPosition := nextEventData.(scared.Position)

			return scared.Position{
				X: bPosition.X,
				Y: bPosition.Y,
			}
		}),
		engine.NewFieldEffectMapping[scared.Position]([]engine.Effect{position.EffectMove}, func(currentData scared.Position, nextEventData interface{}) scared.Position {
			bPosition := nextEventData.(scared.Position)

			return scared.Position{
				X: currentData.X + bPosition.X,
				Y: currentData.Y + bPosition.Y,
			}
		}),
	)

	scared.ProjectorEquippedWeapon = engine.NewStoreProjector[scared.EquippedWeapon](StoreEquippedWeapon,
		engine.NewFieldEffectMapping[scared.EquippedWeapon]([]engine.Effect{engine.EffectInit}, func(_ scared.EquippedWeapon, nextEventData interface{}) scared.EquippedWeapon {
			return nextEventData.(scared.EquippedWeapon)
		}),
	)

	scared.ProjectorEquippedWeaponCoolDown = engine.NewStoreProjector[int](StoreEquippedWeapon,
		engine.NewFieldEffectMapping[int]([]engine.Effect{engine.EffectInit}, func(currentData int, nextEventData interface{}) int {
			equippedWeapon := nextEventData.(scared.EquippedWeapon)
			template := scared.WeaponTemplates[equippedWeapon.TemplateID]

			weaponCoolDown := template.CoolDown

			for _, runeSlotID := range equippedWeapon.RuneSlotIDs {
				runeSlot := scared.ProjectorEquippedWeaponRuneSlot.Project(runeSlotID)
				if runeSlot.Type == scared.WeaponRuneSlotTypeCoolDown && runeSlot.RuneID != uuid.Nil {
					rune := scared.ProjectorEquippedRune.Project(runeSlot.RuneID)
					weaponCoolDown = weaponCoolDown / scared.RuneTemplates[rune.TemplateID].Modifier.Minus
				}
			}

			return weaponCoolDown
		}),
	)

	scared.ProjectorEquippedWeaponCoolDownCount = engine.NewStoreProjector[int](StoreEquippedWeapon,
		engine.NewFieldEffectMapping[int]([]engine.Effect{weapon.EffectHitEnemy}, func(currentData int, nextEventData interface{}) int {
			weaponLog := nextEventData.(scared.WeaponLog)

			return currentData + weaponLog.CoolDown
		}),
		engine.NewFieldEffectMapping[int]([]engine.Effect{weapon.EffectCoolDown}, func(currentData int, nextEventData interface{}) int {
			return currentData - 1
		}),
	)

	scared.ProjectorEquippedWeaponRange = engine.NewStoreProjector[int](StoreEquippedWeapon,
		engine.NewFieldEffectMapping[int]([]engine.Effect{engine.EffectInit}, func(currentData int, nextEventData interface{}) int {
			equippedWeapon := nextEventData.(scared.EquippedWeapon)
			template := scared.WeaponTemplates[equippedWeapon.TemplateID]

			weaponRange := template.Range

			for _, runeSlotID := range equippedWeapon.RuneSlotIDs {
				runeSlot := scared.ProjectorEquippedWeaponRuneSlot.Project(runeSlotID)
				if runeSlot.Type == scared.WeaponRuneSlotTypeRange && runeSlot.RuneID != uuid.Nil {
					rune := scared.ProjectorEquippedRune.Project(runeSlot.RuneID)
					weaponRange = weaponRange * scared.RuneTemplates[rune.TemplateID].Modifier.Plus
				}
			}

			return weaponRange
		}),
	)

	scared.ProjectorEquippedWeaponDamage = engine.NewStoreProjector[scared.Stat](StoreEquippedWeapon,
		engine.NewFieldEffectMapping[scared.Stat]([]engine.Effect{engine.EffectInit}, func(currentData scared.Stat, nextEventData interface{}) scared.Stat {
			equippedWeapon := nextEventData.(scared.EquippedWeapon)
			template := scared.WeaponTemplates[equippedWeapon.TemplateID]

			damage := template.Damage

			for _, runeSlotID := range equippedWeapon.RuneSlotIDs {
				runeSlot := scared.ProjectorEquippedWeaponRuneSlot.Project(runeSlotID)
				if runeSlot.Type == scared.WeaponRuneSlotTypeDamage && runeSlot.RuneID != uuid.Nil {
					rune := scared.ProjectorEquippedRune.Project(runeSlot.RuneID)
					damage.Minus = damage.Minus * scared.RuneTemplates[rune.TemplateID].Modifier.Minus
					damage.Plus = damage.Plus * scared.RuneTemplates[rune.TemplateID].Modifier.Plus
				}
			}

			return damage
		}),
	)

	scared.ProjectorEquippedWeaponRuneSlot = engine.NewStoreProjector[scared.WeaponRuneSlot](StoreEquippedWeaponRuneSlot,
		engine.NewFieldEffectMapping[scared.WeaponRuneSlot]([]engine.Effect{engine.EffectInit}, func(_ scared.WeaponRuneSlot, nextEventData interface{}) scared.WeaponRuneSlot {
			return scared.WeaponRuneSlot{
				Type: nextEventData.(scared.WeaponRuneSlotType),
			}
		}),
		engine.NewFieldEffectMapping[scared.WeaponRuneSlot]([]engine.Effect{weapon.EffectRuneSlotActive}, func(currentData scared.WeaponRuneSlot, nextEventData interface{}) scared.WeaponRuneSlot {
			currentData.RuneID = nextEventData.(uuid.UUID)
			return currentData
		}),
	)

	scared.ProjectorEquippedRune = engine.NewStoreProjector[scared.EquippedRune](StoreEquippedRune,
		engine.NewFieldEffectMapping[scared.EquippedRune]([]engine.Effect{engine.EffectInit}, func(_ scared.EquippedRune, nextEventData interface{}) scared.EquippedRune {
			return nextEventData.(scared.EquippedRune)
		}),
	)

	scared.ProjectorWeaponHitLogs = engine.NewStoreProjector[[]scared.Log[int]](StoreEquippedWeapon,
		engine.NewFieldEffectMapping[[]scared.Log[int]]([]engine.Effect{weapon.EffectHitEnemy}, func(currentData []scared.Log[int], nextEventData interface{}) []scared.Log[int] {
			log := nextEventData.(scared.WeaponLog)
			return append(currentData, log.Log)
		}),
	)

	scared.ProjectorHealth = engine.NewStoreProjector[int](StoreHealth,
		engine.NewFieldEffectMapping[int]([]engine.Effect{health.EffectHit}, func(currentData int, nextEventData interface{}) int {
			logs := nextEventData.([]scared.Log[int])

			minus := 0
			plus := 0

			for _, log := range logs {
				minus += log.Minus
				plus += log.Plus
			}

			return currentData + plus - minus
		}),
		engine.NewFieldEffectMapping[int]([]engine.Effect{engine.EffectInit}, func(currentData int, nextEventData interface{}) int {
			return nextEventData.(int)
		}),
	)

	scared.ProjectorConsumedWeaponHitLogs = engine.NewStoreProjector[[]uuid.UUID](StoreHealth,
		engine.NewFieldEffectMapping[[]uuid.UUID]([]engine.Effect{health.EffectHit}, func(currentData []uuid.UUID, nextEventData interface{}) []uuid.UUID {
			logs := nextEventData.([]scared.Log[int])
			ids := []uuid.UUID{}
			for _, log := range logs {
				ids = append(ids, log.ID)
			}

			return append(currentData, ids...)
		}),
	)

	scared.ProjectorTarget = engine.NewStoreProjector[uuid.UUID](StoreTarget,
		engine.NewFieldEffectMapping[uuid.UUID]([]engine.Effect{target.EffectSelectTarget}, func(_ uuid.UUID, nextEventData interface{}) uuid.UUID {
			return nextEventData.(uuid.UUID)
		}),
	)
}

func InitComposerLifeCycle() []*engine.ComposerLifeCycle {
	return []*engine.ComposerLifeCycle{
		engine.NewComposerLifeCycle(StorePosition, position.StateMachine),
		engine.NewComposerLifeCycle(StoreTarget, target.StateMachine),
		engine.NewComposerLifeCycle(StoreEquippedWeapon, weapon.StateMachineEquippedWeapon),
		engine.NewComposerLifeCycle(StoreEquippedWeaponRuneSlot, weapon.LifeCycleStateMachineEquippedWeaponRuneSlot),
		engine.NewComposerLifeCycle(StoreHealth, health.StateMachine),
		engine.NewComposerLifeCycle(StoreEntityType, entitytype.StateMachine),
		engine.NewComposerLifeCycle(StoreSpawnerSoul, spawnersoul.StateMachine),
	}
}

func InitComposerSpawner() []*engine.ComposerSpawner {
	return []*engine.ComposerSpawner{
		engine.NewComposerSpawner(StoreSpawnerShip, spawnership.EffectSpawn, spawnership.NewSpawner(StorePosition, StoreEntityType, StoreHealth, StoreShipGuard)),
		engine.NewComposerSpawner(StoreSpawnerSoul, spawnersoul.EffectSpawn, spawnersoul.NewSpawner(StorePosition, StoreEntityType, StoreTarget, StoreHealth)),
		engine.NewComposerSpawner(StoreShipGuard, ship.EffectArm, weapon.NewSpawnerEquippedWeapon(StoreEquippedWeapon, StoreEquippedWeaponRuneSlot)),
		engine.NewComposerSpawner(StoreEquippedWeaponRuneSlot, weapon.EffectRuneSlotRequest, scaredrune.NewSpawnerEquippedRune(StoreEquippedRune)),
	}
}

func InitComposerDestroyer() []*engine.ComposerDestroyer {
	return []*engine.ComposerDestroyer{
		engine.NewComposerDestroyer(StoreHealth, health.StateMachine),
		engine.NewComposerDestroyer(StoreTarget, target.StateMachine),
		engine.NewComposerDestroyer(StorePosition, position.StateMachine),
		engine.NewComposerDestroyer(StoreEntityType, entitytype.StateMachine),
	}
}

const (
	SpawnShip   = "SpawnShip"
	EquipWeapon = "EquipWeapon"
	EquipRune   = "EquipRune"
)

func InitComposerExternalInput() map[string]*engine.ComposerExternalInput {
	return map[string]*engine.ComposerExternalInput{
		SpawnShip:   engine.NewComposerExternalInput(StoreSpawnerShip, spawnership.StateMachine),
		EquipWeapon: engine.NewComposerExternalInput(StoreShipGuard, ship.StateMachine),
		EquipRune:   engine.NewComposerExternalInput(StoreEquippedWeaponRuneSlot, weapon.ExternalInputStateMachineEquippedWeaponRuneSlot),
	}
}
