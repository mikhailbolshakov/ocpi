package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgtype"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/storages/pg"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/errors"
)

type platformDetails struct {
	VersionInfo domain.VersionInfo      `json:"ver,omitempty"`
	Endpoints   domain.ModuleEndpoints  `json:"eps,omitempty"`
	Protocol    *domain.ProtocolDetails `json:"protocol,omitempty"`
	TokenBase64 *bool                   `json:"tokenBase64,omitempty"`
}

type platform struct {
	pg.GormDto
	Id      string        `gorm:"column:id"`
	TokenA  string        `gorm:"column:token_a"`
	TokenB  *string       `gorm:"column:token_b"`
	TokenC  *string       `gorm:"column:token_c"`
	Name    string        `gorm:"column:name"`
	Role    string        `gorm:"column:role"`
	Status  string        `gorm:"column:status"`
	Remote  bool          `gorm:"column:remote"`
	Details *pgtype.JSONB `gorm:"column:details"`
}

type platformStorageImpl struct {
	pg *pg.Storage
}

func (s *platformStorageImpl) l() kit.CLogger {
	return ocpi.L().Cmp("platform-storage")
}

func newPlatformStorage(pg *pg.Storage) *platformStorageImpl {
	return &platformStorageImpl{
		pg: pg,
	}
}

func (s *platformStorageImpl) CreatePlatform(ctx context.Context, p *domain.Platform) error {
	s.l().Mth("create").C(ctx).F(kit.KV{"platformId": p.Id}).Dbg()
	err := s.pg.Instance.Create(s.toPlatformDto(p)).Error
	if err != nil {
		return errors.ErrPlatformStorageCreate(ctx, err)
	}
	return nil
}

func (s *platformStorageImpl) UpdatePlatform(ctx context.Context, p *domain.Platform) error {
	s.l().Mth("update").C(ctx).F(kit.KV{"platformId": p.Id}).Dbg()
	err := s.pg.Instance.Scopes(update()).Save(s.toPlatformDto(p)).Error
	if err != nil {
		return errors.ErrPlatformStorageUpdate(ctx, err)
	}
	return nil
}

func (s *platformStorageImpl) DeletePlatform(ctx context.Context, p *domain.Platform) error {
	s.l().Mth("delete").C(ctx).F(kit.KV{"platformId": p.Id}).Dbg()
	err := s.pg.Instance.Delete(&platform{Id: p.Id}).Error
	if err != nil {
		return errors.ErrPlatformStorageDelete(ctx, err)
	}
	return nil
}

func (s *platformStorageImpl) GetPlatform(ctx context.Context, id string) (*domain.Platform, error) {
	s.l().Mth("get").C(ctx).F(kit.KV{"platformId": id}).Dbg()
	if id == "" {
		return nil, nil
	}
	dto := &platform{Id: id}
	res := s.pg.Instance.Limit(1).Find(&dto)
	if res.Error != nil {
		return nil, errors.ErrPlatformStorageGetDb(ctx, res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return s.toPlatformDomain(dto), nil
}

func (s *platformStorageImpl) GetPlatformByTokenA(ctx context.Context, tokenA domain.PlatformToken) (*domain.Platform, error) {
	s.l().Mth("get-by-token-a").C(ctx).Dbg()
	return s.getPlatformByToken(ctx, "token_a", tokenA)
}

func (s *platformStorageImpl) GetPlatformByTokenB(ctx context.Context, token domain.PlatformToken) (*domain.Platform, error) {
	s.l().Mth("get-by-token-b").C(ctx).Dbg()
	return s.getPlatformByToken(ctx, "token_b", token)
}

func (s *platformStorageImpl) GetPlatformByTokenC(ctx context.Context, token domain.PlatformToken) (*domain.Platform, error) {
	s.l().Mth("get-by-token").C(ctx).Dbg()
	return s.getPlatformByToken(ctx, "token_c", token)
}

func (s *platformStorageImpl) GetPlatformByToken(ctx context.Context, token domain.PlatformToken) (*domain.Platform, error) {
	s.l().Mth("get-by-token").C(ctx).Dbg()
	if token == "" {
		return nil, nil
	}
	dto := &platform{}
	res := s.pg.Instance.Where("? in (token_b, token_c)", token).Limit(1).Find(&dto)
	if res.Error != nil {
		return nil, errors.ErrPlatformStorageGetDb(ctx, res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return s.toPlatformDomain(dto), nil
}

func (s *platformStorageImpl) SearchPlatforms(ctx context.Context, cr *domain.PlatformSearchCriteria) ([]*domain.Platform, error) {
	s.l().Mth("search").C(ctx).Dbg()
	if cr == nil {
		return nil, nil
	}
	var dtos []*platform
	q := s.pg.Instance.Model(&platform{})
	if len(cr.Statuses) > 0 {
		q = q.Where("status in (?)", cr.Statuses)
	}
	if cr.Remote != nil {
		q = q.Where("remote = ?", *cr.Remote)
	}
	if len(cr.Roles) > 0 {
		q = q.Where("role in (?)", cr.Roles)
	}
	if len(cr.ExcRoles) > 0 {
		q = q.Where("role not in (?)", cr.ExcRoles)
	}
	if len(cr.IncIds) > 0 {
		q = q.Where("id in (?)", cr.IncIds)
	}
	if len(cr.ExcIds) > 0 {
		q = q.Where("id not in (?)", cr.ExcIds)
	}
	res := q.Find(&dtos)
	if res.Error != nil {
		return nil, errors.ErrPlatformStorageGetDb(ctx, res.Error)
	}
	return s.toPlatformsDomain(dtos), nil
}

func (s *platformStorageImpl) getPlatformByToken(ctx context.Context, field string, token domain.PlatformToken) (*domain.Platform, error) {
	if token == "" {
		return nil, nil
	}
	dto := &platform{}
	res := s.pg.Instance.Where(fmt.Sprintf("%s = ?", field), token).Limit(1).Find(&dto)
	if res.Error != nil {
		return nil, errors.ErrPlatformStorageGetDb(ctx, res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return s.toPlatformDomain(dto), nil
}
