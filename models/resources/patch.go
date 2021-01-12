package modelsResources

import (
	"service-scim/errors"
	"strings"
)

const PatchOpUrn = "urn:ietf:params:scim:api:messages:2.0:PatchOp"

type Modification struct {
	Schemas    []string `json:"schemas" xml:"Schemas"`
	Operations []*Patch `json:"operations" xml:"Operations"`
}

type Patch struct {
	Operation string      `json:"op" xml:"Op"`
	Path      string      `json:"path" xml:"Path"`
	Value     interface{} `json:"value" xml:"Value"`
}

func (m *Modification) OperationIndexForField(path string) int {
	for i, op := range m.Operations {
		if op.Path == path {
			return i
		}
	}
	return -1
}

func (m *Modification) RemoveOperationByIndex(i int) {
	m.Operations = append(m.Operations[:i], m.Operations[i+1:]...)
}

func (m *Modification) RemoveOperationByPath(path string) {
	if i := m.OperationIndexForField(path); i >= 0 {
		m.Operations = append(m.Operations[:i], m.Operations[i+1:]...)
	}
}

func (m *Modification) Validate() error {
	if len(m.Schemas) != 1 && m.Schemas[0] != PatchOpUrn {
		return scimErrors.InvalidParameterError("schemas", PatchOpUrn, m.Schemas)
	}
	if len(m.Operations) == 0 {
		return scimErrors.InvalidParameterError("operations", "at least one patch operation", "none")
	}

	for _, patch := range m.Operations {
		op := strings.ToLower(patch.Operation)
		switch op {
		case "add":
			if patch.Value == nil {
				return scimErrors.InvalidParameterError("value of add op", "to be present", "nil")
			} else if len(patch.Path) == 0 {
				if _, ok := patch.Value.(map[string]interface{}); !ok {
					return scimErrors.InvalidParameterError("value of add op", "to be complex (for implicit path)", "non-complex")
				}
			}
		case "replace":
			if patch.Value == nil {
				return scimErrors.InvalidParameterError("value of replace op", "to be present", "nil")
			} else if len(patch.Path) == 0 {
				if _, ok := patch.Value.(map[string]interface{}); !ok {
					return scimErrors.InvalidParameterError("value of replace op", "to be complex (for implicit path)", "non-complex")
				}
			}
		case "remove":
			if len(patch.Path) == 0 {
				return scimErrors.InvalidParameterError("path", "to be present", "empty")
			}

		default:
			return scimErrors.InvalidParameterError("op", "one of [add|remove|replace]", patch.Operation)
		}
	}

	return nil
}
