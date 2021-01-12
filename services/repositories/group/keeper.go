package repositoriesGroup

import (
	"encoding/json"
	"service-scim/errors/repositories"
	"service-scim/models/resources"
	"service-scim/models/sdks/keeper"
	"service-scim/sdks/restApi"
	"service-scim/system"
	"strconv"
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

func (k *Keeper) SetInstanceDomain(domain string) {
	k.client.SetInstanceDomain(domain)
}

func (k *Keeper) List(offset, limit int) (totalCount int, list []*modelsResources.Group, err errorsRepositories.Interface) {
	if keeperGroups, kErr := k.client.Groups(offset, limit); kErr != nil {
		err = errorsRepositories.NewError(kErr.Error(), errorsRepositories.ApiError)
		return
	} else {
		totalCount = keeperGroups.Total
		list = make([]*modelsResources.Group, 0)
		for _, keeperGroup := range keeperGroups.Data {
			list = append(list, k.resourceFromKeeperGroup(keeperGroup))
		}
		return
	}
}

func (k *Keeper) ById(id string) (group *modelsResources.Group, err errorsRepositories.Interface) {
	if keeperGroup, kErr := k.client.Group(id); kErr != nil {
		if kErr.Response() != nil && kErr.Response().Code == modelsSdkKeeper.NotFoundCode {
			return nil, errorsRepositories.NewError(kErr.Error(), errorsRepositories.NotFoundError)
		} else {
			return nil, errorsRepositories.NewError(kErr.Error(), errorsRepositories.ApiError)
		}
	} else {
		return k.resourceFromKeeperGroup(&keeperGroup), nil
	}
}

func (k *Keeper) Create(data map[string]interface{}) (resultedGroup *modelsResources.Group, err errorsRepositories.Interface) {
	if resultedKeeperGroup, kErr := k.client.CreateGroup(k.keeperGroupMapFromMap(data)); kErr != nil {
		return nil, errorsRepositories.NewError(kErr.Error(), errorsRepositories.ApiError)
	} else {
		return k.resourceFromKeeperGroup(&resultedKeeperGroup), nil
	}
}

func (k *Keeper) Update(id string, data map[string]interface{}) (resultedGroup *modelsResources.Group, err errorsRepositories.Interface) {
	if resultedKeeperGroup, kErr := k.client.UpdateGroup(id, k.keeperGroupMapFromMap(data)); kErr != nil {
		return nil, errorsRepositories.NewError(kErr.Error(), errorsRepositories.ApiError)
	} else {
		return k.resourceFromKeeperGroup(&resultedKeeperGroup), nil
	}
}

func (k *Keeper) Replace(id string, data map[string]interface{}) (resultedGroup *modelsResources.Group, err errorsRepositories.Interface) {
	if resultedKeeperGroup, kErr := k.client.ReplaceGroup(id, k.keeperGroupMapFromMap(data)); kErr != nil {
		return nil, errorsRepositories.NewError(kErr.Error(), errorsRepositories.ApiError)
	} else {
		return k.resourceFromKeeperGroup(&resultedKeeperGroup), nil
	}
}

func (k *Keeper) Disable(id string) errorsRepositories.Interface {
	if kErr := k.client.DisableGroup(id); kErr != nil {
		if kErr.Response().Code == 404 {
			return errorsRepositories.NewError(kErr.Error(), errorsRepositories.NotFoundError)
		} else {
			return errorsRepositories.NewError(kErr.Error(), errorsRepositories.ApiError)
		}
	}
	return nil
}

func (k *Keeper) Count(filter map[string]interface{}, id *string) (count int, err errorsRepositories.Interface) {
	keeperFilter := k.keeperGroupFilter(filter)
	if len(keeperFilter) == 0 {
		return 0, nil
	} else {
		if response, kErr := k.client.CountGroups(keeperFilter, id); kErr != nil {
			return 0, errorsRepositories.NewError(kErr.Error(), errorsRepositories.ApiError)
		} else {
			return response.Count, nil
		}
	}
}

func (k *Keeper) resourceFromKeeperGroup(group *modelsSdkKeeper.Group) *modelsResources.Group {
	resultedGroup := &modelsResources.Group{
		Id:          group.Id,
		DisplayName: group.Name,
		ExternalId:  group.ScimId,
		Updated:     group.UpdatedAt,
		Members:     make([]*modelsResources.GroupMember, 0),
	}
	if group.Users != nil {
		for _, u := range group.Users {
			if u.IsValid() {
				resultedGroup.Members = append(resultedGroup.Members, &modelsResources.GroupMember{
					Value:   strconv.FormatInt(u.Id, 10),
					Ref:     u.Ref(),
					Type:    "User",
					Display: u.Name,
				})
			}
		}
	}
	return resultedGroup
}

func (k *Keeper) keeperGroupMapFromMap(data map[string]interface{}) map[string]interface{} {
	reflGroup := system.ReflectValue(new(modelsResources.Group))
	keeperMap := map[string]interface{}{
		"id":     k.mapValueForField(data, reflGroup, "Id"),
		"scimId": k.mapValueForField(data, reflGroup, "ExternalId"),
		"name":   k.mapValueForField(data, reflGroup, "DisplayName"),
	}
	if members := k.mapSliceOfMapsForField(data, reflGroup, "Members"); members != nil {
		keeperUsers := make([]map[string]interface{}, 0)
		reflMember := system.ReflectValue(new(modelsResources.GroupMember))
		for _, m := range members {
			user := make(map[string]interface{})
			// member id
			memberValue := k.mapValueForField(m, reflMember, "Value")
			switch memberValue.(type) {
			case string:
				if memberId, cnvErr := strconv.ParseInt(memberValue.(string), 10, 64); cnvErr == nil {
					user["id"] = memberId
				}
			case json.Number:
				user["id"] = memberValue
			}
			// member roles
			if memberRoles := k.mapSliceOfMapsForField(m, reflMember, "Roles"); memberRoles != nil {
				keeperRoles := make([]map[string]interface{}, 0)
				reflRole := system.ReflectValue(new(modelsResources.Role))
				for _, r := range memberRoles {
					role := map[string]interface{}{
						"alias": k.mapValueForField(r, reflRole, "Value"),
					}
					if modelsSdkKeeper.RoleIsValid(role) {
						keeperRoles = append(keeperRoles, role)
					}
				}
				user["roles"] = keeperRoles
			}
			// validation
			if modelsSdkKeeper.UserIsValid(user) {
				keeperUsers = append(keeperUsers, user)
			}
		}
		keeperMap["users"] = keeperUsers
	}
	return keeperMap
}

func (k *Keeper) keeperGroupFilter(filter map[string]interface{}) map[string]interface{} {
	m := make(map[string]interface{})

	if externalId, ok := filter["externalId"]; ok {
		m["scimId"] = externalId
	}
	if displayName, ok := filter["displayName"]; ok {
		m["name"] = displayName
	}

	return m
}
