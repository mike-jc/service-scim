package modelsConfig

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
	"service-scim/services/navigation"
	"service-scim/system"
	"strings"
)

type Schema struct {
	Id          string       `json:"id" xml:"Id"`
	Name        string       `json:"name" xml:"Name"`
	Description string       `json:"description,omitempty" xml:"Description,omitempty"`
	Attributes  []*Attribute `json:"attributes,omitempty" xml:"Attributes,omitempty"`
	Meta        *SchemaMeta  `json:"meta,omitempty" xml:"Meta,omitempty"`
	XMLName     struct{}     `json:"-" xml:"Schema"`
}

type Schemas struct {
	Schemas      []string  `json:"schemas" xml:"Schemas"`
	TotalResults int       `json:"totalResults" xml:"totalResults"`
	ItemsPerPage int       `json:"itemsPerPage" xml:"ItemsPerPage"`
	StartIndex   int       `json:"startIndex" xml:"StartIndex"`
	Resources    []*Schema `json:"Resources" xml:"Resources"`
	XMLName      struct{}  `json:"-" xml:"Schemas"`
}

type Attribute struct {
	Name            string               `json:"name" xml:"name,attr"`
	Type            string               `json:"type" xml:"type,attr"`
	MultiValued     bool                 `json:"multiValued" xml:"multiValued,attr"`
	Required        bool                 `json:"required" xml:"required,attr"`
	CaseExact       *bool                `json:"caseExact,omitempty" xml:"caseExact,attr,omitempty"`
	Mutability      string               `json:"mutability,omitempty" xml:"mutability,attr,omitempty"`
	Returned        string               `json:"returned,omitempty" xml:"returned,attr,omitempty"`
	Uniqueness      string               `json:"uniqueness,omitempty" xml:"uniqueness,attr,omitempty"`
	Description     string               `json:"description,omitempty" xml:"Description,omitempty"`
	CanonicalValues []string             `json:"canonicalValues,omitempty" xml:"CanonicalValues,omitempty"`
	ReferenceTypes  []string             `json:"referenceTypes,omitempty" xml:"ReferenceTypes,omitempty"`
	SubAttributes   []*Attribute         `json:"subAttributes,omitempty" xml:"SubAttributes,omitempty"`
	Navigation      *AttributeNavigation `json:"-" xml:"-"`
	XMLName         struct{}             `json:"-" xml:"Attribute"`

	// extension attributes aren't shown in schema representation
	IsExtensionAttribute bool `json:"-" xml:"-"`
}

type AttributeNavigation struct {
	FieldName string
	Path      string
	FullPath  string
	IndexKeys []string
}

type SchemaMeta struct {
	Location     string `json:"location" xml:"Location"`
	ResourceType string `json:"resourceType" xml:"ResourceType"`
}

func (s *Schema) ToAttribute() *Attribute {
	return &Attribute{
		Type:          "complex",
		MultiValued:   false,
		Mutability:    "readWrite",
		SubAttributes: s.Attributes,
		Navigation: &AttributeNavigation{
			IndexKeys: []string{},
		},
	}
}

func (s *Schema) MarshalJSON() ([]byte, error) {
	type Alias Schema

	// skip extension attributes since they're not supposed to be shown in the schema representation,
	// but only used for hydrating resource's data
	return json.Marshal(&struct {
		Attributes []*Attribute `json:"attributes,omitempty"`
		*Alias
	}{
		Attributes: s.getAttributesWithoutExtensions(),
		Alias:      (*Alias)(s),
	})
}

func (s *Schema) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type Alias Schema
	// skip extension attributes since they're not supposed to be shown in the schema representation
	// but only used for hydrating resource's data
	return e.EncodeElement(&struct {
		Attributes []*Attribute `xml:"Attributes,omitempty"`
		*Alias
	}{
		Attributes: s.getAttributesWithoutExtensions(),
		Alias:      (*Alias)(s),
	}, start)
}

// returns the slice of attributes except attributes which're marked as extension ones
func (s *Schema) getAttributesWithoutExtensions() []*Attribute {
	result := make([]*Attribute, 0)
	for _, a := range s.Attributes {
		if !a.IsExtensionAttribute {
			result = append(result, a)
		}
	}

	return result
}

// Check if the given path corresponds to the attribute
func (a *Attribute) HasPath(p navigation.Path) bool {
	for _, levelVal := range p.CollectLevelValues() {
		switch strings.ToLower(levelVal) {
		case strings.ToLower(a.Navigation.FullPath):
			return true
		case strings.ToLower(a.Navigation.Path):
			return true
		}
	}
	return false
}

// Get sub-attribute of this attribute determined by the given path
func (a *Attribute) GetAttribute(p navigation.Path, recursive bool) *Attribute {
	if p == nil {
		return a
	}

	if a.SubAttributes != nil {
		for _, subAttr := range a.SubAttributes {
			if strings.ToLower(subAttr.Name) == strings.ToLower(p.Base()) {
				if recursive {
					return subAttr.GetAttribute(p.Next(), recursive)
				} else {
					return subAttr
				}

			}
		}
	}
	return nil
}

func (a *Attribute) Clone() *Attribute {
	var subAttributes []*Attribute
	if a.SubAttributes != nil {
		subAttributes = make([]*Attribute, len(a.SubAttributes))
		for i, subAttr := range a.SubAttributes {
			subAttributes[i] = subAttr.Clone()
		}
	}

	return &Attribute{
		Name:            a.Name,
		Type:            a.Type,
		MultiValued:     a.MultiValued,
		Required:        a.Required,
		CaseExact:       a.CaseExact,
		Mutability:      a.Mutability,
		Returned:        a.Returned,
		Uniqueness:      a.Uniqueness,
		ReferenceTypes:  a.ReferenceTypes,
		CanonicalValues: a.CanonicalValues,
		Description:     a.Description,
		SubAttributes:   subAttributes,
		Navigation:      a.Navigation,
	}
}

// Check if this attribute should be returned in the service response,
// taking into account the attribute properties and
// the lists of included/excluded attribute given in the request
func (a *Attribute) IsReturned(included, excluded []navigation.Path) bool {
	// Schemas and meta (and their sub-attributes) are ALWAYS included
	for _, path := range []string{"schemas", "meta"} {
		if strings.Index(strings.ToLower(a.Navigation.FullPath), path) == 0 {
			return true
		}
	}

	// Overriding the default set
	if included != nil && len(included) > 0 {
		for _, p := range included {
			if a.HasPath(p) {
				return true
			}
		}
		return false
	}

	// Must not be returned
	if excluded != nil && len(excluded) > 0 {
		for _, p := range excluded {
			if a.HasPath(p) && a.Returned != "always" {
				return false
			}
		}
	}

	// Return-ability,
	// returned=request would not reach here if specified, hence it's not requested
	switch a.Returned {
	case "always":
		return true
	case "never", "request":
		return false
	}

	return true
}

// Get item of map for the given key or get field of struct for the given name.
// Key/name is attribute navigation name
func (a *Attribute) Item(data interface{}) (reflect.StructField, reflect.Value, error) {
	emptyField := reflect.StructField{}
	emptyValue := reflect.Value{}

	if len(a.Navigation.FieldName) == 0 {
		return emptyField, emptyValue, fmt.Errorf("Attribute %s has empty field name", a.Name)
	}

	attrFieldName := strings.ToLower(a.Navigation.FieldName)
	val := system.ReflectValue(data)

	switch val.Kind() {
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			subField := val.Type().Field(i)
			subVal := val.Field(i)
			if attrFieldName == strings.ToLower(subField.Name) {
				return subField, subVal, nil
			}
		}
	case reflect.Map:
		for _, mKey := range val.MapKeys() {
			if attrFieldName == strings.ToLower(mKey.String()) {
				mapVal := val.MapIndex(mKey)
				mapField := reflect.StructField{
					Name: mKey.String(),
					Type: mapVal.Type(),
				}
				return mapField, mapVal, nil
			}
		}
	default:
		return emptyField, emptyValue, fmt.Errorf("Cannot look for field %s for attribute %s since the given data is not a map or a struct but %s", a.Navigation.FieldName, a.Name, val.Kind())
	}

	// if value absent for some attribute, just return empty
	return emptyField, emptyValue, nil
}

// Get item of map for the given key or zero value
// Key is attribute navigation name
func (a *Attribute) MapItem(data interface{}) (reflect.Value, error) {
	emptyValue := reflect.Value{}

	if len(a.Navigation.FieldName) == 0 {
		return emptyValue, fmt.Errorf("Attribute %s has empty field name", a.Name)
	}

	attrFieldName := strings.ToLower(a.Navigation.FieldName)
	val := system.ReflectValue(data)

	switch val.Kind() {
	case reflect.Map:
		for _, mKey := range val.MapKeys() {
			if attrFieldName == strings.ToLower(mKey.String()) {
				return val.MapIndex(mKey), nil
			}
		}
	default:
		return emptyValue, fmt.Errorf("Cannot look for field %s for attribute %s since the given data is not a map but %s", a.Navigation.FieldName, a.Name, val.Kind())
	}

	return emptyValue, nil
}

func (a *Attribute) CorrectPathCase(path navigation.Path, recursive bool) {
	subAttr := a.GetAttribute(path, false)
	if subAttr == nil {
		return
	}

	switch strings.ToLower(path.Base()) {
	case strings.ToLower(subAttr.Name):
		path.SetBase(subAttr.Name)
	case strings.ToLower(subAttr.Navigation.FullPath):
		path.SetBase(subAttr.Navigation.FullPath)
	}

	if path.FilterRoot() != nil {
		subAttr.CorrectFilterNodeCase(path.FilterRoot())
	}
	if recursive && path.Next() != nil {
		subAttr.CorrectPathCase(path.Next(), recursive)
	}
}

func (a *Attribute) CorrectFilterNodeCase(node navigation.FilterNode) {
	nodeReflect := system.ReflectValue(node)
	if nodeReflect.Kind() == reflect.Invalid {
		return
	} else if nodeReflect.Kind() == reflect.Interface && nodeReflect.Elem().Kind() == reflect.Invalid {
		return
	}

	if node.Left() != nil {
		a.CorrectFilterNodeCase(node.Left())
	}
	if node.Type() == navigation.PathOperand {
		a.CorrectPathCase(node.Data().(navigation.Path), true)
	}
	if node.Right() != nil {
		a.CorrectFilterNodeCase(node.Right())
	}
}

func (a *Attribute) TypeExpectation() string {
	expects := ""
	switch a.Type {
	case "string", "binary", "reference", "datetime":
		expects = "string"
	case "integer":
		expects = "integer"
	case "decimal":
		expects = "decimal"
	case "boolean":
		expects = "boolean"
	case "complex":
		expects = "complex"
	}
	if a.MultiValued {
		expects += " array"
	}
	return expects
}

func (a *Attribute) ExpectsString() bool {
	switch a.Type {
	case "string", "datetime", "reference", "binary":
		return !a.MultiValued
	default:
		return false
	}
}

func (a *Attribute) ExpectsStringArray() bool {
	switch a.Type {
	case "string", "datetime", "reference", "binary":
		return a.MultiValued
	default:
		return false
	}
}

func (a *Attribute) ExpectsInteger() bool {
	return !a.MultiValued && a.Type == "integer"
}

func (a *Attribute) ExpectsFloat() bool {
	return !a.MultiValued && a.Type == "decimal"
}

func (a *Attribute) ExpectsBool() bool {
	return !a.MultiValued && a.Type == "boolean"
}

func (a *Attribute) ExpectsBinary() bool {
	return !a.MultiValued && a.Type == "binary"
}

func (a *Attribute) ExpectsComplex() bool {
	return !a.MultiValued && a.Type == "complex"
}

func (a *Attribute) ExpectsComplexArray() bool {
	return a.MultiValued && a.Type == "complex"
}
