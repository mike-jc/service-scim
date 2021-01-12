package system

import (
	"github.com/astaxie/beego"
	"path/filepath"
	"runtime"
)

var appDir string

func SetAppDirToCurrentDir(skipChildDirectories int) {
	_, callerFile, _, _ := runtime.Caller(1)

	dir := filepath.Dir(callerFile)
	for i := 0; i < skipChildDirectories; i++ {
		dir = filepath.Join(dir, "..")
	}

	var err error
	appDir, err = filepath.Abs(dir)
	if err != nil {
		panic("Can not get working directory's path")
	}
}

func AppDir() string {
	return appDir
}

func AppVersion() string {
	version := beego.AppConfig.String("sp.version")
	if len(version) == 0 {
		return "v2"
	}
	return version
}
