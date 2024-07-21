package acolyte

import (
	"thief/base/engine"
	"thief/scared/projectors"

	"github.com/google/uuid"
)

const (
	stateActive    engine.State = "Active"
	stateDestroyed engine.State = engine.StateDestroyed

	EffectInit     engine.Effect[int] = "Init"
	EffectDeposit  engine.Effect[int] = "Deposit"
	EffectWithdraw engine.Effect[int] = "Withdraw"
	EffectDestroy  engine.Effect[any] = "Destroy"
)

var StateMachine = engine.NewStateMachine(EffectInit.ToState(stateActive), engine.Nodes{
	stateActive: {
		EffectDestroy.ToStateWhen(
			stateActive,
			func(selfID uuid.UUID) (any, bool) {
				return nil, projectors.ProjectorEntityType.IsDestroyed(selfID)
			},
		),
	},
})
