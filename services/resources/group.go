package resources

import (
	"fmt"
	"github.com/astaxie/beego/context"
	sdksData "gitlab.com/24sessions/sdk-go-configurator/data"
	"service-scim/models/normalization"
	"service-scim/models/resources"
	"service-scim/resources/schemas"
	"service-scim/services/repositories/group"
)

type Group struct {
	Abstract

	repository repositoriesGroup.Interface
}

func (g *Group) SetRepository(repository repositoriesGroup.Interface) {
	g.repository = repository
	g.uniquenessValidator.SetRepository(repository)
}

func (g *Group) SetInstanceDomain(domain string) {
	g.repository.SetInstanceDomain(domain)
}

func (g *Group) SetFormat(format string) {
	g.Abstract.SetFormat(format)
	g.repository.SetFormat(format)
}

func NewGroupService(r repositoriesGroup.Interface, ctx *context.Context, config *sdksData.ScimContainer, domain, format string) *Group {
	group := new(Group)
	group.Init()
	group.SetRepository(r)
	group.AttributesFromRequest(ctx)
	group.SetConfig(config)
	group.SetInstanceDomain(domain)
	group.SetFormat(format)

	return group
}

func (g *Group) List(offset, limit int) (*modelsResources.Groups, error) {
	normalizationOptions := new(modelsNormalization.Options)
	if err := normalizationOptions.Parse(g.includedAttributes, g.excludedAttributes); err != nil {
		return nil, err
	}

	if totalCount, resources, err := g.repository.List(offset, limit); err != nil {
		return nil, err
	} else {
		if resources == nil {
			resources = make([]*modelsResources.Group, 0)
		}
		groups := &modelsResources.Groups{
			TotalResults: totalCount,
			ItemsPerPage: limit,
			StartIndex:   offset + 1, // offset is 0-based, but startIndex is 1-based
			Resources:    resources,
		}
		groups.SetNormalizationOptions(normalizationOptions)
		groups.AddSchemas()
		groups.AddMeta()
		return groups, nil
	}
}

func (g *Group) ById(id string) (*modelsResources.Group, error) {
	normalizationOptions := new(modelsNormalization.Options)
	if err := normalizationOptions.Parse(g.includedAttributes, g.excludedAttributes); err != nil {
		return nil, err
	}

	if group, err := g.repository.ById(id); err != nil {
		return nil, err
	} else {
		group.SetNormalizationOptions(normalizationOptions)
		group.AddSchemas()
		group.AddMeta()
		return group, nil
	}
}

func (g *Group) Create(data map[string]interface{}) (group *modelsResources.Group, err error) {
	normalizationOptions := new(modelsNormalization.Options)
	if err = normalizationOptions.Parse(g.includedAttributes, g.excludedAttributes); err != nil {
		return
	}

	if data, err = g.middlewareProcessing(data, nil); err != nil {
		return
	}
	if err = g.validate(data, nil); err != nil {
		return
	}

	if group, err = g.repository.Create(data); err != nil {
		return
	}

	group.SetNormalizationOptions(normalizationOptions)
	group.AddSchemas()
	group.AddMeta()
	return
}

func (g *Group) Modify(id string, modification *modelsResources.Modification) (group *modelsResources.Group, err error) {
	normalizationOptions := new(modelsNormalization.Options)
	if err = normalizationOptions.Parse(g.includedAttributes, g.excludedAttributes); err != nil {
		return
	}

	var data map[string]interface{}
	if data, err = g.applyPatch(id, modification); err != nil {
		return
	}
	if data, err = g.middlewareProcessing(data, &id); err != nil {
		return
	}
	if err = g.validate(data, &id); err != nil {
		return
	}

	if group, err = g.repository.Update(id, data); err != nil {
		return
	}

	group.SetNormalizationOptions(normalizationOptions)
	group.AddSchemas()
	group.AddMeta()
	return
}

func (g *Group) Replace(id string, data map[string]interface{}) (group *modelsResources.Group, err error) {
	normalizationOptions := new(modelsNormalization.Options)
	if err = normalizationOptions.Parse(g.includedAttributes, g.excludedAttributes); err != nil {
		return
	}

	if data, err = g.middlewareProcessing(data, &id); err != nil {
		return
	}
	if err = g.validate(data, &id); err != nil {
		return
	}

	if group, err = g.repository.Replace(id, data); err != nil {
		return
	}

	group.SetNormalizationOptions(normalizationOptions)
	group.AddSchemas()
	group.AddMeta()
	return
}

func (g *Group) Disable(id string) error {
	return g.repository.Disable(id)
}

func (g *Group) middlewareProcessing(data map[string]interface{}, id *string) (map[string]interface{}, error) {
	if engine, err := g.MiddlewareEngine(); err != nil {
		return data, fmt.Errorf("Cannot process group data with middleware: %s", err.Error())
	} else {
		return engine.ProcessGroup(data, id)
	}
}

func (g *Group) validate(data map[string]interface{}, id *string) error {
	schema := &resourcesSchemas.GroupSchemaObject

	if err := g.typeValidator.Validate(data, schema); err != nil {
		return err
	}
	if err := g.caseValidator.Validate(data, schema); err != nil {
		return err
	}
	if err := g.readOnlyValidator.ValidateGroup(data, g.format, id); err != nil {
		return err
	}
	if err := g.requireValidator.Validate(data, schema); err != nil {
		return err
	}

	if id != nil {
		if group, rErr := g.repository.ById(*id); rErr != nil {
			return rErr
		} else {
			if err := g.mutabilityValidator.Validate(data, group.ToMap(g.format), schema); err != nil {
				return err
			}
		}
	}

	if err := g.uniquenessValidator.Validate(data, id, schema); err != nil {
		return err
	}
	return nil
}

func (g *Group) applyPatch(id string, modification *modelsResources.Modification) (map[string]interface{}, error) {

	// get original entity
	if origGroup, err := g.repository.ById(id); err != nil {
		return nil, err
	} else {
		groupMap := origGroup.ToMap(g.format)

		// try to apply all operations and patches to the entity
		for _, patch := range modification.Operations {
			if pErr := ApplyPatch(patch, groupMap, &resourcesSchemas.GroupSchemaObject); pErr != nil {
				return nil, pErr
			}
		}
		return groupMap.(map[string]interface{}), nil
	}
}
