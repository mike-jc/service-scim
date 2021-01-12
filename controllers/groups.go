package controllers

import (
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"net/http"
	"service-scim/errors/repositories"
	"service-scim/models/resources"
	"service-scim/services/repositories"
	repositoriesGroup "service-scim/services/repositories/group"
	"service-scim/services/resources"
	"strconv"
)

type GroupsController struct {
	AbstractResourceController
}

// @Title Get group list
// @router /groups [get]
func (c *GroupsController) GetList() {
	if offset, limit, err := c.PaginationParameters(); err != nil {
		c.ShowError("Cannot parse pagination parameters", 400, "app.groups.list.params", err.Error(), true)
	} else {
		groupService := resources.NewGroupService(c.getGroupRepository(), c.Ctx, c.scimConfig, c.domain, c.Format())
		if list, err := groupService.List(offset, limit); err != nil {
			c.ShowError("Cannot get group list", 500, "app.groups.list.error", err.Error(), true)
		} else {
			c.ServeResponse(list)
		}
	}
}

// @Title Get group by its id
// @router /groups [get]
func (c *GroupsController) GetById() {
	id := c.Ctx.Input.Param(":id")

	if len(id) == 0 {
		c.ShowError("Not found", 404, "app.groups.by_id.not_found", "Empty ID", true)
	} else {
		groupService := resources.NewGroupService(c.getGroupRepository(), c.Ctx, c.scimConfig, c.domain, c.Format())
		if group, err := groupService.ById(id); err != nil {
			if rErr, ok := err.(errorsRepositories.Interface); ok && rErr.Code() == errorsRepositories.NotFoundError {
				c.ShowError("Not found", 404, "app.groups.by_id.not_found", err.Error(), false)
			} else {
				c.ShowError("Cannot get group by ID", 500, "app.groups.by_id.error", err.Error(), true)
			}
		} else {
			c.ServeResponse(group)
		}
	}
}

// @Title Create a user group
// @router /groups [post]
func (c *GroupsController) Create() {
	if data, format, sErr := c.UnmarshalRequestBodyToMap(); sErr != nil {
		c.ShowError("Wrong format of the request body", 400, "app.groups.create", sErr.Error(), true)
	} else {
		groupService := resources.NewGroupService(c.getGroupRepository(), c.Ctx, c.scimConfig, c.domain, format)
		if group, cErr := groupService.Create(data); cErr != nil {
			c.ShowError("Cannot create group", 400, "app.groups.create", cErr.Error(), true)
		} else {
			LogMain.Audit(logger.NewAudit("Group created").
				SetInstance(c.domain).
				SetObject("group", strconv.FormatInt(group.Id, 10)).
				SetObjectAfter(group).
				AddData("requestBodyMap", data))
			c.AddResourceLocationHeader(group)
			c.ServeResponseWithStatus(group, http.StatusCreated)
		}
	}
}

// @Title Modify existing group
// @Description Only attributes with not empty value will be changed
// @router /groups/:id [patch]
func (c *GroupsController) Modify() {
	id := c.Ctx.Input.Param(":id")

	if len(id) == 0 {
		c.ShowError("Not found", 404, "app.groups.modify.not_found", "Empty ID", true)
	} else {
		modification := &modelsResources.Modification{}
		if _, format, sErr := c.UnmarshalRequestBody(modification); sErr != nil {
			c.ShowError("Wrong format of the request body", 400, "app.groups.modify.request", sErr.Error(), true)
		} else if mErr := modification.Validate(); mErr != nil {
			c.ShowError("Bad request", 400, "app.groups.modify.request", mErr.Error(), true)
		} else {
			groupService := resources.NewGroupService(c.getGroupRepository(), c.Ctx, c.scimConfig, c.domain, format)

			var groupBefore *modelsResources.Group
			if group, err := groupService.ById(id); err != nil {
				if rErr, ok := err.(errorsRepositories.Interface); ok && rErr.Code() == errorsRepositories.NotFoundError {
					c.ShowError("Not found", 404, "app.groups.modify.not_found", err.Error(), false)
				} else {
					c.ShowError("Cannot modify group", 500, "app.groups.modify.error", err.Error(), true)
				}
				return
			} else {
				groupBefore = group
			}

			if group, mErr := groupService.Modify(id, modification); mErr != nil {
				c.ShowError("Cannot modify group", 500, "app.groups.modify.error", mErr.Error(), true)
			} else {
				LogMain.Audit(logger.NewAudit("Group modified").
					SetInstance(c.domain).
					SetObject("group", strconv.FormatInt(group.Id, 10)).
					SetObjectBefore(groupBefore).
					SetObjectAfter(group).
					AddData("modification", modification))
				c.AddResourceLocationHeader(group)
				c.ServeResponse(group)
			}
		}
	}
}

// @Title Replace existing group
// @Description All attributes will be changed (even if they are not in request body)
// @router /groups/:id [put]
func (c *GroupsController) Replace() {
	id := c.Ctx.Input.Param(":id")

	if len(id) == 0 {
		c.ShowError("Not found", 404, "app.groups.replace.not_found", "Empty ID", true)
	} else {
		if data, format, sErr := c.UnmarshalRequestBodyToMap(); sErr != nil {
			c.ShowError("Wrong format of the request body", 400, "app.groups.replace", sErr.Error(), true)
		} else {
			groupService := resources.NewGroupService(c.getGroupRepository(), c.Ctx, c.scimConfig, c.domain, format)

			var groupBefore *modelsResources.Group
			if group, err := groupService.ById(id); err != nil {
				if rErr, ok := err.(errorsRepositories.Interface); ok && rErr.Code() == errorsRepositories.NotFoundError {
					c.ShowError("Not found", 404, "app.groups.replace.not_found", err.Error(), false)
				} else {
					c.ShowError("Cannot replace group", 500, "app.groups.replace.error", err.Error(), true)
				}
				return
			} else {
				groupBefore = group
			}

			if group, rpErr := groupService.Replace(id, data); rpErr != nil {
				c.ShowError("Cannot replace group", 500, "app.groups.replace.error", rpErr.Error(), true)
			} else {
				LogMain.Audit(logger.NewAudit("Group replaced").
					SetInstance(c.domain).
					SetObject("group", strconv.FormatInt(group.Id, 10)).
					SetObjectBefore(groupBefore).
					SetObjectAfter(group).
					AddData("requestBodyMap", data))
				c.AddResourceLocationHeader(group)
				c.ServeResponse(group)
			}
		}
	}
}

// @Title Disable existing group
// @router /groups/:id [delete]
func (c *GroupsController) Disable() {
	id := c.Ctx.Input.Param(":id")

	if len(id) == 0 {
		c.ShowError("Not found", 404, "app.groups.disable.not_found", "Empty ID", true)
	} else {
		groupService := resources.NewGroupService(c.getGroupRepository(), c.Ctx, c.scimConfig, c.domain, c.Format())
		if bErr := groupService.Disable(id); bErr != nil {
			if rErr, ok := bErr.(errorsRepositories.Interface); ok && rErr.Code() == errorsRepositories.NotFoundError {
				c.ShowError("Not found", 404, "app.groups.disable.not_found", bErr.Error(), false)
			} else {
				c.ShowError("Cannot disable group", 500, "app.groups.disable.error", bErr.Error(), true)
			}
		} else {
			LogMain.Audit(logger.NewAudit("Group disabled").
				SetInstance(c.domain).
				SetObject("group", id))
			c.SuccessResponseWithStatus(http.StatusNoContent)
		}
	}
}

func (c *GroupsController) getGroupRepository() repositoriesGroup.Interface {
	repo, err := new(repositories.Factory).NewGroupEngine()
	if err != nil {
		LogMain.Log(logger.CreateAlert("Cannot create repository engine for groups collection: " + err.Error()).SetCode("conf.group_repo.create"))
		c.ShowError("Cannot create group repository", 500, "conf.group_repo.create", err.Error(), false)
		c.Abort("500")
	}

	return repo
}
