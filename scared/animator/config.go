package animator

import "thief/base/engine"

const (
	FPS = 60
)

const (
	RenderLayerHitMarker engine.RenderLayer = "HitMarker"
	RenderLayerEntity    engine.RenderLayer = "Entity"
)

var (
	ImageShip          = []byte("ship")
	ImageKnight        = []byte("knight")
	ImageSoul          = []byte("soul")
	ImageHitMarker     = []byte("hit")
	ImageRunePlacement = []byte("rune_placement")
	ImageChurch        = []byte("church")
)
