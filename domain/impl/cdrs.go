package impl

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/errors"
)

type cdrService struct {
	base
	storage    domain.CdrStorage
	trfService domain.TariffService
}

func NewCdrService(storage domain.CdrStorage, trfService domain.TariffService) domain.CdrService {
	return &cdrService{
		storage:    storage,
		trfService: trfService,
	}
}

func (s *cdrService) l() kit.CLogger {
	return ocpi.L().Cmp("cdr-svc")
}

func (s *cdrService) PutCdr(ctx context.Context, cdr *domain.Cdr) (*domain.Cdr, error) {
	l := s.l().C(ctx).Mth("put-cdr").F(kit.KV{"cdrId": cdr.Id}).Dbg()

	if cdr.Id == "" {
		return nil, errors.ErrCdrIdEmpty(ctx)
	}

	// search by id
	stored, err := s.storage.GetCdr(ctx, cdr.Id)
	if err != nil {
		return nil, err
	}

	// check last_updated
	if stored != nil && cdr.LastUpdated.Before(stored.LastUpdated) {
		l.Warn("later changes found")
		return nil, nil
	}

	// validate
	err = s.validateAndPopulatePut(ctx, cdr, stored)
	if err != nil {
		return nil, err
	}

	err = s.storage.MergeCdr(ctx, cdr)
	if err != nil {
		return nil, err
	}

	return cdr, nil
}

func (s *cdrService) GetCdr(ctx context.Context, cdrId string) (*domain.Cdr, error) {
	s.l().C(ctx).Mth("get-cdr").Dbg()
	if cdrId == "" {
		return nil, errors.ErrCdrIdEmpty(ctx)
	}
	return s.storage.GetCdr(ctx, cdrId)
}

func (s *cdrService) SearchCdrs(ctx context.Context, cr *domain.CdrSearchCriteria) (*domain.CdrSearchResponse, error) {
	s.l().C(ctx).Mth("search-cdr").Dbg()
	if cr.Limit == nil {
		cr.Limit = kit.IntPtr(20)
	}
	return s.storage.SearchCdrs(ctx, cr)
}

func (s *cdrService) DeleteCdrsByExtId(ctx context.Context, extId domain.PartyExtId) error {
	s.l().C(ctx).Mth("del-ext").F(kit.KV{"partyId": extId.PartyId, "country": extId.CountryCode}).Dbg()
	if err := s.validateExtId(ctx, extId); err != nil {
		return err
	}
	return s.storage.DeleteCdrsByExtId(ctx, extId)
}

func (s *cdrService) validateCdrLocation(ctx context.Context, loc *domain.CdrLocation) error {

	// name
	if err := s.validateMaxLen(ctx, loc.Name, 45, "name"); err != nil {
		return err
	}

	// address
	if loc.Address == "" {
		return errors.ErrCdrLocDetailsEmptyAttr(ctx, "address")
	}
	if loc.City == "" {
		return errors.ErrCdrLocDetailsEmptyAttr(ctx, "city")
	}
	if loc.Country == "" {
		return errors.ErrCdrLocDetailsEmptyAttr(ctx, "country")
	}
	if err := s.validateMaxLen(ctx, loc.Address, 45, "address"); err != nil {
		return err
	}
	if err := s.validateMaxLen(ctx, loc.City, 45, "city"); err != nil {
		return err
	}
	if err := s.validateMaxLen(ctx, loc.PostalCode, 10, "postal_code"); err != nil {
		return err
	}
	if err := s.validateMaxLen(ctx, loc.State, 20, "state"); err != nil {
		return err
	}
	if kit.GetCountryByAlfa3(loc.Country) == nil {
		return errors.ErrCdrLocDetailsInvalidAttr(ctx, "country")
	}

	// coordinates
	if loc.Coordinates.Latitude == "" || loc.Coordinates.Longitude == "" {
		return errors.ErrCdrLocDetailsEmptyAttr(ctx, "coordinates")
	}
	if !kit.IsCoordinateValid(loc.Coordinates.Latitude) || !kit.IsCoordinateValid(loc.Coordinates.Longitude) {
		return errors.ErrCdrLocDetailsInvalidAttr(ctx, "coordinates")
	}

	// location
	if loc.Id == "" {
		return errors.ErrCdrLocDetailsEmptyAttr(ctx, "location_id")
	}
	if err := s.validateId(ctx, loc.Id, "location_id"); err != nil {
		return err
	}
	if loc.EvseId == "" {
		return errors.ErrCdrLocDetailsEmptyAttr(ctx, "evse_id")
	}
	if err := s.validateId(ctx, loc.EvseId, "evse_id"); err != nil {
		return err
	}
	if loc.ConnectorId == "" {
		return errors.ErrCdrLocDetailsEmptyAttr(ctx, "connector_id")
	}
	if err := s.validateId(ctx, loc.ConnectorId, "connector_id"); err != nil {
		return err
	}

	if loc.ConnectorStandard == "" {
		return errors.ErrCdrLocDetailsEmptyAttr(ctx, "connector_standard")
	}
	if _, ok := connectorTypeMap[loc.ConnectorStandard]; !ok {
		return errors.ErrCdrLocDetailsInvalidAttr(ctx, "connector_standard")
	}
	if loc.ConnectorFormat == "" {
		return errors.ErrCdrLocDetailsEmptyAttr(ctx, "connector_format")
	}
	if _, ok := formatMap[loc.ConnectorFormat]; !ok {
		return errors.ErrCdrLocDetailsInvalidAttr(ctx, "connector_format")
	}
	if loc.ConnectorPowerType == "" {
		return errors.ErrCdrLocDetailsEmptyAttr(ctx, "connector_power_type")
	}
	if _, ok := powerTypeMap[loc.ConnectorPowerType]; !ok {
		return errors.ErrCdrLocDetailsInvalidAttr(ctx, "connector_power_type")
	}

	return nil
}

func (s *cdrService) validate(ctx context.Context, cdr *domain.Cdr) error {

	if err := s.validateOcpiItem(ctx, &cdr.OcpiItem); err != nil {
		return err
	}
	if cdr.Id == "" {
		return errors.ErrCdrIdEmpty(ctx)
	}
	if err := s.validateId(ctx, cdr.Id, "cdr_id"); err != nil {
		return err
	}

	// start/end date
	if cdr.Details.StartDateTime.Year() < 2020 {
		return errors.ErrCdrInvalidAttr(ctx, "cdr", "start_date_time")
	}
	if cdr.Details.EndDateTime.Year() < 2020 {
		return errors.ErrCdrInvalidAttr(ctx, "cdr", "start_end_time")
	}

	// session id
	if cdr.Details.SessionId == "" {
		return errors.ErrCdrEmptyAttr(ctx, "cdr", "session_id")
	}
	if err := s.validateId(ctx, cdr.Details.SessionId, "session_id"); err != nil {
		return err
	}

	// token
	if cdr.Details.CdrToken == nil {
		return errors.ErrCdrEmptyAttr(ctx, "cdr", "cdr_token")
	}
	if err := s.validateCdrToken(ctx, "cdr", cdr.Details.CdrToken); err != nil {
		return err
	}

	// auth method
	if err := s.validateAuth(ctx, "cdr", cdr.Details.AuthMethod, cdr.Details.AuthRef); err != nil {
		return err
	}

	// location
	if err := s.validateCdrLocation(ctx, &cdr.Details.CdrLocation); err != nil {
		return err
	}

	// periods
	if err := s.validateMaxLen(ctx, cdr.Details.MeterId, 255, "meter_id"); err != nil {
		return err
	}
	if err := s.validateChargingPeriods(ctx, "cdr", cdr.Details.ChargingPeriods); err != nil {
		return err
	}

	// currency
	if cdr.Details.Currency == "" {
		return errors.ErrCdrEmptyAttr(ctx, "cdr", "currency")
	}
	if !kit.CurrencyValid(cdr.Details.Currency) {
		return errors.ErrCdrInvalidAttr(ctx, "cdr", "currency")
	}

	// tariffs
	for _, t := range cdr.Details.Tariffs {
		t.ExtId = cdr.ExtId
		t.PlatformId = cdr.PlatformId
		t.LastUpdated = cdr.LastUpdated
		if err := s.trfService.Validate(ctx, t); err != nil {
			return err
		}
	}

	// cost
	if err := s.validatePrice(ctx, "cdr", "total_cost", &cdr.Details.TotalCost); err != nil {
		return err
	}
	if err := s.validatePrice(ctx, "cdr", "total_fixed_cost", cdr.Details.TotalFixedCost); err != nil {
		return err
	}
	if err := s.validatePrice(ctx, "cdr", "total_energy_cost", cdr.Details.TotalEnergyCost); err != nil {
		return err
	}
	if err := s.validatePrice(ctx, "cdr", "total_time_cost", cdr.Details.TotalTimeCost); err != nil {
		return err
	}
	if err := s.validatePrice(ctx, "cdr", "total_parking_cost", cdr.Details.TotalParkingCost); err != nil {
		return err
	}
	if err := s.validatePrice(ctx, "cdr", "total_reservation_cost", cdr.Details.TotalReservationCost); err != nil {
		return err
	}

	// amounts
	if err := s.validateAmount(ctx, "cdr", "total_energy", &cdr.Details.TotalEnergy); err != nil {
		return err
	}
	if err := s.validateAmount(ctx, "cdr", "total_time", &cdr.Details.TotalTime); err != nil {
		return err
	}
	if err := s.validateAmount(ctx, "cdr", "total_parking_time", cdr.Details.TotalParkingTime); err != nil {
		return err
	}

	// remark
	if err := s.validateMaxLen(ctx, cdr.Details.Remark, 255, "cdr.remark"); err != nil {
		return err
	}

	// invoice
	if err := s.validateMaxLen(ctx, cdr.Details.InvoiceReferenceId, 39, "cdr.invoice_reference_id"); err != nil {
		return err
	}

	// credit reference
	if err := s.validateMaxLen(ctx, cdr.Details.CreditReferenceId, 39, "cdr.credit_reference_id"); err != nil {
		return err
	}

	// tariffs are provided for all charging periods
	if len(cdr.Details.Tariffs) > 0 {
		// get not empty tariff Ids
		notEmptyTariffsFilter := func(p *domain.ChargingPeriod) bool { return p.TariffId != "" }
		getTariffIdFn := func(p *domain.ChargingPeriod) string { return p.TariffId }

		// build map of tariff Ids
		trfIdMap := kit.SetFromSlice(cdr.Details.Tariffs, func(t *domain.Tariff) string { return t.Id })
		tariffIdsOfChargingPeriods := kit.SetFromSlice(kit.Filter(cdr.Details.ChargingPeriods, notEmptyTariffsFilter), getTariffIdFn)

		// check tariffs
		for tId := range tariffIdsOfChargingPeriods {
			if _, ok := trfIdMap[tId]; !ok {
				return errors.ErrCdrNotFoundTariffForChargingPeriod(ctx, tId)
			}
		}

	}

	return nil
}

func (s *cdrService) validateAndPopulatePut(ctx context.Context, cdr, stored *domain.Cdr) error {

	if stored != nil {
		cdr.LastSent = stored.LastSent
		if cdr.PlatformId == "" {
			cdr.PlatformId = stored.PlatformId
		}
		if cdr.RefId == "" {
			cdr.RefId = stored.RefId
		}
		if cdr.ExtId.PartyId == "" || cdr.ExtId.CountryCode == "" {
			cdr.ExtId = stored.ExtId
		}
	}

	return s.validate(ctx, cdr)
}
