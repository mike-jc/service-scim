package repositoriesEntity

import (
	"service-scim/errors/repositories"
)

type EntityInterface interface {
	Init(url string) errorsRepositories.Interface
	SetInstanceDomain(domain string)
	Count(filter map[string]interface{}, id *string) (count int, err errorsRepositories.Interface)
}

type AbstractEntity struct {
}

func (a *AbstractEntity) Init(url string) errorsRepositories.Interface {
	return nil
}

func (a *AbstractEntity) SetInstanceDomain(domain string) {
}

func (a *AbstractEntity) Count(filter map[string]interface{}, id *string) (count int, err errorsRepositories.Interface) {
	return 0, nil
}
