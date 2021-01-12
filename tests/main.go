package tests

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/context"
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	sdksData "gitlab.com/24sessions/sdk-go-configurator/data"
	"io/ioutil"
	"net/http"
	"service-scim/resources"
	_ "service-scim/tests/init"
)

var LogMain *logger.Logger

func Init() {
	logger.OverrideConfigAsDummy(logger.LOG_LEVEL_DEBUG)
	LogMain = resources.InitLogger(resources.AppTypeTest)
}

func NewContext() *context.Context {
	ctx := new(context.Context)
	ctx.Input = new(context.BeegoInput)
	ctx.Input.Context = new(context.Context)
	ctx.Input.Context.Request = new(http.Request)
	ctx.Input.Context.Request.Header = make(map[string][]string)
	return ctx
}

func NewConfig() *sdksData.ScimContainer {
	roleMapping, _ := json.Marshal(map[string][]int64{
		"operator": []int64{123},
	})
	cont := sdksData.InstanceContainer{
		ScimRabobankRoleMapping: string(roleMapping),
	}
	return sdksData.NewScim(&cont)
}

func ReadFromFixture(data interface{}, filePath string) (interface{}, error) {
	if jsonStr, err := ioutil.ReadFile(filePath); err != nil {
		return nil, fmt.Errorf("Can not read fixture '%s': %s", filePath, err.Error())
	} else {
		if err := json.Unmarshal(jsonStr, data); err != nil {
			return nil, fmt.Errorf("Can not unserialize JSON in '%s': %s", filePath, err.Error())
		} else {
			return data, nil
		}
	}
}
