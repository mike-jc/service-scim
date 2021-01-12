package sdkFile

import (
	"encoding/json"
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"io/ioutil"
	"os"
	"reflect"
	"service-scim/errors/repositories"
	"service-scim/sdks"
	"service-scim/system"
	"strings"
)

type Json struct {
}

func (j *Json) ObjectsFromDir(dirName string, object interface{}) (list map[string]interface{}, err errorsRepositories.Interface) {
	if files, dErr := ioutil.ReadDir(dirName); dErr != nil {
		if _, sErr := os.Stat(dirName); os.IsNotExist(sErr) {
			if cErr := os.MkdirAll(dirName, 0777); cErr != nil {
				sdks.LogMain.Log(logger.CreateError("Can not create directory: " + cErr.Error()))
				return nil, errorsRepositories.NewError(cErr.Error(), errorsRepositories.ApiError)
			} else {
				return nil, nil
			}
		} else {
			sdks.LogMain.Log(logger.CreateError("Can not read directory: " + dErr.Error()))
			return nil, errorsRepositories.NewError(dErr.Error(), errorsRepositories.ApiError)
		}
	} else {
		list = make(map[string]interface{})
		objectType := system.ReflectValue(object).Type()

		for _, f := range files {
			if f.IsDir() {
				continue
			}

			newObjectPtr := reflect.New(objectType).Interface()
			filePath := strings.TrimRight(dirName, "/") + "/" + f.Name()

			if jsonStr, fErr := ioutil.ReadFile(filePath); fErr != nil {
				sdks.LogMain.Log(logger.CreateError("Can not read file " + filePath + ": " + fErr.Error()))
				return nil, errorsRepositories.NewError(fErr.Error(), errorsRepositories.ApiError)
			} else if jErr := json.Unmarshal(jsonStr, newObjectPtr); jErr != nil {
				sdks.LogMain.Log(logger.CreateError("Can not unmarshal json: " + jErr.Error()))
				return nil, errorsRepositories.NewError(jErr.Error(), errorsRepositories.ApiError)
			} else {
				list[filePath] = newObjectPtr
			}
		}
		return
	}
}

func (j *Json) ObjectToFile(filePath string, object interface{}) errorsRepositories.Interface {
	if jsonStr, jErr := json.MarshalIndent(object, "", "  "); jErr != nil {
		sdks.LogMain.Log(logger.CreateError("Can not marshal object: " + jErr.Error()))
		return errorsRepositories.NewError(jErr.Error(), errorsRepositories.ApiError)
	} else if wErr := ioutil.WriteFile(filePath, jsonStr, 0666); wErr != nil {
		sdks.LogMain.Log(logger.CreateError("Can not write object: " + wErr.Error()))
		return errorsRepositories.NewError(wErr.Error(), errorsRepositories.ApiError)
	}
	return nil
}

func (j *Json) RemoveObject(filePath string) errorsRepositories.Interface {
	if rErr := os.Remove(filePath); rErr != nil {
		sdks.LogMain.Log(logger.CreateError("Can not remove object: " + rErr.Error()))
		return errorsRepositories.NewError(rErr.Error(), errorsRepositories.ApiError)
	}
	return nil
}
