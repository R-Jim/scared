package animator

import (
	"thief/base/engine"
	"thief/scared/statemachines/runeplacement"

	"github.com/google/uuid"
)

var (
	stateRunePlacementIdle      engine.State = "Idle"
	stateRunePlacementDamage    engine.State = "Damage"
	stateRunePlacementDestroyed engine.State = "Destroyed"
)

type runePlacement struct {
	state          engine.State
	renderPosition *engine.RenderPosition
	transition     transition
}

type runePlacementAnimator struct {
	runePlacements map[uuid.UUID]runePlacement
}

func NewRunePlacementAnimator() engine.Animator {
	return &runePlacementAnimator{
		runePlacements: map[uuid.UUID]runePlacement{},
	}
}

func (a *runePlacementAnimator) GetHook() engine.Hook {
	return func(e engine.Event) {
		switch e.Effect {
		case string(runeplacement.EffectInit):
			data := runeplacement.EffectInit.ParseData(e)

			position := &engine.RenderPosition{
				X: data.Position.X,
				Y: data.Position.Y,
			}

			a.runePlacements[e.EntityID] = runePlacement{
				state: stateRunePlacementIdle,
				transition: transition{
					animation: newAnimationRunePlacementIdle(position),
					nextState: stateRunePlacementIdle,
				},
				renderPosition: position,
			}
		case string(runeplacement.EffectCollected):
			delete(a.runePlacements, e.EntityID)
		}
	}
}

var (
	newAnimationRunePlacementIdle = func(renderPosition *engine.RenderPosition) engine.Animation {
		return engine.NewAnimation(
			true,
			engine.Frame{RenderLayer: RenderLayerEntity, Image: ImageRunePlacement, RenderPosition: renderPosition},
		)
	}
)

func (a *runePlacementAnimator) Frame() []engine.Frame {
	frames := []engine.Frame{}

	for id, s := range a.runePlacements {
		currentTransition := s.transition

		frames = append(frames, currentTransition.animation.Frame(FPS))

		if currentTransition.animation.IsCompleted() {
			if s.state != currentTransition.nextState {
				var nextTransition transition
				switch currentTransition.nextState {
				case stateRunePlacementIdle:
					nextTransition = transition{
						animation: newAnimationRunePlacementIdle(s.renderPosition),
						nextState: stateRunePlacementIdle,
					}
				}

				a.runePlacements[id] = runePlacement{
					state:          currentTransition.nextState,
					renderPosition: s.renderPosition,
					transition:     nextTransition,
				}
			}
		}
	}

	return frames
}
