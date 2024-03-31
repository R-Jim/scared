package example

import (
	"log"
	"thief/internal/base"

	"github.com/google/uuid"
)

type ControllerProjector struct {
	controllerStore *base.Store
}

func NewControllerProjector(controllerStore *base.Store) base.Projector {
	return ControllerProjector{
		controllerStore: controllerStore,
	}
}

func (p ControllerProjector) Project(identifier uuid.UUID, field string) interface{} {
	events, err := p.controllerStore.GetEventsByEntityID(identifier)
	if err != nil {
		log.Fatalln(err)
	}

	for i := len(events) - 1; i >= 0; i-- {
		event := events[i]
		switch field {
		case fieldControllerThiefInput:
			switch event.Effect {
			case EffectControllerMove:
				moveInputData := base.ParseData[ControllerMoveInput](event)
				moveInputData.Value = base.ParseSystemData[moveInput](event)
				return moveInputData
			default:
				return ControllerMoveInput{}
			}
		}
	}

	return nil
}

func (p ControllerProjector) ListIdentifiers() []uuid.UUID {
	identifiers := []uuid.UUID{}
	for id := range p.controllerStore.GetEvents() {
		identifiers = append(identifiers, id)
	}

	return identifiers
}
