package normalization_test

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"service-scim/models/normalization"
	"service-scim/services/navigation"
	"service-scim/services/normalization"
	"service-scim/system"
	"service-scim/tests"
	"service-scim/tests/fixtures/normalization"
	"testing"
)

var norm *normalization.ByAttribute

func init() {
	tests.Init()
	norm = new(normalization.ByAttribute)
}

func objectFromJson(jsonFile string) (interface{}, error) {
	var object interface{}
	if jsonStr, err := ioutil.ReadFile("fixtures/normalization/" + jsonFile); err != nil {
		return nil, fmt.Errorf("Can not read fixture: %s", err.Error())
	} else if err := json.Unmarshal(jsonStr, &object); err != nil {
		return nil, fmt.Errorf("Can not unmarshal JSON string: %s", err.Error())
	} else {
		return object, nil
	}
}

func TestNormalizeObjectFullOK(t *testing.T) {
	if object, err := objectFromJson("object.json"); err != nil {
		t.Errorf(err.Error())
	} else if normalized, err := norm.Normalize(object, fixturesNormalization.AttributesObject, []navigation.Path{}, []navigation.Path{}); err != nil {
		t.Errorf("Can not normalize full object: %s", err.Error())
	} else {
		assert.EqualValues(t, object, system.StructToMap(normalized.Interface(), "json"))
	}
}

func TestNormalizeObjectPartiallyOK(t *testing.T) {
	included := []string{"emails.value", "name.formatted", "photos.value", "addresses.country", "groups.user.name"}

	if object, err := objectFromJson("object.json"); err != nil {
		t.Error(err.Error())
	} else if partialObject, err := objectFromJson("objectPartial.json"); err != nil {
		t.Error(err.Error())
	} else if includedAttributes, err := modelsNormalization.MakeAttributes(included); err != nil {
		t.Errorf("Can not make included attributes: %s", err.Error())
	} else if normalized, err := norm.Normalize(object, fixturesNormalization.AttributesObject, includedAttributes, []navigation.Path{}); err != nil {
		t.Errorf("Can not normalize partial object: %s", err.Error())
	} else {
		assert.EqualValues(t, partialObject, system.StructToMap(normalized.Interface(), "json"))
	}
}

func TestNormalizeFailed(t *testing.T) {
	if object, err := objectFromJson("objectFailed.json"); err != nil {
		t.Errorf(err.Error())
	} else if _, err := norm.Normalize(object, fixturesNormalization.AttributesObject, []navigation.Path{}, []navigation.Path{}); err == nil {
		t.Errorf("Normalization of failed object is successful, should fail")
	} else if err.Error() != "Got value of type map for attribute emails which is multi-valued" {
		t.Errorf("Got wrong error when normalization of failed object failed: %s", err.Error())
	}
}
