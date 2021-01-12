package fakeResources

import (
	"service-scim/models/resources"
)

var GroupManagersObject = &modelsResources.Group{
	Schemas:     []string{"urn:ietf:params:scim:schemas:core:2.0:Group"},
	DisplayName: "Managers",
	Id:          34,
	ExternalId:  "6de1cc84-d9d0-4fb4-a6a4-c85c52675887",
	Members: []*modelsResources.GroupMember{
		{
			Value: "43562",
			Ref:   "/v2/users/43562",
			Type:  "User",
		},
	},
	Meta: &modelsResources.AbstractMeta{
		Location:     "/v2/groups/34",
		ResourceType: "Group",
	},
}

var GroupOperatorsObject = &modelsResources.Group{
	Schemas:     []string{"urn:ietf:params:scim:schemas:core:2.0:Group"},
	Id:          57,
	DisplayName: "Operators",
	ExternalId:  "38c5d26b-070a-42f5-89c8-882326f50b8b",
	Members: []*modelsResources.GroupMember{
		{
			Value: "43562",
			Ref:   "/v2/users/43562",
			Type:  "User",
		}, {
			Value: "123",
			Ref:   "/v2/users/123",
			Type:  "User",
		},
	},
	Meta: &modelsResources.AbstractMeta{
		Location:     "/v2/groups/57",
		ResourceType: "Group",
	},
}

var GroupsObject = &modelsResources.Groups{
	Schemas:      []string{"urn:ietf:params:scim:api:messages:2.0:ListResponse"},
	TotalResults: 2,
	ItemsPerPage: 2,
	StartIndex:   1,
	Resources:    []*modelsResources.Group{GroupOperatorsObject, GroupManagersObject},
}
