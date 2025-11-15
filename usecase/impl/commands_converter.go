package impl

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/model"
	"github.com/mikhailbolshakov/ocpi/usecase"
)

type commandConverter struct {
	tokenConverter usecase.TokenConverter
}

func NewCommandConverter(tokenConverter usecase.TokenConverter) usecase.CommandConverter {
	return &commandConverter{
		tokenConverter: tokenConverter,
	}
}

func (c *commandConverter) StartSessionCommandDomainToModel(cmd *domain.Command) *model.OcpiStartSession {
	if cmd == nil {
		return nil
	}
	return &model.OcpiStartSession{
		ResponseUrl:            string(cmd.Details.ResponseUrl),
		Token:                  c.tokenConverter.TokenDomainToModel(cmd.Details.StartSession.Token),
		LocationId:             cmd.Details.StartSession.LocationId,
		EvseId:                 cmd.Details.StartSession.EvseId,
		ConnectorId:            cmd.Details.StartSession.ConnectorId,
		AuthorizationReference: cmd.AuthRef,
		KwhLimit:               cmd.Details.StartSession.KwhLimit,
	}
}

func (c *commandConverter) StartSessionCommandModelToDomain(cmd *model.OcpiStartSession, platformId string) *domain.Command {
	if cmd == nil {
		return nil
	}
	return &domain.Command{
		OcpiItem: domain.OcpiItem{
			PlatformId:  platformId,
			LastUpdated: kit.Now(),
		},
		Cmd: domain.CmdStartSession,
		Details: domain.CommandDetails{
			ResponseUrl: domain.Endpoint(cmd.ResponseUrl),
			StartSession: &domain.StartSession{
				Token:       c.tokenConverter.TokenModelToDomain(cmd.Token, platformId),
				LocationId:  cmd.LocationId,
				EvseId:      cmd.EvseId,
				ConnectorId: cmd.ConnectorId,
				KwhLimit:    cmd.KwhLimit,
			},
		},
		AuthRef: cmd.AuthorizationReference,
	}
}

func (c *commandConverter) StartSessionCommandBackendToDomain(cmd *backend.StartSessionRequest, platformId string) *domain.Command {
	if cmd == nil {
		return nil
	}
	return &domain.Command{
		Id: cmd.Id,
		OcpiItem: domain.OcpiItem{
			LastUpdated: kit.Now(),
			PlatformId:  platformId,
			ExtId: domain.PartyExtId{
				PartyId:     cmd.PartyId,
				CountryCode: cmd.CountryCode,
			},
			RefId: cmd.RefId,
		},
		Cmd: domain.CmdStartSession,
		Details: domain.CommandDetails{
			StartSession: &domain.StartSession{
				Token:       c.tokenConverter.TokenBackendToDomain(&cmd.Token, platformId),
				LocationId:  cmd.LocationId,
				EvseId:      cmd.EvseId,
				ConnectorId: cmd.ConnectorId,
				KwhLimit:    cmd.KwhLimit,
			},
		},
		AuthRef: cmd.AuthorizationReference,
	}
}

func (c *commandConverter) StopSessionCommandDomainToModel(cmd *domain.Command) *model.OcpiStopSession {
	if cmd == nil {
		return nil
	}
	return &model.OcpiStopSession{
		ResponseUrl: string(cmd.Details.ResponseUrl),
		SessionId:   cmd.Details.StopSession.SessionId,
	}
}

func (c *commandConverter) StopSessionCommandModelToDomain(cmd *model.OcpiStopSession, platformId string) *domain.Command {
	if cmd == nil {
		return nil
	}
	return &domain.Command{
		OcpiItem: domain.OcpiItem{
			PlatformId:  platformId,
			LastUpdated: kit.Now(),
		},
		Cmd: domain.CmdStopSession,
		Details: domain.CommandDetails{
			ResponseUrl: domain.Endpoint(cmd.ResponseUrl),
			StopSession: &domain.StopSession{
				SessionId: cmd.SessionId,
			},
		},
	}
}

func (c *commandConverter) StopSessionCommandBackendToDomain(cmd *backend.StopSessionRequest, platformId string) *domain.Command {
	if cmd == nil {
		return nil
	}
	return &domain.Command{
		Id: cmd.Id,
		OcpiItem: domain.OcpiItem{
			LastUpdated: kit.Now(),
			PlatformId:  platformId,
			ExtId: domain.PartyExtId{
				PartyId:     cmd.PartyId,
				CountryCode: cmd.CountryCode,
			},
			RefId: cmd.RefId,
		},
		Cmd: domain.CmdStopSession,
		Details: domain.CommandDetails{
			StopSession: &domain.StopSession{
				SessionId: cmd.SessionId,
			},
		},
	}
}

func (c *commandConverter) CommandDomainToBackend(cmd *domain.Command) *backend.Command {
	if cmd == nil {
		return nil
	}
	return &backend.Command{
		Id:          cmd.Id,
		Status:      cmd.Status,
		Cmd:         cmd.Cmd,
		Deadline:    cmd.Deadline,
		Details:     *c.detailsDomainToBackend(cmd, &cmd.Details),
		AuthRef:     cmd.AuthRef,
		PartyId:     cmd.ExtId.PartyId,
		CountryCode: cmd.ExtId.CountryCode,
		RefId:       cmd.RefId,
	}
}

func (c *commandConverter) CommandsDomainToBackend(commands []*domain.Command) []*backend.Command {
	return kit.Select(commands, c.CommandDomainToBackend)
}

func (c *commandConverter) detailsDomainToBackend(cmd *domain.Command, det *domain.CommandDetails) *backend.CommandDetails {
	if det == nil {
		return nil
	}
	r := &backend.CommandDetails{}
	if det.Reserve != nil {
		r.Reserve = &backend.ReserveNow{
			Token:                  c.tokenConverter.TokenDomainToBackend(det.Reserve.Token),
			ExpireDate:             det.Reserve.ExpireDate,
			ReservationId:          det.Reserve.ReservationId,
			LocationId:             det.Reserve.LocationId,
			EvseId:                 det.Reserve.EvseId,
			ConnectorId:            det.Reserve.ConnectorId,
			AuthorizationReference: cmd.AuthRef,
		}
	}
	if det.StartSession != nil {
		r.StartSession = &backend.StartSession{
			Token:                  c.tokenConverter.TokenDomainToBackend(det.StartSession.Token),
			LocationId:             det.StartSession.LocationId,
			EvseId:                 det.StartSession.EvseId,
			ConnectorId:            det.StartSession.ConnectorId,
			AuthorizationReference: cmd.AuthRef,
			KwhLimit:               det.StartSession.KwhLimit,
		}
	}
	if det.StopSession != nil {
		r.StopSession = &backend.StopSession{
			SessionId: det.StopSession.SessionId,
		}
	}
	if det.CancelReservation != nil {
		r.CancelReservation = &backend.CancelReservation{ReservationId: det.CancelReservation.ReservationId}
	}
	if det.UnlockConnector != nil {
		r.UnlockConnector = &backend.UnlockConnector{
			LocationId:  det.UnlockConnector.LocationId,
			EvseId:      det.UnlockConnector.EvseId,
			ConnectorId: det.UnlockConnector.ConnectorId,
		}
	}
	return r
}

func (c *commandConverter) ReserveNowCommandDomainToModel(cmd *domain.Command) *model.OcpiReserveNow {
	if cmd == nil {
		return nil
	}
	return &model.OcpiReserveNow{
		ResponseUrl:            string(cmd.Details.ResponseUrl),
		Token:                  c.tokenConverter.TokenDomainToModel(cmd.Details.Reserve.Token),
		ExpireDate:             cmd.Details.Reserve.ExpireDate,
		ReservationId:          cmd.Details.Reserve.ReservationId,
		LocationId:             cmd.Details.Reserve.LocationId,
		EvseId:                 cmd.Details.Reserve.EvseId,
		ConnectorId:            cmd.Details.Reserve.ConnectorId,
		AuthorizationReference: cmd.AuthRef,
	}
}

func (c *commandConverter) ReserveNowCommandModelToDomain(cmd *model.OcpiReserveNow, platformId string) *domain.Command {
	if cmd == nil {
		return nil
	}
	return &domain.Command{
		OcpiItem: domain.OcpiItem{
			PlatformId:  platformId,
			LastUpdated: kit.Now(),
		},
		Cmd: domain.CmdReserve,
		Details: domain.CommandDetails{
			ResponseUrl: domain.Endpoint(cmd.ResponseUrl),
			Reserve: &domain.ReserveNow{
				Token:         c.tokenConverter.TokenModelToDomain(cmd.Token, platformId),
				ExpireDate:    cmd.ExpireDate,
				ReservationId: cmd.ReservationId,
				LocationId:    cmd.LocationId,
				EvseId:        cmd.EvseId,
				ConnectorId:   cmd.ConnectorId,
			},
		},
		AuthRef: cmd.AuthorizationReference,
	}
}

func (c *commandConverter) ReserveNowCommandBackendToDomain(cmd *backend.ReserveNowRequest, platformId string) *domain.Command {
	if cmd == nil {
		return nil
	}
	return &domain.Command{
		Id: cmd.Id,
		OcpiItem: domain.OcpiItem{
			LastUpdated: kit.Now(),
			PlatformId:  platformId,
			ExtId: domain.PartyExtId{
				PartyId:     cmd.PartyId,
				CountryCode: cmd.CountryCode,
			},
			RefId: cmd.RefId,
		},
		Cmd: domain.CmdReserve,
		Details: domain.CommandDetails{
			Reserve: &domain.ReserveNow{
				Token:         c.tokenConverter.TokenBackendToDomain(&cmd.Token, platformId),
				ExpireDate:    cmd.ExpireDate,
				ReservationId: cmd.ReservationId,
				LocationId:    cmd.LocationId,
				EvseId:        cmd.EvseId,
				ConnectorId:   cmd.ConnectorId,
			},
		},
		AuthRef: cmd.AuthorizationReference,
	}
}

func (c *commandConverter) CancelReservationCommandDomainToModel(cmd *domain.Command) *model.OcpiCancelReservation {
	if cmd == nil {
		return nil
	}
	return &model.OcpiCancelReservation{
		ResponseUrl:   string(cmd.Details.ResponseUrl),
		ReservationId: cmd.Details.CancelReservation.ReservationId,
	}
}

func (c *commandConverter) CancelReservationCommandModelToDomain(cmd *model.OcpiCancelReservation, platformId string) *domain.Command {
	if cmd == nil {
		return nil
	}
	return &domain.Command{
		OcpiItem: domain.OcpiItem{
			PlatformId:  platformId,
			LastUpdated: kit.Now(),
		},
		Cmd: domain.CmdCancelReservation,
		Details: domain.CommandDetails{
			ResponseUrl: domain.Endpoint(cmd.ResponseUrl),
			CancelReservation: &domain.CancelReservation{
				ReservationId: cmd.ReservationId,
			},
		},
	}
}

func (c *commandConverter) CancelReservationCommandBackendToDomain(cmd *backend.CancelReservationRequest, platformId string) *domain.Command {
	if cmd == nil {
		return nil
	}
	return &domain.Command{
		Id: cmd.Id,
		OcpiItem: domain.OcpiItem{
			LastUpdated: kit.Now(),
			PlatformId:  platformId,
			ExtId: domain.PartyExtId{
				PartyId:     cmd.PartyId,
				CountryCode: cmd.CountryCode,
			},
			RefId: cmd.RefId,
		},
		Cmd: domain.CmdCancelReservation,
		Details: domain.CommandDetails{
			CancelReservation: &domain.CancelReservation{
				ReservationId: cmd.ReservationId,
			},
		},
	}
}
