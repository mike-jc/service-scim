package modelsConfig

type ResourceType struct {
	Schemas          []string           `json:"schemas" xml:"Schemas"`
	Id               string             `json:"id" xml:"Id"`
	Name             string             `json:"name,omitempty" xml:"Name,omitempty"`
	Endpoint         string             `json:"endpoint" xml:"Endpoint"`
	Description      string             `json:"description,omitempty" xml:"Description,omitempty"`
	Schema           string             `json:"schema" xml:"Schema"`
	SchemaExtensions []*SchemaExtension `json:"schemaExtensions,omitempty" xml:"SchemaExtensions,omitempty"`
	Meta             *ResourceTypeMeta  `json:"meta,omitempty" xml:"Meta,omitempty"`
	XMLName          struct{}           `json:"-" xml:"ResourceType"`
}

type ResourceTypes struct {
	Schemas      []string        `json:"schemas" xml:"Schemas"`
	TotalResults int             `json:"totalResults" xml:"totalResults"`
	ItemsPerPage int             `json:"itemsPerPage" xml:"ItemsPerPage"`
	StartIndex   int             `json:"startIndex" xml:"StartIndex"`
	Resources    []*ResourceType `json:"Resources" xml:"Resources"`
	XMLName      struct{}        `json:"-" xml:"ResourceTypes"`
}

type SchemaExtension struct {
	Schema   string `json:"schema" xml:"Schema"`
	Required bool   `json:"required" xml:"Required"`
}

type ResourceTypeMeta struct {
	Location     string `json:"location" xml:"Location"`
	ResourceType string `json:"resourceType" xml:"ResourceType"`
}
