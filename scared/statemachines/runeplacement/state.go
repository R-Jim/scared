package runeplacement

import (
	"thief/base/engine"
	"thief/scared/model"
	"thief/scared/projectors"

	"github.com/google/uuid"
)

const (
	stateActive    engine.State = "Active"
	stateCollected engine.State = engine.StateDestroyed

	EffectInit          engine.Effect[model.SpawnRunePlacementData] = "Init"
	EffectSpawnGuardian engine.Effect[model.SpawnSoulData]          = "SpawnGuardian"
	EffectCollected     engine.Effect[uuid.UUID]                    = "Collected"
)

var StateMachine = engine.NewStateMachine(EffectInit.ToState(stateActive), engine.Nodes{
	stateActive: {
		EffectCollected.ToStateWhen(
			stateCollected,
			func(selfID uuid.UUID) (uuid.UUID, bool) {
				runePlacementData := projectors.ProjectorRunePlacement.Project(selfID)

				if len(runePlacementData.SpawnedSoulIDs) <= 0 {
					return uuid.Nil, false
				}

				soulIDs := projectors.ProjectorEntityType.ListIdentifiers(func(et model.EntityType) bool {
					return et == model.EntityTypeSoul
				})

				remainingRunePlacementSpawnedSoulIDs := []uuid.UUID{}
				for _, soulID := range soulIDs {
					for _, runePlacementSpawnedSoulID := range runePlacementData.SpawnedSoulIDs {
						if soulID == runePlacementSpawnedSoulID {
							remainingRunePlacementSpawnedSoulIDs = append(remainingRunePlacementSpawnedSoulIDs, soulID)
						}
					}
				}

				return uuid.Nil, len(remainingRunePlacementSpawnedSoulIDs) <= 0
			},
		),
		EffectSpawnGuardian.ToStateWhen(
			stateActive,
			func(selfID uuid.UUID) (model.SpawnSoulData, bool) {
				// will try to cover the area with a number of souls guardian
				runePlacementData := projectors.ProjectorRunePlacement.Project(selfID)

				soulIDs := projectors.ProjectorEntityType.ListIdentifiers(func(et model.EntityType) bool {
					return et == model.EntityTypeSoul
				})

				remainingRunePlacementSpawnedSoulIDs := []uuid.UUID{}
				for _, soulID := range soulIDs {
					for _, runePlacementSpawnedSoulID := range runePlacementData.SpawnedSoulIDs {
						if soulID == runePlacementSpawnedSoulID {
							remainingRunePlacementSpawnedSoulIDs = append(remainingRunePlacementSpawnedSoulIDs, soulID)
						}
					}
				}

				if len(remainingRunePlacementSpawnedSoulIDs) >= 2 {
					return model.SpawnSoulData{}, false
				}

				// TODO: spawn cooldown
				return model.SpawnSoulData{
					ID:             uuid.New(),
					Position:       runePlacementData.Position,
					SoulTemplateID: model.SoulID,
					RuneOwnerID:    selfID,
				}, true
			},
		),
	},
})
