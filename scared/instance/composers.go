package instance

import (
	"thief/base/engine"
	"thief/scared/statemachines/acolyte"
	"thief/scared/statemachines/church"
	"thief/scared/statemachines/devotion"
	"thief/scared/statemachines/entitytype"
	"thief/scared/statemachines/health"
	"thief/scared/statemachines/knight"
	"thief/scared/statemachines/position"
	"thief/scared/statemachines/runeplacement"
	"thief/scared/statemachines/ship"
	"thief/scared/statemachines/soul"
	"thief/scared/statemachines/spawnership"
	"thief/scared/statemachines/spawnersoul"
	"thief/scared/statemachines/target"
	"thief/scared/statemachines/weapon"
)

var (
	composerLifeCyclePosition               = engine.NewComposerLifeCycle(storePosition, position.StateMachine)
	composerLifeCycleTarget                 = engine.NewComposerLifeCycle(storeTarget, target.StateMachine)
	composerLifeCycleEquippedWeapon         = engine.NewComposerLifeCycle(storeEquippedWeapon, weapon.StateMachineEquippedWeapon)
	composerLifeCycleEquippedWeaponRuneSlot = engine.NewComposerLifeCycle(storeEquippedWeaponRuneSlot, weapon.StateMachineEquippedWeaponRuneSlot)
	composerLifeCycleHealth                 = engine.NewComposerLifeCycle(storeHealth, health.StateMachine)
	composerLifeCycleEntityType             = engine.NewComposerLifeCycle(storeEntityType, entitytype.StateMachine)
	composerLifeCycleSpawnerShip            = engine.NewComposerLifeCycle(storeSpawnerShip, spawnership.StateMachine)
	composerLifeCycleSpawnerSoul            = engine.NewComposerLifeCycle(storeRunePlacement, runeplacement.StateMachine)
	composerLifeCycleSoul                   = engine.NewComposerLifeCycle(storeSoul, soul.StateMachine)
	composerLifeCycleShipGuard              = engine.NewComposerLifeCycle(storeShipGuard, ship.StateMachineShipGuard)
	composerLifeCycleDevotionGenerator      = engine.NewComposerLifeCycle(storeDevotionGenerator, devotion.StateMachineDevotionGenerator)
	composerLifeCycleShipBlessingAltar      = engine.NewComposerLifeCycle(storeShipBlessingAltar, ship.StateMachineShipBlessingAltar)
	composerLifeCycleChurch                 = engine.NewComposerLifeCycle(storeChurch, church.StateMachine)
	composerLifeCycleAcolyte                = engine.NewComposerLifeCycle(storeAcolyte, acolyte.StateMachine)

	composerAfterEffectShip                            = engine.NewComposerAfterEffect(storeSpawnerShip, spawnership.EffectSpawn, spawnership.NewShip(storePosition, storeEntityType, storeHealth, storeShipGuard, storeShipBlessingAltar, storeAcolyte))
	composerAfterEffectShipGuardArmWeapon              = engine.NewComposerAfterEffect(storeShipGuard, ship.EffectArm, weapon.NewEquippedWeapon(storeEquippedWeapon, storeEquippedWeaponRuneSlot))
	composerAfterEffectSoul                            = engine.NewComposerAfterEffect(storeRunePlacement, runeplacement.EffectSpawnGuardian, spawnersoul.NewSoul(storePosition, storeEntityType, storeTarget, storeHealth, storeSoul))
	composerAfterEffectSoulArmWeapon                   = engine.NewComposerAfterEffect(storeSoul, soul.EffectArm, weapon.NewEquippedWeapon(storeEquippedWeapon, storeEquippedWeaponRuneSlot))
	composerAfterEffectDevotionGenerator               = engine.NewComposerAfterEffect(storeDevotionGenerator, devotion.EffectGeneratorRun, devotion.NewAfterEffectDevotionGenerator(storeDevotion))
	composerAfterEffectDevotionWeaponConsumer          = engine.NewComposerAfterEffect(storeEquippedWeapon, weapon.EffectHitEnemy, devotion.NewAfterEffectDevotionWeaponConsumer(storeDevotion))
	composerAfterEffectKnight                          = engine.NewComposerAfterEffect(storeShipBlessingAltar, ship.EffectBlessAcolyteToKnight, knight.NewKnight(storePosition, storeEntityType, storeTarget, storeHealth, storeKnight))
	composerAfterEffectHealthWeaponLogConsumer         = engine.NewComposerAfterEffect(storeEquippedWeapon, weapon.EffectHitEnemy, health.NewWeaponHitLogConsumer(storeHealth))
	composerAfterEffectHealthWeaponHitDevotionConsumer = engine.NewComposerAfterEffect(storeHealth, health.EffectHitDevotion, devotion.NewAfterEffectDevotionHealthConsumer(storeDevotion))
	composerAfterEffectSpawnChurch                     = engine.NewComposerAfterEffect(storeChurch, church.EffectInit, church.NewChurch(storePosition, storeEntityType, storeHealth, storeAcolyte))
	composerAfterEffectAcolyteChurchTransfer           = engine.NewComposerAfterEffect(storeChurch, church.EffectCollect, acolyte.NewAfterEffectAcolyteTransfer(storeAcolyte))
)

func OperateConsumers() {
	composerLifeCyclePosition.Operate()
	composerLifeCycleTarget.Operate()
	composerLifeCycleEquippedWeapon.Operate()
	composerLifeCycleEquippedWeaponRuneSlot.Operate()
	composerLifeCycleHealth.Operate()
	composerLifeCycleEntityType.Operate()
	composerLifeCycleSpawnerShip.Operate()
	composerLifeCycleSpawnerSoul.Operate()
	composerLifeCycleSoul.Operate()
	composerLifeCycleShipGuard.Operate()
	composerLifeCycleDevotionGenerator.Operate()
	composerLifeCycleShipBlessingAltar.Operate()
	composerLifeCycleChurch.Operate()
	composerLifeCycleAcolyte.Operate()

	composerAfterEffectShip.Operate()
	composerAfterEffectSoul.Operate()
	composerAfterEffectShipGuardArmWeapon.Operate()
	composerAfterEffectSoulArmWeapon.Operate()
	composerAfterEffectDevotionGenerator.Operate()
	composerAfterEffectDevotionWeaponConsumer.Operate()
	composerAfterEffectKnight.Operate()
	composerAfterEffectHealthWeaponLogConsumer.Operate()
	composerAfterEffectHealthWeaponHitDevotionConsumer.Operate()
	composerAfterEffectSpawnChurch.Operate()
	composerAfterEffectAcolyteChurchTransfer.Operate()

	composerLifeCyclePosition.CommitDestroyedIDs()
	composerLifeCycleTarget.CommitDestroyedIDs()
	composerLifeCycleEquippedWeapon.CommitDestroyedIDs()
	composerLifeCycleEquippedWeaponRuneSlot.CommitDestroyedIDs()
	composerLifeCycleHealth.CommitDestroyedIDs()
	composerLifeCycleEntityType.CommitDestroyedIDs()
	composerLifeCycleSpawnerShip.CommitDestroyedIDs()
	composerLifeCycleSpawnerSoul.CommitDestroyedIDs()
	composerLifeCycleSoul.CommitDestroyedIDs()
	composerLifeCycleShipGuard.CommitDestroyedIDs()
	composerLifeCycleDevotionGenerator.CommitDestroyedIDs()
	composerLifeCycleShipBlessingAltar.CommitDestroyedIDs()
	composerLifeCycleChurch.CommitDestroyedIDs()
	composerLifeCycleAcolyte.CommitDestroyedIDs()
}
