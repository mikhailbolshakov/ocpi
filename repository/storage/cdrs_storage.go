package storage

import (
	"context"
	"github.com/jackc/pgtype"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/storages/pg"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/errors"
	"gorm.io/gorm"
	"time"
)

type cdr struct {
	pg.GormDto
	Id          string        `gorm:"column:id;primaryKey"`
	SessionId   *string       `gorm:"column:session_id"`
	Details     *pgtype.JSONB `gorm:"column:details"`
	PartyId     string        `gorm:"column:party_id"`
	CountryCode string        `gorm:"column:country_code"`
	PlatformId  string        `gorm:"column:platform_id"`
	RefId       *string       `gorm:"column:ref_id"`
	LastUpdated time.Time     `gorm:"column:last_updated"`
	LastSent    *time.Time    `gorm:"column:last_sent"`
}

type cdrRead struct {
	Cdr        cdr        `gorm:"embedded"`
	TotalCount totalCount `gorm:"embedded"`
}

type cdrStorageImpl struct {
	pg *pg.Storage
}

func (s *cdrStorageImpl) l() kit.CLogger {
	return ocpi.L().Cmp("cdr-storage")
}

func newCdrStorage(pg *pg.Storage) *cdrStorageImpl {
	return &cdrStorageImpl{
		pg: pg,
	}
}

func (s *cdrStorageImpl) GetCdr(ctx context.Context, id string) (*domain.Cdr, error) {
	s.l().C(ctx).Mth("get-cdr").F(kit.KV{"cdrId": id}).Dbg()
	if id == "" {
		return nil, nil
	}
	dto := &cdr{}
	res := s.pg.Instance.Where("id = ?", id).Limit(1).Find(&dto)
	if res.Error != nil {
		return nil, errors.ErrCdrStorageGet(ctx, res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return s.toCdrDomain(dto), nil
}

func (s *cdrStorageImpl) MergeCdr(ctx context.Context, cdr *domain.Cdr) error {
	s.l().C(ctx).Mth("merge-cdr").F(kit.KV{"cdrId": cdr.Id}).Dbg()
	if err := s.pg.Instance.Scopes(merge()).Create(s.toCdrDto(cdr)).Error; err != nil {
		return errors.ErrCdrStorageMerge(ctx, err)
	}
	return nil
}

func (s *cdrStorageImpl) UpdateCdr(ctx context.Context, cdr *domain.Cdr) error {
	s.l().C(ctx).Mth("update-cdr").F(kit.KV{"cdrId": cdr.Id}).Dbg()
	if err := s.pg.Instance.Scopes(update()).Save(s.toCdrDto(cdr)).Error; err != nil {
		return errors.ErrCdrStorageUpdate(ctx, err)
	}
	return nil
}

func (s *cdrStorageImpl) DeleteCdrsByExtId(ctx context.Context, extId domain.PartyExtId) error {
	s.l().C(ctx).Mth("delete-ext").F(kit.KV{"partyId": extId.PartyId, "country": extId.CountryCode}).Dbg()
	if extId.PartyId == "" || extId.CountryCode == "" {
		return nil
	}
	if err := s.pg.Instance.Where(`party_id = ? and country_code = ?`, extId.PartyId, extId.CountryCode).
		Delete(&cdr{}).Error; err != nil {
		return errors.ErrCdrStorageDelete(ctx, err)
	}
	return nil
}

func (s *cdrStorageImpl) SearchCdrs(ctx context.Context, cr *domain.CdrSearchCriteria) (*domain.CdrSearchResponse, error) {
	s.l().Mth("search-cdr").C(ctx).Dbg()

	rs := &domain.CdrSearchResponse{
		PageResponse: domain.PageResponse{
			Limit: pagingLimit(cr.PageRequest.Limit),
			Total: kit.IntPtr(0),
		},
	}

	// make query
	var dtosRead []*cdrRead

	if err := s.pg.Instance.
		Scopes(s.buildSearchQuery(cr), paging(cr.PageRequest), orderByLastUpdated(true)).
		Find(&dtosRead).Error; err != nil {
		return nil, errors.ErrCdrStorageGet(ctx, err)
	}

	if len(dtosRead) == 0 {
		return rs, nil
	}

	dtos := make([]*cdr, 0, len(dtosRead))
	for _, p := range dtosRead {
		dtos = append(dtos, &p.Cdr)
	}

	rs.Items = s.toCdrsDomain(dtos)
	rs.Total = &dtosRead[0].TotalCount.TotalCount
	rs.NextPage = nextPage(cr.PageRequest, rs.Total)

	return rs, nil
}

func (s *cdrStorageImpl) buildSearchQuery(criteria *domain.CdrSearchCriteria) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query := db.Table("cdrs").Select("cdrs.*, count(*) over() total_count")
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
		if len(criteria.IncPlatforms) > 0 {
			query = query.Where("platform_id in (?)", criteria.IncPlatforms)
		}
		if len(criteria.Ids) > 0 {
			query = query.Where("id in (?)", criteria.Ids)
		}
		if len(criteria.ExcPlatforms) > 0 {
			query = query.Where("platform_id not in (?)", criteria.ExcPlatforms)
		}
		if criteria.RefId != "" {
			query = query.Where("ref_id = ?", criteria.RefId)
		}
		return query
	}
}
