package middleware

import (
	sdksData "gitlab.com/24sessions/sdk-go-configurator/data"
	"reflect"
	"service-scim/system"
	"strings"
)

type Interface interface {
	SetFormat(format string)
	SetConfig(config *sdksData.ScimContainer)
	ProcessUser(data map[string]interface{}, id *string) (map[string]interface{}, error)
	ProcessGroup(data map[string]interface{}, id *string) (map[string]interface{}, error)
}

type Abstract struct {
	format string
	config *sdksData.ScimContainer
}

func (a *Abstract) SetFormat(format string) {
	a.format = format
}

func (a *Abstract) SetConfig(config *sdksData.ScimContainer) {
	a.config = config
}

func (a *Abstract) ProcessUser(data map[string]interface{}, id *string) (map[string]interface{}, error) {
	return data, nil
}

func (a *Abstract) ProcessGroup(data map[string]interface{}, id *string) (map[string]interface{}, error) {
	return data, nil
}

// Gives key name for the given struct's field and format
// basing on the information from the corresponding field's tag
func (a *Abstract) mapKeyForField(val reflect.Value, fieldName string) string {
	val = system.ReflectValue(val)
	switch val.Kind() {
	case reflect.Struct:
		if field, ok := val.Type().FieldByName(fieldName); ok {
			if key := strings.Split(field.Tag.Get(a.format), ",")[0]; len(key) > 0 {
				return key
			}
		}
	}
	return ""
}
