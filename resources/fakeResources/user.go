package fakeResources

import (
	"service-scim/models/resources"
)

var UserJackObject = &modelsResources.User{
	Schemas: []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
	Id:      123,
	Emails: []*modelsResources.Email{
		{
			Value: "jack@example.com",
			Type:  "work",
		},
	},
	Name: &modelsResources.UserName{
		Formatted: "Jack Green",
	},
	ExternalId: "51db7ce8-75d7-4c81-ad05-c98a9233811e",
	Title:      "Founder",
	Locale:     "en-US",
	Timezone:   "Europe/Kiev",
	Active:     true,
	Addresses: []*modelsResources.Address{
		{
			StreetAddress: "157 Main Street",
			Locality:      "Kiev",
			Country:       "UA",
		},
	},
	Photos: []*modelsResources.Photo{
		{
			Value: "https://photos.example.com/profilephoto/72930000000/F.png",
			Type:  "photo",
		}, {
			Value: "https://photos.example.com/profilephoto/72930000000/T.png",
			Type:  "thumbnail",
		},
	},
	Roles: []*modelsResources.Role{
		{
			Value:   "manager",
			Display: "Global manager",
		},
	},
	Meta: &modelsResources.AbstractMeta{
		Location:     "/v2/users/123",
		ResourceType: "User",
	},
}

var UserAnneObject = &modelsResources.User{
	Schemas: []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
	Id:      43562,
	Emails: []*modelsResources.Email{
		{
			Value:   "anne@example.com",
			Primary: true,
		},
	},
	Name: &modelsResources.UserName{
		Formatted: "Anne Smith",
	},
	ExternalId: "ef3b507c-d973-4d3f-821a-1fe5277f0af4",
	ProfileUrl: "https://login.example.com/annex",
	PhoneNumbers: []*modelsResources.PhoneNumber{
		{
			Value: "123-456-0789",
		},
	},
	Title:    "Employee",
	Locale:   "en-US",
	Timezone: "Asia/Shanghai",
	Active:   true,
	Addresses: []*modelsResources.Address{
		{
			StreetAddress: "123 Home Street",
			Locality:      "Shanghai",
			PostalCode:    "200124",
			Country:       "CN",
		},
	},
	Photos: []*modelsResources.Photo{
		{
			Value: "https://photos.example.com/profilephoto/75890000001/F.png",
			Type:  "photo",
		},
	},
	Meta: &modelsResources.AbstractMeta{
		Location:     "/v2/users/43562",
		ResourceType: "User",
	},
}

var UsersObject = &modelsResources.Users{
	Schemas:      []string{"urn:ietf:params:scim:api:messages:2.0:ListResponse"},
	TotalResults: 2,
	ItemsPerPage: 2,
	StartIndex:   1,
	Resources:    []*modelsResources.User{UserAnneObject, UserJackObject},
}
