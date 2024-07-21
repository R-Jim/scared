package animator

import (
	"thief/base/engine"
	"thief/scared/model"
	"thief/scared/projectors"
	"thief/scared/statemachines/entitytype"

	"github.com/google/uuid"
)

var (
	stateKnightIdle      engine.State = "Idle"
	stateKnightMove      engine.State = "Move"
	stateKnightDamage    engine.State = "Damage"
	stateKnightDestroyed engine.State = "Destroyed"
)

type knight struct {
	state          engine.State
	renderPosition *engine.RenderPosition
	transition     transition
}

type knightAnimator struct {
	knights map[uuid.UUID]knight
}

func NewKnightAnimator() engine.Animator {
	return &knightAnimator{
		knights: map[uuid.UUID]knight{},
	}
}

func (a *knightAnimator) GetHook() engine.Hook {
	return func(e engine.Event) {
		switch e.Effect {
		case string(entitytype.EffectInit):
			entityType, ok := e.Data.(model.EntityType)
			if ok && entityType == model.EntityTypeKnight {
				position := &engine.RenderPosition{}

				a.knights[e.EntityID] = knight{
					state: stateKnightIdle,
					transition: transition{
						animation: newAnimationKnightIdle(position),
						nextState: stateKnightIdle,
					},
					renderPosition: position,
				}
			}
		case string(entitytype.EffectDestroy):
			delete(a.knights, e.EntityID)
		}
	}
}

var (
	newAnimationKnightIdle = func(renderPosition *engine.RenderPosition) engine.Animation {
		return engine.NewAnimation(
			true,
			engine.Frame{RenderLayer: RenderLayerEntity, Image: ImageKnight, RenderPosition: renderPosition},
		)
	}
)

func (a *knightAnimator) Frame() []engine.Frame {
	frames := []engine.Frame{}

	for id, s := range a.knights {
		position := projectors.ProjectorPosition.Project(id)

		s.renderPosition.X = position.X
		s.renderPosition.Y = position.Y

		currentTransition := s.transition

		frames = append(frames, currentTransition.animation.Frame(FPS))

		if currentTransition.animation.IsCompleted() {
			if s.state != currentTransition.nextState {
				var nextTransition transition
				switch currentTransition.nextState {
				case stateKnightIdle:
					nextTransition = transition{
						animation: newAnimationKnightIdle(s.renderPosition),
						nextState: stateKnightIdle,
					}
				}

				a.knights[id] = knight{
					state:          currentTransition.nextState,
					renderPosition: s.renderPosition,
					transition:     nextTransition,
				}
			}
		}
	}

	return frames
}
