package sp_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gitlab.com/24sessions/sdk-go-configurator/data"
	"io/ioutil"
	"service-scim/services"
	"service-scim/services/sp"
	"service-scim/tests"
	"testing"
)

var sp *serviceSp.ServiceConfig
var scimConfig *sdksData.ScimContainer

func init() {
	tests.Init()

	var err error
	if scimConfig, err = services.NewConfig(""); err != nil {
		panic("Can not get pseudo instance configuration")
	}

	sp = new(serviceSp.ServiceConfig)
	sp.SetScimConfig(scimConfig)
	sp.SetBaseUrl("http://127.0.0.1:8101")
}

func TestSpConfigOK(t *testing.T) {
	if expected, err := ioutil.ReadFile("fixtures/spConfig.json"); err != nil {
		t.Errorf("Can not read fixture: %s", err.Error())
	} else {
		if spConfig, err := sp.Config(); err != nil {
			t.Errorf("Can get SP configuration: %s", err.Error())
		} else if spConfigStr, err := json.Marshal(spConfig); err != nil {
			t.Errorf("Can not serialize SP configuration to JSON: %s", err.Error())
		} else {
			assert.JSONEq(t, string(expected), string(spConfigStr))
		}
	}
}
