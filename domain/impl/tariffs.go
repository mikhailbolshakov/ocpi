package impl

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/errors"
)

type tariffService struct {
	base
	storage domain.TariffStorage
}

func NewTariffService(storage domain.TariffStorage) domain.TariffService {
	return &tariffService{
		storage: storage,
	}
}

func (s *tariffService) l() kit.CLogger {
	return ocpi.L().Cmp("trf-svc")
}

var (
	tariffTypeMap = map[string]struct{}{
		domain.TariffTypeProfileCheap: {},
		domain.TariffTypeProfileGreen: {},
		domain.TariffTypeProfileFast:  {},
		domain.TariffTypeAdHocPay:     {},
		domain.TariffTypeReg:          {},
	}

	tariffDimType = map[string]struct{}{
		domain.TariffDimParkingType: {},
		domain.TariffDimFlat:        {},
		domain.TariffDimEnergy:      {},
		domain.TariffDimTime:        {},
	}

	daysMap = map[string]struct{}{
		domain.DayMon: {},
		domain.DayTue: {},
		domain.DayWed: {},
		domain.DayThu: {},
		domain.DayFri: {},
		domain.DaySat: {},
		domain.DaySun: {},
	}
)

func (s *tariffService) PutTariff(ctx context.Context, trf *domain.Tariff) (*domain.Tariff, error) {
	l := s.l().C(ctx).Mth("put-trf").F(kit.KV{"trfId": trf.Id}).Dbg()

	if trf.Id == "" {
		return nil, errors.ErrTrfIdEmpty(ctx)
	}

	// search by id
	stored, err := s.storage.GetTariff(ctx, trf.Id)
	if err != nil {
		return nil, err
	}

	// check last_updated
	if stored != nil && trf.LastUpdated.Before(stored.LastUpdated) {
		l.Warn("later changes found")
		return nil, nil
	}

	// validate
	err = s.validateAndPopulatePut(ctx, trf, stored)
	if err != nil {
		return nil, err
	}

	err = s.storage.MergeTariff(ctx, trf)
	if err != nil {
		return nil, err
	}

	return trf, nil
}

func (s *tariffService) MergeTariff(ctx context.Context, trf *domain.Tariff) (*domain.Tariff, error) {
	l := s.l().C(ctx).Mth("merge-trf").F(kit.KV{"trfId": trf.Id}).Dbg()

	if trf.Id == "" {
		return nil, errors.ErrTrfIdEmpty(ctx)
	}
	if err := s.validateLastUpdated(ctx, trf.LastUpdated); err != nil {
		return nil, err
	}

	// search by id
	stored, err := s.storage.GetTariff(ctx, trf.Id)
	if err != nil {
		return nil, err
	}
	if stored == nil {
		return nil, errors.ErrTrfNotFound(ctx)
	}

	// check last_updated
	if trf.LastUpdated.Before(stored.LastUpdated) {
		l.Warn("later changes found")
		return nil, nil
	}

	// validate
	err = s.validateAndPopulateMerge(ctx, trf, stored)
	if err != nil {
		return nil, err
	}

	err = s.storage.UpdateTariff(ctx, stored)
	if err != nil {
		return nil, err
	}

	return stored, nil
}

func (s *tariffService) GetTariff(ctx context.Context, trfId string) (*domain.Tariff, error) {
	s.l().C(ctx).Mth("get-trf").Dbg()
	if trfId == "" {
		return nil, errors.ErrTrfIdEmpty(ctx)
	}
	return s.storage.GetTariff(ctx, trfId)
}

func (s *tariffService) SearchTariffs(ctx context.Context, cr *domain.TariffSearchCriteria) (*domain.TariffSearchResponse, error) {
	s.l().C(ctx).Mth("search-trf").Dbg()
	if cr.Limit == nil {
		cr.Limit = kit.IntPtr(20)
	}
	return s.storage.SearchTariffs(ctx, cr)
}

func (s *tariffService) DeleteTariffsByExtId(ctx context.Context, extId domain.PartyExtId) error {
	s.l().C(ctx).Mth("del-ext").F(kit.KV{"partyId": extId.PartyId, "country": extId.CountryCode}).Dbg()
	if err := s.validateExtId(ctx, extId); err != nil {
		return err
	}
	return s.storage.DeleteTariffsByExtId(ctx, extId)
}

func (s *tariffService) Validate(ctx context.Context, trf *domain.Tariff) error {

	if err := s.validateOcpiItem(ctx, &trf.OcpiItem); err != nil {
		return err
	}
	if trf.Id == "" {
		return errors.ErrTrfIdEmpty(ctx)
	}
	if err := s.validateId(ctx, trf.Id, "id"); err != nil {
		return err
	}

	// currency
	if trf.Details.Currency == "" {
		return errors.ErrTrfEmptyAttr(ctx, "tariff", "currency")
	}
	if !kit.CurrencyValid(trf.Details.Currency) {
		return errors.ErrTrfInvalidAttr(ctx, "tariff", "currency")
	}

	// type
	if trf.Details.Type != "" {
		if _, ok := tariffTypeMap[trf.Details.Type]; !ok {
			return errors.ErrTrfInvalidAttr(ctx, "tariff", "type")
		}
	}

	// info
	for _, altTxt := range trf.Details.TariffAltText {
		if err := s.validateDisplayText(ctx, "alt_text", altTxt); err != nil {
			return err
		}
	}
	if trf.Details.TariffAltUrl != "" && !kit.IsUrlValid(trf.Details.TariffAltUrl) {
		return errors.ErrTrfInvalidAttr(ctx, "tariff", "tariff_alt_url")
	}

	// min/max price
	if err := s.validatePrice(ctx, "tariff", "min_price", trf.Details.MinPrice); err != nil {
		return err
	}
	if err := s.validatePrice(ctx, "tariff", "max_price", trf.Details.MaxPrice); err != nil {
		return err
	}

	// elements
	if err := s.validateElements(ctx, &trf.Details); err != nil {
		return err
	}

	// period
	if trf.Details.StartDateTime != nil && trf.Details.EndDateTime != nil && trf.Details.StartDateTime.After(*trf.Details.EndDateTime) {
		return errors.ErrTrfInvalidAttr(ctx, "element", "start_date_time")
	}

	//TODO: Energy Mix is currently not validated and not used
	return nil
}

func (s *tariffService) validatePriceComponent(ctx context.Context, pc *domain.PriceComponent) error {
	if pc.Type == "" {
		return errors.ErrTrfEmptyAttr(ctx, "price_component", "type")
	}
	if _, ok := tariffDimType[pc.Type]; !ok {
		return errors.ErrTrfInvalidAttr(ctx, "price_component", "type")
	}
	if pc.Price < 0 {
		return errors.ErrTrfInvalidAttr(ctx, "price_component", "price")
	}
	if pc.Vat != nil && *pc.Vat < 0 {
		return errors.ErrTrfInvalidAttr(ctx, "price_component", "vat")
	}
	if pc.StepSize < 0 {
		return errors.ErrTrfInvalidAttr(ctx, "price_component", "step_size")
	}
	return nil
}

func (s *tariffService) validateRestriction(ctx context.Context, r *domain.TariffRestrictions) error {
	if r == nil {
		return nil
	}
	if err := s.validateTimePeriod(ctx, "restriction", "start_time", r.StartTime); err != nil {
		return err
	}
	if err := s.validateTimePeriod(ctx, "restriction", "end_time", r.EndTime); err != nil {
		return err
	}
	if r.MinKwh != nil && *r.MinKwh < 0 {
		return errors.ErrTrfInvalidAttr(ctx, "restriction", "min_kwh")
	}
	if r.MaxKwh != nil && *r.MaxKwh < 0 {
		return errors.ErrTrfInvalidAttr(ctx, "restriction", "max_kwh")
	}
	if r.MinCurrent != nil && *r.MinCurrent < 0 {
		return errors.ErrTrfInvalidAttr(ctx, "restriction", "min_current")
	}
	if r.MaxCurrent != nil && *r.MaxCurrent < 0 {
		return errors.ErrTrfInvalidAttr(ctx, "restriction", "max_current")
	}
	if r.MinPower != nil && *r.MinPower < 0 {
		return errors.ErrTrfInvalidAttr(ctx, "restriction", "min_power")
	}
	if r.MaxPower != nil && *r.MaxPower < 0 {
		return errors.ErrTrfInvalidAttr(ctx, "restriction", "max_power")
	}
	if r.MinDuration != nil && *r.MinDuration < 0 {
		return errors.ErrTrfInvalidAttr(ctx, "restriction", "min_duration")
	}
	if r.MaxDuration != nil && *r.MaxDuration < 0 {
		return errors.ErrTrfInvalidAttr(ctx, "restriction", "max_duration")
	}
	for _, d := range r.DayOfWeek {
		if _, ok := daysMap[d]; !ok {
			return errors.ErrTrfInvalidAttr(ctx, "restriction", "day_of_week")
		}
	}
	if r.Reservation != "" && r.Reservation != domain.Reservation && r.Reservation != domain.ReservationExp {
		return errors.ErrTrfInvalidAttr(ctx, "restriction", "reservation")
	}
	return nil
}

func (s *tariffService) validateElements(ctx context.Context, details *domain.TariffDetails) error {

	if len(details.Elements) == 0 {
		return errors.ErrTrfEmptyAttr(ctx, "tariff", "elements")
	}

	for _, elem := range details.Elements {

		// validate price components
		if len(elem.PriceComponents) == 0 {
			return errors.ErrTrfEmptyAttr(ctx, "element", "price_components")
		}
		for _, pc := range elem.PriceComponents {
			if err := s.validatePriceComponent(ctx, pc); err != nil {
				return err
			}
		}

		// validate restrictions
		if err := s.validateRestriction(ctx, elem.Restrictions); err != nil {
			return err
		}
	}
	return nil
}

func (s *tariffService) validateAndPopulatePut(ctx context.Context, trf, stored *domain.Tariff) error {

	if stored != nil {
		trf.LastSent = stored.LastSent
		if trf.PlatformId == "" {
			trf.PlatformId = stored.PlatformId
		}
		if trf.RefId == "" {
			trf.RefId = stored.RefId
		}
		if trf.ExtId.PartyId == "" || trf.ExtId.CountryCode == "" {
			trf.ExtId = stored.ExtId
		}
	}

	return s.Validate(ctx, trf)
}

func (s *tariffService) validateAndPopulateMerge(ctx context.Context, trf, stored *domain.Tariff) error {

	stored.LastUpdated = trf.LastUpdated

	if trf.PlatformId != "" {
		stored.PlatformId = trf.PlatformId
	}
	if trf.RefId != "" {
		stored.RefId = trf.RefId
	}
	if trf.ExtId.PartyId != "" && trf.ExtId.CountryCode != "" {
		stored.ExtId = trf.ExtId
	}
	if trf.Details.Currency != "" {
		stored.Details.Currency = trf.Details.Currency
	}
	if trf.Details.Type != "" {
		stored.Details.Type = trf.Details.Type
	}
	if len(trf.Details.TariffAltText) > 0 {
		stored.Details.TariffAltText = trf.Details.TariffAltText
	}
	if trf.Details.TariffAltUrl != "" {
		stored.Details.TariffAltUrl = trf.Details.TariffAltUrl
	}
	if trf.Details.MinPrice != nil {
		stored.Details.MinPrice = trf.Details.MinPrice
	}
	if trf.Details.MaxPrice != nil {
		stored.Details.MaxPrice = trf.Details.MaxPrice
	}
	if len(trf.Details.Elements) > 0 {
		stored.Details.Elements = trf.Details.Elements
	}
	if trf.Details.StartDateTime != nil {
		stored.Details.StartDateTime = trf.Details.StartDateTime
	}
	if trf.Details.EndDateTime != nil {
		stored.Details.EndDateTime = trf.Details.EndDateTime
	}
	if trf.Details.EnergyMix != nil {
		stored.Details.EnergyMix = trf.Details.EnergyMix
	}
	return s.Validate(ctx, stored)
}
