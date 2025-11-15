package model

import "time"

type OcpiCdrToken struct {
	OcpiPartyId
	Id         string `json:"uid"`         // Id unique ID by which this Token can be identified
	Type       string `json:"type"`        // Type token type
	ContractId string `json:"contract_id"` // ContractId uniquely identifies the EV driver contract token
}

type OcpiCdrDimension struct {
	Type   string  `json:"type"`   // Type of CDR dimension
	Volume float64 `json:"volume"` // Volume of the dimension
}

type OcpiSignedValue struct {
	Nature     string `json:"nature"`      // Nature of the value
	PlainData  string `json:"plain_data"`  // PlainData un-encoded string of data
	SignedData string `json:"signed_data"` // SignedData blob of signed data, base64 encoded
}

type OcpiChargingPeriod struct {
	StartDateTime         time.Time           `json:"start_date_time"`                   // StartDateTime start timestamp of the charging period
	Dimensions            []*OcpiCdrDimension `json:"dimensions"`                        // Dimensions list of relevant values for this charging period
	TariffId              string              `json:"tariff_id"`                         // TariffId unique identifier of the Tariff that is relevant for this Charging Period
	EncodingMethod        string              `json:"encoding_method"`                   // EncodingMethod name of the encoding used in the SignedData field
	EncodingMethodVersion *int                `json:"encoding_method_version,omitempty"` // EncodingMethodVersion version of the EncodingMethod
	PublicKey             string              `json:"public_key"`                        // PublicKey used to sign the data, base64 encoded
	SignedValues          []*OcpiSignedValue  `json:"signed_values"`                     // SignedValues one or more signed values
	Url                   string              `json:"url,omitempty"`                     // Url that can be shown to an EV driver
}

type OcpiSession struct {
	OcpiPartyId
	Id              string                `json:"id"`                                // Id unique id that identifies the charging session
	StartDateTime   *time.Time            `json:"start_date_time"`                   // StartDateTime timestamp when the session became ACTIVE
	EndDateTime     *time.Time            `json:"end_date_time"`                     // EndDateTime timestamp when the session was completed/finished
	Kwh             *float64              `json:"kwh"`                               // Kwh how many kWh were charged
	CdrToken        *OcpiCdrToken         `json:"cdr_token"`                         // CdrToken token used to start this charging session
	AuthMethod      string                `json:"auth_method"`                       // AuthMethod method used for authentication
	AuthRef         string                `json:"authorization_reference,omitempty"` // AuthRef reference to the authorization given by the eMSP
	LocationId      string                `json:"location_id"`                       // LocationId id of the location obj
	EvseId          string                `json:"evse_uid"`                          // EvseId id of the evse obj
	ConnectorId     string                `json:"connector_id"`                      // ConnectorId id of the connector
	MeterId         string                `json:"meter_id,omitempty"`                // MeterId id of the kWh meter
	Currency        string                `json:"currency"`                          // Currency ISO-4217 currency code
	ChargingPeriods []*OcpiChargingPeriod `json:"charging_periods,omitempty"`        // ChargingPeriods  list of Charging Periods that can be used to calculate and verify the total cost
	TotalCost       *OcpiPrice            `json:"total_cost,omitempty"`              // TotalCost total cost of the session in the specified currency
	Status          string                `json:"status"`                            // Status session status
	LastUpdated     time.Time             `json:"last_updated"`                      // LastUpdated when updated or created
}

type OcpiSessionsResponse struct {
	OcpiResponse
	Data []*OcpiSession `json:"data"`
}
