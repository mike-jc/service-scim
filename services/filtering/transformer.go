package filtering

import (
	"fmt"
	"service-scim/errors"
	"service-scim/models/config"
	"service-scim/services/navigation"
)

func StringToFieldsMap(filterStr string, schema *modelsConfig.Schema) (map[string]interface{}, error) {
	if len(filterStr) == 0 {
		return nil, nil
	} else if filter, err := navigation.NewFilter(filterStr); err != nil {
		return nil, err
	} else {
		var attr *modelsConfig.Attribute
		guide := schema.ToAttribute()

		switch filter.Type() {
		case navigation.RelationalOperator:
			path := filter.Left().Data().(navigation.Path)
			attr = guide.GetAttribute(path, true)

			if attr == nil {
				return nil, scimErrors.NoAttributeError(path.CollectValue())
			} else if attr.ExpectsComplex() && filter.Data() != navigation.Pr {
				return nil, scimErrors.InvalidFilterError("", fmt.Sprintf("Cannot perform %v on complex attribute", filter.Data()))
			}

			switch filter.Data() {
			case navigation.Ge, navigation.Gt, navigation.Le, navigation.Lt:
				if attr.ExpectsBool() || attr.ExpectsBinary() {
					return nil, scimErrors.InvalidFilterError("", fmt.Sprintf("Cannot determine order on %s attribute", attr.Type))
				}
			}
		}

		filterMap := make(map[string]interface{})
		switch filter.Data() {
		case navigation.Eq:
			if attr != nil {
				filterMap[attr.Navigation.Path] = filter.Right().Data()
			}
		default:
			return nil, scimErrors.InvalidFilterError("", fmt.Sprintf("unknown operator %v", filter.Data()))
		}
		return filterMap, nil
	}
}
