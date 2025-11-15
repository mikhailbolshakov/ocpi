package storage

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/storages/pg"
	"github.com/mikhailbolshakov/ocpi/domain"
)

func (s *tokenStorageImpl) toTokenDto(tkn *domain.Token) *token {
	if tkn == nil {
		return nil
	}
	dto := &token{
		Id:          tkn.Id,
		PartyId:     tkn.ExtId.PartyId,
		CountryCode: tkn.ExtId.CountryCode,
		PlatformId:  tkn.PlatformId,
		RefId:       pg.StringToNull(tkn.RefId),
		LastUpdated: tkn.LastUpdated,
		LastSent:    tkn.LastSent,
	}
	dto.Details, _ = pg.ToJsonb(&tkn.Details)
	return dto
}

func (s *tokenStorageImpl) toTokenDomain(dto *token) *domain.Token {
	if dto == nil {
		return nil
	}
	tkn := &domain.Token{
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
	det, _ := pg.FromJsonb[domain.TokenDetails](dto.Details)
	if det != nil {
		tkn.Details = *det
	}
	return tkn
}

func (s *tokenStorageImpl) toTokensDomain(dtos []*token) []*domain.Token {
	return kit.Select(dtos, s.toTokenDomain)
}
