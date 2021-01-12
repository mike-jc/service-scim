package modelsResources

import (
	"service-scim/models/normalization"
	"service-scim/services/normalization"
)

type Interface interface {
	NormalizationOptions() *modelsNormalization.Options
	SetNormalizationOptions(options *modelsNormalization.Options)
	AddSchemas()
	AddMeta()
}

type Abstract struct {
	Interface `json:"-" xml:"-"`

	normalizationOptions *modelsNormalization.Options
	normalizator         *normalization.ByAttribute
}

type AbstractMeta struct {
	Created      string `json:"created,omitempty" xml:"Created,omitempty"`
	LastModified string `json:"lastModified,omitempty" xml:"LastModified,omitempty"`
	Location     string `json:"location" xml:"Location"`
	ResourceType string `json:"resourceType" xml:"ResourceType"`
}

func (a *Abstract) NormalizationOptions() *modelsNormalization.Options {
	return a.normalizationOptions
}

func (a *Abstract) SetNormalizationOptions(options *modelsNormalization.Options) {
	if a.normalizator == nil {
		a.normalizator = new(normalization.ByAttribute)
	}
	a.normalizationOptions = options
}

func (a *Abstract) AddSchemas() {
}

func (a *Abstract) AddMeta() {
}
