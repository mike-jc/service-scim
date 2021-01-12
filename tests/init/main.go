package testsInit

import (
	"github.com/astaxie/beego"
	"service-scim/system"
)

func init() {
	system.SetAppDirToCurrentDir(1)
	beego.TestBeegoInit(system.AppDir())
}
