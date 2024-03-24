package state

import (
	"log"
	"reflect"

	"github.com/google/uuid"
)

func FieldValue[T interface{}](pm ProjectorManager, id uuid.UUID, entityName, fieldName string) T {
	projector := pm.GetEntityProjector(entityName)

	data, ok := projector.Project(id, fieldName).(T)
	if !ok {
		log.Fatalf("failed to parse data for effect: %s", reflect.TypeOf(data))
	}
	return data
}

func FieldValues[T interface{}](pm ProjectorManager, entityName, fieldName string) ([]T, []uuid.UUID) {
	projector := pm.GetEntityProjector(entityName)

	ids := projector.ListIdentifiers()
	results := make([]T, len(ids))

	for i, id := range ids {
		results[i] = FieldValue[T](pm, id, entityName, fieldName)
	}

	return results, ids
}
