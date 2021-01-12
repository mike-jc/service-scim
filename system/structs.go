package system

import (
	"reflect"
	"strings"
)

// Gets keys for the given struct field, basing on its tag info
func MapKeysForStructField(field reflect.StructField) []string {
	keys := []string{field.Name}
	for _, tagName := range []string{"json", "xml"} {
		if key := strings.Split(field.Tag.Get(tagName), ",")[0]; len(key) > 0 {
			keys = append(keys, key)
		}
	}
	return keys
}

// Gets value of corresponding map's key for the given struct field,
// basing on its name or tag info. Returns nil if nothing is found
func MapValueForStructField(field reflect.StructField, data map[string]interface{}) interface{} {
	// collect possible map keys
	keys := MapKeysForStructField(field)
	// get value of the first existing key
	for _, key := range keys {
		if val, ok := data[key]; ok {
			return val
		}
	}
	return nil
}

// Gets keys for the given struct field name, basing on its name or tag info
func MapKeysForStruct(s interface{}, fieldName string) []string {
	keys := []string{fieldName}
	if field, fOk := ReflectValue(s).Type().FieldByName(fieldName); fOk {
		for _, tagName := range []string{"json", "xml"} {
			if key := strings.Split(field.Tag.Get(tagName), ",")[0]; len(key) > 0 {
				keys = append(keys, key)
			}
		}
	}
	return keys
}

// Gets value of corresponding map's key basing on struct's field name or field's tag info.
// Returns nil if nothing is found
func MapValueForStruct(s interface{}, data map[string]interface{}, fieldName string) interface{} {
	// collect possible map keys
	keys := MapKeysForStruct(s, fieldName)
	// get value of the first existing key
	for _, key := range keys {
		if val, ok := data[key]; ok {
			return val
		}
	}
	return nil
}

// Gets value (slice of maps) of corresponding map's key basing on struct's field name or field's tag info.
// Returns empty slice if nothing is found
func SliceOfMapsForStruct(s interface{}, data map[string]interface{}, fieldName string) []map[string]interface{} {
	if val := MapValueForStruct(s, data, fieldName); val != nil {
		if slice, ok := val.([]interface{}); ok {
			sliceOfMaps := make([]map[string]interface{}, 0)
			for _, item := range slice {
				if itemAsMap, ok := item.(map[string]interface{}); ok {
					sliceOfMaps = append(sliceOfMaps, itemAsMap)
				}
			}
			return sliceOfMaps
		}
	}
	return []map[string]interface{}{}
}

// Checks if struct is empty, i.e. all its fields have zero value
func StructIsEmpty(s interface{}) bool {
	val := ReflectValue(s)

	switch val.Kind() {
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			fieldVal := val.Field(i)
			if fieldVal.Interface() != reflect.Zero(fieldVal.Type()).Interface() {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func SetStructFieldByName(s interface{}, fieldName string, newVal interface{}) {
	sVal := ReflectValue(s)

	switch sVal.Kind() {
	case reflect.Struct:
		sVal.FieldByName(fieldName).Set(ReflectValue(newVal))
	}
}

func SetStructStringFields(s interface{}, data map[string]interface{}, strFields []string, skipNil bool) {
	for _, fieldName := range strFields {
		if val := MapValueForStruct(s, data, fieldName); val != nil || !skipNil {
			strVal, _ := val.(string)
			SetStructFieldByName(s, fieldName, strVal)
		}
	}
}

func SetStructBoolFields(s interface{}, data map[string]interface{}, bFields []string, skipNil bool) {
	for _, fieldName := range bFields {
		if val := MapValueForStruct(s, data, fieldName); val != nil || !skipNil {
			bVal, _ := val.(bool)
			SetStructFieldByName(s, fieldName, bVal)
		}
	}
}

func SetStructInt64Fields(s interface{}, data map[string]interface{}, iFields []string, skipNil bool) {
	for _, fieldName := range iFields {
		if val := MapValueForStruct(s, data, fieldName); val != nil || !skipNil {
			bVal, _ := val.(int64)
			SetStructFieldByName(s, fieldName, bVal)
		}
	}
}

func StructFilterPassed(s interface{}, filter map[string]interface{}) bool {
	if len(filter) == 0 {
		return true
	}

	val := ReflectValue(s)
	switch val.Kind() {
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			sVal := val.Field(i)
			fVal := MapValueForStructField(val.Type().Field(i), filter)
			if fVal != nil && fVal != sVal.Interface() {
				return false
			}
		}
	}
	return true
}
