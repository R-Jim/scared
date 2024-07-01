package scaredrune

import (
	"log"
	"thief/base/engine"
	"thief/scared"

	"github.com/google/uuid"
)

const (
	StoreEquippedRune           = "StoreEquippedRune"
	StoreEquippedWeaponRuneSlot = "StoreEquippedWeaponRuneSlot"
)

type spawnerEquippedRune struct {
	storeEquippedRune *engine.Store
}

func NewSpawnerEquippedRune(storeEquippedRune *engine.Store) engine.Spawner {
	return spawnerEquippedRune{
		storeEquippedRune: storeEquippedRune,
	}
}

func (s spawnerEquippedRune) GetInitEvents(eventData interface{}) map[string]engine.Event {
	inputData := eventData.(engine.InputData)
	templateID, err := uuid.Parse(inputData.Key)
	if err != nil {
		log.Fatalln(err)
	}

	weaponRuneSlotID := inputData.TransitionData.(uuid.UUID)

	log.Println("equipped rune")

	return map[string]engine.Event{
		StoreEquippedRune: {
			EntityID: uuid.New(),
			Effect:   engine.EffectInit,
			Data: scared.EquippedRune{
				TemplateID:       templateID,
				WeaponRuneSlotID: weaponRuneSlotID,
			},
		},
	}
}

func (s spawnerEquippedRune) GetStore(store string) *engine.Store {
	switch store {
	case StoreEquippedRune:
		return s.storeEquippedRune
	default:
		return nil
	}
}
