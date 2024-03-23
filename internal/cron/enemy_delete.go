package cron

import (
	model "thief/internal/model"
)

type enemyDeleteCron struct {
	HitLogs *[]model.HitLog
}

func (c enemyDeleteCron) Run() (bool, error) {
	return len(*c.HitLogs) >= 2, nil
}
