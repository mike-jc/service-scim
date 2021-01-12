package repositoriesManagement

import (
	"service-scim/errors/repositories"
	"service-scim/sdks/restApi"
)

type Keeper struct {
	Abstract

	client *restApi.Keeper
}

func (k *Keeper) Init(url string) errorsRepositories.Interface {
	k.client = new(restApi.Keeper)
	k.client.SetBaseUrl(url)
	return nil
}

func (k *Keeper) Ping() errorsRepositories.Interface {
	if err := k.client.Ping(); err != nil {
		return errorsRepositories.NewError(err.Error(), errorsRepositories.ApiError)
	}
	return nil
}
