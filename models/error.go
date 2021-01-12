package models

type Error struct {
	Schemas  []string `json:"schemas" xml:"Schemas"`
	Status   string   `json:"status" xml:"Status"`
	ScimType string   `json:"scimType" xml:"ScimType"`
	Detail   string   `json:"detail" xml:"Detail"`
}
