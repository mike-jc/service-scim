package controllers

import (
	"service-scim/services/auth"
	"sync"
)

type AbstractAuthController struct {
	AbstractInstanceController

	authFactory     *auth.Factory
	authFactoryOnce sync.Once
}

func (c *AbstractAuthController) Prepare() {
	c.AbstractInstanceController.Prepare()

	c.authFactoryOnce.Do(func() {
		c.authFactory = new(auth.Factory)
	})
	if engine, err := c.authFactory.Engine(c.scimConfig); err != nil {
		c.ShowError("Can not get auth engine", 500, "app.auth.invalid-engine", err.Error(), false)
		c.Finish()
		c.StopRun()
	} else if engine != nil {
		if aErr := engine.Auth(c.Ctx); aErr != nil {
			c.ShowError("Unauthorized", 401, "app.unauthorized", aErr.Error(), false)
			c.Finish()
			c.StopRun()
		}
	}
}
