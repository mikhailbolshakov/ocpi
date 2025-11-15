package storage

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/storages/pg"
	"github.com/mikhailbolshakov/ocpi/domain"
)

func (s *platformStorageImpl) toPlatformDto(p *domain.Platform) *platform {
	if p == nil {
		return nil
	}
	dto := &platform{
		Id:     p.Id,
		TokenA: string(p.TokenA),
		TokenB: pg.StringToNull(string(p.TokenB)),
		TokenC: pg.StringToNull(string(p.TokenC)),
		Name:   p.Name,
		Role:   p.Role,
		Status: p.Status,
		Remote: p.Remote,
	}
	dto.Details, _ = pg.ToJsonb(&platformDetails{
		VersionInfo: p.VersionInfo,
		Endpoints:   p.Endpoints,
		Protocol:    p.Protocol,
		TokenBase64: p.TokenBase64,
	})
	return dto
}

func (s *platformStorageImpl) toPlatformDomain(dto *platform) *domain.Platform {
	if dto == nil {
		return nil
	}
	p := &domain.Platform{
		Id:     dto.Id,
		TokenA: domain.PlatformToken(dto.TokenA),
		TokenB: domain.PlatformToken(pg.NullToString(dto.TokenB)),
		TokenC: domain.PlatformToken(pg.NullToString(dto.TokenC)),
		Name:   dto.Name,
		Role:   dto.Role,
		Status: dto.Status,
		Remote: dto.Remote,
	}
	det, _ := pg.FromJsonb[platformDetails](dto.Details)
	if det != nil {
		p.Endpoints = det.Endpoints
		p.VersionInfo = det.VersionInfo
		p.Protocol = det.Protocol
		p.TokenBase64 = det.TokenBase64
	}
	return p
}

func (s *platformStorageImpl) toPlatformsDomain(dtos []*platform) []*domain.Platform {
	return kit.Select(dtos, s.toPlatformDomain)
}
