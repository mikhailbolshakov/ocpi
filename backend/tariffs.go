package backend

import "time"

const (
	TariffTypeAdHocPay     = "AD_HOC_PAYMENT"
	TariffTypeProfileCheap = "PROFILE_CHEAP"
	TariffTypeProfileFast  = "PROFILE_FAST"
	TariffTypeProfileGreen = "PROFILE_GREEN"
	TariffTypeReg          = "REGULAR"

	TariffDimEnergy      = "ENERGY"
	TariffDimFlat        = "FLAT"
	TariffDimParkingType = "PARKING_TIME"
	TariffDimTime        = "TIME"

	DayMon = "MONDAY"
	DayTue = "TUESDAY"
	DayWed = "WEDNESDAY"
	DayThu = "THURSDAY"
	DayFri = "FRIDAY"
	DaySat = "SATURDAY"
	DaySun = "SUNDAY"

	Reservation    = "RESERVATION"
	ReservationExp = "RESERVATION_EXPIRES"
)

type Price struct {
	ExclVat float64  `json:"exclVat"`           // ExclVat price/Cost excluding VAT
	InclVat *float64 `json:"inclVat,omitempty"` // InclVat Price/Cost including VAT
}

type TariffRestrictions struct {
	StartTime   string     `json:"startTime,omitempty"`   // StartTime time of day in local time
	EndTime     string     `json:"endTime,omitempty"`     // EndTime time of day in local time
	StartDate   *time.Time `json:"startDate,omitempty"`   // StartDate in local time (2015-12-24)
	EndDate     *time.Time `json:"endDate,omitempty"`     // EndDate in local time (2015-12-24)
	MinKwh      *float64   `json:"minKwh,omitempty"`      // MinKwh minimum consumed energy
	MaxKwh      *float64   `json:"maxKwh,omitempty"`      // MaxKwh maximum consumed energy
	MinCurrent  *float64   `json:"minCurrent,omitempty"`  // MinCurrent sum of the minimum current (in Amperes)
	MaxCurrent  *float64   `json:"maxCurrent,omitempty"`  // MaxCurrent sum of the maximum current (in Amperes)
	MinPower    *float64   `json:"minPower,omitempty"`    // MinPower minimum power in kW
	MaxPower    *float64   `json:"maxPower,omitempty"`    // MaxPower maximum power in kW
	MinDuration *float64   `json:"minDuration,omitempty"` // MinDuration minimum duration in second
	MaxDuration *float64   `json:"maxDuration,omitempty"` // MaxDuration maximum duration in second
	DayOfWeek   []string   `json:"dayOfWeek,omitempty"`   // DayOfWeek which day(s) of the week this TariffElement is active
	Reservation string     `json:"reservation,omitempty"` // Reservation populated if the element describes reservation costs
}

type PriceComponent struct {
	Type     string   `json:"type"`     // Type of tariff dimension
	Price    float64  `json:"price"`    // Price per unit (excl. VAT) for this tariff dimension
	Vat      *float64 `json:"vat"`      // Vat percentage for this tariff dimension
	StepSize int      `json:"stepSize"` // StepSize minimum amount to be billed
}

type TariffElement struct {
	PriceComponents []*PriceComponent   `json:"priceComponents"` // PriceComponents list of price components
	Restrictions    *TariffRestrictions `json:"restrictions"`    // Restrictions describe the applicability of a tariff
}

type Tariff struct {
	Id            string           `json:"id"`                      // Id uniquely identifies the tariff
	Currency      string           `json:"currency"`                // Currency ISO-4217 code
	Type          string           `json:"type,omitempty"`          // Type of the tariff
	TariffAltText []*DisplayText   `json:"tariffAltText,omitempty"` // TariffAltText list of multi-language alternative tariff info texts
	TariffAltUrl  string           `json:"tariffAltUrl,omitempty"`  // TariffAltUrl web page that contains an explanation of the tariff
	MinPrice      *Price           `json:"minPrice,omitempty"`      // MinPrice minimum possible price
	MaxPrice      *Price           `json:"maxPrice,omitempty"`      // MaxPrice maximum possible price
	Elements      []*TariffElement `json:"elements"`                // Elements list of tariff elements
	StartDateTime *time.Time       `json:"startDateTime,omitempty"` // StartDateTime when this tariff becomes active in UTC
	EndDateTime   *time.Time       `json:"endDateTime,omitempty"`   // EndDateTime when this tariff no longer valid in UTC
	EnergyMix     *EnergyMix       `json:"energyMix,omitempty"`     // EnergyMix details on the energy supplied with this tariff
	LastUpdated   time.Time        `json:"lastUpdated"`             // LastUpdated when this Tariff was last updated
	PlatformId    string           `json:"platformId"`              // PlatformId rel to platform
	RefId         string           `json:"refId"`                   // RefId any external relation
	PartyId       string           `json:"partyId,omitempty"`       // PartyId should be unique within country
	CountryCode   string           `json:"countryCode,omitempty"`   // CountryCode alfa-2 code
}

type TariffSearchResponse struct {
	PageInfo *PageResponse `json:"pageInfo,omitempty"`
	Items    []*Tariff     `json:"items,omitempty"`
}
