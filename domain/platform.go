package domain

import (
	"context"
	"github.com/mikhailbolshakov/ocpi"
)

const (
	ModuleIdCredentials   = "credentials"
	ModuleIdCdrs          = "cdrs"
	ModuleIdCommands      = "commands"
	ModuleIdHubClientInfo = "hubclientinfo"
	ModuleIdLocations     = "locations"
	ModuleIdSessions      = "sessions"
	ModuleIdTariffs       = "tariffs"
	ModuleIdTokens        = "tokens"

	ConnectionStatusConnected = "CONNECTED"
	ConnectionStatusOffLine   = "OFFLINE"
	ConnectionStatusPlanned   = "PLANNED"
	ConnectionStatusSuspended = "SUSPENDED"

	RoleCPO   = "CPO"
	RoleEMSP  = "EMSP"
	RoleHUB   = "HUB"
	RoleNAP   = "NAP"
	RoleNSP   = "NSP"
	RoleOTHER = "OTHER"
	RoleSCSP  = "SCSP"
)

var RoleMap = map[string]bool{
	RoleHUB:   true,
	RoleCPO:   true,
	RoleEMSP:  true,
	RoleNAP:   true,
	RoleNSP:   true,
	RoleSCSP:  false,
	RoleOTHER: false,
}

// Endpoint url
type Endpoint string

// PlatformToken OCPI token
type PlatformToken string

// RoleEndpoint mapping role and endpoint
type RoleEndpoint map[string]Endpoint

// ModuleEndpoints mapping module code and endpoints per role
type ModuleEndpoints map[string]RoleEndpoint

// Versions mapping version num and version endpoint
type Versions map[string]Endpoint

// VersionInfo platform version info
type VersionInfo struct {
	Current   string   `json:"current,omitempty"`   // Current is a current platform version
	Available Versions `json:"available,omitempty"` // Available full list of the available versions
	VersionEp Endpoint `json:"versionEp,omitempty"` // VersionEp is endpoint to get list of versions
}

// PushSupport specifies if platform supports pushing for the particular module
type PushSupport struct {
	Credentials   bool `json:"credentials"`   // Credentials pushing supported
	Cdrs          bool `json:"cdrs"`          // Cdrs pushing supported
	Commands      bool `json:"commands"`      // Commands pushing supported
	HubClientInfo bool `json:"hubClientInfo"` // HubClientInfo pushing supported
	Locations     bool `json:"locations"`     // Locations pushing supported
	Sessions      bool `json:"sessions"`      // Sessions pushing supported
	Tariffs       bool `json:"tariffs"`       // Tariffs pushing supported
	Tokens        bool `json:"tokens"`        // Tokens pushing supported
}

type ProtocolDetails struct {
	PushSupport PushSupport `json:"pushSupport"` // PushSupport specifies pushing support by module
}

// Platform OCPI platform (either local or remote)
type Platform struct {
	Id          string           `json:"id"`                    // Id unique ID
	TokenA      PlatformToken    `json:"tokenA"`                // TokenA is an initial token used for handshake
	TokenB      PlatformToken    `json:"tokenB,omitempty"`      // TokenB token received from sender
	TokenC      PlatformToken    `json:"tokenC,omitempty"`      // TokenC token received from receiver
	TokenBase64 *bool            `json:"tokenBase64,omitempty"` // TokenBase64 if true, token is base64 encoded
	Name        string           `json:"name,omitempty"`        // Name platform name
	Role        string           `json:"role,omitempty"`        // Role platform role
	VersionInfo VersionInfo      `json:"versionInfo"`           // VersionInfo platform version info
	Endpoints   ModuleEndpoints  `json:"endpoints,omitempty"`   // Endpoints list of available endpoints
	Status      string           `json:"status,omitempty"`      // Status platform status
	Remote      bool             `json:"remote"`                // Remote is true, if it's a remote platform (not local)
	Protocol    *ProtocolDetails `json:"protocol,omitempty"`    // Protocol specifies protocol details
}

type PlatformSearchCriteria struct {
	Roles    []string // Roles roles
	ExcRoles []string // ExcRoles exclude roles
	Statuses []string // Statuses search by statuses
	IncIds   []string // IncIds include IDs
	ExcIds   []string // ExcIds exclude IDs
	Remote   *bool    // Remote search by remote flag
}

// TokenGenerator responsible for token generation
type TokenGenerator interface {
	// Generate generates a new token
	Generate(ctx context.Context) (PlatformToken, error)
	// TryBase64Decode tries to decode token from base64
	TryBase64Decode(token PlatformToken) (PlatformToken, bool)
	// Base64Encode encodes token to base64
	Base64Encode(tkn PlatformToken) PlatformToken
}

type PlatformService interface {
	// Init initializes service
	Init(ctx context.Context, cfg *ocpi.CfgOcpiConfig) error
	// Merge creates or updates platform
	Merge(ctx context.Context, platform *Platform) (*Platform, error)
	// SetStatus sets status
	SetStatus(ctx context.Context, platformId string, status string) (*Platform, error)
	// Get retrieves platform by ID
	Get(ctx context.Context, platformId string) (*Platform, error)
	// GetByTokenA retrieves platform by token A
	GetByTokenA(ctx context.Context, token PlatformToken) (*Platform, error)
	// GetByTokenB retrieves platform by token B
	GetByTokenB(ctx context.Context, token PlatformToken) (*Platform, error)
	// GetByTokenC retrieves platform by token C
	GetByTokenC(ctx context.Context, token PlatformToken) (*Platform, error)
	// Search searches platforms based on criteria
	Search(ctx context.Context, cr *PlatformSearchCriteria) ([]*Platform, error)
	// RoleEndpoint returns endpoint of the requested module and role. If roles isn't supported, empty string is returned
	RoleEndpoint(ctx context.Context, platform *Platform, module, role string) Endpoint
}

type LocalPlatformService interface {
	// Init initializes service
	Init(ctx context.Context, cfg *ocpi.CfgOcpiConfig) error
	// Get retrieves local platform
	Get(ctx context.Context) (*Platform, error)
	// InitializePlatform initializes local platform
	InitializePlatform(ctx context.Context) error
	// GetPlatformId retrieves local platform ID
	GetPlatformId(ctx context.Context) string
	// GetEndpoints retrieves all endpoints supported by requested version
	GetEndpoints(ctx context.Context, version string) ModuleEndpoints
}

type PlatformStorage interface {
	// CreatePlatform creates platform
	CreatePlatform(ctx context.Context, p *Platform) error
	// UpdatePlatform updates platform
	UpdatePlatform(ctx context.Context, p *Platform) error
	// DeletePlatform deletes platform
	DeletePlatform(ctx context.Context, p *Platform) error
	// GetPlatform retrieves platform
	GetPlatform(ctx context.Context, id string) (*Platform, error)
	// GetPlatformByTokenA retrieves platform by token A
	GetPlatformByTokenA(ctx context.Context, tokenA PlatformToken) (*Platform, error)
	// GetPlatformByTokenB retrieves platform by token B
	GetPlatformByTokenB(ctx context.Context, token PlatformToken) (*Platform, error)
	// GetPlatformByTokenC retrieves platform by token C
	GetPlatformByTokenC(ctx context.Context, token PlatformToken) (*Platform, error)
	// SearchPlatforms searches platforms by criteria
	SearchPlatforms(ctx context.Context, cr *PlatformSearchCriteria) ([]*Platform, error)
}
