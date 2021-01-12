package validation

import (
	"service-scim/models/config"
)

type Interface interface {
	Validate(data interface{}, schema *modelsConfig.Schema) error
}

type Abstract struct {
	Interface
}

func (a *Abstract) Validate(data interface{}, schema *modelsConfig.Schema) error {
	return nil
}
