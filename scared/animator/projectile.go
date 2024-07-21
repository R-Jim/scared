package animator

import (
	"thief/base/engine"
	"thief/scared/model"
	"thief/scared/projectors"
	"thief/scared/statemachines/weapon"
)

type projectile struct {
	targetPosition model.Position
}

type projectileAnimator struct {
	pendingProjectiles  []projectile
	remainingAnimations []engine.Animation
}

func NewProjectileAnimator() engine.Animator {
	return &projectileAnimator{}
}

func (a *projectileAnimator) GetHook() engine.Hook {
	return func(e engine.Event) {
		if e.Effect != string(weapon.EffectHitEnemy) {
			return
		}

		weaponLog := e.Data.(model.WeaponLog)

		for _, targetID := range weaponLog.Log.Targets {
			targetPosition := projectors.ProjectorPosition.Project(targetID)

			a.pendingProjectiles = append(a.pendingProjectiles, projectile{
				targetPosition: targetPosition,
			})
		}
	}
}

var (
	newAnimationHitMarker = func(x, y int) engine.Animation {
		return engine.NewAnimation(
			false,
			engine.Frame{RenderLayer: RenderLayerHitMarker, Image: ImageHitMarker, RenderPosition: &engine.RenderPosition{X: x, Y: y}},
			engine.Frame{RenderLayer: RenderLayerHitMarker, Image: ImageHitMarker, RenderPosition: &engine.RenderPosition{X: x, Y: y}},
			engine.Frame{RenderLayer: RenderLayerHitMarker, Image: ImageHitMarker, RenderPosition: &engine.RenderPosition{X: x, Y: y}},
			engine.Frame{RenderLayer: RenderLayerHitMarker, Image: ImageHitMarker, RenderPosition: &engine.RenderPosition{X: x, Y: y}},
			engine.Frame{RenderLayer: RenderLayerHitMarker, Image: ImageHitMarker, RenderPosition: &engine.RenderPosition{X: x, Y: y}},
			engine.Frame{RenderLayer: RenderLayerHitMarker, Image: ImageHitMarker, RenderPosition: &engine.RenderPosition{X: x, Y: y}},
		)
	}
)

func (a *projectileAnimator) Frame() []engine.Frame {
	animations := a.remainingAnimations

	for _, projectile := range a.pendingProjectiles {
		animations = append(animations, newAnimationHitMarker(projectile.targetPosition.X, projectile.targetPosition.Y))
	}
	a.pendingProjectiles = []projectile{}

	remainingAnimations := []engine.Animation{}

	frames := []engine.Frame{}
	for _, animation := range animations {
		frames = append(frames, animation.Frame(FPS))
		if !animation.IsCompleted() {
			remainingAnimations = append(remainingAnimations, animation)
		}
	}

	a.remainingAnimations = remainingAnimations

	return frames
}
