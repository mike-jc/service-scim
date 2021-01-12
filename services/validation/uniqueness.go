package validation

import (
	"reflect"
	"service-scim/errors"
	"service-scim/models/config"
	"service-scim/services/repositories/entity"
	"service-scim/system"
)

type Uniqueness struct {
	Abstract

	repository repositoriesEntity.EntityInterface
}

func (v *Uniqueness) SetRepository(repository repositoriesEntity.EntityInterface) {
	v.repository = repository
}

func (v *Uniqueness) Validate(data interface{}, id *string, schema *modelsConfig.Schema) error {
	return v.validateValue(system.ReflectValue(data), id, schema.ToAttribute())
}

func (v *Uniqueness) validateValue(val reflect.Value, id *string, attr *modelsConfig.Attribute) error {
	switch val.Kind() {
	case reflect.Map:
		for _, subAttr := range attr.SubAttributes {
			mapVal := system.ReflectValue(val.MapIndex(reflect.ValueOf(subAttr.Name)))
			if !system.ReflectValueIsAssigned(mapVal) {
				continue
			}

			switch subAttr.Uniqueness {
			case "server", "global":
				filter := map[string]interface{}{
					subAttr.Navigation.Path: mapVal.Interface(),
				}
				if count, cErr := v.repository.Count(filter, id); cErr != nil {
					return cErr
				} else if count > 0 {
					return scimErrors.DuplicationError(mapVal.Interface(), subAttr.Navigation.Path)
				}
			}

			if subAttr.ExpectsComplex() && mapVal.Kind() == reflect.Map {
				if vErr := v.validateValue(mapVal, id, subAttr); vErr != nil {
					return vErr
				}
			}
		}
	}
	return nil
}
