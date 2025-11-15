package domain

import (
	"context"
	"time"
)

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
	ExclVat float64  `json:"exclVat"`            // ExclVat price/Cost excluding VAT
	InclVat *float64 `json:"incl_vat,omitempty"` // InclVat Price/Cost including VAT
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

type TariffDetails struct {
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
}

type Tariff struct {
	OcpiItem
	Id      string        `json:"id"`      // Id uniquely identifies the tariff
	Details TariffDetails `json:"details"` // Details tariff details
}

type TariffSearchCriteria struct {
	PageRequest
	ExtId        *PartyExtId // ExtId by party ext ID
	RefId        string      // RefId by ref id
	IncPlatforms []string    // IncPlatforms includes platform Ids
	ExcPlatforms []string    // ExcPlatforms exclude platform Ids
	Ids          []string    // Ids by list of Ids
}

type TariffSearchResponse struct {
	PageResponse
	Items []*Tariff
}

type TariffService interface {
	// PutTariff creates or updates tariff
	PutTariff(ctx context.Context, trf *Tariff) (*Tariff, error)
	// MergeTariff merges tariff
	MergeTariff(ctx context.Context, trf *Tariff) (*Tariff, error)
	// GetTariff retrieves tariff by ID
	GetTariff(ctx context.Context, trfId string) (*Tariff, error)
	// DeleteTariffsByExtId deletes all tariffs by party ext id
	DeleteTariffsByExtId(ctx context.Context, extId PartyExtId) error
	// SearchTariffs searches tariffs
	SearchTariffs(ctx context.Context, cr *TariffSearchCriteria) (*TariffSearchResponse, error)
	// Validate validates tariffs
	Validate(ctx context.Context, trf *Tariff) error
}

type TariffStorage interface {
	// MergeTariff creates or updates tariff
	MergeTariff(ctx context.Context, trf *Tariff) error
	// UpdateTariff updates tariff
	UpdateTariff(ctx context.Context, trf *Tariff) error
	// GetTariff retrieves tariff by ID
	GetTariff(ctx context.Context, trfId string) (*Tariff, error)
	// DeleteTariffsByExtId deletes all tariffs by party ext id
	DeleteTariffsByExtId(ctx context.Context, extId PartyExtId) error
	// SearchTariffs searches tariffs
	SearchTariffs(ctx context.Context, cr *TariffSearchCriteria) (*TariffSearchResponse, error)
}
