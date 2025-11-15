package domain

import (
	"context"
	"time"
)

const (
	AuthMethodRequest   = "AUTH_REQUEST"
	AuthMethodCommand   = "COMMAND"
	AuthMethodWhitelist = "WHITELIST"

	SessionStatusActive      = "ACTIVE"
	SessionStatusCompleted   = "COMPLETED"
	SessionStatusInvalid     = "INVALID"
	SessionStatusPending     = "PENDING"
	SessionStatusReservation = "RESERVATION"

	DimensionTypeCurrent         = "CURRENT"
	DimensionTypeEnergy          = "ENERGY"
	DimensionTypeEnergyExport    = "ENERGY_EXPORT"
	DimensionTypeEnergyImport    = "ENERGY_IMPORT"
	DimensionTypeMaxCurrent      = "MAX_CURRENT"
	DimensionTypeMinCurrent      = "MIN_CURRENT"
	DimensionTypeMaxPower        = "MAX_POWER"
	DimensionTypeMinPower        = "MIN_POWER"
	DimensionTypeParkingTime     = "PARKING_TIME"
	DimensionTypePower           = "POWER"
	DimensionTypeReservationTime = "RESERVATION_TIME"
	DimensionTypeStateOfCharge   = "STATE_OF_CHARGE"
	DimensionTypeTime            = "TIME"
)

type CdrToken struct {
	PartyExtId
	Id         string `json:"uid"`        // Id unique ID by which this Token can be identified
	Type       string `json:"type"`       // Type token type
	ContractId string `json:"contractId"` // ContractId uniquely identifies the EV driver contract token
}

type CdrDimension struct {
	Type   string  `json:"type"`   // Type of CDR dimension
	Volume float64 `json:"volume"` // Volume of the dimension
}

type SignedValue struct {
	Nature     string `json:"nature"`     // Nature of the value
	PlainData  string `json:"plainData"`  // PlainData un-encoded string of data
	SignedData string `json:"signedData"` // SignedData blob of signed data, base64 encoded
}

type ChargingPeriod struct {
	StartDateTime         time.Time       `json:"startDateTime"`                   // StartDateTime start timestamp of the charging period
	Dimensions            []*CdrDimension `json:"dimensions"`                      // Dimensions list of relevant values for this charging period
	TariffId              string          `json:"tariffId,omitempty"`              // TariffId unique identifier of the Tariff that is relevant for this Charging Period
	EncodingMethod        string          `json:"encodingMethod,omitempty"`        // EncodingMethod name of the encoding used in the SignedData field
	EncodingMethodVersion *int            `json:"encodingMethodVersion,omitempty"` // EncodingMethodVersion version of the EncodingMethod
	PublicKey             string          `json:"publicKey,omitempty"`             // PublicKey used to sign the data, base64 encoded
	SignedValues          []*SignedValue  `json:"signedValues,omitempty"`          // SignedValues one or more signed values
	Url                   string          `json:"url,omitempty"`                   // Url that can be shown to an EV driver
}

type SessionDetails struct {
	StartDateTime *time.Time `json:"startDateTime"`       // StartDateTime timestamp when the session became ACTIVE
	EndDateTime   *time.Time `json:"endDateTime"`         // EndDateTime timestamp when the session was completed/finished
	Kwh           *float64   `json:"kwh"`                 // Kwh how many kWh were charged
	CdrToken      *CdrToken  `json:"cdrToken"`            // CdrToken token used to start this charging session
	AuthMethod    string     `json:"authMethod"`          // AuthMethod method used for authentication
	AuthRef       string     `json:"authRef,omitempty"`   // AuthRef reference to the authorization given by the eMSP
	LocationId    string     `json:"locationId"`          // LocationId id of the location obj
	EvseId        string     `json:"evseId"`              // EvseId id of the evse obj
	ConnectorId   string     `json:"connectorId"`         // ConnectorId id of the connector
	MeterId       string     `json:"meterId,omitempty"`   // MeterId id of the kWh meter
	Currency      string     `json:"currency"`            // Currency ISO-4217 currency code
	TotalCost     *Price     `json:"totalCost,omitempty"` // TotalCost total cost of the session in the specified currency
	Status        string     `json:"status"`              // Status session status
}

type Session struct {
	OcpiItem
	Id              string            `json:"id"`                        // Id uniquely identifies the session
	Details         SessionDetails    `json:"details"`                   // Details session details
	ChargingPeriods []*ChargingPeriod `json:"chargingPeriods,omitempty"` // ChargingPeriods  list of Charging Periods that can be used to calculate and verify the total cost
}

type SessionSearchCriteria struct {
	PageRequest
	ExtId               *PartyExtId // ExtId by party ext ID
	RefId               string      // RefId by ref id
	IncPlatforms        []string    // IncPlatforms includes platform Ids
	ExcPlatforms        []string    // ExcPlatforms exclude platform Ids
	Ids                 []string    // Ids by list of Ids
	AuthRef             string      // AuthRef by auth ref
	WithChargingPeriods bool        // WithChargingPeriods if true, retrieve charging periods for ech item
}

type SessionSearchResponse struct {
	PageResponse
	Items []*Session
}

type SessionService interface {
	// PutSession creates or updates session
	PutSession(ctx context.Context, sess *Session) (*Session, error)
	// MergeSession merges session
	MergeSession(ctx context.Context, sess *Session) (*Session, error)
	// GetSession retrieves session by ID
	GetSession(ctx context.Context, sessId string) (*Session, error)
	// GetSessionWithPeriods retrieves session by ID with periods
	GetSessionWithPeriods(ctx context.Context, sessId string) (*Session, error)
	// DeleteSessionsByExtId deletes all sessions by party ext id
	DeleteSessionsByExtId(ctx context.Context, extId PartyExtId) error
	// SearchSessions searches sessions
	SearchSessions(ctx context.Context, cr *SessionSearchCriteria) (*SessionSearchResponse, error)
}

type SessionStorage interface {
	// MergeSession creates or updates session
	MergeSession(ctx context.Context, sess *Session) error
	// UpdateSession updates session
	UpdateSession(ctx context.Context, sess *Session) error
	// GetSession retrieves session by ID
	GetSession(ctx context.Context, sessId string, withChargingPeriods bool) (*Session, error)
	// DeleteSessionsByExtId deletes all sessions by party ext id
	DeleteSessionsByExtId(ctx context.Context, extId PartyExtId) error
	// SearchSessions searches sessions
	SearchSessions(ctx context.Context, cr *SessionSearchCriteria) (*SessionSearchResponse, error)
	// CreateChargingPeriods creates new entries for charging periods
	CreateChargingPeriods(ctx context.Context, sess *Session, periods []*ChargingPeriod) error
	// UpdateChargingPeriods updates (delete -> insert) all the charging periods for the session
	UpdateChargingPeriods(ctx context.Context, sess *Session, periods []*ChargingPeriod) error
	// GetChargingPeriods retrieves charging periods by ID
	GetChargingPeriods(ctx context.Context, sessIds ...string) (map[string][]*ChargingPeriod, error)
}
