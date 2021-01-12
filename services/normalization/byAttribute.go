package normalization

import (
	"fmt"
	"reflect"
	"service-scim/errors"
	"service-scim/models/config"
	"service-scim/models/normalization"
	"service-scim/services/navigation"
	"service-scim/system"
	"time"
)

type ByAttribute struct {
}

// Normalization of entire given data by scanning all given attributes
func (a *ByAttribute) Normalize(data interface{}, attributes []*modelsConfig.Attribute, included, excluded []navigation.Path) (reflect.Value, error) {
	emptyValue := reflect.Value{}
	if attributes == nil || len(attributes) == 0 {
		return emptyValue, nil
	}

	reflectItems := make([]*modelsNormalization.ReflectStructItem, 0)

	// Check each attribute from given scheme
	for _, attr := range attributes {
		// Shouldn't be returned, skip
		if !attr.IsReturned(included, excluded) {
			continue
		}

		// Get data's item by attribute (by path to key/field name)
		if field, val, err := attr.Item(data); err != nil {
			return emptyValue, err
		} else if len(field.Name) == 0 || !val.IsValid() {
			// no value found in data for the attribute, skip it
			continue
		} else {
			if attr.MultiValued {
				val = system.ReflectValue(val)
				// For multi-valued item we can store only multi-valued values
				switch val.Kind() {
				case reflect.Array, reflect.Slice:
					// Can be pointer (or interface), get value in that case
					var elem reflect.Value
					if valArray := system.ReflectArray(val); valArray.IsValid() {
						elem = valArray
					} else if valSlice := system.ReflectSlice(val); valSlice.IsValid() {
						elem = valSlice
					} else {
						continue
					}

					// For multi-valued items normalize all their values
					var slice reflect.Value
					for i := 0; i < elem.Len(); i++ {
						if normalized, err := a.NormalizeValue(elem.Index(i), attr, included, excluded); err != nil {
							return emptyValue, err
						} else if normalized.IsValid() {
							if !slice.IsValid() {
								slice = reflect.Zero(reflect.SliceOf(normalized.Type()))
							}
							slice = reflect.Append(slice, normalized)
						}
					}
					if slice.IsValid() {
						reflectItems = append(reflectItems, modelsNormalization.MakeReflectStructItem(field, slice))
					}
				case reflect.Invalid:
					continue
				default:
					return emptyValue, scimErrors.MultiValueError(val, attr)
				}
			} else {
				// For ordinal items just normalize value
				if normalized, err := a.NormalizeValue(val, attr, included, excluded); err != nil {
					return emptyValue, err
				} else if normalized.IsValid() {
					reflectItems = append(reflectItems, modelsNormalization.MakeReflectStructItem(field, normalized))
				}
			}
		}
	}

	return a.newValue(data, reflectItems)
}

// Normalize one value according to the given attribute's properties
func (a *ByAttribute) NormalizeValue(val reflect.Value, attr *modelsConfig.Attribute, included, excluded []navigation.Path) (reflect.Value, error) {
	emptyValue := reflect.Value{}
	val = system.ReflectValue(val)

	switch attr.Type {
	case "complex":
		// For complex item we can store only complex value or its pointer
		switch val.Kind() {
		case reflect.Struct, reflect.Map, reflect.Ptr, reflect.Interface:
			// Can be pointer (or interface), get value in that case
			var elem reflect.Value
			if valStruct := system.ReflectStruct(val); valStruct.IsValid() {
				elem = valStruct
			} else if valMap := system.ReflectMap(val); valMap.IsValid() {
				elem = valMap
			} else {
				return emptyValue, nil
			}

			// For struct/map we should recursively scan all sub-attributes
			if normalized, err := a.Normalize(elem.Interface(), attr.SubAttributes, included, excluded); err != nil {
				return emptyValue, err
			} else {
				return normalized, nil
			}
		case reflect.Invalid:
			return emptyValue, nil
		default:
			return emptyValue, scimErrors.InvalidValueTypeError(val, attr)
		}
	case "string", "binary", "decimal", "integer", "reference":
		// For the most of basic types do nothing
		return val, nil
	case "boolean":
		// Boolean should store boolean
		switch val.Kind() {
		case reflect.Bool:
			return val, nil
		case reflect.Invalid:
			return emptyValue, nil
		default:
			return emptyValue, scimErrors.InvalidValueTypeError(val, attr)
		}
	case "datetime":
		// Datetime should store date/time
		switch val.Interface().(type) {
		case string, time.Time:
			return val, nil
		default:
			return emptyValue, scimErrors.InvalidValueTypeError(val, attr)
		}
	default:
		return val, nil
	}
}

// Build new data (struct or map) basing on collected values
func (a *ByAttribute) newValue(data interface{}, reflectItems []*modelsNormalization.ReflectStructItem) (reflect.Value, error) {
	emptyValue := reflect.Value{}
	dataKind := system.ReflectValue(data).Kind()

	switch dataKind {
	case reflect.Struct:
		return a.newStruct(data, reflectItems)
	case reflect.Map:
		return a.newMap(data, reflectItems)
	default:
		return emptyValue, fmt.Errorf("God %s. Only struct or map can be normalized.", dataKind)
	}
}

// Build new struct if original data is struct
func (a *ByAttribute) newStruct(data interface{}, reflectItems []*modelsNormalization.ReflectStructItem) (reflect.Value, error) {
	if len(reflectItems) == 0 {
		return reflect.ValueOf(struct{}{}), nil
	}

	fields := make([]reflect.StructField, len(reflectItems))
	for i, item := range reflectItems {
		fields[i] = item.Field
	}

	newStruct := reflect.New(reflect.StructOf(fields)).Elem()
	for i, item := range reflectItems {
		newStruct.Field(i).Set(item.Value)
	}
	return newStruct, nil
}

// Build new map if original data is map.
// If collected values are of different types, build struct
func (a *ByAttribute) newMap(data interface{}, reflectItems []*modelsNormalization.ReflectStructItem) (reflect.Value, error) {
	if len(reflectItems) == 0 {
		return reflect.ValueOf(map[string]interface{}{}), nil
	}

	if mapElemType, equal := modelsNormalization.EqualStructItemType(reflectItems); equal {
		mapKeyType := reflect.TypeOf(string(""))
		newMap := reflect.MakeMap(reflect.MapOf(mapKeyType, mapElemType))
		for _, item := range reflectItems {
			newMap.SetMapIndex(reflect.ValueOf(item.MapKeyName), item.Value)
		}
		return newMap, nil
	} else {
		return a.newStruct(data, reflectItems)
	}
}
