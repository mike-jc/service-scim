package resources

import (
	"fmt"
	"reflect"
	"service-scim/errors"
	"service-scim/models/config"
	"service-scim/models/resources"
	"service-scim/services/filtering"
	"service-scim/services/navigation"
	"service-scim/system"
	"service-scim/system/mergemap"
	"strconv"
	"strings"
)

type patchState struct {
	patch    *modelsResources.Patch
	schema   *modelsConfig.Schema
	attr     *modelsConfig.Attribute
	destAttr *modelsConfig.Attribute
}

func (ps *patchState) applyPatchAdd(path navigation.Path, val reflect.Value, dataComplex filtering.Complex) error {
	if path == nil {
		if val.Kind() != reflect.Map {
			return scimErrors.InvalidParameterError("value of add op", "to be complex (for implicit path)", "non-complex")
		}
		for _, k := range val.MapKeys() {
			v0 := val.MapIndex(k)
			patch0 := &modelsResources.Patch{
				Operation: "add",
				Path:      k.String(),
				Value:     v0.Interface(),
			}
			if err := ApplyPatch(patch0, map[string]interface{}(dataComplex), ps.schema); err != nil {
				return err
			}
		}
	} else {
		basePath, lastPath := path.SeparateAtLast()
		baseChannel := make(chan interface{}, 1)

		if basePath == nil {
			go func() {
				baseChannel <- dataComplex
				close(baseChannel)
			}()
		} else {
			baseChannel = dataComplex.Get(basePath, ps.attr)
		}

		for base := range baseChannel {
			baseVal := reflect.ValueOf(base)
			if baseVal.IsNil() {
				continue
			}
			if baseVal.Kind() == reflect.Interface {
				baseVal = baseVal.Elem()
			}

			switch baseVal.Kind() {
			case reflect.Map:
				keyVal := reflect.ValueOf(lastPath.Base())
				if ps.destAttr.MultiValued {
					origVal := baseVal.MapIndex(keyVal)
					if !origVal.IsValid() {
						switch val.Kind() {
						case reflect.Array, reflect.Slice:
							baseVal.SetMapIndex(keyVal, val)
						default:
							baseVal.SetMapIndex(keyVal, reflect.ValueOf([]interface{}{val.Interface()}))
						}

					} else {
						if origVal.Kind() == reflect.Interface {
							origVal = origVal.Elem()
						}
						newArr := filtering.MultiValued(origVal.Interface().([]interface{}))
						switch val.Kind() {
						case reflect.Array, reflect.Slice:
							for i := 0; i < val.Len(); i++ {
								newArr = newArr.Add(val.Index(i).Interface())
							}
						default:
							newArr = newArr.Add(val.Interface())
						}
						baseVal.SetMapIndex(keyVal, reflect.ValueOf([]interface{}(newArr)))
					}
				} else {
					baseVal.SetMapIndex(keyVal, val)
				}
			case reflect.Array, reflect.Slice:
				for i := 0; i < baseVal.Len(); i++ {
					elemVal := baseVal.Index(i)
					if elemVal.Kind() == reflect.Interface {
						elemVal = elemVal.Elem()
					}
					switch elemVal.Kind() {
					case reflect.Map:
						elemVal.SetMapIndex(reflect.ValueOf(lastPath.Base()), val)
					default:
						return scimErrors.InvalidPathError(ps.patch.Path, "array base contains non-map")
					}
				}
			default:
				return scimErrors.InvalidPathError(ps.patch.Path, "base evaluated to non-map and non-array.")
			}
		}
	}
	return nil
}

func (ps *patchState) applyPatchReplace(path navigation.Path, val reflect.Value, dataComplex filtering.Complex) error {
	if path == nil {
		if val.Kind() != reflect.Map {
			return scimErrors.InvalidParameterError("value of replace op", "to be complex (for implicit path)", "non-complex")
		}
		for _, k := range val.MapKeys() {
			v0 := val.MapIndex(k)
			patch0 := &modelsResources.Patch{
				Operation: "replace",
				Path:      k.String(),
				Value:     v0.Interface(),
			}
			if err := ApplyPatch(patch0, map[string]interface{}(dataComplex), ps.schema); err != nil {
				return err
			}
		}
	} else {
		basePath, lastPath := path.SeparateAtLast()
		baseChannel := make(chan interface{}, 1)
		if basePath == nil {
			go func() {
				baseChannel <- dataComplex
				close(baseChannel)
			}()
		} else {
			baseChannel = dataComplex.Get(basePath, ps.attr)
		}

		for base := range baseChannel {
			baseVal := reflect.ValueOf(base)
			if baseVal.IsNil() {
				continue
			}
			if baseVal.Kind() == reflect.Interface {
				baseVal = baseVal.Elem()
			}
			switch baseVal.Kind() {
			case reflect.Map:
				switch val.Kind() {
				case reflect.Map:
					currentEl := baseVal.MapIndex(reflect.ValueOf(lastPath.Base()))
					// check if the the current value is a map (complex)
					// we have to update only the fields which were passed and not overwrite the other ones
					// So we merge new value with existed value
					if currentEl.IsValid() && !currentEl.IsNil() {
						currentVal, currentValOk := currentEl.Interface().(map[string]interface{})
						newVal, newValOk := val.Interface().(map[string]interface{})
						if !currentValOk {
							return scimErrors.InvalidParameterError("value of replace op", "to be non-complex", "complex")
						} else if !newValOk {
							return scimErrors.InvalidParameterError("value of replace op", "to be complex", "non-complex")
						}

						mergemap.Merge(currentVal, newVal)
						baseVal.SetMapIndex(reflect.ValueOf(lastPath.Base()), reflect.ValueOf(currentVal))
					} else {
						baseVal.SetMapIndex(reflect.ValueOf(lastPath.Base()), val)
					}
				default:
					baseVal.SetMapIndex(reflect.ValueOf(lastPath.Base()), val)
				}
			case reflect.Array, reflect.Slice:
				if index, iErr := strconv.Atoi(lastPath.Base()); iErr == nil && index >= 0 && index < baseVal.Len() {
					baseVal.Index(index).Set(val)
				}
			}
		}
	}
	return nil
}

func (ps *patchState) applyPatchRemove(path navigation.Path, val reflect.Value, dataComplex filtering.Complex) error {
	basePath, lastPath := path.SeparateAtLast()
	baseChannel := make(chan interface{}, 1)
	if basePath == nil {
		go func() {
			baseChannel <- dataComplex
			close(baseChannel)
		}()
	} else {
		baseChannel = dataComplex.Get(basePath, ps.attr)
	}

	var baseAttr = ps.attr
	if basePath != nil {
		baseAttr = ps.attr.GetAttribute(basePath, true)
	}

	var filterMapFunc = func(baseVal, keyVal reflect.Value, filter navigation.FilterNode) {
		origVal := baseVal.MapIndex(keyVal)
		baseAttr = baseAttr.GetAttribute(lastPath, false)
		reverseRoot := navigation.NewFilterNode().
			SetData(navigation.Not).
			SetType(navigation.LogicalOperator).
			SetLeft(filter)
		newElemChannel := filtering.MultiValued(origVal.Interface().([]interface{})).
			Filter(reverseRoot, baseAttr)
		newArr := make([]interface{}, 0)
		for newElem := range newElemChannel {
			newArr = append(newArr, newElem)
		}
		if len(newArr) == 0 {
			baseVal.SetMapIndex(keyVal, reflect.Value{})
		} else {
			baseVal.SetMapIndex(keyVal, reflect.ValueOf(newArr))
		}
	}

	for base := range baseChannel {
		baseVal := reflect.ValueOf(base)
		if baseVal.IsNil() {
			continue
		}
		if baseVal.Kind() == reflect.Interface {
			baseVal = baseVal.Elem()
		}

		switch baseVal.Kind() {
		case reflect.Map:
			keyVal := reflect.ValueOf(lastPath.Base())
			if ps.destAttr.MultiValued {
				if lastPath.FilterRoot() == nil {
					if eqFilter, err := navigation.NewEqFilterNodeFromValue(val); err != nil {
						return err
					} else if eqFilter != nil {
						filterMapFunc(baseVal, keyVal, eqFilter)
					} else {
						baseVal.SetMapIndex(keyVal, reflect.Value{})
					}
				} else {
					filterMapFunc(baseVal, keyVal, lastPath.FilterRoot())
				}
			} else {
				baseVal.SetMapIndex(keyVal, reflect.Value{})
			}
		case reflect.Array, reflect.Slice:
			keyVal := reflect.ValueOf(lastPath.Base())
			for i := 0; i < baseVal.Len(); i++ {
				elemVal := baseVal.Index(i)
				if elemVal.Kind() == reflect.Interface {
					elemVal = elemVal.Elem()
				}
				switch elemVal.Kind() {
				case reflect.Map:
					elemVal.SetMapIndex(keyVal, reflect.Value{})
				default:
					return scimErrors.InvalidPathError(ps.patch.Path, "array base contains non-map")
				}
			}
		default:
			return scimErrors.InvalidPathError(ps.patch.Path, "base evaluated to non-map and non-array.")
		}
	}
	return nil
}

func ApplyPatch(patch *modelsResources.Patch, data interface{}, schema *modelsConfig.Schema) error {
	ps := patchState{
		patch:  patch,
		schema: schema,
		attr:   schema.ToAttribute(),
	}

	var path navigation.Path
	var err error

	if len(patch.Path) == 0 {
		path = nil
	} else {
		if path, err = navigation.NewPath(patch.Path); err != nil {
			return err
		}
		if ps.attr != nil {
			ps.attr.CorrectPathCase(path, true)

			if subAttr := ps.attr.GetAttribute(path, true); subAttr != nil {
				ps.destAttr = subAttr
			} else {
				return scimErrors.InvalidPathError(patch.Path, "no attribute found for path")
			}
		} else {
			return scimErrors.InvalidPathError(patch.Path, "no attribute found for path")
		}
	}

	if dataMap, ok := data.(map[string]interface{}); !ok {
		return fmt.Errorf("Entity, which the patch is applied to, should be a map")
	} else {
		dataComplex := filtering.Complex(dataMap)
		val := system.ReflectValue(patch.Value)
		op := strings.ToLower(patch.Operation)

		switch op {
		case "add":
			if aErr := ps.applyPatchAdd(path, val, dataComplex); aErr != nil {
				return aErr
			}
		case "replace":
			if aErr := ps.applyPatchReplace(path, val, dataComplex); aErr != nil {
				return aErr
			}
		case "remove":
			if aErr := ps.applyPatchRemove(path, val, dataComplex); aErr != nil {
				return aErr
			}
		default:
			err = scimErrors.InvalidParameterError("op", "one of [add|remove|replace]", patch.Operation)
		}

		data = map[string]interface{}(dataComplex)
	}

	return nil
}
