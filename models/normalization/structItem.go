package modelsNormalization

import (
	"fmt"
	"reflect"
	"service-scim/system"
)

type ReflectStructItem struct {
	Field      reflect.StructField
	Value      reflect.Value
	MapKeyName string
}

func MakeReflectStructItem(f reflect.StructField, v reflect.Value) *ReflectStructItem {
	upperFirstFieldName := system.UpperCaseFirst(f.Name)
	lowerFirstFieldName := system.LowerCaseFirst(f.Name)

	tag := f.Tag
	tagStr := string(tag)
	if len(tag.Get("json")) == 0 {
		tagStr += fmt.Sprintf(` json:"%s"`, lowerFirstFieldName)
	}
	if len(tag.Get("xml")) == 0 {
		tagStr += fmt.Sprintf(` xml:"%s"`, upperFirstFieldName)
	}
	if len(tag.Get("map")) == 0 {
		tagStr += fmt.Sprintf(` map:"%s"`, f.Name)
	}
	if len(tagStr) > 0 {
		tag = reflect.StructTag(tagStr)
	}

	return &ReflectStructItem{
		Field: reflect.StructField{
			Name: upperFirstFieldName, // to make struct field exported
			Type: v.Type(),
			Tag:  tag,
		},
		Value:      system.ReflectValue(v),
		MapKeyName: f.Name,
	}
}

// Gets type of first value in item list.
// The second returned value indicates if types of each value are equal
func EqualStructItemType(items []*ReflectStructItem) (reflect.Type, bool) {
	if len(items) == 0 {
		return nil, false
	}

	t := items[0].Value.Type()
	if len(items) == 1 {
		return t, true
	}

	for i := 1; i < len(items); i++ {
		if t != items[i].Value.Type() {
			return nil, false
		}
	}
	return t, true
}
