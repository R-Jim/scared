package state

import "github.com/google/uuid"

type Projector interface {
	Project(identifier uuid.UUID, field string) interface{}
	ListIdentifiers() []uuid.UUID
}

type ProjectorManager struct {
	entityProjectorMap map[string]Projector
}

func NewProjectorManager(projectorMapping map[string]Projector) ProjectorManager {
	return ProjectorManager{
		entityProjectorMap: projectorMapping,
	}
}

func (m ProjectorManager) GetEntityProjector(entityType string) Projector {
	return m.entityProjectorMap[entityType]
}
