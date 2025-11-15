package storage

import (
	"context"
	"github.com/jackc/pgtype"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/goroutine"
	"github.com/mikhailbolshakov/kit/storages/pg"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/errors"
	"gorm.io/gorm"
	"time"
)

type session struct {
	pg.GormDto
	Id          string        `gorm:"column:id;primaryKey"`
	Details     *pgtype.JSONB `gorm:"column:details"`
	PartyId     string        `gorm:"column:party_id"`
	CountryCode string        `gorm:"column:country_code"`
	PlatformId  string        `gorm:"column:platform_id"`
	RefId       *string       `gorm:"column:ref_id"`
	LastUpdated time.Time     `gorm:"column:last_updated"`
	LastSent    *time.Time    `gorm:"column:last_sent"`
}

type sessionChargingPeriod struct {
	SessionId   string        `gorm:"column:session_id"`
	Details     *pgtype.JSONB `gorm:"column:details"`
	LastUpdated time.Time     `gorm:"column:last_updated"`
}

type sessionRead struct {
	Session    session    `gorm:"embedded"`
	TotalCount totalCount `gorm:"embedded"`
}

type sessionStorageImpl struct {
	pg *pg.Storage
}

func (s *sessionStorageImpl) l() kit.CLogger {
	return ocpi.L().Cmp("sess-storage")
}

func newSessionStorage(pg *pg.Storage) *sessionStorageImpl {
	return &sessionStorageImpl{
		pg: pg,
	}
}

func (s *sessionStorageImpl) GetSession(ctx context.Context, id string, withChargingPeriods bool) (*domain.Session, error) {
	s.l().C(ctx).Mth("get-sess").F(kit.KV{"sessId": id}).Dbg()

	if id == "" {
		return nil, nil
	}

	dto := &session{}
	res := s.pg.Instance.Where("id = ?", id).Limit(1).Find(&dto)
	if res.Error != nil {
		return nil, errors.ErrSessStorageGet(ctx, res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	sess := s.toSessionDomain(dto)

	// retrieve charging periods if requested
	if withChargingPeriods {
		periods, err := s.GetChargingPeriods(ctx, sess.Id)
		if err != nil {
			return nil, err
		}
		sess.ChargingPeriods = periods[sess.Id]
	}

	return sess, nil
}

func (s *sessionStorageImpl) MergeSession(ctx context.Context, sess *domain.Session) error {
	s.l().C(ctx).Mth("merge-sess").F(kit.KV{"sessId": sess.Id}).Dbg()
	if err := s.pg.Instance.Scopes(merge()).Create(s.toSessionDto(sess)).Error; err != nil {
		return errors.ErrSessStorageMerge(ctx, err)
	}
	return nil
}

func (s *sessionStorageImpl) UpdateSession(ctx context.Context, sess *domain.Session) error {
	s.l().C(ctx).Mth("update-sess").F(kit.KV{"sessId": sess.Id}).Dbg()
	if err := s.pg.Instance.Scopes(update()).Save(s.toSessionDto(sess)).Error; err != nil {
		return errors.ErrSessStorageUpdate(ctx, err)
	}
	return nil
}

func (s *sessionStorageImpl) DeleteSessionsByExtId(ctx context.Context, extId domain.PartyExtId) error {
	s.l().C(ctx).Mth("delete-ext").F(kit.KV{"partyId": extId.PartyId, "country": extId.CountryCode}).Dbg()
	if extId.PartyId == "" || extId.CountryCode == "" {
		return nil
	}
	if err := s.pg.Instance.Where(`party_id = ? and country_code = ?`, extId.PartyId, extId.CountryCode).
		Delete(&session{}).Error; err != nil {
		return errors.ErrSessStorageDelete(ctx, err)
	}
	return nil
}

func (s *sessionStorageImpl) CreateChargingPeriods(ctx context.Context, sess *domain.Session, periods []*domain.ChargingPeriod) error {
	l := s.l().C(ctx).Mth("create-charging-periods").F(kit.KV{"sessId": sess.Id}).Dbg()

	tx := s.pg.Instance.Begin()
	if tx.Error != nil {
		return errors.ErrSessStorageCreateTx(ctx, tx.Error)
	}

	eg := goroutine.NewGroup(ctx).WithLogger(l)

	eg.Go(func() error {
		return tx.Model(&session{Id: sess.Id}).Update("last_updated", sess.LastUpdated).Error
	})
	eg.Go(func() error {
		return tx.Create(s.toSessionChargingPeriodsDto(sess, periods, sess.LastUpdated)).Error
	})

	if err := eg.Wait(); err != nil {
		tx.Rollback()
		return errors.ErrSessStorageChargingPeriodsCreate(ctx, err)
	}

	tx.Commit()

	return nil
}

func (s *sessionStorageImpl) UpdateChargingPeriods(ctx context.Context, sess *domain.Session, periods []*domain.ChargingPeriod) error {
	l := s.l().C(ctx).Mth("update-charging-periods").F(kit.KV{"sessId": sess.Id}).Dbg()

	tx := s.pg.Instance.Begin()
	if tx.Error != nil {
		return errors.ErrSessStorageCreateTx(ctx, tx.Error)
	}

	eg := goroutine.NewGroup(ctx).WithLogger(l)

	eg.Go(func() error {
		return tx.Model(&session{Id: sess.Id}).Update("last_updated", sess.LastUpdated).Error
	})

	eg.Go(func() error {

		// delete existent periods
		err := tx.Where("session_id = ?", sess.Id).Delete(&sessionChargingPeriod{}).Error
		if err != nil {
			return err
		}

		// insert new ones if provided
		if len(periods) == 0 {
			return nil
		}
		return tx.Create(s.toSessionChargingPeriodsDto(sess, periods, sess.LastUpdated)).Error
	})

	if err := eg.Wait(); err != nil {
		tx.Rollback()
		return errors.ErrSessStorageChargingPeriodsUpdate(ctx, err)
	}

	tx.Commit()

	return nil
}

func (s *sessionStorageImpl) GetChargingPeriods(ctx context.Context, sessIds ...string) (map[string][]*domain.ChargingPeriod, error) {
	s.l().C(ctx).Mth("get-charging-periods").Dbg()

	if len(sessIds) == 0 {
		return map[string][]*domain.ChargingPeriod{}, nil
	}

	var dtos []*sessionChargingPeriod
	if err := s.pg.Instance.Where("session_id in (?)", sessIds).Find(&dtos).Error; err != nil {
		return nil, errors.ErrSessStorageChargingPeriodsGet(ctx, err)
	}

	return s.toSessionChargingPeriodsDomain(dtos), nil
}

func (s *sessionStorageImpl) SearchSessions(ctx context.Context, cr *domain.SessionSearchCriteria) (*domain.SessionSearchResponse, error) {
	s.l().Mth("search-sess").C(ctx).Dbg()

	rs := &domain.SessionSearchResponse{
		PageResponse: domain.PageResponse{
			Limit: pagingLimit(cr.PageRequest.Limit),
			Total: kit.IntPtr(0),
		},
	}

	// make query
	var dtosRead []*sessionRead

	if err := s.pg.Instance.
		Scopes(s.buildSearchQuery(cr), paging(cr.PageRequest), orderByLastUpdated(true)).
		Find(&dtosRead).Error; err != nil {
		return nil, errors.ErrSessStorageGet(ctx, err)
	}

	if len(dtosRead) == 0 {
		return rs, nil
	}

	dtos := make([]*session, 0, len(dtosRead))
	for _, p := range dtosRead {
		dtos = append(dtos, &p.Session)
	}

	rs.Items = s.toSessionsDomain(dtos)
	rs.Total = &dtosRead[0].TotalCount.TotalCount
	rs.NextPage = nextPage(cr.PageRequest, rs.Total)

	// retrieve charging periods if requested
	if cr.WithChargingPeriods && len(rs.Items) > 0 {

		// retrieve map for all sessions
		periodsMap, err := s.GetChargingPeriods(ctx, kit.SelectSlice(rs.Items, func(i *domain.Session) string { return i.Id })...)
		if err != nil {
			return nil, err
		}

		// populate response
		for _, sess := range rs.Items {
			sess.ChargingPeriods = periodsMap[sess.Id]
		}
	}

	return rs, nil
}

func (s *sessionStorageImpl) buildSearchQuery(criteria *domain.SessionSearchCriteria) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query := db.Table("sessions").Select("sessions.*, count(*) over() total_count")
		// populate conditions
		if criteria.ExtId != nil {
			query = query.Where("party_id = ? and country_code = ?", criteria.ExtId.PartyId, criteria.ExtId.CountryCode)
		}
		if criteria.DateTo != nil {
			query = query.Where("last_updated <= ?", *criteria.DateTo)
		}
		if criteria.DateFrom != nil {
			query = query.Where("last_updated >= ?", *criteria.DateFrom)
		}
		if len(criteria.Ids) > 0 {
			query = query.Where("id in (?)", criteria.Ids)
		}
		if len(criteria.IncPlatforms) > 0 {
			query = query.Where("platform_id in (?)", criteria.IncPlatforms)
		}
		if len(criteria.ExcPlatforms) > 0 {
			query = query.Where("platform_id not in (?)", criteria.ExcPlatforms)
		}
		if criteria.RefId != "" {
			query = query.Where("ref_id = ?", criteria.RefId)
		}
		if criteria.AuthRef != "" {
			query = query.Where("auth_ref = ?", criteria.AuthRef)
		}
		return query
	}
}
