package animator

import (
	"thief/base/engine"
	"thief/scared/model"
	"thief/scared/projectors"
	"thief/scared/statemachines/entitytype"

	"github.com/google/uuid"
)

var (
	stateShipIdle      engine.State = "Idle"
	stateShipMove      engine.State = "Move"
	stateShipDamage    engine.State = "Damage"
	stateShipDestroyed engine.State = "Destroyed"
)

type ship struct {
	state          engine.State
	renderPosition *engine.RenderPosition
	transition     transition
}

type shipAnimator struct {
	ships map[uuid.UUID]ship
}

func NewShipAnimator() engine.Animator {
	return &shipAnimator{
		ships: map[uuid.UUID]ship{},
	}
}

func (a *shipAnimator) GetHook() engine.Hook {
	return func(e engine.Event) {
		switch e.Effect {
		case string(entitytype.EffectInit):
			entityType, ok := e.Data.(model.EntityType)
			if ok && entityType == model.EntityTypeShip {
				position := &engine.RenderPosition{}

				a.ships[e.EntityID] = ship{
					state: stateShipIdle,
					transition: transition{
						animation: newAnimationShipIdle(position),
						nextState: stateShipIdle,
					},
					renderPosition: position,
				}
			}
		case string(entitytype.EffectDestroy):
			delete(a.ships, e.EntityID)
		}
	}
}

var (
	newAnimationShipIdle = func(renderPosition *engine.RenderPosition) engine.Animation {
		return engine.NewAnimation(
			true,
			engine.Frame{RenderLayer: RenderLayerEntity, Image: ImageShip, RenderPosition: renderPosition},
		)
	}
)

func (a *shipAnimator) Frame() []engine.Frame {
	frames := []engine.Frame{}

	for id, s := range a.ships {
		position := projectors.ProjectorPosition.Project(id)

		s.renderPosition.X = position.X
		s.renderPosition.Y = position.Y

		currentTransition := s.transition

		frames = append(frames, currentTransition.animation.Frame(FPS))

		if currentTransition.animation.IsCompleted() {
			if s.state != currentTransition.nextState {
				var nextTransition transition
				switch currentTransition.nextState {
				case stateShipIdle:
					nextTransition = transition{
						animation: newAnimationShipIdle(s.renderPosition),
						nextState: stateShipIdle,
					}
				}

				a.ships[id] = ship{
					state:          currentTransition.nextState,
					renderPosition: s.renderPosition,
					transition:     nextTransition,
				}
			}
		}
	}

	return frames
}
