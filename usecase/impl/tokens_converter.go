package impl

import (
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/model"
	"github.com/mikhailbolshakov/ocpi/usecase"
)

type tokenConverter struct {
}

func NewTokenConverter() usecase.TokenConverter {
	return &tokenConverter{}
}

func (t *tokenConverter) TokenDomainToModel(tkn *domain.Token) *model.OcpiToken {
	if tkn == nil {
		return nil
	}
	r := &model.OcpiToken{
		OcpiPartyId: model.OcpiPartyId{
			PartyId:     tkn.ExtId.PartyId,
			CountryCode: tkn.ExtId.CountryCode,
		},
		Id:                 tkn.Id,
		Type:               tkn.Details.Type,
		ContractId:         tkn.Details.ContractId,
		VisualNumber:       tkn.Details.VisualNumber,
		Issuer:             tkn.Details.Issuer,
		GroupId:            tkn.Details.GroupId,
		Valid:              tkn.Details.Valid,
		WhiteList:          tkn.Details.WhiteList,
		Lang:               tkn.Details.Lang,
		DefaultProfileType: tkn.Details.DefaultProfileType,
		LastUpdated:        tkn.LastUpdated,
	}
	if tkn.Details.EnergyContract != nil {
		r.EnergyContract = &model.OcpiEnergyContract{
			SupplierName: tkn.Details.EnergyContract.SupplierName,
			ContractId:   tkn.Details.EnergyContract.ContractId,
		}
	}
	return r
}

func (t *tokenConverter) TokensDomainToModel(ts []*domain.Token) []*model.OcpiToken {
	var r []*model.OcpiToken
	for _, tkn := range ts {
		r = append(r, t.TokenDomainToModel(tkn))
	}
	return r
}

func (t *tokenConverter) TokenModelToDomain(tkn *model.OcpiToken, platformId string) *domain.Token {
	if tkn == nil {
		return nil
	}
	r := &domain.Token{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     tkn.PartyId,
				CountryCode: tkn.CountryCode,
			},
			PlatformId:  platformId,
			LastUpdated: tkn.LastUpdated,
		},
		Id: tkn.Id,
		Details: domain.TokenDetails{
			Type:               tkn.Type,
			ContractId:         tkn.ContractId,
			VisualNumber:       tkn.VisualNumber,
			Issuer:             tkn.Issuer,
			GroupId:            tkn.GroupId,
			Valid:              tkn.Valid,
			WhiteList:          tkn.WhiteList,
			Lang:               tkn.Lang,
			DefaultProfileType: tkn.DefaultProfileType,
		},
	}
	if tkn.EnergyContract != nil {
		r.Details.EnergyContract = &domain.EnergyContract{
			SupplierName: tkn.EnergyContract.SupplierName,
			ContractId:   tkn.EnergyContract.ContractId,
		}
	}
	return r
}

func (t *tokenConverter) TokenDomainToBackend(tkn *domain.Token) *backend.Token {
	if tkn == nil {
		return nil
	}
	r := &backend.Token{
		Id:                 tkn.Id,
		Type:               tkn.Details.Type,
		ContractId:         tkn.Details.ContractId,
		VisualNumber:       tkn.Details.VisualNumber,
		Issuer:             tkn.Details.Issuer,
		GroupId:            tkn.Details.GroupId,
		Valid:              tkn.Details.Valid,
		WhiteList:          tkn.Details.WhiteList,
		Lang:               tkn.Details.Lang,
		DefaultProfileType: tkn.Details.DefaultProfileType,
		LastUpdated:        tkn.LastUpdated,
		PlatformId:         tkn.PlatformId,
		RefId:              tkn.RefId,
		PartyId:            tkn.ExtId.PartyId,
		CountryCode:        tkn.ExtId.CountryCode,
	}
	if tkn.Details.EnergyContract != nil {
		r.EnergyContract = &backend.EnergyContract{
			SupplierName: tkn.Details.EnergyContract.SupplierName,
			ContractId:   tkn.Details.EnergyContract.ContractId,
		}
	}
	return r
}

func (t *tokenConverter) TokensDomainToBackend(ts []*domain.Token) []*backend.Token {
	var r []*backend.Token
	for _, tkn := range ts {
		r = append(r, t.TokenDomainToBackend(tkn))
	}
	return r
}

func (t *tokenConverter) TokenBackendToDomain(tkn *backend.Token, platformId string) *domain.Token {
	if tkn == nil {
		return nil
	}
	r := &domain.Token{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     tkn.PartyId,
				CountryCode: tkn.CountryCode,
			},
			PlatformId:  platformId,
			LastUpdated: tkn.LastUpdated,
			RefId:       tkn.RefId,
		},
		Id: tkn.Id,
		Details: domain.TokenDetails{
			Type:               tkn.Type,
			ContractId:         tkn.ContractId,
			VisualNumber:       tkn.VisualNumber,
			Issuer:             tkn.Issuer,
			GroupId:            tkn.GroupId,
			Valid:              tkn.Valid,
			WhiteList:          tkn.WhiteList,
			Lang:               tkn.Lang,
			DefaultProfileType: tkn.DefaultProfileType,
		},
	}
	if tkn.EnergyContract != nil {
		r.Details.EnergyContract = &domain.EnergyContract{
			SupplierName: tkn.EnergyContract.SupplierName,
			ContractId:   tkn.EnergyContract.ContractId,
		}
	}
	return r
}
