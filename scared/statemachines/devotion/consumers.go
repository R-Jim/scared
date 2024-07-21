package devotion

import (
	"thief/base/engine"
	"thief/scared"
	"thief/scared/model"

	"github.com/google/uuid"
)

type devotionGenerator struct {
	storeDevotion *engine.Store
}

func NewAfterEffectDevotionGenerator(storeDevotion *engine.Store) engine.Receiver[int] {
	return devotionGenerator{
		storeDevotion: storeDevotion,
	}
}

func (a devotionGenerator) GetEvents(entityID uuid.UUID, data int) map[*engine.Store][]engine.Event {
	return map[*engine.Store][]engine.Event{
		a.storeDevotion: {
			EffectAdd.NewEvent(scared.DevotionID, 2),
		},
	}
}

type devotionWeaponConsumer struct {
	storeDevotion *engine.Store
}

func NewAfterEffectDevotionWeaponConsumer(storeDevotion *engine.Store) engine.Receiver[model.WeaponLog] {
	return devotionWeaponConsumer{
		storeDevotion: storeDevotion,
	}
}

func (a devotionWeaponConsumer) GetEvents(entityID uuid.UUID, data model.WeaponLog) map[*engine.Store][]engine.Event {
	return map[*engine.Store][]engine.Event{
		a.storeDevotion: {
			EffectConsume.NewEvent(scared.DevotionID, data.DevotionCost),
		},
	}
}

type devotionHealthConsumer struct {
	storeDevotion *engine.Store
}

func NewAfterEffectDevotionHealthConsumer(storeDevotion *engine.Store) engine.Receiver[int] {
	return devotionHealthConsumer{
		storeDevotion: storeDevotion,
	}
}

func (a devotionHealthConsumer) GetEvents(entityID uuid.UUID, data int) map[*engine.Store][]engine.Event {
	return map[*engine.Store][]engine.Event{
		a.storeDevotion: {
			EffectConsume.NewEvent(scared.DevotionID, -data),
		},
	}
}
