package cron

import (
	"thief/internal/constant"
	model "thief/internal/model"
)

type enemyPatrolCron struct {
	Position       *model.Position
	PlayerPosition *model.Position
}

func (c enemyPatrolCron) Run() (bool, error) {
	if c.PlayerPosition.X > c.Position.X {
		c.Position.X += constant.ENEMY_MOVE_DISTANCE
	} else if c.PlayerPosition.X < c.Position.X {
		c.Position.X -= constant.ENEMY_MOVE_DISTANCE
	}

	return true, nil
}
