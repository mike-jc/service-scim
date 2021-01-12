package controllers

import (
	"github.com/astaxie/beego"
	"service-scim/models/resources"
)

type AbstractResourceController struct {
	AbstractAuthController
}

func (c *AbstractResourceController) PaginationParameters() (offset, limit int, err error) {
	if offset, err = c.GetInt("startIndex", 1); err != nil {
		return
	}
	offset -= 1 // startIndex is 1-based, but offset is 0-based
	if offset < 0 {
		offset = 0
	}

	var maxResults int
	if maxResults, err = beego.AppConfig.Int("sp.filter.maxResults"); err != nil {
		return
	}
	if limit, err = c.GetInt("count", maxResults); err != nil {
		return
	}
	return
}

func (c *AbstractResourceController) AddResourceLocationHeader(data interface{}) {
	if c.Ctx.Output.Status >= 300 && c.Ctx.Output.Status < 400 {
		// it's not redirecting location, so do nothing in that case
		return
	}

	var location string

	switch data.(type) {
	case modelsResources.User:
		user := data.(modelsResources.User)
		location = c.BasicURL() + user.Meta.Location
	case *modelsResources.User:
		user := data.(*modelsResources.User)
		location = c.BasicURL() + user.Meta.Location
	case modelsResources.Group:
		group := data.(modelsResources.Group)
		location = c.BasicURL() + group.Meta.Location
	case *modelsResources.Group:
		group := data.(*modelsResources.Group)
		location = c.BasicURL() + group.Meta.Location
	default:
		location = ""
	}

	if len(location) > 0 {
		c.Ctx.Output.Header("Location", location)
	}
}
