package validation

import (
	"fmt"
	"reflect"
	"service-scim/errors"
	"service-scim/models/config"
	"service-scim/services/navigation"
	"service-scim/system"
)

type AttributeType struct {
	Abstract
}

func (v *AttributeType) Validate(data interface{}, schema *modelsConfig.Schema) error {
	return v.validateValue(system.ReflectValue(data), schema.ToAttribute())
}

func (v *AttributeType) validateValue(val reflect.Value, attr *modelsConfig.Attribute) error {
	if attr.Mutability == "readOnly" {
		return nil
	}

	val = system.ReflectValue(val)
	if !val.IsValid() {
		return nil
	}

	switch val.Kind() {
	case reflect.String:
		if !attr.ExpectsString() {
			return scimErrors.InvalidAttrTypeError(attr, val.Kind().String())
		}

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		if !attr.ExpectsInteger() {
			return scimErrors.InvalidAttrTypeError(attr, val.Kind().String())
		}

	case reflect.Float32, reflect.Float64:
		if !attr.ExpectsFloat() {
			return scimErrors.InvalidAttrTypeError(attr, val.Kind().String())
		}

	case reflect.Bool:
		if !attr.ExpectsBool() {
			return scimErrors.InvalidAttrTypeError(attr, val.Kind().String())
		}

	case reflect.Array, reflect.Slice:
		if !attr.MultiValued {
			return scimErrors.InvalidAttrTypeError(attr, val.Kind().String())
		}

		subAttr := attr.Clone()
		subAttr.MultiValued = false
		for i := 0; i < val.Len(); i++ {
			if arrErr := v.validateValue(val.Index(i), subAttr); arrErr != nil {
				return arrErr
			}
		}

	case reflect.Map:
		if !attr.ExpectsComplex() {
			return scimErrors.InvalidAttrTypeError(attr, val.Kind().String())
		}

		for _, k := range val.MapKeys() {
			if p, pErr := navigation.NewPath(k.String()); pErr != nil {
				return pErr
			} else if subAttr := attr.GetAttribute(p, false); subAttr == nil {
				continue // just ignore redundant attributes in input data
				//return scimErrors.NoAttributeError(p.Value())
			} else {
				if mapErr := v.validateValue(val.MapIndex(k), subAttr); mapErr != nil {
					return mapErr
				}
			}
		}

	default:
		return scimErrors.InvalidAttrTypeError(attr, fmt.Sprintf("unhandled type (%s)", val.Kind().String()))
	}
	return nil
}
