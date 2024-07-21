package health

import (
	"thief/base/engine"
	"thief/scared"
	"thief/scared/model"
	"thief/scared/projectors"

	"github.com/google/uuid"
)

type weaponHitLogConsumer struct {
	storeHealth *engine.Store
}

func NewWeaponHitLogConsumer(storeHealth *engine.Store) engine.Receiver[model.WeaponLog] {
	return weaponHitLogConsumer{
		storeHealth: storeHealth,
	}
}

func (s weaponHitLogConsumer) GetEvents(entityID uuid.UUID, data model.WeaponLog) map[*engine.Store][]engine.Event {
	log := data.Log

	hitEvents := []engine.Event{}
	for _, targetID := range log.Targets {
		entity := projectors.ProjectorEntityType.Project(targetID)
		switch entity {
		case model.EntityTypeKnight, model.EntityTypeShip:
			hitEvents = append(hitEvents, hitShipOrKnight(targetID, log.Value.Plus-log.Value.Minus))
		case model.EntityTypeSoul:
			hitEvents = append(hitEvents, hitSoul(targetID, log.Value.Plus-log.Value.Minus))
		}
	}

	return map[*engine.Store][]engine.Event{
		s.storeHealth: hitEvents,
	}
}

func hitShipOrKnight(targetID uuid.UUID, totalDamage int) engine.Event {
	devotion := projectors.ProjectorDevotion.Project(scared.DevotionID)
	if devotion >= totalDamage {
		return EffectHitDevotion.NewEvent(targetID, totalDamage)
	}

	return EffectHit.NewEvent(targetID, totalDamage)
}

func hitSoul(targetID uuid.UUID, totalDamage int) engine.Event {
	return EffectHit.NewEvent(targetID, totalDamage)
}
