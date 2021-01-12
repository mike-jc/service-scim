package modelsSdkKeeper

import (
	"net/url"
)

type Photo struct {
	Url  string `json:"url,omitempty"`
	Type string `json:"type,omitempty"`
}

func (ph *Photo) IsValid() bool {
	if _, err := url.ParseRequestURI(ph.Url); err != nil {
		return false
	}

	switch ph.Type {
	case "photo", "thumbnail":
		return true
	default:
		return false
	}
}

func PhotoIsValid(photo map[string]interface{}) bool {
	if val, exists := photo["url"]; exists {
		if _, err := url.ParseRequestURI(val.(string)); err == nil {
			return true
		}
	}
	return false
}
