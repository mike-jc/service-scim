package controllers

import (
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"net/http"
	"service-scim/errors/repositories"
	"service-scim/models/resources"
	"service-scim/services/repositories"
	repositoriesUser "service-scim/services/repositories/user"
	"service-scim/services/resources"
	"strconv"
)

type UsersController struct {
	AbstractResourceController
}

// @Title Get user list
// @router /users [get]
func (c *UsersController) GetList() {
	if offset, limit, err := c.PaginationParameters(); err != nil {
		c.ShowError("Cannot parse pagination parameters", 400, "app.users.list.params", err.Error(), true)
	} else {
		userService := resources.NewUserService(c.getUserRepository(), c.Ctx, c.scimConfig, c.domain, c.Format())
		filter := c.GetString("filter", "")
		if list, err := userService.List(offset, limit, filter); err != nil {
			c.ShowError("Cannot get user list", 500, "app.users.list.error", err.Error(), true)
		} else {
			c.ServeResponse(list)
		}
	}
}

// @Title Get user by his id
// @router /users/:id [get]
func (c *UsersController) GetById() {
	id := c.Ctx.Input.Param(":id")

	if len(id) == 0 {
		c.ShowError("Not found", 404, "app.users.by_id.not_found", "Empty ID", true)
	} else {
		userService := resources.NewUserService(c.getUserRepository(), c.Ctx, c.scimConfig, c.domain, c.Format())
		if user, err := userService.ById(id); err != nil {
			if rErr, ok := err.(errorsRepositories.Interface); ok && rErr.Code() == errorsRepositories.NotFoundError {
				c.ShowError("Not found", 404, "app.users.by_id.not_found", err.Error(), false)
			} else {
				c.ShowError("Cannot get user by ID", 500, "app.users.by_id.error", err.Error(), true)
			}
		} else {
			c.ServeResponse(user)
		}
	}
}

// @Title Create a user
// @router /users [post]
func (c *UsersController) Create() {
	if data, format, sErr := c.UnmarshalRequestBodyToMap(); sErr != nil {
		c.ShowError("Wrong format of the request body", 400, "app.users.create", sErr.Error(), true)
	} else {
		userService := resources.NewUserService(c.getUserRepository(), c.Ctx, c.scimConfig, c.domain, format)
		if user, cErr := userService.Create(data); cErr != nil {
			c.ShowError("Cannot create user", 500, "app.users.create", cErr.Error(), true)
		} else {
			LogMain.Audit(logger.NewAudit("User created").
				SetInstance(c.domain).
				SetObject("user", strconv.FormatInt(user.Id, 10)).
				SetObjectAfter(user).
				AddData("requestBodyMap", data))
			c.AddResourceLocationHeader(user)
			c.ServeResponseWithStatus(user, http.StatusCreated)
		}
	}
}

// @Title Modify existing user
// @Description Only attributes with not empty value will be changed
// @router /users/:id [patch]
func (c *UsersController) Modify() {
	id := c.Ctx.Input.Param(":id")

	if len(id) == 0 {
		c.ShowError("Not found", 404, "app.users.modify.not_found", "Empty ID", true)
	} else {
		modification := &modelsResources.Modification{}
		if _, format, sErr := c.UnmarshalRequestBody(modification); sErr != nil {
			c.ShowError("Wrong format of the request body", 400, "app.users.modify.request", sErr.Error(), true)
		} else if mErr := modification.Validate(); mErr != nil {
			c.ShowError("Bad request", 400, "app.users.modify.request", mErr.Error(), true)
		} else {
			userService := resources.NewUserService(c.getUserRepository(), c.Ctx, c.scimConfig, c.domain, format)

			var userBefore *modelsResources.User
			if user, err := userService.ById(id); err != nil {
				if rErr, ok := err.(errorsRepositories.Interface); ok && rErr.Code() == errorsRepositories.NotFoundError {
					c.ShowError("Not found", 404, "app.users.modify.not_found", err.Error(), false)
				} else {
					c.ShowError("Cannot modify user", 500, "app.users.modify.error", err.Error(), true)
				}
				return
			} else {
				userBefore = user
			}

			if user, mErr := userService.Modify(id, modification); mErr != nil {
				c.ShowError("Cannot modify user", 500, "app.users.modify.error", mErr.Error(), true)
			} else {
				LogMain.Audit(logger.NewAudit("User modified").
					SetInstance(c.domain).
					SetObject("user", strconv.FormatInt(user.Id, 10)).
					SetObjectBefore(userBefore).
					SetObjectAfter(user).
					AddData("modification", modification))
				c.AddResourceLocationHeader(user)
				c.ServeResponse(user)
			}
		}
	}
}

// @Title Replace existing user
// @Description All attributes will be changed (even if they are not in request body)
// @router /users/:id [put]
func (c *UsersController) Replace() {
	id := c.Ctx.Input.Param(":id")

	if len(id) == 0 {
		c.ShowError("Not found", 404, "app.users.replace.not_found", "Empty ID", true)
	} else {
		if data, format, sErr := c.UnmarshalRequestBodyToMap(); sErr != nil {
			c.ShowError("Wrong format of the request body", 400, "app.users.replace", sErr.Error(), true)
		} else {
			userService := resources.NewUserService(c.getUserRepository(), c.Ctx, c.scimConfig, c.domain, format)

			var userBefore *modelsResources.User
			if user, err := userService.ById(id); err != nil {
				if rErr, ok := err.(errorsRepositories.Interface); ok && rErr.Code() == errorsRepositories.NotFoundError {
					c.ShowError("Not found", 404, "app.users.replace.not_found", err.Error(), false)
				} else {
					c.ShowError("Cannot replace user", 500, "app.users.replace.error", err.Error(), true)
				}
				return
			} else {
				userBefore = user
			}

			if user, rpErr := userService.Replace(id, data); rpErr != nil {
				c.ShowError("Cannot replace user", 500, "app.users.replace.error", rpErr.Error(), true)
			} else {
				LogMain.Audit(logger.NewAudit("User replaced").
					SetInstance(c.domain).
					SetObject("user", strconv.FormatInt(user.Id, 10)).
					SetObjectBefore(userBefore).
					SetObjectAfter(user).
					AddData("requestBodyMap", data))
				c.AddResourceLocationHeader(user)
				c.ServeResponse(user)
			}
		}
	}
}

// @Title Block existing user
// @router /users/:id [delete]
func (c *UsersController) Block() {
	id := c.Ctx.Input.Param(":id")

	if len(id) == 0 {
		c.ShowError("Not found", 404, "app.users.block.not_found", "Empty ID", true)
	} else {
		userService := resources.NewUserService(c.getUserRepository(), c.Ctx, c.scimConfig, c.domain, c.Format())

		var userBefore *modelsResources.User
		if user, err := userService.ById(id); err != nil {
			if rErr, ok := err.(errorsRepositories.Interface); ok && rErr.Code() == errorsRepositories.NotFoundError {
				c.ShowError("Not found", 404, "app.users.block.not_found", err.Error(), false)
			} else {
				c.ShowError("Cannot disable user account", 500, "app.users.block.error", err.Error(), true)
			}
			return
		} else {
			userBefore = user
		}

		if bErr := userService.Block(id); bErr != nil {
			if rErr, ok := bErr.(errorsRepositories.Interface); ok && rErr.Code() == errorsRepositories.NotFoundError {
				c.ShowError("Not found", 404, "app.users.block.not_found", bErr.Error(), false)
			} else {
				c.ShowError("Cannot disable user account", 500, "app.users.block.error", bErr.Error(), true)
			}
		} else {
			LogMain.Audit(logger.NewAudit("User blocked").
				SetInstance(c.domain).
				SetObject("user", id).
				SetObjectBefore(userBefore))
			c.SuccessResponseWithStatus(http.StatusNoContent)
		}
	}
}

func (c *UsersController) getUserRepository() repositoriesUser.Interface {
	userRepository, err := new(repositories.Factory).NewUserEngine()
	if err != nil {
		LogMain.Log(logger.CreateAlert("Cannot create repository engine for users collection: " + err.Error()).SetCode("conf.user_repo.create"))
		c.ShowError("Cannot create user repository", 500, "conf.user_repo.create", err.Error(), false)
		c.Abort("500")
	}

	return userRepository
}
