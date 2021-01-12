package controllers

import (
	"github.com/astaxie/beego"
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"service-scim/system"
)

var LogMain *logger.Logger

func RegisterRoutes() {
	prefix := "/" + system.AppVersion()

	indexController := new(IndexController)
	beego.Router("/", indexController, "get:Home")

	healthCheckController := new(HealthCheckController)
	beego.Router("/healthcheck", healthCheckController, "get:HealthCheck")
	beego.Router(prefix+"/healthcheck", healthCheckController, "get:HealthCheck")

	serviceController := new(ServiceController)
	beego.Router(prefix+"/ServiceProviderConfig", serviceController, "get:SpConfiguration")
	beego.Router(prefix+"/serviceproviderconfig", serviceController, "get:SpConfiguration")
	beego.Router(prefix+"/ResourceTypes", serviceController, "get:ResourceTypes")
	beego.Router(prefix+"/ResourceTypes/:name", serviceController, "get:ResourceTypeByName")
	beego.Router(prefix+"/resourcetypes", serviceController, "get:ResourceTypes")
	beego.Router(prefix+"/resourcetypes/:name", serviceController, "get:ResourceTypeByName")
	beego.Router(prefix+"/Schemas", serviceController, "get:Schemas")
	beego.Router(prefix+"/Schemas/:urn", serviceController, "get:SchemaByUrn")
	beego.Router(prefix+"/schemas", serviceController, "get:Schemas")
	beego.Router(prefix+"/schemas/:urn", serviceController, "get:SchemaByUrn")

	usersController := new(UsersController)
	beego.Router(prefix+"/users", usersController, "get:GetList")
	beego.Router(prefix+"/users/:id", usersController, "get:GetById")
	beego.Router(prefix+"/users", usersController, "post:Create")
	beego.Router(prefix+"/users/:id", usersController, "patch:Modify")
	beego.Router(prefix+"/users/:id", usersController, "put:Replace")
	beego.Router(prefix+"/users/:id", usersController, "delete:Block")
	beego.Router(prefix+"/Users", usersController, "get:GetList")
	beego.Router(prefix+"/Users/:id", usersController, "get:GetById")
	beego.Router(prefix+"/Users", usersController, "post:Create")
	beego.Router(prefix+"/Users/:id", usersController, "patch:Modify")
	beego.Router(prefix+"/Users/:id", usersController, "put:Replace")
	beego.Router(prefix+"/Users/:id", usersController, "delete:Block")

	groupsController := new(GroupsController)
	beego.Router(prefix+"/groups", groupsController, "get:GetList")
	beego.Router(prefix+"/groups/:id", groupsController, "get:GetById")
	beego.Router(prefix+"/groups", groupsController, "post:Create")
	beego.Router(prefix+"/groups/:id", groupsController, "patch:Modify")
	beego.Router(prefix+"/groups/:id", groupsController, "put:Replace")
	beego.Router(prefix+"/groups/:id", groupsController, "delete:Disable")
	beego.Router(prefix+"/Groups", groupsController, "get:GetList")
	beego.Router(prefix+"/Groups/:id", groupsController, "get:GetById")
	beego.Router(prefix+"/Groups", groupsController, "post:Create")
	beego.Router(prefix+"/Groups/:id", groupsController, "patch:Modify")
	beego.Router(prefix+"/Groups/:id", groupsController, "put:Replace")
	beego.Router(prefix+"/Groups/:id", groupsController, "delete:Disable")
}
