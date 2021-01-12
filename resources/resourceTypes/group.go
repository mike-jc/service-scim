package resourcesTypes

import (
	"service-scim/models/config"
)

var GroupTypeObject = modelsConfig.ResourceType{
	Schemas:     []string{"urn:ietf:params:scim:schemas:core:2.0:ResourceType"},
	Id:          "Group",
	Name:        "Group",
	Endpoint:    "/Groups",
	Description: "Group of users: https://tools.ietf.org/html/rfc7643#section-8.7.1",
	Schema:      "urn:ietf:params:scim:schemas:core:2.0:Group",
	Meta: &modelsConfig.ResourceTypeMeta{
		Location:     "/ResourceTypes/Group",
		ResourceType: "ResourceType",
	},
}
