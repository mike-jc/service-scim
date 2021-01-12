package main

import (
	"github.com/astaxie/beego"
	"service-scim/controllers"
	"service-scim/resources"
	"service-scim/sdks"
	"service-scim/services"
	"service-scim/system"
)

func init() {
	system.SetAppDirToCurrentDir(0)

	// logger
	lg := resources.InitLogger(resources.AppTypeRegular)
	controllers.LogMain = lg
	services.LogMain = lg
	sdks.LogMain = lg

	// routing
	controllers.RegisterRoutes()
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.Run()
}
