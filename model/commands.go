package model

import "time"

const (
	CommandStartSession      = "START_SESSION"
	CommandStopSession       = "STOP_SESSION"
	CommandReserve           = "RESERVE_NOW"
	CommandCancelReservation = "CANCEL_RESERVATION"
	CommandUnlockConnector   = "UNLOCK_CONNECTOR"
)

type OcpiReserveNow struct {
	ResponseUrl            string     `json:"response_url"`                      // ResponseUrl URL that the CommandResult POST should be sent to
	Token                  *OcpiToken `json:"token"`                             // Token object the Charge Point has to use to start a new session
	ExpireDate             time.Time  `json:"expire_date"`                       // ExpireDate when this reservation ends, in UTC
	ReservationId          string     `json:"reservation_id"`                    // ReservationId unique for this reservation
	LocationId             string     `json:"location_id"`                       // LocationId on which a session is to be started
	EvseId                 string     `json:"evse_uid,omitempty"`                // EvseId of the EVSE of this Location on which a session is to be started
	ConnectorId            string     `json:"connector_id,omitempty"`            // ConnectorId  of the Connector of the EVSE on which a session 	is to be started
	AuthorizationReference string     `json:"authorization_reference,omitempty"` // AuthorizationReference reference to the authorization given by the eMSP
}

type OcpiCancelReservation struct {
	ResponseUrl   string `json:"response_url"`   // ResponseUrl URL that the CommandResult POST should be sent to
	ReservationId string `json:"reservation_id"` // ReservationId unique for this reservation
}

type OcpiStartSession struct {
	ResponseUrl            string     `json:"response_url"`                      // ResponseUrl URL that the CommandResult POST should be sent to
	Token                  *OcpiToken `json:"token"`                             // Token object the Charge Point has to use to start a new session
	LocationId             string     `json:"location_id"`                       // LocationId on which a session is to be started
	EvseId                 string     `json:"evse_uid,omitempty"`                // EvseId of the EVSE of this Location on which a session is to be started
	ConnectorId            string     `json:"connector_id,omitempty"`            // ConnectorId  of the Connector of the EVSE on which a session 	is to be started
	AuthorizationReference string     `json:"authorization_reference,omitempty"` // AuthorizationReference reference to the authorization given by the eMSP
	KwhLimit               *float64   `json:"kwh_limit,omitempty"`               // KwhLimit allows setting limit on maximum Kwh (Yandex extension, isn't a part of OCPI protocol)
}

type OcpiStopSession struct {
	ResponseUrl string `json:"response_url"` // ResponseUrl URL that the CommandResult POST should be sent to
	SessionId   string `json:"session_id"`   // SessionId of the Session that is requested to be stopped
}

type OcpiUnlockConnector struct {
	ResponseUrl string `json:"response_url"`           // ResponseUrl URL that the CommandResult POST should be sent to
	LocationId  string `json:"location_id"`            // LocationId on which a session is to be started
	EvseId      string `json:"evse_uid,omitempty"`     // EvseId of the EVSE of this Location on which a session is to be started
	ConnectorId string `json:"connector_id,omitempty"` // ConnectorId  of the Connector of the EVSE on which a session 	is to be started
}

type OcpiCommandResponse struct {
	Result  string             `json:"result"`            // Result from the CPO on the command request
	Timeout int                `json:"timeout"`           // Timeout for this command in seconds
	Message []*OcpiDisplayText `json:"message,omitempty"` // Message human-readable description of the result
}

type OcpiCommandResult struct {
	Result  string             `json:"result"`            // Result of the command request as sent by the Charge Point to the CPO
	Message []*OcpiDisplayText `json:"message,omitempty"` // Message human-readable description of the result
}
