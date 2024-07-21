package model

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
	PositionRunePlacement = Position{
		X: -3000,
		Y: -3000,
	}
)

func (from Position) DistanceOf(to Position) float64 {
	return math.Sqrt(math.Pow(float64(to.X-from.X), 2) + math.Pow(float64(to.Y-from.Y), 2))
}

type EntityType string

const (
	EntityTypeShip   EntityType = "Ship"
	EntityTypeKnight EntityType = "Knight"

	EntityTypeSoul EntityType = "Soul"

	EntityTypeChurch EntityType = "Church"
)

var (
	EntityTypeMoveTargetMapping = map[EntityType][]EntityType{
		EntityTypeSoul:   {EntityTypeShip, EntityTypeKnight},
		EntityTypeKnight: {EntityTypeSoul},
	}
)

var (
	EntityTypeAttackTargetMapping = map[EntityType][]EntityType{
		EntityTypeShip:   {EntityTypeSoul},
		EntityTypeSoul:   {EntityTypeShip, EntityTypeKnight},
		EntityTypeKnight: {EntityTypeSoul},
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
	DevotionCost int

	Value    Stat
	Modifier Stat
}

var (
	TimeTwoRuneID  = uuid.New()
	TimeFourRuneID = uuid.New()

	RuneTemplates = map[uuid.UUID]RuneTemplate{
		TimeTwoRuneID: {
			DevotionCost: 2,
			Modifier: Stat{
				Plus:  2,
				Minus: 2,
			},
		},
		TimeFourRuneID: {
			DevotionCost: 4,
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

type TemplateWeapon struct {
	CoolDown int
	Range    int
	Damage   Stat
	Patterns []Position
	RuneSlot []WeaponRuneSlot
}

type ArmWeapon struct {
	TemplateID uuid.UUID
	OwnerID    uuid.UUID
}

// Only equipped, user input to equip weapon -> spawn composer -> spawn new weapon with input info, destroyed when unequipped
type EquippedWeapon struct {
	TemplateID  uuid.UUID
	OwnerID     uuid.UUID
	RuneSlotIDs []uuid.UUID
}

type Log[model any] struct {
	ID      uuid.UUID
	Value   model
	Targets []uuid.UUID
}

type WeaponLog struct {
	CoolDown     int
	DevotionCost int
	Log          Log[Stat]
}

var (
	// Soul weapons
	ClawID = uuid.New()

	// Blessed weapons
	WeaponMusketID = uuid.New()

	TemplateWeapons = map[uuid.UUID]TemplateWeapon{
		// Soul weapons
		ClawID: {
			CoolDown: 60,
			Range:    10,
			Damage: Stat{
				Minus: 1,
			},
			Patterns: []Position{
				{X: 0, Y: 0},
			},
		},

		// Blessed weapons
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

type TemplateShip struct {
	GuardQuantity   int
	AcolyteQuantity int
}

var (
	ShipDawnBreakID = uuid.New()

	TemplateShips = map[uuid.UUID]TemplateShip{
		ShipDawnBreakID: {
			GuardQuantity:   2,
			AcolyteQuantity: 5,
		},
	}
)

type SpawnShipData struct {
	TemplateID uuid.UUID
	Position   Position
}

type Soul struct {
	TemplateID uuid.UUID
}

type SpawnSoulData struct {
	ID             uuid.UUID
	Position       Position
	SoulTemplateID uuid.UUID
	RuneOwnerID    uuid.UUID
}

type TemplateSoul struct {
	WeaponTemplateID uuid.UUID
}

var (
	SoulID = uuid.New()

	TemplateSouls = map[uuid.UUID]TemplateSoul{
		SoulID: {
			WeaponTemplateID: ClawID,
		},
	}
)

type WaypointType string

const (
	WaypointTypeIdle       WaypointType = "Idle"
	WaypointTypeToPosition WaypointType = "ToPosition"
)

type Waypoint struct {
	OwnerID  uuid.UUID
	Type     WaypointType
	Position *Position
}

type SpawnRunePlacementData struct {
	RuneTemplateID uuid.UUID
	Position       Position
}

type RunePlacementData struct {
	Position       Position
	SpawnedSoulIDs []uuid.UUID
}

type BlessingAltarData struct {
	OwnerID                uuid.UUID
	NumberOfBlessedAcolyte int
}

type SpawnKnightData struct {
	Position Position
}

type TransferData struct {
	From  uuid.UUID
	To    uuid.UUID
	Value int
}
