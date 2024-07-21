package scene

import "thief/scared/model"

type areaMap struct {
	start         model.Position
	end           model.Position
	churches      []model.Position
	runePlacement []model.SpawnRunePlacementData
}
