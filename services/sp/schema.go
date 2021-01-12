package serviceSp

import (
	"errors"
	"fmt"
	"gitlab.com/24sessions/sdk-go-configurator/data"
	"service-scim/models/config"
	"service-scim/resources/schemas"
	"service-scim/system"
)

const SchemaUnknown = 0
const SchemaUser = 1
const SchemaUserEnterprise = 2
const SchemaGroup = 3

var Schemas = map[int]string{
	SchemaUser:           "urn:ietf:params:scim:schemas:core:2.0:User",
	SchemaUserEnterprise: "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User",
	SchemaGroup:          "urn:ietf:params:scim:schemas:core:2.0:Group",
}

type Schema struct {
	scimConfig *sdksData.ScimContainer
	baseUrl    string
}

func (s *Schema) SetScimConfig(scimConfig *sdksData.ScimContainer) {
	s.scimConfig = scimConfig
}

func (s *Schema) SetBaseUrl(baseUrl string) {
	s.baseUrl = baseUrl
}

func (s *Schema) IdFromUrn(givenUrn string) (id int, err error) {
	var urn string
	for id, urn = range Schemas {
		if urn == givenUrn {
			return
		}
	}

	id = SchemaUnknown
	err = errors.New(fmt.Sprintf("Unknown schema URN: %s", givenUrn))
	return
}

func (s *Schema) SchemaById(id int) (sc *modelsConfig.Schema, err error) {
	if id == SchemaUser {
		resourcesSchemas.UserSchemaObject.Meta.Location = s.baseUrl + "/" + system.AppVersion() + "/Schemas/" + Schemas[SchemaUser]
		return &resourcesSchemas.UserSchemaObject, nil
	} else if id == SchemaUserEnterprise {
		resourcesSchemas.EnterpriseUserSchemaObject.Meta.Location = s.baseUrl + "/" + system.AppVersion() + "/Schemas/" + Schemas[SchemaUserEnterprise]
		return &resourcesSchemas.EnterpriseUserSchemaObject, nil
	} else if id == SchemaGroup {
		resourcesSchemas.GroupSchemaObject.Meta.Location = s.baseUrl + "/" + system.AppVersion() + "/Schemas/" + Schemas[SchemaGroup]
		return &resourcesSchemas.GroupSchemaObject, nil
	} else {
		return nil, errors.New(fmt.Sprintf("Unknown schema ID: %d", id))
	}
}

func (s *Schema) Schemas() (list *modelsConfig.Schemas, err error) {
	schemas := make([]*modelsConfig.Schema, 0)

	if data, err := s.SchemaById(SchemaUser); err != nil {
		return nil, err
	} else {
		schemas = append(schemas, data)
	}

	if data, err := s.SchemaById(SchemaUserEnterprise); err != nil {
		return nil, err
	} else {
		schemas = append(schemas, data)
	}

	if data, err := s.SchemaById(SchemaGroup); err != nil {
		return nil, err
	} else {
		schemas = append(schemas, data)
	}

	return &modelsConfig.Schemas{
		Schemas:      []string{"urn:ietf:params:scim:api:messages:2.0:ListResponse"},
		TotalResults: len(schemas),
		ItemsPerPage: len(schemas),
		StartIndex:   1,
		Resources:    schemas,
	}, nil
}
