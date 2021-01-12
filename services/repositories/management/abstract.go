package repositoriesManagement

import (
	"service-scim/errors/repositories"
)

type Interface interface {
	Init(url string) errorsRepositories.Interface
	Ping() errorsRepositories.Interface
}

type Abstract struct {
	Interface
}

func (a *Abstract) Init(url string) errorsRepositories.Interface {
	return nil
}

func (a *Abstract) Ping() errorsRepositories.Interface {
	return nil
}
