package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgtype"
	"github.com/lib/pq"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/storages/pg"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/errors"
	"gorm.io/gorm"
	"strings"
	"time"
)

type partyDetails struct {
	BusinessDetails *domain.BusinessDetails `json:"businessDetails,omitempty"`
}

type party struct {
	pg.GormDto
	Id          string         `gorm:"column:id"`
	PlatformId  string         `gorm:"column:platform_id"`
	PartyId     string         `gorm:"column:party_id"`
	CountryCode string         `gorm:"column:country_code"`
	Status      string         `gorm:"column:status"`
	RefId       *string        `gorm:"column:ref_id"`
	Details     *pgtype.JSONB  `gorm:"column:details"`
	Roles       pq.StringArray `gorm:"column:roles;type:varchar[]"`
	LastUpdated time.Time      `gorm:"column:last_updated"`
	LastSent    *time.Time     `gorm:"column:last_sent"`
}

type partyRead struct {
	Party      party      `gorm:"embedded"`
	TotalCount totalCount `gorm:"embedded"`
}

type partyStorageImpl struct {
	pg *pg.Storage
}

func (s *partyStorageImpl) l() kit.CLogger {
	return ocpi.L().Cmp("party-storage")
}

func newPartyStorage(pg *pg.Storage) *partyStorageImpl {
	return &partyStorageImpl{
		pg: pg,
	}
}

func (s *partyStorageImpl) CreateParty(ctx context.Context, p *domain.Party) error {
	s.l().Mth("create").C(ctx).F(kit.KV{"partyId": p.Id}).Dbg()
	err := s.pg.Instance.Create(s.toPartyDto(p)).Error
	if err != nil {
		return errors.ErrPartyStorageCreate(ctx, err)
	}
	return nil
}

func (s *partyStorageImpl) UpdateParty(ctx context.Context, p *domain.Party) error {
	s.l().Mth("update").C(ctx).F(kit.KV{"partyId": p.Id}).Dbg()
	err := s.pg.Instance.Scopes(update()).Save(s.toPartyDto(p)).Error
	if err != nil {
		return errors.ErrPartyStorageUpdate(ctx, err)
	}
	return nil
}

func (s *partyStorageImpl) GetPartyByExtId(ctx context.Context, extId domain.PartyExtId) (*domain.Party, error) {
	s.l().Mth("get-by-party").C(ctx).Dbg()
	if extId.PartyId == "" || extId.CountryCode == "" {
		return nil, nil
	}
	dto := &party{}
	res := s.pg.Instance.
		Where("party_id = ?", extId.PartyId).
		Where("country_code = ?", extId.CountryCode).
		Limit(1).Find(&dto)
	if res.Error != nil {
		return nil, errors.ErrPartyStorageGetDb(ctx, res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return s.toPartyDomain(dto), nil
}

func (s *partyStorageImpl) GetPartiesByPlatform(ctx context.Context, platformId string) ([]*domain.Party, error) {
	s.l().Mth("get-by-platform").C(ctx).Dbg()
	if platformId == "" {
		return nil, nil
	}
	var dtos []*party
	res := s.pg.Instance.Where("platform_id = ?", platformId).Find(&dtos)
	if res.Error != nil {
		return nil, errors.ErrPartyStorageGetDb(ctx, res.Error)
	}
	return s.toPartiesDomain(dtos), nil
}

func (s *partyStorageImpl) GetByRefId(ctx context.Context, refId string) (*domain.Party, error) {
	s.l().Mth("get-by-ref").C(ctx).Dbg()
	if refId == "" {
		return nil, nil
	}
	dto := &party{}
	res := s.pg.Instance.Where("ref_id = ?", refId).Limit(1).Find(&dto)
	if res.Error != nil {
		return nil, errors.ErrPartyStorageGetDb(ctx, res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return s.toPartyDomain(dto), nil
}

func (s *partyStorageImpl) GetParty(ctx context.Context, id string) (*domain.Party, error) {
	s.l().Mth("get-by-id").C(ctx).Dbg()
	if id == "" {
		return nil, nil
	}
	dto := &party{}
	res := s.pg.Instance.Where("id = ?", id).Limit(1).Find(&dto)
	if res.Error != nil {
		return nil, errors.ErrPartyStorageGetDb(ctx, res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return s.toPartyDomain(dto), nil
}

func (s *partyStorageImpl) MarkSentParties(ctx context.Context, date time.Time, partyIds ...string) error {
	s.l().Mth("mark-sent").C(ctx).Dbg()
	res := s.pg.Instance.Model(&party{}).Where("id in (?)", partyIds).Update("last_sent", date)
	if res.Error != nil {
		return errors.ErrPartyStorageUpdate(ctx, res.Error)
	}
	return nil
}

func (s *partyStorageImpl) Search(ctx context.Context, criteria *domain.PartySearchCriteria) (*domain.PartySearchResponse, error) {
	s.l().Mth("search").C(ctx).Dbg()

	rs := &domain.PartySearchResponse{
		PageResponse: domain.PageResponse{
			Limit: pagingLimit(criteria.PageRequest.Limit),
			Total: kit.IntPtr(0),
		},
	}

	// make query
	var dtosRead []*partyRead

	if err := s.pg.Instance.
		Scopes(s.buildSearchQuery(criteria), paging(criteria.PageRequest), orderByLastUpdated(true)).
		Find(&dtosRead).Error; err != nil {
		return nil, errors.ErrPartyStorageGetDb(ctx, err)
	}

	if len(dtosRead) == 0 {
		return rs, nil
	}

	dtos := make([]*party, 0, len(dtosRead))
	for _, p := range dtosRead {
		dtos = append(dtos, &p.Party)
	}

	rs.Items = s.toPartiesDomain(dtos)
	rs.Total = &dtosRead[0].TotalCount.TotalCount
	rs.NextPage = nextPage(criteria.PageRequest, rs.Total)

	return rs, nil
}

func (s *partyStorageImpl) DeletePartyByExtId(ctx context.Context, extId domain.PartyExtId) error {
	s.l().C(ctx).Mth("delete-ext").F(kit.KV{"partyId": extId.PartyId, "country": extId.CountryCode}).Dbg()
	if extId.PartyId == "" || extId.CountryCode == "" {
		return nil
	}
	if err := s.pg.Instance.Where(`party_id = ? and country_code = ?`, extId.PartyId, extId.CountryCode).
		Delete(&party{}).Error; err != nil {
		return errors.ErrPartyStorageDelete(ctx, err)
	}
	return nil
}

func (s *partyStorageImpl) buildSearchQuery(criteria *domain.PartySearchCriteria) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query := db.Table("parties").Select("parties.*, count(*) over() total_count")
		// populate conditions
		if len(criteria.Ids) > 0 {
			query = query.Where("id in (?)", criteria.Ids)
		}
		if len(criteria.IncPlatforms) > 0 {
			query = query.Where("platform_id in (?)", criteria.IncPlatforms)
		}
		if len(criteria.ExcPlatforms) > 0 {
			query = query.Where("platform_id not in (?)", criteria.ExcPlatforms)
		}
		if len(criteria.IncRoles) > 0 {
			query = query.Where(fmt.Sprintf("roles && '{%s}'::varchar[]", strings.Join(criteria.IncRoles, ",")))
		}
		if len(criteria.ExcRoles) > 0 {
			query = query.Where(fmt.Sprintf("not(roles && '{%s}'::varchar[])", strings.Join(criteria.ExcRoles, ",")))
		}
		if criteria.ExtId != nil {
			query = query.Where("party_id = ? and country_code = ?", criteria.ExtId.PartyId, criteria.ExtId.CountryCode)
		}
		if criteria.DateTo != nil {
			query = query.Where("last_updated <= ?", *criteria.DateTo)
		}
		if criteria.DateFrom != nil {
			query = query.Where("last_updated >= ?", *criteria.DateFrom)
		}
		if criteria.RefId != "" {
			query = query.Where("ref_id = ?", criteria.RefId)
		}
		return query
	}
}
