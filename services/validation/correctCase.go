package validation

import (
	"reflect"
	"service-scim/models/config"
	"service-scim/services/navigation"
	"service-scim/system"
)

type CorrectCase struct {
	Abstract
}

func (v *CorrectCase) Validate(data interface{}, schema *modelsConfig.Schema) error {
	return v.validateValue(system.ReflectValue(data), schema.ToAttribute())
}

func (v *CorrectCase) validateValue(val reflect.Value, attr *modelsConfig.Attribute) error {
	val = system.ReflectValue(val)
	if !val.IsValid() {
		return nil
	}

	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if arrErr := v.validateValue(val.Index(i), attr); arrErr != nil {
				return arrErr
			}
		}

	case reflect.Map:
		// fix case of map keys
		for _, key := range val.MapKeys() {
			if p, kErr := navigation.NewPath(key.String()); kErr != nil {
				continue // just ignore
				//return scimErrors.NoAttributeError(p.Value())
			} else {
				if subAttr := attr.GetAttribute(p, false); subAttr == nil {
					continue // just ignore redundant attributes in input data
					//return scimErrors.NoAttributeError(p.Value())
				} else {
					subVal := val.MapIndex(key)
					if key.String() != subAttr.Name {
						val.SetMapIndex(reflect.ValueOf(subAttr.Name), subVal)
						val.SetMapIndex(key, reflect.Value{})
					}

					if subErr := v.validateValue(subVal, subAttr); subErr != nil {
						return subErr
					}
				}
			}
		}
	}
	return nil
}
