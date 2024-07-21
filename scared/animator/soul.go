package animator

import (
	"thief/base/engine"
	"thief/scared/model"
	"thief/scared/projectors"
	"thief/scared/statemachines/entitytype"

	"github.com/google/uuid"
)

var (
	stateSoulIdle      engine.State = "Idle"
	stateSoulMove      engine.State = "Move"
	stateSoulDamage    engine.State = "Damage"
	stateSoulDestroyed engine.State = "Destroyed"
)

type soul struct {
	state          engine.State
	renderPosition *engine.RenderPosition
	transition     transition
}

type soulAnimator struct {
	souls map[uuid.UUID]soul
}

func NewSoulAnimator() engine.Animator {
	return &soulAnimator{
		souls: map[uuid.UUID]soul{},
	}
}

func (a *soulAnimator) GetHook() engine.Hook {
	return func(e engine.Event) {
		switch e.Effect {
		case string(entitytype.EffectInit):
			entityType, ok := e.Data.(model.EntityType)
			if ok && entityType == model.EntityTypeSoul {
				position := &engine.RenderPosition{}

				a.souls[e.EntityID] = soul{
					state: stateSoulIdle,
					transition: transition{
						animation: newAnimationSoulIdle(position),
						nextState: stateSoulIdle,
					},
					renderPosition: position,
				}
			}
		case string(entitytype.EffectDestroy):
			delete(a.souls, e.EntityID)
		}
	}
}

var (
	newAnimationSoulIdle = func(renderPosition *engine.RenderPosition) engine.Animation {
		return engine.NewAnimation(
			true,
			engine.Frame{RenderLayer: RenderLayerEntity, Image: ImageSoul, RenderPosition: renderPosition},
		)
	}
)

func (a *soulAnimator) Frame() []engine.Frame {
	frames := []engine.Frame{}

	for id, s := range a.souls {
		position := projectors.ProjectorPosition.Project(id)

		s.renderPosition.X = position.X
		s.renderPosition.Y = position.Y

		currentTransition := s.transition

		frames = append(frames, currentTransition.animation.Frame(FPS))

		if currentTransition.animation.IsCompleted() {
			if s.state != currentTransition.nextState {
				var nextTransition transition
				switch currentTransition.nextState {
				case stateSoulIdle:
					nextTransition = transition{
						animation: newAnimationSoulIdle(s.renderPosition),
						nextState: stateSoulIdle,
					}
				}

				a.souls[id] = soul{
					state:          currentTransition.nextState,
					renderPosition: s.renderPosition,
					transition:     nextTransition,
				}
			}
		}
	}

	return frames
}
