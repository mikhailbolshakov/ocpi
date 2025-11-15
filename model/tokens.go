package model

import "time"

type OcpiLocationRef struct {
	LocationId string   `json:"location_id"`         // LocationId unique identifier for the location
	EvseUids   []string `json:"evse_uids,omitempty"` //  EvseUids unique identifiers for EVSEs within the CPO’s platform
}

type OcpiEnergyContract struct {
	SupplierName string `json:"supplier_name"`         // SupplierName name of the energy supplier for this token
	ContractId   string `json:"contract_id,omitempty"` // ContractId at the energy supplier, that belongs to the owner of this token
}

type OcpiToken struct {
	OcpiPartyId
	Id                 string              `json:"uid"`                            // Uid unique ID by which this Token can be identified
	Type               string              `json:"type"`                           // Type of the token
	ContractId         string              `json:"contract_id"`                    // ContractId EV Driver contract token within the eMSP’s platform
	VisualNumber       string              `json:"visual_number,omitempty"`        // VisualNumber number/identification as printed on the Token (RFID card)
	Issuer             string              `json:"issuer"`                         // Issuer issuing company
	GroupId            string              `json:"group_id,omitempty"`             // GroupId to group a couple of tokens
	Valid              *bool               `json:"valid"`                          // Valid is this Token valid
	WhiteList          string              `json:"whitelist"`                      // WhiteList indicates what type of white-listing is allowed
	Lang               string              `json:"language,omitempty"`             // Lang code ISO 639-1
	DefaultProfileType string              `json:"default_profile_type,omitempty"` // DefaultProfileType default Charging Preference
	EnergyContract     *OcpiEnergyContract `json:"energy_contract,omitempty"`      // EnergyContract energy supplier/contract
	LastUpdated        time.Time           `json:"last_updated"`                   // LastUpdated timestamp when this Token was last updated
}

type OcpiTokenAuthorizationInfo struct {
	Token    OcpiToken        `json:"token"`                             // Token complete Token object for which this authorization was requested
	Allowed  string           `json:"allowed"`                           // Allowed status of the Token, and whether charging is allowed
	Location *OcpiLocationRef `json:"location,omitempty"`                // Location optional reference to the location
	AuthRef  string           `json:"authorization_reference,omitempty"` // AuthRef reference to the authorization given by the eMSP
	Info     *OcpiDisplayText `json:"info,omitempty"`                    // Info display tex
}

type OcpiTokensResponse struct {
	OcpiResponse
	Data []*OcpiToken `json:"data"`
}
