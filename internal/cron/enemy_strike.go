package cron

import (
	"fmt"
	model "thief/internal/model"
)

type enemyStrikeCron struct {
	Position       *model.Position
	PlayerPosition *model.Position
}

func (c enemyStrikeCron) Run() (bool, error) {
	// fmt.Printf("%d, %d", c.PlayerPosition.X, c.Position.X)
	if c.PlayerPosition.X == c.Position.X {
		fmt.Println("hit")
		return true, nil
	}

	return false, nil
}
