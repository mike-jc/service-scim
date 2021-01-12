package repositories

import (
	"fmt"
	"github.com/astaxie/beego"
	"service-scim/services/repositories/group"
	"service-scim/services/repositories/management"
	"service-scim/services/repositories/user"
)

type Factory struct {
}

func (f *Factory) NewUserEngine() (e repositoriesUser.Interface, err error) {
	repositoryEngine := beego.AppConfig.String("repository.engine")
	repositoryUrl := beego.AppConfig.String("repository.api.url")

	switch repositoryEngine {
	case "fake":
		return new(repositoriesUser.Fake), nil
	case "file":
		engine := new(repositoriesUser.File)
		if err := engine.Init(repositoryUrl); err != nil {
			return nil, err
		}
		return engine, nil
	case "keeper":
		engine := new(repositoriesUser.Keeper)
		if err := engine.Init(repositoryUrl); err != nil {
			return nil, err
		}
		return engine, nil
	default:
		return nil, fmt.Errorf("Unknown repository engine %s for users", repositoryEngine)
	}
}

func (f *Factory) NewUserEngineOrPanic() repositoriesUser.Interface {
	engine, err := f.NewUserEngine()
	if err != nil {
		panic("Cannot create user engine: " + err.Error())
	}

	return engine
}

func (f *Factory) NewGroupEngine() (e repositoriesGroup.Interface, err error) {
	repositoryEngine := beego.AppConfig.String("repository.engine")
	repositoryUrl := beego.AppConfig.String("repository.api.url")

	switch repositoryEngine {
	case "fake":
		return new(repositoriesGroup.Fake), nil
	case "file":
		engine := new(repositoriesGroup.File)
		if err := engine.Init(repositoryUrl); err != nil {
			return nil, err
		}
		return engine, nil
	case "keeper":
		engine := new(repositoriesGroup.Keeper)
		if err := engine.Init(repositoryUrl); err != nil {
			return nil, err
		}
		return engine, nil
	default:
		return nil, fmt.Errorf("Unknown repository engine %s for groups", repositoryEngine)
	}
}

func (f *Factory) NewGroupEngineOrPanic() repositoriesGroup.Interface {
	engine, err := f.NewGroupEngine()
	if err != nil {
		panic("Cannot create group engine: " + err.Error())
	}

	return engine
}

func (f *Factory) NewManagementEngine() (e repositoriesManagement.Interface, err error) {
	repositoryEngine := beego.AppConfig.String("repository.engine")
	repositoryUrl := beego.AppConfig.String("repository.management.url")

	switch repositoryEngine {
	case "fake":
		return new(repositoriesManagement.Fake), nil
	case "file":
		engine := new(repositoriesManagement.File)
		if err := engine.Init(repositoryUrl); err != nil {
			return nil, err
		}
		return engine, nil
	case "keeper":
		engine := new(repositoriesManagement.Keeper)
		if err := engine.Init(repositoryUrl); err != nil {
			return nil, err
		}
		return engine, nil
	default:
		return nil, fmt.Errorf("Unknown repository engine %s for management", repositoryEngine)
	}
}

func (f *Factory) NewManagementEngineOrPanic() repositoriesManagement.Interface {
	engine, err := f.NewManagementEngine()
	if err != nil {
		panic("Cannot create management engine: " + err.Error())
	}

	return engine
}
