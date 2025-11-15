package storage

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/storages/pg"
	"github.com/mikhailbolshakov/ocpi/domain"
	"time"
)

func (s *sessionStorageImpl) toSessionDto(sess *domain.Session) *session {
	if sess == nil {
		return nil
	}
	dto := &session{
		Id:          sess.Id,
		PartyId:     sess.ExtId.PartyId,
		CountryCode: sess.ExtId.CountryCode,
		PlatformId:  sess.PlatformId,
		RefId:       pg.StringToNull(sess.RefId),
		LastUpdated: sess.LastUpdated,
		LastSent:    sess.LastSent,
	}
	dto.Details, _ = pg.ToJsonb(&sess.Details)
	return dto
}

func (s *sessionStorageImpl) toSessionDomain(dto *session) *domain.Session {
	if dto == nil {
		return nil
	}
	sess := &domain.Session{
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
	det, _ := pg.FromJsonb[domain.SessionDetails](dto.Details)
	if det != nil {
		sess.Details = *det
	}
	return sess
}

func (s *sessionStorageImpl) toSessionsDomain(dtos []*session) []*domain.Session {
	return kit.Select(dtos, s.toSessionDomain)
}

func (s *sessionStorageImpl) toSessionChargingPeriodsDto(sess *domain.Session, periods []*domain.ChargingPeriod, lastUpdated time.Time) []*sessionChargingPeriod {
	return kit.Select(periods, func(p *domain.ChargingPeriod) *sessionChargingPeriod {
		dto := &sessionChargingPeriod{
			SessionId:   sess.Id,
			LastUpdated: lastUpdated,
		}
		dto.Details, _ = pg.ToJsonb(p)
		return dto
	})
}

func (s *sessionStorageImpl) toSessionChargingPeriodDomain(dto *sessionChargingPeriod) *domain.ChargingPeriod {
	if dto == nil {
		return nil
	}
	det, _ := pg.FromJsonb[domain.ChargingPeriod](dto.Details)
	if det != nil {
		return det
	}
	return nil
}

func (s *sessionStorageImpl) toSessionChargingPeriodsDomain(dtos []*sessionChargingPeriod) map[string][]*domain.ChargingPeriod {
	res := make(map[string][]*domain.ChargingPeriod)
	for _, dto := range dtos {
		p, _ := pg.FromJsonb[domain.ChargingPeriod](dto.Details)
		if p != nil {
			res[dto.SessionId] = append(res[dto.SessionId], p)
		}
	}
	return res
}
