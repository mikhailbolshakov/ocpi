package storage

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/storages/pg"
	"github.com/mikhailbolshakov/ocpi/domain"
)

func (s *locationStorageImpl) toLocationDto(loc *domain.Location) *location {
	if loc == nil {
		return nil
	}
	dto := &location{
		Id:          loc.Id,
		PartyId:     loc.ExtId.PartyId,
		CountryCode: loc.ExtId.CountryCode,
		PlatformId:  loc.PlatformId,
		RefId:       pg.StringToNull(loc.RefId),
		LastUpdated: loc.LastUpdated,
		LastSent:    loc.LastSent,
	}
	dto.Details, _ = pg.ToJsonb(&loc.Details)
	return dto
}

func (s *locationStorageImpl) toEvseDto(e *domain.Evse) *evse {
	if e == nil {
		return nil
	}
	dto := &evse{
		Id:          e.Id,
		LocationId:  e.LocationId,
		Status:      e.Status,
		PartyId:     e.ExtId.PartyId,
		CountryCode: e.ExtId.CountryCode,
		PlatformId:  e.PlatformId,
		RefId:       pg.StringToNull(e.RefId),
		LastUpdated: e.LastUpdated,
		LastSent:    e.LastSent,
	}
	dto.Details, _ = pg.ToJsonb(&e.Details)
	return dto
}

func (s *locationStorageImpl) toConnectorDto(c *domain.Connector) *connector {
	if c == nil {
		return nil
	}
	dto := &connector{
		Id:          c.Id,
		LocationId:  c.LocationId,
		EvseId:      c.EvseId,
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

func (s *locationStorageImpl) toConnectorDomain(dto *connector) *domain.Connector {
	if dto == nil {
		return nil
	}
	c := &domain.Connector{
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
		Id:         dto.Id,
		LocationId: dto.LocationId,
		EvseId:     dto.EvseId,
	}
	det, _ := pg.FromJsonb[domain.ConnectorDetails](dto.Details)
	if det != nil {
		c.Details = *det
	}
	return c
}

func (s *locationStorageImpl) toConnectorsDomain(conDtos []*connector) []*domain.Connector {
	return kit.Select(conDtos, s.toConnectorDomain)
}

func (s *locationStorageImpl) toEvseDomain(dto *evse, conDtos []*connector) *domain.Evse {
	if dto == nil {
		return nil
	}
	e := &domain.Evse{
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
		Id:         dto.Id,
		LocationId: dto.LocationId,
		Status:     dto.Status,
		Connectors: s.toConnectorsDomain(conDtos),
	}
	det, _ := pg.FromJsonb[domain.EvseDetails](dto.Details)
	if det != nil {
		e.Details = *det
	}
	return e
}

func (s *locationStorageImpl) toEvsesDomain(evseDtos []*evse, conDtos map[string][]*connector) []*domain.Evse {
	var r []*domain.Evse
	for _, e := range evseDtos {
		r = append(r, s.toEvseDomain(e, conDtos[e.Id]))
	}
	return r
}

func (s *locationStorageImpl) toLocationDomain(dto *location, evseDtos []*evse, conDtos []*connector) *domain.Location {
	if dto == nil {
		return nil
	}

	conMap := make(map[string][]*connector)
	for _, con := range conDtos {
		conMap[con.EvseId] = append(conMap[con.EvseId], con)
	}

	loc := &domain.Location{
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
		Id:    dto.Id,
		Evses: s.toEvsesDomain(evseDtos, conMap),
	}
	det, _ := pg.FromJsonb[domain.LocationDetails](dto.Details)
	if det != nil {
		loc.Details = *det
	}
	return loc
}

func (s *locationStorageImpl) toLocationsDomain(dtos []*location, evseDtos []*evse, conDtos []*connector) []*domain.Location {
	evseMap := make(map[string][]*evse)
	for _, evse := range evseDtos {
		evseMap[evse.LocationId] = append(evseMap[evse.LocationId], evse)
	}
	var r []*domain.Location
	for _, loc := range dtos {
		r = append(r, s.toLocationDomain(loc, evseMap[loc.Id], conDtos))
	}
	return r
}
