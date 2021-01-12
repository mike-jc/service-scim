package modelsResources

import (
	"service-scim/models/normalization"
)

type Users struct {
	Schemas      []string `json:"schemas" xml:"Schemas"`
	TotalResults int      `json:"totalResults" xml:"TotalResults"`
	ItemsPerPage int      `json:"itemsPerPage" xml:"itemsPerPage"`
	StartIndex   int      `json:"startIndex" xml:"startIndex"`
	Resources    []*User  `json:"Resources" xml:"Resources"`
	XMLName      struct{} `json:"-" xml:"Users"`
}

func (l *Users) SetNormalizationOptions(options *modelsNormalization.Options) {
	for _, u := range l.Resources {
		u.SetNormalizationOptions(options)
	}
}

func (l *Users) AddSchemas() {
	l.Schemas = []string{"urn:ietf:params:scim:api:messages:2.0:ListResponse"}
	for _, u := range l.Resources {
		u.AddSchemas()
	}
}

func (l *Users) AddMeta() {
	for _, u := range l.Resources {
		u.AddMeta()
	}
}
