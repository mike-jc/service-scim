package modelsResources

import (
	"encoding/json"
	"encoding/xml"
	"service-scim/resources/schemas"
)

type EnterpriseUser struct {
	Abstract
	EmployeeNumber string                 `json:"employeeNumber,omitempty" xml:"EmployeeNumber,omitempty"`
	Organization   string                 `json:"organization,omitempty" xml:"Organization,omitempty"`
	Division       string                 `json:"division,omitempty" xml:"Division,omitempty"`
	Department     string                 `json:"department,omitempty" xml:"Department,omitempty"`
	Manager        *EnterpriseUserManager `json:"manager,omitempty" xml:"Manager,omitempty"`
	XMLName        struct{}               `json:"-" xml:"EnterpriseUser"`
}

type EnterpriseUserManager struct {
	Value       string   `json:"value,omitempty" xml:"Value,omitempty"`
	Ref         string   `json:"$ref,omitempty" xml:"Ref,omitempty"`
	DisplayName string   `json:"displayName,omitempty" xml:"DisplayName,omitempty"`
	XMLName     struct{} `json:"-" xml:"EnterpriseUserManager"`
}

func (u *EnterpriseUser) IsEmpty() bool {
	return (EnterpriseUser{}) == *u
}

func (u *EnterpriseUser) MarshalJSON() ([]byte, error) {
	if normalized, err := u.normalized(); err != nil {
		return nil, err
	} else {
		return json.Marshal(normalized)
	}
}

func (u *EnterpriseUser) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if normalized, err := u.normalized(); err != nil {
		return err
	} else {
		return e.EncodeElement(normalized, start)
	}
}

// Filter excluded attributes and keep included and required ones
func (u *EnterpriseUser) normalized() (interface{}, error) {
	if normalized, err := u.normalizator.Normalize(u, resourcesSchemas.EnterpriseUserSchemaObject.Attributes, u.normalizationOptions.Included(), u.normalizationOptions.Excluded()); err != nil {
		return nil, err
	} else {
		return normalized.Interface(), nil
	}
}
