package model

const (
	ModuleIdCredentials   = "credentials"
	ModuleIdCdrs          = "cdrs"
	ModuleIdCommands      = "commands"
	ModuleIdHubClientInfo = "hubclientinfo"
	ModuleIdLocations     = "locations"
	ModuleIdSessions      = "sessions"
	ModuleIdTariffs       = "tariffs"
	ModuleIdTokens        = "tokens"

	OcpiStatusCodeOk                   = 1000
	OcpiStatusGenClientError           = 2000
	OcpiStatusInvalidParamError        = 2001
	OcpiStatusNotEnoughInfoError       = 2002
	OcpiStatusUnknownLocationError     = 2003
	OcpiStatusUnknownTokenError        = 2004
	OcpiStatusGenServerError           = 3000
	OcpiStatusUnableUseApiError        = 3001
	OcpiStatusUnsupportedVersionError  = 3002
	OcpiStatusNoMatchingEndpointsError = 3003
	OcpiStatusHubGenErrorError         = 4000
	OcpiStatusUnknownReceiverError     = 4001
	OcpiStatusTimeoutError             = 4002
	OcpiStatusConnectionError          = 4003

	OcpiSender   = "SENDER"
	OcpiReceiver = "RECEIVER"

	OcpiStatusMessageSuccess = "Success"

	OcpiHeaderRequestId       = "X-Request-Id"
	OcpiHeaderCorrelationId   = "X-correlation-id"
	OcpiHeaderFromCountryCode = "OCPI-from-country-code"
	OcpiHeaderFromPartyId     = "OCPI-from-party-id"
	OcpiHeaderToCountryCode   = "OCPI-to-country-code"
	OcpiHeaderToPartyId       = "OCPI-to-party-id"
	OcpiHeaderAuth            = "authorization"
	OcpiHeaderTotalCount      = "X-Total-Count"
	OcpiHeaderLimit           = "X-Limit"
	OcpiHeaderLink            = "Link"

	OcpiCtxCorrelationId   = "ocpi-corr-Id"
	OcpiCtxCountryCode     = "ocpi-from-cc"
	OcpiCtxFromCountryCode = "ocpi-from-cc"
	OcpiCtxFromParty       = "ocpi-from-party"
	OcpiCtxToCountryCode   = "ocpi-to-cc"
	OcpiCtxToParty         = "ocpi-to-party"
	OcpiCtxPlatform        = "ocpi-platform"

	OcpiQueryParamDateFrom = "date_from"
	OcpiQueryParamDateTo   = "date_to"
	OcpiQueryParamOffset   = "offset"
	OcpiQueryParamLimit    = "limit"

	OcpiQueryParamPartyId     = "party_id"
	OcpiQueryParamCountryCode = "country_code"

	OcpiQueryParamCommand = "command"
	OcpiQueryParamUid     = "uid"

	OcpiStatusField = "ocpi-status"

	OcpiRoleCPO   = "CPO"
	OcpiRoleEMSP  = "EMSP"
	OcpiRoleHUB   = "HUB"
	OcpiRoleNAP   = "NAP"
	OcpiRoleNSP   = "NSP"
	OcpiRoleOTHER = "OTHER"
	OcpiRoleSCSP  = "SCSP"

	OcpiStatusConnected = "CONNECTED"
	OcpiStatusOffLine   = "OFFLINE"
	OcpiStatusPlanned   = "PLANNED"
	OcpiStatusSuspended = "SUSPENDED"

	OcpiEvseStatusAvailable   = "AVAILABLE"
	OcpiEvseStatusBlocked     = "BLOCKED"
	OcpiEvseStatusCharging    = "CHARGING"
	OcpiEvseStatusInOperative = "INOPERATIVE"
	OcpiEvseStatusOutOfOrder  = "OUTOFORDER"
	OcpiEvseStatusPlanned     = "PLANNED"
	OcpiEvseStatusRemoved     = "REMOVED"
	OcpiEvseStatusReserved    = "RESERVED"
	OcpiEvseStatusUnknown     = "UNKNOWN"
)

var (
	OcpiModules = map[string]struct {
		SenderOnly bool
	}{
		ModuleIdCredentials:   {SenderOnly: true},
		ModuleIdCdrs:          {},
		ModuleIdCommands:      {},
		ModuleIdHubClientInfo: {SenderOnly: true},
		ModuleIdLocations:     {},
		ModuleIdSessions:      {},
		ModuleIdTariffs:       {},
		ModuleIdTokens:        {},
	}
)
