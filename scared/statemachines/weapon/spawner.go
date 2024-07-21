package weapon

import (
	"thief/base/engine"
	"thief/scared/model"

	"github.com/google/uuid"
)

type spawnerEquippedWeapon struct {
	storeEquippedWeapon         *engine.Store
	storeEquippedWeaponRuneSlot *engine.Store
}

func NewEquippedWeapon(storeEquippedWeapon, storeEquippedWeaponRuneSlot *engine.Store) engine.Receiver[model.ArmWeapon] {
	return spawnerEquippedWeapon{
		storeEquippedWeapon:         storeEquippedWeapon,
		storeEquippedWeaponRuneSlot: storeEquippedWeaponRuneSlot,
	}
}

func (s spawnerEquippedWeapon) GetEvents(entityID uuid.UUID, data model.ArmWeapon) map[*engine.Store][]engine.Event {
	template := model.TemplateWeapons[data.TemplateID]

	weaponID := uuid.New()

	runeSlotIDs := make([]uuid.UUID, len(template.RuneSlot))
	for i := range template.RuneSlot {
		runeSlotIDs[i] = uuid.New()
	}

	result := map[*engine.Store][]engine.Event{
		s.storeEquippedWeapon: {
			EffectInit.NewEvent(
				weaponID,
				model.EquippedWeapon{
					TemplateID:  data.TemplateID,
					OwnerID:     data.OwnerID,
					RuneSlotIDs: runeSlotIDs,
				},
			),
		},
	}

	for i, runeSlotID := range runeSlotIDs {
		result[s.storeEquippedWeaponRuneSlot] = append(result[s.storeEquippedWeaponRuneSlot],
			EffectRuneSlotInit.NewEvent(runeSlotID, template.RuneSlot[i].Type))
	}

	return result
}
