package main

import (
	"image/color"
	"log"
	"thief/base/engine"
	"thief/scared"
	"thief/scared/instance"
	"thief/scared/ship"
	"thief/scared/spawnership"
	"thief/scared/weapon"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	SHIP_CHARACTER = "P"

	SOUL_CHARACTER = "E"
)

var (
	PLAYER_COLOR = color.White
	ENEMY_COLOR  = color.RGBA{0xff, 0x0, 0x0, 0xff}
	COIN_COLOR   = color.RGBA{0xff, 165, 0x0, 0xff}
)

type EnemyCrons struct {
}

type Game struct {
	ComposerLifeCycles     []*engine.ComposerLifeCycle
	ComposerExternalInputs map[string]*engine.ComposerExternalInput
	ComposerSpawners       []*engine.ComposerSpawner
	ComposerDestroyers     []*engine.ComposerDestroyer

	SpawnerIDMappings map[scared.EntityType]uuid.UUID
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		shipIDs := scared.ProjectorEntityType.ListIdentifiers(func(et scared.EntityType) bool {
			return et == scared.EntityTypeShip
		})
		if len(shipIDs) != 0 {
			g.ComposerExternalInputs[instance.EquipWeapon].TransitionByInput(shipIDs[0], ship.EffectArm, scared.WeaponMusketID.String())
		}
	}
	if ebiten.IsKeyPressed(ebiten.Key1) {
		equippedWeaponRuleSlotIDs := scared.ProjectorEquippedWeaponRuneSlot.ListIdentifiers()
		for _, runeSlotID := range equippedWeaponRuleSlotIDs {
			runeSlot := scared.ProjectorEquippedWeaponRuneSlot.Project(runeSlotID)
			if runeSlot.Type == scared.WeaponRuneSlotTypeDamage && runeSlot.RuneID == uuid.Nil {
				g.ComposerExternalInputs[instance.EquipRune].TransitionByInput(runeSlotID, weapon.EffectRuneSlotRequest, scared.TimeTwoRuneID.String())
				break
			}
		}
	}
	if ebiten.IsKeyPressed(ebiten.Key2) {
		equippedWeaponRuleSlotIDs := scared.ProjectorEquippedWeaponRuneSlot.ListIdentifiers()
		for _, runeSlotID := range equippedWeaponRuleSlotIDs {
			runeSlot := scared.ProjectorEquippedWeaponRuneSlot.Project(runeSlotID)
			if runeSlot.Type == scared.WeaponRuneSlotTypeRange && runeSlot.RuneID == uuid.Nil {
				g.ComposerExternalInputs[instance.EquipRune].TransitionByInput(runeSlotID, weapon.EffectRuneSlotRequest, scared.TimeFourRuneID.String())
				break
			}
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		equippedWeaponRuleSlotIDs := scared.ProjectorEquippedWeaponRuneSlot.ListIdentifiers()
		for _, runeSlotID := range equippedWeaponRuleSlotIDs {
			runeSlot := scared.ProjectorEquippedWeaponRuneSlot.Project(runeSlotID)
			if runeSlot.Type == scared.WeaponRuneSlotTypeCoolDown && runeSlot.RuneID == uuid.Nil {
				g.ComposerExternalInputs[instance.EquipRune].TransitionByInput(runeSlotID, weapon.EffectRuneSlotRequest, scared.TimeFourRuneID.String())
				break
			}
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyP) {
		g.ComposerExternalInputs[instance.SpawnShip].TransitionByInput(g.SpawnerIDMappings[scared.EntityTypeShip], spawnership.EffectSpawn, scared.ShipDawnBreakID.String())
	}

	for _, composer := range g.ComposerLifeCycles {
		composer.Operate()
	}

	for _, composer := range g.ComposerDestroyers {
		composer.PlanDestroyIDs()
	}
	for _, composer := range g.ComposerDestroyers {
		composer.Commit()
	}

	for _, composer := range g.ComposerSpawners {
		composer.Operate()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	positionProjector := scared.ProjectorPosition

	for _, id := range positionProjector.ListIdentifiers() {
		entityType := scared.ProjectorEntityType.Project(id)

		var character string
		switch entityType {
		case scared.EntityTypeShip:
			character = SHIP_CHARACTER
		case scared.EntityTypeSoul:
			character = SOUL_CHARACTER
		}

		p := positionProjector.Project(id)
		ebitenutil.DebugPrintAt(screen, character, p.X, p.Y)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")

	scared.PositionSoulSpawn = scared.Position{
		X: 10,
		Y: 10,
	}

	scared.PositionShipSpawn = scared.Position{
		X: 200,
		Y: 10,
	}

	instance.InitProjector()

	if err := ebiten.RunGame(&Game{
		ComposerLifeCycles:     instance.InitComposerLifeCycle(),
		ComposerExternalInputs: instance.InitComposerExternalInput(),
		ComposerSpawners:       instance.InitComposerSpawner(),
		ComposerDestroyers:     instance.InitComposerDestroyer(),

		SpawnerIDMappings: instance.InitSpawner(),
	}); err != nil {
		log.Fatal(err)
	}
}
