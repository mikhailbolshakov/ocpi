package impl

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/model"
	"github.com/mikhailbolshakov/ocpi/usecase"
)

type tariffConverter struct {
	baseConverter
}

func NewTariffConverter() usecase.TariffConverter {
	return &tariffConverter{}
}

func (t *tariffConverter) TariffDomainToModel(trf *domain.Tariff) *model.OcpiTariff {
	if trf == nil {
		return nil
	}
	return &model.OcpiTariff{
		OcpiPartyId: model.OcpiPartyId{
			PartyId:     trf.OcpiItem.ExtId.PartyId,
			CountryCode: trf.OcpiItem.ExtId.CountryCode,
		},
		Id:            trf.Id,
		Currency:      trf.Details.Currency,
		Type:          trf.Details.Type,
		TariffAltText: t.displayTextsDomainToModel(trf.Details.TariffAltText),
		TariffAltUrl:  trf.Details.TariffAltUrl,
		MinPrice:      t.priceDomainToModel(trf.Details.MinPrice),
		MaxPrice:      t.priceDomainToModel(trf.Details.MaxPrice),
		Elements:      t.elementsDomainToModel(trf.Details.Elements),
		StartDateTime: trf.Details.StartDateTime,
		EndDateTime:   trf.Details.EndDateTime,
		EnergyMix:     t.energyMixDomainToModel(trf.Details.EnergyMix),
		LastUpdated:   trf.LastUpdated,
	}
}

func (t *tariffConverter) TariffsDomainToModel(ts []*domain.Tariff) []*model.OcpiTariff {
	return kit.Select(ts, t.TariffDomainToModel)
}

func (t *tariffConverter) TariffModelToDomain(trf *model.OcpiTariff, platformId string) *domain.Tariff {
	if trf == nil {
		return nil
	}
	return &domain.Tariff{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     trf.PartyId,
				CountryCode: trf.CountryCode,
			},
			PlatformId:  platformId,
			LastUpdated: trf.LastUpdated,
		},
		Id: trf.Id,
		Details: domain.TariffDetails{
			Currency:      trf.Currency,
			Type:          trf.Type,
			TariffAltText: t.displayTextsModelToDomain(trf.TariffAltText),
			TariffAltUrl:  trf.TariffAltUrl,
			MinPrice:      t.priceModelToDomain(trf.MinPrice),
			MaxPrice:      t.priceModelToDomain(trf.MaxPrice),
			Elements:      t.elementsModelToDomain(trf.Elements),
			StartDateTime: trf.StartDateTime,
			EndDateTime:   trf.EndDateTime,
			EnergyMix:     t.energyMixModelToDomain(trf.EnergyMix),
		},
	}
}

func (t *tariffConverter) TariffsModelToDomain(ts []*model.OcpiTariff, platformId string) []*domain.Tariff {
	var r []*domain.Tariff
	for _, i := range ts {
		r = append(r, t.TariffModelToDomain(i, platformId))
	}
	return r
}

func (t *tariffConverter) TariffDomainToBackend(trf *domain.Tariff) *backend.Tariff {
	if trf == nil {
		return nil
	}
	return &backend.Tariff{
		Id:            trf.Id,
		Currency:      trf.Details.Currency,
		Type:          trf.Details.Type,
		TariffAltText: t.displayTextsDomainToBackend(trf.Details.TariffAltText),
		TariffAltUrl:  trf.Details.TariffAltUrl,
		MinPrice:      t.priceDomainToBackend(trf.Details.MinPrice),
		MaxPrice:      t.priceDomainToBackend(trf.Details.MaxPrice),
		Elements:      t.elementsDomainToBackend(trf.Details.Elements),
		StartDateTime: trf.Details.StartDateTime,
		EndDateTime:   trf.Details.EndDateTime,
		EnergyMix:     t.energyMixDomainToBackend(trf.Details.EnergyMix),
		LastUpdated:   trf.LastUpdated,
		PlatformId:    trf.PlatformId,
		RefId:         trf.RefId,
		PartyId:       trf.OcpiItem.ExtId.PartyId,
		CountryCode:   trf.OcpiItem.ExtId.CountryCode,
	}
}

func (t *tariffConverter) TariffsDomainToBackend(ts []*domain.Tariff) []*backend.Tariff {
	return kit.Select(ts, t.TariffDomainToBackend)
}

func (t *tariffConverter) TariffBackendToDomain(trf *backend.Tariff) *domain.Tariff {
	if trf == nil {
		return nil
	}
	return &domain.Tariff{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     trf.PartyId,
				CountryCode: trf.CountryCode,
			},
			PlatformId:  trf.PlatformId,
			RefId:       trf.RefId,
			LastUpdated: trf.LastUpdated,
		},
		Id: trf.Id,
		Details: domain.TariffDetails{
			Currency:      trf.Currency,
			Type:          trf.Type,
			TariffAltText: t.displayTextsBackendToDomain(trf.TariffAltText),
			TariffAltUrl:  trf.TariffAltUrl,
			MinPrice:      t.priceBackendToDomain(trf.MinPrice),
			MaxPrice:      t.priceBackendToDomain(trf.MaxPrice),
			Elements:      t.elementsBackendToDomain(trf.Elements),
			StartDateTime: trf.StartDateTime,
			EndDateTime:   trf.EndDateTime,
			EnergyMix:     t.energyMixBackendToDomain(trf.EnergyMix),
		},
	}
}

func (t *tariffConverter) TariffsBackendToDomain(ts []*backend.Tariff) []*domain.Tariff {
	return kit.Select(ts, t.TariffBackendToDomain)
}

func (t *tariffConverter) priceDomainToModel(p *domain.Price) *model.OcpiPrice {
	if p == nil {
		return nil
	}
	return &model.OcpiPrice{
		ExclVat: p.ExclVat,
		InclVat: p.InclVat,
	}
}

func (t *tariffConverter) energyMixDomainToModel(em *domain.EnergyMix) *model.OcpiEnergyMix {
	if em == nil {
		return nil
	}
	r := &model.OcpiEnergyMix{
		IsGreenEnergy:     em.IsGreenEnergy,
		EnvironImpact:     nil,
		SupplierName:      em.SupplierName,
		EnergyProductName: em.EnergyProductName,
	}
	for _, es := range em.EnergySources {
		r.EnergySources = append(r.EnergySources, &model.OcpiEnergySource{
			Source:     es.Source,
			Percentage: es.Percentage,
		})
	}
	for _, imp := range em.EnvironImpact {
		r.EnvironImpact = append(r.EnvironImpact, &model.OcpiEnvironmentalImpact{
			Category: imp.Category,
			Amount:   imp.Amount,
		})
	}
	return r
}

func (t *tariffConverter) priceComponentsDomainToModel(pc []*domain.PriceComponent) []*model.OcpiPriceComponent {
	var r []*model.OcpiPriceComponent
	for _, c := range pc {
		r = append(r, &model.OcpiPriceComponent{
			Type:     c.Type,
			Price:    c.Price,
			Vat:      c.Vat,
			StepSize: c.StepSize,
		})
	}
	return r
}

func (t *tariffConverter) restrictionsDomainToModel(rs *domain.TariffRestrictions) *model.OcpiTariffRestrictions {
	if rs == nil {
		return nil
	}
	return &model.OcpiTariffRestrictions{
		StartTime:   rs.StartTime,
		EndTime:     rs.EndTime,
		StartDate:   rs.StartDate,
		EndDate:     rs.EndDate,
		MinKwh:      rs.MinKwh,
		MaxKwh:      rs.MaxKwh,
		MinCurrent:  rs.MinCurrent,
		MaxCurrent:  rs.MaxCurrent,
		MinPower:    rs.MinPower,
		MaxPower:    rs.MaxPower,
		MinDuration: rs.MinDuration,
		MaxDuration: rs.MaxDuration,
		DayOfWeek:   rs.DayOfWeek,
		Reservation: rs.Reservation,
	}
}

func (t *tariffConverter) elementsDomainToModel(te []*domain.TariffElement) []*model.OcpiTariffElement {
	var r []*model.OcpiTariffElement
	for _, e := range te {
		r = append(r, &model.OcpiTariffElement{
			PriceComponents: t.priceComponentsDomainToModel(e.PriceComponents),
			Restrictions:    t.restrictionsDomainToModel(e.Restrictions),
		})
	}
	return r
}

func (t *tariffConverter) priceModelToDomain(p *model.OcpiPrice) *domain.Price {
	if p == nil {
		return nil
	}
	return &domain.Price{
		ExclVat: p.ExclVat,
		InclVat: p.InclVat,
	}
}

func (t *tariffConverter) energyMixModelToDomain(em *model.OcpiEnergyMix) *domain.EnergyMix {
	if em == nil {
		return nil
	}
	r := &domain.EnergyMix{
		IsGreenEnergy:     em.IsGreenEnergy,
		EnvironImpact:     nil,
		SupplierName:      em.SupplierName,
		EnergyProductName: em.EnergyProductName,
	}
	for _, es := range em.EnergySources {
		r.EnergySources = append(r.EnergySources, &domain.EnergySource{
			Source:     es.Source,
			Percentage: es.Percentage,
		})
	}
	for _, imp := range em.EnvironImpact {
		r.EnvironImpact = append(r.EnvironImpact, &domain.EnvironmentalImpact{
			Category: imp.Category,
			Amount:   imp.Amount,
		})
	}
	return r
}

func (t *tariffConverter) priceComponentsModelToDomain(pc []*model.OcpiPriceComponent) []*domain.PriceComponent {
	var r []*domain.PriceComponent
	for _, c := range pc {
		r = append(r, &domain.PriceComponent{
			Type:     c.Type,
			Price:    c.Price,
			Vat:      c.Vat,
			StepSize: c.StepSize,
		})
	}
	return r
}

func (t *tariffConverter) restrictionsModelToDomain(rs *model.OcpiTariffRestrictions) *domain.TariffRestrictions {
	if rs == nil {
		return nil
	}
	return &domain.TariffRestrictions{
		StartTime:   rs.StartTime,
		EndTime:     rs.EndTime,
		StartDate:   rs.StartDate,
		EndDate:     rs.EndDate,
		MinKwh:      rs.MinKwh,
		MaxKwh:      rs.MaxKwh,
		MinCurrent:  rs.MinCurrent,
		MaxCurrent:  rs.MaxCurrent,
		MinPower:    rs.MinPower,
		MaxPower:    rs.MaxPower,
		MinDuration: rs.MinDuration,
		MaxDuration: rs.MaxDuration,
		DayOfWeek:   rs.DayOfWeek,
		Reservation: rs.Reservation,
	}
}

func (t *tariffConverter) elementsModelToDomain(te []*model.OcpiTariffElement) []*domain.TariffElement {
	var r []*domain.TariffElement
	for _, e := range te {
		r = append(r, &domain.TariffElement{
			PriceComponents: t.priceComponentsModelToDomain(e.PriceComponents),
			Restrictions:    t.restrictionsModelToDomain(e.Restrictions),
		})
	}
	return r
}

func (t *tariffConverter) priceDomainToBackend(p *domain.Price) *backend.Price {
	if p == nil {
		return nil
	}
	return &backend.Price{
		ExclVat: p.ExclVat,
		InclVat: p.InclVat,
	}
}

func (t *tariffConverter) energyMixDomainToBackend(em *domain.EnergyMix) *backend.EnergyMix {
	if em == nil {
		return nil
	}
	r := &backend.EnergyMix{
		IsGreenEnergy:     em.IsGreenEnergy,
		SupplierName:      em.SupplierName,
		EnergyProductName: em.EnergyProductName,
	}
	for _, es := range em.EnergySources {
		r.EnergySources = append(r.EnergySources, &backend.EnergySource{
			Source:     es.Source,
			Percentage: es.Percentage,
		})
	}
	for _, imp := range em.EnvironImpact {
		r.EnvironImpact = append(r.EnvironImpact, &backend.EnvironmentalImpact{
			Category: imp.Category,
			Amount:   imp.Amount,
		})
	}

	return r
}

func (t *tariffConverter) priceComponentsDomainToBackend(pc []*domain.PriceComponent) []*backend.PriceComponent {
	var r []*backend.PriceComponent
	for _, c := range pc {
		r = append(r, &backend.PriceComponent{
			Type:     c.Type,
			Price:    c.Price,
			Vat:      c.Vat,
			StepSize: c.StepSize,
		})
	}
	return r
}

func (t *tariffConverter) restrictionsDomainToBackend(rs *domain.TariffRestrictions) *backend.TariffRestrictions {
	if rs == nil {
		return nil
	}
	return &backend.TariffRestrictions{
		StartTime:   rs.StartTime,
		EndTime:     rs.EndTime,
		StartDate:   rs.StartDate,
		EndDate:     rs.EndDate,
		MinKwh:      rs.MinKwh,
		MaxKwh:      rs.MaxKwh,
		MinCurrent:  rs.MinCurrent,
		MaxCurrent:  rs.MaxCurrent,
		MinPower:    rs.MinPower,
		MaxPower:    rs.MaxPower,
		MinDuration: rs.MinDuration,
		MaxDuration: rs.MaxDuration,
		DayOfWeek:   rs.DayOfWeek,
		Reservation: rs.Reservation,
	}
}

func (t *tariffConverter) elementsDomainToBackend(te []*domain.TariffElement) []*backend.TariffElement {
	var r []*backend.TariffElement
	for _, e := range te {
		r = append(r, &backend.TariffElement{
			PriceComponents: t.priceComponentsDomainToBackend(e.PriceComponents),
			Restrictions:    t.restrictionsDomainToBackend(e.Restrictions),
		})
	}
	return r
}

func (t *tariffConverter) priceBackendToDomain(p *backend.Price) *domain.Price {
	if p == nil {
		return nil
	}
	return &domain.Price{
		ExclVat: p.ExclVat,
		InclVat: p.InclVat,
	}
}

func (t *tariffConverter) energyMixBackendToDomain(em *backend.EnergyMix) *domain.EnergyMix {
	if em == nil {
		return nil
	}
	r := &domain.EnergyMix{
		IsGreenEnergy:     em.IsGreenEnergy,
		SupplierName:      em.SupplierName,
		EnergyProductName: em.EnergyProductName,
	}
	for _, es := range em.EnergySources {
		r.EnergySources = append(r.EnergySources, &domain.EnergySource{
			Source:     es.Source,
			Percentage: es.Percentage,
		})
	}
	for _, imp := range em.EnvironImpact {
		r.EnvironImpact = append(r.EnvironImpact, &domain.EnvironmentalImpact{
			Category: imp.Category,
			Amount:   imp.Amount,
		})
	}
	return r
}

func (t *tariffConverter) priceComponentsBackendToDomain(pc []*backend.PriceComponent) []*domain.PriceComponent {
	var r []*domain.PriceComponent
	for _, c := range pc {
		r = append(r, &domain.PriceComponent{
			Type:     c.Type,
			Price:    c.Price,
			Vat:      c.Vat,
			StepSize: c.StepSize,
		})
	}
	return r
}

func (t *tariffConverter) restrictionsBackendToDomain(rs *backend.TariffRestrictions) *domain.TariffRestrictions {
	if rs == nil {
		return nil
	}
	return &domain.TariffRestrictions{
		StartTime:   rs.StartTime,
		EndTime:     rs.EndTime,
		StartDate:   rs.StartDate,
		EndDate:     rs.EndDate,
		MinKwh:      rs.MinKwh,
		MaxKwh:      rs.MaxKwh,
		MinCurrent:  rs.MinCurrent,
		MaxCurrent:  rs.MaxCurrent,
		MinPower:    rs.MinPower,
		MaxPower:    rs.MaxPower,
		MinDuration: rs.MinDuration,
		MaxDuration: rs.MaxDuration,
		DayOfWeek:   rs.DayOfWeek,
		Reservation: rs.Reservation,
	}
}

func (t *tariffConverter) elementsBackendToDomain(te []*backend.TariffElement) []*domain.TariffElement {
	var r []*domain.TariffElement
	for _, e := range te {
		r = append(r, &domain.TariffElement{
			PriceComponents: t.priceComponentsBackendToDomain(e.PriceComponents),
			Restrictions:    t.restrictionsBackendToDomain(e.Restrictions),
		})
	}
	return r
}
