package usecase

import (
	"context"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/model"
)

type CommandConverter interface {
	// CommandDomainToBackend converts a command from domain to backend
	CommandDomainToBackend(cmd *domain.Command) *backend.Command
	CommandsDomainToBackend(commands []*domain.Command) []*backend.Command

	// StartSessionCommandDomainToModel converts start session command from domain to ocpi
	StartSessionCommandDomainToModel(cmd *domain.Command) *model.OcpiStartSession
	// StartSessionCommandModelToDomain converts start session command from ocpi to domain
	StartSessionCommandModelToDomain(cmd *model.OcpiStartSession, platformId string) *domain.Command
	// StartSessionCommandBackendToDomain converts start session command from backend to domain
	StartSessionCommandBackendToDomain(cmd *backend.StartSessionRequest, platformId string) *domain.Command

	// StopSessionCommandDomainToModel converts stop session command from domain to ocpi
	StopSessionCommandDomainToModel(cmd *domain.Command) *model.OcpiStopSession
	// StopSessionCommandModelToDomain converts stop session command from ocpi to domain
	StopSessionCommandModelToDomain(cmd *model.OcpiStopSession, platformId string) *domain.Command
	// StopSessionCommandBackendToDomain converts stop session command from backend to domain
	StopSessionCommandBackendToDomain(cmd *backend.StopSessionRequest, platformId string) *domain.Command

	// ReserveNowCommandDomainToModel converts reserve now command from domain to ocpi
	ReserveNowCommandDomainToModel(cmd *domain.Command) *model.OcpiReserveNow
	// ReserveNowCommandModelToDomain converts reserve now command from ocpi to domain
	ReserveNowCommandModelToDomain(cmd *model.OcpiReserveNow, platformId string) *domain.Command
	// ReserveNowCommandBackendToDomain converts reserve now command from backend to domain
	ReserveNowCommandBackendToDomain(cmd *backend.ReserveNowRequest, platformId string) *domain.Command

	// CancelReservationCommandDomainToModel converts cancel reservation command from domain to ocpi
	CancelReservationCommandDomainToModel(cmd *domain.Command) *model.OcpiCancelReservation
	// CancelReservationCommandModelToDomain converts cancel reservation command from ocpi to domain
	CancelReservationCommandModelToDomain(cmd *model.OcpiCancelReservation, platformId string) *domain.Command
	// CancelReservationCommandBackendToDomain converts cancel reservation command from backend to domain
	CancelReservationCommandBackendToDomain(cmd *backend.CancelReservationRequest, platformId string) *domain.Command
}

type CommandUc interface {
	// OnRemoteStartSession fires when a remote platform requests starting session
	OnRemoteStartSession(ctx context.Context, platformId string, rq *model.OcpiStartSession) (*model.OcpiCommandResponse, error)
	// OnRemoteStopSession fires when a remote platform requests stopping session
	OnRemoteStopSession(ctx context.Context, platformId string, rq *model.OcpiStopSession) (*model.OcpiCommandResponse, error)
	// OnRemoteReserve fires when a remote platform requests reservation
	OnRemoteReserve(ctx context.Context, platformId string, rq *model.OcpiReserveNow) (*model.OcpiCommandResponse, error)
	// OnRemoteCancelReservation fires when a remote platform requests reservation cancellation
	OnRemoteCancelReservation(ctx context.Context, platformId string, rq *model.OcpiCancelReservation) (*model.OcpiCommandResponse, error)
	// OnRemoteUnlockConnector fires when a remote platform requests unlock connector
	OnRemoteUnlockConnector(ctx context.Context, platformId string, rq *model.OcpiUnlockConnector) (*model.OcpiCommandResponse, error)
	// OnRemoteSetResponse fires when a remote platform sets response
	OnRemoteSetResponse(ctx context.Context, platformId, uid string, rq *model.OcpiCommandResult) error
	// RemoteCommandsDeadlineCronHandler handles all hanging local commands
	RemoteCommandsDeadlineCronHandler(ctx context.Context)

	// OnLocalStartSession fires when a local platform requests starting session
	OnLocalStartSession(ctx context.Context, rq *domain.Command) error
	// OnLocalStopSession fires when a local platform requests stopping session
	OnLocalStopSession(ctx context.Context, rq *domain.Command) error
	// OnLocalReserve fires when a local platform requests reservation
	OnLocalReserve(ctx context.Context, rq *domain.Command) error
	// OnLocalCancelReservation fires when a local platform requests reservation cancellation
	OnLocalCancelReservation(ctx context.Context, rq *domain.Command) error
	// OnLocalUnlockConnector fires when a local platform requests unlock connector
	OnLocalUnlockConnector(ctx context.Context, rq *domain.Command) error
	// OnLocalCommandSetResponse fires when a local platform sets command response
	OnLocalCommandSetResponse(ctx context.Context, uid, status, errMsg string) error
	// LocalCommandsDeadlineCronHandler handles all hanging remote commands
	LocalCommandsDeadlineCronHandler(ctx context.Context)
}

type RemoteCommandRepository interface {
	// PostCommandAsync sends a command to the remote platform
	PostCommandAsync(ctx context.Context, rq *OcpiRepositoryErrHandlerRequest, cmdType string, cmd any)
	// PostCommandResponseAsync sends response to the command request
	PostCommandResponseAsync(ctx context.Context, rq *OcpiRepositoryErrHandlerRequestG[*model.OcpiCommandResult])
}
