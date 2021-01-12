package repositoriesUser

import (
	"fmt"
	"service-scim/errors/repositories"
	"service-scim/models/resources"
	"service-scim/sdks/file"
	"service-scim/system"
	"strings"
)

type File struct {
	Abstract

	jsonReader *sdkFile.Json
	jsonDir    string
}

func (f *File) Init(dir string) errorsRepositories.Interface {
	f.jsonReader = new(sdkFile.Json)
	f.jsonDir = strings.TrimRight(dir, "/")
	return nil
}

func (f *File) List(offset, limit int, filterMap map[string]interface{}) (totalCount int, list []*modelsResources.User, err errorsRepositories.Interface) {
	if list, err = f.readUsers(); err != nil {
		list = nil
		return
	} else {
		totalCount = len(list)
		if offset < totalCount {
			if offset+limit < totalCount {
				list = list[offset : offset+limit]
			} else {
				list = list[offset:]
			}
		} else {
			list = nil
		}
		return
	}
}

func (f *File) ById(id string) (user *modelsResources.User, err errorsRepositories.Interface) {
	if users, rErr := f.readUsers(); rErr != nil {
		return nil, rErr
	} else {
		for _, u := range users {
			if u.ScimId() == id {
				return u, nil
			}
		}
		return nil, errorsRepositories.NewError("User not found by ID "+id, errorsRepositories.NotFoundError)
	}
}

func (f *File) Create(data map[string]interface{}) (resultedUser *modelsResources.User, err errorsRepositories.Interface) {
	user := modelsResources.NewUserFromMap(data, nil)
	if wErr := f.writeUser(user); wErr != nil {
		return nil, wErr
	}
	return user, nil
}

func (f *File) Update(id string, data map[string]interface{}) (resultedUser *modelsResources.User, err errorsRepositories.Interface) {
	if users, rErr := f.readUsersAndFilePaths(); rErr != nil {
		return nil, rErr
	} else {
		for filePath, u := range users {
			if u.ScimId() == id {
				newUser := modelsResources.NewUserFromMap(data, u)
				if wErr := f.jsonReader.ObjectToFile(filePath, newUser); wErr != nil {
					return nil, wErr
				}
				return newUser, nil
			}
		}
		return nil, errorsRepositories.NewError("User not found by ID "+id, errorsRepositories.NotFoundError)
	}
}

func (f *File) Replace(id string, data map[string]interface{}) (resultedUser *modelsResources.User, err errorsRepositories.Interface) {
	if users, rErr := f.readUsersAndFilePaths(); rErr != nil {
		return nil, rErr
	} else {
		for filePath, u := range users {
			if u.ScimId() == id {
				newUser := modelsResources.NewUserFromMap(data, nil)
				if wErr := f.jsonReader.ObjectToFile(filePath, newUser); wErr != nil {
					return nil, wErr
				}
				return newUser, nil
			}
		}
		return nil, errorsRepositories.NewError("User not found by ID "+id, errorsRepositories.NotFoundError)
	}
}

func (f *File) Block(id string) errorsRepositories.Interface {
	if users, rErr := f.readUsersAndFilePaths(); rErr != nil {
		return rErr
	} else {
		for filePath, u := range users {
			if u.ScimId() == id {
				if rmErr := f.jsonReader.RemoveObject(filePath); rmErr != nil {
					return rmErr
				}
				return nil
			}
		}
		return errorsRepositories.NewError("User not found by ID "+id, errorsRepositories.NotFoundError)
	}
}

func (f *File) Count(filter map[string]interface{}, id *string) (count int, err errorsRepositories.Interface) {
	if users, rErr := f.readUsers(); rErr != nil {
		return 0, rErr
	} else {
		count := 0
		for _, u := range users {
			if id == nil || u.ScimId() != *id {
				if system.StructFilterPassed(u, filter) {
					count++
				}
			}
		}
		return count, nil
	}
}

func (f *File) readUsers() ([]*modelsResources.User, errorsRepositories.Interface) {
	if objects, rErr := f.jsonReader.ObjectsFromDir(f.jsonDir, new(modelsResources.User)); rErr != nil {
		return nil, rErr
	} else {
		list := make([]*modelsResources.User, 0)
		for _, object := range objects {
			if user, ok := object.(*modelsResources.User); !ok {
				return nil, errorsRepositories.NewError("Files repository. Got object that is not modelsResources.User", errorsRepositories.GeneralError)
			} else {
				list = append(list, user)
			}
		}
		return list, nil
	}
}

func (f *File) readUsersAndFilePaths() (map[string]*modelsResources.User, errorsRepositories.Interface) {
	if objects, rErr := f.jsonReader.ObjectsFromDir(f.jsonDir, new(modelsResources.User)); rErr != nil {
		return nil, rErr
	} else {
		list := make(map[string]*modelsResources.User, 0)
		for filePath, object := range objects {
			if user, ok := object.(*modelsResources.User); !ok {
				return nil, errorsRepositories.NewError("Files repository. Got object that is not modelsResources.User", errorsRepositories.GeneralError)
			} else {
				list[filePath] = user
			}
		}
		return list, nil
	}
}

func (f *File) writeUser(user *modelsResources.User) errorsRepositories.Interface {
	filePath := fmt.Sprintf("%s/%s.json", f.jsonDir, user.ScimId())

	if wErr := f.jsonReader.ObjectToFile(filePath, user); wErr != nil {
		return wErr
	}
	return nil
}
