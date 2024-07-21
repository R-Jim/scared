package acolyte

import (
	"thief/base/engine"
	"thief/scared/model"

	"github.com/google/uuid"
)

type acolyteTransformer struct {
	storeAcolyte *engine.Store
}

func NewAfterEffectAcolyteTransfer(storeAcolyte *engine.Store) engine.Receiver[model.TransferData] {
	return acolyteTransformer{
		storeAcolyte: storeAcolyte,
	}
}

func (a acolyteTransformer) GetEvents(entityID uuid.UUID, data model.TransferData) map[*engine.Store][]engine.Event {
	return map[*engine.Store][]engine.Event{
		a.storeAcolyte: {
			EffectDeposit.NewEvent(data.To, data.Value),
			EffectWithdraw.NewEvent(data.From, data.Value),
		},
	}
}
