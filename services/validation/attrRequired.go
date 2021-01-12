package validation

import (
	"reflect"
	"service-scim/errors"
	"service-scim/models/config"
	"service-scim/system"
)

type RequiredAttribute struct {
	Abstract
}

func (v *RequiredAttribute) Validate(data interface{}, schema *modelsConfig.Schema) error {
	return v.validateValue(system.ReflectValue(data), schema.ToAttribute())
}

func (v *RequiredAttribute) validateValue(val reflect.Value, attr *modelsConfig.Attribute) error {
	val = system.ReflectValue(val)
	if !val.IsValid() {
		return v.checkValue(val, attr)
	}

	switch val.Kind() {
	case reflect.String:
		return v.checkValue(val, attr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.checkValue(val, attr)
	case reflect.Float32, reflect.Float64:
		return v.checkValue(val, attr)
	case reflect.Bool:
		return v.checkValue(val, attr)

	case reflect.Slice, reflect.Array:
		if arrErr := v.checkValue(val, attr); arrErr != nil {
			return arrErr
		}
		if attr.ExpectsComplexArray() {
			subAttr := attr.Clone()
			subAttr.MultiValued = false
			for i := 0; i < val.Len(); i++ {
				if subErr := v.validateValue(val.Index(i), subAttr); subErr != nil {
					return subErr
				}
			}
		}

	case reflect.Map:
		if mapErr := v.checkValue(val, attr); mapErr != nil {
			return mapErr
		}
		for _, subAttr := range attr.SubAttributes {
			if subVal, subErr := subAttr.MapItem(val); subErr != nil {
				return subErr
			} else {
				if subErr := v.validateValue(subVal, subAttr); subErr != nil {
					return subErr
				}
			}
		}
	}
	return nil
}

func (v *RequiredAttribute) checkValue(val reflect.Value, attr *modelsConfig.Attribute) error {
	// required property should be assigned with non-zero value
	if attr.Required && !system.ReflectValueIsAssigned(val) {
		switch attr.Mutability {
		case "readOnly":
			// property that is required, but readOnly, is allowed to be unassigned
		case "immutable":
			// property that is required, but immutable (non-changeable),
			// is allowed to have zero value (meaning that it has such value initially)
			if val.IsValid() {
				return scimErrors.MissingRequiredPropertyError(attr.Navigation.FullPath)
			}
		default:
			return scimErrors.MissingRequiredPropertyError(attr.Navigation.FullPath)
		}
	}
	return nil
}
