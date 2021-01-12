package validation

import (
	"reflect"
	"service-scim/models/resources"
	"service-scim/system"
	"strings"
)

type ReadOnlyAttribute struct {
	Abstract

	format        string
	reflectEntity reflect.Value
}

func (v *ReadOnlyAttribute) ValidateUser(data map[string]interface{}, format string, id *string) error {
	v.format = format
	v.reflectEntity = system.ReflectValue(new(modelsResources.User))

	v.validateId(data, id)
	v.validateMeta(data, "User", id)
	v.validateUserGroups(data)
	return nil
}

func (v *ReadOnlyAttribute) ValidateGroup(data map[string]interface{}, format string, id *string) error {
	v.format = format
	v.reflectEntity = system.ReflectValue(new(modelsResources.Group))

	v.validateId(data, id)
	v.validateMeta(data, "Group", id)
	return nil
}

// Remove ID value from data that's supposed to be new entity (for POST request).
// Keep current ID value for data of existing entity.
// ID field is determined by field's tags in `s` structure
func (v *ReadOnlyAttribute) validateId(data map[string]interface{}, id *string) {
	if key := v.mapKeyForField("Id"); len(key) > 0 {
		if _, ok := data[key]; ok && id == nil {
			delete(data, key)
		} else if id != nil {
			data[key] = *id
		}
	}
}

func (v *ReadOnlyAttribute) validateMeta(data map[string]interface{}, resourceType string, id *string) {
	if id == nil {
		return
	}
	if key := v.mapKeyForField("Meta"); len(key) > 0 {
		// TODO: here may be done refilling of the information of the entity meta
	}
}

func (v *ReadOnlyAttribute) validateUserGroups(data map[string]interface{}) {
	if key := v.mapKeyForField("Groups"); len(key) > 0 {
		// TODO: here may be done refilling of the information of the user groups
	}
}

func (v *ReadOnlyAttribute) mapKeyForField(fieldName string) string {
	val := system.ReflectValue(v.reflectEntity)
	switch val.Kind() {
	case reflect.Struct:
		if field, ok := val.Type().FieldByName(fieldName); ok {
			if key := strings.Split(field.Tag.Get(v.format), ",")[0]; len(key) > 0 {
				return key
			}
		}
	}
	return ""
}
