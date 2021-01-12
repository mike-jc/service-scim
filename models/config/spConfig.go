package modelsConfig

type SpConfig struct {
	Schemas               []string                 `json:"schemas" xml:"Schemas"`
	Id                    string                   `json:"id,omitempty" xml:"id,omitempty"`
	DocumentationUri      string                   `json:"documentationUri,omitempty" xml:"DocumentationUri,omitempty"`
	Patch                 *PatchOperation          `json:"patch,omitempty" xml:"Patch,omitempty"`
	Bulk                  *BulkOperation           `json:"bulk,omitempty" xml:"Bulk,omitempty"`
	Filter                *FilterOperation         `json:"filter,omitempty" xml:"Filter,omitempty"`
	ChangePassword        *ChangePasswordOperation `json:"changePassword,omitempty" xml:"ChangePassword,omitempty"`
	Sort                  *SortOperation           `json:"sort,omitempty" xml:"Sort,omitempty"`
	ETag                  *ETag                    `json:"etag,omitempty" xml:"ETag,omitempty"`
	AuthenticationSchemes []*AuthenticationScheme  `json:"authenticationSchemes,omitempty" xml:"AuthenticationSchemes,omitempty"`
	Meta                  *SpConfigMeta            `json:"meta,omitempty" xml:"Meta,omitempty"`
	XMLName               struct{}                 `json:"-" xml:"ServiceProviderConfig"`
}

type PatchOperation struct {
	Supported bool `json:"supported" xml:"supported,attr"`
}

type BulkOperation struct {
	Supported      bool  `json:"supported" xml:"supported,attr"`
	MaxOperations  int   `json:"maxOperations" xml:"maxOperations,attr"`
	MaxPayloadSize int64 `json:"maxPayloadSize" xml:"maxPayloadSize,attr"`
}

type FilterOperation struct {
	Supported  bool `json:"supported" xml:"supported,attr"`
	MaxResults int  `json:"maxResults" xml:"maxResults,attr"`
}

type ChangePasswordOperation struct {
	Supported bool `json:"supported" xml:"supported,attr"`
}

type SortOperation struct {
	Supported bool `json:"supported" xml:"supported,attr"`
}

type ETag struct {
	Supported bool `json:"supported" xml:"supported,attr"`
}

type AuthenticationScheme struct {
	Name             string `json:"name" xml:"Name"`
	Description      string `json:"description,omitempty" xml:"Description,omitempty"`
	SpecUri          string `json:"specUri,omitempty" xml:"SpecUri,omitempty"`
	DocumentationUri string `json:"documentationUri,omitempty" xml:"DocumentationUri,omitempty"`
	Type             string `json:"type" xml:"Type"`
	Primary          bool   `json:"primary,omitempty" xml:"Primary,omitempty"`
}

type SpConfigMeta struct {
	ResourceType string `json:"resourceType" xml:"ResourceType"`
	Location     string `json:"location" xml:"Location"`
	Created      string `json:"created" xml:"Created"`
	LastModified string `json:"lastModified" xml:"LastModified"`
	Version      string `json:"version,omitempty" xml:"Version,omitempty"`
}
