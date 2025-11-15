package backend

import "time"

type Cdr struct {
	Id                       string            `json:"id"`                             // Id uniquely identifies the session
	StartDateTime            time.Time         `json:"startDateTime"`                  // StartDateTime timestamp of the charging cdr
	EndDateTime              time.Time         `json:"endDateTime"`                    // EndDateTime timestamp when the cdr was completed/finished
	SessionId                string            `json:"sessionId"`                      // SessionId unique ID of the Session for which this CDR is sent
	MeterId                  string            `json:"meterId,omitempty"`              // MeterId identification of the Meter inside the Charge Point
	Currency                 string            `json:"currency"`                       // Currency of the CDR in ISO 4217 Code
	Tariffs                  []*Tariff         `json:"tariffs,omitempty"`              // Tariffs list of relevant Tariff Elements
	ChargingPeriods          []*ChargingPeriod `json:"chargingPeriods,omitempty"`      // ChargingPeriods  list of Charging Periods that make up this charging cdr
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
	LastUpdated              time.Time         `json:"lastUpdated"`                    // LastUpdated when this Tariff was last updated
	PlatformId               string            `json:"platformId"`                     // PlatformId rel to platform
	RefId                    string            `json:"refId"`                          // RefId any external relation
	PartyId                  string            `json:"partyId,omitempty"`              // PartyId should be unique within country
	CountryCode              string            `json:"countryCode,omitempty"`          // CountryCode alfa-2 code
}

type CdrSearchResponse struct {
	PageInfo *PageResponse `json:"pageInfo,omitempty"`
	Items    []*Cdr        `json:"items,omitempty"`
}
