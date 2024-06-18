package model

import "math"

type Position struct {
	X int
	Y int
}

func (from Position) DistanceOf(to Position) float64 {
	return math.Sqrt(math.Pow(float64(to.X-from.X), 2) + math.Pow(float64(to.Y-from.Y), 2))
}
