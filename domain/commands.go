package domain

import (
	"context"
	"time"
)

const (
	CmdResponseTypeNotSupported   = "NOT_SUPPORTED"
	CmdResponseTypeRejected       = "REJECTED"
	CmdResponseTypeAccepted       = "ACCEPTED"
	CmdResponseTypeUnknownSession = "UNKNOWN_SESSION"

	CmdResultTypeAccepted             = "ACCEPTED"
	CmdResultTypeCancelledReservation = "CANCELED_RESERVATION"
	CmdResultTypeEvseOccupied         = "EVSE_OCCUPIED"
	CmdResultTypeEvseInoperative      = "EVSE_INOPERATIVE"
	CmdResultTypeFailed               = "FAILED"
	CmdResultTypeNotSupported         = "NOT_SUPPORTED"
	CmdResultTypeRejected             = "REJECTED"
	CmdResultTypeTimeout              = "TIMEOUT"
	CmdResultTypeUnknownReservation   = "UNKNOWN_RESERVATION"

	CmdStatusRequestAccepted        = "accepted-request"
	CmdStatusRequestRejected        = "rejected-request"
	CmdStatusRequestProcessedOk     = "ok"
	CmdStatusRequestProcessedFailed = "failed"
	CmdStatusRequestExpired         = "expired"

	CmdStartSession      = "START_SESSION"
	CmdStopSession       = "STOP_SESSION"
	CmdReserve           = "RESERVE_NOW"
	CmdCancelReservation = "CANCEL_RESERVATION"
	CmdUnlockConnector   = "UNLOCK_CONNECTOR"
)

type ReserveNow struct {
	Token         *Token    `json:"token"`                 // Token object the Charge Point has to use to start a new session
	ExpireDate    time.Time `json:"expireDate"`            // ExpireDate when this reservation ends, in UTC
	ReservationId string    `json:"reservationId"`         // ReservationId unique for this reservation
	LocationId    string    `json:"locationId"`            // LocationId on which a session is to be started
	EvseId        string    `json:"evseUid,omitempty"`     // EvseId of the EVSE of this Location on which a session is to be started
	ConnectorId   string    `json:"connectorId,omitempty"` // ConnectorId  of the Connector of the EVSE on which a session 	is to be started
}

type CancelReservation struct {
	ReservationId string `json:"reservationId"` // ReservationId unique for this reservation
}

type StartSession struct {
	Token       *Token   `json:"token"`                 // Token object the Charge Point has to use to start a new session
	LocationId  string   `json:"locationId"`            // LocationId on which a session is to be started
	EvseId      string   `json:"evseId,omitempty"`      // EvseId of the EVSE of this Location on which a session is to be started
	ConnectorId string   `json:"connectorId,omitempty"` // ConnectorId  of the Connector of the EVSE on which a session 	is to be started
	KwhLimit    *float64 `json:"kwhLimit,omitempty"`    // KwhLimit allows setting limit on maximum Kwh (Yandex extension, isn't a part of OCPI protocol)
}

type StopSession struct {
	SessionId string `json:"sessionId"` // SessionId of the Session that is requested to be stopped
}

type UnlockConnector struct {
	LocationId  string `json:"locationId"`            // LocationId on which a session is to be started
	EvseId      string `json:"evseUid,omitempty"`     // EvseId of the EVSE of this Location on which a session is to be started
	ConnectorId string `json:"connectorId,omitempty"` // ConnectorId  of the Connector of the EVSE on which a session 	is to be started
}

type CommandDetails struct {
	ResponseUrl       Endpoint           `json:"responseUrl"`                 // ResponseUrl URL that the CommandResult POST should be sent to
	Processing        Processing         `json:"proc,omitempty"`              // Processing command processing details
	Reserve           *ReserveNow        `json:"reserve,omitempty"`           // Reserve reservation details
	CancelReservation *CancelReservation `json:"cancelReservation,omitempty"` // CancelReservation cancel reservation details
	StartSession      *StartSession      `json:"startSession,omitempty"`      // StartSession start session details
	StopSession       *StopSession       `json:"stopSession,omitempty"`       // StopSession stop session
	UnlockConnector   *UnlockConnector   `json:"unlockConnector,omitempty"`   // UnlockConnector unlock connector details
}

type Processing struct {
	Status string `json:"status,omitempty"`
	ErrMsg string `json:"errMsg,omitempty"`
}

type Command struct {
	OcpiItem
	Id       string         `json:"id"`       // Id unique identifier
	Status   string         `json:"status"`   // Status command status
	Cmd      string         `json:"cmd"`      // Cmd command type
	Deadline time.Time      `json:"deadline"` // Deadline timestamp after that command request is no longer active
	Details  CommandDetails `json:"details"`  // Details command details
	AuthRef  string         `json:"authRef"`  // AuthRef identified command
}

type CommandSearchCriteria struct {
	PageRequest
	ExtId         *PartyExtId // ExtId by party ext ID
	IncPlatforms  []string    // IncPlatforms includes platform Ids
	ExcPlatforms  []string    // ExcPlatforms exclude platform Ids
	Ids           []string    // Ids by list Ids
	RefId         string      // RefId by ref id
	AuthRef       string      // AuthRef by auth ref
	Cmd           string      // Cmd command type
	Statuses      []string    // Statuses by command statuses
	DeadlineLE    *time.Time  // DeadlineLE retrieves items where deadline less or equal the given value
	RetrieveAll   bool        // RetrieveAll if true, skip paging (only for backend usage)
	ReservationId string      // ReservationId search by reservation id
}

type CommandSearchResponse struct {
	PageResponse
	Items []*Command
}

type CommandService interface {
	// Create creates a new command
	Create(ctx context.Context, cmd *Command) (*Command, error)
	// Update updates existent command
	Update(ctx context.Context, cmd *Command) (*Command, error)
	// Get retrieves command by id
	Get(ctx context.Context, id string) (*Command, error)
	// DeleteCommandsByExt deletes all commands by ext party id
	DeleteCommandsByExt(ctx context.Context, extId PartyExtId) error
	// GetByAuthRef retrieves command by auth_ref
	GetByAuthRef(ctx context.Context, authRef string) (*Command, error)
	// SearchCommands searches commands
	SearchCommands(ctx context.Context, cr *CommandSearchCriteria) (*CommandSearchResponse, error)
}

type CommandStorage interface {
	// CreateCommand creates a command
	CreateCommand(ctx context.Context, cmd *Command) error
	// UpdateCommand updates a command
	UpdateCommand(ctx context.Context, cmd *Command) error
	// GetCommand retrieves a command
	GetCommand(ctx context.Context, cmdId string) (*Command, error)
	// GetCommandByAuthRef retrieves command by auth_ref
	GetCommandByAuthRef(ctx context.Context, authRef string) (*Command, error)
	// DeleteCommand deletes a command
	DeleteCommand(ctx context.Context, cmdId string) error
	// DeleteCommandsByExt deletes all commands by ext party id
	DeleteCommandsByExt(ctx context.Context, extId PartyExtId) error
	// SearchCommands searches commands
	SearchCommands(ctx context.Context, cr *CommandSearchCriteria) (*CommandSearchResponse, error)
}
