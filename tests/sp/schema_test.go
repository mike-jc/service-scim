package sp_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"service-scim/services/sp"
	"service-scim/tests"
	"testing"
)

var s *serviceSp.Schema

func init() {
	tests.Init()

	s = new(serviceSp.Schema)
	s.SetBaseUrl("http://127.0.0.1:8101")
}

func TestUserSchemaOK(t *testing.T) {
	if expected, err := ioutil.ReadFile("fixtures/schemas/user.json"); err != nil {
		t.Errorf("Can not read fixture: %s", err.Error())
	} else {
		user, _ := s.SchemaById(serviceSp.SchemaUser)
		if userJson, err := json.Marshal(user); err != nil {
			t.Errorf("Can not serialize user schema to JSON: %s", err.Error())
		} else {
			assert.JSONEq(t, string(expected), string(userJson))
		}
	}
}

func TestUserEnterpriseSchemaOK(t *testing.T) {
	if expected, err := ioutil.ReadFile("fixtures/schemas/userEnterprise.json"); err != nil {
		t.Errorf("Can not read fixture: %s", err.Error())
	} else {
		user, _ := s.SchemaById(serviceSp.SchemaUserEnterprise)
		if userJson, err := json.Marshal(user); err != nil {
			t.Errorf("Can not serialize enterprise user schema to JSON: %s", err.Error())
		} else {
			assert.JSONEq(t, string(expected), string(userJson))
		}
	}
}

func TestGroupSchemaOK(t *testing.T) {
	if expected, err := ioutil.ReadFile("fixtures/schemas/group.json"); err != nil {
		t.Errorf("Can not read fixture: %s", err.Error())
	} else {
		group, _ := s.SchemaById(serviceSp.SchemaGroup)
		if groupJson, err := json.Marshal(group); err != nil {
			t.Errorf("Can not serialize group schema to JSON: %s", err.Error())
		} else {
			assert.JSONEq(t, string(expected), string(groupJson))
		}
	}
}

func TestSchemaListOK(t *testing.T) {
	if expected, err := ioutil.ReadFile("fixtures/schemas/list.json"); err != nil {
		t.Errorf("Can not read fixture: %s", err.Error())
	} else {
		list, _ := s.Schemas()
		if listJson, err := json.Marshal(list); err != nil {
			t.Errorf("Can not serialize list of schemas to JSON: %s", err.Error())
		} else {
			assert.JSONEq(t, string(expected), string(listJson))
		}
	}
}
