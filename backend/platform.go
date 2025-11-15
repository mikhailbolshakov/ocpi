package backend

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

type PlatformRequest struct {
	Id           string           `json:"id,omitempty"`           // Id platform id
	TokenA       string           `json:"tokenA,omitempty"`       // TokenA platform token A
	Name         string           `json:"name,omitempty"`         // Name platform name
	Role         string           `json:"role,omitempty"`         // Role platform role
	GetVersionEp string           `json:"getVersionEp,omitempty"` // GetVersionEp versions endpoint
	TokenBase64  *bool            `json:"tokenBase64,omitempty"`  // TokenBase64 if true, token is base64 encoded
	Protocol     *ProtocolDetails `json:"protocol,omitempty"`     // Protocol details
}

type VersionInfo struct {
	Current      string            `json:"current,omitempty"`      // Current is a current platform version
	Available    map[string]string `json:"available,omitempty"`    // Available full list of the available versions
	GetVersionEp string            `json:"getVersionEp,omitempty"` // GetVersionEp is endpoint to get list of versions
}

type RoleEndpoint struct {
	Val map[string]string `json:"val,omitempty"`
}

type Platform struct {
	Id          string                   `json:"id,omitempty"`          // Id of platform
	TokenA      string                   `json:"tokenA,omitempty"`      // TokenA platform token A
	TokenB      string                   `json:"tokenB,omitempty"`      // TokenB platform token B
	TokenC      string                   `json:"tokenC,omitempty"`      // TokenC platform token C
	Name        string                   `json:"name,omitempty"`        // Name platform name
	Role        string                   `json:"role,omitempty"`        // Role platform role
	VersionInfo *VersionInfo             `json:"versionInfo,omitempty"` // VersionInfo platform versions
	Endpoints   map[string]*RoleEndpoint `json:"endpoints,omitempty"`   // Endpoints platform endpoints
	Status      string                   `json:"status,omitempty"`      // Status platform status
	Remote      bool                     `json:"remote,omitempty"`      // Remote if platform remote
	TokenBase64 *bool                    `json:"tokenBase64,omitempty"` // TokenBase64 if true, token is base64 encoded
	Protocol    *ProtocolDetails         `json:"protocol,omitempty"`    // Protocol details
}
