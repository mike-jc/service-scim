package modelsNormalization

import (
	"service-scim/services/navigation"
)

type Options struct {
	includedAttributes []navigation.Path
	excludedAttributes []navigation.Path
}

func (o *Options) Included() []navigation.Path {
	return o.includedAttributes
}

func (o *Options) Excluded() []navigation.Path {
	return o.excludedAttributes
}

func (o *Options) Parse(includedAttributes []string, excludedAttributes []string) error {
	if attributes, err := MakeAttributes(includedAttributes); err != nil {
		return err
	} else {
		o.includedAttributes = attributes
	}
	if attributes, err := MakeAttributes(excludedAttributes); err != nil {
		return err
	} else {
		o.excludedAttributes = attributes
	}
	return nil
}

func MakeAttributes(attributes []string) ([]navigation.Path, error) {
	normAttributes := make([]navigation.Path, 0)
	for _, attr := range attributes {
		if len(attr) > 0 {
			if p, err := navigation.NewPath(attr); err != nil {
				return normAttributes, err
			} else {
				normAttributes = append(normAttributes, p)
			}
		}
	}
	return normAttributes, nil
}

func NewEmptyNormalizationOption() *Options {
	return &Options{
		includedAttributes: []navigation.Path{},
		excludedAttributes: []navigation.Path{},
	}
}
