package modelsResources

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"service-scim/models/normalization"
	"service-scim/resources/schemas"
	"service-scim/system"
	"strconv"
	"time"
)

type Group struct {
	Abstract

	Schemas     []string       `json:"schemas" xml:"Schemas"`
	Id          int64          `json:"id" xml:"Id"`
	DisplayName string         `json:"displayName" xml:"DisplayName"`
	ExternalId  string         `json:"externalId,omitempty" xml:"ExternalId,omitempty"`
	Updated     string         `json:"-" xml:"-"`
	Members     []*GroupMember `json:"members,omitempty" xml:"Members,omitempty"`
	Meta        *AbstractMeta  `json:"meta" xml:"Meta"`
	XMLName     struct{}       `json:"-" xml:"Group"`
}

type GroupMember struct {
	Value   string   `json:"value,omitempty" xml:"Value,omitempty"`
	Ref     string   `json:"$ref,omitempty" xml:"Ref,omitempty"`
	Type    string   `json:"type,omitempty" xml:"Type,omitempty"`
	Display string   `json:"display,omitempty" xml:"Display,omitempty"`
	Roles   []*Role  `json:"roles,omitempty" xml:"Roles,omitempty"`
	XMLName struct{} `json:"-" xml:"Member"`
}

// Filter excluded attributes and keep included and required ones
func (g *Group) normalized() (interface{}, error) {
	if normalized, err := g.normalizator.Normalize(g, resourcesSchemas.GroupSchemaObject.Attributes, g.normalizationOptions.Included(), g.normalizationOptions.Excluded()); err != nil {
		return nil, err
	} else {
		return normalized.Interface(), nil
	}
}

func (g *Group) MarshalJSON() ([]byte, error) {
	if normalized, err := g.normalized(); err != nil {
		return nil, err
	} else {
		return json.Marshal(normalized)
	}
}

func (g *Group) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if normalized, err := g.normalized(); err != nil {
		return err
	} else {
		return e.EncodeElement(normalized, start)
	}
}

func (g *Group) AddSchemas() {
	g.Schemas = []string{"urn:ietf:params:scim:schemas:core:2.0:Group"}
}

func (g *Group) AddMeta() {
	modified := g.Updated
	if len(modified) > 0 {
		if modifiedParsed, err := time.ParseInLocation("2006-01-02 15:04:05", modified, time.UTC); err == nil {
			modified = modifiedParsed.Format("2006-01-02T15:04:05-07:00")
		}
	}

	g.Meta = &AbstractMeta{
		LastModified: modified,
		Location:     fmt.Sprintf("/%s/groups/%s", system.AppVersion(), g.ScimId()),
		ResourceType: "Group",
	}
}

func (g *Group) ScimId() string {
	return strconv.FormatInt(g.Id, 10)
}

func (g *Group) ToMap(tagWithFieldName string) interface{} {
	return system.StructToMap(g, tagWithFieldName)
}

func NewGroupFromMap(data map[string]interface{}, existingGroup *Group) *Group {
	var skipNil bool
	var group *Group

	if existingGroup == nil {
		group = new(Group)
		skipNil = false
	} else {
		group = existingGroup
		skipNil = true
	}

	system.SetStructInt64Fields(group, data, []string{"Id"}, skipNil)
	system.SetStructStringFields(group, data, []string{"DisplayName", "ExternalId"}, skipNil)

	members := make([]*GroupMember, 0)
	for _, val := range system.SliceOfMapsForStruct(group, data, "Members") {
		member := new(GroupMember)
		system.SetStructStringFields(member, val, []string{"Value", "Ref", "Type"}, skipNil)
		if !system.StructIsEmpty(member) {
			members = append(members, member)
		}
	}
	if len(members) > 0 || group.Members == nil {
		group.Members = members
	}

	group.SetNormalizationOptions(modelsNormalization.NewEmptyNormalizationOption())
	group.AddSchemas()
	group.AddMeta()
	return group
}
