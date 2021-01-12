package resources

import (
	"fmt"
	"github.com/astaxie/beego/context"
	sdksData "gitlab.com/24sessions/sdk-go-configurator/data"
	"service-scim/models/normalization"
	"service-scim/models/resources"
	"service-scim/resources/schemas"
	"service-scim/services/filtering"
	"service-scim/services/repositories/user"
)

type User struct {
	Abstract

	repository repositoriesUser.Interface
}

func (u *User) SetRepository(repository repositoriesUser.Interface) {
	u.repository = repository
	u.uniquenessValidator.SetRepository(repository)
}

func (u *User) SetInstanceDomain(domain string) {
	u.repository.SetInstanceDomain(domain)
}

func (u *User) SetFormat(format string) {
	u.Abstract.SetFormat(format)
	u.repository.SetFormat(format)
}

func NewUserService(r repositoriesUser.Interface, ctx *context.Context, config *sdksData.ScimContainer, domain, format string) *User {
	user := new(User)
	user.Init()
	user.SetRepository(r)
	user.AttributesFromRequest(ctx)
	user.SetConfig(config)
	user.SetInstanceDomain(domain)
	user.SetFormat(format)

	return user
}

func (u *User) List(offset, limit int, filterStr string) (*modelsResources.Users, error) {
	normalizationOptions := new(modelsNormalization.Options)
	if err := normalizationOptions.Parse(u.includedAttributes, u.excludedAttributes); err != nil {
		return nil, err
	}

	if filterMap, err := u.mapFromFilterString(filterStr); err != nil {
		return nil, err
	} else if totalCount, resources, err := u.repository.List(offset, limit, filterMap); err != nil {
		return nil, err
	} else {
		if resources == nil {
			resources = make([]*modelsResources.User, 0)
		}
		users := &modelsResources.Users{
			TotalResults: totalCount,
			ItemsPerPage: limit,
			StartIndex:   offset + 1, // offset is 0-based, but startIndex is 1-based
			Resources:    resources,
		}
		users.SetNormalizationOptions(normalizationOptions)
		users.AddSchemas()
		users.AddMeta()
		return users, nil
	}
}

func (u *User) ById(id string) (*modelsResources.User, error) {
	normalizationOptions := new(modelsNormalization.Options)
	if err := normalizationOptions.Parse(u.includedAttributes, u.excludedAttributes); err != nil {
		return nil, err
	}

	if user, err := u.repository.ById(id); err != nil {
		return nil, err
	} else {
		user.SetNormalizationOptions(normalizationOptions)
		user.AddSchemas()
		user.AddMeta()
		return user, nil
	}
}

func (u *User) Create(data map[string]interface{}) (user *modelsResources.User, err error) {
	normalizationOptions := new(modelsNormalization.Options)
	if err = normalizationOptions.Parse(u.includedAttributes, u.excludedAttributes); err != nil {
		return
	}

	if data, err = u.middlewareProcessing(data, nil); err != nil {
		return
	}
	if err = u.validate(data, nil); err != nil {
		return
	}

	if user, err = u.repository.Create(data); err != nil {
		return
	}

	user.SetNormalizationOptions(normalizationOptions)
	user.AddSchemas()
	user.AddMeta()
	return
}

func (u *User) Modify(id string, modification *modelsResources.Modification) (user *modelsResources.User, err error) {
	normalizationOptions := new(modelsNormalization.Options)
	if err = normalizationOptions.Parse(u.includedAttributes, u.excludedAttributes); err != nil {
		return
	}

	var data map[string]interface{}
	if data, err = u.applyPatch(id, modification); err != nil {
		return
	}
	if data, err = u.middlewareProcessing(data, &id); err != nil {
		return
	}
	if err = u.validate(data, &id); err != nil {
		return
	}

	if user, err = u.repository.Update(id, data); err != nil {
		return
	}

	user.SetNormalizationOptions(normalizationOptions)
	user.AddSchemas()
	user.AddMeta()
	return
}

func (u *User) Replace(id string, data map[string]interface{}) (user *modelsResources.User, err error) {
	normalizationOptions := new(modelsNormalization.Options)
	if err = normalizationOptions.Parse(u.includedAttributes, u.excludedAttributes); err != nil {
		return
	}

	if data, err = u.middlewareProcessing(data, &id); err != nil {
		return
	}
	if err = u.validate(data, &id); err != nil {
		return
	}

	if user, err = u.repository.Replace(id, data); err != nil {
		return
	}

	user.SetNormalizationOptions(normalizationOptions)
	user.AddSchemas()
	user.AddMeta()
	return
}

func (u *User) Block(id string) error {
	return u.repository.Block(id)
}

func (u *User) mapFromFilterString(filterStr string) (map[string]interface{}, error) {
	if parsedFilter, err := filtering.StringToFieldsMap(filterStr, &resourcesSchemas.UserSchemaObject); err != nil {
		return nil, err
	} else if parsedFilter != nil {
		resultFilter := make(map[string]interface{})

		if userName, ok := parsedFilter["userName"].(string); ok && len(userName) > 0 {
			resultFilter["scimId"] = userName
		} else if externalId, ok := parsedFilter["externalId"].(string); ok && len(externalId) > 0 {
			resultFilter["scimId"] = externalId
		}

		return resultFilter, nil
	}
	return nil, nil
}

func (u *User) middlewareProcessing(data map[string]interface{}, id *string) (map[string]interface{}, error) {
	if engine, err := u.MiddlewareEngine(); err != nil {
		return data, fmt.Errorf("Cannot process user data with middleware: %s", err.Error())
	} else {
		return engine.ProcessUser(data, id)
	}
}

func (u *User) validate(data map[string]interface{}, id *string) error {
	schema := &resourcesSchemas.UserSchemaObject

	if err := u.typeValidator.Validate(data, schema); err != nil {
		return err
	}
	if err := u.caseValidator.Validate(data, schema); err != nil {
		return err
	}
	if err := u.readOnlyValidator.ValidateUser(data, u.format, id); err != nil {
		return err
	}
	if err := u.requireValidator.Validate(data, schema); err != nil {
		return err
	}

	if id != nil {
		if user, rErr := u.repository.ById(*id); rErr != nil {
			return rErr
		} else {
			if err := u.mutabilityValidator.Validate(data, user.ToMap(u.format), schema); err != nil {
				return err
			}
		}
	}

	if err := u.uniquenessValidator.Validate(data, id, schema); err != nil {
		return err
	}
	return nil
}

func (u *User) applyPatch(id string, modification *modelsResources.Modification) (map[string]interface{}, error) {

	// get original entity
	if origUser, err := u.repository.ById(id); err != nil {
		return nil, err
	} else {
		userMap := origUser.ToMap(u.format)

		// try to apply all operations and patches to the entity
		for _, patch := range modification.Operations {
			if pErr := ApplyPatch(patch, userMap, &resourcesSchemas.UserSchemaObject); pErr != nil {
				return nil, pErr
			}
		}
		return userMap.(map[string]interface{}), nil
	}
}
