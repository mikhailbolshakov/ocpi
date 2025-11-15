package impl

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/model"
	"github.com/mikhailbolshakov/ocpi/usecase"
)

type cdrConverter struct {
	baseConverter
	trfConverter usecase.TariffConverter
}

func NewCdrConverter(trfConverter usecase.TariffConverter) usecase.CdrConverter {
	return &cdrConverter{
		trfConverter: trfConverter,
	}
}

func (t *cdrConverter) CdrDomainToModel(cdr *domain.Cdr) *model.OcpiCdr {
	if cdr == nil {
		return nil
	}
	return &model.OcpiCdr{
		OcpiPartyId: model.OcpiPartyId{
			PartyId:     cdr.OcpiItem.ExtId.PartyId,
			CountryCode: cdr.OcpiItem.ExtId.CountryCode,
		},
		Id:            cdr.Id,
		StartDateTime: cdr.Details.StartDateTime,
		EndDateTime:   cdr.Details.EndDateTime,
		SessionId:     cdr.Details.SessionId,
		CdrToken:      t.cdrTokenDomainToModel(cdr.Details.CdrToken),
		AuthMethod:    cdr.Details.AuthMethod,
		AuthRef:       cdr.Details.AuthRef,
		CdrLocation: model.OcpiCdrLocation{
			Id:                 cdr.Id,
			Name:               cdr.Details.CdrLocation.Name,
			Address:            cdr.Details.CdrLocation.Address,
			City:               cdr.Details.CdrLocation.City,
			PostalCode:         cdr.Details.CdrLocation.PostalCode,
			State:              cdr.Details.CdrLocation.State,
			Country:            cdr.Details.CdrLocation.Country,
			Coordinates:        *t.coordinatesDomainToModel(&cdr.Details.CdrLocation.Coordinates),
			EvseUid:            cdr.Details.CdrLocation.EvseId,
			EvseId:             cdr.Details.CdrLocation.Evse,
			ConnectorId:        cdr.Details.CdrLocation.ConnectorId,
			ConnectorStandard:  cdr.Details.CdrLocation.ConnectorStandard,
			ConnectorFormat:    cdr.Details.CdrLocation.ConnectorFormat,
			ConnectorPowerType: cdr.Details.CdrLocation.ConnectorPowerType,
		},
		MeterId:                  cdr.Details.MeterId,
		Currency:                 cdr.Details.Currency,
		Tariffs:                  t.trfConverter.TariffsDomainToModel(cdr.Details.Tariffs),
		ChargingPeriods:          t.chargingPeriodDomainToModel(cdr.Details.ChargingPeriods),
		SignedData:               t.SignedDataDomainToModel(cdr.Details.SignedData),
		TotalCost:                *t.priceDomainToModel(&cdr.Details.TotalCost),
		TotalFixedCost:           t.priceDomainToModel(cdr.Details.TotalFixedCost),
		TotalEnergy:              cdr.Details.TotalEnergy,
		TotalEnergyCost:          t.priceDomainToModel(cdr.Details.TotalEnergyCost),
		TotalTime:                cdr.Details.TotalTime,
		TotalTimeCost:            t.priceDomainToModel(cdr.Details.TotalTimeCost),
		TotalParkingTime:         cdr.Details.TotalParkingTime,
		TotalParkingCost:         t.priceDomainToModel(cdr.Details.TotalParkingCost),
		TotalReservationCost:     t.priceDomainToModel(cdr.Details.TotalReservationCost),
		Remark:                   cdr.Details.Remark,
		InvoiceReferenceId:       cdr.Details.InvoiceReferenceId,
		Credit:                   cdr.Details.Credit,
		CreditReferenceId:        cdr.Details.CreditReferenceId,
		HomeChargingCompensation: cdr.Details.HomeChargingCompensation,
		LastUpdated:              cdr.LastUpdated,
	}
}

func (t *cdrConverter) CdrsDomainToModel(ts []*domain.Cdr) []*model.OcpiCdr {
	return kit.Select(ts, t.CdrDomainToModel)
}

func (t *cdrConverter) CdrModelToDomain(cdr *model.OcpiCdr, platformId string) *domain.Cdr {
	if cdr == nil {
		return nil
	}
	return &domain.Cdr{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     cdr.PartyId,
				CountryCode: cdr.CountryCode,
			},
			PlatformId:  platformId,
			LastUpdated: cdr.LastUpdated,
		},
		Id: cdr.Id,
		Details: domain.CdrDetails{
			StartDateTime: cdr.StartDateTime,
			EndDateTime:   cdr.EndDateTime,
			SessionId:     cdr.SessionId,
			CdrToken:      t.cdrTokenModelToDomain(cdr.CdrToken),
			AuthMethod:    cdr.AuthMethod,
			AuthRef:       cdr.AuthRef,
			CdrLocation: domain.CdrLocation{
				Id:                 cdr.CdrLocation.Id,
				Name:               cdr.CdrLocation.Name,
				Address:            cdr.CdrLocation.Address,
				City:               cdr.CdrLocation.City,
				PostalCode:         cdr.CdrLocation.PostalCode,
				State:              cdr.CdrLocation.State,
				Country:            cdr.CdrLocation.Country,
				Coordinates:        *t.coordinatesModelToDomain(&cdr.CdrLocation.Coordinates),
				EvseId:             cdr.CdrLocation.EvseUid,
				Evse:               cdr.CdrLocation.EvseId,
				ConnectorId:        cdr.CdrLocation.ConnectorId,
				ConnectorStandard:  cdr.CdrLocation.ConnectorStandard,
				ConnectorFormat:    cdr.CdrLocation.ConnectorFormat,
				ConnectorPowerType: cdr.CdrLocation.ConnectorPowerType,
			},
			MeterId:                  cdr.MeterId,
			Currency:                 cdr.Currency,
			Tariffs:                  t.trfConverter.TariffsModelToDomain(cdr.Tariffs, platformId),
			ChargingPeriods:          t.chargingPeriodModelToDomain(cdr.ChargingPeriods),
			SignedData:               t.SignedDataModelToDomain(cdr.SignedData),
			TotalCost:                *t.priceModelToDomain(&cdr.TotalCost),
			TotalFixedCost:           t.priceModelToDomain(cdr.TotalFixedCost),
			TotalEnergy:              cdr.TotalEnergy,
			TotalEnergyCost:          t.priceModelToDomain(cdr.TotalEnergyCost),
			TotalTime:                cdr.TotalTime,
			TotalTimeCost:            t.priceModelToDomain(cdr.TotalTimeCost),
			TotalParkingTime:         cdr.TotalParkingTime,
			TotalParkingCost:         t.priceModelToDomain(cdr.TotalParkingCost),
			TotalReservationCost:     t.priceModelToDomain(cdr.TotalReservationCost),
			Remark:                   cdr.Remark,
			InvoiceReferenceId:       cdr.InvoiceReferenceId,
			Credit:                   cdr.Credit,
			CreditReferenceId:        cdr.CreditReferenceId,
			HomeChargingCompensation: cdr.HomeChargingCompensation,
		},
	}
}

func (t *cdrConverter) CdrDomainToBackend(cdr *domain.Cdr) *backend.Cdr {
	if cdr == nil {
		return nil
	}
	return &backend.Cdr{
		Id:                       cdr.Id,
		StartDateTime:            cdr.Details.StartDateTime,
		EndDateTime:              cdr.Details.EndDateTime,
		SessionId:                cdr.Details.SessionId,
		MeterId:                  cdr.Details.MeterId,
		Currency:                 cdr.Details.Currency,
		Tariffs:                  t.trfConverter.TariffsDomainToBackend(cdr.Details.Tariffs),
		ChargingPeriods:          t.chargingPeriodDomainToBackend(cdr.Details.ChargingPeriods),
		TotalCost:                *t.priceDomainToBackend(&cdr.Details.TotalCost),
		TotalFixedCost:           t.priceDomainToBackend(cdr.Details.TotalFixedCost),
		TotalEnergy:              cdr.Details.TotalEnergy,
		TotalEnergyCost:          t.priceDomainToBackend(cdr.Details.TotalEnergyCost),
		TotalTime:                cdr.Details.TotalTime,
		TotalTimeCost:            t.priceDomainToBackend(cdr.Details.TotalTimeCost),
		TotalParkingTime:         cdr.Details.TotalParkingTime,
		TotalParkingCost:         t.priceDomainToBackend(cdr.Details.TotalParkingCost),
		TotalReservationCost:     t.priceDomainToBackend(cdr.Details.TotalReservationCost),
		Remark:                   cdr.Details.Remark,
		InvoiceReferenceId:       cdr.Details.InvoiceReferenceId,
		Credit:                   cdr.Details.Credit,
		CreditReferenceId:        cdr.Details.CreditReferenceId,
		HomeChargingCompensation: cdr.Details.HomeChargingCompensation,
		LastUpdated:              cdr.LastUpdated,
		PlatformId:               cdr.PlatformId,
		RefId:                    cdr.RefId,
		PartyId:                  cdr.OcpiItem.ExtId.PartyId,
		CountryCode:              cdr.OcpiItem.ExtId.CountryCode,
	}
}

func (t *cdrConverter) CdrsDomainToBackend(ts []*domain.Cdr) []*backend.Cdr {
	var r []*backend.Cdr
	for _, sess := range ts {
		r = append(r, t.CdrDomainToBackend(sess))
	}
	return r
}

func (t *cdrConverter) CdrBackendToDomain(cdr *backend.Cdr, sess *domain.Session, loc *domain.Location, evse *domain.Evse, con *domain.Connector) *domain.Cdr {
	if cdr == nil {
		return nil
	}
	return &domain.Cdr{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     cdr.PartyId,
				CountryCode: cdr.CountryCode,
			},
			PlatformId:  loc.PlatformId,
			RefId:       cdr.RefId,
			LastUpdated: cdr.LastUpdated,
		},
		Id: cdr.Id,
		Details: domain.CdrDetails{
			StartDateTime: cdr.StartDateTime,
			EndDateTime:   cdr.EndDateTime,
			SessionId:     cdr.SessionId,
			CdrToken:      sess.Details.CdrToken,
			AuthMethod:    sess.Details.AuthMethod,
			AuthRef:       sess.Details.AuthRef,
			CdrLocation: domain.CdrLocation{
				Id:                 loc.Id,
				Name:               loc.Details.Name,
				Address:            loc.Details.Address,
				City:               loc.Details.City,
				PostalCode:         loc.Details.PostalCode,
				State:              loc.Details.State,
				Country:            loc.Details.Country,
				Coordinates:        loc.Details.Coordinates,
				EvseId:             evse.Id,
				Evse:               evse.Details.EvseId,
				ConnectorId:        con.Id,
				ConnectorStandard:  con.Details.Standard,
				ConnectorFormat:    con.Details.Format,
				ConnectorPowerType: con.Details.PowerType,
			},
			MeterId:                  cdr.MeterId,
			Currency:                 cdr.Currency,
			Tariffs:                  t.trfConverter.TariffsBackendToDomain(cdr.Tariffs),
			ChargingPeriods:          t.chargingPeriodBackendToDomain(cdr.ChargingPeriods),
			SignedData:               nil,
			TotalCost:                *t.priceBackendToDomain(&cdr.TotalCost),
			TotalFixedCost:           t.priceBackendToDomain(cdr.TotalFixedCost),
			TotalEnergy:              cdr.TotalEnergy,
			TotalEnergyCost:          t.priceBackendToDomain(cdr.TotalEnergyCost),
			TotalTime:                cdr.TotalTime,
			TotalTimeCost:            t.priceBackendToDomain(cdr.TotalTimeCost),
			TotalParkingTime:         cdr.TotalParkingTime,
			TotalParkingCost:         t.priceBackendToDomain(cdr.TotalParkingCost),
			TotalReservationCost:     t.priceBackendToDomain(cdr.TotalReservationCost),
			Remark:                   cdr.Remark,
			InvoiceReferenceId:       cdr.InvoiceReferenceId,
			Credit:                   cdr.Credit,
			CreditReferenceId:        cdr.CreditReferenceId,
			HomeChargingCompensation: cdr.HomeChargingCompensation,
		},
	}
}

func (t *cdrConverter) SignedDataDomainToModel(sd *domain.SignedData) *model.OcpiSignedData {
	if sd == nil {
		return nil
	}
	return &model.OcpiSignedData{
		EncodingMethod:        sd.EncodingMethod,
		EncodingMethodVersion: sd.EncodingMethodVersion,
		PublicKey:             sd.PublicKey,
		SignedValues:          t.signedValueDomainToModel(sd.SignedValues),
		Url:                   sd.Url,
	}
}

func (t *cdrConverter) SignedDataModelToDomain(sd *model.OcpiSignedData) *domain.SignedData {
	if sd == nil {
		return nil
	}
	return &domain.SignedData{
		EncodingMethod:        sd.EncodingMethod,
		EncodingMethodVersion: sd.EncodingMethodVersion,
		PublicKey:             sd.PublicKey,
		SignedValues:          t.signedValueModelToDomain(sd.SignedValues),
		Url:                   sd.Url,
	}
}
