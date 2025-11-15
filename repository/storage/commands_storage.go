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

type command struct {
	pg.GormDto
	Id          string        `gorm:"column:id;primaryKey"`
	Status      string        `gorm:"column:status"`
	Cmd         string        `gorm:"column:cmd"`
	Deadline    time.Time     `gorm:"column:deadline"`
	AuthRef     string        `gorm:"column:auth_ref"`
	Details     *pgtype.JSONB `gorm:"column:details"`
	PartyId     string        `gorm:"column:party_id"`
	CountryCode string        `gorm:"column:country_code"`
	PlatformId  string        `gorm:"column:platform_id"`
	RefId       *string       `gorm:"column:ref_id"`
	LastUpdated time.Time     `gorm:"column:last_updated"`
	LastSent    *time.Time    `gorm:"column:last_sent"`
}

type commandRead struct {
	Command    command    `gorm:"embedded"`
	TotalCount totalCount `gorm:"embedded"`
}

type commandStorageImpl struct {
	pg *pg.Storage
}

func (s *commandStorageImpl) l() kit.CLogger {
	return ocpi.L().Cmp("cmd-storage")
}

func newCommandStorage(pg *pg.Storage) *commandStorageImpl {
	return &commandStorageImpl{
		pg: pg,
	}
}

func (s *commandStorageImpl) GetCommand(ctx context.Context, id string) (*domain.Command, error) {
	s.l().C(ctx).Mth("get-cmd").F(kit.KV{"cmdId": id}).Dbg()
	if id == "" {
		return nil, nil
	}
	dto := &command{}
	res := s.pg.Instance.Where("id = ?", id).Limit(1).Find(&dto)
	if res.Error != nil {
		return nil, errors.ErrCmdStorageGet(ctx, res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return s.toCommandDomain(dto), nil
}

func (s *commandStorageImpl) GetCommandByAuthRef(ctx context.Context, authRef string) (*domain.Command, error) {
	s.l().C(ctx).Mth("get-by-ref").F(kit.KV{"authRef": authRef}).Dbg()
	if authRef == "" {
		return nil, nil
	}
	dto := &command{}
	res := s.pg.Instance.Where("auth_ref = ?", authRef).Limit(1).Find(&dto)
	if res.Error != nil {
		return nil, errors.ErrCmdStorageGet(ctx, res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return s.toCommandDomain(dto), nil
}

func (s *commandStorageImpl) CreateCommand(ctx context.Context, cmd *domain.Command) error {
	s.l().C(ctx).Mth("create-cmd").F(kit.KV{"cmdId": cmd.Id}).Dbg()
	if err := s.pg.Instance.Create(s.toCommandDto(cmd)).Error; err != nil {
		return errors.ErrCmdStorageCreate(ctx, err)
	}
	return nil
}

func (s *commandStorageImpl) UpdateCommand(ctx context.Context, cmd *domain.Command) error {
	s.l().C(ctx).Mth("update-cmd").F(kit.KV{"cmdId": cmd.Id}).Dbg()
	if err := s.pg.Instance.Scopes(update()).Save(s.toCommandDto(cmd)).Error; err != nil {
		return errors.ErrCmdStorageUpdate(ctx, err)
	}
	return nil
}

func (s *commandStorageImpl) DeleteCommand(ctx context.Context, cmdId string) error {
	s.l().C(ctx).Mth("delete-cmd").F(kit.KV{"cmdId": cmdId}).Dbg()
	if cmdId == "" {
		return nil
	}
	if err := s.pg.Instance.Delete(&command{Id: cmdId}).Error; err != nil {
		return errors.ErrCmdStorageDelete(ctx, err)
	}
	return nil
}

func (s *commandStorageImpl) DeleteCommandsByExt(ctx context.Context, extId domain.PartyExtId) error {
	s.l().C(ctx).Mth("delete-ext").F(kit.KV{"partyId": extId.PartyId, "country": extId.CountryCode}).Dbg()
	if extId.PartyId == "" || extId.CountryCode == "" {
		return nil
	}
	if err := s.pg.Instance.Where(`party_id = ? and country_code = ?`, extId.PartyId, extId.CountryCode).
		Delete(&command{}).Error; err != nil {
		return errors.ErrCmdStorageDelete(ctx, err)
	}
	return nil
}

func (s *commandStorageImpl) SearchCommands(ctx context.Context, cr *domain.CommandSearchCriteria) (*domain.CommandSearchResponse, error) {
	s.l().Mth("search-cmd").C(ctx).Dbg()

	rs := &domain.CommandSearchResponse{
		PageResponse: domain.PageResponse{
			Total: kit.IntPtr(0),
		},
	}

	q := s.pg.Instance.Scopes(s.buildSearchQuery(cr), orderByLastUpdated(true))

	if !cr.RetrieveAll {
		rs.PageResponse.Limit = pagingLimit(cr.PageRequest.Limit)
		q = q.Scopes(paging(cr.PageRequest))
	}

	// make query
	var dtosRead []*commandRead

	if err := q.Find(&dtosRead).Error; err != nil {
		return nil, errors.ErrCmdStorageGet(ctx, err)
	}

	if len(dtosRead) == 0 {
		return rs, nil
	}

	dtos := make([]*command, 0, len(dtosRead))
	for _, p := range dtosRead {
		dtos = append(dtos, &p.Command)
	}

	rs.Items = s.toCommandsDomain(dtos)
	rs.Total = &dtosRead[0].TotalCount.TotalCount

	if !cr.RetrieveAll {
		rs.NextPage = nextPage(cr.PageRequest, rs.Total)
	}

	return rs, nil
}

func (s *commandStorageImpl) buildSearchQuery(criteria *domain.CommandSearchCriteria) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query := db.Table("commands").Select("commands.*, count(*) over() total_count")
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
		if criteria.AuthRef != "" {
			query = query.Where("auth_ref = ?", criteria.AuthRef)
		}
		if criteria.Cmd != "" {
			query = query.Where("cmd = ?", criteria.Cmd)
		}
		if len(criteria.Statuses) > 0 {
			query = query.Where("status in (?)", criteria.Statuses)
		}
		if criteria.DeadlineLE != nil {
			query = query.Where("deadline <= ?", *criteria.DeadlineLE)
		}
		if criteria.ReservationId != "" {
			query = query.Where("reservation_id = ?", criteria.ReservationId)
		}
		return query
	}
}
