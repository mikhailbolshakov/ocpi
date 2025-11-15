package backend

import "time"

const (
	TokenTypeAdHocUser = "AD_HOC_USER"
	TokenTypeAppUser   = "APP_USER"
	TokenTypeRfid      = "RFID"
	TokenTypeOther     = "OTHER"

	TokenWLTypeAlways         = "ALWAYS"
	TokenWLTypeAllowed        = "ALLOWED"
	TokenWLTypeAllowedOffline = "ALLOWED_OFFLINE"
	TokenWLTypeNever          = "NEVER"
)

type EnergyContract struct {
	SupplierName string `json:"supplierName"`         // SupplierName name of the energy supplier for this token
	ContractId   string `json:"contractId,omitempty"` // ContractId at the energy supplier, that belongs to the owner of this token
}

type Token struct {
	Id                 string          `json:"id"`                           // Id unique ID by which this Token can be identified
	Type               string          `json:"type"`                         // Type of the token
	ContractId         string          `json:"contractId"`                   // ContractId EV Driver contract token within the eMSP’s platform
	VisualNumber       string          `json:"visualNumber,omitempty"`       // VisualNumber number/identification as printed on the Token (RFID card)
	Issuer             string          `json:"issuer"`                       // Issuer issuing company
	GroupId            string          `json:"groupId,omitempty"`            // GroupId to group a couple of tokens
	Valid              *bool           `json:"valid"`                        // Valid is this Token valid
	WhiteList          string          `json:"whitelist"`                    // WhiteList indicates what type of white-listing is allowed
	Lang               string          `json:"language,omitempty"`           // Lang code ISO 639-1
	DefaultProfileType string          `json:"defaultProfileType,omitempty"` // DefaultProfileType default Charging Preference
	EnergyContract     *EnergyContract `json:"energyContract,omitempty"`     // EnergyContract energy supplier/contract
	LastUpdated        time.Time       `json:"lastUpdated"`                  // LastUpdated when this Tariff was last updated
	PlatformId         string          `json:"platformId"`                   // PlatformId rel to platform
	RefId              string          `json:"refId"`                        // RefId any external relation
	PartyId            string          `json:"partyId,omitempty"`            // PartyId should be unique within country
	CountryCode        string          `json:"countryCode,omitempty"`        // CountryCode alfa-2 code
}

type TokenSearchResponse struct {
	PageInfo *PageResponse `json:"pageInfo,omitempty"`
	Items    []*Token      `json:"items,omitempty"`
}
