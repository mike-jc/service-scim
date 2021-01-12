package scimErrors

import (
	"fmt"
	"reflect"
	"service-scim/models/config"
)

func InvalidAttrTypeError(attr *modelsConfig.Attribute, actualType string) error {
	return fmt.Errorf("Invalid type at '%s', expected '%s', got '%s'", attr.Navigation.FullPath, attr.TypeExpectation(), actualType)
}
func InvalidValueTypeError(val reflect.Value, attr *modelsConfig.Attribute) error {
	return fmt.Errorf("Got value of type %s for attribute %s which is %s", val.Kind(), attr.Name, attr.Type)
}

func MultiValueError(val reflect.Value, attr *modelsConfig.Attribute) error {
	return fmt.Errorf("Got value of type %s for attribute %s which is multi-valued", val.Kind(), attr.Name)
}

func NoAttributeError(path string) error {
	return fmt.Errorf("No attribute defined for path (segment) '%s'", path)
}

func MissingRequiredPropertyError(path string) error {
	return fmt.Errorf("Missing required property value at '%s'", path)
}

func MutabilityViolationError(path string) error {
	return fmt.Errorf("Violated mutability rule at '%s'", path)
}

func DuplicationError(val interface{}, path string) error {
	return fmt.Errorf("Resource has duplicate value '%v' at path '%s'", val, path)
}

func InvalidParameterError(name, expect string, got interface{}) error {
	return fmt.Errorf("Invalid parameter for '%s', expect %s, but got %+v", name, expect, got)
}

func InvalidPathError(path, details string) error {
	return fmt.Errorf("Path [%s] is invalid: %s", path, details)
}

func InvalidFilterError(filter, details string) error {
	if len(filter) > 0 {
		return fmt.Errorf("Filter [%s] is invalid: %s", filter, details)
	} else {
		return fmt.Errorf("Filter is invalid: %s", details)
	}
}
