package storage

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/storages/pg"
	"github.com/mikhailbolshakov/ocpi/domain"
)

func (s *partyStorageImpl) toPartyDto(p *domain.Party) *party {
	if p == nil {
		return nil
	}
	dto := &party{
		Id:          p.Id,
		PlatformId:  p.PlatformId,
		PartyId:     p.ExtId.PartyId,
		CountryCode: p.ExtId.CountryCode,
		Status:      p.Status,
		RefId:       pg.StringToNull(p.RefId),
		Roles:       p.Roles,
		LastUpdated: p.LastUpdated,
		LastSent:    p.LastSent,
	}
	dto.Details, _ = pg.ToJsonb(&partyDetails{BusinessDetails: p.BusinessDetails})
	return dto
}

func (s *partyStorageImpl) toPartyDomain(dto *party) *domain.Party {
	if dto == nil {
		return nil
	}
	p := &domain.Party{
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
		Id:     dto.Id,
		Roles:  dto.Roles,
		Status: dto.Status,
	}
	det, _ := pg.FromJsonb[partyDetails](dto.Details)
	if det != nil {
		p.BusinessDetails = det.BusinessDetails
	}
	return p
}

func (s *partyStorageImpl) toPartiesDomain(dtos []*party) []*domain.Party {
	return kit.Select(dtos, s.toPartyDomain)
}
