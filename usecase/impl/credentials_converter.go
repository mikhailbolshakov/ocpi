package impl

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/model"
	"github.com/mikhailbolshakov/ocpi/usecase"
)

type credentialsConverter struct {
	baseConverter
}

func NewCredentialsConverter() usecase.CredentialsConverter {
	return &credentialsConverter{}
}

func (c *credentialsConverter) PartiesModelToDomain(platformId, status string, roles ...*model.OcpiCredentialRole) []*domain.Party {
	var partiesMap = map[domain.PartyExtId]*domain.Party{}

	for _, role := range roles {
		extId := domain.PartyExtId{
			PartyId:     role.PartyId,
			CountryCode: role.CountryCode,
		}
		if party, ok := partiesMap[extId]; ok {
			party.Roles = append(party.Roles, role.Role)
		} else {
			partiesMap[extId] = &domain.Party{
				OcpiItem: domain.OcpiItem{
					ExtId:       extId,
					PlatformId:  platformId,
					LastUpdated: kit.Now(),
				},
				Roles:           []string{role.Role},
				BusinessDetails: c.businessDetailsModelToDomain(role.BusinessDetails),
				Status:          status,
			}
		}
	}

	var res []*domain.Party
	for _, v := range partiesMap {
		res = append(res, v)
	}

	return res
}

func (c *credentialsConverter) PartyBackendToDomain(platformId string, p *backend.Party) *domain.Party {
	if p == nil {
		return nil
	}
	return &domain.Party{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     p.PartyId,
				CountryCode: p.CountryCode,
			},
			PlatformId:  platformId,
			RefId:       p.RefId,
			LastUpdated: p.LastUpdated,
		},
		Id:              p.Id,
		Roles:           p.Roles,
		BusinessDetails: c.businessDetailsBackendToDomain(p.BusinessDetails),
		Status:          p.Status,
	}
}

func (c *credentialsConverter) PartiesDomainToModel(parties ...*domain.Party) []*model.OcpiCredentialRole {
	var res []*model.OcpiCredentialRole
	for _, p := range parties {
		for _, r := range p.Roles {
			res = append(res, &model.OcpiCredentialRole{
				OcpiPartyId: model.OcpiPartyId{
					PartyId:     p.ExtId.PartyId,
					CountryCode: p.ExtId.CountryCode,
				},
				Role:            r,
				BusinessDetails: c.businessDetailsDomainToModel(p.BusinessDetails),
			})
		}
	}
	return res
}

func (c *credentialsConverter) toBusinessDetailsBackend(bd *domain.BusinessDetails) *backend.BusinessDetails {
	if bd == nil {
		return nil
	}
	r := &backend.BusinessDetails{
		Name:    bd.Name,
		Website: bd.Website,
	}
	if bd.Logo != nil {
		r.Logo = &backend.Image{
			Url:       bd.Logo.Url,
			Thumbnail: bd.Logo.Thumbnail,
			Category:  bd.Logo.Category,
			Type:      bd.Logo.Type,
			Width:     bd.Logo.Width,
			Height:    bd.Logo.Height,
		}
	}
	return r
}

func (c *credentialsConverter) PartyDomainToBackend(party *domain.Party) *backend.Party {
	if party == nil {
		return nil
	}
	return &backend.Party{
		Id:              party.Id,
		Roles:           party.Roles,
		PartyId:         party.ExtId.PartyId,
		CountryCode:     party.ExtId.CountryCode,
		BusinessDetails: c.toBusinessDetailsBackend(party.BusinessDetails),
		RefId:           party.RefId,
		Status:          party.Status,
		LastUpdated:     party.LastUpdated,
	}
}

func (c *credentialsConverter) PartiesDomainToBackend(parties []*domain.Party) []*backend.Party {
	var res []*backend.Party
	for _, p := range parties {
		res = append(res, c.PartyDomainToBackend(p))
	}
	return res
}

func (c *credentialsConverter) PlatformBackendToDomain(rq *backend.PlatformRequest) *domain.Platform {
	if rq == nil {
		return nil
	}
	r := &domain.Platform{
		Id:          rq.Id,
		TokenA:      domain.PlatformToken(rq.TokenA),
		TokenBase64: rq.TokenBase64,
		Name:        rq.Name,
		Role:        rq.Role,
		VersionInfo: domain.VersionInfo{
			VersionEp: domain.Endpoint(rq.GetVersionEp),
		},
		Remote: true,
	}
	if rq.Protocol != nil {
		r.Protocol = &domain.ProtocolDetails{
			PushSupport: domain.PushSupport{
				Credentials:   rq.Protocol.PushSupport.Credentials,
				Cdrs:          rq.Protocol.PushSupport.Cdrs,
				Commands:      rq.Protocol.PushSupport.Commands,
				HubClientInfo: rq.Protocol.PushSupport.HubClientInfo,
				Locations:     rq.Protocol.PushSupport.Locations,
				Sessions:      rq.Protocol.PushSupport.Sessions,
				Tariffs:       rq.Protocol.PushSupport.Tariffs,
				Tokens:        rq.Protocol.PushSupport.Tokens,
			},
		}
	}
	return r
}

func (c *credentialsConverter) endpointsDomainToBackend(eps domain.ModuleEndpoints) map[string]*backend.RoleEndpoint {
	r := make(map[string]*backend.RoleEndpoint)
	for k, roles := range eps {
		v := &backend.RoleEndpoint{Val: map[string]string{}}
		for r, ep := range roles {
			v.Val[r] = string(ep)
		}
		r[k] = v
	}
	return r
}

func (c *credentialsConverter) versionsDomainToBackend(eps domain.Versions) map[string]string {
	r := make(map[string]string)
	for k, v := range eps {
		r[k] = string(v)
	}
	return r
}

func (c *credentialsConverter) PlatformDomainToBackend(p *domain.Platform) *backend.Platform {
	if p == nil {
		return nil
	}
	r := &backend.Platform{
		Id:          p.Id,
		TokenA:      string(p.TokenA),
		TokenB:      string(p.TokenB),
		TokenC:      string(p.TokenC),
		TokenBase64: p.TokenBase64,
		Name:        p.Name,
		Role:        p.Role,
		VersionInfo: &backend.VersionInfo{
			Current:      p.VersionInfo.Current,
			Available:    c.versionsDomainToBackend(p.VersionInfo.Available),
			GetVersionEp: string(p.VersionInfo.VersionEp),
		},
		Endpoints: c.endpointsDomainToBackend(p.Endpoints),
		Status:    p.Status,
		Remote:    p.Remote,
	}
	if p.Protocol != nil {
		r.Protocol = &backend.ProtocolDetails{
			PushSupport: backend.PushSupport{
				Credentials:   p.Protocol.PushSupport.Credentials,
				Cdrs:          p.Protocol.PushSupport.Cdrs,
				Commands:      p.Protocol.PushSupport.Commands,
				HubClientInfo: p.Protocol.PushSupport.HubClientInfo,
				Locations:     p.Protocol.PushSupport.Locations,
				Sessions:      p.Protocol.PushSupport.Sessions,
				Tariffs:       p.Protocol.PushSupport.Tariffs,
				Tokens:        p.Protocol.PushSupport.Tokens,
			},
		}
	}
	return r
}
