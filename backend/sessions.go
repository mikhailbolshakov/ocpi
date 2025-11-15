package backend

import "time"

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
	PartyId     string `json:"partyId,omitempty"`     // PartyId should be unique within country
	CountryCode string `json:"countryCode,omitempty"` // CountryCode alfa-2 code
	Id          string `json:"uid"`                   // Id unique ID by which this Token can be identified
	Type        string `json:"type"`                  // Type token type
	ContractId  string `json:"contractId"`            // ContractId uniquely identifies the EV driver contract token
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
	TariffId              string          `json:"tariffId"`                        // TariffId unique identifier of the Tariff that is relevant for this Charging Period
	EncodingMethod        string          `json:"encodingMethod"`                  // EncodingMethod name of the encoding used in the SignedData field
	EncodingMethodVersion *int            `json:"encodingMethodVersion,omitempty"` // EncodingMethodVersion version of the EncodingMethod
	PublicKey             string          `json:"publicKey"`                       // PublicKey used to sign the data, base64 encoded
	SignedValues          []*SignedValue  `json:"signedValues"`                    // SignedValues one or more signed values
	Url                   string          `json:"url,omitempty"`                   // Url that can be shown to an EV driver
}

type Session struct {
	Id              string            `json:"id"`                        // Id uniquely identifies the session
	StartDateTime   *time.Time        `json:"startDateTime"`             // StartDateTime timestamp when the session became ACTIVE
	EndDateTime     *time.Time        `json:"endDateTime"`               // EndDateTime timestamp when the session was completed/finished
	Kwh             *float64          `json:"kwh"`                       // Kwh how many kWh were charged
	CdrToken        *CdrToken         `json:"cdrToken"`                  // CdrToken token used to start this charging session
	AuthMethod      string            `json:"authMethod"`                // AuthMethod method used for authentication
	AuthRef         string            `json:"authRef,omitempty"`         // AuthRef reference to the authorization given by the eMSP
	LocationId      string            `json:"locationId"`                // LocationId id of the location obj
	EvseId          string            `json:"evseId"`                    // EvseId id of the evse obj
	ConnectorId     string            `json:"connectorId"`               // ConnectorId id of the connector
	MeterId         string            `json:"meterId,omitempty"`         // MeterId id of the kWh meter
	Currency        string            `json:"currency"`                  // Currency ISO-4217 currency code
	ChargingPeriods []*ChargingPeriod `json:"chargingPeriods,omitempty"` // ChargingPeriods  list of Charging Periods that can be used to calculate and verify the total cost
	TotalCost       *Price            `json:"totalCost,omitempty"`       // TotalCost total cost of the session in the specified currency
	Status          string            `json:"status"`                    // Status session status
	LastUpdated     time.Time         `json:"lastUpdated"`               // LastUpdated when this Tariff was last updated
	PlatformId      string            `json:"platformId"`                // PlatformId rel to platform
	RefId           string            `json:"refId"`                     // RefId any external relation
	PartyId         string            `json:"partyId,omitempty"`         // PartyId should be unique within country
	CountryCode     string            `json:"countryCode,omitempty"`     // CountryCode alfa-2 code
}

type SessionSearchResponse struct {
	PageInfo *PageResponse `json:"pageInfo,omitempty"`
	Items    []*Session    `json:"items,omitempty"`
}
