package impl

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/errors"
)

type sessionService struct {
	base
	storage domain.SessionStorage
}

func NewSessionService(storage domain.SessionStorage) domain.SessionService {
	return &sessionService{
		storage: storage,
	}
}

func (s *sessionService) l() kit.CLogger {
	return ocpi.L().Cmp("sess-svc")
}

var (
	sessionStatusMap = map[string]struct{}{
		domain.SessionStatusActive:      {},
		domain.SessionStatusCompleted:   {},
		domain.SessionStatusInvalid:     {},
		domain.SessionStatusPending:     {},
		domain.SessionStatusReservation: {},
	}

	dimensionMap = map[string]struct{}{
		domain.DimensionTypeCurrent:         {},
		domain.DimensionTypeEnergy:          {},
		domain.DimensionTypeEnergyExport:    {},
		domain.DimensionTypeEnergyImport:    {},
		domain.DimensionTypeMaxCurrent:      {},
		domain.DimensionTypeMinCurrent:      {},
		domain.DimensionTypeMaxPower:        {},
		domain.DimensionTypeMinPower:        {},
		domain.DimensionTypeParkingTime:     {},
		domain.DimensionTypePower:           {},
		domain.DimensionTypeReservationTime: {},
		domain.DimensionTypeStateOfCharge:   {},
		domain.DimensionTypeTime:            {},
	}
)

func (s *sessionService) PutSession(ctx context.Context, sess *domain.Session) (*domain.Session, error) {
	l := s.l().C(ctx).Mth("put-sess").F(kit.KV{"sessId": sess.Id}).Dbg()

	if sess.Id == "" {
		return nil, errors.ErrSessIdEmpty(ctx)
	}

	// search by id
	stored, err := s.storage.GetSession(ctx, sess.Id, false)
	if err != nil {
		return nil, err
	}

	// check last_updated
	if stored != nil && sess.LastUpdated.Before(stored.LastUpdated) {
		l.Warn("later changes found")
		return nil, nil
	}

	// validate
	err = s.validateAndPopulatePut(ctx, sess, stored)
	if err != nil {
		return nil, err
	}

	// update storage
	err = s.storage.MergeSession(ctx, sess)
	if err != nil {
		return nil, err
	}

	// update charging periods
	err = s.storage.UpdateChargingPeriods(ctx, sess, sess.ChargingPeriods)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (s *sessionService) MergeSession(ctx context.Context, sess *domain.Session) (*domain.Session, error) {
	l := s.l().C(ctx).Mth("merge-sess").F(kit.KV{"sessId": sess.Id}).Dbg()

	if sess.Id == "" {
		return nil, errors.ErrSessIdEmpty(ctx)
	}

	// search by id
	stored, err := s.storage.GetSession(ctx, sess.Id, false)
	if err != nil {
		return nil, err
	}
	if stored == nil {
		return nil, errors.ErrSessNotFound(ctx)
	}

	// check last_updated
	if sess.LastUpdated.Before(stored.LastUpdated) {
		l.Warn("later changes found")
		return nil, nil
	}

	// validate
	err = s.validateAndPopulateMerge(ctx, sess, stored)
	if err != nil {
		return nil, err
	}

	// update session
	err = s.storage.UpdateSession(ctx, stored)
	if err != nil {
		return nil, err
	}

	// if there are new charging periods, add them
	if len(stored.ChargingPeriods) > 0 {
		err = s.storage.CreateChargingPeriods(ctx, sess, stored.ChargingPeriods)
		if err != nil {
			return nil, err
		}
	}

	return stored, nil
}

func (s *sessionService) GetSession(ctx context.Context, sessId string) (*domain.Session, error) {
	s.l().C(ctx).Mth("get-sess").Dbg()
	if sessId == "" {
		return nil, errors.ErrSessIdEmpty(ctx)
	}
	return s.storage.GetSession(ctx, sessId, false)
}

func (s *sessionService) GetSessionWithPeriods(ctx context.Context, sessId string) (*domain.Session, error) {
	s.l().C(ctx).Mth("get-sess").Dbg()
	if sessId == "" {
		return nil, errors.ErrSessIdEmpty(ctx)
	}
	return s.storage.GetSession(ctx, sessId, true)
}

func (s *sessionService) SearchSessions(ctx context.Context, cr *domain.SessionSearchCriteria) (*domain.SessionSearchResponse, error) {
	s.l().C(ctx).Mth("search-sess").Dbg()
	if cr.Limit == nil {
		cr.Limit = kit.IntPtr(20)
	}
	// search sessions
	return s.storage.SearchSessions(ctx, cr)
}

func (s *sessionService) DeleteSessionsByExtId(ctx context.Context, extId domain.PartyExtId) error {
	s.l().C(ctx).Mth("del-ext").F(kit.KV{"partyId": extId.PartyId, "country": extId.CountryCode}).Dbg()
	if err := s.validateExtId(ctx, extId); err != nil {
		return err
	}
	return s.storage.DeleteSessionsByExtId(ctx, extId)
}

func (s *sessionService) validate(ctx context.Context, sess *domain.Session) error {

	if err := s.validateOcpiItem(ctx, &sess.OcpiItem); err != nil {
		return err
	}
	if sess.Id == "" {
		return errors.ErrSessIdEmpty(ctx)
	}
	if err := s.validateId(ctx, sess.Id, "session_id"); err != nil {
		return err
	}

	// start/end date
	if sess.Details.StartDateTime == nil {
		return errors.ErrSessEmptyAttr(ctx, "session", "start_date_time")
	}
	if sess.Details.StartDateTime.Year() < 2020 {
		return errors.ErrSessInvalidAttr(ctx, "session", "start_date_time")
	}
	if sess.Details.EndDateTime != nil {
		if sess.Details.EndDateTime.Year() < 2020 {
			return errors.ErrSessInvalidAttr(ctx, "session", "start_end_time")
		}
	}

	// kwh
	if sess.Details.Kwh == nil {
		return errors.ErrSessEmptyAttr(ctx, "session", "kwh")
	}
	if *sess.Details.Kwh < 0 {
		return errors.ErrSessInvalidAttr(ctx, "session", "kwh")
	}

	// token
	if sess.Details.CdrToken == nil {
		return errors.ErrSessEmptyAttr(ctx, "session", "cdr_token")
	}
	if err := s.validateCdrToken(ctx, "session", sess.Details.CdrToken); err != nil {
		return err
	}

	// auth method
	if err := s.validateAuth(ctx, "session", sess.Details.AuthMethod, sess.Details.AuthRef); err != nil {
		return err
	}

	// location
	if sess.Details.LocationId == "" {
		return errors.ErrSessEmptyAttr(ctx, "session", "location_id")
	}
	if err := s.validateId(ctx, sess.Details.LocationId, "location_id"); err != nil {
		return err
	}
	if sess.Details.EvseId == "" {
		return errors.ErrSessEmptyAttr(ctx, "session", "evse_id")
	}
	if err := s.validateId(ctx, sess.Details.EvseId, "evse_id"); err != nil {
		return err
	}
	if sess.Details.ConnectorId == "" {
		return errors.ErrSessEmptyAttr(ctx, "session", "connector_id")
	}
	if err := s.validateId(ctx, sess.Details.ConnectorId, "connector_id"); err != nil {
		return err
	}

	// currency
	if sess.Details.Currency == "" {
		return errors.ErrSessEmptyAttr(ctx, "session", "currency")
	}
	if !kit.CurrencyValid(sess.Details.Currency) {
		return errors.ErrSessInvalidAttr(ctx, "session", "currency")
	}

	// periods
	if err := s.validateMaxLen(ctx, sess.Details.MeterId, 255, "meter_id"); err != nil {
		return err
	}
	if err := s.validateChargingPeriods(ctx, "session", sess.ChargingPeriods); err != nil {
		return err
	}

	// total cost
	if err := s.validatePrice(ctx, "session", "total_cost", sess.Details.TotalCost); err != nil {
		return err
	}

	if sess.Details.Status == "" {
		return errors.ErrSessEmptyAttr(ctx, "session", "status")
	}
	if _, ok := sessionStatusMap[sess.Details.Status]; !ok {
		return errors.ErrSessInvalidAttr(ctx, "session", "status")
	}
	return nil
}

func (s *sessionService) validateAndPopulatePut(ctx context.Context, sess, stored *domain.Session) error {

	if stored != nil {
		sess.LastSent = stored.LastSent
		if sess.PlatformId == "" {
			sess.PlatformId = stored.PlatformId
		}
		if sess.RefId == "" {
			sess.RefId = stored.RefId
		}
		if sess.ExtId.PartyId == "" || sess.ExtId.CountryCode == "" {
			sess.ExtId = stored.ExtId
		}
	}

	return s.validate(ctx, sess)
}

func (s *sessionService) validateAndPopulateMerge(ctx context.Context, sess, stored *domain.Session) error {
	// LatUpdated must be always specified
	stored.LastUpdated = sess.LastUpdated

	if sess.PlatformId != "" {
		stored.PlatformId = sess.PlatformId
	}
	if sess.RefId != "" {
		stored.RefId = sess.RefId
	}
	if sess.ExtId.PartyId != "" && sess.ExtId.CountryCode != "" {
		stored.ExtId = sess.ExtId
	}
	if sess.Details.StartDateTime != nil {
		stored.Details.StartDateTime = sess.Details.StartDateTime
	}
	if sess.Details.EndDateTime != nil {
		stored.Details.EndDateTime = sess.Details.EndDateTime
	}
	if sess.Details.Kwh != nil {
		stored.Details.Kwh = sess.Details.Kwh
	}
	if sess.Details.CdrToken != nil {
		stored.Details.CdrToken = sess.Details.CdrToken
	}
	if sess.Details.AuthMethod != "" {
		stored.Details.AuthMethod = sess.Details.AuthMethod
	}
	if sess.Details.AuthRef != "" {
		stored.Details.AuthRef = sess.Details.AuthRef
	}
	if sess.Details.LocationId != "" {
		stored.Details.LocationId = sess.Details.LocationId
	}
	if sess.Details.EvseId != "" {
		stored.Details.EvseId = sess.Details.EvseId
	}
	if sess.Details.ConnectorId != "" {
		stored.Details.ConnectorId = sess.Details.ConnectorId
	}
	if sess.Details.MeterId != "" {
		stored.Details.MeterId = sess.Details.MeterId
	}
	if sess.Details.Currency != "" {
		stored.Details.Currency = sess.Details.Currency
	}
	// when charging periods appear in PATCH, they must be added to the existent charging periods
	if len(sess.ChargingPeriods) > 0 {
		stored.ChargingPeriods = sess.ChargingPeriods
	}
	if sess.Details.TotalCost != nil {
		stored.Details.TotalCost = sess.Details.TotalCost
	}
	if sess.Details.Status != "" {
		stored.Details.Status = sess.Details.Status
	}
	return s.validate(ctx, stored)
}
