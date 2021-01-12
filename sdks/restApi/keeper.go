package restApi

import (
	"fmt"
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"service-scim/errors/sdks"
	"service-scim/models/sdks/jsonApi"
	"service-scim/models/sdks/jwt"
	"service-scim/models/sdks/keeper"
	"service-scim/sdks"
	"strings"
	"time"
)

const TokenExpirationPeriod = 30 * time.Minute

type Keeper struct {
	baseUrl        string
	instanceDomain string

	token struct {
		value    string
		expireAt time.Time
	}

	client *RestJson
	jwt    *sdks.Jwt
}

func (k *Keeper) Client() *RestJson {
	if k.client == nil {
		k.client = new(RestJson)
	}
	return k.client
}

func (k *Keeper) Jwt() *sdks.Jwt {
	if k.jwt == nil {
		k.jwt = new(sdks.Jwt)
	}
	return k.jwt
}

func (k *Keeper) SetBaseUrl(url string) {
	k.token.value = "" // force token re-request since API base URL is changed
	k.baseUrl = url
}

func (k *Keeper) SetInstanceDomain(domain string) {
	k.token.value = "" // force token re-request since API instance is changed
	k.instanceDomain = domain
	k.Client().SetInstanceDomain(domain)
}

func (k *Keeper) absoluteUrl(url string) string {
	return fmt.Sprintf("%s/%s", strings.TrimRight(k.baseUrl, "/"), strings.TrimLeft(url, "/"))
}

func (k *Keeper) IsExpired(t time.Time) bool {
	if t.Unix() <= 0 {
		return false
	}
	return time.Now().Add(5 * time.Minute).After(t)
}

// Generate new token if there's no token or it's expired
func (k *Keeper) TokenValue() (value string, err error) {
	if k.token.value == "" || k.IsExpired(k.token.expireAt) {
		payload := k.Jwt().CreateSystemUserPayload(k.instanceDomain)
		ttl := TokenExpirationPeriod + 1*time.Minute

		if value, vErr := k.Jwt().GenerateToken(payload, ttl, modelsSdkJwt.ServiceKeeper); vErr != nil {
			return "", vErr
		} else {
			k.token.value = value
			k.token.expireAt = time.Now().Add(TokenExpirationPeriod)
		}
	}
	return k.token.value, nil
}

func (k *Keeper) Users(offset, limit int, filter map[string]interface{}) (users modelsSdkKeeper.Users, err errorsSdks.KeeperInterface) {
	if token, tErr := k.TokenValue(); tErr != nil {
		err = errorsSdks.NewKeeperError("Can not get users from Keeper: "+tErr.Error(), nil)
		return
	} else {
		var url = k.absoluteUrl("instance/scim/users")
		var params = map[string]interface{}{
			"offset": offset,
			"limit":  limit,
			"filter": filter,
		}
		var headers = map[string]string{
			"Authorization": "Bearer " + token,
		}
		if cErr := k.Client().Get(url, params, headers, &users); cErr != nil {
			cErr.SetError("Can not get users from Keeper: " + cErr.Error())
			err = errorsSdks.NewKeeperErrorFromAbstract(cErr)
			return
		}
	}
	return
}

func (k *Keeper) User(id string) (user modelsSdkKeeper.User, err errorsSdks.KeeperInterface) {
	if token, tErr := k.TokenValue(); tErr != nil {
		err = errorsSdks.NewKeeperError("Can not get user by ID from Keeper: "+tErr.Error(), nil)
		return
	} else {
		var url = k.absoluteUrl("instance/scim/users/" + id)
		var headers = map[string]string{
			"Authorization": "Bearer " + token,
		}
		if cErr := k.Client().Get(url, nil, headers, &user); cErr != nil {
			cErr.SetError("Can not get user by ID from Keeper: " + cErr.Error())
			err = errorsSdks.NewKeeperErrorFromAbstract(cErr)
			return
		}
	}
	return
}

func (k *Keeper) CreateUser(data map[string]interface{}) (resultedUser modelsSdkKeeper.User, err errorsSdks.KeeperInterface) {
	if token, tErr := k.TokenValue(); tErr != nil {
		err = errorsSdks.NewKeeperError("Can not create user in Keeper: "+tErr.Error(), nil)
		return
	} else {
		var url = k.absoluteUrl("instance/scim/users")
		var headers = map[string]string{
			"Authorization": "Bearer " + token,
		}
		if cErr := k.Client().Post(url, data, headers, &resultedUser); cErr != nil {
			cErr.SetError("Can not create user in Keeper: " + cErr.Error())
			err = errorsSdks.NewKeeperErrorFromAbstract(cErr)
			return
		}
	}
	return
}

func (k *Keeper) UpdateUser(id string, data map[string]interface{}) (resultedUser modelsSdkKeeper.User, err errorsSdks.KeeperInterface) {
	if token, tErr := k.TokenValue(); tErr != nil {
		err = errorsSdks.NewKeeperError("Can not update user in Keeper: "+tErr.Error(), nil)
		return
	} else {
		var url = k.absoluteUrl("instance/scim/users/" + id)
		var headers = map[string]string{
			"Authorization": "Bearer " + token,
		}
		if cErr := k.Client().Patch(url, data, headers, &resultedUser); cErr != nil {
			cErr.SetError("Can not update user in Keeper: " + cErr.Error())
			err = errorsSdks.NewKeeperErrorFromAbstract(cErr)
			return
		}
	}
	return
}

func (k *Keeper) ReplaceUser(id string, data map[string]interface{}) (resultedUser modelsSdkKeeper.User, err errorsSdks.KeeperInterface) {
	if token, tErr := k.TokenValue(); tErr != nil {
		err = errorsSdks.NewKeeperError("Can not replace user in Keeper: "+tErr.Error(), nil)
		return
	} else {
		var url = k.absoluteUrl("instance/scim/users/" + id)
		var headers = map[string]string{
			"Authorization": "Bearer " + token,
		}
		if cErr := k.Client().Put(url, data, headers, &resultedUser); cErr != nil {
			cErr.SetError("Can not replace user in Keeper: " + cErr.Error())
			err = errorsSdks.NewKeeperErrorFromAbstract(cErr)
			return
		}
	}
	return
}

func (k *Keeper) BlockUser(id string) errorsSdks.KeeperInterface {
	if token, tErr := k.TokenValue(); tErr != nil {
		return errorsSdks.NewKeeperError("Can not block user in Keeper: "+tErr.Error(), nil)
	} else {
		var url = k.absoluteUrl("instance/scim/users/" + id)
		var headers = map[string]string{
			"Authorization": "Bearer " + token,
		}
		if cErr := k.Client().Delete(url, headers); cErr != nil {
			cErr.SetError("Can not disable user in Keeper: " + cErr.Error())
			return errorsSdks.NewKeeperErrorFromAbstract(cErr)
		}
	}
	return nil
}

func (k *Keeper) CountUsers(filter map[string]interface{}, id *string) (*modelsSdkKeeper.CountResponse, errorsSdks.KeeperInterface) {
	if token, tErr := k.TokenValue(); tErr != nil {
		return nil, errorsSdks.NewKeeperError("Can not count users in Keeper: "+tErr.Error(), nil)
	} else {
		var url string
		if id != nil {
			url = k.absoluteUrl("instance/scim/users/count/" + (*id))
		} else {
			url = k.absoluteUrl("instance/scim/users/count")
		}
		var headers = map[string]string{
			"Authorization": "Bearer " + token,
		}
		var response modelsSdkKeeper.CountResponse
		if cErr := k.Client().Get(url, filter, headers, &response); cErr != nil {
			cErr.SetError("Can not count users in Keeper: " + cErr.Error())
			return nil, errorsSdks.NewKeeperErrorFromAbstract(cErr)
		} else {
			return &response, nil
		}
	}
}

func (k *Keeper) Groups(offset, limit int) (groups modelsSdkKeeper.Groups, err errorsSdks.KeeperInterface) {
	if token, tErr := k.TokenValue(); tErr != nil {
		err = errorsSdks.NewKeeperError("Can not get groups from Keeper: "+tErr.Error(), nil)
		return
	} else {
		var url = k.absoluteUrl("instance/scim/groups")
		var params = map[string]interface{}{
			"offset": offset,
			"limit":  limit,
		}
		var headers = map[string]string{
			"Authorization": "Bearer " + token,
		}
		if cErr := k.Client().Get(url, params, headers, &groups); cErr != nil {
			cErr.SetError("Can not get groups from Keeper: " + cErr.Error())
			err = errorsSdks.NewKeeperErrorFromAbstract(cErr)
			return
		}
	}
	return
}

func (k *Keeper) Group(id string) (group modelsSdkKeeper.Group, err errorsSdks.KeeperInterface) {
	if token, tErr := k.TokenValue(); tErr != nil {
		err = errorsSdks.NewKeeperError("Can not get group by ID from Keeper: "+tErr.Error(), nil)
		return
	} else {
		var url = k.absoluteUrl("instance/scim/groups/" + id)
		var headers = map[string]string{
			"Authorization": "Bearer " + token,
		}
		if cErr := k.Client().Get(url, nil, headers, &group); cErr != nil {
			cErr.SetError("Can not get group by ID from Keeper: " + cErr.Error())
			err = errorsSdks.NewKeeperErrorFromAbstract(cErr)
			return
		}
	}
	return
}

func (k *Keeper) CreateGroup(data map[string]interface{}) (resultedGroup modelsSdkKeeper.Group, err errorsSdks.KeeperInterface) {
	if token, tErr := k.TokenValue(); tErr != nil {
		err = errorsSdks.NewKeeperError("Can not create group in Keeper: "+tErr.Error(), nil)
		return
	} else {
		var url = k.absoluteUrl("instance/scim/groups")
		var headers = map[string]string{
			"Authorization": "Bearer " + token,
		}
		if cErr := k.Client().Post(url, data, headers, &resultedGroup); cErr != nil {
			cErr.SetError("Can not create group in Keeper: " + cErr.Error())
			err = errorsSdks.NewKeeperErrorFromAbstract(cErr)
			return
		}
	}
	return
}

func (k *Keeper) UpdateGroup(id string, data map[string]interface{}) (resultedGroup modelsSdkKeeper.Group, err errorsSdks.KeeperInterface) {
	if token, tErr := k.TokenValue(); tErr != nil {
		err = errorsSdks.NewKeeperError("Can not update group in Keeper: "+tErr.Error(), nil)
		return
	} else {
		var url = k.absoluteUrl("instance/scim/groups/" + id)
		var headers = map[string]string{
			"Authorization": "Bearer " + token,
		}
		if cErr := k.Client().Patch(url, data, headers, &resultedGroup); cErr != nil {
			cErr.SetError("Can not update group in Keeper: " + cErr.Error())
			err = errorsSdks.NewKeeperErrorFromAbstract(cErr)
			return
		}
	}
	return
}

func (k *Keeper) ReplaceGroup(id string, data map[string]interface{}) (resultedGroup modelsSdkKeeper.Group, err errorsSdks.KeeperInterface) {
	if token, tErr := k.TokenValue(); tErr != nil {
		err = errorsSdks.NewKeeperError("Can not replace group in Keeper: "+tErr.Error(), nil)
		return
	} else {
		var url = k.absoluteUrl("instance/scim/groups/" + id)
		var headers = map[string]string{
			"Authorization": "Bearer " + token,
		}
		if cErr := k.Client().Put(url, data, headers, &resultedGroup); cErr != nil {
			cErr.SetError("Can not replace group in Keeper: " + cErr.Error())
			err = errorsSdks.NewKeeperErrorFromAbstract(cErr)
			return
		}
	}
	return
}

func (k *Keeper) DisableGroup(id string) errorsSdks.KeeperInterface {
	if token, tErr := k.TokenValue(); tErr != nil {
		return errorsSdks.NewKeeperError("Can not disable group in Keeper: "+tErr.Error(), nil)
	} else {
		var url = k.absoluteUrl("instance/scim/groups/" + id)
		var headers = map[string]string{
			"Authorization": "Bearer " + token,
		}
		if cErr := k.Client().Delete(url, headers); cErr != nil {
			cErr.SetError("Can not disable group in Keeper: " + cErr.Error())
			return errorsSdks.NewKeeperErrorFromAbstract(cErr)
		}
	}
	return nil
}

func (k *Keeper) CountGroups(filter map[string]interface{}, id *string) (*modelsSdkKeeper.CountResponse, errorsSdks.KeeperInterface) {
	if token, tErr := k.TokenValue(); tErr != nil {
		return nil, errorsSdks.NewKeeperError("Can not count groups in Keeper: "+tErr.Error(), nil)
	} else {
		var url string
		if id != nil {
			url = k.absoluteUrl("instance/scim/groups/count/" + (*id))
		} else {
			url = k.absoluteUrl("instance/scim/groups/count")
		}
		var headers = map[string]string{
			"Authorization": "Bearer " + token,
		}
		var response modelsSdkKeeper.CountResponse
		if cErr := k.Client().Get(url, filter, headers, &response); cErr != nil {
			cErr.SetError("Can not count groups in Keeper: " + cErr.Error())
			return nil, errorsSdks.NewKeeperErrorFromAbstract(cErr)
		} else {
			return &response, nil
		}
	}
}

func (k *Keeper) UserNameFromTemplate(nameAttributes map[string]interface{}) (name string, err errorsSdks.KeeperInterface) {
	if token, tErr := k.TokenValue(); tErr != nil {
		err = errorsSdks.NewKeeperError("Can not get user name from template in Keeper: "+tErr.Error(), nil)
		return
	} else {
		var url = k.absoluteUrl("instance/storage/scim/template")
		var headers = map[string]string{
			"Authorization": "Bearer " + token,
		}
		attributes := modelsSdkKeeper.TemplateCompilationRequest{
			TemplateName: "scim.user.name",
			Compile:      true,
			User: map[string]interface{}{
				"name": nameAttributes,
			},
		}
		request := k.Client().builtJsonApiRequest(attributes, "scim-resource-template-params")
		response := new(modelsSdkJsonApi.Response)
		if cErr := k.Client().Post(url, request, headers, response); cErr != nil {
			cErr.SetError("Can not get user name from template in Keeper: " + cErr.Error())
			err = errorsSdks.NewKeeperErrorFromAbstract(cErr)
			return
		}
		result := new(modelsSdkKeeper.TemplateCompilationResponse)
		if pErr := k.Client().parseJsonApiResponse(response, result, "scim-resource-settings"); pErr != nil {
			err = errorsSdks.NewKeeperError("Can not get user name from template in Keeper: "+pErr.Error(), nil)
			sdks.LogMain.Log(logger.CreateError(err.Error()).
				AddData("request", request).
				AddData("response", response))
			return
		} else {
			name = result.Value
			return
		}
	}
}

func (k *Keeper) Ping() errorsSdks.KeeperInterface {
	var url = k.absoluteUrl("/") // just index page
	if cErr := k.Client().Get(url, nil, nil, nil); cErr != nil {
		cErr.SetError("Error while ping Keeper: " + cErr.Error())
		return errorsSdks.NewKeeperErrorFromAbstract(cErr)
	}
	return nil
}
