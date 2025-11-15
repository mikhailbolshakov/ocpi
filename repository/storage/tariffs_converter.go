package storage

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/storages/pg"
	"github.com/mikhailbolshakov/ocpi/domain"
)

func (s *tariffStorageImpl) toTariffDto(trf *domain.Tariff) *tariff {
	if trf == nil {
		return nil
	}
	dto := &tariff{
		Id:          trf.Id,
		PartyId:     trf.ExtId.PartyId,
		CountryCode: trf.ExtId.CountryCode,
		PlatformId:  trf.PlatformId,
		RefId:       pg.StringToNull(trf.RefId),
		LastUpdated: trf.LastUpdated,
		LastSent:    trf.LastSent,
	}
	dto.Details, _ = pg.ToJsonb(&trf.Details)
	return dto
}

func (s *tariffStorageImpl) toTariffDomain(dto *tariff) *domain.Tariff {
	if dto == nil {
		return nil
	}
	trf := &domain.Tariff{
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
	det, _ := pg.FromJsonb[domain.TariffDetails](dto.Details)
	if det != nil {
		trf.Details = *det
	}
	return trf
}

func (s *tariffStorageImpl) toTariffsDomain(dtos []*tariff) []*domain.Tariff {
	return kit.Select(dtos, s.toTariffDomain)
}
