package model

import "time"

type OcpiCdrLocation struct {
	Id                 string          `json:"id"`                    // Id Uniquely identifies the location within the CPO’s platform
	Name               string          `json:"name,omitempty"`        // Name of the location
	Address            string          `json:"address"`               // Address street/block name and house number if available
	City               string          `json:"city"`                  // City or town
	PostalCode         string          `json:"postal_code,omitempty"` // PostalCode of the location
	State              string          `json:"state,omitempty"`       // State or province of the location,
	Country            string          `json:"country"`               // Country alpha-3 code for the country
	Coordinates        OcpiGeoLocation `json:"coordinates"`           // Coordinates of the location
	EvseUid            string          `json:"evse_uid"`              // EvseUid identifies the EVSE within the CPOs platform
	EvseId             string          `json:"evse_id"`               // EvseId following specification for EVSE ID from "eMI3 standard version V1.0". Can be reused
	ConnectorId        string          `json:"connector_id"`          // ConnectorId Identifier of the connector within the EVSE
	ConnectorStandard  string          `json:"connector_standard"`    // ConnectorStandard of the installed connector
	ConnectorFormat    string          `json:"connector_format"`      // ConnectorFormat socket/cable
	ConnectorPowerType string          `json:"connector_power_type"`  // ConnectorPowerType power type of connector
}

type OcpiSignedData struct {
	EncodingMethod        string             `json:"encoding_method"`                   // EncodingMethod name of the encoding used in the SignedData field
	EncodingMethodVersion *int               `json:"encoding_method_version,omitempty"` // EncodingMethodVersion version of the EncodingMethod
	PublicKey             string             `json:"public_key,omitempty"`              // PublicKey used to sign the data, base64 encoded
	SignedValues          []*OcpiSignedValue `json:"signed_values"`                     // SignedValues one or more signed values
	Url                   string             `json:"url,omitempty"`                     // Url that can be shown to an EV driver
}

type OcpiCdr struct {
	OcpiPartyId
	Id                       string                `json:"id"`                                // Id unique id that identifies the cdr
	StartDateTime            time.Time             `json:"start_date_time"`                   // StartDateTime timestamp of the charging session
	EndDateTime              time.Time             `json:"end_date_time"`                     // EndDateTime timestamp when the session was completed/finished
	SessionId                string                `json:"session_id"`                        // SessionId unique ID of the Session for which this CDR is sent
	CdrToken                 *OcpiCdrToken         `json:"cdr_token"`                         // CdrToken token used to start this charging session
	AuthMethod               string                `json:"auth_method"`                       // AuthMethod method used for authentication
	AuthRef                  string                `json:"authorization_reference,omitempty"` // AuthRef reference to the authorization given by the eMSP
	CdrLocation              OcpiCdrLocation       `json:"cdr_location"`                      // CdrLocation location where the charging session took place
	MeterId                  string                `json:"meter_id,omitempty"`                // MeterId identification of the Meter inside the Charge Point
	Currency                 string                `json:"currency"`                          // Currency of the CDR in ISO 4217 Code
	Tariffs                  []*OcpiTariff         `json:"tariffs,omitempty"`                 // Tariffs list of relevant Tariff Elements
	ChargingPeriods          []*OcpiChargingPeriod `json:"charging_periods,omitempty"`        // ChargingPeriods  list of Charging Periods that make up this charging session
	SignedData               *OcpiSignedData       `json:"signed_data,omitempty"`             // SignedData that belongs to this charging Session
	TotalCost                OcpiPrice             `json:"total_cost"`                        // TotalCost total sum of all the costs of this transaction in the specified currency
	TotalFixedCost           *OcpiPrice            `json:"total_fixed_cost,omitempty"`        // TotalFixedCost total sum of all the fixed costs in the specified currency
	TotalEnergy              float64               `json:"total_energy"`                      // TotalEnergy charged, in kWh
	TotalEnergyCost          *OcpiPrice            `json:"total_energy_cost,omitempty"`       // TotalEnergyCost total sum of all the cost of all the energy used, in the specified currency
	TotalTime                float64               `json:"total_time"`                        // TotalTime total duration of the charging session in hours
	TotalTimeCost            *OcpiPrice            `json:"total_time_cost,omitempty"`         // TotalTimeCost total sum of all the cost related to duration of charging during this transaction, in the specified currency
	TotalParkingTime         *float64              `json:"total_parking_time,omitempty"`      // TotalParkingTime total duration of the charging session where the EV was not charging in hours
	TotalParkingCost         *OcpiPrice            `json:"total_parking_cost,omitempty"`      // TotalParkingCost total sum of all the cost related to parking of this transaction in the specified currency
	TotalReservationCost     *OcpiPrice            `json:"total_reservation_cost,omitempty"`  // TotalReservationCost total sum of all the cost related to a reservation of a Charge Point in the specified currency
	Remark                   string                `json:"remark,omitempty"`                  // Remark can be used to provide additional human readable information
	InvoiceReferenceId       string                `json:"invoice_reference_id,omitempty"`    // InvoiceReferenceId  can be used to reference an invoice
	Credit                   bool                  `json:"credit"`                            // Credit when set to true, this is a Credit CDR, and the field credit_reference_id needs to be set as wel
	CreditReferenceId        string                `json:"credit_reference_id,omitempty"`     // CreditReferenceId to be set for a Credit CDR
	HomeChargingCompensation bool                  `json:"home_charging_compensation"`        // HomeChargingCompensation when set to true, this CDR is for a charging session using the home charge
	LastUpdated              time.Time             `json:"last_updated"`                      // LastUpdated when updated or created
}

type OcpiCdrsResponse struct {
	OcpiResponse
	Data []*OcpiCdr `json:"data"`
}
