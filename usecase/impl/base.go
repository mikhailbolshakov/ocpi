package impl

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/errors"
	"github.com/mikhailbolshakov/ocpi/model"
	"github.com/mikhailbolshakov/ocpi/usecase"
	"strings"
)

type channelData struct {
	platformId string
	data       any
}

type ucBase struct {
	platformService domain.PlatformService
	partyService    domain.PartyService
	tokenGen        domain.TokenGenerator
}

func newBase(platformService domain.PlatformService, partyService domain.PartyService, tokenGen domain.TokenGenerator) ucBase {
	return ucBase{
		platformService: platformService,
		partyService:    partyService,
		tokenGen:        tokenGen,
	}
}

func (u *ucBase) tokenC(platform *domain.Platform) domain.PlatformToken {
	if platform.TokenBase64 != nil && *platform.TokenBase64 {
		return u.tokenGen.Base64Encode(platform.TokenC)
	}
	return platform.TokenC
}

func (u *ucBase) getConnectedPlatform(ctx context.Context, platformId string) (*domain.Platform, error) {
	platform, err := u.platformService.Get(ctx, platformId)
	if err != nil {
		return nil, err
	}
	if platform == nil {
		return nil, errors.ErrPlatformNotFound(ctx, platformId)
	}
	if platform.Status != domain.ConnectionStatusConnected {
		return nil, errors.ErrPlatformNotConnected(ctx)
	}
	return platform, nil
}

func (u *ucBase) getCreateParty(ctx context.Context, platformId, partyId, countryCode string) (*domain.Party, error) {
	party, err := u.partyService.GetByExtId(ctx, domain.PartyExtId{PartyId: partyId, CountryCode: countryCode})
	if err != nil {
		return nil, err
	}
	if party == nil {
		// if party not found, create a new party
		party, err = u.partyService.Merge(ctx, &domain.Party{
			OcpiItem: domain.OcpiItem{
				ExtId: domain.PartyExtId{
					PartyId:     partyId,
					CountryCode: countryCode,
				},
				PlatformId:  platformId,
				LastUpdated: kit.Now(),
			},
			Roles: []string{domain.RoleCPO},
		})
		if err != nil {
			return nil, err
		}
	}
	return party, nil
}

func (u *ucBase) getLocalParties(ctx context.Context, localPlatform *domain.Platform) ([]*domain.Party, error) {
	searchRq := &domain.PartySearchCriteria{
		PageRequest:  domain.PageRequest{Limit: kit.IntPtr(999)},
		IncPlatforms: []string{localPlatform.Id},
	}
	rs, err := u.partyService.Search(ctx, searchRq)
	if err != nil {
		return nil, err
	}
	return rs.Items, nil
}

func (u *ucBase) setToPartyCtx(ctx context.Context, extId domain.PartyExtId) context.Context {
	if rqCtx, ok := kit.Request(ctx); ok {
		rqCtx.WithKv(model.OcpiCtxToParty, extId.PartyId)
		rqCtx.WithKv(model.OcpiCtxToCountryCode, extId.CountryCode)
		return rqCtx.ToContext(ctx)
	}
	return ctx
}

func (u *ucBase) getFromPartyCtx(ctx context.Context) domain.PartyExtId {
	if rqCtx, ok := kit.Request(ctx); ok && rqCtx.GetKv() != nil {
		partyId, countryCode := rqCtx.Kv[model.OcpiCtxFromParty], rqCtx.Kv[model.OcpiCtxFromCountryCode]
		return domain.PartyExtId{
			PartyId:     partyId.(string),
			CountryCode: countryCode.(string),
		}
	}
	return domain.PartyExtId{}
}

func (u *ucBase) setFromPartyCtx(ctx context.Context, extId domain.PartyExtId) context.Context {
	if rqCtx, ok := kit.Request(ctx); ok {
		rqCtx.WithKv(model.OcpiCtxFromParty, extId.PartyId)
		rqCtx.WithKv(model.OcpiCtxFromCountryCode, extId.CountryCode)
		return rqCtx.ToContext(ctx)
	}
	return ctx
}

func (u *ucBase) tokenToCdrToken(tkn *domain.Token) *domain.CdrToken {
	return &domain.CdrToken{
		PartyExtId: tkn.ExtId,
		Id:         tkn.Id,
		Type:       tkn.Details.Type,
		ContractId: tkn.Details.ContractId,
	}
}

type baseConverter struct{}

func (c *baseConverter) businessDetailsModelToDomain(bd *model.OcpiBusinessDetails) *domain.BusinessDetails {
	if bd == nil {
		return nil
	}
	r := &domain.BusinessDetails{
		Name:    bd.Name,
		Website: bd.Website,
		Inn:     bd.Inn,
	}
	if bd.Logo != nil {
		r.Logo = &domain.Image{
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

func (c *baseConverter) businessDetailsBackendToDomain(bd *backend.BusinessDetails) *domain.BusinessDetails {
	if bd == nil {
		return nil
	}
	r := &domain.BusinessDetails{
		Name:    bd.Name,
		Website: bd.Website,
		Inn:     bd.Inn,
	}
	if bd.Logo != nil {
		r.Logo = &domain.Image{
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

func (c *baseConverter) businessDetailsDomainToModel(bd *domain.BusinessDetails) *model.OcpiBusinessDetails {
	if bd == nil {
		return nil
	}
	r := &model.OcpiBusinessDetails{
		Name:    bd.Name,
		Website: bd.Website,
		Inn:     bd.Inn,
	}
	if bd.Logo != nil {
		r.Logo = &model.OcpiImage{
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

func (c *baseConverter) displayTextModelToDomain(dt *model.OcpiDisplayText) *domain.DisplayText {
	if dt == nil {
		return nil
	}
	return &domain.DisplayText{
		Language: strings.ToLower(dt.Language),
		Text:     dt.Text,
	}
}

func (c *baseConverter) displayTextsModelToDomain(dts []*model.OcpiDisplayText) []*domain.DisplayText {
	var r []*domain.DisplayText
	for _, i := range dts {
		r = append(r, c.displayTextModelToDomain(i))
	}
	return r
}

func (c *baseConverter) displayTextDomainToModel(dt *domain.DisplayText) *model.OcpiDisplayText {
	if dt == nil {
		return nil
	}
	return &model.OcpiDisplayText{
		Language: dt.Language,
		Text:     dt.Text,
	}
}

func (c *baseConverter) displayTextsDomainToModel(dts []*domain.DisplayText) []*model.OcpiDisplayText {
	var r []*model.OcpiDisplayText
	for _, i := range dts {
		r = append(r, c.displayTextDomainToModel(i))
	}
	return r
}

func (c *baseConverter) displayTextDomainToBackend(dt *domain.DisplayText) *backend.DisplayText {
	if dt == nil {
		return nil
	}
	return &backend.DisplayText{
		Language: dt.Language,
		Text:     dt.Text,
	}
}

func (c *baseConverter) displayTextsDomainToBackend(dts []*domain.DisplayText) []*backend.DisplayText {
	var r []*backend.DisplayText
	for _, i := range dts {
		r = append(r, c.displayTextDomainToBackend(i))
	}
	return r
}

func (c *baseConverter) displayTextBackendToDomain(dt *backend.DisplayText) *domain.DisplayText {
	if dt == nil {
		return nil
	}
	return &domain.DisplayText{
		Language: strings.ToLower(dt.Language),
		Text:     dt.Text,
	}
}

func (c *baseConverter) displayTextsBackendToDomain(dts []*backend.DisplayText) []*domain.DisplayText {
	var r []*domain.DisplayText
	for _, i := range dts {
		r = append(r, c.displayTextBackendToDomain(i))
	}
	return r
}

func (c *baseConverter) priceBackendToDomain(p *backend.Price) *domain.Price {
	if p == nil {
		return nil
	}
	return &domain.Price{
		ExclVat: p.ExclVat,
		InclVat: p.InclVat,
	}
}

func (c *baseConverter) priceDomainToModel(p *domain.Price) *model.OcpiPrice {
	if p == nil {
		return nil
	}
	return &model.OcpiPrice{
		ExclVat: p.ExclVat,
		InclVat: p.InclVat,
	}
}

func (c *baseConverter) priceModelToDomain(p *model.OcpiPrice) *domain.Price {
	if p == nil {
		return nil
	}
	return &domain.Price{
		ExclVat: p.ExclVat,
		InclVat: p.InclVat,
	}
}

func (c *baseConverter) priceDomainToBackend(p *domain.Price) *backend.Price {
	if p == nil {
		return nil
	}
	return &backend.Price{
		ExclVat: p.ExclVat,
		InclVat: p.InclVat,
	}
}

func (c *baseConverter) cdrTokenDomainToModel(tkn *domain.CdrToken) *model.OcpiCdrToken {
	if tkn == nil {
		return nil
	}
	return &model.OcpiCdrToken{
		OcpiPartyId: model.OcpiPartyId{
			PartyId:     tkn.PartyId,
			CountryCode: tkn.CountryCode,
		},
		Id:         tkn.Id,
		Type:       tkn.Type,
		ContractId: tkn.ContractId,
	}
}

func (c *baseConverter) chargingPeriodDomainToModel(cps []*domain.ChargingPeriod) []*model.OcpiChargingPeriod {
	var r []*model.OcpiChargingPeriod
	for _, cp := range cps {
		r = append(r, &model.OcpiChargingPeriod{
			StartDateTime:         cp.StartDateTime,
			Dimensions:            c.dimensionsDomainToModel(cp.Dimensions),
			TariffId:              cp.TariffId,
			EncodingMethod:        cp.EncodingMethod,
			EncodingMethodVersion: cp.EncodingMethodVersion,
			PublicKey:             cp.PublicKey,
			SignedValues:          c.signedValueDomainToModel(cp.SignedValues),
			Url:                   cp.Url,
		})
	}
	return r
}

func (c *baseConverter) dimensionsDomainToModel(dms []*domain.CdrDimension) []*model.OcpiCdrDimension {
	var r []*model.OcpiCdrDimension
	for _, dm := range dms {
		r = append(r, &model.OcpiCdrDimension{
			Type:   dm.Type,
			Volume: dm.Volume,
		})
	}
	return r
}

func (c *baseConverter) signedValueDomainToModel(svs []*domain.SignedValue) []*model.OcpiSignedValue {
	var r []*model.OcpiSignedValue
	for _, sv := range svs {
		r = append(r, &model.OcpiSignedValue{
			Nature:     sv.Nature,
			PlainData:  sv.PlainData,
			SignedData: sv.SignedData,
		})
	}
	return r
}

func (c *baseConverter) cdrTokenModelToDomain(tkn *model.OcpiCdrToken) *domain.CdrToken {
	if tkn == nil {
		return nil
	}
	return &domain.CdrToken{
		PartyExtId: domain.PartyExtId{
			PartyId:     tkn.PartyId,
			CountryCode: tkn.CountryCode,
		},
		Id:         tkn.Id,
		Type:       tkn.Type,
		ContractId: tkn.ContractId,
	}
}

func (c *baseConverter) chargingPeriodModelToDomain(cps []*model.OcpiChargingPeriod) []*domain.ChargingPeriod {
	var r []*domain.ChargingPeriod
	for _, cp := range cps {
		r = append(r, &domain.ChargingPeriod{
			StartDateTime:         cp.StartDateTime,
			Dimensions:            c.dimensionsModelToDomain(cp.Dimensions),
			TariffId:              cp.TariffId,
			EncodingMethod:        cp.EncodingMethod,
			EncodingMethodVersion: cp.EncodingMethodVersion,
			PublicKey:             cp.PublicKey,
			SignedValues:          c.signedValueModelToDomain(cp.SignedValues),
			Url:                   cp.Url,
		})
	}
	return r
}

func (c *baseConverter) dimensionsModelToDomain(dms []*model.OcpiCdrDimension) []*domain.CdrDimension {
	var r []*domain.CdrDimension
	for _, dm := range dms {
		r = append(r, &domain.CdrDimension{
			Type:   dm.Type,
			Volume: dm.Volume,
		})
	}
	return r
}

func (c *baseConverter) signedValueModelToDomain(svs []*model.OcpiSignedValue) []*domain.SignedValue {
	var r []*domain.SignedValue
	for _, sv := range svs {
		r = append(r, &domain.SignedValue{
			Nature:     sv.Nature,
			PlainData:  sv.PlainData,
			SignedData: sv.SignedData,
		})
	}
	return r
}

func (c *baseConverter) cdrTokenDomainToBackend(tkn *domain.CdrToken) *backend.CdrToken {
	if tkn == nil {
		return nil
	}
	return &backend.CdrToken{
		PartyId:     tkn.PartyId,
		CountryCode: tkn.CountryCode,
		Id:          tkn.Id,
		Type:        tkn.Type,
		ContractId:  tkn.ContractId,
	}
}

func (c *baseConverter) chargingPeriodDomainToBackend(cps []*domain.ChargingPeriod) []*backend.ChargingPeriod {
	var r []*backend.ChargingPeriod
	for _, cp := range cps {
		r = append(r, &backend.ChargingPeriod{
			StartDateTime:         cp.StartDateTime,
			Dimensions:            c.dimensionsDomainToBackend(cp.Dimensions),
			TariffId:              cp.TariffId,
			EncodingMethod:        cp.EncodingMethod,
			EncodingMethodVersion: cp.EncodingMethodVersion,
			PublicKey:             cp.PublicKey,
			SignedValues:          c.signedValueDomainToBackend(cp.SignedValues),
			Url:                   cp.Url,
		})
	}
	return r
}

func (c *baseConverter) dimensionsDomainToBackend(dms []*domain.CdrDimension) []*backend.CdrDimension {
	var r []*backend.CdrDimension
	for _, dm := range dms {
		r = append(r, &backend.CdrDimension{
			Type:   dm.Type,
			Volume: dm.Volume,
		})
	}
	return r
}

func (c *baseConverter) signedValueDomainToBackend(svs []*domain.SignedValue) []*backend.SignedValue {
	var r []*backend.SignedValue
	for _, sv := range svs {
		r = append(r, &backend.SignedValue{
			Nature:     sv.Nature,
			PlainData:  sv.PlainData,
			SignedData: sv.SignedData,
		})
	}
	return r
}

func (c *baseConverter) cdrTokenBackendToDomain(tkn *backend.CdrToken) *domain.CdrToken {
	if tkn == nil {
		return nil
	}
	return &domain.CdrToken{
		PartyExtId: domain.PartyExtId{
			PartyId:     tkn.PartyId,
			CountryCode: tkn.CountryCode,
		},
		Id:         tkn.Id,
		Type:       tkn.Type,
		ContractId: tkn.ContractId,
	}
}

func (c *baseConverter) chargingPeriodBackendToDomain(cps []*backend.ChargingPeriod) []*domain.ChargingPeriod {
	var r []*domain.ChargingPeriod
	for _, cp := range cps {
		r = append(r, &domain.ChargingPeriod{
			StartDateTime:         cp.StartDateTime,
			Dimensions:            c.dimensionsBackendToDomain(cp.Dimensions),
			TariffId:              cp.TariffId,
			EncodingMethod:        cp.EncodingMethod,
			EncodingMethodVersion: cp.EncodingMethodVersion,
			PublicKey:             cp.PublicKey,
			SignedValues:          c.signedValueBackendToDomain(cp.SignedValues),
			Url:                   cp.Url,
		})
	}
	return r
}

func (c *baseConverter) dimensionsBackendToDomain(dms []*backend.CdrDimension) []*domain.CdrDimension {
	var r []*domain.CdrDimension
	for _, dm := range dms {
		r = append(r, &domain.CdrDimension{
			Type:   dm.Type,
			Volume: dm.Volume,
		})
	}
	return r
}

func (c *baseConverter) signedValueBackendToDomain(svs []*backend.SignedValue) []*domain.SignedValue {
	var r []*domain.SignedValue
	for _, sv := range svs {
		r = append(r, &domain.SignedValue{
			Nature:     sv.Nature,
			PlainData:  sv.PlainData,
			SignedData: sv.SignedData,
		})
	}
	return r
}

func (c *baseConverter) coordinatesModelToDomain(g *model.OcpiGeoLocation) *domain.GeoLocation {
	if g == nil {
		return nil
	}
	return &domain.GeoLocation{
		Latitude:  g.Latitude,
		Longitude: g.Longitude,
	}
}

func (c *baseConverter) coordinatesDomainToModel(g *domain.GeoLocation) *model.OcpiGeoLocation {
	if g == nil {
		return nil
	}
	return &model.OcpiGeoLocation{
		Latitude:  g.Latitude,
		Longitude: g.Longitude,
	}
}

func (c *baseConverter) coordinatesDomainToBackend(g *domain.GeoLocation) *backend.GeoLocation {
	if g == nil {
		return nil
	}
	return &backend.GeoLocation{
		Latitude:  g.Latitude,
		Longitude: g.Longitude,
	}
}

func (c *baseConverter) coordinatesBackendToDomain(g *backend.GeoLocation) *domain.GeoLocation {
	if g == nil {
		return nil
	}
	return &domain.GeoLocation{
		Latitude:  g.Latitude,
		Longitude: g.Longitude,
	}
}

func buildOcpiRepositoryErrHandlerRequestG[T any](ep domain.Endpoint, tkn domain.PlatformToken, fromPlatform, toPlatform *domain.Platform, rq T, l kit.CLogger) *usecase.OcpiRepositoryErrHandlerRequestG[T] {
	return &usecase.OcpiRepositoryErrHandlerRequestG[T]{
		OcpiRepositoryErrHandlerRequest: usecase.OcpiRepositoryErrHandlerRequest{
			OcpiRepositoryBaseRequest: usecase.OcpiRepositoryBaseRequest{
				Endpoint:       ep,
				Token:          tkn,
				FromPlatformId: fromPlatform.Id,
				ToPlatformId:   toPlatform.Id,
			},
			Handler: func(err error) { l.F(kit.KV{"platform": toPlatform.Id}).E(err).St().Err() },
		},
		Request: rq,
	}
}

func buildOcpiRepositoryErrHandlerRequest(ep domain.Endpoint, tkn domain.PlatformToken, fromPlatform, toPlatform *domain.Platform, l kit.CLogger) *usecase.OcpiRepositoryErrHandlerRequest {
	return &usecase.OcpiRepositoryErrHandlerRequest{
		OcpiRepositoryBaseRequest: usecase.OcpiRepositoryBaseRequest{
			Endpoint:       ep,
			Token:          tkn,
			FromPlatformId: fromPlatform.Id,
			ToPlatformId:   toPlatform.Id,
		},
		Handler: func(err error) { l.F(kit.KV{"platform": toPlatform.Id}).E(err).St().Err() },
	}
}

func buildOcpiRepositoryRequest(ep domain.Endpoint, tkn domain.PlatformToken, fromPlatform, toPlatform *domain.Platform) *usecase.OcpiRepositoryBaseRequest {
	return &usecase.OcpiRepositoryBaseRequest{
		Endpoint:       ep,
		Token:          tkn,
		FromPlatformId: fromPlatform.Id,
		ToPlatformId:   toPlatform.Id,
	}
}

func buildOcpiRepositoryRequestG[T any](ep domain.Endpoint, tkn domain.PlatformToken, fromPlatform, toPlatform *domain.Platform, rq T) *usecase.OcpiRepositoryRequestG[T] {
	return &usecase.OcpiRepositoryRequestG[T]{
		OcpiRepositoryBaseRequest: usecase.OcpiRepositoryBaseRequest{
			Endpoint:       ep,
			Token:          tkn,
			FromPlatformId: fromPlatform.Id,
			ToPlatformId:   toPlatform.Id,
		},
		Request: rq,
	}
}

func buildOcpiRepositoryIdRequest(ep domain.Endpoint, tkn domain.PlatformToken, fromPlatform, toPlatform *domain.Platform, id string) *usecase.OcpiRepositoryIdRequest {
	return &usecase.OcpiRepositoryIdRequest{
		OcpiRepositoryBaseRequest: usecase.OcpiRepositoryBaseRequest{
			Endpoint:       ep,
			Token:          tkn,
			FromPlatformId: fromPlatform.Id,
			ToPlatformId:   toPlatform.Id,
		},
		Id: id,
	}
}
