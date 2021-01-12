package validation

import (
	"reflect"
	"service-scim/errors"
	"service-scim/models/config"
	"service-scim/services/navigation"
	"service-scim/system"
)

type AttributeMutability struct {
	Abstract

	subjectStack   navigation.Stack
	referenceStack navigation.Stack
}

func (v *AttributeMutability) Validate(subject, reference interface{}, schema *modelsConfig.Schema) error {
	v.subjectStack = navigation.NewStackWithoutLimit()
	v.referenceStack = navigation.NewStackWithoutLimit()

	return v.stepIn(system.ReflectValue(subject), system.ReflectValue(reference), schema.ToAttribute())
}

func (v *AttributeMutability) stepIn(subjVal, refVal reflect.Value, attr *modelsConfig.Attribute) error {
	if !system.ReflectValueIsAssigned(subjVal) || !system.ReflectValueIsAssigned(refVal) {
		return nil
	}

	v.subjectStack.Push(system.ReflectValue(subjVal))
	v.referenceStack.Push(system.ReflectValue(refVal))

	err := v.validateBySubAttributes(attr)

	v.subjectStack.Pop()
	v.referenceStack.Pop()
	return err
}

func (v *AttributeMutability) validateBySubAttributes(attr *modelsConfig.Attribute) error {
	for _, subAttr := range attr.SubAttributes {
		subjVal := v.subjectValue(subAttr)
		refVal := v.referenceValue(subAttr)

		switch subAttr.Type {
		case "complex":
			if subAttr.MultiValued {
				// for multi-valued properties, first compare their own values
				if cErr := v.compareAndCopy(subjVal, refVal, subAttr); cErr != nil {
					return cErr
				}
				// then if those values are not empty,
				// compare all their sub-values (array items)
				if system.ReflectValueIsAssigned(subjVal) && system.ReflectValueIsAssigned(refVal) {
					subjVal = system.ReflectValue(subjVal)
					refVal = system.ReflectValue(refVal)

					singledAttr := subAttr.Clone()
					singledAttr.MultiValued = false

					// look for coincident items and compare their sub-values
					for i := 0; i < subjVal.Len(); i++ {
						for j := 0; j < refVal.Len(); j++ {
							subjItemVal := subjVal.Index(i)
							refItemVal := refVal.Index(j)
							if v.mapsAreEqualByKeys(subjItemVal, refItemVal, singledAttr.Navigation.IndexKeys) {
								if sErr := v.stepIn(subjItemVal, refItemVal, singledAttr); sErr != nil {
									return sErr
								}
							}
						}
					}
				}
			} else {
				// for complex values, first compare them
				if cErr := v.compareAndCopy(subjVal, refVal, subAttr); cErr != nil {
					return cErr
				}
				// then compare their properties' values
				if sErr := v.stepIn(subjVal, refVal, subAttr); sErr != nil {
					return sErr
				}
			}
		default:
			// for not-complex values, just compare them
			if cErr := v.compareAndCopy(subjVal, refVal, subAttr); cErr != nil {
				return cErr
			}
		}
	}
	return nil
}

func (v *AttributeMutability) mapValue(val reflect.Value, attr *modelsConfig.Attribute) reflect.Value {
	switch val.Kind() {
	case reflect.Map:
		return val.MapIndex(reflect.ValueOf(attr.Name))
	default:
		return reflect.Value{}
	}
}

func (v *AttributeMutability) subjectValue(attr *modelsConfig.Attribute) reflect.Value {
	return v.mapValue(v.subjectStack.Peek().(reflect.Value), attr)
}

func (v *AttributeMutability) referenceValue(attr *modelsConfig.Attribute) reflect.Value {
	return v.mapValue(v.referenceStack.Peek().(reflect.Value), attr)
}

func (v *AttributeMutability) compareAndCopy(subjVal, refVal reflect.Value, attr *modelsConfig.Attribute) error {
	switch attr.Mutability {
	case "readOnly":
		// restore original value if the property is read-only
		baseVal := v.subjectStack.Peek().(reflect.Value)
		baseVal.SetMapIndex(reflect.ValueOf(attr.Name), refVal)

	case "immutable":
		// for unchangeable property, the new value should be equal to the original
		if !system.ReflectSafeIsNil(refVal) {
			if !system.ReflectSafeIsNil(subjVal) {
				if !reflect.DeepEqual(subjVal.Interface(), refVal.Interface()) {
					return scimErrors.MutabilityViolationError(attr.Navigation.FullPath)
				}
			} else {
				return scimErrors.MutabilityViolationError(attr.Navigation.FullPath)
			}
		}
	}
	return nil
}

func (v *AttributeMutability) mapsAreEqualByKeys(subjVal, refVal reflect.Value, keys []string) bool {
	if !subjVal.IsValid() || !refVal.IsValid() {
		return false
	}

	subjVal = system.ReflectValue(subjVal)
	refVal = system.ReflectValue(refVal)

	for _, key := range keys {
		subjKeyVal := subjVal.MapIndex(reflect.ValueOf(key))
		refKeyVal := refVal.MapIndex(reflect.ValueOf(key))

		if !subjKeyVal.IsValid() || !refKeyVal.IsValid() {
			return false
		} else if !reflect.DeepEqual(subjKeyVal.Interface(), refKeyVal.Interface()) {
			return false
		}
	}
	return true
}
