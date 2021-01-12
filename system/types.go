package system

import (
	"reflect"
	"strings"
)

func StructToMap(data interface{}, tagWithFieldName string) interface{} {
	val := reflect.ValueOf(data)

	switch val.Kind() {
	case reflect.Struct:
		m := make(map[string]interface{})
		for i := 0; i < val.NumField(); i++ {
			f := val.Type().Field(i)
			fName := ""
			if mapKey := strings.Split(f.Tag.Get(tagWithFieldName), ",")[0]; len(mapKey) > 0 && mapKey != "-" {
				fName = mapKey
			}

			if fName != "" && val.Field(i).CanInterface() {
				m[fName] = StructToMap(val.Field(i).Interface(), tagWithFieldName)
			}
		}
		return m
	case reflect.Map:
		m := make(map[string]interface{})
		for _, k := range val.MapKeys() {
			m[k.Interface().(string)] = StructToMap(val.MapIndex(k).Interface(), tagWithFieldName)
		}
		return m
	case reflect.Ptr, reflect.Interface:
		if val.IsValid() && !val.IsNil() {
			return StructToMap(val.Elem().Interface(), tagWithFieldName)
		} else {
			return nil
		}
	case reflect.Array, reflect.Slice:
		s := make([]interface{}, 0)
		for i := 0; i < val.Len(); i++ {
			s = append(s, StructToMap(val.Index(i).Interface(), tagWithFieldName))
		}
		return s
	default:
		return data
	}
}

func TypeIsEmpty(obj interface{}) bool {
	val := ReflectValue(obj)
	switch val.Kind() {
	case reflect.Invalid:
		return true
	case reflect.Map, reflect.Array, reflect.Slice:
		return val.Len() == 0
	case reflect.Bool:
		return false
	default:
		return obj == reflect.Zero(val.Type()).Interface()
	}
}

func ToString(val interface{}) string {
	switch val.(type) {
	case string:
		return val.(string)
	default:
		return ""
	}
}

func ToBool(val interface{}) bool {
	switch val.(type) {
	case bool:
		return val.(bool)
	default:
		return false
	}
}
