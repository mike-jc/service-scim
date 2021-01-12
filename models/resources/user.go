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

type User struct {
	Abstract

	Schemas             []string        `json:"schemas" xml:"Schemas"`
	Id                  int64           `json:"id" xml:"Id"`
	Emails              []*Email        `json:"emails" xml:"Emails"`
	Name                *UserName       `json:"name,omitempty" xml:"Name,omitempty"`
	ExternalId          string          `json:"externalId,omitempty" xml:"ExternalId,omitempty"`
	UserName            string          `json:"userName,omitempty" xml:"UserName,omitempty"`
	ProfileUrl          string          `json:"profileUrl,omitempty" xml:"ProfileUrl,omitempty"`
	Title               string          `json:"title,omitempty" xml:"Title,omitempty"`
	Locale              string          `json:"locale,omitempty" xml:"Locale,omitempty"`
	Timezone            string          `json:"timezone,omitempty" xml:"Timezone,omitempty"`
	Active              bool            `json:"active,omitempty" xml:"Active,omitempty"`
	Password            string          `json:"password,omitempty" xml:"Password,omitempty"`
	Created             string          `json:"-" xml:"-"`
	Updated             string          `json:"-" xml:"-"`
	PhoneNumbers        []*PhoneNumber  `json:"phoneNumbers,omitempty" xml:"PhoneNumbers,omitempty"`
	Photos              []*Photo        `json:"photos,omitempty" xml:"Photos,omitempty"`
	Addresses           []*Address      `json:"addresses,omitempty" xml:"Addresses,omitempty"`
	Groups              []*UserGroup    `json:"groups,omitempty" xml:"Groups,omitempty"`
	Roles               []*Role         `json:"roles,omitempty" xml:"Roles,omitempty"`
	EnterpriseExtension *EnterpriseUser `json:"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User,omitempty" xml:"EnterpriseExtension"`
	Meta                *AbstractMeta   `json:"meta" xml:"Meta"`
	XMLName             struct{}        `json:"-" xml:"User"`
}

type Email struct {
	Value   string   `json:"value" xml:"Value"`
	Display string   `json:"display,omitempty" xml:"Display,omitempty"`
	Type    string   `json:"type,omitempty" xml:"Type,omitempty"`
	Primary bool     `json:"primary,omitempty" xml:"Primary,omitempty"`
	XMLName struct{} `json:"-" xml:"Email"`
}

type UserName struct {
	Formatted       string   `json:"formatted,omitempty" xml:"Formatted,omitempty"`
	FamilyName      string   `json:"familyName,omitempty" xml:"FamilyName,omitempty"`
	GivenName       string   `json:"givenName,omitempty" xml:"GivenName,omitempty"`
	MiddleName      string   `json:"middleName,omitempty" xml:"MiddleName,omitempty"`
	HonorificPrefix string   `json:"honorificPrefix,omitempty" xml:"HonorificPrefix,omitempty"`
	HonorificSuffix string   `json:"honorificSuffix,omitempty" xml:"HonorificSuffix,omitempty"`
	XMLName         struct{} `json:"-" xml:"Name"`
}

type PhoneNumber struct {
	Value   string   `json:"value" xml:"Value"`
	Display string   `json:"display,omitempty" xml:"Display,omitempty"`
	Type    string   `json:"type,omitempty" xml:"Type,omitempty"`
	Primary bool     `json:"primary,omitempty" xml:"Primary,omitempty"`
	XMLName struct{} `json:"-" xml:"PhoneNumber"`
}

type Photo struct {
	Value   string   `json:"value,omitempty" xml:"Value,omitempty"`
	Type    string   `json:"type,omitempty" xml:"Type,omitempty"`
	XMLName struct{} `json:"-" xml:"Photo"`
}

type Address struct {
	StreetAddress string   `json:"streetAddress,omitempty" xml:"StreetAddress,omitempty"`
	Locality      string   `json:"locality,omitempty" xml:"Locality,omitempty"`
	Region        string   `json:"region,omitempty" xml:"Region,omitempty"`
	PostalCode    string   `json:"postalCode,omitempty" xml:"PostalCode,omitempty"`
	Country       string   `json:"country,omitempty" xml:"Country,omitempty"`
	Type          string   `json:"type" xml:"Type"`
	Primary       bool     `json:"primary" xml:"Primary"`
	XMLName       struct{} `json:"-" xml:"Address"`
}

type UserGroup struct {
	Value   string   `json:"value,omitempty" xml:"Value,omitempty"`
	Ref     string   `json:"$ref,omitempty" xml:"Ref,omitempty"`
	Display string   `json:"display,omitempty" xml:"Display,omitempty"`
	XMLName struct{} `json:"-" xml:"Group"`
}

type Role struct {
	Value   string   `json:"value,omitempty" xml:"Value,omitempty"`
	Display string   `json:"display,omitempty" xml:"Display,omitempty"`
	XMLName struct{} `json:"-" xml:"Role"`
}

// Filter excluded attributes and keep included and required ones
func (u *User) normalized() (interface{}, error) {
	if u.EnterpriseExtension != nil && u.EnterpriseExtension.IsEmpty() {
		u.EnterpriseExtension = nil
	}

	if normalized, err := u.normalizator.Normalize(u, resourcesSchemas.UserSchemaObject.Attributes, u.normalizationOptions.Included(), u.normalizationOptions.Excluded()); err != nil {
		return nil, err
	} else {
		return normalized.Interface(), nil
	}
}

func (u *User) MarshalJSON() ([]byte, error) {
	if normalized, err := u.normalized(); err != nil {
		return nil, err
	} else {
		return json.Marshal(normalized)
	}
}

func (u *User) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if normalized, err := u.normalized(); err != nil {
		return err
	} else {
		return e.EncodeElement(normalized, start)
	}
}

func (u *User) AddSchemas() {
	u.Schemas = []string{"urn:ietf:params:scim:schemas:core:2.0:User"}
	if u.EnterpriseExtension != nil && !u.EnterpriseExtension.IsEmpty() {
		u.Schemas = append(u.Schemas, "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User")
	}
}

func (u *User) AddMeta() {
	created := u.Created
	if len(created) > 0 {
		if createdParsed, err := time.ParseInLocation("2006-01-02 15:04:05", created, time.UTC); err == nil {
			created = createdParsed.Format("2006-01-02T15:04:05-07:00")
		}
	}

	modified := u.Updated
	if len(modified) > 0 {
		if modifiedParsed, err := time.ParseInLocation("2006-01-02 15:04:05", modified, time.UTC); err == nil {
			modified = modifiedParsed.Format("2006-01-02T15:04:05-07:00")
		}
	}

	u.Meta = &AbstractMeta{
		Created:      created,
		LastModified: modified,
		Location:     fmt.Sprintf("/%s/users/%s", system.AppVersion(), u.ScimId()),
		ResourceType: "User",
	}
}

func (u *User) ScimId() string {
	return strconv.FormatInt(u.Id, 10)
}

func (u *User) PrimaryEmailValue() string {
	if u.Emails != nil && len(u.Emails) > 0 {
		// try to find primary, or working, or any first email
		primaryEmail := ""
		workingEmail := ""
		firstEmail := ""
		for _, email := range u.Emails {
			if email.Primary && len(primaryEmail) == 0 {
				primaryEmail = email.Value
			}
			if email.Type == "work" && len(workingEmail) == 0 {
				workingEmail = email.Value
			}
			if len(firstEmail) == 0 {
				firstEmail = email.Value
			}
		}
		if len(primaryEmail) > 0 {
			return primaryEmail
		} else if len(workingEmail) > 0 {
			return workingEmail
		} else {
			return firstEmail
		}

	}
	return ""
}

func (u *User) NameValue() string {
	if u.Name != nil {
		if len(u.Name.Formatted) > 0 {
			return u.Name.Formatted
		} else if len(u.Name.MiddleName) > 0 {
			return fmt.Sprintf("%s %s %s", u.Name.GivenName, u.Name.MiddleName, u.Name.FamilyName)
		} else {
			return fmt.Sprintf("%s %s", u.Name.GivenName, u.Name.FamilyName)
		}
	}
	return ""
}

func (u *User) PrimaryPhoneNumberValue() string {
	if u.PhoneNumbers != nil && len(u.PhoneNumbers) > 0 {
		// try to find primary, or working, or any first phone number
		primaryPhoneNumber := ""
		workingPhoneNumber := ""
		firstPhoneNumber := ""
		for _, number := range u.PhoneNumbers {
			if number.Primary && len(primaryPhoneNumber) == 0 {
				primaryPhoneNumber = number.Value
			}
			if number.Type == "work" && len(workingPhoneNumber) == 0 {
				workingPhoneNumber = number.Value
			}
			if len(firstPhoneNumber) == 0 {
				firstPhoneNumber = number.Value
			}
		}
		if len(primaryPhoneNumber) > 0 {
			return primaryPhoneNumber
		} else if len(workingPhoneNumber) > 0 {
			return workingPhoneNumber
		} else {
			return firstPhoneNumber
		}

	}
	return ""
}

func (u *User) PrimaryAddress() *Address {
	if u.Addresses != nil && len(u.Addresses) > 0 {
		// try to find primary, or working, or any first phone number
		var primaryAddress *Address
		var workingAddress *Address
		var firstAddress *Address
		for _, address := range u.Addresses {
			if address.Primary && primaryAddress == nil {
				primaryAddress = address
			}
			if address.Type == "work" && workingAddress == nil {
				workingAddress = address
			}
			if firstAddress == nil {
				firstAddress = address
			}
		}
		if primaryAddress == nil {
			return primaryAddress
		} else if workingAddress == nil {
			return workingAddress
		} else {
			return firstAddress
		}

	}
	return nil
}

func (u *User) ToMap(tagWithFieldName string) interface{} {
	return system.StructToMap(u, tagWithFieldName)
}

func NewUserFromMap(data map[string]interface{}, existingUser *User) *User {
	var skipNil bool
	var user *User

	if existingUser == nil {
		user = new(User)
		skipNil = false
	} else {
		user = existingUser
		skipNil = true
	}

	system.SetStructInt64Fields(user, data, []string{"Id"}, skipNil)

	system.SetStructStringFields(user, data, []string{
		"ExternalId", "UserName", "ProfileUrl", "Title",
		"Locale", "Timezone", "Password", "PhoneNumber",
	}, skipNil)

	system.SetStructBoolFields(user, data, []string{"Active"}, skipNil)

	emails := make([]*Email, 0)
	for _, val := range system.SliceOfMapsForStruct(user, data, "Emails") {
		email := new(Email)
		system.SetStructStringFields(email, val, []string{"Value", "Type"}, skipNil)
		system.SetStructBoolFields(email, val, []string{"Primary"}, skipNil)
		if !system.StructIsEmpty(email) {
			emails = append(emails, email)
		}
	}
	if len(emails) > 0 || user.Emails == nil {
		user.Emails = emails
	}

	if nameInterface := system.MapValueForStruct(user, data, "Name"); nameInterface != nil {
		if nameVal, ok := nameInterface.(map[string]interface{}); ok {
			name := new(UserName)
			system.SetStructStringFields(name, nameVal, []string{
				"Formatted", "FamilyName", "GivenName", "MiddleName",
			}, skipNil)
			if !system.StructIsEmpty(name) || user.Name == nil {
				user.Name = name
			}
		}
	}

	photos := make([]*Photo, 0)
	for _, val := range system.SliceOfMapsForStruct(user, data, "Photos") {
		photo := new(Photo)
		system.SetStructStringFields(photo, val, []string{"Value", "Type"}, skipNil)
		if !system.StructIsEmpty(photo) {
			photos = append(photos, photo)
		}
	}
	if len(photos) > 0 || user.Photos == nil {
		user.Photos = photos
	}

	addresses := make([]*Address, 0)
	for _, val := range system.SliceOfMapsForStruct(user, data, "Addresses") {
		addr := new(Address)
		system.SetStructStringFields(addr, val, []string{
			"StreetAddress", "Locality", "Region", "PostalCode", "Country",
		}, skipNil)
		if !system.StructIsEmpty(addr) {
			addresses = append(addresses, addr)
		}
	}
	if len(addresses) > 0 || user.Addresses == nil {
		user.Addresses = addresses
	}

	groups := make([]*UserGroup, 0)
	for _, val := range system.SliceOfMapsForStruct(user, data, "Groups") {
		group := new(UserGroup)
		system.SetStructStringFields(group, val, []string{"Value", "Ref", "Display"}, skipNil)
		if !system.StructIsEmpty(group) {
			groups = append(groups, group)
		}
	}
	if len(groups) > 0 || user.Groups == nil {
		user.Groups = groups
	}

	roles := make([]*Role, 0)
	for _, val := range system.SliceOfMapsForStruct(user, data, "Roles") {
		role := new(Role)
		system.SetStructStringFields(role, val, []string{"Value", "Display"}, skipNil)
		if !system.StructIsEmpty(role) {
			roles = append(roles, role)
		}
	}
	if len(roles) > 0 || user.Roles == nil {
		user.Roles = roles
	}

	user.SetNormalizationOptions(modelsNormalization.NewEmptyNormalizationOption())
	user.AddSchemas()
	user.AddMeta()
	return user
}
