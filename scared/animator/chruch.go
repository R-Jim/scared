package animator

import (
	"thief/base/engine"
	"thief/scared/statemachines/church"

	"github.com/google/uuid"
)

var (
	stateChurchIdle engine.State = "Idle"
	// stateChurchDestroyed engine.State = "Destroyed"
)

type churchData struct {
	state          engine.State
	renderPosition *engine.RenderPosition
	transition     transition
}

type churchAnimator struct {
	churches map[uuid.UUID]churchData
}

func NewChurchAnimator() engine.Animator {
	return &churchAnimator{
		churches: map[uuid.UUID]churchData{},
	}
}

func (a *churchAnimator) GetHook() engine.Hook {
	return func(e engine.Event) {
		switch e.Effect {
		case string(church.EffectInit):
			data := church.EffectInit.ParseData(e)

			position := &engine.RenderPosition{
				X: data.X,
				Y: data.Y,
			}

			a.churches[e.EntityID] = churchData{
				state: stateChurchIdle,
				transition: transition{
					animation: newAnimationChurchIdle(position),
					nextState: stateChurchIdle,
				},
				renderPosition: position,
			}
		case string(church.EffectCollect):
			delete(a.churches, e.EntityID)
		}
	}
}

var (
	newAnimationChurchIdle = func(renderPosition *engine.RenderPosition) engine.Animation {
		return engine.NewAnimation(
			true,
			engine.Frame{RenderLayer: RenderLayerEntity, Image: ImageChurch, RenderPosition: renderPosition},
		)
	}
)

func (a *churchAnimator) Frame() []engine.Frame {
	frames := []engine.Frame{}

	for id, s := range a.churches {
		currentTransition := s.transition

		frames = append(frames, currentTransition.animation.Frame(FPS))

		if currentTransition.animation.IsCompleted() {
			if s.state != currentTransition.nextState {
				var nextTransition transition
				switch currentTransition.nextState {
				case stateChurchIdle:
					nextTransition = transition{
						animation: newAnimationChurchIdle(s.renderPosition),
						nextState: stateChurchIdle,
					}
				}

				a.churches[id] = churchData{
					state:          currentTransition.nextState,
					renderPosition: s.renderPosition,
					transition:     nextTransition,
				}
			}
		}
	}

	return frames
}
