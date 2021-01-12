package middleware

import (
	"service-scim/models/resources"
	"service-scim/system"
	"strconv"
	"strings"
)

// mapping of Rabobank role names (in lower case!) to our role aliases
var globalRoleMapping = map[string]string{
	"videogesprek-admin": "manager",
	"administrator":      "admin",
	"admin":              "admin",
}

type Rabobank struct {
	Default
}

func (m *Rabobank) ProcessUser(data map[string]interface{}, id *string) (map[string]interface{}, error) {
	var err error
	reflectedUser := system.ReflectValue(new(modelsResources.User))

	// parse roles
	// (only global ones since in SCIM there no group-specific roles)
	if mapKey := m.mapKeyForField(reflectedUser, "Roles"); len(mapKey) > 0 {
		if newRoles := m.applyRoleAliasMapping(data[mapKey], globalRoleMapping); newRoles != nil {
			data[mapKey] = newRoles
		} else {
			delete(data, mapKey)
		}
	}

	// reconstruct some attributes
	if val, ok := data["formattedName"]; ok {
		data["name"] = map[string]interface{}{"formatted": val}
	}

	// make other stuff
	if data, err = m.Default.ProcessUser(data, id); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *Rabobank) ProcessGroup(data map[string]interface{}, id *string) (map[string]interface{}, error) {
	var err error
	reflectedGroup := system.ReflectValue(new(modelsResources.Group))

	// set group roles
	if mapKey := m.mapKeyForField(reflectedGroup, "Members"); len(mapKey) > 0 {
		if members, ok := data[mapKey]; ok {
			data[mapKey] = m.setGroupRoles(members, id)
		}
	}

	// make other stuff
	if data, err = m.Default.ProcessGroup(data, id); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *Rabobank) applyRoleAliasMapping(roles interface{}, roleMapping map[string]string) []map[string]interface{} {
	newRoles := make(map[string]map[string]interface{})

	// parse giving roles
	if !system.TypeIsEmpty(roles) {
		switch roles.(type) {
		case []interface{}:
			for _, role := range roles.([]interface{}) {
				switch role.(type) {
				case map[string]interface{}:
					if givenValue, ok := role.(map[string]interface{})["value"]; ok {
						givenValueStr := strings.ToLower(givenValue.(string))
						if newValue, ok := roleMapping[givenValueStr]; ok {
							newRoles[newValue] = map[string]interface{}{
								"value": newValue,
							}
						}
					}
				}
			}
		}
	}

	// make a slice
	newRolesSlice := make([]map[string]interface{}, 0)
	for _, role := range newRoles {
		newRolesSlice = append(newRolesSlice, role)
	}
	return newRolesSlice
}

func (m *Rabobank) setGroupRoles(val interface{}, idStr *string) interface{} {
	// check type of the struct of group members
	var members []map[string]interface{}
	switch val.(type) {
	case []interface{}:
		valSlice := val.([]interface{})
		members = make([]map[string]interface{}, len(valSlice))
		for i, item := range valSlice {
			switch item.(type) {
			case map[string]interface{}:
				members[i] = item.(map[string]interface{})
			default:
				return val
			}
		}
	default:
		return val
	}

	// check group id
	if idStr == nil {
		return val
	}
	var id int64
	var err error
	if id, err = strconv.ParseInt(*idStr, 10, 64); err != nil {
		return val
	}

	// get the key for the group member's role and for role's value
	reflectedMember := system.ReflectValue(new(modelsResources.GroupMember))
	var memberRolesKey string
	if memberRolesKey = m.mapKeyForField(reflectedMember, "Roles"); len(memberRolesKey) == 0 {
		return val
	}
	reflectedRole := system.ReflectValue(new(modelsResources.Role))
	var roleValKey string
	if roleValKey = m.mapKeyForField(reflectedRole, "Value"); len(roleValKey) == 0 {
		return val
	}

	// check if the group is in role mapping
	roleMapping := m.config.ScimRabobankRoleMappingByGroup()
	if groupRoles, ok := roleMapping[id]; ok {
		memberRoles := make([]map[string]interface{}, len(groupRoles))
		for i, role := range groupRoles {
			memberRoles[i] = map[string]interface{}{
				roleValKey: role,
			}
		}
		for i, _ := range members {
			members[i][memberRolesKey] = memberRoles
		}
		return members
	} else {
		return val
	}
}
