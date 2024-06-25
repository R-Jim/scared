package engine

import (
	"log"

	"github.com/google/uuid"
)

const (
	EffectSpawn Effect = "Spawn"
)

type spawnData struct {
	OriginID      uuid.UUID
	OriginEventID uuid.UUID
}

type Spawner interface {
	GetInitEvents(eventData interface{}) map[string]Event
	GetStore(store string) *Store
}

type ComposerSpawner struct {
	eventCounter int

	spawnOriginStore *Store
	spawnEffect      Effect
	spawner          Spawner
	spawnLogStore    *Store
}

func NewComposerSpawner(spawnerStore *Store, spawnEffect Effect, spawner Spawner) *ComposerSpawner {
	logStore := NewStore()

	c := *spawnerStore.counter
	return &ComposerSpawner{
		eventCounter:     c,
		spawnOriginStore: spawnerStore,
		spawnEffect:      spawnEffect,
		spawner:          spawner,
		spawnLogStore:    logStore,
	}
}

func (c *ComposerSpawner) Operate() {
	if c.eventCounter == *c.spawnOriginStore.counter {
		return
	}

	for id, events := range c.spawnOriginStore.GetEvents() {
		spawnLogEvents, err := c.spawnLogStore.GetEventsByEntityID(id)
		if err != nil {
			log.Fatalln(err)
		}

		lastSkipIndex := 0

		if len(spawnLogEvents) > 0 {
			lastSpawnData := spawnLogEvents[len(spawnLogEvents)-1].Data.(spawnData)
			for i := 0; i < len(events); i++ {
				event := events[i]
				if lastSpawnData.OriginEventID == event.ID {
					lastSkipIndex = i + 1
					break
				}
			}
		}

		for _, event := range events[lastSkipIndex:] {
			if event.Effect == c.spawnEffect {
				for store, initEvent := range c.spawner.GetInitEvents(event.Data) {
					err := c.spawner.GetStore(store).AppendEvent(initEvent)
					if err != nil {
						log.Fatalln(err)
					}
				}

				err := c.spawnLogStore.AppendEvent(initEvent(EffectSpawn, id, spawnData{
					OriginID:      id,
					OriginEventID: event.ID,
				}))
				if err != nil {
					log.Fatalln(err)
				}
			}
		}
	}

	c.eventCounter = *c.spawnOriginStore.counter
}
