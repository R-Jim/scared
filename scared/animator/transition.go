package animator

import "thief/base/engine"

type transition struct {
	animation engine.Animation
	nextState engine.State
}
