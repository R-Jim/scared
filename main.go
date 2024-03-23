package main

import (
	"fmt"
	"image/color"
	"log"
	"thief/internal/constant"
	"thief/internal/model"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	PLAYER_CHARACTER = "P"

	ENEMY_CHARACTER = "E"

	COIN_CHARACTER = "C"
)

var (
	PLAYER_COLOR = color.White
	ENEMY_COLOR  = color.RGBA{0xff, 0x0, 0x0, 0xff}
	COIN_COLOR   = color.RGBA{0xff, 165, 0x0, 0xff}
)

type EnemyCrons struct {
}

type Game struct {
	MainFloor string

	PlayerPosition    *model.Position
	PlayerFacingRight bool
	PlayerHitLog      []model.HitLog

	TestEnemyID uuid.UUID

	EnemyPositions map[uuid.UUID]*model.Position
	EnemyCronsMap  map[uuid.UUID]*EnemyCrons
	EnemyHitLogMap *model.HitLogMap

	CoinPositions map[uuid.UUID]*model.Position

	EnemyMap map[uuid.UUID]model.Enemy
	CoinMap  []uuid.UUID
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.PlayerPosition.X -= constant.PLAYER_MOVE_DISTANCE
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.PlayerPosition.X += constant.PLAYER_MOVE_DISTANCE
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		// var nearestEnemyID uuid.UUID
		// var nearestEnemyPosition model.Position
		fmt.Println("yo1")

		tmp := (*g.EnemyHitLogMap)[g.TestEnemyID]
		tmp2 := append(*tmp, model.HitLog{})

		(*g.EnemyHitLogMap)[g.TestEnemyID] = &tmp2
	}

	// for _, enemyCrons := range g.EnemyCronsMap {
	// 	if isRun, _ := enemyCrons.Delete.Run(); isRun {
	// 		fmt.Println("yo")
	// 	}
	// 	enemyCrons.Patrol.Run()
	// 	// enemyCrons.Strike.Run()
	// }

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// ebitenutil.DebugPrint(screen, "Hello, World!")
	ebitenutil.DebugPrintAt(screen, PLAYER_CHARACTER, g.PlayerPosition.X, 0)

	for _, enemyPosition := range g.EnemyPositions {
		ebitenutil.DebugPrintAt(screen, ENEMY_CHARACTER, enemyPosition.X, 0)
	}

	ebitenutil.DebugPrintAt(screen, COIN_CHARACTER, 100, 0)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")

	playerPosition := model.Position{
		X: 0,
	}

	enemyPosition := model.Position{
		X: 50,
	}
	enemyID := uuid.New()
	enemyHitLogs := []model.HitLog{}

	enemyHitLogMap := model.HitLogMap{
		enemyID: &enemyHitLogs,
	}

	if err := ebiten.RunGame(&Game{
		PlayerPosition: &playerPosition,
		EnemyPositions: map[uuid.UUID]*model.Position{
			enemyID: &enemyPosition,
		},
		// EnemyCronsMap: map[uuid.UUID]*EnemyCrons{
		// 	enemyID: {
		// 		cron.EnemyPatrolCron{
		// 			Position:       &enemyPosition,
		// 			PlayerPosition: &playerPosition,
		// 		},
		// 		cron.EnemyStrikeCron{
		// 			Position:       &enemyPosition,
		// 			PlayerPosition: &playerPosition,
		// 		},
		// 		cron.EnemyDeleteCron{
		// 			HitLogs: &enemyHitLogs,
		// 		},
		// 	},
		// },

		EnemyHitLogMap: &enemyHitLogMap,

		TestEnemyID: enemyID,
	}); err != nil {
		log.Fatal(err)
	}
}
