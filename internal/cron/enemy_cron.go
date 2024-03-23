package cron

type EnemyCron struct {
	patrol enemyPatrolCron
	strike enemyStrikeCron
	delete enemyDeleteCron
}

func (c EnemyCron) Run() {
	c.patrol.Run()
	c.strike.Run()
	c.delete.Run()
}
