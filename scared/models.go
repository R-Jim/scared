package scared

import (
	"math"

	"github.com/google/uuid"
)

type Position struct {
	X int
	Y int
}

var (
	PositionShipSpawn = Position{
		X: -3000,
		Y: -3000,
	}
	PositionSoulSpawn = Position{
		X: -3000,
		Y: -3000,
	}
)

func (from Position) DistanceOf(to Position) float64 {
	return math.Sqrt(math.Pow(float64(to.X-from.X), 2) + math.Pow(float64(to.Y-from.Y), 2))
}

type EntityType string

const (
	EntityTypeShip EntityType = "Ship"
	EntityTypeSoul EntityType = "Soul"
)

var (
	EntityTypeTargetMapping = map[EntityType][]EntityType{
		EntityTypeShip: {EntityTypeSoul},
	}
)

type WeaponRuneSlotType string

const (
	WeaponRuneSlotTypeCoolDown WeaponRuneSlotType = "CoolDown"
	WeaponRuneSlotTypeRange    WeaponRuneSlotType = "Range"
	WeaponRuneSlotTypeDamage   WeaponRuneSlotType = "Damage"
	WeaponRuneSlotTypePattern  WeaponRuneSlotType = "Pattern"
)

type RuneTemplate struct {
	Value    Stat
	Modifier Stat
}

var (
	TimeTwoRuneID  = uuid.New()
	TimeFourRuneID = uuid.New()

	RuneTemplates = map[uuid.UUID]RuneTemplate{
		TimeTwoRuneID: {
			Modifier: Stat{
				Plus:  2,
				Minus: 2,
			},
		},
		TimeFourRuneID: {
			Modifier: Stat{
				Plus:  4,
				Minus: 4,
			},
		},
	}
)

type EquippedRune struct {
	TemplateID       uuid.UUID
	WeaponRuneSlotID uuid.UUID
}

type WeaponRuneSlot struct {
	Type   WeaponRuneSlotType
	RuneID uuid.UUID
}

type Stat struct {
	Minus int
	Plus  int
}

type WeaponTemplate struct {
	CoolDown int
	Range    int
	Damage   Stat
	Patterns []Position
	RuneSlot []WeaponRuneSlot
}

// Only equipped, user input to equip weapon -> spawn composer -> spawn new weapon with input info, destroyed when unequipped
type EquippedWeapon struct {
	TemplateID  uuid.UUID
	OwnerID     uuid.UUID
	RuneSlotIDs []uuid.UUID
}

type Log[model any] struct {
	ID      uuid.UUID
	Minus   model
	Plus    model
	Targets []uuid.UUID
}

type WeaponLog struct {
	CoolDown int
	Log      Log[int]
}

var (
	HunterTargetMapping = map[EntityType][]EntityType{
		EntityTypeSoul: {EntityTypeShip},
	}
)

var (
	WeaponMusketID = uuid.New()

	WeaponTemplates = map[uuid.UUID]WeaponTemplate{
		WeaponMusketID: {
			CoolDown: 60,
			Range:    200,
			Damage: Stat{
				Minus: 1,
			},
			Patterns: []Position{
				{X: 0, Y: 0},
			},
			RuneSlot: []WeaponRuneSlot{
				{
					Type: WeaponRuneSlotTypeDamage,
				},
				{
					Type: WeaponRuneSlotTypeRange,
				},
				{
					Type: WeaponRuneSlotTypeCoolDown,
				},
				{
					Type: WeaponRuneSlotTypeCoolDown,
				},
			},
		},
	}
)

type ShipGuard struct {
	WeaponID uuid.UUID
	Quantity int
}

type ShipTemplate struct {
	GuardQuantity int
}

var (
	ShipDawnBreakID = uuid.New()

	ShipTemplates = map[uuid.UUID]ShipTemplate{
		ShipDawnBreakID: {
			GuardQuantity: 2,
		},
	}
)
