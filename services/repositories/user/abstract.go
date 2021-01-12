package repositoriesUser

import (
	"reflect"
	"service-scim/errors/repositories"
	"service-scim/models/resources"
	"service-scim/services/repositories/entity"
	"service-scim/system"
)

type Interface interface {
	repositoriesEntity.EntityInterface

	SetFormat(format string)

	List(offset, limit int, filterMap map[string]interface{}) (totalCount int, list []*modelsResources.User, err errorsRepositories.Interface)
	ById(id string) (user *modelsResources.User, err errorsRepositories.Interface)
	Create(data map[string]interface{}) (resultedUser *modelsResources.User, err errorsRepositories.Interface)
	Update(id string, data map[string]interface{}) (resultedUser *modelsResources.User, err errorsRepositories.Interface)
	Replace(id string, data map[string]interface{}) (resultedUser *modelsResources.User, err errorsRepositories.Interface)
	Block(id string) errorsRepositories.Interface
	Search() errorsRepositories.Interface
	Bulk() errorsRepositories.Interface
}

type Abstract struct {
	format string

	repositoriesEntity.AbstractEntity
}

func (a *Abstract) SetFormat(format string) {
	a.format = format
}

func (a *Abstract) List(offset, limit int, filterMap map[string]interface{}) (totalCount int, list []*modelsResources.User, err errorsRepositories.Interface) {
	return 0, nil, nil
}

func (a *Abstract) ById(id string) (user *modelsResources.User, err errorsRepositories.Interface) {
	return nil, nil
}

func (a *Abstract) Create(data map[string]interface{}) (resultedUser *modelsResources.User, err errorsRepositories.Interface) {
	return nil, nil
}

func (a *Abstract) Update(id string, data map[string]interface{}) (resultedUser *modelsResources.User, err errorsRepositories.Interface) {
	return nil, nil
}

func (a *Abstract) Replace(id string, data map[string]interface{}) (resultedUser *modelsResources.User, err errorsRepositories.Interface) {
	return nil, nil
}

func (a *Abstract) Block(id string) errorsRepositories.Interface {
	return nil
}

func (a *Abstract) Search() errorsRepositories.Interface {
	return nil
}

func (a *Abstract) Bulk() errorsRepositories.Interface {
	return nil
}

func (a *Abstract) mapValueForField(data map[string]interface{}, val reflect.Value, fieldName string) interface{} {
	return system.MapValueForField(data, val, a.format, fieldName)
}

func (a *Abstract) mapSliceOfMapsForField(data map[string]interface{}, val reflect.Value, fieldName string) []map[string]interface{} {
	return system.MapSliceOfMapsForField(data, val, a.format, fieldName)
}
