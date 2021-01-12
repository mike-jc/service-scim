package modelsSdkKeeper

import (
	"service-scim/system"
)

type Role struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

func RoleIsValid(role map[string]interface{}) bool {
	return !system.MapValueIsEmpty(role, "alias")
}
