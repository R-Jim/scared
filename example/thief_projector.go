package example

import (
	"log"
	"thief/internal/base"
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

	for i := len(events) - 1; i >= 0; i-- {
		event := events[i]
		switch field {
		case FieldThiefPosition:
			switch event.Effect {
			case effectThiefMove:
				return base.ParseData[thiefMoveOutput](event).position
			default:
				return model.Position{}
			}
		case fieldThiefMoveInput:
			switch event.Effect {
			case effectThiefMove:
				return base.ParseData[thiefMoveOutput](event).inputID
			default:
				return uuid.Nil
			}
		}
	}

	return nil
}

func (p ThiefProjector) ListIdentifiers() []uuid.UUID {
	identifiers := []uuid.UUID{}
	for id := range p.thiefStore.GetEvents() {
		identifiers = append(identifiers, id)
	}

	return identifiers
}
