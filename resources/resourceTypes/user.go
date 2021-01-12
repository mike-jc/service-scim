package resourcesTypes

import (
	"service-scim/models/config"
)

var UserTypeObject = modelsConfig.ResourceType{
	Schemas:     []string{"urn:ietf:params:scim:schemas:core:2.0:ResourceType"},
	Id:          "User",
	Name:        "User",
	Endpoint:    "/Users",
	Description: "User account: https://tools.ietf.org/html/rfc7643#section-8.7.1",
	Schema:      "urn:ietf:params:scim:schemas:core:2.0:User",
	SchemaExtensions: []*modelsConfig.SchemaExtension{
		{
			Schema:   "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User",
			Required: false,
		},
	},
	Meta: &modelsConfig.ResourceTypeMeta{
		Location:     "/ResourceTypes/User",
		ResourceType: "ResourceType",
	},
}
