package controllers

import (
	"fmt"
	"service-scim/services"
)

type AbstractInstanceController struct {
	AbstractController
}

func (c *AbstractInstanceController) Prepare() {
	c.scimConfig = nil
	c.scimDomain = c.Ctx.Input.Domain()
	c.domain = services.InstanceDomainFromScimDomain(c.scimDomain)

	c.GetLogger().SetInstance(c.domain)

	if config, err := services.NewConfig(c.domain); err != nil {
		c.ShowError(fmt.Sprintf("Can not get configuration for domain %s, request domain is %s", c.domain, c.scimDomain), 500, "conf.not_readable", err.Error(), false)
		c.Finish()
		c.StopRun()
	} else {
		c.scimConfig = config
	}

	if !c.scimConfig.IsEnabled() {
		c.ShowError(fmt.Sprintf("SCIM is not enabled on instance %s", c.domain), 500, "conf.not_enabled", "", false)
		c.Finish()
		c.StopRun()
	}
}
