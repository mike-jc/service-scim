package modelsResources

import (
	"service-scim/models/normalization"
)

type Groups struct {
	Schemas      []string `json:"schemas" xml:"Schemas"`
	TotalResults int      `json:"totalResults" xml:"TotalResults"`
	ItemsPerPage int      `json:"itemsPerPage" xml:"itemsPerPage"`
	StartIndex   int      `json:"startIndex" xml:"startIndex"`
	Resources    []*Group `json:"Resources" xml:"Resources"`
	XMLName      struct{} `json:"-" xml:"Groups"`
}

func (l *Groups) SetNormalizationOptions(options *modelsNormalization.Options) {
	for _, g := range l.Resources {
		g.SetNormalizationOptions(options)
	}
}

func (l *Groups) AddSchemas() {
	l.Schemas = []string{"urn:ietf:params:scim:api:messages:2.0:ListResponse"}
	for _, g := range l.Resources {
		g.AddSchemas()
	}
}

func (l *Groups) AddMeta() {
	for _, g := range l.Resources {
		g.AddMeta()
	}
}
