package backend

import "time"

const (
	CmdResultTypeAccepted             = "ACCEPTED"
	CmdResultTypeCancelledReservation = "CANCELED_RESERVATION"
	CmdResultTypeEvseOccupied         = "EVSE_OCCUPIED"
	CmdResultTypeEvseInoperative      = "EVSE_INOPERATIVE"
	CmdResultTypeFailed               = "FAILED"
	CmdResultTypeNotSupported         = "NOT_SUPPORTED"
	CmdResultTypeRejected             = "REJECTED"
	CmdResultTypeTimeout              = "TIMEOUT"
	CmdResultTypeUnknownReservation   = "UNKNOWN_RESERVATION"

	CmdStartSession      = "START_SESSION"
	CmdStopSession       = "STOP_SESSION"
	CmdReserve           = "RESERVE_NOW"
	CmdCancelReservation = "CANCEL_RESERVATION"
	CmdUnlockConnector   = "UNLOCK_CONNECTOR"

	CmdStatusRequestAccepted        = "accepted-request"
	CmdStatusRequestRejected        = "rejected-request"
	CmdStatusRequestProcessedOk     = "ok"
	CmdStatusRequestProcessedFailed = "failed"
	CmdStatusRequestExpired         = "expired"
)

type ReserveNow struct {
	Token                  *Token    `json:"token"`                 // Token object the Charge Point has to use to start a new session
	ExpireDate             time.Time `json:"expireDate"`            // ExpireDate when this reservation ends, in UTC
	ReservationId          string    `json:"reservationId"`         // ReservationId unique for this reservation
	LocationId             string    `json:"locationId"`            // LocationId on which a session is to be started
	EvseId                 string    `json:"evseId,omitempty"`      // EvseId of the EVSE of this Location on which a session is to be started
	ConnectorId            string    `json:"connectorId,omitempty"` // ConnectorId  of the Connector of the EVSE on which a session 	is to be started
	AuthorizationReference string    `json:"authRef,omitempty"`     // AuthorizationReference reference to the authorization given by the eMSP
}

type ReserveNowRequest struct {
	Id                     string    `json:"id"`                    // Id unique identifier
	LocationId             string    `json:"locationId"`            // LocationId on which a session is to be started
	EvseId                 string    `json:"evseId"`                // EvseId of the EVSE of this Location on which a session is to be started
	ConnectorId            string    `json:"connectorId"`           // ConnectorId  of the Connector of the EVSE on which a session 	is to be started
	AuthorizationReference string    `json:"authRef,omitempty"`     // AuthorizationReference reference to the authorization given by the eMSP
	Token                  Token     `json:"token"`                 // Token token info
	PartyId                string    `json:"partyId,omitempty"`     // PartyId should be unique within country
	CountryCode            string    `json:"countryCode,omitempty"` // CountryCode alfa-2 code
	RefId                  string    `json:"refId,omitempty"`       // RefId any external relation
	ReservationId          string    `json:"reservationId"`         // ReservationId unique for this reservation
	ExpireDate             time.Time `json:"expireDate"`            // ExpireDate when this reservation ends, in UTC
}

type CancelReservation struct {
	ReservationId string `json:"reservationId"` // ReservationId unique for this reservation
}

type CancelReservationRequest struct {
	Id            string `json:"id"`                    // Id unique identifier
	PartyId       string `json:"partyId,omitempty"`     // PartyId should be unique within country
	CountryCode   string `json:"countryCode,omitempty"` // CountryCode alfa-2 code
	RefId         string `json:"refId,omitempty"`       // RefId any external relation
	ReservationId string `json:"reservationId"`         // ReservationId unique for this reservation
}

type StartSession struct {
	Token                  *Token   `json:"token"`                 // Token object the Charge Point has to use to start a new session
	LocationId             string   `json:"locationId"`            // LocationId on which a session is to be started
	EvseId                 string   `json:"evseId,omitempty"`      // EvseId of the EVSE of this Location on which a session is to be started
	ConnectorId            string   `json:"connectorId,omitempty"` // ConnectorId  of the Connector of the EVSE on which a session 	is to be started
	AuthorizationReference string   `json:"authRef,omitempty"`     // AuthorizationReference reference to the authorization given by the eMSP
	KwhLimit               *float64 `json:"kwhLimit,omitempty"`    // KwhLimit allows setting limit on maximum Kwh (Yandex extension, isn't a part of OCPI protocol)
}

type StopSession struct {
	SessionId string `json:"sessionId"` // SessionId of the Session that is requested to be stopped
}

type UnlockConnector struct {
	LocationId  string `json:"locationId"`            // LocationId on which a session is to be started
	EvseId      string `json:"evseId,omitempty"`      // EvseId of the EVSE of this Location on which a session is to be started
	ConnectorId string `json:"connectorId,omitempty"` // ConnectorId  of the Connector of the EVSE on which a session is to be started
}

type CommandDetails struct {
	Reserve           *ReserveNow        `json:"reserve,omitempty"`           // Reserve reservation details
	CancelReservation *CancelReservation `json:"cancelReservation,omitempty"` // CancelReservation cancel reservation details
	StartSession      *StartSession      `json:"startSession,omitempty"`      // StartSession start session details
	StopSession       *StopSession       `json:"stopSession,omitempty"`       // StopSession stop session
	UnlockConnector   *UnlockConnector   `json:"unlockConnector,omitempty"`   // UnlockConnector unlock connector details
}

type Command struct {
	Id          string         `json:"id"`                    // Id unique identifier
	Status      string         `json:"status"`                // Status command status
	Cmd         string         `json:"cmd"`                   // Cmd command type
	Deadline    time.Time      `json:"deadline"`              // Deadline timestamp after that command request is no longer active
	Details     CommandDetails `json:"details"`               // Details command details
	AuthRef     string         `json:"authRef"`               // AuthRef identified command
	PartyId     string         `json:"partyId,omitempty"`     // PartyId should be unique within country
	CountryCode string         `json:"countryCode,omitempty"` // CountryCode alfa-2 code
	RefId       string         `json:"refId,omitempty"`       // RefId any external relation
}

type CommandResponse struct {
	Status string `json:"status"` // Status command status
	ErrMsg string `json:"errMsg"` // ErrMsg error message
}

type CmdTokenRq struct {
	Id          string `json:"id"`          // Id token ID. If not exists, it'll be created with default params
	PartyId     string `json:"partyId"`     // PartyId owner of the token
	CountryCode string `json:"countryCode"` // CountryCode owner of the token
	Num         string `json:"num"`         // Num token number
	Type        string `json:"type"`        // Type token type
}

type StartSessionRequest struct {
	Id                     string   `json:"id"`                    // Id unique identifier
	LocationId             string   `json:"locationId"`            // LocationId on which a session is to be started
	EvseId                 string   `json:"evseUid"`               // EvseId of the EVSE of this Location on which a session is to be started
	ConnectorId            string   `json:"connectorId"`           // ConnectorId  of the Connector of the EVSE on which a session 	is to be started
	AuthorizationReference string   `json:"authRef,omitempty"`     // AuthorizationReference reference to the authorization given by the eMSP
	Token                  Token    `json:"token"`                 // Token token info
	PartyId                string   `json:"partyId,omitempty"`     // PartyId should be unique within country
	CountryCode            string   `json:"countryCode,omitempty"` // CountryCode alfa-2 code
	RefId                  string   `json:"refId,omitempty"`       // RefId any external relation
	KwhLimit               *float64 `json:"kwhLimit,omitempty"`    // KwhLimit allows setting limit on maximum Kwh (Yandex extension, isn't a part of OCPI protocol)
}

type StopSessionRequest struct {
	Id          string `json:"id"`                    // Id unique identifier
	SessionId   string `json:"sessionId"`             // SessionId of the Session that is requested to be stopped
	PartyId     string `json:"partyId,omitempty"`     // PartyId should be unique within country
	CountryCode string `json:"countryCode,omitempty"` // CountryCode alfa-2 code
	RefId       string `json:"refId,omitempty"`       // RefId any external relation
}

type CommandSearchResponse struct {
	PageInfo *PageResponse `json:"pageInfo,omitempty"`
	Items    []*Command    `json:"items,omitempty"`
}
