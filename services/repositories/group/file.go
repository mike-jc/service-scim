package repositoriesGroup

import (
	"fmt"
	"service-scim/errors/repositories"
	"service-scim/models/resources"
	"service-scim/sdks/file"
	"service-scim/system"
)

type File struct {
	Abstract

	jsonReader *sdkFile.Json
	jsonDir    string
}

func (f *File) Init(dir string) errorsRepositories.Interface {
	f.jsonReader = new(sdkFile.Json)
	f.jsonDir = dir
	return nil
}

func (f *File) List(offset, limit int) (totalCount int, list []*modelsResources.Group, err errorsRepositories.Interface) {
	if list, err = f.readGroups(); err != nil {
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

func (f *File) ById(id string) (group *modelsResources.Group, err errorsRepositories.Interface) {
	if groups, rErr := f.readGroups(); rErr != nil {
		return nil, rErr
	} else {
		for _, g := range groups {
			if g.ScimId() == id {
				return g, nil
			}
		}
		return nil, errorsRepositories.NewError("Group not found by ID "+id, errorsRepositories.NotFoundError)
	}
}

func (f *File) Create(data map[string]interface{}) (resultedGroup *modelsResources.Group, err errorsRepositories.Interface) {
	group := modelsResources.NewGroupFromMap(data, nil)
	if wErr := f.writeGroup(group); wErr != nil {
		return nil, wErr
	}
	return group, nil
}

func (f *File) Update(id string, data map[string]interface{}) (resultedGroup *modelsResources.Group, err errorsRepositories.Interface) {
	if groups, rErr := f.readGroupsAndFilePaths(); rErr != nil {
		return nil, rErr
	} else {
		for filePath, g := range groups {
			if g.ScimId() == id {
				newGroup := modelsResources.NewGroupFromMap(data, g)
				if wErr := f.jsonReader.ObjectToFile(filePath, newGroup); wErr != nil {
					return nil, wErr
				}
				return newGroup, nil
			}
		}
		return nil, errorsRepositories.NewError("Group not found by ID "+id, errorsRepositories.NotFoundError)
	}
}

func (f *File) Replace(id string, data map[string]interface{}) (resultedGroup *modelsResources.Group, err errorsRepositories.Interface) {
	if groups, rErr := f.readGroupsAndFilePaths(); rErr != nil {
		return nil, rErr
	} else {
		for filePath, g := range groups {
			if g.ScimId() == id {
				newGroup := modelsResources.NewGroupFromMap(data, nil)
				if wErr := f.jsonReader.ObjectToFile(filePath, newGroup); wErr != nil {
					return nil, wErr
				}
				return newGroup, nil
			}
		}
		return nil, errorsRepositories.NewError("Group not found by ID "+id, errorsRepositories.NotFoundError)
	}
}

func (f *File) Disable(id string) errorsRepositories.Interface {
	if groups, rErr := f.readGroupsAndFilePaths(); rErr != nil {
		return rErr
	} else {
		for filePath, g := range groups {
			if g.ScimId() == id {
				if rmErr := f.jsonReader.RemoveObject(filePath); rmErr != nil {
					return rmErr
				}
				return nil
			}
		}
		return errorsRepositories.NewError("Group not found by ID "+id, errorsRepositories.NotFoundError)
	}
}

func (f *File) Count(filter map[string]interface{}, id *string) (count int, err errorsRepositories.Interface) {
	if groups, rErr := f.readGroups(); rErr != nil {
		return 0, rErr
	} else {
		count := 0
		for _, g := range groups {
			if id == nil || g.ScimId() != *id {
				if system.StructFilterPassed(g, filter) {
					count++
				}
			}
		}
		return count, nil
	}
}

func (f *File) readGroups() ([]*modelsResources.Group, errorsRepositories.Interface) {
	if objects, rErr := f.jsonReader.ObjectsFromDir(f.jsonDir, new(modelsResources.Group)); rErr != nil {
		return nil, rErr
	} else {
		list := make([]*modelsResources.Group, 0)
		for _, object := range objects {
			if group, ok := object.(*modelsResources.Group); !ok {
				return nil, errorsRepositories.NewError("Files repository. Got object that is not modelsResources.Group", errorsRepositories.GeneralError)
			} else {
				list = append(list, group)
			}
		}
		return list, nil
	}
}

func (f *File) readGroupsAndFilePaths() (map[string]*modelsResources.Group, errorsRepositories.Interface) {
	if objects, rErr := f.jsonReader.ObjectsFromDir(f.jsonDir, new(modelsResources.Group)); rErr != nil {
		return nil, rErr
	} else {
		list := make(map[string]*modelsResources.Group, 0)
		for filePath, object := range objects {
			if group, ok := object.(*modelsResources.Group); !ok {
				return nil, errorsRepositories.NewError("Files repository. Got object that is not modelsResources.Group", errorsRepositories.GeneralError)
			} else {
				list[filePath] = group
			}
		}
		return list, nil
	}
}

func (f *File) writeGroup(group *modelsResources.Group) errorsRepositories.Interface {
	filePath := fmt.Sprintf("%s/%s.json", f.jsonDir, group.ScimId())

	if wErr := f.jsonReader.ObjectToFile(filePath, group); wErr != nil {
		return wErr
	}
	return nil
}
