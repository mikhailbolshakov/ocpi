package domain

import (
	"context"
	"time"
)

type CdrLocation struct {
	Id                 string      `json:"id"`                   // Id Uniquely identifies the location within the CPO’s platform
	Name               string      `json:"name,omitempty"`       // Name of the location
	Address            string      `json:"address"`              // Address street/block name and house number if available
	City               string      `json:"city"`                 // City or town
	PostalCode         string      `json:"postalCode,omitempty"` // PostalCode of the location
	State              string      `json:"state,omitempty"`      // State or province of the location,
	Country            string      `json:"country"`              // Country alpha-3 code for the country
	Coordinates        GeoLocation `json:"coordinates"`          // Coordinates of the location
	EvseId             string      `json:"evseId"`               // EvseId identifies the EVSE within the CPOs platform
	Evse               string      `json:"evse"`                 // Evse following specification for EVSE ID from "eMI3 standard version V1.0". Can be reused
	ConnectorId        string      `json:"connectorId"`          // ConnectorId Identifier of the connector within the EVSE
	ConnectorStandard  string      `json:"connectorStandard"`    // ConnectorStandard of the installed connector
	ConnectorFormat    string      `json:"connectorFormat"`      // ConnectorFormat socket/cable
	ConnectorPowerType string      `json:"connectorPowerType"`   // ConnectorPowerType power type of connector
}

type SignedData struct {
	EncodingMethod        string         `json:"encodingMethod"`                  // EncodingMethod name of the encoding used in the SignedData field
	EncodingMethodVersion *int           `json:"encodingMethodVersion,omitempty"` // EncodingMethodVersion version of the EncodingMethod
	PublicKey             string         `json:"publicKey,omitempty"`             // PublicKey used to sign the data, base64 encoded
	SignedValues          []*SignedValue `json:"signedValues"`                    // SignedValues one or more signed values
	Url                   string         `json:"url,omitempty"`                   // Url that can be shown to an EV driver
}

type CdrDetails struct {
	StartDateTime            time.Time         `json:"startDateTime"`                  // StartDateTime timestamp of the charging cdr
	EndDateTime              time.Time         `json:"endDateTime"`                    // EndDateTime timestamp when the cdr was completed/finished
	SessionId                string            `json:"sessionId"`                      // SessionId unique ID of the Session for which this CDR is sent
	CdrToken                 *CdrToken         `json:"cdrToken"`                       // CdrToken token used to start this charging cdr
	AuthMethod               string            `json:"authMethod"`                     // AuthMethod method used for authentication
	AuthRef                  string            `json:"authRef,omitempty"`              // AuthRef reference to the authorization given by the eMSP
	CdrLocation              CdrLocation       `json:"cdrLocation"`                    // CdrLocation location where the charging cdr took place
	MeterId                  string            `json:"meterId,omitempty"`              // MeterId identification of the Meter inside the Charge Point
	Currency                 string            `json:"currency"`                       // Currency of the CDR in ISO 4217 Code
	Tariffs                  []*Tariff         `json:"tariffs,omitempty"`              // Tariffs list of relevant Tariff Elements
	ChargingPeriods          []*ChargingPeriod `json:"chargingPeriods,omitempty"`      // ChargingPeriods  list of Charging Periods that make up this charging cdr
	SignedData               *SignedData       `json:"signedData,omitempty"`           // SignedData that belongs to this charging Cdr
	TotalCost                Price             `json:"totalCost"`                      // TotalCost total sum of all the costs of this transaction in the specified currency
	TotalFixedCost           *Price            `json:"totalFixedCost,omitempty"`       // TotalFixedCost total sum of all the fixed costs in the specified currency
	TotalEnergy              float64           `json:"totalEnergy"`                    // TotalEnergy charged, in kWh
	TotalEnergyCost          *Price            `json:"totalEnergyCost,omitempty"`      // TotalEnergyCost total sum of all the cost of all the energy used, in the specified currency
	TotalTime                float64           `json:"totalTime"`                      // TotalTime total duration of the charging cdr in hours
	TotalTimeCost            *Price            `json:"totalTimeCost,omitempty"`        // TotalTimeCost total sum of all the cost related to duration of charging during this transaction, in the specified currency
	TotalParkingTime         *float64          `json:"totalParkingTime,omitempty"`     // TotalParkingTime total duration of the charging cdr where the EV was not charging in hours
	TotalParkingCost         *Price            `json:"totalParkingCost,omitempty"`     // TotalParkingCost total sum of all the cost related to parking of this transaction in the specified currency
	TotalReservationCost     *Price            `json:"totalReservationCost,omitempty"` // TotalReservationCost total sum of all the cost related to a reservation of a Charge Point in the specified currency
	Remark                   string            `json:"remark,omitempty"`               // Remark can be used to provide additional human readable information
	InvoiceReferenceId       string            `json:"invoiceReferenceId,omitempty"`   // InvoiceReferenceId  can be used to reference an invoice
	Credit                   bool              `json:"credit"`                         // Credit when set to true, this is a Credit CDR, and the field credit_reference_id needs to be set as wel
	CreditReferenceId        string            `json:"creditReferenceId,omitempty"`    // CreditReferenceId to be set for a Credit CDR
	HomeChargingCompensation bool              `json:"homeChargingCompensation"`       // HomeChargingCompensation when set to true, this CDR is for a charging cdr using the home charge

}

type Cdr struct {
	OcpiItem
	Id      string     `json:"id"`      // Id uniquely identifies the cdr
	Details CdrDetails `json:"details"` // Details cdr details
}

type CdrSearchCriteria struct {
	PageRequest
	ExtId        *PartyExtId // ExtId by party ext ID
	RefId        string      // RefId by ref id
	IncPlatforms []string    // IncPlatforms includes platform Ids
	ExcPlatforms []string    // ExcPlatforms exclude platform Ids
	Ids          []string    // Ids by list Ids
}

type CdrSearchResponse struct {
	PageResponse
	Items []*Cdr
}

type CdrService interface {
	// PutCdr creates or updates cdr
	PutCdr(ctx context.Context, sess *Cdr) (*Cdr, error)
	// GetCdr retrieves cdr by ID
	GetCdr(ctx context.Context, sessId string) (*Cdr, error)
	// DeleteCdrsByExtId deletes all cdrs by party ext id
	DeleteCdrsByExtId(ctx context.Context, extId PartyExtId) error
	// SearchCdrs searches cdrs
	SearchCdrs(ctx context.Context, cr *CdrSearchCriteria) (*CdrSearchResponse, error)
}

type CdrStorage interface {
	// MergeCdr creates or updates cdr
	MergeCdr(ctx context.Context, sess *Cdr) error
	// UpdateCdr updates cdr
	UpdateCdr(ctx context.Context, sess *Cdr) error
	// GetCdr retrieves cdr by ID
	GetCdr(ctx context.Context, sessId string) (*Cdr, error)
	// DeleteCdrsByExtId deletes all cdrs by party ext id
	DeleteCdrsByExtId(ctx context.Context, extId PartyExtId) error
	// SearchCdrs searches cdrs
	SearchCdrs(ctx context.Context, cr *CdrSearchCriteria) (*CdrSearchResponse, error)
}
