package system

import (
	"reflect"
	"strings"
)

func MapValueForField(data map[string]interface{}, val reflect.Value, format, fieldName string) interface{} {
	// look for key corresponding to field name
	val = ReflectValue(val)
	mapKey := ""
	if field, fExists := val.Type().FieldByName(fieldName); fExists {
		mapKey = strings.Split(field.Tag.Get(format), ",")[0]
	}

	// get map value for the key if it's not absent and not empty
	if len(mapKey) > 0 {
		if mapVal, vExists := data[mapKey]; vExists && !TypeIsEmpty(mapVal) {
			return mapVal
		}
	}
	return nil
}

func MapSliceOfMapsForField(data map[string]interface{}, val reflect.Value, format, fieldName string) []map[string]interface{} {
	if mapVal := MapValueForField(data, val, format, fieldName); mapVal != nil {
		switch mapVal.(type) {
		case []map[string]interface{}:
			return mapVal.([]map[string]interface{})
		case []interface{}:
			slice := mapVal.([]interface{})
			sliceOfMaps := make([]map[string]interface{}, len(slice))
			for i, item := range slice {
				switch item.(type) {
				case map[string]interface{}:
					sliceOfMaps[i] = item.(map[string]interface{})
				default:
					return nil
				}
			}
			return sliceOfMaps
		}
	}
	return nil
}

func MapValueIsEmpty(data map[string]interface{}, key string) bool {
	if val, exists := data[key]; exists {
		return TypeIsEmpty(val)
	}
	return true
}
