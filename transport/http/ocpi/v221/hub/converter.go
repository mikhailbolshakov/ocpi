package hub

import (
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/model"
)

func (c *ctrlImpl) toClientInfoApi(p *domain.Party) *model.OcpiClientInfo {
	if p == nil {
		return nil
	}
	return &model.OcpiClientInfo{
		OcpiPartyId: model.OcpiPartyId{
			PartyId:     p.ExtId.PartyId,
			CountryCode: p.ExtId.CountryCode,
		},
		Role:        p.Roles[0],
		Status:      p.Status,
		LastUpdated: p.LastUpdated,
	}
}

func (c *ctrlImpl) toClientInfosApi(rs *domain.PartySearchResponse) []*model.OcpiClientInfo {
	var r []*model.OcpiClientInfo
	for _, i := range rs.Items {
		for _, role := range i.Roles {
			r = append(r, &model.OcpiClientInfo{
				OcpiPartyId: model.OcpiPartyId{
					PartyId:     i.ExtId.PartyId,
					CountryCode: i.ExtId.CountryCode,
				},
				Role:        role,
				Status:      i.Status,
				LastUpdated: i.LastUpdated,
			})
		}
	}
	return r
}
