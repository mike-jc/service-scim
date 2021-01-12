package controllers

import (
	"service-scim/services/sp"
)

type ServiceController struct {
	AbstractInstanceController
}

// @Title Service Provider configuration
// @Description List of supported methods, authorization types etc.
// @router /ServiceProviderConfigs [get]
func (c *ServiceController) SpConfiguration() {
	s := new(serviceSp.ServiceConfig)
	s.SetScimConfig(c.scimConfig)
	s.SetBaseUrl(c.BasicURL())

	if config, err := s.Config(); err != nil {
		c.ShowError("Can not get SP configuration", 500, "app.config.invalid", err.Error(), false)
	} else {
		c.ServeResponse(config)
	}
}

// @Title Types of resources
// @router /ResourceTypes [get]
func (c *ServiceController) ResourceTypes() {
	t := new(serviceSp.ResourceType)
	t.SetBaseUrl(c.BasicURL())

	data, _ := t.Types()
	c.ServeResponse(data)
}

// @Title Types of resources
// @router /ResourceTypes/:name [get]
func (c *ServiceController) ResourceTypeByName() {
	t := new(serviceSp.ResourceType)
	t.SetBaseUrl(c.BasicURL())

	name := c.Ctx.Input.Param(":name")
	if id, err := t.IdFromName(name); err != nil {
		c.ShowError("Unknown resource type", 400, "app.unknown_resource_type", "Name of requested type is "+name, false)
	} else {
		data, _ := t.TypeById(id)
		c.ServeResponse(data)
	}
}

// @Title Schemas of users and groups
// @router /schemas [get]
func (c *ServiceController) Schemas() {
	s := new(serviceSp.Schema)
	s.SetBaseUrl(c.BasicURL())

	data, _ := s.Schemas()
	c.ServeResponse(data)
}

// @Title Schemas of users and groups
// @router /schemas/:urn [get]
func (c *ServiceController) SchemaByUrn() {
	s := new(serviceSp.Schema)
	s.SetBaseUrl(c.BasicURL())

	urn := c.Ctx.Input.Param(":urn")
	if id, err := s.IdFromUrn(urn); err != nil {
		c.ShowError("Unknown schema", 400, "app.unknown_schema", "URN of requested schema is "+urn, false)
	} else {
		data, _ := s.SchemaById(id)
		c.ServeResponse(data)
	}
}
