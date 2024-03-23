package state

import "github.com/google/uuid"

type Projector interface {
	Project(identifier uuid.UUID, field string) interface{}
	ListIdentifiers() []uuid.UUID
}

type ProjectorManager struct {
	EntityProjectorMap map[string]Projector
}

func (m ProjectorManager) GetEntityProjector(entityType string) Projector {
	return m.EntityProjectorMap[entityType]
}
