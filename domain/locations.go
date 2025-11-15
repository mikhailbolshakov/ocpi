package domain

import (
	"context"
	"time"
)

const (
	EvseStatusAvailable   = "AVAILABLE"
	EvseStatusBlocked     = "BLOCKED"
	EvseStatusCharging    = "CHARGING"
	EvseStatusInOperative = "INOPERATIVE"
	EvseStatusOutOfOrder  = "OUTOFORDER"
	EvseStatusPlanned     = "PLANNED"
	EvseStatusRemoved     = "REMOVED"
	EvseStatusReserved    = "RESERVED"
	EvseStatusUnknown     = "UNKNOWN"

	ParkingTypeMotorway          = "ALONG_MOTORWAY"
	ParkingTypeGarage            = "PARKING_GARAGE"
	ParkingTypeParkingLot        = "PARKING_LOT"
	ParkingTypeDriveway          = "ON_DRIVEWAY"
	ParkingTypeStreet            = "ON_STREET"
	ParkingTypeUndergroundGarage = "UNDERGROUND_GARAGE"

	ImageCategoryCharger  = "CHARGER"
	ImageCategoryEntrance = "ENTRANCE"
	ImageCategoryLocation = "LOCATION"
	ImageCategoryNetwork  = "NETWORK"
	ImageCategoryOperator = "OPERATOR"
	ImageCategoryOther    = "OTHER"
	ImageCategoryOwner    = "OWNER"

	FacilityHotel          = "HOTEL"
	FacilityRestaurant     = "RESTAURANT"
	FacilityCafe           = "CAFE"
	FacilityMall           = "MALL"
	FacilitySupermarket    = "SUPERMARKET"
	FacilitySport          = "SPORT"
	FacilityRecreationArea = "RECREATION_AREA"
	FacilityNature         = "NATURE"
	FacilityMuseum         = "MUSEUM"
	FacilityBikeSharing    = "BIKE_SHARING"
	FacilityBusStop        = "BUS_STOP"
	FacilityTaxiStand      = "TAXI_STAND"
	FacilityTramStop       = "TRAM_STOP"
	FacilityMetroStation   = "METRO_STATION"
	FacilityTrainStation   = "TRAIN_STATION"
	FacilityAirport        = "AIRPORT"
	FacilityParkingLot     = "PARKING_LOT"
	FacilityCarpoolParking = "CARPOOL_PARKING"
	FacilityFuelStation    = "FUEL_STATION"
	FacilityWifi           = "WIFI"

	CapabilityChargingProfile    = "CHARGING_PROFILE_CAPABLE"
	CapabilityChargingPreference = "CHARGING_PREFERENCES_CAPABLE"
	CapabilityChipCard           = "CHIP_CARD_SUPPORT"
	CapabilityContactlessCard    = "CONTACTLESS_CARD_SUPPORT"
	CapabilityCreditCard         = "CREDIT_CARD_PAYABLE"
	CapabilityDebitCard          = "DEBIT_CARD_PAYABLE"
	CapabilityPedTerminal        = "PED_TERMINAL"
	CapabilityRemoteStartStop    = "REMOTE_START_STOP_CAPABLE"
	CapabilityReservable         = "RESERVABLE"
	CapabilityRfid               = "RFID_READER"
	CapabilityStartSessionConReq = "START_SESSION_CONNECTOR_REQUIRED"
	CapabilityTokenGroup         = "TOKEN_GROUP_CAPABLE"
	CapabilityUnlock             = "UNLOCK_CAPABLE"

	ParkingRestrictionEvOnly   = "EV_ONLY"
	ParkingRestrictionPlugged  = "PLUGGED"
	ParkingRestrictionDisable  = "DISABLED"
	ParkingRestrictionCustomer = "CUSTOMERS"
	ParkingRestrictionMoto     = "MOTORCYCLES"

	ConnectorTypeChademo            = "CHADEMO"
	ConnectorTypeChaoji             = "CHAOJI"
	ConnectorTypeDomesticA          = "DOMESTIC_A"
	ConnectorTypeDomesticB          = "DOMESTIC_B"
	ConnectorTypeDomesticC          = "DOMESTIC_C"
	ConnectorTypeDomesticD          = "DOMESTIC_D"
	ConnectorTypeDomesticE          = "DOMESTIC_E"
	ConnectorTypeDomesticF          = "DOMESTIC_F"
	ConnectorTypeDomesticG          = "DOMESTIC_G"
	ConnectorTypeDomesticH          = "DOMESTIC_H"
	ConnectorTypeDomesticI          = "DOMESTIC_I"
	ConnectorTypeDomesticJ          = "DOMESTIC_J"
	ConnectorTypeDomesticK          = "DOMESTIC_K"
	ConnectorTypeDomesticL          = "DOMESTIC_L"
	ConnectorTypeDomesticM          = "DOMESTIC_M"
	ConnectorTypeDomesticN          = "DOMESTIC_N"
	ConnectorTypeDomesticO          = "DOMESTIC_O"
	ConnectorTypeGbtAc              = "GBT_AC"
	ConnectorTypeGbtDc              = "GBT_DC"
	ConnectorTypeSingle16           = "IEC_60309_2_single_16"
	ConnectorTypeThree16            = "IEC_60309_2_three_16"
	ConnectorTypeThree32            = "IEC_60309_2_three_32"
	ConnectorTypeThree64            = "IEC_60309_2_three_64"
	ConnectorTypeT1                 = "IEC_62196_T1"
	ConnectorTypeT1Combo            = "IEC_62196_T1_COMBO"
	ConnectorTypeT2                 = "IEC_62196_T2"
	ConnectorTypeT2Combo            = "IEC_62196_T2_COMBO"
	ConnectorTypeT3A                = "IEC_62196_T3A"
	ConnectorTypeT3C                = "IEC_62196_T3C"
	ConnectorTypeNema520            = "NEMA_5_20"
	ConnectorTypeNema630            = "NEMA_6_30"
	ConnectorTypeNema650            = "NEMA_6_50"
	ConnectorTypeNema1030           = "NEMA_10_30"
	ConnectorTypeNema1050           = "NEMA_10_50"
	ConnectorTypeNema1430           = "NEMA_14_30"
	ConnectorTypeNema1450           = "NEMA_14_50"
	ConnectorTypePantographBottomUp = "PANTOGRAPH_BOTTOM_UP"
	ConnectorTypePantographTopDown  = "PANTOGRAPH_TOP_DOWN"
	ConnectorTypeTeslaR             = "TESLA_R"
	ConnectorTypeTeslaS             = "TESLA_S"

	FormatSocket = "SOCKET"
	FormatCable  = "CABLE"

	PowerTypeAc1Phase      = "AC_1_PHASE"
	PowerTypeAc2Phase      = "AC_2_PHASE"
	PowerTypeAc2PhaseSplit = "AC_2_PHASE_SPLIT"
	PowerTypeAc3Phase      = "AC_3_PHASE"
	PowerTypeDc            = "DC"
)

type PublishTokenType struct {
	Uid          string `json:"uid,omitempty"`          // Uid unique ID by which this Token can be identified
	Type         string `json:"type,omitempty"`         // Type of the token
	VisualNumber string `json:"visualNumber,omitempty"` // VisualNumber readable number/identification as printed on the Token (RFID card)
	Issuer       string `json:"issuer,omitempty"`       // Issuer issuing company
	GroupId      string `json:"groupId,omitempty"`      // GroupId can be used to make two or more tokens work as one
}

type GeoLocation struct {
	Latitude  string `json:"latitude"`  // Latitude of the location
	Longitude string `json:"longitude"` // Longitude of the location
}

type AdditionalGeoLocation struct {
	GeoLocation
	Name *DisplayText `json:"name"` // Name of the point in local language
}

type EnergySource struct {
	Source     string `json:"source"`     // Source type of energy source
	Percentage int    `json:"percentage"` // Percentage of this source (0-100) in the mix
}

type EnvironmentalImpact struct {
	Category string  `json:"category"` // Category environmental impact category
	Amount   float64 `json:"amount"`   // Amount of this portion in g/kWh
}

type EnergyMix struct {
	IsGreenEnergy     bool                   `json:"isGreenEnergy"`               // IsGreenEnergy true if 100% from regenerative sources
	EnergySources     []*EnergySource        `json:"energySources,omitempty"`     // EnergySources energy sources of this location’s tariff
	EnvironImpact     []*EnvironmentalImpact `json:"environImpact,omitempty"`     // EnvironImpact key-value pairs (enum + percentage) of nuclear waste and CO2 exhaust
	SupplierName      string                 `json:"supplierName,omitempty"`      // SupplierName of the energy supplier
	EnergyProductName string                 `json:"energyProductName,omitempty"` // EnergyProductName energy suppliers product/tariff plan
}

type StatusSchedule struct {
	PeriodBegin time.Time  `json:"periodBegin"`         // PeriodBegin begin of scheduled period
	PeriodEnd   *time.Time `json:"periodEnd,omitempty"` // PeriodEnd end of schedule period
	Status      string     `json:"status"`              // Status value during the scheduled period.
}

type RegularHours struct {
	Weekday     int    `json:"weekday"`     // Weekday from Monday (1) till Sunday (7)
	PeriodBegin string `json:"periodBegin"` // PeriodBegin in 24h format with leading zeros. Example: "18:15"
	PeriodEnd   string `json:"periodEnd"`   // PeriodEnd in 24h format with leading zeros. Example: "18:15"
}

type ExceptionalPeriod struct {
	PeriodBegin time.Time `json:"periodBegin"` // PeriodBegin begin of the exception. In UTC, time_zone field can be used to convert to local time
	PeriodEnd   time.Time `json:"periodEnd"`   // PeriodEnd end of the exception. In UTC, time_zone field can be used to convert to local time
}

type Hours struct {
	TwentyFourSeven     bool                 `json:"twentyfourseven"`               // TwentyFourSeven true to represent 24 hours a day and 7 days a week
	RegularHours        []*RegularHours      `json:"regularHours,omitempty"`        // RegularHours weekday-based, used if twentyfourseven=false
	ExceptionalOpenings []*ExceptionalPeriod `json:"exceptionalOpenings,omitempty"` // ExceptionalOpenings for specified calendar dates, time-range based
	ExceptionalClosings []*ExceptionalPeriod `json:"exceptionalClosings,omitempty"` // ExceptionalClosings for specified calendar dates, time-range based
}

type ConnectorDetails struct {
	Standard           string   `json:"standard"`                     // Standard of the installed connector
	Format             string   `json:"format"`                       // Format socket/cable
	PowerType          string   `json:"powerType"`                    // PowerType
	MaxVoltage         float64  `json:"maxVoltage"`                   // MaxVoltage maximum voltage in V
	MaxAmperage        float64  `json:"maxAmperage"`                  // MaxAmperage maximum amperage in A
	MaxElectricPower   *float64 `json:"maxElectricPower,omitempty"`   // MaxElectricPower maximum electric power in W
	TariffIds          []string `json:"tariffIds,omitempty"`          // TariffIds charging tariffs
	TermsAndConditions string   `json:"termsAndConditions,omitempty"` // TermsAndConditions url of operator’s terms and conditions
}

type Connector struct {
	OcpiItem
	Id         string           `json:"id"`               // Id identifier of the Connector
	LocationId string           `json:"locationId"`       // LocationId link to location id
	EvseId     string           `json:"evseId,omitempty"` // EvseId following specification for EVSE ID from "eMI3 standard version V1.0". Can be reused
	Details    ConnectorDetails `json:"details"`          // Details connector details
}

type EvseDetails struct {
	EvseId              string            `json:"evseId,omitempty"`              // EvseId following specification for EVSE ID from "eMI3 standard version V1.0". Can be reused
	StatusSchedule      []*StatusSchedule `json:"statusSchedule,omitempty"`      // StatusSchedule indicates a planned status update of the EVSE
	Capabilities        []string          `json:"capabilities,omitempty"`        // Capabilities list of functionalities that the EVSE is capable of
	FloorLevel          string            `json:"floorLevel,omitempty"`          // FloorLevel level on which the Charge Point is located (in garage buildings)
	Coordinates         *GeoLocation      `json:"coordinates,omitempty"`         // Coordinates of the EVSE
	PhysicalReference   string            `json:"physicalReference,omitempty"`   // PhysicalReference number/string printed on the outside of the EVSE
	Directions          []*DisplayText    `json:"directions,omitempty"`          // Directions human-readable directions
	ParkingRestrictions []string          `json:"parkingRestrictions,omitempty"` // ParkingRestrictions restrictions that apply to the parking spot
	Images              []*Image          `json:"images,omitempty"`              // Images related to the EVSE
}

type Evse struct {
	OcpiItem
	Id         string       `json:"id"`                   // Id identifies the EVSE within the CPOs platform
	LocationId string       `json:"locationId"`           // LocationId link to location id
	Status     string       `json:"status"`               // Status current status of the EVSE
	Details    EvseDetails  `json:"details"`              // Details evse details
	Connectors []*Connector `json:"connectors,omitempty"` // Connectors related connectors
}

type LocationDetails struct {
	Publish            *bool                    `json:"publish"`                      // Publish if a Location may be published
	PublishAllowedTo   []*PublishTokenType      `json:"publishAllowedTo,omitempty"`   // PublishAllowedTo the list are allowed to be shown this location
	Name               string                   `json:"name,omitempty"`               // Name of the location
	Address            string                   `json:"address"`                      // Address street/block name and house number
	City               string                   `json:"city"`                         // City or town
	PostalCode         string                   `json:"postalCode,omitempty"`         // PostalCode of the location
	State              string                   `json:"state,omitempty"`              // State or province of the location,
	Country            string                   `json:"country"`                      // Country alpha-3 code for the country
	Coordinates        GeoLocation              `json:"coordinates"`                  // Coordinates of the location
	RelatedLocations   []*AdditionalGeoLocation `json:"relatedLocations,omitempty"`   // RelatedLocations related points relevant to the user
	ParkingType        string                   `json:"parkingType,omitempty"`        // ParkingType type of parking at the charge point location
	Directions         []*DisplayText           `json:"directions,omitempty"`         // Directions human-readable directions on how to reach the location
	Operator           *BusinessDetails         `json:"operator,omitempty"`           // Operator of the location
	SubOperator        *BusinessDetails         `json:"subOperator,omitempty"`        // SubOperator of the location
	Owner              *BusinessDetails         `json:"owner,omitempty"`              // Owner of the location
	Facilities         []string                 `json:"facilities,omitempty"`         // Facilities this charging location directly belongs
	TimeZone           string                   `json:"timeZone"`                     // TimeZone IANA tzdata’s TZ-values
	OpeningTimes       *Hours                   `json:"openingTimes,omitempty"`       // OpeningTimes times when the EVSEs at the location can be accessed
	ChargingWhenClosed *bool                    `json:"chargingWhenClosed,omitempty"` // ChargingWhenClosed if the EVSEs are still charging outside the opening hours of the location
	Images             []*Image                 `json:"images,omitempty"`             // Images links to images related to the location
	EnergyMix          *EnergyMix               `json:"energyMix,omitempty"`          // EnergyMix energy supplied at this location
}

type Location struct {
	OcpiItem
	Id      string          `json:"id"`              // Id uniquely identifies the location within the CPOs platform
	Details LocationDetails `json:"details"`         // Details locations details
	Evses   []*Evse         `json:"evses,omitempty"` // Evses list of evses
}

type LocationSearchCriteria struct {
	PageRequest
	ExtId        *PartyExtId // ExtId by party ext ID
	RefId        string      // RefId by ref id
	IncPlatforms []string    // IncPlatforms includes platform Ids
	ExcPlatforms []string    // ExcPlatforms exclude platform Ids
	Ids          []string    // Ids by list of ids
}

type LocationSearchResponse struct {
	PageResponse
	Items []*Location
}

type EvseSearchCriteria struct {
	PageRequest
	ExtId *PartyExtId // ExtId by party ext ID
	RefId string      // RefId by ref id
}

type EvseSearchResponse struct {
	PageResponse
	Items []*Evse
}

type ConnectorSearchCriteria struct {
	PageRequest
	ExtId *PartyExtId // ExtId by party ext ID
	RefId string      // RefId by ref id
}

type ConnectorSearchResponse struct {
	PageResponse
	Items []*Connector
}

type LocationService interface {
	// PutLocation creates or updates location
	PutLocation(ctx context.Context, loc *Location) (*Location, error)
	// MergeLocation merges location
	MergeLocation(ctx context.Context, loc *Location) (*Location, error)
	// GetLocation retrieves location by ID
	GetLocation(ctx context.Context, locId string, withEvse bool) (*Location, error)
	// SearchLocations searches locations
	SearchLocations(ctx context.Context, cr *LocationSearchCriteria) (*LocationSearchResponse, error)
	// DeleteLocationsByExtId deletes locations (evse + connectors) by party ext id
	DeleteLocationsByExtId(ctx context.Context, extId PartyExtId) error
	// PutEvse creates or updates evse
	PutEvse(ctx context.Context, evse *Evse) (*Evse, error)
	// MergeEvse merges evse
	MergeEvse(ctx context.Context, evse *Evse) (*Evse, error)
	// GetEvse retrieves evse by ID
	GetEvse(ctx context.Context, locationId, evseId string, withConnectors bool) (*Evse, error)
	// SearchEvses searches evses
	SearchEvses(ctx context.Context, cr *EvseSearchCriteria) (*EvseSearchResponse, error)
	// PutConnector creates or updates connector
	PutConnector(ctx context.Context, con *Connector) (*Connector, error)
	// MergeConnector merges connector
	MergeConnector(ctx context.Context, con *Connector) (*Connector, error)
	// GetConnector retrieves connector by ID
	GetConnector(ctx context.Context, locationId, evseId, conId string) (*Connector, error)
	// SearchConnectors searches connectors
	SearchConnectors(ctx context.Context, cr *ConnectorSearchCriteria) (*ConnectorSearchResponse, error)
}

type LocationStorage interface {
	// GetLocation retrieves location by id
	GetLocation(ctx context.Context, id string, withEvse bool) (*Location, error)
	// MergeLocation merges location
	MergeLocation(ctx context.Context, loc *Location) error
	// UpdateLocation updates location
	UpdateLocation(ctx context.Context, loc *Location) error
	// DeleteLocationsByExtId deletes locations (evse + connectors) by party ext id
	DeleteLocationsByExtId(ctx context.Context, extId PartyExtId) error
	// SearchLocations searches locations
	SearchLocations(ctx context.Context, cr *LocationSearchCriteria) (*LocationSearchResponse, error)
	// GetEvse retrieves evse by id
	GetEvse(ctx context.Context, locId, evseId string, withConnectors bool) (*Evse, error)
	// MergeEvse merges evse
	MergeEvse(ctx context.Context, evse *Evse) error
	// UpdateEvse update evse
	UpdateEvse(ctx context.Context, evse *Evse) error
	// SearchEvses searches evses
	SearchEvses(ctx context.Context, cr *EvseSearchCriteria) (*EvseSearchResponse, error)
	// MergeConnector merges connector
	MergeConnector(ctx context.Context, con *Connector) error
	// UpdateConnector updates connector
	UpdateConnector(ctx context.Context, con *Connector) error
	// SearchConnectors searches connectors
	SearchConnectors(ctx context.Context, cr *ConnectorSearchCriteria) (*ConnectorSearchResponse, error)
	// GetConnector retrieves connector by id
	GetConnector(ctx context.Context, locId, evseId, conId string) (*Connector, error)
}
