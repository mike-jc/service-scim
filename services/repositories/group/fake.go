package repositoriesGroup

import (
	"service-scim/errors/repositories"
	"service-scim/models/resources"
	"service-scim/resources/fakeResources"
)

type Fake struct {
	Abstract
}

func (f *Fake) ById(id string) (group *modelsResources.Group, err errorsRepositories.Interface) {
	return fakeResources.GroupOperatorsObject, nil
}

func (f *Fake) List(offset, limit int) (totalCount int, list []*modelsResources.Group, err errorsRepositories.Interface) {
	return fakeResources.GroupsObject.TotalResults, fakeResources.GroupsObject.Resources, nil
}

func (f *Fake) Count(filter map[string]interface{}, id *string) (count int, err errorsRepositories.Interface) {
	if externalId, ok := filter["externalId"]; ok {
		if fakeResources.GroupOperatorsObject.ExternalId == externalId {
			return 1, nil
		} else {
			return 0, nil
		}
	} else {
		return 0, errorsRepositories.NewError("Only know how to filter entity by 'name'", errorsRepositories.ApiError)
	}
}
