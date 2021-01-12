package repositoriesUser

import (
	"service-scim/errors/repositories"
	"service-scim/models/resources"
	"service-scim/resources/fakeResources"
)

type Fake struct {
	Abstract
}

func (f *Fake) ById(id string) (user *modelsResources.User, err errorsRepositories.Interface) {
	return fakeResources.UserJackObject, nil
}

func (f *Fake) List(offset, limit int, filterMap map[string]interface{}) (totalCount int, list []*modelsResources.User, err errorsRepositories.Interface) {
	return fakeResources.UsersObject.TotalResults, fakeResources.UsersObject.Resources, nil
}

func (f *Fake) Count(filter map[string]interface{}, id *string) (count int, err errorsRepositories.Interface) {
	if externalId, ok := filter["externalId"]; ok {
		if externalId == fakeResources.UserJackObject.ExternalId {
			return 1, nil
		} else {
			return 0, nil
		}
	} else {
		return 0, errorsRepositories.NewError("Only know how to filter entity by 'email'", errorsRepositories.ApiError)
	}
}
