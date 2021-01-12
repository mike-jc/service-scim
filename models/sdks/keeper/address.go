package modelsSdkKeeper

import (
	"service-scim/system"
)

type Address struct {
	Street     string `json:"street,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	Country    string `json:"country,omitempty"`
	PostalCode string `json:"postalCode,omitempty"`
}

func (a *Address) IsValid() bool {
	return len(a.Street) > 0 || len(a.City) > 0 || len(a.Country) > 0
}

func AddressIsValid(address map[string]interface{}) bool {
	return !system.MapValueIsEmpty(address, "street") ||
		!system.MapValueIsEmpty(address, "city") ||
		!system.MapValueIsEmpty(address, "country")
}
