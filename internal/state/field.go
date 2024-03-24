package state

import (
	"fmt"
	"reflect"

	"github.com/google/uuid"
	pkgerrors "github.com/pkg/errors"
)

func FieldValue[T interface{}](pm ProjectorManager, id uuid.UUID, entityName, fieldName string) (T, error) {
	projector := pm.GetEntityProjector(entityName)

	data, ok := projector.Project(id, fieldName).(T)
	if !ok {
		return data, pkgerrors.WithStack(fmt.Errorf("failed to parse data for effect: %s", reflect.TypeOf(data)))
	}
	return data, nil
}

func FieldValues[T interface{}](pm ProjectorManager, entityName, fieldName string) ([]T, []uuid.UUID, error) {
	projector := pm.GetEntityProjector(entityName)

	ids := projector.ListIdentifiers()
	results := make([]T, len(ids))

	for i, id := range ids {
		data, err := FieldValue[T](pm, id, entityName, fieldName)
		if err != nil {
			return nil, nil, err
		}

		results[i] = data
	}

	return results, ids, nil
}
