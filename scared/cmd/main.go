package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"thief/base/engine"
	"thief/scared/animator"
	"thief/scared/instance"
	"thief/scared/model"
	"thief/scared/projectors"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	SHIP_CHARACTER   = "P"
	KNIGHT_CHARACTER = "K"

	RUNE_CHARACTER = "R"

	SOUL_CHARACTER = "E"

	CHURCH_CHARACTER = "C"
)

var (
	PLAYER_COLOR = color.White
	ENEMY_COLOR  = color.RGBA{0xff, 0x0, 0x0, 0xff}
	COIN_COLOR   = color.RGBA{0xff, 165, 0x0, 0xff}
)

type EnemyCrons struct {
}

type Game struct {
	Animators []engine.Animator
}

func (g *Game) Update() error {
	instance.OperateConsumers()

	var shipID uuid.UUID
	{
		shipIDs := projectors.ProjectorEntityType.ListIdentifiers(func(et model.EntityType) bool {
			return et == model.EntityTypeShip
		})
		if len(shipIDs) != 0 {
			shipID = shipIDs[0]
		}
	}

	if shipID == uuid.Nil {
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		equippedWeaponIDs := projectors.ProjectorEquippedWeapon.ListIdentifiers(func(ew model.EquippedWeapon) bool {
			return projectors.ProjectorEntityType.Project(ew.OwnerID) == model.EntityTypeShip
		})

		if len(equippedWeaponIDs) > 0 {
			instance.SetActiveWeapon(equippedWeaponIDs[0])
			log.Println("weapon activated")
		} else {
			log.Println("weapon assigned")
			instance.AssignWeaponToEntityFunc(model.WeaponMusketID, shipID)
		}
	}
	// if inpututil.IsKeyJustPressed(ebiten.KeyA) {
	// 	instance.SetAcolyte(shipID, 2)
	// }
	// if inpututil.IsKeyJustPressed(ebiten.KeyS) {
	// 	blessingAltarIDs := projectors.ProjectorBlessingAltar.ListIdentifiers()

	// 	if len(blessingAltarIDs) > 0 {
	// 		err := instance.TransferAcolyte(shipID, blessingAltarIDs[0], 1)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }
	if ebiten.IsKeyPressed(ebiten.Key1) {
		equippedWeaponRuleSlotIDs := projectors.ProjectorEquippedWeaponRuneSlot.ListIdentifiers()
		for _, runeSlotID := range equippedWeaponRuleSlotIDs {
			runeSlot := projectors.ProjectorEquippedWeaponRuneSlot.Project(runeSlotID)
			if runeSlot.Type == model.WeaponRuneSlotTypeDamage && runeSlot.RuneID == uuid.Nil {
				instance.SetRuneToEquippedWeaponRuneSlotFunc(model.TimeFourRuneID, runeSlotID)
				break
			}
		}
	}
	if ebiten.IsKeyPressed(ebiten.Key2) {
		equippedWeaponRuleSlotIDs := projectors.ProjectorEquippedWeaponRuneSlot.ListIdentifiers()
		for _, runeSlotID := range equippedWeaponRuleSlotIDs {
			runeSlot := projectors.ProjectorEquippedWeaponRuneSlot.Project(runeSlotID)
			if runeSlot.Type == model.WeaponRuneSlotTypeRange && runeSlot.RuneID == uuid.Nil {
				instance.SetRuneToEquippedWeaponRuneSlotFunc(model.TimeFourRuneID, runeSlotID)
				break
			}
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		equippedWeaponRuleSlotIDs := projectors.ProjectorEquippedWeaponRuneSlot.ListIdentifiers()
		for _, runeSlotID := range equippedWeaponRuleSlotIDs {
			runeSlot := projectors.ProjectorEquippedWeaponRuneSlot.Project(runeSlotID)
			if runeSlot.Type == model.WeaponRuneSlotTypeCoolDown && runeSlot.RuneID == uuid.Nil {
				instance.SetRuneToEquippedWeaponRuneSlotFunc(model.TimeTwoRuneID, runeSlotID)
				break
			}
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		x, y := ebiten.CursorPosition()
		instance.SetWaypointFunc(model.Waypoint{
			OwnerID: shipID,
			Type:    model.WaypointTypeToPosition,
			Position: &model.Position{
				X: x,
				Y: y,
			},
		})
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton2) {
		instance.SetWaypointFunc(model.Waypoint{
			OwnerID:  shipID,
			Type:     model.WaypointTypeIdle,
			Position: nil,
		})
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	devotionLayer := ebiten.NewImageFromImage(screen)
	hitLayer := ebiten.NewImageFromImage(screen)
	entityLayer := ebiten.NewImageFromImage(screen)

	for _, a := range g.Animators {
		for _, frame := range a.Frame() {
			var layer *ebiten.Image
			switch frame.RenderLayer {
			case animator.RenderLayerHitMarker:
				layer = hitLayer
			case animator.RenderLayerEntity:
				layer = entityLayer
			}

			var character string
			switch true {
			case bytes.Equal(frame.Image, animator.ImageHitMarker):
				character = "x"
			case bytes.Equal(frame.Image, animator.ImageSoul):
				character = SOUL_CHARACTER
			case bytes.Equal(frame.Image, animator.ImageShip):
				character = SHIP_CHARACTER
			case bytes.Equal(frame.Image, animator.ImageKnight):
				character = KNIGHT_CHARACTER
			case bytes.Equal(frame.Image, animator.ImageRunePlacement):
				character = RUNE_CHARACTER
			case bytes.Equal(frame.Image, animator.ImageChurch):
				character = CHURCH_CHARACTER
			}

			ebitenutil.DebugPrintAt(layer, character, frame.RenderPosition.X, frame.RenderPosition.Y)
		}
	}

	var totalDevotion int
	for _, id := range projectors.ProjectorDevotion.ListIdentifiers() {
		totalDevotion += projectors.ProjectorDevotion.Project(id)
	}

	var shipID uuid.UUID
	{
		shipIDs := projectors.ProjectorEntityType.ListIdentifiers(func(et model.EntityType) bool {
			return et == model.EntityTypeShip
		})
		if len(shipIDs) != 0 {
			shipID = shipIDs[0]
		}
	}
	ebitenutil.DebugPrintAt(devotionLayer, fmt.Sprintf("Devotion: %d, Acolytes: %d", totalDevotion, projectors.ProjectorAcolyte.Project(shipID)), 2, 2)

	screen.DrawImage(hitLayer, nil)
	screen.DrawImage(entityLayer, nil)
	screen.DrawImage(devotionLayer, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")

	model.PositionRunePlacement = model.Position{
		X: 200,
		Y: 200,
	}

	model.PositionShipSpawn = model.Position{
		X: 200,
		Y: 10,
	}

	animators := instance.InitAnimators()

	instance.InitProjector()
	instance.InitEntities()

	if err := ebiten.RunGame(&Game{
		Animators: animators,
	}); err != nil {
		log.Fatal(err)
	}
}
