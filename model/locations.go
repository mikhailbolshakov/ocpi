package model

import "time"

type OcpiPublishTokenType struct {
	Uid          string `json:"uid,omitempty"`           // Uid unique ID by which this Token can be identified
	Type         string `json:"type,omitempty"`          // Type of the token
	VisualNumber string `json:"visual_number,omitempty"` // VisualNumber readable number/identification as printed on the Token (RFID card)
	Issuer       string `json:"issuer,omitempty"`        // Issuer issuing company
	GroupId      string `json:"group_id,omitempty"`      // GroupId can be used to make two or more tokens work as one
}

type OcpiGeoLocation struct {
	Latitude  string `json:"latitude"`  // Latitude of the location
	Longitude string `json:"longitude"` // Longitude of the location
}

type OcpiAdditionalGeoLocation struct {
	OcpiGeoLocation
	Name *OcpiDisplayText `json:"name"` // Name of the point in local language
}

type OcpiEnergySource struct {
	Source     string `json:"source"`     // Source type of energy source
	Percentage int    `json:"percentage"` // Percentage of this source (0-100) in the mix
}

type OcpiEnvironmentalImpact struct {
	Category string  `json:"category"` // Category environmental impact category
	Amount   float64 `json:"amount"`   // Amount of this portion in g/kWh
}

type OcpiEnergyMix struct {
	IsGreenEnergy     bool                       `json:"is_green_energy"`               // IsGreenEnergy true if 100% from regenerative sources
	EnergySources     []*OcpiEnergySource        `json:"energy_sources,omitempty"`      // EnergySources energy sources of this location’s tariff
	EnvironImpact     []*OcpiEnvironmentalImpact `json:"environ_impact,omitempty"`      // EnvironImpact key-value pairs (enum + percentage) of nuclear waste and CO2 exhaust
	SupplierName      string                     `json:"supplier_name,omitempty"`       // SupplierName of the energy supplier
	EnergyProductName string                     `json:"energy_product_name,omitempty"` // EnergyProductName energy suppliers product/tariff plan
}

type OcpiStatusSchedule struct {
	PeriodBegin time.Time  `json:"period_begin"`         // PeriodBegin begin of scheduled period
	PeriodEnd   *time.Time `json:"period_end,omitempty"` // PeriodEnd end of schedule period
	Status      string     `json:"status"`               // Status value during the scheduled period.
}

type OcpiConnector struct {
	Id                 string    `json:"id"`                             // Id identifier of the Connector
	Standard           string    `json:"standard"`                       // Standard of the installed connector
	Format             string    `json:"format"`                         // Format socket/cable
	PowerType          string    `json:"power_type"`                     // PowerType
	MaxVoltage         float64   `json:"max_voltage"`                    // MaxVoltage maximum voltage in V
	MaxAmperage        float64   `json:"max_amperage"`                   // MaxAmperage maximum amperage in A
	MaxElectricPower   *float64  `json:"max_electric_power,omitempty"`   // MaxElectricPower maximum electric power in W
	TariffIds          []string  `json:"tariff_ids,omitempty"`           // TariffIds charging tariffs
	TermsAndConditions string    `json:"terms_and_conditions,omitempty"` // TermsAndConditions url of operator’s terms and conditions
	LastUpdated        time.Time `json:"last_updated"`                   // LastUpdated when updated or created
}

type OcpiEvse struct {
	Uid                 string                `json:"uid"`                            // Uid identifies the EVSE within the CPOs platform
	EvseId              string                `json:"evse_id,omitempty"`              // EvseId following specification for EVSE ID from "eMI3 standard version V1.0". Can be reused
	Status              string                `json:"status"`                         // Status current status of the EVSE
	StatusSchedule      []*OcpiStatusSchedule `json:"status_schedule,omitempty"`      // StatusSchedule indicates a planned status update of the EVSE
	Capabilities        []string              `json:"capabilities,omitempty"`         // Capabilities list of functionalities that the EVSE is capable of
	Connectors          []*OcpiConnector      `json:"connectors,omitempty"`           // Connectors list of available connectors on the EVSE
	FloorLevel          string                `json:"floor_level,omitempty"`          // FloorLevel level on which the Charge Point is located (in garage buildings)
	Coordinates         *OcpiGeoLocation      `json:"coordinates,omitempty"`          // Coordinates of the EVSE
	PhysicalReference   string                `json:"physical_reference,omitempty"`   // PhysicalReference number/string printed on the outside of the EVSE
	Directions          []*OcpiDisplayText    `json:"directions,omitempty"`           // Directions human-readable directions
	ParkingRestrictions []string              `json:"parking_restrictions,omitempty"` // ParkingRestrictions restrictions that apply to the parking spot
	Images              []*OcpiImage          `json:"images,omitempty"`               // Images related to the EVSE
	LastUpdated         time.Time             `json:"last_updated"`                   // LastUpdated when updated or created
}

type OcpiLocation struct {
	OcpiPartyId
	Id                 string                       `json:"id"`                             // Id uniquely identifies the location within the CPOs platform
	Publish            *bool                        `json:"publish"`                        // Publish if a Location may be published
	PublishAllowedTo   []*OcpiPublishTokenType      `json:"publish_allowed_to,omitempty"`   // PublishAllowedTo the list are allowed to be shown this location
	Name               string                       `json:"name,omitempty"`                 // Name of the location
	Address            string                       `json:"address"`                        // Address street/block name and house number
	City               string                       `json:"city"`                           // City or town
	PostalCode         string                       `json:"postal_code,omitempty"`          // PostalCode of the location
	State              string                       `json:"state,omitempty"`                // State or province of the location,
	Country            string                       `json:"country"`                        // Country alpha-3 code for the country
	Coordinates        OcpiGeoLocation              `json:"coordinates"`                    // Coordinates of the location
	RelatedLocations   []*OcpiAdditionalGeoLocation `json:"related_locations,omitempty"`    // RelatedLocations related points relevant to the user
	ParkingType        string                       `json:"parking_type,omitempty"`         // ParkingType type of parking at the charge point location
	Evses              []*OcpiEvse                  `json:"evses,omitempty"`                // Evses list of EVSEs that belong to this Location
	Directions         []*OcpiDisplayText           `json:"directions,omitempty"`           // Directions human-readable directions on how to reach the location
	Operator           *OcpiBusinessDetails         `json:"operator,omitempty"`             // Operator of the location
	SubOperator        *OcpiBusinessDetails         `json:"sub_operator,omitempty"`         // SubOperator of the location
	Owner              *OcpiBusinessDetails         `json:"owner,omitempty"`                // Owner of the location
	Facilities         []string                     `json:"facilities,omitempty"`           // Facilities this charging location directly belongs
	TimeZone           string                       `json:"time_zone"`                      // TimeZone IANA tzdata’s TZ-values
	OpeningTimes       *OcpiHours                   `json:"opening_times,omitempty"`        // OpeningTimes times when the EVSEs at the location can be accessed
	ChargingWhenClosed *bool                        `json:"charging_when_closed,omitempty"` // ChargingWhenClosed if the EVSEs are still charging outside the opening hours of the location
	Images             []*OcpiImage                 `json:"images,omitempty"`               // Images links to images related to the location
	EnergyMix          *OcpiEnergyMix               `json:"energy_mix,omitempty"`           // EnergyMix energy supplied at this location
	LastUpdated        time.Time                    `json:"last_updated"`                   // LastUpdated when updated or created
}

type OcpiLocationsResponse struct {
	OcpiResponse
	Data []*OcpiLocation `json:"data"`
}
