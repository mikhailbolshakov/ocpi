package impl

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/errors"
)

type locationService struct {
	base
	storage domain.LocationStorage
}

func NewLocationService(storage domain.LocationStorage) domain.LocationService {
	return &locationService{
		storage: storage,
	}
}

func (s *locationService) l() kit.CLogger {
	return ocpi.L().Cmp("loc-svc")
}

var (
	statusMap = map[string]struct{}{
		domain.EvseStatusAvailable:   {},
		domain.EvseStatusBlocked:     {},
		domain.EvseStatusCharging:    {},
		domain.EvseStatusInOperative: {},
		domain.EvseStatusOutOfOrder:  {},
		domain.EvseStatusPlanned:     {},
		domain.EvseStatusRemoved:     {},
		domain.EvseStatusReserved:    {},
		domain.EvseStatusUnknown:     {},
	}

	parkingTypeMap = map[string]struct{}{
		domain.ParkingTypeMotorway:          {},
		domain.ParkingTypeParkingLot:        {},
		domain.ParkingTypeDriveway:          {},
		domain.ParkingTypeGarage:            {},
		domain.ParkingTypeStreet:            {},
		domain.ParkingTypeUndergroundGarage: {},
	}

	imgCatMap = map[string]struct{}{
		domain.ImageCategoryEntrance: {},
		domain.ImageCategoryLocation: {},
		domain.ImageCategoryOther:    {},
		domain.ImageCategoryCharger:  {},
		domain.ImageCategoryNetwork:  {},
		domain.ImageCategoryOperator: {},
		domain.ImageCategoryOwner:    {},
	}

	facilityMap = map[string]struct{}{
		domain.FacilityHotel:          {},
		domain.FacilityRestaurant:     {},
		domain.FacilityCafe:           {},
		domain.FacilityMall:           {},
		domain.FacilitySupermarket:    {},
		domain.FacilitySport:          {},
		domain.FacilityRecreationArea: {},
		domain.FacilityNature:         {},
		domain.FacilityMuseum:         {},
		domain.FacilityBikeSharing:    {},
		domain.FacilityBusStop:        {},
		domain.FacilityTaxiStand:      {},
		domain.FacilityTramStop:       {},
		domain.FacilityMetroStation:   {},
		domain.FacilityTrainStation:   {},
		domain.FacilityAirport:        {},
		domain.FacilityParkingLot:     {},
		domain.FacilityCarpoolParking: {},
		domain.FacilityFuelStation:    {},
		domain.FacilityWifi:           {},
	}

	capabilityMap = map[string]struct{}{
		domain.CapabilityChargingProfile:    {},
		domain.CapabilityChargingPreference: {},
		domain.CapabilityChipCard:           {},
		domain.CapabilityContactlessCard:    {},
		domain.CapabilityCreditCard:         {},
		domain.CapabilityDebitCard:          {},
		domain.CapabilityPedTerminal:        {},
		domain.CapabilityRemoteStartStop:    {},
		domain.CapabilityReservable:         {},
		domain.CapabilityRfid:               {},
		domain.CapabilityStartSessionConReq: {},
		domain.CapabilityTokenGroup:         {},
		domain.CapabilityUnlock:             {},
	}

	parkingRestrictionMap = map[string]struct{}{
		domain.ParkingRestrictionEvOnly:   {},
		domain.ParkingRestrictionPlugged:  {},
		domain.ParkingRestrictionDisable:  {},
		domain.ParkingRestrictionCustomer: {},
		domain.ParkingRestrictionMoto:     {},
	}

	connectorTypeMap = map[string]struct{}{
		domain.ConnectorTypeChademo:            {},
		domain.ConnectorTypeChaoji:             {},
		domain.ConnectorTypeDomesticA:          {},
		domain.ConnectorTypeDomesticB:          {},
		domain.ConnectorTypeDomesticC:          {},
		domain.ConnectorTypeDomesticD:          {},
		domain.ConnectorTypeDomesticE:          {},
		domain.ConnectorTypeDomesticF:          {},
		domain.ConnectorTypeDomesticG:          {},
		domain.ConnectorTypeDomesticH:          {},
		domain.ConnectorTypeDomesticI:          {},
		domain.ConnectorTypeDomesticJ:          {},
		domain.ConnectorTypeDomesticK:          {},
		domain.ConnectorTypeDomesticL:          {},
		domain.ConnectorTypeDomesticM:          {},
		domain.ConnectorTypeDomesticN:          {},
		domain.ConnectorTypeDomesticO:          {},
		domain.ConnectorTypeGbtAc:              {},
		domain.ConnectorTypeGbtDc:              {},
		domain.ConnectorTypeSingle16:           {},
		domain.ConnectorTypeThree16:            {},
		domain.ConnectorTypeThree32:            {},
		domain.ConnectorTypeThree64:            {},
		domain.ConnectorTypeT1:                 {},
		domain.ConnectorTypeT1Combo:            {},
		domain.ConnectorTypeT2:                 {},
		domain.ConnectorTypeT2Combo:            {},
		domain.ConnectorTypeT3A:                {},
		domain.ConnectorTypeT3C:                {},
		domain.ConnectorTypeNema520:            {},
		domain.ConnectorTypeNema630:            {},
		domain.ConnectorTypeNema650:            {},
		domain.ConnectorTypeNema1030:           {},
		domain.ConnectorTypeNema1050:           {},
		domain.ConnectorTypeNema1430:           {},
		domain.ConnectorTypeNema1450:           {},
		domain.ConnectorTypePantographBottomUp: {},
		domain.ConnectorTypePantographTopDown:  {},
		domain.ConnectorTypeTeslaR:             {},
		domain.ConnectorTypeTeslaS:             {},
	}

	formatMap = map[string]struct{}{
		domain.FormatSocket: {},
		domain.FormatCable:  {},
	}

	powerTypeMap = map[string]struct{}{
		domain.PowerTypeAc1Phase:      {},
		domain.PowerTypeAc2Phase:      {},
		domain.PowerTypeAc2PhaseSplit: {},
		domain.PowerTypeAc3Phase:      {},
		domain.PowerTypeDc:            {},
	}
)

func (s *locationService) PutLocation(ctx context.Context, loc *domain.Location) (*domain.Location, error) {
	l := s.l().C(ctx).Mth("put-loc").F(kit.KV{"locId": loc.Id}).Dbg()

	if loc.Id == "" {
		return nil, errors.ErrLocIdEmpty(ctx)
	}

	// search by id
	stored, err := s.storage.GetLocation(ctx, loc.Id, true)
	if err != nil {
		return nil, err
	}

	// check last_updated
	if stored != nil && loc.LastUpdated.Before(stored.LastUpdated) {
		l.Warn("later changes found")
		return nil, nil
	}

	// validate
	err = s.validateAndPopulatePutLocation(ctx, loc, stored)
	if err != nil {
		return nil, err
	}

	err = s.storage.MergeLocation(ctx, loc)
	if err != nil {
		return nil, err
	}

	return loc, nil
}

func (s *locationService) GetLocation(ctx context.Context, locId string, withEvse bool) (*domain.Location, error) {
	s.l().C(ctx).Mth("get-loc").Dbg()
	if locId == "" {
		return nil, errors.ErrLocIdEmpty(ctx)
	}
	return s.storage.GetLocation(ctx, locId, withEvse)
}

func (s *locationService) MergeLocation(ctx context.Context, loc *domain.Location) (*domain.Location, error) {
	l := s.l().C(ctx).Mth("merge-loc").F(kit.KV{"locId": loc.Id}).Dbg()

	if loc.Id == "" {
		return nil, errors.ErrLocIdEmpty(ctx)
	}
	if err := s.validateLastUpdated(ctx, loc.LastUpdated); err != nil {
		return nil, err
	}

	// search by id
	stored, err := s.storage.GetLocation(ctx, loc.Id, false)
	if err != nil {
		return nil, err
	}
	if stored == nil {
		return nil, errors.ErrLocationNotFound(ctx)
	}

	// check last_updated
	if loc.LastUpdated.Before(stored.LastUpdated) {
		l.Warn("later changes found")
		return nil, nil
	}

	// validate
	err = s.validateAndPopulateMergeLocation(ctx, loc, stored)
	if err != nil {
		return nil, err
	}

	err = s.storage.MergeLocation(ctx, stored)
	if err != nil {
		return nil, err
	}

	return stored, nil
}

func (s *locationService) SearchLocations(ctx context.Context, cr *domain.LocationSearchCriteria) (*domain.LocationSearchResponse, error) {
	s.l().C(ctx).Mth("search-loc").Dbg()
	if cr.Limit == nil {
		cr.Limit = kit.IntPtr(20)
	}
	return s.storage.SearchLocations(ctx, cr)
}

func (s *locationService) DeleteLocationsByExtId(ctx context.Context, extId domain.PartyExtId) error {
	s.l().C(ctx).Mth("del-ext").F(kit.KV{"partyId": extId.PartyId, "country": extId.CountryCode}).Dbg()
	if err := s.validateExtId(ctx, extId); err != nil {
		return err
	}
	return s.storage.DeleteLocationsByExtId(ctx, extId)
}

func (s *locationService) PutEvse(ctx context.Context, evse *domain.Evse) (*domain.Evse, error) {
	l := s.l().C(ctx).Mth("put-evse").F(kit.KV{"evseId": evse.Id}).Dbg()

	if evse.Id == "" {
		return nil, errors.ErrEvseIdEmpty(ctx)
	}
	if evse.LocationId == "" {
		return nil, errors.ErrLocIdEmpty(ctx)
	}

	// search by id
	stored, err := s.storage.GetEvse(ctx, evse.LocationId, evse.Id, true)
	if err != nil {
		return nil, err
	}

	// check last_updated
	if stored != nil && evse.LastUpdated.Before(stored.LastUpdated) {
		l.Warn("later changes found")
		return nil, nil
	}

	// validate
	err = s.validateAndPopulatePutEvse(ctx, nil, evse, stored)
	if err != nil {
		return nil, err
	}

	err = s.storage.MergeEvse(ctx, evse)
	if err != nil {
		return nil, err
	}

	return evse, nil
}

func (s *locationService) GetEvse(ctx context.Context, locId, evseId string, withConnectors bool) (*domain.Evse, error) {
	s.l().C(ctx).Mth("get-evse").Dbg()
	if locId == "" {
		return nil, errors.ErrLocIdEmpty(ctx)
	}
	if evseId == "" {
		return nil, errors.ErrEvseIdEmpty(ctx)
	}
	return s.storage.GetEvse(ctx, locId, evseId, withConnectors)
}

func (s *locationService) MergeEvse(ctx context.Context, evse *domain.Evse) (*domain.Evse, error) {
	l := s.l().C(ctx).Mth("merge-evse").F(kit.KV{"evseId": evse.Id}).Dbg()

	if evse.Id == "" {
		return nil, errors.ErrEvseIdEmpty(ctx)
	}
	if evse.LocationId == "" {
		return nil, errors.ErrLocIdEmpty(ctx)
	}
	if err := s.validateLastUpdated(ctx, evse.LastUpdated); err != nil {
		return nil, err
	}

	// search by id
	stored, err := s.storage.GetEvse(ctx, evse.LocationId, evse.Id, false)
	if err != nil {
		return nil, err
	}
	if stored == nil {
		return nil, errors.ErrEvseNotFound(ctx)
	}

	// check last_updated
	if evse.LastUpdated.Before(stored.LastUpdated) {
		l.Warn("later changes found")
		return nil, nil
	}

	// validate
	err = s.validateAndPopulateMergeEvse(ctx, evse, stored)
	if err != nil {
		return nil, err
	}

	err = s.storage.MergeEvse(ctx, stored)
	if err != nil {
		return nil, err
	}

	return stored, nil
}

func (s *locationService) SearchEvses(ctx context.Context, cr *domain.EvseSearchCriteria) (*domain.EvseSearchResponse, error) {
	s.l().C(ctx).Mth("search-evse").Dbg()
	if cr.Limit == nil {
		cr.Limit = kit.IntPtr(20)
	}
	return s.storage.SearchEvses(ctx, cr)
}

func (s *locationService) PutConnector(ctx context.Context, con *domain.Connector) (*domain.Connector, error) {
	l := s.l().C(ctx).Mth("put-con").F(kit.KV{"conId": con.Id}).Dbg()

	if con.LocationId == "" {
		return nil, errors.ErrLocIdEmpty(ctx)
	}
	if con.EvseId == "" {
		return nil, errors.ErrEvseIdEmpty(ctx)
	}
	if con.Id == "" {
		return nil, errors.ErrConIdEmpty(ctx)
	}

	// search by id
	stored, err := s.storage.GetConnector(ctx, con.LocationId, con.EvseId, con.Id)
	if err != nil {
		return nil, err
	}

	// check last_updated
	if stored != nil && con.LastUpdated.Before(stored.LastUpdated) {
		l.Warn("later changes found")
		return nil, nil
	}

	// validate
	err = s.validateAndPopulatePutConnector(ctx, nil, con, stored)
	if err != nil {
		return nil, err
	}

	err = s.storage.MergeConnector(ctx, con)
	if err != nil {
		return nil, err
	}

	return con, nil
}

func (s *locationService) GetConnector(ctx context.Context, locId, evseId, conId string) (*domain.Connector, error) {
	s.l().C(ctx).Mth("get-con").Dbg()
	if locId == "" {
		return nil, errors.ErrLocIdEmpty(ctx)
	}
	if evseId == "" {
		return nil, errors.ErrEvseIdEmpty(ctx)
	}
	if conId == "" {
		return nil, errors.ErrConIdEmpty(ctx)
	}
	return s.storage.GetConnector(ctx, locId, evseId, conId)
}

func (s *locationService) MergeConnector(ctx context.Context, con *domain.Connector) (*domain.Connector, error) {
	l := s.l().C(ctx).Mth("merge-con").F(kit.KV{"conId": con.Id}).Dbg()

	if con.Id == "" {
		return nil, errors.ErrConIdEmpty(ctx)
	}
	if con.LocationId == "" {
		return nil, errors.ErrLocIdEmpty(ctx)
	}
	if con.EvseId == "" {
		return nil, errors.ErrEvseIdEmpty(ctx)
	}
	if err := s.validateLastUpdated(ctx, con.LastUpdated); err != nil {
		return nil, err
	}

	// search by id
	stored, err := s.storage.GetConnector(ctx, con.LocationId, con.EvseId, con.Id)
	if err != nil {
		return nil, err
	}
	if stored == nil {
		return nil, errors.ErrConNotFound(ctx)
	}

	// check last_updated
	if con.LastUpdated.Before(stored.LastUpdated) {
		l.Warn("later changes found")
		return nil, nil
	}

	// validate
	err = s.validateAndPopulateMergeConnector(ctx, con, stored)
	if err != nil {
		return nil, err
	}

	err = s.storage.MergeConnector(ctx, stored)
	if err != nil {
		return nil, err
	}

	return stored, nil
}

func (s *locationService) SearchConnectors(ctx context.Context, cr *domain.ConnectorSearchCriteria) (*domain.ConnectorSearchResponse, error) {
	s.l().C(ctx).Mth("search-con").Dbg()
	if cr.Limit == nil {
		cr.Limit = kit.IntPtr(20)
	}
	return s.storage.SearchConnectors(ctx, cr)
}

func (s *locationService) validateHours(ctx context.Context, entity string, hours *domain.Hours) error {
	if hours == nil {
		return nil
	}
	if !hours.TwentyFourSeven && len(hours.RegularHours) == 0 {
		return errors.ErrLocDetailsHoursEmptyAttr(ctx, entity, "regular_hours")
	}
	for _, rh := range hours.RegularHours {
		if rh.Weekday < 1 || rh.Weekday > 7 {
			return errors.ErrLocDetailsHoursInvalidAttr(ctx, entity, "weekday")
		}
		if err := s.validateTimePeriod(ctx, entity, "regular_hours.begin", rh.PeriodBegin); err != nil {
			return err
		}
		if err := s.validateTimePeriod(ctx, entity, "regular_hours.end", rh.PeriodEnd); err != nil {
			return err
		}
	}
	curYear := kit.Now().Year()
	for _, ep := range hours.ExceptionalOpenings {
		if ep.PeriodEnd.Year() < curYear-3 || ep.PeriodBegin.Year() < curYear-3 {
			return errors.ErrLocDetailsHoursInvalidAttr(ctx, entity, "exp_openings")
		}
		if ep.PeriodEnd.Before(ep.PeriodBegin) {
			return errors.ErrLocDetailsHoursInvalidAttr(ctx, entity, "exp_openings")
		}
	}
	for _, ep := range hours.ExceptionalClosings {
		if ep.PeriodEnd.Year() < curYear-3 || ep.PeriodBegin.Year() < curYear-3 {
			return errors.ErrLocDetailsHoursInvalidAttr(ctx, entity, "exp_openings")
		}
		if ep.PeriodEnd.Before(ep.PeriodBegin) {
			return errors.ErrLocDetailsHoursInvalidAttr(ctx, entity, "exp_closings")
		}
	}
	return nil
}

func (s *locationService) validateAddress(ctx context.Context, details *domain.LocationDetails) error {
	if details.Address == "" {
		return errors.ErrLocDetailsEmptyAttr(ctx, "address")
	}
	if details.City == "" {
		return errors.ErrLocDetailsEmptyAttr(ctx, "city")
	}
	if details.Country == "" {
		return errors.ErrLocDetailsEmptyAttr(ctx, "country")
	}
	if err := s.validateMaxLen(ctx, details.Address, 45, "location.address"); err != nil {
		return err
	}
	if err := s.validateMaxLen(ctx, details.City, 45, "location.city"); err != nil {
		return err
	}
	if err := s.validateMaxLen(ctx, details.PostalCode, 10, "location.postal_code"); err != nil {
		return err
	}
	if err := s.validateMaxLen(ctx, details.State, 20, "location.state"); err != nil {
		return err
	}
	if kit.GetCountryByAlfa3(details.Country) == nil {
		return errors.ErrLocDetailsInvalidAttr(ctx, "country")
	}
	return nil
}

func (s *locationService) validateCoordinates(ctx context.Context, details *domain.LocationDetails) error {
	if details.Coordinates.Latitude == "" || details.Coordinates.Longitude == "" {
		return errors.ErrLocDetailsEmptyAttr(ctx, "coordinates")
	}
	if !kit.IsCoordinateValid(details.Coordinates.Latitude) || !kit.IsCoordinateValid(details.Coordinates.Longitude) {
		return errors.ErrLocDetailsInvalidAttr(ctx, "coordinates")
	}
	for _, addGeo := range details.RelatedLocations {
		if addGeo.Latitude == "" || addGeo.Longitude == "" {
			return errors.ErrLocDetailsEmptyAttr(ctx, "related_location.coordinates")
		}
		if !kit.IsCoordinateValid(addGeo.Latitude) || !kit.IsCoordinateValid(addGeo.Longitude) {
			return errors.ErrLocDetailsInvalidAttr(ctx, "related_location.coordinates")
		}
		if err := s.validateDisplayText(ctx, "related_location", addGeo.Name); err != nil {
			return err
		}
	}
	return nil
}

func (s *locationService) validateOperators(ctx context.Context, details *domain.LocationDetails) error {
	if details.Operator != nil {
		if err := s.validateBusinessDetails(ctx, "operator", details.Operator); err != nil {
			return err
		}
	}
	if details.SubOperator != nil {
		if err := s.validateBusinessDetails(ctx, "sub-operator", details.SubOperator); err != nil {
			return err
		}
	}
	if details.Owner != nil {
		if err := s.validateBusinessDetails(ctx, "owner", details.Owner); err != nil {
			return err
		}
	}
	return nil
}
func (s *locationService) validateTZ(ctx context.Context, details *domain.LocationDetails) error {
	if details.TimeZone == "" {
		return errors.ErrLocDetailsEmptyAttr(ctx, "time_zone")
	}
	if !kit.IsTimeZoneIANA(details.TimeZone) {
		return errors.ErrLocDetailsInvalidAttr(ctx, "time_zone")
	}
	return nil
}

func (s *locationService) validateLocation(ctx context.Context, loc *domain.Location) error {

	if err := s.validateId(ctx, loc.Id, "location_id"); err != nil {
		return err
	}
	if err := s.validateMaxLen(ctx, loc.Details.Name, 255, "name"); err != nil {
		return err
	}
	if err := s.validateOcpiItem(ctx, &loc.OcpiItem); err != nil {
		return err
	}
	if err := s.validateAddress(ctx, &loc.Details); err != nil {
		return err
	}
	if err := s.validateCoordinates(ctx, &loc.Details); err != nil {
		return err
	}
	if err := s.validateOperators(ctx, &loc.Details); err != nil {
		return err
	}
	if err := s.validateTZ(ctx, &loc.Details); err != nil {
		return err
	}
	if err := s.validateHours(ctx, "opening_times", loc.Details.OpeningTimes); err != nil {
		return err
	}

	if loc.Details.ParkingType != "" {
		if _, ok := parkingTypeMap[loc.Details.ParkingType]; !ok {
			return errors.ErrLocDetailsInvalidAttr(ctx, "parking_type")
		}
	}
	for _, dir := range loc.Details.Directions {
		if err := s.validateDisplayText(ctx, "directions", dir); err != nil {
			return err
		}
	}
	for _, fac := range loc.Details.Facilities {
		if _, ok := facilityMap[fac]; !ok {
			return errors.ErrLocDetailsInvalidAttr(ctx, "facility")
		}
	}

	for _, img := range loc.Details.Images {
		if err := s.validateImage(ctx, "images", img); err != nil {
			return err
		}
	}

	//TODO: Energy Mix is currently not validated and not used
	return nil
}

func (s *locationService) validateAndPopulatePutLocation(ctx context.Context, loc, stored *domain.Location) error {

	// pre populate
	evseStored := make(map[string]*domain.Evse)
	if stored != nil {
		loc.LastSent = stored.LastSent
		if loc.PlatformId == "" {
			loc.PlatformId = stored.PlatformId
		}
		if loc.RefId == "" {
			loc.RefId = stored.RefId
		}
		if loc.ExtId.PartyId == "" || loc.ExtId.CountryCode == "" {
			loc.ExtId = stored.ExtId
		}
		for _, evse := range stored.Evses {
			evseStored[evse.Id] = evse
		}
	}

	// validate location object
	if err := s.validateLocation(ctx, loc); err != nil {
		return err
	}

	// validate evse
	for _, evse := range loc.Evses {
		if err := s.validateAndPopulatePutEvse(ctx, loc, evse, evseStored[evse.Id]); err != nil {
			return err
		}
	}

	return nil
}

func (s *locationService) validateAndPopulateMergeLocation(ctx context.Context, loc, stored *domain.Location) error {

	if len(loc.Evses) > 0 {
		return errors.ErrLocCannotMergeEvses(ctx)
	}

	stored.LastUpdated = loc.LastUpdated

	if loc.PlatformId != "" {
		stored.PlatformId = loc.PlatformId
	}
	if loc.RefId != "" {
		stored.RefId = loc.RefId
	}
	if loc.ExtId.PartyId != "" && loc.ExtId.CountryCode != "" {
		stored.ExtId = loc.ExtId
	}
	if loc.Details.Publish != nil {
		stored.Details.Publish = loc.Details.Publish
	}
	if len(loc.Details.Directions) > 0 {
		stored.Details.Directions = loc.Details.Directions
	}
	if loc.Details.Coordinates.Latitude != "" {
		stored.Details.Coordinates.Latitude = loc.Details.Coordinates.Latitude
	}
	if loc.Details.Coordinates.Longitude != "" {
		stored.Details.Coordinates.Longitude = loc.Details.Coordinates.Longitude
	}
	if loc.Details.TimeZone != "" {
		stored.Details.TimeZone = loc.Details.TimeZone
	}
	if loc.Details.ParkingType != "" {
		stored.Details.ParkingType = loc.Details.ParkingType
	}
	if loc.Details.Country != "" {
		stored.Details.Country = loc.Details.Country
	}
	if loc.Details.Name != "" {
		stored.Details.Name = loc.Details.Name
	}
	if loc.Details.City != "" {
		stored.Details.City = loc.Details.City
	}
	if loc.Details.Address != "" {
		stored.Details.Address = loc.Details.Address
	}
	if loc.Details.State != "" {
		stored.Details.State = loc.Details.State
	}
	if loc.Details.PostalCode != "" {
		stored.Details.PostalCode = loc.Details.PostalCode
	}
	if loc.Details.OpeningTimes != nil {
		stored.Details.OpeningTimes = loc.Details.OpeningTimes
	}
	if len(loc.Details.Images) > 0 {
		stored.Details.Images = loc.Details.Images
	}
	if len(loc.Details.Facilities) > 0 {
		stored.Details.Facilities = loc.Details.Facilities
	}
	if len(loc.Details.RelatedLocations) > 0 {
		stored.Details.RelatedLocations = loc.Details.RelatedLocations
	}
	if loc.Details.Owner != nil {
		stored.Details.Owner = loc.Details.Owner
	}
	if loc.Details.SubOperator != nil {
		stored.Details.SubOperator = loc.Details.SubOperator
	}
	if loc.Details.Operator != nil {
		stored.Details.Operator = loc.Details.Operator
	}
	if loc.Details.ChargingWhenClosed != nil {
		stored.Details.ChargingWhenClosed = loc.Details.ChargingWhenClosed
	}
	if loc.Details.EnergyMix != nil {
		stored.Details.EnergyMix = loc.Details.EnergyMix
	}

	return s.validateLocation(ctx, stored)
}

func (s *locationService) validateEvse(ctx context.Context, evse *domain.Evse) error {

	if err := s.validateId(ctx, evse.Id, "evse_id"); err != nil {
		return err
	}
	if err := s.validateOcpiItem(ctx, &evse.OcpiItem); err != nil {
		return err
	}
	if evse.Id == "" {
		return errors.ErrEvseIdEmpty(ctx)
	}
	if err := s.validateMaxLen(ctx, evse.Details.EvseId, 48, "evse_id"); err != nil {
		return err
	}
	if evse.LocationId == "" {
		return errors.ErrLocIdEmpty(ctx)
	}
	if _, ok := statusMap[evse.Status]; !ok {
		return errors.ErrEvseDetailsInvalidAttr(ctx, "status")
	}
	for _, ss := range evse.Details.StatusSchedule {
		if _, ok := statusMap[ss.Status]; !ok {
			return errors.ErrEvseDetailsInvalidAttr(ctx, "status_schedule.status")
		}
		if ss.PeriodEnd != nil && ss.PeriodBegin.After(*ss.PeriodEnd) {
			return errors.ErrEvseDetailsInvalidAttr(ctx, "status_schedule.period_end")
		}
	}
	for _, cp := range evse.Details.Capabilities {
		if _, ok := capabilityMap[cp]; !ok {
			return errors.ErrEvseDetailsInvalidAttr(ctx, "capabilities")
		}
	}
	if err := s.validateMaxLen(ctx, evse.Details.FloorLevel, 4, "floor_level"); err != nil {
		return err
	}
	if err := s.validateMaxLen(ctx, evse.Details.PhysicalReference, 16, "physical_reference"); err != nil {
		return err
	}
	if evse.Details.Coordinates != nil {
		if !kit.IsCoordinateValid(evse.Details.Coordinates.Latitude) || !kit.IsCoordinateValid(evse.Details.Coordinates.Longitude) {
			return errors.ErrEvseDetailsInvalidAttr(ctx, "coordinates")
		}
	}
	for _, dir := range evse.Details.Directions {
		if err := s.validateDisplayText(ctx, "directions", dir); err != nil {
			return err
		}
	}
	for _, pr := range evse.Details.ParkingRestrictions {
		if _, ok := parkingRestrictionMap[pr]; !ok {
			return errors.ErrEvseDetailsInvalidAttr(ctx, "parking_restriction")
		}
	}
	for _, img := range evse.Details.Images {
		if err := s.validateImage(ctx, "images", img); err != nil {
			return err
		}
	}
	return nil
}

func (s *locationService) validateAndPopulatePutEvse(ctx context.Context, loc *domain.Location, evse, stored *domain.Evse) error {

	// if parent location specified
	if loc != nil {
		evse.LocationId = loc.Id
		evse.ExtId = loc.ExtId
		evse.PlatformId = loc.PlatformId
		evse.LastUpdated = loc.LastUpdated
	}

	conStored := make(map[string]*domain.Connector)
	if stored != nil {
		evse.LastSent = stored.LastSent
		if evse.PlatformId == "" {
			evse.PlatformId = stored.PlatformId
		}
		if evse.RefId == "" {
			evse.RefId = stored.RefId
		}
		if evse.ExtId.PartyId == "" || evse.ExtId.CountryCode == "" {
			evse.ExtId = stored.ExtId
		}
		for _, con := range stored.Connectors {
			conStored[con.Id] = con
		}
		if evse.LocationId == "" {
			evse.LocationId = stored.LocationId
		}
	}

	// validate evse object
	if err := s.validateEvse(ctx, evse); err != nil {
		return err
	}

	// validate connectors
	for _, con := range evse.Connectors {
		if err := s.validateAndPopulatePutConnector(ctx, evse, con, conStored[con.Id]); err != nil {
			return err
		}
	}

	return nil
}

func (s *locationService) validateAndPopulateMergeEvse(ctx context.Context, evse, stored *domain.Evse) error {

	if len(evse.Connectors) > 0 {
		return errors.ErrEvseCannotMergeConnectors(ctx)
	}

	stored.LastUpdated = evse.LastUpdated

	if evse.PlatformId != "" {
		stored.PlatformId = evse.PlatformId
	}
	if evse.RefId != "" {
		stored.RefId = evse.RefId
	}
	if evse.ExtId.PartyId != "" && evse.ExtId.CountryCode != "" {
		stored.ExtId = evse.ExtId
	}
	if evse.LocationId != "" {
		stored.LocationId = evse.LocationId
	}
	if evse.Status != "" {
		stored.Status = evse.Status
	}
	if evse.Details.EvseId != "" {
		stored.Details.EvseId = evse.Details.EvseId
	}
	if len(evse.Details.StatusSchedule) > 0 {
		stored.Details.StatusSchedule = evse.Details.StatusSchedule
	}
	if len(evse.Details.Capabilities) > 0 {
		stored.Details.Capabilities = evse.Details.Capabilities
	}
	if evse.Details.Coordinates != nil {
		stored.Details.Coordinates = evse.Details.Coordinates
	}
	if evse.Details.PhysicalReference != "" {
		stored.Details.PhysicalReference = evse.Details.PhysicalReference
	}
	if len(evse.Details.Directions) > 0 {
		stored.Details.Directions = evse.Details.Directions
	}
	if len(evse.Details.ParkingRestrictions) > 0 {
		stored.Details.ParkingRestrictions = evse.Details.ParkingRestrictions
	}
	if len(evse.Details.Images) > 0 {
		stored.Details.Images = evse.Details.Images
	}
	return s.validateEvse(ctx, stored)
}

func (s *locationService) validateConnector(ctx context.Context, con *domain.Connector) error {

	if err := s.validateId(ctx, con.Id, "connector_id"); err != nil {
		return err
	}
	if err := s.validateOcpiItem(ctx, &con.OcpiItem); err != nil {
		return err
	}
	if con.Id == "" {
		return errors.ErrConIdEmpty(ctx)
	}
	if con.LocationId == "" {
		return errors.ErrLocIdEmpty(ctx)
	}
	if con.EvseId == "" {
		return errors.ErrEvseIdEmpty(ctx)
	}

	if con.Details.Standard == "" {
		return errors.ErrConDetailsEmptyAttr(ctx, "standard")
	}
	if _, ok := connectorTypeMap[con.Details.Standard]; !ok {
		return errors.ErrConDetailsInvalidAttr(ctx, "standard")
	}
	if con.Details.Format == "" {
		return errors.ErrConDetailsEmptyAttr(ctx, "format")
	}
	if _, ok := formatMap[con.Details.Format]; !ok {
		return errors.ErrConDetailsInvalidAttr(ctx, "format")
	}
	if con.Details.PowerType == "" {
		return errors.ErrConDetailsEmptyAttr(ctx, "power_type")
	}
	if _, ok := powerTypeMap[con.Details.PowerType]; !ok {
		return errors.ErrConDetailsInvalidAttr(ctx, "power_type")
	}
	if con.Details.MaxVoltage < 0 || con.Details.MaxVoltage > 10000 {
		return errors.ErrConDetailsInvalidAttr(ctx, "max_voltage")
	}
	if con.Details.MaxAmperage < 0 || con.Details.MaxAmperage > 10000 {
		return errors.ErrConDetailsInvalidAttr(ctx, "max_amperage")
	}
	if con.Details.MaxElectricPower != nil && (*con.Details.MaxElectricPower < 0 || *con.Details.MaxElectricPower > 100000000) {
		return errors.ErrConDetailsInvalidAttr(ctx, "max_electric_power")
	}
	// TODO: tariffs isn't obligatory
	//if len(con.Details.TariffIds) == 0 {
	//	return errors.ErrConDetailsEmptyAttr(ctx, "tariff_ids")
	//}
	if con.Details.TermsAndConditions != "" && !kit.IsUrlValid(con.Details.TermsAndConditions) {
		return errors.ErrConDetailsInvalidAttr(ctx, "terms_and_conditions")
	}

	return nil
}

func (s *locationService) validateAndPopulatePutConnector(ctx context.Context, evse *domain.Evse, con, stored *domain.Connector) error {

	if evse != nil {
		con.LocationId = evse.LocationId
		con.EvseId = evse.Id
		con.LastUpdated = evse.LastUpdated
		con.ExtId = evse.ExtId
		con.PlatformId = evse.PlatformId
	}

	if stored != nil {
		con.LastSent = stored.LastSent
		if con.PlatformId == "" {
			con.PlatformId = stored.PlatformId
		}
		if con.RefId == "" {
			con.RefId = stored.RefId
		}
		if con.ExtId.PartyId == "" || con.ExtId.CountryCode == "" {
			con.ExtId = stored.ExtId
		}
		if con.LocationId == "" {
			con.LocationId = stored.LocationId
		}
		if con.EvseId == "" {
			con.EvseId = stored.EvseId
		}
	}

	return s.validateConnector(ctx, con)
}

func (s *locationService) validateAndPopulateMergeConnector(ctx context.Context, con, stored *domain.Connector) error {

	stored.LastUpdated = con.LastUpdated

	if con.PlatformId != "" {
		stored.PlatformId = con.PlatformId
	}
	if con.RefId != "" {
		stored.RefId = con.RefId
	}
	if con.ExtId.PartyId != "" && con.ExtId.CountryCode != "" {
		stored.ExtId = con.ExtId
	}
	if con.LocationId != "" {
		stored.LocationId = con.LocationId
	}
	if con.EvseId != "" {
		stored.EvseId = con.EvseId
	}
	if con.Details.Standard != "" {
		stored.Details.Standard = con.Details.Standard
	}
	if con.Details.Format != "" {
		stored.Details.Format = con.Details.Format
	}
	if con.Details.PowerType != "" {
		stored.Details.PowerType = con.Details.PowerType
	}
	if con.Details.MaxVoltage != 0 {
		stored.Details.MaxVoltage = con.Details.MaxVoltage
	}
	if con.Details.MaxAmperage != 0 {
		stored.Details.MaxAmperage = con.Details.MaxAmperage
	}
	if con.Details.MaxElectricPower != nil {
		stored.Details.MaxElectricPower = con.Details.MaxElectricPower
	}
	if len(con.Details.TariffIds) > 0 {
		stored.Details.TariffIds = con.Details.TariffIds
	}
	if con.Details.TermsAndConditions != "" {
		stored.Details.TermsAndConditions = con.Details.TermsAndConditions
	}
	return s.validateConnector(ctx, stored)
}
