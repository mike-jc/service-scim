package system

import (
	"reflect"
)

func ReflectValue(data interface{}) reflect.Value {
	var val reflect.Value

	// get reflect value if necessary
	switch data.(type) {
	case reflect.Value:
		val = data.(reflect.Value)
	default:
		val = reflect.ValueOf(data)
	}

	// get value from pointer if necessary
	switch val.Kind() {
	case reflect.Ptr, reflect.Interface:
		val = val.Elem()
	}

	return val
}

// Checks value if it's array. Returns invalid value if not.
// For pointer tries to check the value it points to.
// For interface tries to cast type to array.
func ReflectArray(val reflect.Value) reflect.Value {
	switch val.Kind() {
	case reflect.Array:
		return val
	case reflect.Ptr, reflect.Interface:
		return ReflectArray(val.Elem())
	default:
		return reflect.Value{}
	}
}

// Checks value if it's slice. Returns invalid value if not.
// For pointer tries to check the value it points to.
// For interface tries to cast type to slice.
func ReflectSlice(val reflect.Value) reflect.Value {
	switch val.Kind() {
	case reflect.Slice:
		return val
	case reflect.Ptr, reflect.Interface:
		return ReflectSlice(val.Elem())
	default:
		return reflect.Value{}
	}
}

// Checks value if it's struct. Returns invalid value if not.
// For pointer tries to check the value it points to.
// For interface tries to cast type to struct.
func ReflectStruct(val reflect.Value) reflect.Value {
	switch val.Kind() {
	case reflect.Struct:
		return val
	case reflect.Ptr, reflect.Interface:
		return ReflectStruct(val.Elem())
	default:
		return reflect.Value{}
	}
}

// Checks value if it's map. Returns invalid value if not.
// For pointer tries to check the value it points to.
// For interface tries to cast type to map.
func ReflectMap(val reflect.Value) reflect.Value {
	switch val.Kind() {
	case reflect.Map:
		return val
	case reflect.Ptr, reflect.Interface:
		return ReflectMap(val.Elem())
	default:
		return reflect.Value{}
	}
}

func ReflectValueIsAssigned(val reflect.Value) bool {
	val = ReflectValue(val)
	if !val.IsValid() {
		return false
	}

	switch val.Kind() {
	case reflect.String, reflect.Map, reflect.Array, reflect.Slice:
		return val.Len() > 0
	default:
		return true
	}
}

func ReflectSafeIsNil(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return val.IsNil()
	default:
		return !val.IsValid() || val.IsNil()
	}
}
