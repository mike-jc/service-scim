package sp_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"service-scim/services/sp"
	"service-scim/tests"
	"testing"
)

var rt *serviceSp.ResourceType

func init() {
	tests.Init()

	rt = new(serviceSp.ResourceType)
	rt.SetBaseUrl("http://127.0.0.1:8101")
}

func TestUserTypeOK(t *testing.T) {
	if expected, err := ioutil.ReadFile("fixtures/resourceTypes/user.json"); err != nil {
		t.Errorf("Can not read fixture: %s", err.Error())
	} else {
		user, _ := rt.TypeById(serviceSp.ResourceTypeUser)
		if userJson, err := json.Marshal(user); err != nil {
			t.Errorf("Can not serialize user resource type to JSON: %s", err.Error())
		} else {
			assert.JSONEq(t, string(expected), string(userJson))
		}
	}
}

func TestGroupTypeOK(t *testing.T) {
	if expected, err := ioutil.ReadFile("fixtures/resourceTypes/group.json"); err != nil {
		t.Errorf("Can not read fixture: %s", err.Error())
	} else {
		group, _ := rt.TypeById(serviceSp.ResourceTypeGroup)
		if groupJson, err := json.Marshal(group); err != nil {
			t.Errorf("Can not serialize group resource type to JSON: %s", err.Error())
		} else {
			assert.JSONEq(t, string(expected), string(groupJson))
		}
	}
}

func TestResourceTypeListOK(t *testing.T) {
	if expected, err := ioutil.ReadFile("fixtures/resourceTypes/list.json"); err != nil {
		t.Errorf("Can not read fixture: %s", err.Error())
	} else {
		list, _ := rt.Types()
		if listJson, err := json.Marshal(list); err != nil {
			t.Errorf("Can not serialize list of resource types to JSON: %s", err.Error())
		} else {
			assert.JSONEq(t, string(expected), string(listJson))
		}
	}
}
