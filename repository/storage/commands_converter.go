package storage

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/storages/pg"
	"github.com/mikhailbolshakov/ocpi/domain"
)

func (s *commandStorageImpl) toCommandDto(cmd *domain.Command) *command {
	if cmd == nil {
		return nil
	}
	dto := &command{
		Id:          cmd.Id,
		Status:      cmd.Status,
		Cmd:         cmd.Cmd,
		Deadline:    cmd.Deadline,
		AuthRef:     cmd.AuthRef,
		PartyId:     cmd.ExtId.PartyId,
		CountryCode: cmd.ExtId.CountryCode,
		PlatformId:  cmd.PlatformId,
		RefId:       pg.StringToNull(cmd.RefId),
		LastUpdated: cmd.LastUpdated,
		LastSent:    cmd.LastSent,
	}
	dto.Details, _ = pg.ToJsonb(&cmd.Details)
	return dto
}

func (s *commandStorageImpl) toCommandDomain(dto *command) *domain.Command {
	if dto == nil {
		return nil
	}
	c := &domain.Command{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     dto.PartyId,
				CountryCode: dto.CountryCode,
			},
			PlatformId:  dto.PlatformId,
			RefId:       pg.NullToString(dto.RefId),
			LastUpdated: dto.LastUpdated,
			LastSent:    dto.LastSent,
		},
		Id:       dto.Id,
		Status:   dto.Status,
		Cmd:      dto.Cmd,
		Deadline: dto.Deadline,
		AuthRef:  dto.AuthRef,
	}
	det, _ := pg.FromJsonb[domain.CommandDetails](dto.Details)
	if det != nil {
		c.Details = *det
	}
	return c
}

func (s *commandStorageImpl) toCommandsDomain(dtos []*command) []*domain.Command {
	return kit.Select(dtos, s.toCommandDomain)
}
