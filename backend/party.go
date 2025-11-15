package backend

import "time"

type Party struct {
	// Id unique key
	Id string `json:"id,omitempty"`
	// Role
	Roles []string `json:"roles,omitempty"`
	// Party identified party
	PartyId string `json:"partyId,omitempty"`
	// Party country code
	CountryCode string `json:"countryCode,omitempty"`
	// BusinessDetails party company details
	BusinessDetails *BusinessDetails `json:"businessDetails,omitempty"`
	// RefId any external relation
	RefId string `json:"refId,omitempty"`
	// Status role status
	Status string `json:"status,omitempty"`
	// LastUpdated last updated time
	LastUpdated time.Time `json:"lastUpdated,omitempty"`
}

type PartySearchResponse struct {
	// page info
	PageInfo *PageResponse `json:"pageInfo,omitempty"`
	// items
	Items []*Party `json:"items,omitempty"`
}
