package impl

import (
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/model"
	"github.com/mikhailbolshakov/ocpi/usecase"
)

type sessionConverter struct {
	baseConverter
}

func NewSessionConverter() usecase.SessionConverter {
	return &sessionConverter{}
}

func (t *sessionConverter) SessionDomainToModel(sess *domain.Session) *model.OcpiSession {
	if sess == nil {
		return nil
	}
	return &model.OcpiSession{
		OcpiPartyId: model.OcpiPartyId{
			PartyId:     sess.OcpiItem.ExtId.PartyId,
			CountryCode: sess.OcpiItem.ExtId.CountryCode,
		},
		Id:              sess.Id,
		StartDateTime:   sess.Details.StartDateTime,
		EndDateTime:     sess.Details.EndDateTime,
		Kwh:             sess.Details.Kwh,
		CdrToken:        t.cdrTokenDomainToModel(sess.Details.CdrToken),
		AuthMethod:      sess.Details.AuthMethod,
		AuthRef:         sess.Details.AuthRef,
		LocationId:      sess.Details.LocationId,
		EvseId:          sess.Details.EvseId,
		ConnectorId:     sess.Details.ConnectorId,
		MeterId:         sess.Details.MeterId,
		Currency:        sess.Details.Currency,
		ChargingPeriods: t.chargingPeriodDomainToModel(sess.ChargingPeriods),
		TotalCost:       t.priceDomainToModel(sess.Details.TotalCost),
		Status:          sess.Details.Status,
		LastUpdated:     sess.LastUpdated,
	}
}

func (t *sessionConverter) TokenToCdrTokenDomain(tkn *domain.Token) *domain.CdrToken {
	if tkn == nil {
		return nil
	}
	return &domain.CdrToken{
		PartyExtId: tkn.ExtId,
		Id:         tkn.Id,
		Type:       tkn.Details.Type,
		ContractId: tkn.Details.ContractId,
	}
}

func (t *sessionConverter) SessionsDomainToModel(ts []*domain.Session) []*model.OcpiSession {
	var r []*model.OcpiSession
	for _, sess := range ts {
		r = append(r, t.SessionDomainToModel(sess))
	}
	return r
}

func (t *sessionConverter) SessionModelToDomain(sess *model.OcpiSession, platformId string) *domain.Session {
	if sess == nil {
		return nil
	}
	return &domain.Session{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     sess.PartyId,
				CountryCode: sess.CountryCode,
			},
			PlatformId:  platformId,
			LastUpdated: sess.LastUpdated,
		},
		Id:              sess.Id,
		ChargingPeriods: t.chargingPeriodModelToDomain(sess.ChargingPeriods),
		Details: domain.SessionDetails{
			StartDateTime: sess.StartDateTime,
			EndDateTime:   sess.EndDateTime,
			Kwh:           sess.Kwh,
			CdrToken:      t.cdrTokenModelToDomain(sess.CdrToken),
			AuthMethod:    sess.AuthMethod,
			AuthRef:       sess.AuthRef,
			LocationId:    sess.LocationId,
			EvseId:        sess.EvseId,
			ConnectorId:   sess.ConnectorId,
			MeterId:       sess.MeterId,
			Currency:      sess.Currency,
			TotalCost:     t.priceModelToDomain(sess.TotalCost),
			Status:        sess.Status,
		},
	}
}

func (t *sessionConverter) SessionDomainToBackend(sess *domain.Session) *backend.Session {
	if sess == nil {
		return nil
	}
	return &backend.Session{
		Id:              sess.Id,
		StartDateTime:   sess.Details.StartDateTime,
		EndDateTime:     sess.Details.EndDateTime,
		Kwh:             sess.Details.Kwh,
		CdrToken:        t.cdrTokenDomainToBackend(sess.Details.CdrToken),
		AuthMethod:      sess.Details.AuthMethod,
		AuthRef:         sess.Details.AuthRef,
		LocationId:      sess.Details.LocationId,
		EvseId:          sess.Details.EvseId,
		ConnectorId:     sess.Details.ConnectorId,
		MeterId:         sess.Details.MeterId,
		Currency:        sess.Details.Currency,
		ChargingPeriods: t.chargingPeriodDomainToBackend(sess.ChargingPeriods),
		TotalCost:       t.priceDomainToBackend(sess.Details.TotalCost),
		Status:          sess.Details.Status,
		LastUpdated:     sess.LastUpdated,
		PlatformId:      sess.PlatformId,
		RefId:           sess.RefId,
		PartyId:         sess.OcpiItem.ExtId.PartyId,
		CountryCode:     sess.OcpiItem.ExtId.CountryCode,
	}
}

func (t *sessionConverter) SessionsDomainToBackend(ts []*domain.Session) []*backend.Session {
	var r []*backend.Session
	for _, sess := range ts {
		r = append(r, t.SessionDomainToBackend(sess))
	}
	return r
}

func (t *sessionConverter) SessionBackendToDomain(sess *backend.Session, platformId string) *domain.Session {
	if sess == nil {
		return nil
	}
	return &domain.Session{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     sess.PartyId,
				CountryCode: sess.CountryCode,
			},
			PlatformId:  platformId,
			RefId:       sess.RefId,
			LastUpdated: sess.LastUpdated,
		},
		Id:              sess.Id,
		ChargingPeriods: t.chargingPeriodBackendToDomain(sess.ChargingPeriods),
		Details: domain.SessionDetails{
			StartDateTime: sess.StartDateTime,
			EndDateTime:   sess.EndDateTime,
			Kwh:           sess.Kwh,
			CdrToken:      t.cdrTokenBackendToDomain(sess.CdrToken),
			AuthMethod:    sess.AuthMethod,
			AuthRef:       sess.AuthRef,
			LocationId:    sess.LocationId,
			EvseId:        sess.EvseId,
			ConnectorId:   sess.ConnectorId,
			MeterId:       sess.MeterId,
			Currency:      sess.Currency,
			TotalCost:     t.priceBackendToDomain(sess.TotalCost),
			Status:        sess.Status,
		},
	}
}
