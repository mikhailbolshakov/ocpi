package domain

import "time"

const (
	PageSizeMaxLimit = 100
	PageSizeDefault  = 20
)

type PageRequest struct {
	Offset   *int       // Offset paging offset
	Limit    *int       // Limit paging limit
	DateFrom *time.Time // DateFrom updating period
	DateTo   *time.Time // DateTo updating period
}

type PageResponse struct {
	Total    *int         // Total number of objects available in the server
	Limit    *int         // Limit maximum number of objects that the server can return
	NextPage *PageRequest // NextPage paging criteria to request the next page
}

type DisplayText struct {
	Language string
	Text     string
}

type PartyExtId struct {
	PartyId     string `json:"partyId"`     // PartyId should be unique within country
	CountryCode string `json:"countryCode"` // CountryCode alfa-2 code
}

type OcpiItem struct {
	ExtId       PartyExtId `json:"extId"`       // ExtId party external ID
	PlatformId  string     `json:"platformId"`  // PlatformId rel to platform
	RefId       string     `json:"refId"`       // RefId any external relation
	LastUpdated time.Time  `json:"lastUpdated"` // LastUpdated last updated
	LastSent    *time.Time `json:"lastSent"`    // LastSent last sent
}
