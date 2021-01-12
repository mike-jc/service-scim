package middleware

import (
	"reflect"
	"service-scim/models/resources"
	"service-scim/system"
	"strings"
)

var userMapping = map[string][]string{
	"Name": {
		"formatted", "Formatted", "formattedName", "FormattedName",
		"displayName", "DisplayName", "display", "Display",
		"userName", "UserName", "fullName", "FullName",
	},
	"ExternalId": {
		"userName", "UserName",
	},
}
var userConcatMapping = map[string][][]string{
	"Name": {
		{"givenName", "GivenName", "firstName", "FirstName", "name", "Name"},
		{"middleName", "MiddleName"},
		{"familyName", "FamilyName", "lastName", "LastName", "surname", "Surname"},
	},
}
var userMultiValued = []string{"Emails"}

var groupMapping = map[string][]string{
	"DisplayName": {"name", "Name"},
}

type Default struct {
	Abstract
}

func (m *Default) ProcessUser(data map[string]interface{}, id *string) (map[string]interface{}, error) {
	reflectedUser := system.ReflectValue(new(modelsResources.User))

	if err := m.applyMappingToMap(reflectedUser, data, userMapping); err != nil {
		return nil, err
	}
	if err := m.applyConcatenationToMap(reflectedUser, data, userConcatMapping); err != nil {
		return nil, err
	}
	if err := m.applyMultiValued(reflectedUser, data, userMultiValued); err != nil {
		return nil, err
	}
	if err := m.applyEntitlements(reflectedUser, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *Default) ProcessGroup(data map[string]interface{}, id *string) (map[string]interface{}, error) {
	reflectedGroup := system.ReflectValue(new(modelsResources.Group))

	if err := m.applyMappingToMap(reflectedGroup, data, groupMapping); err != nil {
		return nil, err
	}

	return data, nil
}

// Applies attribute mapping to `data` basing on `val` structure.
// The mapping itself is a map [field name => array of alternative names]
func (m *Default) applyMappingToMap(val reflect.Value, data map[string]interface{}, mapping map[string][]string) error {
	for fieldName, alternatives := range mapping {
		// look for key corresponding to field name
		if mapKey := m.mapKeyForField(val, fieldName); len(mapKey) > 0 {
			// if map value for the key is absent or empty, try to find alternative value
			if mapVal, mOk := data[mapKey]; !mOk || system.TypeIsEmpty(mapVal) {
				// check each given alternative
				for _, altKey := range alternatives {
					// set map value for the found key with the first found alternative value
					if altVal, aOk := data[altKey]; aOk && !system.TypeIsEmpty(altVal) {
						data[mapKey] = altVal
						break
					}
				}
			}
		}
	}
	return nil
}

// Forcibly set given default values for the given fields.
// Mapping to `data` is based on `val` structure's tags.
// The mapping itself is a map [field name => default value]
func (m *Default) applyDefaultsForciblyToMap(val reflect.Value, data map[string]interface{}, defaults map[string]interface{}) error {
	for fieldName, defaultVal := range defaults {
		// look for key corresponding to field name
		if mapKey := m.mapKeyForField(val, fieldName); len(mapKey) > 0 {
			// set given value
			data[mapKey] = defaultVal
		}
	}
	return nil
}

// Concat attribute values to the one in `data` basing on mapping.
// The mapping itself is a map [field name => array of parts names]
// where part name is list of alternative names
func (m *Default) applyConcatenationToMap(val reflect.Value, data map[string]interface{}, mapping map[string][][]string) error {
	for fieldName, parts := range mapping {
		// look for key corresponding to field name
		if mapKey := m.mapKeyForField(val, fieldName); len(mapKey) > 0 {
			concatValues := make([]string, 0)
			// check each given alternative for each concatenation part
			for _, alternatives := range parts {
				for _, altKey := range alternatives {
					if altVal, aOk := data[altKey]; aOk {
						switch altVal.(type) {
						case string:
							if str := altVal.(string); len(str) > 0 {
								concatValues = append(concatValues, str)
								break
							}
						}
					}
				}
			}
			if len(concatValues) > 0 {
				data[mapKey] = strings.Join(concatValues, " ")
			}
		}
	}
	return nil
}

// If attribute is multi-valued but data value is not, then
// make from data value a slice with one item to pass type validation afterwards
func (m *Default) applyMultiValued(val reflect.Value, data map[string]interface{}, fields []string) error {
	for _, fieldName := range fields {
		// look for key corresponding to field name and map value for that key
		if mapKey := m.mapKeyForField(val, fieldName); len(mapKey) > 0 {
			if mapVal, mOk := data[mapKey]; mOk {
				// check value type
				mapValRefl := system.ReflectValue(mapVal)
				switch mapValRefl.Kind() {
				case reflect.Array, reflect.Slice:
					// as should be, so do nothing
				default:
					// make it slice
					slice := reflect.Zero(reflect.SliceOf(mapValRefl.Type()))
					slice = reflect.Append(slice, mapValRefl)
					data[mapKey] = slice.Interface()
				}
			}
		}
	}
	return nil
}

func (m *Default) applyEntitlements(val reflect.Value, data map[string]interface{}) error {
	// get entitlements from data if they exist there
	field := "entitlements"
	if m.format == "xml" {
		field = "Entitlements"
	}
	entitlements := make([]map[string]interface{}, 0)
	if sliceInterface, ok := data[field]; ok {
		switch sliceInterface.(type) {
		case []interface{}:
			for _, item := range sliceInterface.([]interface{}) {
				switch item.(type) {
				case map[string]interface{}:
					entitlements = append(entitlements, item.(map[string]interface{}))
				}
			}
		case []map[string]interface{}:
			entitlements = sliceInterface.([]map[string]interface{})
		}
	}
	// email
	if mapKey := m.mapKeyForField(val, "Emails"); len(mapKey) > 0 {
		if mapVal, mOk := data[mapKey]; !mOk || system.TypeIsEmpty(mapVal) {
			for _, entitlement := range entitlements {
				if email, ok := entitlement[mapKey]; ok {
					data[mapKey] = []map[string]interface{}{
						{
							"value": email,
						},
					}
					break
				}
			}
		}
	}
	return nil
}
