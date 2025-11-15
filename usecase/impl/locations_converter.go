package impl

import (
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/model"
	"github.com/mikhailbolshakov/ocpi/usecase"
)

type locationConverter struct {
	baseConverter
}

func NewLocationConverter() usecase.LocationConverter {
	return &locationConverter{}
}

func (l *locationConverter) ConnectorModelToDomain(con *model.OcpiConnector, cc, partyId, platformId, locId, evseId string) *domain.Connector {
	if con == nil {
		return nil
	}
	return &domain.Connector{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     partyId,
				CountryCode: cc,
			},
			PlatformId:  platformId,
			LastUpdated: con.LastUpdated,
		},
		Id:         con.Id,
		LocationId: locId,
		EvseId:     evseId,
		Details: domain.ConnectorDetails{
			Standard:           con.Standard,
			Format:             con.Format,
			PowerType:          con.PowerType,
			MaxVoltage:         con.MaxVoltage,
			MaxAmperage:        con.MaxAmperage,
			MaxElectricPower:   con.MaxElectricPower,
			TariffIds:          con.TariffIds,
			TermsAndConditions: con.TermsAndConditions,
		},
	}
}

func (l *locationConverter) ConnectorsModelToDomain(cons []*model.OcpiConnector, cc, partyId, platformId, locId, evseId string) []*domain.Connector {
	var r []*domain.Connector
	for _, c := range cons {
		r = append(r, l.ConnectorModelToDomain(c, cc, partyId, platformId, locId, evseId))
	}
	return r
}

func (l *locationConverter) statusScheduleModelToDomain(ss []*model.OcpiStatusSchedule) []*domain.StatusSchedule {
	var r []*domain.StatusSchedule
	for _, i := range ss {
		r = append(r, &domain.StatusSchedule{
			PeriodBegin: i.PeriodBegin,
			PeriodEnd:   i.PeriodEnd,
			Status:      i.Status,
		})
	}
	return r
}

func (l *locationConverter) imagesModelToDomain(dts []*model.OcpiImage) []*domain.Image {
	var r []*domain.Image
	for _, i := range dts {
		r = append(r, &domain.Image{
			Url:       i.Url,
			Thumbnail: i.Thumbnail,
			Category:  i.Category,
			Type:      i.Type,
			Width:     i.Width,
			Height:    i.Height,
		})
	}
	return r
}

func (l *locationConverter) EvseModelToDomain(e *model.OcpiEvse, cc, partyId, platformId, locId string) *domain.Evse {
	if e == nil {
		return nil
	}
	return &domain.Evse{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     partyId,
				CountryCode: cc,
			},
			PlatformId:  platformId,
			LastUpdated: e.LastUpdated,
		},
		Id:         e.Uid,
		LocationId: locId,
		Status:     e.Status,
		Details: domain.EvseDetails{
			EvseId:              e.EvseId,
			StatusSchedule:      l.statusScheduleModelToDomain(e.StatusSchedule),
			Capabilities:        e.Capabilities,
			FloorLevel:          e.FloorLevel,
			Coordinates:         l.coordinatesModelToDomain(e.Coordinates),
			PhysicalReference:   e.PhysicalReference,
			Directions:          l.displayTextsModelToDomain(e.Directions),
			ParkingRestrictions: e.ParkingRestrictions,
			Images:              l.imagesModelToDomain(e.Images),
		},
		Connectors: l.ConnectorsModelToDomain(e.Connectors, cc, partyId, platformId, locId, e.Uid),
	}
}

func (l *locationConverter) EvsesModelToDomain(evses []*model.OcpiEvse, cc, partyId, platformId, locId string) []*domain.Evse {
	var r []*domain.Evse
	for _, e := range evses {
		r = append(r, l.EvseModelToDomain(e, cc, partyId, platformId, locId))
	}
	return r
}

func (l *locationConverter) publishTokenTypesModelToDomain(tt []*model.OcpiPublishTokenType) []*domain.PublishTokenType {
	var r []*domain.PublishTokenType
	for _, t := range tt {
		r = append(r, &domain.PublishTokenType{
			Uid:          t.Uid,
			Type:         t.Type,
			VisualNumber: t.VisualNumber,
			Issuer:       t.Issuer,
			GroupId:      t.GroupId,
		})
	}
	return r
}

func (l *locationConverter) additionalGetLocationsModelToDomain(gl []*model.OcpiAdditionalGeoLocation) []*domain.AdditionalGeoLocation {
	var r []*domain.AdditionalGeoLocation
	for _, i := range gl {
		r = append(r, &domain.AdditionalGeoLocation{
			GeoLocation: domain.GeoLocation{
				Latitude:  i.Latitude,
				Longitude: i.Longitude,
			},
			Name: l.displayTextModelToDomain(i.Name),
		})
	}
	return r
}

func (l *locationConverter) exceptionPeriodsModelToDomain(ep []*model.OcpiExceptionalPeriod) []*domain.ExceptionalPeriod {
	var r []*domain.ExceptionalPeriod
	for _, e := range ep {
		r = append(r, &domain.ExceptionalPeriod{
			PeriodBegin: e.PeriodBegin,
			PeriodEnd:   e.PeriodEnd,
		})
	}
	return r
}

func (l *locationConverter) hoursModelToDomain(h *model.OcpiHours) *domain.Hours {
	if h == nil {
		return nil
	}
	r := &domain.Hours{
		TwentyFourSeven:     h.TwentyFourSeven,
		ExceptionalOpenings: l.exceptionPeriodsModelToDomain(h.ExceptionalOpenings),
		ExceptionalClosings: l.exceptionPeriodsModelToDomain(h.ExceptionalClosings),
	}
	for _, rh := range h.RegularHours {
		r.RegularHours = append(r.RegularHours, &domain.RegularHours{
			Weekday:     int(rh.Weekday),
			PeriodBegin: rh.PeriodBegin,
			PeriodEnd:   rh.PeriodEnd,
		})
	}
	return r
}

func (l *locationConverter) energyMixModelToDomain(em *model.OcpiEnergyMix) *domain.EnergyMix {
	if em == nil {
		return nil
	}
	r := &domain.EnergyMix{
		IsGreenEnergy:     em.IsGreenEnergy,
		SupplierName:      em.SupplierName,
		EnergyProductName: em.EnergyProductName,
	}
	for _, imp := range em.EnvironImpact {
		r.EnvironImpact = append(r.EnvironImpact, &domain.EnvironmentalImpact{
			Category: imp.Category,
			Amount:   imp.Amount,
		})
	}
	for _, es := range em.EnergySources {
		r.EnergySources = append(r.EnergySources, &domain.EnergySource{
			Source:     es.Source,
			Percentage: es.Percentage,
		})
	}
	return r
}

func (l *locationConverter) LocationModelToDomain(loc *model.OcpiLocation, platformId string) *domain.Location {
	if loc == nil {
		return nil
	}
	return &domain.Location{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     loc.PartyId,
				CountryCode: loc.CountryCode,
			},
			PlatformId:  platformId,
			LastUpdated: loc.LastUpdated,
		},
		Id: loc.Id,
		Details: domain.LocationDetails{
			Publish:            loc.Publish,
			PublishAllowedTo:   l.publishTokenTypesModelToDomain(loc.PublishAllowedTo),
			Name:               loc.Name,
			Address:            loc.Address,
			City:               loc.City,
			PostalCode:         loc.PostalCode,
			State:              loc.State,
			Country:            loc.Country,
			Coordinates:        *l.coordinatesModelToDomain(&loc.Coordinates),
			RelatedLocations:   l.additionalGetLocationsModelToDomain(loc.RelatedLocations),
			ParkingType:        loc.ParkingType,
			Directions:         l.displayTextsModelToDomain(loc.Directions),
			Operator:           l.businessDetailsModelToDomain(loc.Operator),
			SubOperator:        l.businessDetailsModelToDomain(loc.SubOperator),
			Owner:              l.businessDetailsModelToDomain(loc.Owner),
			Facilities:         loc.Facilities,
			TimeZone:           loc.TimeZone,
			OpeningTimes:       l.hoursModelToDomain(loc.OpeningTimes),
			ChargingWhenClosed: loc.ChargingWhenClosed,
			Images:             l.imagesModelToDomain(loc.Images),
			EnergyMix:          l.energyMixModelToDomain(loc.EnergyMix),
		},
		Evses: l.EvsesModelToDomain(loc.Evses, loc.CountryCode, loc.PartyId, platformId, loc.Id),
	}
}

func (l *locationConverter) statusScheduleDomainToModel(ss []*domain.StatusSchedule) []*model.OcpiStatusSchedule {
	var r []*model.OcpiStatusSchedule
	for _, i := range ss {
		r = append(r, &model.OcpiStatusSchedule{
			PeriodBegin: i.PeriodBegin,
			PeriodEnd:   i.PeriodEnd,
			Status:      i.Status,
		})
	}
	return r
}

func (l *locationConverter) imagesDomainToModelModel(dts []*domain.Image) []*model.OcpiImage {
	var r []*model.OcpiImage
	for _, i := range dts {
		r = append(r, &model.OcpiImage{
			Url:       i.Url,
			Thumbnail: i.Thumbnail,
			Category:  i.Category,
			Type:      i.Type,
			Width:     i.Width,
			Height:    i.Height,
		})
	}
	return r
}

func (l *locationConverter) publishTokenTypesDomainToModel(tt []*domain.PublishTokenType) []*model.OcpiPublishTokenType {
	var r []*model.OcpiPublishTokenType
	for _, t := range tt {
		r = append(r, &model.OcpiPublishTokenType{
			Uid:          t.Uid,
			Type:         t.Type,
			VisualNumber: t.VisualNumber,
			Issuer:       t.Issuer,
			GroupId:      t.GroupId,
		})
	}
	return r
}

func (l *locationConverter) additionalGetLocationsDomainToModel(gl []*domain.AdditionalGeoLocation) []*model.OcpiAdditionalGeoLocation {
	var r []*model.OcpiAdditionalGeoLocation
	for _, i := range gl {
		r = append(r, &model.OcpiAdditionalGeoLocation{
			OcpiGeoLocation: model.OcpiGeoLocation{
				Latitude:  i.Latitude,
				Longitude: i.Longitude,
			},
			Name: l.displayTextDomainToModel(i.Name),
		})
	}
	return r
}

func (l *locationConverter) exceptionPeriodsDomainToModel(ep []*domain.ExceptionalPeriod) []*model.OcpiExceptionalPeriod {
	var r []*model.OcpiExceptionalPeriod
	for _, e := range ep {
		r = append(r, &model.OcpiExceptionalPeriod{
			PeriodBegin: e.PeriodBegin,
			PeriodEnd:   e.PeriodEnd,
		})
	}
	return r
}

func (l *locationConverter) hoursDomainToModel(h *domain.Hours) *model.OcpiHours {
	if h == nil {
		return nil
	}
	r := &model.OcpiHours{
		TwentyFourSeven:     h.TwentyFourSeven,
		ExceptionalOpenings: l.exceptionPeriodsDomainToModel(h.ExceptionalOpenings),
		ExceptionalClosings: l.exceptionPeriodsDomainToModel(h.ExceptionalClosings),
	}
	for _, rh := range h.RegularHours {
		r.RegularHours = append(r.RegularHours, &model.OcpiRegularHours{
			Weekday:     rh.Weekday,
			PeriodBegin: rh.PeriodBegin,
			PeriodEnd:   rh.PeriodEnd,
		})
	}
	return r
}

func (l *locationConverter) energyMixDomainToModel(em *domain.EnergyMix) *model.OcpiEnergyMix {
	if em == nil {
		return nil
	}
	r := &model.OcpiEnergyMix{
		IsGreenEnergy:     em.IsGreenEnergy,
		EnergySources:     nil,
		EnvironImpact:     nil,
		SupplierName:      em.SupplierName,
		EnergyProductName: em.EnergyProductName,
	}
	for _, imp := range em.EnvironImpact {
		r.EnvironImpact = append(r.EnvironImpact, &model.OcpiEnvironmentalImpact{
			Category: imp.Category,
			Amount:   imp.Amount,
		})
	}
	for _, es := range em.EnergySources {
		r.EnergySources = append(r.EnergySources, &model.OcpiEnergySource{
			Source:     es.Source,
			Percentage: es.Percentage,
		})
	}
	return r
}

func (l *locationConverter) LocationDomainToModel(loc *domain.Location) *model.OcpiLocation {
	if loc == nil {
		return nil
	}
	return &model.OcpiLocation{
		OcpiPartyId: model.OcpiPartyId{
			PartyId:     loc.ExtId.PartyId,
			CountryCode: loc.ExtId.CountryCode,
		},
		Id:                 loc.Id,
		Publish:            loc.Details.Publish,
		PublishAllowedTo:   l.publishTokenTypesDomainToModel(loc.Details.PublishAllowedTo),
		Name:               loc.Details.Name,
		Address:            loc.Details.Address,
		City:               loc.Details.City,
		PostalCode:         loc.Details.PostalCode,
		State:              loc.Details.State,
		Country:            loc.Details.Country,
		Coordinates:        *l.coordinatesDomainToModel(&loc.Details.Coordinates),
		RelatedLocations:   l.additionalGetLocationsDomainToModel(loc.Details.RelatedLocations),
		ParkingType:        loc.Details.ParkingType,
		Directions:         l.displayTextsDomainToModel(loc.Details.Directions),
		Operator:           l.businessDetailsDomainToModel(loc.Details.Operator),
		SubOperator:        l.businessDetailsDomainToModel(loc.Details.SubOperator),
		Owner:              l.businessDetailsDomainToModel(loc.Details.Owner),
		Facilities:         loc.Details.Facilities,
		TimeZone:           loc.Details.TimeZone,
		OpeningTimes:       l.hoursDomainToModel(loc.Details.OpeningTimes),
		ChargingWhenClosed: loc.Details.ChargingWhenClosed,
		Images:             l.imagesDomainToModelModel(loc.Details.Images),
		EnergyMix:          l.energyMixDomainToModel(loc.Details.EnergyMix),
		Evses:              l.EvsesDomainToModel(loc.Evses),
		LastUpdated:        loc.LastUpdated,
	}
}

func (l *locationConverter) LocationsDomainToModel(locs []*domain.Location) []*model.OcpiLocation {
	var r []*model.OcpiLocation
	for _, loc := range locs {
		r = append(r, l.LocationDomainToModel(loc))
	}
	return r
}

func (l *locationConverter) EvseDomainToModel(evse *domain.Evse) *model.OcpiEvse {
	if evse == nil {
		return nil
	}
	return &model.OcpiEvse{
		Uid:                 evse.Id,
		Status:              evse.Status,
		EvseId:              evse.Details.EvseId,
		StatusSchedule:      l.statusScheduleDomainToModel(evse.Details.StatusSchedule),
		Capabilities:        evse.Details.Capabilities,
		FloorLevel:          evse.Details.FloorLevel,
		Coordinates:         l.coordinatesDomainToModel(evse.Details.Coordinates),
		PhysicalReference:   evse.Details.PhysicalReference,
		Directions:          l.displayTextsDomainToModel(evse.Details.Directions),
		ParkingRestrictions: evse.Details.ParkingRestrictions,
		Images:              l.imagesDomainToModelModel(evse.Details.Images),
		Connectors:          l.ConnectorsDomainToModel(evse.Connectors),
		LastUpdated:         evse.LastUpdated,
	}
}

func (l *locationConverter) EvsesDomainToModel(evses []*domain.Evse) []*model.OcpiEvse {
	var r []*model.OcpiEvse
	for _, e := range evses {
		r = append(r, l.EvseDomainToModel(e))
	}
	return r
}

func (l *locationConverter) ConnectorDomainToModel(con *domain.Connector) *model.OcpiConnector {
	if con == nil {
		return nil
	}
	return &model.OcpiConnector{
		Id:                 con.Id,
		Standard:           con.Details.Standard,
		Format:             con.Details.Format,
		PowerType:          con.Details.PowerType,
		MaxVoltage:         con.Details.MaxVoltage,
		MaxAmperage:        con.Details.MaxAmperage,
		MaxElectricPower:   con.Details.MaxElectricPower,
		TariffIds:          con.Details.TariffIds,
		TermsAndConditions: con.Details.TermsAndConditions,
		LastUpdated:        con.LastUpdated,
	}
}

func (l *locationConverter) ConnectorsDomainToModel(cons []*domain.Connector) []*model.OcpiConnector {
	var r []*model.OcpiConnector
	for _, c := range cons {
		r = append(r, l.ConnectorDomainToModel(c))
	}
	return r
}

func (l *locationConverter) statusScheduleDomainToBackend(ss []*domain.StatusSchedule) []*backend.StatusSchedule {
	var r []*backend.StatusSchedule
	for _, i := range ss {
		r = append(r, &backend.StatusSchedule{
			PeriodBegin: i.PeriodBegin,
			PeriodEnd:   i.PeriodEnd,
			Status:      i.Status,
		})
	}
	return r
}

func (l *locationConverter) displayTextDomainToBackend(dt *domain.DisplayText) *backend.DisplayText {
	if dt == nil {
		return nil
	}
	return &backend.DisplayText{
		Language: dt.Language,
		Text:     dt.Text,
	}
}

func (l *locationConverter) displayTextsDomainToBackend(dts []*domain.DisplayText) []*backend.DisplayText {
	var r []*backend.DisplayText
	for _, i := range dts {
		r = append(r, l.displayTextDomainToBackend(i))
	}
	return r
}

func (l *locationConverter) imagesDomainToBackend(dts []*domain.Image) []*backend.Image {
	var r []*backend.Image
	for _, i := range dts {
		r = append(r, &backend.Image{
			Url:       i.Url,
			Thumbnail: i.Thumbnail,
			Category:  i.Category,
			Type:      i.Type,
			Width:     i.Width,
			Height:    i.Height,
		})
	}
	return r
}

func (l *locationConverter) publishTokenTypesDomainToBackend(tt []*domain.PublishTokenType) []*backend.PublishTokenType {
	var r []*backend.PublishTokenType
	for _, t := range tt {
		r = append(r, &backend.PublishTokenType{
			Uid:          t.Uid,
			Type:         t.Type,
			VisualNumber: t.VisualNumber,
			Issuer:       t.Issuer,
			GroupId:      t.GroupId,
		})
	}
	return r
}

func (l *locationConverter) additionalGetLocationsDomainToBackend(gl []*domain.AdditionalGeoLocation) []*backend.AdditionalGeoLocation {
	var r []*backend.AdditionalGeoLocation
	for _, i := range gl {
		r = append(r, &backend.AdditionalGeoLocation{
			Latitude:  i.Latitude,
			Longitude: i.Longitude,
			Name:      l.displayTextDomainToBackend(i.Name),
		})
	}
	return r
}

func (l *locationConverter) exceptionPeriodsDomainToBackend(ep []*domain.ExceptionalPeriod) []*backend.ExceptionalPeriod {
	var r []*backend.ExceptionalPeriod
	for _, e := range ep {
		r = append(r, &backend.ExceptionalPeriod{
			PeriodBegin: e.PeriodBegin,
			PeriodEnd:   e.PeriodEnd,
		})
	}
	return r
}

func (l *locationConverter) hoursDomainToBackend(h *domain.Hours) *backend.Hours {
	if h == nil {
		return nil
	}
	r := &backend.Hours{
		TwentyFourSeven:     h.TwentyFourSeven,
		ExceptionalOpenings: l.exceptionPeriodsDomainToBackend(h.ExceptionalOpenings),
		ExceptionalClosings: l.exceptionPeriodsDomainToBackend(h.ExceptionalClosings),
	}
	for _, rh := range h.RegularHours {
		r.RegularHours = append(r.RegularHours, &backend.RegularHours{
			Weekday:     rh.Weekday,
			PeriodBegin: rh.PeriodBegin,
			PeriodEnd:   rh.PeriodEnd,
		})
	}
	return r
}

func (l *locationConverter) energyMixDomainToBackend(em *domain.EnergyMix) *backend.EnergyMix {
	if em == nil {
		return nil
	}
	r := &backend.EnergyMix{
		IsGreenEnergy:     em.IsGreenEnergy,
		EnergySources:     nil,
		EnvironImpact:     nil,
		SupplierName:      em.SupplierName,
		EnergyProductName: em.EnergyProductName,
	}
	for _, imp := range em.EnvironImpact {
		r.EnvironImpact = append(r.EnvironImpact, &backend.EnvironmentalImpact{
			Category: imp.Category,
			Amount:   imp.Amount,
		})
	}
	for _, es := range em.EnergySources {
		r.EnergySources = append(r.EnergySources, &backend.EnergySource{
			Source:     es.Source,
			Percentage: es.Percentage,
		})
	}
	return r
}

func (l *locationConverter) businessDetailsDomainToBackend(bd *domain.BusinessDetails) *backend.BusinessDetails {
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

func (l *locationConverter) LocationDomainToBackend(loc *domain.Location) *backend.Location {
	if loc == nil {
		return nil
	}
	return &backend.Location{
		PartyId:            loc.ExtId.PartyId,
		CountryCode:        loc.ExtId.CountryCode,
		Id:                 loc.Id,
		Publish:            loc.Details.Publish,
		PublishAllowedTo:   l.publishTokenTypesDomainToBackend(loc.Details.PublishAllowedTo),
		Name:               loc.Details.Name,
		Address:            loc.Details.Address,
		City:               loc.Details.City,
		PostalCode:         loc.Details.PostalCode,
		State:              loc.Details.State,
		Country:            loc.Details.Country,
		Coordinates:        *l.coordinatesDomainToBackend(&loc.Details.Coordinates),
		RelatedLocations:   l.additionalGetLocationsDomainToBackend(loc.Details.RelatedLocations),
		ParkingType:        loc.Details.ParkingType,
		Directions:         l.displayTextsDomainToBackend(loc.Details.Directions),
		Operator:           l.businessDetailsDomainToBackend(loc.Details.Operator),
		SubOperator:        l.businessDetailsDomainToBackend(loc.Details.SubOperator),
		Owner:              l.businessDetailsDomainToBackend(loc.Details.Owner),
		Facilities:         loc.Details.Facilities,
		TimeZone:           loc.Details.TimeZone,
		OpeningTimes:       l.hoursDomainToBackend(loc.Details.OpeningTimes),
		ChargingWhenClosed: loc.Details.ChargingWhenClosed,
		Images:             l.imagesDomainToBackend(loc.Details.Images),
		EnergyMix:          l.energyMixDomainToBackend(loc.Details.EnergyMix),
		Evses:              l.EvsesDomainToBackend(loc.Evses),
		RefId:              loc.RefId,
		LastUpdated:        loc.LastUpdated,
	}
}

func (l *locationConverter) LocationsDomainToBackend(locs []*domain.Location) []*backend.Location {
	var r []*backend.Location
	for _, loc := range locs {
		r = append(r, l.LocationDomainToBackend(loc))
	}
	return r
}

func (l *locationConverter) EvseDomainToBackend(evse *domain.Evse) *backend.Evse {
	if evse == nil {
		return nil
	}
	return &backend.Evse{
		Id:                  evse.Id,
		Status:              evse.Status,
		EvseId:              evse.Details.EvseId,
		StatusSchedule:      l.statusScheduleDomainToBackend(evse.Details.StatusSchedule),
		Capabilities:        evse.Details.Capabilities,
		FloorLevel:          evse.Details.FloorLevel,
		Coordinates:         l.coordinatesDomainToBackend(evse.Details.Coordinates),
		PhysicalReference:   evse.Details.PhysicalReference,
		Directions:          l.displayTextsDomainToBackend(evse.Details.Directions),
		ParkingRestrictions: evse.Details.ParkingRestrictions,
		Images:              l.imagesDomainToBackend(evse.Details.Images),
		Connectors:          l.ConnectorsDomainToBackend(evse.Connectors),
		RefId:               evse.RefId,
		LocationId:          evse.LocationId,
		CountryCode:         evse.ExtId.CountryCode,
		PartyId:             evse.ExtId.PartyId,
		LastUpdated:         evse.LastUpdated,
	}
}

func (l *locationConverter) EvsesDomainToBackend(evses []*domain.Evse) []*backend.Evse {
	var r []*backend.Evse
	for _, e := range evses {
		r = append(r, l.EvseDomainToBackend(e))
	}
	return r
}

func (l *locationConverter) ConnectorDomainToBackend(con *domain.Connector) *backend.Connector {
	if con == nil {
		return nil
	}
	return &backend.Connector{
		Id:                 con.Id,
		LocationId:         con.LocationId,
		EvseId:             con.EvseId,
		Standard:           con.Details.Standard,
		Format:             con.Details.Format,
		PowerType:          con.Details.PowerType,
		MaxVoltage:         con.Details.MaxVoltage,
		MaxAmperage:        con.Details.MaxAmperage,
		MaxElectricPower:   con.Details.MaxElectricPower,
		TariffIds:          con.Details.TariffIds,
		TermsAndConditions: con.Details.TermsAndConditions,
		PartyId:            con.ExtId.PartyId,
		CountryCode:        con.ExtId.CountryCode,
		RefId:              con.RefId,
		LastUpdated:        con.LastUpdated,
	}
}

func (l *locationConverter) ConnectorsDomainToBackend(cons []*domain.Connector) []*backend.Connector {
	var r []*backend.Connector
	for _, c := range cons {
		r = append(r, l.ConnectorDomainToBackend(c))
	}
	return r
}

func (l *locationConverter) ConnectorBackendToDomain(con *backend.Connector, platformId, locId, evseId string) *domain.Connector {
	if con == nil {
		return nil
	}
	return &domain.Connector{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     con.PartyId,
				CountryCode: con.CountryCode,
			},
			PlatformId:  platformId,
			RefId:       con.RefId,
			LastUpdated: con.LastUpdated,
		},
		Id:         con.Id,
		LocationId: locId,
		EvseId:     evseId,
		Details: domain.ConnectorDetails{
			Standard:           con.Standard,
			Format:             con.Format,
			PowerType:          con.PowerType,
			MaxVoltage:         con.MaxVoltage,
			MaxAmperage:        con.MaxAmperage,
			MaxElectricPower:   con.MaxElectricPower,
			TariffIds:          con.TariffIds,
			TermsAndConditions: con.TermsAndConditions,
		},
	}
}

func (l *locationConverter) ConnectorsBackendToDomain(cons []*backend.Connector, platformId, locId, evseId string) []*domain.Connector {
	var r []*domain.Connector
	for _, c := range cons {
		r = append(r, l.ConnectorBackendToDomain(c, platformId, locId, evseId))
	}
	return r
}

func (l *locationConverter) statusScheduleBackendToDomain(ss []*backend.StatusSchedule) []*domain.StatusSchedule {
	var r []*domain.StatusSchedule
	for _, i := range ss {
		r = append(r, &domain.StatusSchedule{
			PeriodBegin: i.PeriodBegin,
			PeriodEnd:   i.PeriodEnd,
			Status:      i.Status,
		})
	}
	return r
}

func (l *locationConverter) displayTextBackendToDomain(dt *backend.DisplayText) *domain.DisplayText {
	if dt == nil {
		return nil
	}
	return &domain.DisplayText{
		Language: dt.Language,
		Text:     dt.Text,
	}
}

func (l *locationConverter) displayTextsBackendToDomain(dts []*backend.DisplayText) []*domain.DisplayText {
	var r []*domain.DisplayText
	for _, i := range dts {
		r = append(r, l.displayTextBackendToDomain(i))
	}
	return r
}

func (l *locationConverter) imagesBackendToDomain(dts []*backend.Image) []*domain.Image {
	var r []*domain.Image
	for _, i := range dts {
		r = append(r, &domain.Image{
			Url:       i.Url,
			Thumbnail: i.Thumbnail,
			Category:  i.Category,
			Type:      i.Type,
			Width:     i.Width,
			Height:    i.Height,
		})
	}
	return r
}

func (l *locationConverter) EvseBackendToDomain(e *backend.Evse, platformId, locId string) *domain.Evse {
	if e == nil {
		return nil
	}
	return &domain.Evse{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     e.PartyId,
				CountryCode: e.CountryCode,
			},
			PlatformId:  platformId,
			RefId:       e.RefId,
			LastUpdated: e.LastUpdated,
		},
		Id:         e.Id,
		LocationId: locId,
		Status:     e.Status,
		Details: domain.EvseDetails{
			EvseId:              e.EvseId,
			StatusSchedule:      l.statusScheduleBackendToDomain(e.StatusSchedule),
			Capabilities:        e.Capabilities,
			FloorLevel:          e.FloorLevel,
			Coordinates:         l.coordinatesBackendToDomain(e.Coordinates),
			PhysicalReference:   e.PhysicalReference,
			Directions:          l.displayTextsBackendToDomain(e.Directions),
			ParkingRestrictions: e.ParkingRestrictions,
			Images:              l.imagesBackendToDomain(e.Images),
		},
		Connectors: l.ConnectorsBackendToDomain(e.Connectors, platformId, locId, e.Id),
	}
}

func (l *locationConverter) EvsesBackendToDomain(evses []*backend.Evse, platformId, locId string) []*domain.Evse {
	var r []*domain.Evse
	for _, e := range evses {
		r = append(r, l.EvseBackendToDomain(e, platformId, locId))
	}
	return r
}

func (l *locationConverter) publishTokenTypesBackendToDomain(tt []*backend.PublishTokenType) []*domain.PublishTokenType {
	var r []*domain.PublishTokenType
	for _, t := range tt {
		r = append(r, &domain.PublishTokenType{
			Uid:          t.Uid,
			Type:         t.Type,
			VisualNumber: t.VisualNumber,
			Issuer:       t.Issuer,
			GroupId:      t.GroupId,
		})
	}
	return r
}

func (l *locationConverter) additionalGetLocationsBackendToDomain(gl []*backend.AdditionalGeoLocation) []*domain.AdditionalGeoLocation {
	var r []*domain.AdditionalGeoLocation
	for _, i := range gl {
		r = append(r, &domain.AdditionalGeoLocation{
			GeoLocation: domain.GeoLocation{
				Latitude:  i.Latitude,
				Longitude: i.Longitude,
			},
			Name: l.displayTextBackendToDomain(i.Name),
		})
	}
	return r
}

func (l *locationConverter) exceptionPeriodsBackendToDomain(ep []*backend.ExceptionalPeriod) []*domain.ExceptionalPeriod {
	var r []*domain.ExceptionalPeriod
	for _, e := range ep {
		r = append(r, &domain.ExceptionalPeriod{
			PeriodBegin: e.PeriodBegin,
			PeriodEnd:   e.PeriodEnd,
		})
	}
	return r
}

func (l *locationConverter) hoursBackendToDomain(h *backend.Hours) *domain.Hours {
	if h == nil {
		return nil
	}
	r := &domain.Hours{
		TwentyFourSeven:     h.TwentyFourSeven,
		ExceptionalOpenings: l.exceptionPeriodsBackendToDomain(h.ExceptionalOpenings),
		ExceptionalClosings: l.exceptionPeriodsBackendToDomain(h.ExceptionalClosings),
	}
	for _, rh := range h.RegularHours {
		r.RegularHours = append(r.RegularHours, &domain.RegularHours{
			Weekday:     rh.Weekday,
			PeriodBegin: rh.PeriodBegin,
			PeriodEnd:   rh.PeriodEnd,
		})
	}
	return r
}

func (l *locationConverter) energyMixBackendToDomain(em *backend.EnergyMix) *domain.EnergyMix {
	if em == nil {
		return nil
	}
	r := &domain.EnergyMix{
		IsGreenEnergy:     em.IsGreenEnergy,
		SupplierName:      em.SupplierName,
		EnergyProductName: em.EnergyProductName,
	}
	for _, imp := range em.EnvironImpact {
		r.EnvironImpact = append(r.EnvironImpact, &domain.EnvironmentalImpact{
			Category: imp.Category,
			Amount:   imp.Amount,
		})
	}
	for _, es := range em.EnergySources {
		r.EnergySources = append(r.EnergySources, &domain.EnergySource{
			Source:     es.Source,
			Percentage: es.Percentage,
		})
	}
	return r
}

func (l *locationConverter) LocationBackendToDomain(loc *backend.Location, platformId string) *domain.Location {
	if loc == nil {
		return nil
	}
	return &domain.Location{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     loc.PartyId,
				CountryCode: loc.CountryCode,
			},
			PlatformId:  platformId,
			RefId:       loc.RefId,
			LastUpdated: loc.LastUpdated,
		},
		Id: loc.Id,
		Details: domain.LocationDetails{
			Publish:            loc.Publish,
			PublishAllowedTo:   l.publishTokenTypesBackendToDomain(loc.PublishAllowedTo),
			Name:               loc.Name,
			Address:            loc.Address,
			City:               loc.City,
			PostalCode:         loc.PostalCode,
			State:              loc.State,
			Country:            loc.Country,
			Coordinates:        *l.coordinatesBackendToDomain(&loc.Coordinates),
			RelatedLocations:   l.additionalGetLocationsBackendToDomain(loc.RelatedLocations),
			ParkingType:        loc.ParkingType,
			Directions:         l.displayTextsBackendToDomain(loc.Directions),
			Operator:           l.businessDetailsBackendToDomain(loc.Operator),
			SubOperator:        l.businessDetailsBackendToDomain(loc.SubOperator),
			Owner:              l.businessDetailsBackendToDomain(loc.Owner),
			Facilities:         loc.Facilities,
			TimeZone:           loc.TimeZone,
			OpeningTimes:       l.hoursBackendToDomain(loc.OpeningTimes),
			ChargingWhenClosed: loc.ChargingWhenClosed,
			Images:             l.imagesBackendToDomain(loc.Images),
			EnergyMix:          l.energyMixBackendToDomain(loc.EnergyMix),
		},
		Evses: l.EvsesBackendToDomain(loc.Evses, platformId, loc.Id),
	}
}
