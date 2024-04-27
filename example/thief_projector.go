package example

import (
	"log"
	base "thief/internal/engine"
	"thief/internal/model"

	"github.com/google/uuid"
)

type ThiefProjector struct {
	thiefStore *base.Store
}

func NewThiefProjector(thiefStore *base.Store) base.Projector {
	return ThiefProjector{
		thiefStore: thiefStore,
	}
}

func (p ThiefProjector) Project(identifier uuid.UUID, field string) interface{} {
	events, err := p.thiefStore.GetEventsByEntityID(identifier)
	if err != nil {
		log.Fatalln(err)
	}

	var defaultValue interface{}
	for i := len(events) - 1; i >= 0; i-- {
		event := events[i]
		switch field {
		case FieldThiefPosition:
			defaultValue = model.Position{}
			switch event.Effect {
			case effectThiefMove:
				return base.ParseData[model.Position](event)
			}
		case fieldThiefMoveInput:
			defaultValue = uuid.Nil
			switch event.Effect {
			case effectThiefAddEnergy:
				return base.ParseData[thiefMoveOutput](event).inputID
			}
		}
	}

	switch field {
	case FieldThiefEnergy:
		var energy thiefEnergy
		for _, event := range events {
			switch event.Effect {
			case effectThiefAddEnergy:
				energy = base.ParseData[thiefMoveOutput](event).energy
			case effectThiefMove:
				if energy.X > 0 {
					energy.X--
				} else if energy.X < 0 {
					energy.X++
				}

				if energy.Y > 0 {
					energy.Y--
				}
			}
		}
		return energy
	}

	return defaultValue
}

func (p ThiefProjector) ListIdentifiers() []uuid.UUID {
	identifiers := []uuid.UUID{}
	for id := range p.thiefStore.GetEvents() {
		identifiers = append(identifiers, id)
	}

	return identifiers
}
