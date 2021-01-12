package controllers

import (
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"service-scim/services/repositories"
	"service-scim/services/repositories/management"
)

type HealthCheckController struct {
	AbstractController

	repository repositoriesManagement.Interface
}

// @Title Health check
// @Description check system
// @router /healthcheck [get]
func (c *HealthCheckController) HealthCheck() {

	// check that repository is alive
	if repository, err := c.Repository(); err != nil {
		LogMain.Log(logger.CreateError("Can not create repository engine for management: " + err.Error()))
	} else {
		if err := repository.Ping(); err != nil {
			c.ShowError("Something wrong with repository", 500, "app.healthcheck.repository", err.Error(), false)
		}
	}

	c.SuccessResponse()
}

func (c *HealthCheckController) Repository() (repository repositoriesManagement.Interface, err error) {
	if c.repository == nil {
		repositoryFactory := new(repositories.Factory)
		if repository, err := repositoryFactory.NewManagementEngine(); err != nil {
			return nil, err
		} else {
			c.repository = repository
		}
	}
	return c.repository, nil
}
