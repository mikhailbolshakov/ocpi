package storage

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/storages/pg"
	"github.com/mikhailbolshakov/ocpi/domain"
)

func (s *cdrStorageImpl) toCdrDto(c *domain.Cdr) *cdr {
	if c == nil {
		return nil
	}
	dto := &cdr{
		Id:          c.Id,
		SessionId:   pg.StringToNull(c.Details.SessionId),
		PartyId:     c.ExtId.PartyId,
		CountryCode: c.ExtId.CountryCode,
		PlatformId:  c.PlatformId,
		RefId:       pg.StringToNull(c.RefId),
		LastUpdated: c.LastUpdated,
		LastSent:    c.LastSent,
	}
	dto.Details, _ = pg.ToJsonb(&c.Details)
	return dto
}

func (s *cdrStorageImpl) toCdrDomain(dto *cdr) *domain.Cdr {
	if dto == nil {
		return nil
	}
	c := &domain.Cdr{
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
		Id: dto.Id,
	}
	det, _ := pg.FromJsonb[domain.CdrDetails](dto.Details)
	if det != nil {
		c.Details = *det
	}
	return c
}

func (s *cdrStorageImpl) toCdrsDomain(dtos []*cdr) []*domain.Cdr {
	return kit.Select(dtos, s.toCdrDomain)
}
