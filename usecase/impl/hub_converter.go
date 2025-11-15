package impl

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/model"
)

func (h *hubUc) partyDomainToModel(party *domain.Party) []*model.OcpiClientInfo {
	if party == nil {
		return nil
	}
	var r []*model.OcpiClientInfo
	for _, role := range party.Roles {
		r = append(r, &model.OcpiClientInfo{
			OcpiPartyId: model.OcpiPartyId{
				PartyId:     party.ExtId.PartyId,
				CountryCode: party.ExtId.CountryCode,
			},
			Role:        role,
			Status:      party.Status,
			LastUpdated: party.LastUpdated,
		})
	}
	return r
}

func (h *hubUc) partyModelToDomain(platformId string, ci *model.OcpiClientInfo, existent *domain.Party) *domain.Party {
	if ci == nil {
		return nil
	}

	// build list of roles
	roles := []string{ci.Role}
	if existent != nil {
		roles = append(roles, existent.Roles...)
	}

	return &domain.Party{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     ci.PartyId,
				CountryCode: ci.CountryCode,
			},
			PlatformId:  platformId,
			LastUpdated: ci.LastUpdated,
		},
		Roles:  kit.Strings(roles).Distinct(),
		Status: ci.Status,
	}
}

func (h *hubUc) partyDomainToBackend(party *domain.Party) *backend.Party {
	if party == nil {
		return nil
	}
	return &backend.Party{
		PartyId:     party.ExtId.PartyId,
		CountryCode: party.ExtId.CountryCode,
		Roles:       party.Roles,
		Status:      party.Status,
		Id:          party.Id,
		RefId:       party.RefId,
		LastUpdated: party.LastUpdated,
	}
}
