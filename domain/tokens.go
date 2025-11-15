package domain

import "context"

const (
	TokenTypeAdHocUser = "AD_HOC_USER"
	TokenTypeAppUser   = "APP_USER"
	TokenTypeRfid      = "RFID"
	TokenTypeOther     = "OTHER"

	TokenWLTypeAlways         = "ALWAYS"
	TokenWLTypeAllowed        = "ALLOWED"
	TokenWLTypeAllowedOffline = "ALLOWED_OFFLINE"
	TokenWLTypeNever          = "NEVER"

	ProfileTypeCheap   = "CHEAP"
	ProfileTypeFast    = "FAST"
	ProfileTypeGreen   = "GREEN"
	ProfileTypeRegular = "REGULAR"
)

type EnergyContract struct {
	SupplierName string `json:"supplierName"`         // SupplierName name of the energy supplier for this token
	ContractId   string `json:"contractId,omitempty"` // ContractId at the energy supplier, that belongs to the owner of this token
}

type TokenDetails struct {
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
}

type Token struct {
	OcpiItem
	Id      string       `json:"id"`      // Id unique ID by which this Token can be identified
	Details TokenDetails `json:"details"` // Details token details
}

type TokenSearchCriteria struct {
	PageRequest
	ExtId        *PartyExtId // ExtId by party ext ID
	IncPlatforms []string    // IncPlatforms includes platform Ids
	ExcPlatforms []string    // ExcPlatforms exclude platform Ids
	Ids          []string    // Ids by list of Ids
	RefId        string      // RefId search by ref id
}

type TokenSearchResponse struct {
	PageResponse
	Items []*Token
}

type TokenService interface {
	// PutToken creates or updates Token
	PutToken(ctx context.Context, tkn *Token) (*Token, error)
	// MergeToken merges Token
	MergeToken(ctx context.Context, tkn *Token) (*Token, error)
	// GetToken retrieves Token by ID
	GetToken(ctx context.Context, tknId string) (*Token, error)
	// SearchTokens searches Tokens
	SearchTokens(ctx context.Context, cr *TokenSearchCriteria) (*TokenSearchResponse, error)
	// DeleteTokensByExtId deletes all tokens by party ext id
	DeleteTokensByExtId(ctx context.Context, extId PartyExtId) error
	// ValidateToken validate token
	ValidateToken(ctx context.Context, tkn *Token) error
}

type TokenStorage interface {
	// MergeToken creates or updates Token
	MergeToken(ctx context.Context, tkn *Token) error
	// UpdateToken updates Token
	UpdateToken(ctx context.Context, tkn *Token) error
	// GetToken retrieves Token by ID
	GetToken(ctx context.Context, tknId string) (*Token, error)
	// DeleteTokensByExtId deletes all tokens by party ext id
	DeleteTokensByExtId(ctx context.Context, extId PartyExtId) error
	// SearchTokens searches Tokens
	SearchTokens(ctx context.Context, cr *TokenSearchCriteria) (*TokenSearchResponse, error)
}
