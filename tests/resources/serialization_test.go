package resources_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"service-scim/services/repositories"
	"service-scim/services/resources"
	"service-scim/tests"
	"testing"
)

func init() {
	tests.Init()
}

func TestUserOK(t *testing.T) {
	if expected, err := ioutil.ReadFile("fixtures/resources/users/jack.json"); err != nil {
		t.Errorf("Can not read fixture: %s", err.Error())
	} else {
		id := "123"
		ctx := tests.NewContext()
		ctx.Input.SetParam("attributes", "emails.value, name.formatted, active, addresses.country, photos.value")

		userRepository := new(repositories.Factory).NewUserEngineOrPanic()
		userService := resources.NewUserService(userRepository, ctx, tests.NewConfig(), "local.24sessions.com", "json")

		if user, err := userService.ById(id); err != nil {
			t.Errorf("Can not get user by ID: %s", err.Error())
		} else if userJson, err := json.Marshal(user); err != nil {
			t.Errorf("Can not serialize user to JSON: %s", err.Error())
		} else {
			assert.JSONEq(t, string(expected), string(userJson))
		}
	}
}

func TestUserListOK(t *testing.T) {
	if expected, err := ioutil.ReadFile("fixtures/resources/users/list.json"); err != nil {
		t.Errorf("Can not read fixture: %s", err.Error())
	} else {
		ctx := tests.NewContext()
		ctx.Input.SetParam("attributes", "emails.value, name.formatted, active, addresses.country, photos.value")

		userRepository := new(repositories.Factory).NewUserEngineOrPanic()
		userService := resources.NewUserService(userRepository, ctx, tests.NewConfig(), "local.24sessions.com", "json")

		if list, err := userService.List(0, 2, ""); err != nil {
			t.Errorf("Cannot get user list: %s", err.Error())
		} else if listJson, err := json.Marshal(list); err != nil {
			t.Errorf("Can not serialize user list to JSON: %s", err.Error())
		} else {
			assert.JSONEq(t, string(expected), string(listJson))
		}
	}
}

func TestGroupOK(t *testing.T) {
	if expected, err := ioutil.ReadFile("fixtures/resources/groups/operators.json"); err != nil {
		t.Errorf("Can not read fixture: %s", err.Error())
	} else {
		id := "Operators"
		ctx := tests.NewContext()
		ctx.Input.SetParam("attributes", "displayName, members.value, members.$ref")

		groupRepository := new(repositories.Factory).NewGroupEngineOrPanic()
		groupService := resources.NewGroupService(groupRepository, ctx, tests.NewConfig(), "local.24sessions.com", "json")

		if group, err := groupService.ById(id); err != nil {
			t.Errorf("Can not get group by ID: %s", err.Error())
		} else if groupJson, err := json.Marshal(group); err != nil {
			t.Errorf("Can not serialize group to JSON: %s", err.Error())
		} else {
			assert.JSONEq(t, string(expected), string(groupJson))
		}
	}
}

func TestGroupListOK(t *testing.T) {
	if expected, err := ioutil.ReadFile("fixtures/resources/groups/list.json"); err != nil {
		t.Errorf("Can not read fixture: %s", err.Error())
	} else {
		ctx := tests.NewContext()
		ctx.Input.SetParam("attributes", "displayName, members.value, members.$ref")

		groupRepository := new(repositories.Factory).NewGroupEngineOrPanic()
		groupService := resources.NewGroupService(groupRepository, ctx, tests.NewConfig(), "local.24sessions.com", "json")

		if list, err := groupService.List(0, 2); err != nil {
			t.Errorf("Cannot get group list: %s", err.Error())
		} else if listJson, err := json.Marshal(list); err != nil {
			t.Errorf("Can not serialize group list to JSON: %s", err.Error())
		} else {
			assert.JSONEq(t, string(expected), string(listJson))
		}
	}
}
