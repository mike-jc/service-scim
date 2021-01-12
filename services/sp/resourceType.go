package serviceSp

import (
	"errors"
	"fmt"
	"gitlab.com/24sessions/sdk-go-configurator/data"
	"service-scim/models/config"
	"service-scim/resources/resourceTypes"
	"service-scim/system"
	"strings"
)

const ResourceTypeUnknown = 0
const ResourceTypeUser = 1
const ResourceTypeGroup = 2

var ResourceTypes = map[int]string{
	ResourceTypeUser:  "User",
	ResourceTypeGroup: "Group",
}

type ResourceType struct {
	scimConfig *sdksData.ScimContainer
	baseUrl    string
}

func (t *ResourceType) SetScimConfig(scimConfig *sdksData.ScimContainer) {
	t.scimConfig = scimConfig
}

func (t *ResourceType) SetBaseUrl(baseUrl string) {
	t.baseUrl = baseUrl
}

func (t *ResourceType) IdFromName(givenName string) (id int, err error) {
	givenName = strings.ToLower(givenName)

	var name string
	for id, name = range ResourceTypes {
		if strings.ToLower(name) == givenName {
			return
		}
	}

	id = ResourceTypeUnknown
	err = errors.New(fmt.Sprintf("Unknown resource type name: %s", givenName))
	return
}

func (t *ResourceType) TypeById(id int) (tp *modelsConfig.ResourceType, err error) {
	if id == ResourceTypeUser {
		resourcesTypes.UserTypeObject.Meta.Location = t.baseUrl + "/" + system.AppVersion() + "/ResourceTypes/" + ResourceTypes[ResourceTypeUser]
		return &resourcesTypes.UserTypeObject, nil
	} else if id == ResourceTypeGroup {
		resourcesTypes.GroupTypeObject.Meta.Location = t.baseUrl + "/" + system.AppVersion() + "/ResourceTypes/" + ResourceTypes[ResourceTypeGroup]
		return &resourcesTypes.GroupTypeObject, nil
	} else {
		return nil, errors.New(fmt.Sprintf("Unknown resource type ID: %d", id))
	}
}

func (t *ResourceType) Types() (list *modelsConfig.ResourceTypes, err error) {
	resources := make([]*modelsConfig.ResourceType, 0)

	if data, err := t.TypeById(ResourceTypeUser); err != nil {
		return nil, err
	} else {
		resources = append(resources, data)
	}

	if data, err := t.TypeById(ResourceTypeGroup); err != nil {
		return nil, err
	} else {
		resources = append(resources, data)
	}

	return &modelsConfig.ResourceTypes{
		Schemas:      []string{"urn:ietf:params:scim:api:messages:2.0:ListResponse"},
		TotalResults: len(resources),
		ItemsPerPage: len(resources),
		StartIndex:   1,
		Resources:    resources,
	}, nil
}
