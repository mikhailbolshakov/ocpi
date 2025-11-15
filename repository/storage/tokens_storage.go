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

type token struct {
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

type tokenRead struct {
	Token      token      `gorm:"embedded"`
	TotalCount totalCount `gorm:"embedded"`
}

type tokenStorageImpl struct {
	pg *pg.Storage
}

func (s *tokenStorageImpl) l() kit.CLogger {
	return ocpi.L().Cmp("tkn-storage")
}

func newTokenStorage(pg *pg.Storage) *tokenStorageImpl {
	return &tokenStorageImpl{
		pg: pg,
	}
}

func (s *tokenStorageImpl) GetToken(ctx context.Context, id string) (*domain.Token, error) {
	s.l().C(ctx).Mth("get-tkn").F(kit.KV{"tknId": id}).Dbg()
	if id == "" {
		return nil, nil
	}
	dto := &token{}
	res := s.pg.Instance.Where("id = ?", id).Limit(1).Find(&dto)
	if res.Error != nil {
		return nil, errors.ErrTknStorageGet(ctx, res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return s.toTokenDomain(dto), nil
}

func (s *tokenStorageImpl) MergeToken(ctx context.Context, tkn *domain.Token) error {
	s.l().C(ctx).Mth("merge-tkn").F(kit.KV{"tknId": tkn.Id}).Dbg()
	if err := s.pg.Instance.Scopes(merge()).Create(s.toTokenDto(tkn)).Error; err != nil {
		return errors.ErrTknStorageMerge(ctx, err)
	}
	return nil
}

func (s *tokenStorageImpl) UpdateToken(ctx context.Context, tkn *domain.Token) error {
	s.l().C(ctx).Mth("update-tkn").F(kit.KV{"tknId": tkn.Id}).Dbg()
	if err := s.pg.Instance.Scopes(update()).Save(s.toTokenDto(tkn)).Error; err != nil {
		return errors.ErrTknStorageUpdate(ctx, err)
	}
	return nil
}

func (s *tokenStorageImpl) SearchTokens(ctx context.Context, cr *domain.TokenSearchCriteria) (*domain.TokenSearchResponse, error) {
	s.l().Mth("search-tkn").C(ctx).Dbg()

	rs := &domain.TokenSearchResponse{
		PageResponse: domain.PageResponse{
			Limit: pagingLimit(cr.PageRequest.Limit),
			Total: kit.IntPtr(0),
		},
	}

	// make query
	var dtosRead []*tokenRead

	if err := s.pg.Instance.
		Scopes(s.buildSearchQuery(cr), paging(cr.PageRequest), orderByLastUpdated(true)).
		Find(&dtosRead).Error; err != nil {
		return nil, errors.ErrTknStorageGet(ctx, err)
	}

	if len(dtosRead) == 0 {
		return rs, nil
	}

	dtos := make([]*token, 0, len(dtosRead))
	for _, p := range dtosRead {
		dtos = append(dtos, &p.Token)
	}

	rs.Items = s.toTokensDomain(dtos)
	rs.Total = &dtosRead[0].TotalCount.TotalCount
	rs.NextPage = nextPage(cr.PageRequest, rs.Total)

	return rs, nil
}

func (s *tokenStorageImpl) DeleteTokensByExtId(ctx context.Context, extId domain.PartyExtId) error {
	s.l().C(ctx).Mth("delete-ext").F(kit.KV{"partyId": extId.PartyId, "country": extId.CountryCode}).Dbg()
	if extId.PartyId == "" || extId.CountryCode == "" {
		return nil
	}
	if err := s.pg.Instance.Where(`party_id = ? and country_code = ?`, extId.PartyId, extId.CountryCode).
		Delete(&token{}).Error; err != nil {
		return errors.ErrSessStorageDelete(ctx, err)
	}
	return nil
}

func (s *tokenStorageImpl) buildSearchQuery(criteria *domain.TokenSearchCriteria) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query := db.Table("tokens").Select("tokens.*, count(*) over() total_count")
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
		return query
	}
}
