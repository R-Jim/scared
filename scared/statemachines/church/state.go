package church

import (
	"thief/base/engine"
	"thief/scared/model"
	"thief/scared/projectors"

	"github.com/google/uuid"
)

const (
	stateActive    engine.State = "Active"
	stateDestroyed              = engine.StateDestroyed

	EffectInit    engine.Effect[model.Position]     = "Init"
	EffectCollect engine.Effect[model.TransferData] = "Collect"
	EffectDestroy engine.Effect[any]                = "Destroy"
)

var StateMachine = engine.NewStateMachine(EffectInit.ToState(stateActive), engine.Nodes{
	stateActive: {
		EffectCollect.ToStateWhen(
			stateActive,
			func(selfID uuid.UUID) (model.TransferData, bool) {
				if projectors.ProjectorAcolyte.Project(selfID) <= 0 {
					return model.TransferData{}, false
				}

				shipIDs := projectors.ProjectorEntityType.ListIdentifiers(func(et model.EntityType) bool {
					return et == model.EntityTypeShip
				})

				if len(shipIDs) == 0 {
					return model.TransferData{}, false
				}

				churchPosition := projectors.ProjectorPosition.Project(selfID)

				for _, shipID := range shipIDs {
					shipPosition := projectors.ProjectorPosition.Project(shipID)
					if churchPosition.DistanceOf(shipPosition) < 10.0 {
						return model.TransferData{
							From:  selfID,
							To:    shipID,
							Value: 1,
						}, true
					}
				}

				return model.TransferData{}, false
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
