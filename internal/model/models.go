package model

import "github.com/google/uuid"

type Position struct {
	X int
}

type Player struct {
	HealthCap   int
	HealthValue int

	CoinValue int
}

type Enemy struct {
	ID uuid.UUID

	HealthCap   int
	HealthValue int
}

type AttackType string

const (
	AttackTypeShot  AttackType = "SHOT"
	AttackTypeMelee AttackType = "MELEE"
)

type HitLog struct {
	Source     uuid.UUID
	AttackType AttackType
}

type HitLogMap map[uuid.UUID]*[]HitLog
