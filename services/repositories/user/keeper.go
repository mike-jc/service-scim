package repositoriesUser

import (
	"encoding/json"
	"fmt"
	"service-scim/errors/repositories"
	"service-scim/models/resources"
	"service-scim/models/sdks/keeper"
	"service-scim/sdks/restApi"
	"service-scim/system"
	"strconv"
	"strings"
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

func (k *Keeper) List(offset, limit int, filterMap map[string]interface{}) (totalCount int, list []*modelsResources.User, err errorsRepositories.Interface) {
	if keeperUsers, kErr := k.client.Users(offset, limit, filterMap); kErr != nil {
		err = errorsRepositories.NewError(kErr.Error(), errorsRepositories.ApiError)
		return
	} else {
		totalCount = keeperUsers.Total
		list = make([]*modelsResources.User, 0)
		for _, keeperUser := range keeperUsers.Data {
			list = append(list, k.resourceFromKeeperUser(keeperUser))
		}
		return
	}
}

func (k *Keeper) ById(id string) (user *modelsResources.User, err errorsRepositories.Interface) {
	if keeperUser, kErr := k.client.User(id); kErr != nil {
		if kErr.Response() != nil && kErr.Response().Code == modelsSdkKeeper.NotFoundCode {
			return nil, errorsRepositories.NewError(kErr.Error(), errorsRepositories.NotFoundError)
		} else {
			return nil, errorsRepositories.NewError(kErr.Error(), errorsRepositories.ApiError)
		}
	} else {
		return k.resourceFromKeeperUser(&keeperUser), nil
	}
}

func (k *Keeper) Create(data map[string]interface{}) (resultedUser *modelsResources.User, err errorsRepositories.Interface) {
	keeperMap, mErr := k.keeperUserMapFromMap(data)
	if mErr != nil {
		return nil, mErr
	}

	if resultedKeeperUser, kErr := k.client.CreateUser(keeperMap); kErr != nil {
		return nil, errorsRepositories.NewError(kErr.Error(), errorsRepositories.ApiError)
	} else {
		return k.resourceFromKeeperUser(&resultedKeeperUser), nil
	}
}

func (k *Keeper) Update(id string, data map[string]interface{}) (resultedUser *modelsResources.User, err errorsRepositories.Interface) {
	keeperMap, mErr := k.keeperUserMapFromMap(data)
	if mErr != nil {
		return nil, mErr
	}

	if resultedKeeperUser, kErr := k.client.UpdateUser(id, keeperMap); kErr != nil {
		return nil, errorsRepositories.NewError(kErr.Error(), errorsRepositories.ApiError)
	} else {
		return k.resourceFromKeeperUser(&resultedKeeperUser), nil
	}
}

func (k *Keeper) Replace(id string, data map[string]interface{}) (resultedUser *modelsResources.User, err errorsRepositories.Interface) {
	keeperMap, mErr := k.keeperUserMapFromMap(data)
	if mErr != nil {
		return nil, mErr
	}

	if resultedKeeperUser, kErr := k.client.ReplaceUser(id, keeperMap); kErr != nil {
		return nil, errorsRepositories.NewError(kErr.Error(), errorsRepositories.ApiError)
	} else {
		return k.resourceFromKeeperUser(&resultedKeeperUser), nil
	}
}

func (k *Keeper) Block(id string) errorsRepositories.Interface {
	if kErr := k.client.BlockUser(id); kErr != nil {
		if kErr.Response().Code == 404 {
			return errorsRepositories.NewError(kErr.Error(), errorsRepositories.NotFoundError)
		} else {
			return errorsRepositories.NewError(kErr.Error(), errorsRepositories.ApiError)
		}
	}
	return nil
}

func (k *Keeper) Count(filter map[string]interface{}, id *string) (count int, err errorsRepositories.Interface) {
	keeperFilter := k.keeperUserFilter(filter)
	if len(keeperFilter) == 0 {
		return 0, nil
	} else {
		if response, kErr := k.client.CountUsers(keeperFilter, id); kErr != nil {
			return 0, errorsRepositories.NewError(kErr.Error(), errorsRepositories.ApiError)
		} else {
			return response.Count, nil
		}
	}
}

func (k *Keeper) resourceFromKeeperUser(user *modelsSdkKeeper.User) *modelsResources.User {
	resultedUser := &modelsResources.User{
		Id:         user.Id,
		ExternalId: user.ScimId,
		Locale:     user.Locale,
		Timezone:   user.Timezone,
		Active:     user.Active,
		ProfileUrl: user.ProfileUrl,
		Title:      user.JobTitle,
		Created:    user.CreatedAt,
		Updated:    user.UpdatedAt,
		Addresses:  make([]*modelsResources.Address, 0),
		Photos:     make([]*modelsResources.Photo, 0),
		Groups:     make([]*modelsResources.UserGroup, 0),
		Roles:      make([]*modelsResources.Role, 0),
	}
	if len(user.ScimId) > 0 {
		resultedUser.UserName = user.ScimId
	} else {
		resultedUser.UserName = strconv.FormatInt(user.Id, 10)
	}
	if len(user.Email) > 0 {
		resultedUser.Emails = []*modelsResources.Email{
			{
				Value:   user.Email,
				Primary: true,
			},
		}
	}
	if len(user.Name) > 0 {
		// now parsing is skipped since name for 24sessions user is built basing on template
		// so it's quite tricky (if possible at all) to parse it back to name attributes
		// modelsSdkKeeper.ParseUserName(user.Name, resultedUser)
		resultedUser.Name = &modelsResources.UserName{
			Formatted: strings.TrimSpace(user.Name),
		}
	}
	if len(user.Phone) > 0 {
		resultedUser.PhoneNumbers = []*modelsResources.PhoneNumber{
			{
				Value:   user.Phone,
				Primary: true,
			},
		}
	}
	if user.Address != nil && user.Address.IsValid() {
		resultedUser.Addresses = append(resultedUser.Addresses, &modelsResources.Address{
			StreetAddress: user.Address.Street,
			Locality:      user.Address.City,
			Region:        user.Address.State,
			Country:       user.Address.Country,
			PostalCode:    user.Address.PostalCode,
		})
	}
	if user.Photos != nil {
		for _, ph := range user.Photos {
			if ph.IsValid() {
				resultedUser.Photos = append(resultedUser.Photos, &modelsResources.Photo{
					Value: ph.Url,
					Type:  ph.Type,
				})
			}
		}
	}
	if user.Groups != nil {
		for _, g := range user.Groups {
			if g.IsValid() {
				resultedUser.Groups = append(resultedUser.Groups, &modelsResources.UserGroup{
					Value:   strconv.FormatInt(g.Id, 10),
					Ref:     g.Ref(),
					Display: g.Name,
				})
			}
		}
	}
	if user.Roles != nil {
		for _, r := range user.Roles {
			resultedUser.Roles = append(resultedUser.Roles, &modelsResources.Role{
				Value:   r.Alias,
				Display: r.Name,
			})
		}
	}

	resultedUser.EnterpriseExtension = k.getEnterpriseUserResourceFromKeeper(user)

	return resultedUser
}

func (k *Keeper) keeperUserMapFromMap(data map[string]interface{}) (map[string]interface{}, errorsRepositories.Interface) {
	reflUser := system.ReflectValue(new(modelsResources.User))
	keeperMap := map[string]interface{}{
		"id":       k.mapValueForField(data, reflUser, "Id"),
		"scimId":   k.mapValueForField(data, reflUser, "ExternalId"),
		"timezone": k.mapValueForField(data, reflUser, "Timezone"),
		"jobTitle": k.mapValueForField(data, reflUser, "Title"),
		"active":   k.mapValueForField(data, reflUser, "Active"),
	}
	if name, ok := k.mapValueForField(data, reflUser, "Name").(map[string]interface{}); ok && name != nil {
		reflName := system.ReflectValue(new(modelsResources.UserName))
		normalizedName := map[string]interface{}{
			"formatted":       k.mapValueForField(name, reflName, "Formatted"),
			"familyName":      k.mapValueForField(name, reflName, "FamilyName"),
			"givenName":       k.mapValueForField(name, reflName, "GivenName"),
			"middleName":      k.mapValueForField(name, reflName, "MiddleName"),
			"honorificPrefix": k.mapValueForField(name, reflName, "HonorificPrefix"),
			"honorificSuffix": k.mapValueForField(name, reflName, "HonorificSuffix"),
		}
		if name, nErr := k.userNameFromTemplate(normalizedName); nErr != nil {
			nErr.SetError(fmt.Sprintf("Can not build user name: %s", nErr.Error()))
			return map[string]interface{}{}, nErr
		} else {
			keeperMap["name"] = name
		}
	}
	if emails := k.mapSliceOfMapsForField(data, reflUser, "Emails"); emails != nil {
		reflEmail := system.ReflectValue(new(modelsResources.Email))
		normalizedEmails := make([]map[string]interface{}, 0)
		for _, email := range emails {
			normalizedEmails = append(normalizedEmails, map[string]interface{}{
				"value":   k.mapValueForField(email, reflEmail, "Value"),
				"type":    k.mapValueForField(email, reflEmail, "Type"),
				"primary": k.mapValueForField(email, reflEmail, "Primary"),
			})
		}
		keeperMap["email"] = modelsSdkKeeper.PrimaryEmail(normalizedEmails)
	}
	if addresses := k.mapSliceOfMapsForField(data, reflUser, "Addresses"); addresses != nil {
		reflAddress := system.ReflectValue(new(modelsResources.Address))
		for _, addr := range addresses {
			address := map[string]interface{}{
				"street":     k.mapValueForField(addr, reflAddress, "StreetAddress"),
				"city":       k.mapValueForField(addr, reflAddress, "Locality"),
				"state":      k.mapValueForField(addr, reflAddress, "Region"),
				"country":    k.mapValueForField(addr, reflAddress, "Country"),
				"postalCode": k.mapValueForField(addr, reflAddress, "PostalCode"),
			}
			if modelsSdkKeeper.AddressIsValid(address) {
				keeperMap["address"] = address
				break
			}
		}
	}
	if phones := k.mapSliceOfMapsForField(data, reflUser, "PhoneNumbers"); phones != nil {
		reflPhone := system.ReflectValue(new(modelsResources.PhoneNumber))
		normalizedPhones := make([]map[string]interface{}, 0)
		for _, ph := range phones {
			normalizedPhones = append(normalizedPhones, map[string]interface{}{
				"value":   k.mapValueForField(ph, reflPhone, "Value"),
				"type":    k.mapValueForField(ph, reflPhone, "Type"),
				"primary": k.mapValueForField(ph, reflPhone, "Primary"),
			})
			keeperMap["phone"] = modelsSdkKeeper.PrimaryPhone(normalizedPhones)
		}
	}
	if groups := k.mapSliceOfMapsForField(data, reflUser, "Groups"); groups != nil {
		keeperGroups := make([]map[string]interface{}, 0)
		reflGroup := system.ReflectValue(new(modelsResources.UserGroup))
		for _, g := range groups {
			group := map[string]interface{}{
				"name": k.mapValueForField(g, reflGroup, "Display"),
			}
			// group id
			groupValue := k.mapValueForField(g, reflGroup, "Value")
			switch groupValue.(type) {
			case string:
				if groupId, cnvErr := strconv.ParseInt(groupValue.(string), 10, 64); cnvErr == nil {
					group["id"] = groupId
				}
			case json.Number:
				group["id"] = groupValue
			}
			// validation
			if modelsSdkKeeper.GroupIsValid(group) {
				keeperGroups = append(keeperGroups, group)
			}
		}
		keeperMap["groups"] = keeperGroups
	}
	if roles := k.mapSliceOfMapsForField(data, reflUser, "Roles"); roles != nil {
		keeperRoles := make([]map[string]interface{}, 0)
		reflRole := system.ReflectValue(new(modelsResources.Role))
		for _, r := range roles {
			role := map[string]interface{}{
				"alias": k.mapValueForField(r, reflRole, "Value"),
				"name":  k.mapValueForField(r, reflRole, "Display"),
			}
			if modelsSdkKeeper.RoleIsValid(role) {
				keeperRoles = append(keeperRoles, role)
			}
			keeperMap["roles"] = keeperRoles
		}
	}

	if enterpriseUser, ok := k.mapValueForField(data, reflUser, "EnterpriseExtension").(map[string]interface{}); ok && enterpriseUser != nil {
		reflName := system.ReflectValue(new(modelsResources.EnterpriseUser))
		enterpriseData := map[string]interface{}{
			"organization":   k.mapValueForField(enterpriseUser, reflName, "Organization"),
			"department":     k.mapValueForField(enterpriseUser, reflName, "Department"),
			"employeeNumber": k.mapValueForField(enterpriseUser, reflName, "EmployeeNumber"),
			"division":       k.mapValueForField(enterpriseUser, reflName, "Division"),
		}
		keeperMap["locationName"] = modelsSdkKeeper.LocationName(enterpriseData)
	}

	return keeperMap, nil
}

func (k *Keeper) keeperUserFilter(filter map[string]interface{}) map[string]interface{} {
	m := make(map[string]interface{})

	if externalId, ok := filter["externalId"]; ok {
		m["scimId"] = externalId
	} else if userName, ok := filter["userName"]; ok {
		m["scimId"] = userName
	}
	if name, ok := filter["name"]; ok {
		m["name"] = name
	} else if name, ok := filter["name.formatted"]; ok {
		m["name"] = name
	}
	if email, ok := filter["email"]; ok {
		m["email"] = email
	} else if email, ok := filter["emails.value"]; ok {
		m["email"] = email
	}
	if phone, ok := filter["phone"]; ok {
		m["phone"] = phone
	} else if phone, ok := filter["phoneNumbers.value"]; ok {
		m["phone"] = phone
	}
	if jobTitle, ok := filter["title"]; ok {
		m["jobTitle"] = jobTitle
	} else if jobTitle, ok := filter["jobTitle"]; ok {
		m["jobTitle"] = jobTitle
	}

	return m
}

func (k *Keeper) getEnterpriseUserResourceFromKeeper(user *modelsSdkKeeper.User) *modelsResources.EnterpriseUser {
	enterpriseUser := &modelsResources.EnterpriseUser{}
	if len(user.LocationName) > 0 {
		modelsSdkKeeper.ParseLocationName(user.LocationName, enterpriseUser)
	}

	return enterpriseUser
}

func (k *Keeper) userNameFromTemplate(nameAttributes map[string]interface{}) (string, errorsRepositories.Interface) {
	if name, err := k.client.UserNameFromTemplate(nameAttributes); err != nil {
		return "", errorsRepositories.NewError(err.Error(), errorsRepositories.ApiError)
	} else {
		return name, nil
	}
}
