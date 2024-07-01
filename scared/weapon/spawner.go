package weapon

import (
	"fmt"
	"log"
	"strings"
	"thief/base/engine"
	"thief/scared"

	"github.com/google/uuid"
)

const (
	StoreEquippedWeapon         = "StoreEquippedWeapon"
	StoreEquippedWeaponRuneSlot = "StoreEquippedWeaponRuneSlot"
)

type spawnerEquippedWeapon struct {
	storeEquippedWeapon         *engine.Store
	storeEquippedWeaponRuneSlot *engine.Store
}

func NewSpawnerEquippedWeapon(storeEquippedWeapon, storeEquippedWeaponRuneSlot *engine.Store) engine.Spawner {
	return spawnerEquippedWeapon{
		storeEquippedWeapon:         storeEquippedWeapon,
		storeEquippedWeaponRuneSlot: storeEquippedWeaponRuneSlot,
	}
}

func (s spawnerEquippedWeapon) GetInitEvents(eventData interface{}) map[string]engine.Event {
	inputData := eventData.(engine.InputData)
	templateID, err := uuid.Parse(inputData.Key)
	if err != nil {
		log.Fatalln(err)
	}

	ownerID := inputData.TransitionData.(uuid.UUID)
	template := scared.WeaponTemplates[templateID]

	log.Println("equipped weapon")

	runeSlotIDs := make([]uuid.UUID, len(template.RuneSlot))
	for i := range template.RuneSlot {
		runeSlotIDs[i] = uuid.New()
	}

	result := map[string]engine.Event{
		StoreEquippedWeapon: {
			EntityID: uuid.New(),
			Effect:   engine.EffectInit,
			Data: scared.EquippedWeapon{
				TemplateID:  templateID,
				OwnerID:     ownerID,
				RuneSlotIDs: runeSlotIDs,
			},
		},
	}

	for i, runeSlotID := range runeSlotIDs {
		result[fmt.Sprintf("%s%d", StoreEquippedWeaponRuneSlot, i)] = engine.Event{
			EntityID: runeSlotID,
			Effect:   engine.EffectInit,
			Data:     template.RuneSlot[i].Type,
		}
	}

	return result
}

func (s spawnerEquippedWeapon) GetStore(store string) *engine.Store {
	switch true {
	case store == StoreEquippedWeapon:
		return s.storeEquippedWeapon
	case strings.Contains(store, StoreEquippedWeaponRuneSlot):
		return s.storeEquippedWeaponRuneSlot
	default:
		return nil
	}
}
