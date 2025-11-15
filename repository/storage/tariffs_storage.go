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

type tariff struct {
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

type tariffRead struct {
	Tariff     tariff     `gorm:"embedded"`
	TotalCount totalCount `gorm:"embedded"`
}

type tariffStorageImpl struct {
	pg *pg.Storage
}

func (s *tariffStorageImpl) l() kit.CLogger {
	return ocpi.L().Cmp("trf-storage")
}

func newTariffStorage(pg *pg.Storage) *tariffStorageImpl {
	return &tariffStorageImpl{
		pg: pg,
	}
}

func (s *tariffStorageImpl) GetTariff(ctx context.Context, id string) (*domain.Tariff, error) {
	s.l().C(ctx).Mth("get-trf").F(kit.KV{"trfId": id}).Dbg()
	if id == "" {
		return nil, nil
	}
	dto := &tariff{}
	res := s.pg.Instance.Where("id = ?", id).Limit(1).Find(&dto)
	if res.Error != nil {
		return nil, errors.ErrTrfStorageGet(ctx, res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return s.toTariffDomain(dto), nil
}

func (s *tariffStorageImpl) MergeTariff(ctx context.Context, trf *domain.Tariff) error {
	s.l().C(ctx).Mth("merge-trf").F(kit.KV{"trfId": trf.Id}).Dbg()
	if err := s.pg.Instance.Scopes(merge()).Create(s.toTariffDto(trf)).Error; err != nil {
		return errors.ErrTrfStorageMerge(ctx, err)
	}
	return nil
}

func (s *tariffStorageImpl) UpdateTariff(ctx context.Context, trf *domain.Tariff) error {
	s.l().C(ctx).Mth("update-trf").F(kit.KV{"trfId": trf.Id}).Dbg()
	if err := s.pg.Instance.Scopes(update()).Save(s.toTariffDto(trf)).Error; err != nil {
		return errors.ErrTrfStorageUpdate(ctx, err)
	}
	return nil
}

func (s *tariffStorageImpl) DeleteTariffsByExtId(ctx context.Context, extId domain.PartyExtId) error {
	s.l().C(ctx).Mth("delete-ext").F(kit.KV{"partyId": extId.PartyId, "country": extId.CountryCode}).Dbg()
	if extId.PartyId == "" || extId.CountryCode == "" {
		return nil
	}
	if err := s.pg.Instance.Where(`party_id = ? and country_code = ?`, extId.PartyId, extId.CountryCode).
		Delete(&tariff{}).Error; err != nil {
		return errors.ErrTrfStorageDelete(ctx, err)
	}
	return nil
}

func (s *tariffStorageImpl) SearchTariffs(ctx context.Context, cr *domain.TariffSearchCriteria) (*domain.TariffSearchResponse, error) {
	s.l().Mth("search-trf").C(ctx).Dbg()

	rs := &domain.TariffSearchResponse{
		PageResponse: domain.PageResponse{
			Limit: pagingLimit(cr.PageRequest.Limit),
			Total: kit.IntPtr(0),
		},
	}

	// make query
	var dtosRead []*tariffRead

	if err := s.pg.Instance.
		Scopes(s.buildSearchQuery(cr), paging(cr.PageRequest), orderByLastUpdated(true)).
		Find(&dtosRead).Error; err != nil {
		return nil, errors.ErrTrfStorageGet(ctx, err)
	}

	if len(dtosRead) == 0 {
		return rs, nil
	}

	dtos := make([]*tariff, 0, len(dtosRead))
	for _, p := range dtosRead {
		dtos = append(dtos, &p.Tariff)
	}

	rs.Items = s.toTariffsDomain(dtos)
	rs.Total = &dtosRead[0].TotalCount.TotalCount
	rs.NextPage = nextPage(cr.PageRequest, rs.Total)

	return rs, nil
}

func (s *tariffStorageImpl) buildSearchQuery(criteria *domain.TariffSearchCriteria) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query := db.Table("tariffs").Select("tariffs.*, count(*) over() total_count")
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
