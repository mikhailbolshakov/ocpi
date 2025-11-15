package domain

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	"time"
)

const (
	LogEventGetVersions         = "versions.get"
	LogEventGetVersionDetails   = "version-details.get"
	LogEventPostCredentials     = "credentials.post"
	LogEventPutCredentials      = "credentials.put"
	LogEventGetCredentials      = "credentials.get"
	LogEventDelCredentials      = "credentials.del"
	LogEventGetHubClients       = "hubclientinfo.get"
	LogEventPutClientInfo       = "client-info.put"
	LogEventPutLocation         = "locations.put"
	LogEventPatchLocation       = "locations.patch"
	LogEventGetLocations        = "locations.get"
	LogEventGetLocation         = "location.get"
	LogEventPutEvse             = "evse.put"
	LogEventPatchEvse           = "evse.patch"
	LogEventGetEvse             = "evse.get"
	LogEventPutCon              = "con.put"
	LogEventPatchCon            = "con.patch"
	LogEventGetCon              = "con.get"
	LogEventPutTariff           = "tariff.put"
	LogEventPatchTariff         = "tariff.patch"
	LogEventGetTariffs          = "tariffs.get"
	LogEventGetTariff           = "tariff.get"
	LogEventPutToken            = "token.put"
	LogEventPatchToken          = "token.patch"
	LogEventGetTokens           = "tokens.get"
	LogEventGetToken            = "token.get"
	LogEventPutSession          = "session.put"
	LogEventPatchSession        = "session.patch"
	LogEventGetSessions         = "sessions.get"
	LogEventGetSession          = "session.get"
	LogEventPostCommand         = "command.post"
	LogEventPostCommandResponse = "command-rs.post"
	LogEventPostCdr             = "cdr.post"
	LogEventGetCdrs             = "cdrs.get"
	LogEventGetCdr              = "cdr.get"
)

type LogMessage struct {
	Event          string // Event log event
	Url            string // Url
	Token          string // Token used for auth
	RequestId      string // RequestId request ID
	CorrelationId  string // CorrelationId correlation ID
	FromPlatform   string // FromPlatform platform sender
	ToPlatform     string // FromPlatform platform receiver
	RequestBody    any    // RequestBody request body
	ResponseBody   any    // ResponseBody response body
	Headers        any    // Headers
	ResponseStatus int    // ResponseStatus http status of the response
	OcpiStatus     int    // OcpiStatus OCPI status of the response
	Err            error  // Err error
	In             bool   // In if true, incoming call, otherwise outgoing
	DurationMs     int64  // DurationMs request duration ms
}

type SearchLogCriteria struct {
	kit.PagingRequest
	Events       []string   // Events filter by list of events
	RequestId    string     // RequestId filter by request id
	FromPlatform string     // FromPlatform filter by source platform
	ToPlatform   string     // ToPlatform filter by target platform
	OcpiStatus   *int       // OcpiStatus filter by ocpi status
	HttpStatus   *int       // HttpStatus filter by HTTP status
	Incoming     *bool      // Incoming if true, only incoming requests are retrieved
	Error        *bool      // Error if true, only requests with errors are retrieved
	DateFrom     *time.Time // DateFrom filter by date period
	DateTo       *time.Time // DateTo filter by date period
}

// OcpiLogService logs OCPI events
type OcpiLogService interface {
	// Init inits log with severity
	Init(severity string)
	// Log logs incoming requests
	Log(ctx context.Context, msg *LogMessage)
	// Search retrieves log entries by criteria
	Search(ctx context.Context, criteria *SearchLogCriteria) ([]*LogMessage, error)
}

type OcpiLogStorage interface {
	// Save saves log messages async
	Save(ctx context.Context, msg *LogMessage)
	// SearchLog retrieves log entries by criteria
	SearchLog(ctx context.Context, criteria *SearchLogCriteria) ([]*LogMessage, error)
}
