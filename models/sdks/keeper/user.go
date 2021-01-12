package modelsSdkKeeper

import (
	"fmt"
	"service-scim/models/resources"
	"service-scim/system"
	"strings"
)

type User struct {
	Id           int64    `json:"id"`
	ScimId       string   `json:"scimId"`
	Email        string   `json:"email"`
	Password     string   `json:"password,omitempty"`
	Name         string   `json:"name"`
	Active       bool     `json:"active"`
	Locale       string   `json:"locale,omitempty"`
	Timezone     string   `json:"timezone,omitempty"`
	Phone        string   `json:"phone,omitempty"`
	ProfileUrl   string   `json:"profileUrl,omitempty"`
	JobTitle     string   `json:"jobTitle,omitempty"`
	CreatedAt    string   `json:"createdAt,omitempty"`
	UpdatedAt    string   `json:"updatedAt,omitempty"`
	LocationName string   `json:"locationName,omitempty"`
	Address      *Address `json:"address,omitempty"`
	Photos       []*Photo `json:"photos,omitempty"`
	Groups       []*Group `json:"groups,omitempty"`
	Roles        []*Role  `json:"roles,omitempty"`
}

func (u *User) IsValid() bool {
	return u.Id > 0
}

func (u *User) Ref() string {
	return fmt.Sprintf("/%s/users/%d", system.AppVersion(), u.Id)
}

func UserIsValid(user map[string]interface{}) bool {
	return !system.MapValueIsEmpty(user, "id")
}

// build user name as space-separated list of all name parts
func UserName(name map[string]interface{}) string {
	// parts of the name have priority over the formatted name
	familyName := system.ToString(name["familyName"])
	middleName := system.ToString(name["middleName"])
	givenName := system.ToString(name["givenName"])
	if familyName != "" || givenName != "" {
		if middleName != "" {
			return strings.TrimSpace(fmt.Sprintf("%s %s %s", givenName, middleName, familyName))
		} else {
			return strings.TrimSpace(fmt.Sprintf("%s %s", givenName, familyName))
		}
	}

	formatted := system.ToString(name["formatted"])
	return strings.TrimSpace(formatted)
}

// name is space-separated list of all name parts (parts are family, middle and given name; all parts are optional)
func ParseUserName(userName string, user *modelsResources.User) {
	userName = strings.TrimSpace(userName)
	user.Name = &modelsResources.UserName{
		Formatted: userName,
	}

	// given and family names are first and last parts correspondingly, the rest is middle name
	parts := strings.Split(userName, " ")
	if strings.Index(userName, ",") == -1 {
		// name is most probably built from name parts on SCIM service side so given name is first part
		user.Name.GivenName = strings.TrimSpace(parts[0])
		if len(parts) > 1 {
			user.Name.FamilyName = strings.TrimSpace(parts[len(parts)-1])
			if len(parts) > 2 {
				user.Name.MiddleName = strings.TrimSpace(strings.Join(parts[1:len(parts)-1], " "))
			}
		}
	} else {
		// name is most probably formatted by SCIM client and family name is first part
		cutset := " ,()"
		user.Name.FamilyName = strings.Trim(parts[0], cutset)
		if len(parts) > 1 {
			user.Name.GivenName = strings.Trim(parts[len(parts)-1], cutset)
			if len(parts) > 2 {
				user.Name.MiddleName = strings.Trim(strings.Join(parts[1:len(parts)-1], " "), cutset)
			}
		}
	}
}

func LocationName(enterpriseData map[string]interface{}) string {
	var parts []string
	if organization := system.ToString(enterpriseData["organization"]); len(organization) > 0 {
		parts = []string{organization}
	} else {
		parts = []string{"-"}
	}

	if department := system.ToString(enterpriseData["department"]); len(department) > 0 {
		parts = append(parts, department)
	}

	return strings.Join(parts, " / ")
}

func ParseLocationName(locationName string, user *modelsResources.EnterpriseUser) {
	s := strings.SplitN(locationName, "/", 2)
	user.Organization = strings.TrimSpace(s[0])
	if len(s) > 1 {
		user.Department = strings.TrimSpace(s[1])
	}
}

func PrimaryEmail(emails []map[string]interface{}) string {
	if emails != nil && len(emails) > 0 {
		primaryEmail := ""
		workingEmail := ""
		firstEmail := ""
		for _, email := range emails {
			if system.ToBool(email["primary"]) && len(primaryEmail) == 0 {
				primaryEmail = system.ToString(email["value"])
			}
			if system.ToString(email["type"]) == "work" && len(workingEmail) == 0 {
				workingEmail = system.ToString(email["value"])
			}
			if len(firstEmail) == 0 {
				firstEmail = system.ToString(email["value"])
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

func PrimaryPhone(phones []map[string]interface{}) string {
	if phones != nil && len(phones) > 0 {
		primaryPhone := ""
		workingPhone := ""
		firstPhone := ""
		for _, ph := range phones {
			if system.ToBool(ph["primary"]) && len(primaryPhone) == 0 {
				primaryPhone = system.ToString(ph["value"])
			}
			if system.ToString(ph["type"]) == "work" && len(workingPhone) == 0 {
				workingPhone = system.ToString(ph["value"])
			}
			if len(firstPhone) == 0 {
				firstPhone = system.ToString(ph["value"])
			}
		}
		if len(primaryPhone) > 0 {
			return primaryPhone
		} else if len(workingPhone) > 0 {
			return workingPhone
		} else {
			return firstPhone
		}
	}
	return ""
}
