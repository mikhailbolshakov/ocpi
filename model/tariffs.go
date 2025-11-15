package model

import "time"

type OcpiPrice struct {
	ExclVat float64  `json:"excl_vat"`           // ExclVat price/Cost excluding VAT
	InclVat *float64 `json:"incl_vat,omitempty"` // InclVat Price/Cost including VAT
}

type OcpiTariffRestrictions struct {
	StartTime   string     `json:"start_time,omitempty"`   // StartTime time of day in local time
	EndTime     string     `json:"end_time,omitempty"`     // EndTime time of day in local time
	StartDate   *time.Time `json:"start_date,omitempty"`   // StartDate in local time (2015-12-24)
	EndDate     *time.Time `json:"end_date,omitempty"`     // EndDate in local time (2015-12-24)
	MinKwh      *float64   `json:"min_kwh,omitempty"`      // MinKwh minimum consumed energy
	MaxKwh      *float64   `json:"max_kwh,omitempty"`      // MaxKwh maximum consumed energy
	MinCurrent  *float64   `json:"min_current,omitempty"`  // MinCurrent sum of the minimum current (in Amperes)
	MaxCurrent  *float64   `json:"max_current,omitempty"`  // MaxCurrent sum of the maximum current (in Amperes)
	MinPower    *float64   `json:"min_power,omitempty"`    // MinPower minimum power in kW
	MaxPower    *float64   `json:"max_power,omitempty"`    // MaxPower maximum power in kW
	MinDuration *float64   `json:"min_duration,omitempty"` // MinDuration minimum duration in second
	MaxDuration *float64   `json:"max_duration,omitempty"` // MaxDuration maximum duration in second
	DayOfWeek   []string   `json:"day_of_week,omitempty"`  // DayOfWeek which day(s) of the week this TariffElement is active
	Reservation string     `json:"reservation,omitempty"`  // Reservation populated if the element describes reservation costs
}

type OcpiPriceComponent struct {
	Type     string   `json:"type"`      // Type of tariff dimension
	Price    float64  `json:"price"`     // Price per unit (excl. VAT) for this tariff dimension
	Vat      *float64 `json:"vat"`       // Vat percentage for this tariff dimension
	StepSize int      `json:"step_size"` // StepSize minimum amount to be billed
}

type OcpiTariffElement struct {
	PriceComponents []*OcpiPriceComponent   `json:"price_components"` // PriceComponents list of price components
	Restrictions    *OcpiTariffRestrictions `json:"restrictions"`     // Restrictions describe the applicability of a tariff
}

type OcpiTariff struct {
	OcpiPartyId
	Id            string               `json:"id"`                        // Id uniquely identifies the tariff
	Currency      string               `json:"currency"`                  // Currency ISO-4217 code
	Type          string               `json:"type,omitempty"`            // Type of the tariff
	TariffAltText []*OcpiDisplayText   `json:"tariff_alt_text,omitempty"` // TariffAltText list of multi-language alternative tariff info texts
	TariffAltUrl  string               `json:"tariff_alt_url,omitempty"`  // TariffAltUrl web page that contains an explanation of the tariff
	MinPrice      *OcpiPrice           `json:"min_price,omitempty"`       // MinPrice minimum possible price
	MaxPrice      *OcpiPrice           `json:"max_price,omitempty"`       // MaxPrice maximum possible price
	Elements      []*OcpiTariffElement `json:"elements"`                  // Elements list of tariff elements
	StartDateTime *time.Time           `json:"start_date_time,omitempty"` // StartDateTime when this tariff becomes active in UTC
	EndDateTime   *time.Time           `json:"end_date_time,omitempty"`   // EndDateTime when this tariff no longer valid in UTC
	EnergyMix     *OcpiEnergyMix       `json:"energy_mix,omitempty"`      // EnergyMix details on the energy supplied with this tariff
	LastUpdated   time.Time            `json:"last_updated"`              // LastUpdated when this Tariff was last updated
}

type OcpiTariffsResponse struct {
	OcpiResponse
	Data []*OcpiTariff `json:"data"`
}
