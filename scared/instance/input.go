package instance

import (
	"thief/scared/model"

	"github.com/google/uuid"
)

var (
	SetRuneToEquippedWeaponRuneSlotFunc func(runeID, slotID uuid.UUID)
	AssignWeaponToEntityFunc            func(weaponTemplateID, entityID uuid.UUID)
	SetWaypointFunc                     func(waypoint model.Waypoint)
	SetActiveWeapon                     func(weaponID uuid.UUID)
	// SetAcolyte                          func(containerID uuid.UUID, numberOfAcolyte int)
	// TransferAcolyte                     func(fromContainerID, toContainerID uuid.UUID, numberOfAcolyte int) error
)

type projectorRuneToEquippedWeaponRuneSlot struct {
	equippedWeaponRuleSlotWithRuneTemplateMapping map[uuid.UUID]uuid.UUID
}

func newProjectorRuneToEquippedWeaponRuneSlot() projectorRuneToEquippedWeaponRuneSlot {
	return projectorRuneToEquippedWeaponRuneSlot{
		equippedWeaponRuleSlotWithRuneTemplateMapping: map[uuid.UUID]uuid.UUID{},
	}
}

func (p projectorRuneToEquippedWeaponRuneSlot) AddRuneToRuneSlot(runeID, slotID uuid.UUID) {
	p.equippedWeaponRuleSlotWithRuneTemplateMapping[slotID] = runeID
}

func (p projectorRuneToEquippedWeaponRuneSlot) Project(slotID uuid.UUID) uuid.UUID {
	return p.equippedWeaponRuleSlotWithRuneTemplateMapping[slotID]
}

func (p projectorRuneToEquippedWeaponRuneSlot) ListIdentifiers(filters ...func(m uuid.UUID) bool) []uuid.UUID {
	runeTemplateID := []uuid.UUID{}
	for entityID := range p.equippedWeaponRuleSlotWithRuneTemplateMapping {
		projection := p.Project(entityID)

		isMatchedFilters := true
		for _, filter := range filters {
			if !filter(projection) {
				isMatchedFilters = false
				break
			}
		}

		if isMatchedFilters {
			runeTemplateID = append(runeTemplateID, entityID)
		}
	}

	return runeTemplateID
}

func (p projectorRuneToEquippedWeaponRuneSlot) ListDeletedIdentifiers(filters ...func(m uuid.UUID) bool) []uuid.UUID {
	return nil
}

func (p projectorRuneToEquippedWeaponRuneSlot) IsDestroyed(slotID uuid.UUID) bool {
	_, isExist := p.equippedWeaponRuleSlotWithRuneTemplateMapping[slotID]
	return !isExist
}

type projectorWaypoint struct {
	waypointMapping map[uuid.UUID]model.Waypoint
}

func newProjectorWaypoint() projectorWaypoint {
	return projectorWaypoint{
		waypointMapping: map[uuid.UUID]model.Waypoint{},
	}
}

func (p projectorWaypoint) SetWaypoint(waypoint model.Waypoint) {
	p.waypointMapping[waypoint.OwnerID] = waypoint
}

func (p projectorWaypoint) Project(slotID uuid.UUID) model.Waypoint {
	return p.waypointMapping[slotID]
}

func (p projectorWaypoint) ListIdentifiers(filters ...func(m model.Waypoint) bool) []uuid.UUID {
	runeTemplateID := []uuid.UUID{}
	for entityID := range p.waypointMapping {
		projection := p.Project(entityID)

		isMatchedFilters := true
		for _, filter := range filters {
			if !filter(projection) {
				isMatchedFilters = false
				break
			}
		}

		if isMatchedFilters {
			runeTemplateID = append(runeTemplateID, entityID)
		}
	}

	return runeTemplateID
}

func (p projectorWaypoint) ListDeletedIdentifiers(filters ...func(m model.Waypoint) bool) []uuid.UUID {
	return nil
}

func (p projectorWaypoint) IsDestroyed(slotID uuid.UUID) bool {
	_, isExist := p.waypointMapping[slotID]
	return !isExist
}

type projectorEntityAssignedWeapon struct {
	assignedWeaponMapping map[uuid.UUID]uuid.UUID
}

func newProjectorEntityAssignedWeapon() projectorEntityAssignedWeapon {
	return projectorEntityAssignedWeapon{
		assignedWeaponMapping: map[uuid.UUID]uuid.UUID{},
	}
}

func (p projectorEntityAssignedWeapon) AssignWeaponToEntity(weaponTemplateID, entityID uuid.UUID) {
	p.assignedWeaponMapping[entityID] = weaponTemplateID
}

func (p projectorEntityAssignedWeapon) Project(entityID uuid.UUID) uuid.UUID {
	return p.assignedWeaponMapping[entityID]
}

func (p projectorEntityAssignedWeapon) ListIdentifiers(filters ...func(m uuid.UUID) bool) []uuid.UUID {
	runeTemplateID := []uuid.UUID{}
	for entityID := range p.assignedWeaponMapping {
		projection := p.Project(entityID)

		isMatchedFilters := true
		for _, filter := range filters {
			if !filter(projection) {
				isMatchedFilters = false
				break
			}
		}

		if isMatchedFilters {
			runeTemplateID = append(runeTemplateID, entityID)
		}
	}

	return runeTemplateID
}

func (p projectorEntityAssignedWeapon) ListDeletedIdentifiers(filters ...func(m uuid.UUID) bool) []uuid.UUID {
	return nil
}

func (p projectorEntityAssignedWeapon) IsDestroyed(entityID uuid.UUID) bool {
	_, isExist := p.assignedWeaponMapping[entityID]
	return !isExist
}

type projectorActiveWeapon struct {
	activeWeaponMapping map[uuid.UUID]bool
}

func newProjectorActiveWeapon() projectorActiveWeapon {
	return projectorActiveWeapon{
		activeWeaponMapping: map[uuid.UUID]bool{},
	}
}

func (p projectorActiveWeapon) SetActiveWeapon(weaponID uuid.UUID) {
	p.activeWeaponMapping[weaponID] = true
}

func (p projectorActiveWeapon) Project(weaponID uuid.UUID) bool {
	return p.activeWeaponMapping[weaponID]
}

func (p projectorActiveWeapon) ListIdentifiers(filters ...func(m bool) bool) []uuid.UUID {
	runeTemplateID := []uuid.UUID{}
	for entityID := range p.activeWeaponMapping {
		projection := p.Project(entityID)

		isMatchedFilters := true
		for _, filter := range filters {
			if !filter(projection) {
				isMatchedFilters = false
				break
			}
		}

		if isMatchedFilters {
			runeTemplateID = append(runeTemplateID, entityID)
		}
	}

	return runeTemplateID
}

func (p projectorActiveWeapon) ListDeletedIdentifiers(filters ...func(m bool) bool) []uuid.UUID {
	return nil
}

func (p projectorActiveWeapon) IsDestroyed(slotID uuid.UUID) bool {
	_, isExist := p.activeWeaponMapping[slotID]
	return !isExist
}

// type projectorAcolyte struct {
// 	acolyteMapping map[uuid.UUID]int
// }

// func newProjectorAcolyte() projectorAcolyte {
// 	return projectorAcolyte{
// 		acolyteMapping: map[uuid.UUID]int{},
// 	}
// }

// func (p projectorAcolyte) SetAcolyte(containerID uuid.UUID, numberOfAcolyte int) {
// 	p.acolyteMapping[containerID] = numberOfAcolyte
// }

// func (p projectorAcolyte) TransformerAcolyte(fromContainerID, toContainerID uuid.UUID, numberOfAcolyte int) error {
// 	requestedNumberOfAcolyte := p.acolyteMapping[fromContainerID]
// 	if requestedNumberOfAcolyte < numberOfAcolyte {
// 		return errors.New("insufficient number of requested acolyte")
// 	}
// 	p.acolyteMapping[toContainerID] += numberOfAcolyte
// 	p.acolyteMapping[fromContainerID] -= numberOfAcolyte
// 	return nil
// }

// func (p projectorAcolyte) Project(weaponID uuid.UUID) int {
// 	return p.acolyteMapping[weaponID]
// }

// func (p projectorAcolyte) ListIdentifiers(filters ...func(m int) bool) []uuid.UUID {
// 	runeTemplateID := []uuid.UUID{}
// 	for entityID := range p.acolyteMapping {
// 		projection := p.Project(entityID)

// 		isMatchedFilters := true
// 		for _, filter := range filters {
// 			if !filter(projection) {
// 				isMatchedFilters = false
// 				break
// 			}
// 		}

// 		if isMatchedFilters {
// 			runeTemplateID = append(runeTemplateID, entityID)
// 		}
// 	}

// 	return runeTemplateID
// }

// func (p projectorAcolyte) ListDeletedIdentifiers(filters ...func(m int) bool) []uuid.UUID {
// 	return nil
// }

// func (p projectorAcolyte) IsDestroyed(slotID uuid.UUID) bool {
// 	_, isExist := p.acolyteMapping[slotID]
// 	return !isExist
// }
