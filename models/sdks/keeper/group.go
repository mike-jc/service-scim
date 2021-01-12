package modelsSdkKeeper

import (
	"fmt"
	"service-scim/system"
)

type Group struct {
	Id        int64   `json:"id"`
	ScimId    string  `json:"scimId"`
	Name      string  `json:"name"`
	UpdatedAt string  `json:"updatedAt,omitempty"`
	Users     []*User `json:"users,omitempty"`
}

func (g *Group) IsValid() bool {
	return g.Id > 0
}

func (g *Group) Ref() string {
	return fmt.Sprintf("/%s/groups/%d", system.AppVersion(), g.Id)
}

func GroupIsValid(group map[string]interface{}) bool {
	return !system.MapValueIsEmpty(group, "id")
}
