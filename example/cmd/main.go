package main

import (
	"image/color"
	"log"
	"thief/example"
	"thief/internal/base"
	"thief/internal/model"

	"github.com/google/uuid"
	"github.com/hajimehoshi/bitmapfont/v2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Game struct {
	thiefID uuid.UUID

	projectorManager base.ProjectorManager

	lifeCycleComposers   []base.LifeCycleComposer
	systemInputComposers []base.SystemInputComposer

	count int
}

func (g *Game) Update() error {
	var controllerMoveInput []example.MoveInput
	var horizontalControllerMoveInput example.MoveInput

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		horizontalControllerMoveInput = example.MoveInputLeft
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if horizontalControllerMoveInput != "" {
			horizontalControllerMoveInput = ""
		}
		horizontalControllerMoveInput = example.MoveInputRight
	}

	if horizontalControllerMoveInput != "" {
		controllerMoveInput = append(controllerMoveInput, horizontalControllerMoveInput)
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		controllerMoveInput = append(controllerMoveInput, example.MoveInputJump)
	}

	if len(controllerMoveInput) > 0 {
		for _, composer := range g.systemInputComposers {
			composer.TransitionByInput(g.thiefID, example.EffectControllerMove, controllerMoveInput)
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyP) {
		log.Println(g.projectorManager.Get(example.EntityTypeThief).Project(g.thiefID, "Position"))
	}

	if g.count%2 == 0 {
		for _, composer := range g.lifeCycleComposers {
			composer.Operate()
		}
	}
	g.count++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	thiefProjector := g.projectorManager.Get(example.EntityTypeThief)
	for _, identifier := range thiefProjector.ListIdentifiers() {
		position := thiefProjector.Project(identifier, example.FieldThiefPosition).(model.Position)
		text.Draw(screen, "A", bitmapfont.Face, position.X, 100-position.Y, color.White)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")

	thiefStore := base.NewStore()
	controllerStore := base.NewStore()

	thiefID := uuid.New()

	initEvent := base.Event{
		ID:       uuid.New(),
		EntityID: thiefID,
		Effect:   "INIT",
	}
	thiefStore.AppendEvent(initEvent)
	controllerStore.AppendEvent(initEvent)

	pm := base.NewProjectorManager(
		map[string]base.Projector{
			example.EntityTypeThief:      example.NewThiefProjector(&thiefStore),
			example.EntityTypeController: example.NewControllerProjector(&controllerStore),
		},
	)

	thiefComposer := base.LifeCycleComposer(base.NewComposer(&thiefStore, pm, &example.ThiefStateMachine))

	controllerComposer := base.SystemInputComposer(base.NewComposer(&controllerStore, pm, &example.ControllerStateMachine))

	if err := ebiten.RunGame(&Game{
		thiefID: thiefID,

		projectorManager: pm,

		lifeCycleComposers:   []base.LifeCycleComposer{thiefComposer},
		systemInputComposers: []base.SystemInputComposer{controllerComposer},
	}); err != nil {
		log.Fatal(err)
	}
}
