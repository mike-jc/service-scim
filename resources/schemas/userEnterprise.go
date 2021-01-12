package resourcesSchemas

import (
	"service-scim/models/config"
)

var EnterpriseUserSchemaObject = modelsConfig.Schema{
	Id:          "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User",
	Name:        "EnterpriseUser",
	Description: "Enterprise User Schema",
	Attributes: []*modelsConfig.Attribute{
		{
			Name:        "employeeNumber",
			Type:        "string",
			MultiValued: false,
			Required:    false,
			CaseExact:   &False,
			Mutability:  "readWrite",
			Returned:    "default",
			Uniqueness:  "none",
			Description: "Numeric or alphanumeric identifier assigned to a person, typically based on order of hire or association with an organization.",
			Navigation: &modelsConfig.AttributeNavigation{
				FieldName: "EmployeeNumber",
				Path:      "employeeNumber",
				FullPath:  "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:employeeNumber",
				IndexKeys: []string{},
			},
		}, {
			Name:        "organization",
			Type:        "string",
			MultiValued: false,
			Required:    false,
			CaseExact:   &False,
			Mutability:  "readWrite",
			Returned:    "default",
			Uniqueness:  "none",
			Description: "Identifies the name of an organization.",
			Navigation: &modelsConfig.AttributeNavigation{
				FieldName: "Organization",
				Path:      "organization",
				FullPath:  "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:organization",
				IndexKeys: []string{},
			},
		}, {
			Name:        "division",
			Type:        "string",
			MultiValued: false,
			Required:    false,
			CaseExact:   &False,
			Mutability:  "readWrite",
			Returned:    "default",
			Uniqueness:  "none",
			Description: "Identifies the name of a division.",
			Navigation: &modelsConfig.AttributeNavigation{
				FieldName: "Division",
				Path:      "division",
				FullPath:  "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:division",
				IndexKeys: []string{},
			},
		}, {
			Name:        "department",
			Type:        "string",
			MultiValued: false,
			Required:    false,
			CaseExact:   &False,
			Mutability:  "readWrite",
			Returned:    "default",
			Uniqueness:  "none",
			Description: "Identifies the name of a department.",
			Navigation: &modelsConfig.AttributeNavigation{
				FieldName: "Department",
				Path:      "department",
				FullPath:  "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:department",
				IndexKeys: []string{""},
			},
		}, {
			Name:        "manager",
			Type:        "complex",
			MultiValued: false,
			Required:    false,
			Mutability:  "readWrite",
			Returned:    "default",
			Description: "The User's manager.  A complex type that optionally allows service providers to represent organizational hierarchy by referencing the 'id' attribute of another User.",
			SubAttributes: []*modelsConfig.Attribute{
				{
					Name:        "value",
					Type:        "string",
					MultiValued: false,
					Required:    true,
					CaseExact:   &False,
					Mutability:  "readWrite",
					Returned:    "default",
					Uniqueness:  "none",
					Description: "The id of the SCIM resource representing the User's manager.",
					Navigation: &modelsConfig.AttributeNavigation{
						FieldName: "Value",
						Path:      "manager.value",
						FullPath:  "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:manager.value",
						IndexKeys: []string{},
					},
				}, {
					Name:           "$ref",
					Type:           "reference",
					MultiValued:    false,
					Required:       true,
					CaseExact:      &False,
					Mutability:     "readWrite",
					Returned:       "default",
					Uniqueness:     "none",
					ReferenceTypes: []string{"User"},
					Description:    "The URI of the SCIM resource representing the User's manager.",
					Navigation: &modelsConfig.AttributeNavigation{
						FieldName: "Ref",
						Path:      "manager.$ref",
						FullPath:  "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:manager.$ref",
						IndexKeys: []string{},
					},
				}, {
					Name:        "displayName",
					Type:        "string",
					MultiValued: false,
					Required:    false,
					CaseExact:   &False,
					Mutability:  "readOnly",
					Returned:    "default",
					Uniqueness:  "none",
					Description: "The displayName of the User's manager.",
					Navigation: &modelsConfig.AttributeNavigation{
						FieldName: "DisplayName",
						Path:      "manager.displayName",
						FullPath:  "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:manager.displayName",
						IndexKeys: []string{},
					},
				},
			},
			Navigation: &modelsConfig.AttributeNavigation{
				FieldName: "Manager",
				Path:      "manager",
				FullPath:  "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:manager",
				IndexKeys: []string{},
			},
		},
	},
	Meta: &modelsConfig.SchemaMeta{
		Location:     "/Schemas/urn:ietf:params:scim:schemas:extension:enterprise:2.0:User",
		ResourceType: "Schema",
	},
}

var EnterpriseExtensionAttribute = &modelsConfig.Attribute{
	Name:                 EnterpriseUserSchemaObject.Id,
	Type:                 "complex",
	MultiValued:          false,
	Required:             false,
	IsExtensionAttribute: true,
	Mutability:           "readWrite",
	Returned:             "default",
	Uniqueness:           "none",
	SubAttributes:        EnterpriseUserSchemaObject.Attributes,
	Navigation: &modelsConfig.AttributeNavigation{
		FieldName: "EnterpriseExtension",
		Path:      EnterpriseUserSchemaObject.Id,
		FullPath:  EnterpriseUserSchemaObject.Id,
		IndexKeys: []string{},
	},
}