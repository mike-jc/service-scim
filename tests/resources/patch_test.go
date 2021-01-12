package resources

import (
	"github.com/stretchr/testify/assert"
	"service-scim/models/resources"
	"service-scim/resources/schemas"
	"service-scim/services/resources"
	"service-scim/tests"
	"testing"
)

func init() {
	tests.Init()
}

func TestPatchOK(t *testing.T) {
	if originalUser, err := readUserForPatching("userOriginal.json"); err != nil {
		t.Errorf(err.Error())
	} else if patchedUser, err := readUserForPatching("userPatched.json"); err != nil {
		t.Errorf(err.Error())
	} else if modification, err := readPatch("patchOk.json"); err != nil {
		t.Errorf(err.Error())
	} else {
		userMap := originalUser.ToMap("json")
		for _, op := range modification.Operations {
			if pErr := resources.ApplyPatch(op, userMap, &resourcesSchemas.UserSchemaObject); pErr != nil {
				t.Errorf("Cannot apply patch to user: %s", pErr.Error())
			}
		}
		assert.Equal(t, patchedUser.ToMap("json"), userMap)
	}
}

func TestPatchEnterpriseExtension(t *testing.T) {
	if originalUser, err := readUserForPatching("userOriginalEnterpriseExtension.json"); err != nil {
		t.Errorf(err.Error())
	} else if patchedUser, err := readUserForPatching("userPatchedEnterpriseExtension.json"); err != nil {
		t.Errorf(err.Error())
	} else if modification, err := readPatch("patchOkEnterpriseExtension.json"); err != nil {
		t.Errorf(err.Error())
	} else {
		userMap := originalUser.ToMap("json")
		for _, op := range modification.Operations {
			if pErr := resources.ApplyPatch(op, userMap, &resourcesSchemas.UserSchemaObject); pErr != nil {
				t.Errorf("Cannot apply patch to user: %s", pErr.Error())
			}
		}

		m := patchedUser.ToMap("json")
		assert.Equal(t, m, userMap)
	}
}

func TestPatchFailed(t *testing.T) {
	if modification, err := readPatch("patchInvalid.json"); err != nil {
		t.Errorf(err.Error())
	} else if vErr := modification.Validate(); vErr == nil {
		t.Errorf("Patch validation should faile but is successful")
	} else if vErr.Error() != "Invalid parameter for 'value of replace op', expect to be present, but got nil" {
		t.Errorf("Patch validation error is not as expected: %s", vErr.Error())
	}
}

func readUserForPatching(fileName string) (*modelsResources.User, error) {
	var user modelsResources.User
	if _, err := tests.ReadFromFixture(&user, "fixtures/patching/"+fileName); err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

func readPatch(fileName string) (*modelsResources.Modification, error) {
	var modification modelsResources.Modification
	if _, err := tests.ReadFromFixture(&modification, "fixtures/patching/"+fileName); err != nil {
		return nil, err
	} else {
		return &modification, nil
	}
}
