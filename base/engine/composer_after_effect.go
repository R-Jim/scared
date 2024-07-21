package engine

import (
	"log"

	"github.com/google/uuid"
)

const (
	effectConsume Effect[consumeData] = "consume"
)

type consumeData struct {
	OriginID      uuid.UUID
	OriginEventID uuid.UUID
}

type Receiver[model any] interface {
	GetEvents(entityID uuid.UUID, data model) map[*Store][]Event
}

type ComposerAfterEffect[model any] struct {
	eventCounter int

	originStore           *Store
	effect                Effect[model]
	receiver              Receiver[model]
	consumedEventLogStore *Store
}

func NewComposerAfterEffect[model any](originStore *Store, effect Effect[model], receiver Receiver[model]) *ComposerAfterEffect[model] {
	logStore := NewStore("consumeEventLogStore")

	c := *originStore.counter
	return &ComposerAfterEffect[model]{
		eventCounter:          c,
		originStore:           originStore,
		effect:                effect,
		receiver:              receiver,
		consumedEventLogStore: logStore,
	}
}

func (c *ComposerAfterEffect[model]) Operate() {
	if c.eventCounter == *c.originStore.counter {
		return
	}

	for id, events := range c.originStore.GetEvents() {
		spawnLogEvents, err := c.consumedEventLogStore.GetEventsByEntityID(id)
		if err != nil {
			log.Fatalln(err)
		}

		lastSkipIndex := 0

		if len(spawnLogEvents) > 0 {
			lastSpawnData := spawnLogEvents[len(spawnLogEvents)-1].Data.(consumeData)
			for i := 0; i < len(events); i++ {
				event := events[i]
				if lastSpawnData.OriginEventID == event.ID {
					lastSkipIndex = i + 1
					break
				}
			}
		}

		for _, event := range events[lastSkipIndex:] {
			if string(c.effect) == event.Effect {
				for store, initEvents := range c.receiver.GetEvents(event.EntityID, event.Data.(model)) {
					for _, event := range initEvents {
						log.Printf("spawn, store: %s, [%s]id: %s\n", store.name, event.Effect, event.EntityID)
						err := store.AppendEvent(event)
						if err != nil {
							log.Fatalln(err)
						}
					}
				}

				err := c.consumedEventLogStore.AppendEvent(effectConsume.NewEvent(id, consumeData{
					OriginID:      id,
					OriginEventID: event.ID,
				}))
				if err != nil {
					log.Fatalln(err)
				}
			}
		}
	}

	c.eventCounter = *c.originStore.counter
}
