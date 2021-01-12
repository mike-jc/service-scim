package resources_test

import (
	"github.com/stretchr/testify/assert"
	"service-scim/resources/schemas"
	"service-scim/services/repositories/user"
	"service-scim/services/validation"
	"service-scim/tests"
	"testing"
)

func init() {
	tests.Init()
}

func TestAttributeTypeValidationOK(t *testing.T) {
	if user, err := readValidatingMap("typeOk.json"); err != nil {
		t.Errorf(err.Error())
	} else {
		v := new(validation.AttributeType)
		if err := v.Validate(user, &resourcesSchemas.UserSchemaObject); err != nil {
			t.Errorf("Type validation should be successful but failed: %s", err.Error())
		}
	}
}

func TestAttributeTypeValidationFailed(t *testing.T) {
	if user, err := readValidatingMap("typeInvalid.json"); err != nil {
		t.Errorf(err.Error())
	} else {
		v := new(validation.AttributeType)
		if err := v.Validate(user, &resourcesSchemas.UserSchemaObject); err == nil {
			t.Errorf("Type validation should failed but is successful")
		} else if err.Error() != "Invalid type at 'urn:ietf:params:scim:schemas:core:2.0:User:addresses', expected 'complex array', got 'string'" {
			t.Errorf("Type validation error is not as expected: %s", err.Error())
		}
	}
}

func TestCorrectCaseValidationOK(t *testing.T) {
	if user, err := readValidatingMap("caseOk.json"); err != nil {
		t.Errorf(err.Error())
	} else {
		v := new(validation.CorrectCase)
		if err := v.Validate(user, &resourcesSchemas.UserSchemaObject); err != nil {
			t.Errorf("Case validation should be successful but failed: %s", err.Error())
		}
	}
}

func TestRequiredAttributeValidationOK(t *testing.T) {
	if user, err := readValidatingMap("requiredOk.json"); err != nil {
		t.Errorf(err.Error())
	} else {
		v := new(validation.RequiredAttribute)
		if err := v.Validate(user, &resourcesSchemas.UserSchemaObject); err != nil {
			t.Errorf("Requirement validation should be successful but failed: %s", err.Error())
		}
	}
}

func TestRequiredAttributeValidationFailed(t *testing.T) {
	if user, err := readValidatingMap("requiredInvalid.json"); err != nil {
		t.Errorf(err.Error())
	} else {
		v := new(validation.RequiredAttribute)
		if err := v.Validate(user, &resourcesSchemas.UserSchemaObject); err == nil {
			t.Errorf("Requirement validation should failed but is successful")
		} else if err.Error() != "Missing required property value at 'urn:ietf:params:scim:schemas:core:2.0:User:emails'" {
			t.Errorf("Requirement validation error is not as expected: %s", err.Error())
		}
	}
}

func TestMutabilityValidationOK(t *testing.T) {
	if origin, err := readValidatingMap("mutabilityOrigin.json"); err != nil {
		t.Errorf(err.Error())
	} else {
		if user, err := readValidatingMap("mutabilityOk.json"); err != nil {
			t.Errorf(err.Error())
		} else {
			v := new(validation.AttributeMutability)
			if err := v.Validate(user, origin, &resourcesSchemas.UserSchemaObject); err != nil {
				t.Errorf("Mutability validation should be successful but failed: %s", err.Error())
			}
		}
	}
}

func TestMutabilityValidationResetToOrigin(t *testing.T) {
	if origin, err := readValidatingMap("mutabilityOrigin.json"); err != nil {
		t.Errorf(err.Error())
	} else {
		if user, err := readValidatingMap("mutabilityInvalid.json"); err != nil {
			t.Errorf(err.Error())
		} else {
			v := new(validation.AttributeMutability)
			if err := v.Validate(user, origin, &resourcesSchemas.UserSchemaObject); err != nil {
				t.Errorf("Mutability validation should be successful but failed")
			} else {
				assert.Equal(t, origin["groups"], user["groups"])
			}
		}
	}
}

func TestUniquenessValidationOK(t *testing.T) {
	if user, err := readValidatingMap("uniquenessOk.json"); err != nil {
		t.Errorf(err.Error())
	} else {
		fake := new(repositoriesUser.Fake)
		v := new(validation.Uniqueness)
		v.SetRepository(fake)

		if err := v.Validate(user, nil, &resourcesSchemas.UserSchemaObject); err != nil {
			t.Errorf("Uniqueness validation should be successful but failed: %s", err.Error())
		}
	}
}

func TestUniquenessValidationFailed(t *testing.T) {
	if user, err := readValidatingMap("uniquenessInvalid.json"); err != nil {
		t.Errorf(err.Error())
	} else {
		fake := new(repositoriesUser.Fake)
		v := new(validation.Uniqueness)
		v.SetRepository(fake)

		id := "428e7353-204b-4e1a-8ac2-3b2bfa0616f5"
		if err := v.Validate(user, &id, &resourcesSchemas.UserSchemaObject); err == nil {
			t.Errorf("Uniqueness validation should failed but is successful")
		} else if err.Error() != "Resource has duplicate value '51db7ce8-75d7-4c81-ad05-c98a9233811e' at path 'externalId'" {
			t.Errorf("Uniqueness validation error is not as expected: %s", err.Error())
		}
	}
}

func TestReadOnlyAttributeValidationOK(t *testing.T) {
	if user, err := readValidatingMap("readOnlyOk.json"); err != nil {
		t.Errorf(err.Error())
	} else {
		v := new(validation.ReadOnlyAttribute)
		if err := v.ValidateUser(user, "json", nil); err != nil {
			t.Errorf("Read-only attributes validation should be successful but failed: %s", err.Error())
		}
	}
}

func readValidatingMap(fileName string) (map[string]interface{}, error) {
	var data map[string]interface{}
	if _, err := tests.ReadFromFixture(&data, "fixtures/validation/"+fileName); err != nil {
		return nil, err
	} else {
		return data, nil
	}
}
