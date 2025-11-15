package domain

import (
	"context"
	"time"
)

// Image linked to party
type Image struct {
	Url       string `json:"url,omitempty"`       // Url of the image
	Thumbnail string `json:"thumbnail,omitempty"` // Thumbnail of the image
	Category  string `json:"category,omitempty"`  // Category
	Type      string `json:"type,omitempty"`      // Type file type (jpeg, giff etc)
	Width     int    `json:"width,omitempty"`     // Width size
	Height    int    `json:"height,omitempty"`    // Height size
}

// BusinessDetails is party additional data
type BusinessDetails struct {
	Name    string `json:"name,omitempty"`    // Name company name
	Website string `json:"website,omitempty"` // Website company website
	Logo    *Image `json:"logo,omitempty"`    // Logo company logo
	Inn     string `json:"inn,omitempty"`     // Inn (tax number) of the company (Yandex extension of the protocol. Isn't supported by OCPI)
}

// Party OCPI party
type Party struct {
	OcpiItem
	Id              string           `json:"id"`                        // Id unique key
	Roles           []string         `json:"roles"`                     // Roles party roles
	BusinessDetails *BusinessDetails `json:"businessDetails,omitempty"` // BusinessDetails party company details
	Status          string           `json:"status"`                    // Status party status
}

type PartySearchCriteria struct {
	PageRequest
	IncRoles     []string    // IncRoles includes list of roles
	ExcRoles     []string    // ExcRoles excludes list of roles
	IncPlatforms []string    // IncPlatforms includes platform Ids
	ExcPlatforms []string    // ExcPlatforms exclude platform Ids
	Ids          []string    // Ids list of Ids
	ExtId        *PartyExtId // ExtId by party ext ID
	RefId        string      // RefId by ref id
}

type PartySearchResponse struct {
	PageResponse
	Items []*Party
}

type PartyService interface {
	// GetByExtId retrieves parties by ext ID
	GetByExtId(ctx context.Context, extId PartyExtId) (*Party, error)
	// GetByPlatform retrieves parties by platform
	GetByPlatform(ctx context.Context, platformId string) ([]*Party, error)
	// Search searches parties
	Search(ctx context.Context, cr *PartySearchCriteria) (*PartySearchResponse, error)
	// Get retrieves by id
	Get(ctx context.Context, id string) (*Party, error)
	// Merge merges party
	Merge(ctx context.Context, party *Party) (*Party, error)
	// MergeMany merges multiple parties
	MergeMany(ctx context.Context, parties ...*Party) error
	// MarkSent marks parties as sent by setting up sent date = now
	MarkSent(ctx context.Context, partyIds ...string) error
	// DeletePartyByExtId deletes party by ext id
	DeletePartyByExtId(ctx context.Context, extId PartyExtId) error
}

type PartyStorage interface {
	// GetPartyByExtId retrieves party by ext ID
	GetPartyByExtId(ctx context.Context, extId PartyExtId) (*Party, error)
	// GetPartiesByPlatform retrieves parties by platform
	GetPartiesByPlatform(ctx context.Context, platformId string) ([]*Party, error)
	// GetByRefId retrieves by refId
	GetByRefId(ctx context.Context, refId string) (*Party, error)
	// GetParty retrieves by id
	GetParty(ctx context.Context, id string) (*Party, error)
	// CreateParty creates a new party
	CreateParty(ctx context.Context, party *Party) error
	// UpdateParty updates a party
	UpdateParty(ctx context.Context, party *Party) error
	// MarkSentParties marks parties as sent by setting up last_sent
	MarkSentParties(ctx context.Context, date time.Time, partyIds ...string) error
	// Search searches parties
	Search(ctx context.Context, criteria *PartySearchCriteria) (*PartySearchResponse, error)
	// DeletePartyByExtId deletes party by ext id
	DeletePartyByExtId(ctx context.Context, extId PartyExtId) error
}
