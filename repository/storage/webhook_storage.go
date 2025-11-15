package storage

import (
	"context"
	"github.com/lib/pq"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/storages/pg"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/errors"
)

type webhook struct {
	pg.GormDto
	Id     string         `gorm:"column:id"`
	ApiKey string         `gorm:"column:api_key"`
	Url    string         `gorm:"column:url"`
	Events pq.StringArray `gorm:"column:events;type:varchar[]"`
}

type webhookStorageImpl struct {
	pg *pg.Storage
}

func (s *webhookStorageImpl) l() kit.CLogger {
	return ocpi.L().Cmp("wh-storage")
}

func newWebhookStorage(pg *pg.Storage) *webhookStorageImpl {
	return &webhookStorageImpl{
		pg: pg,
	}
}

func (s *webhookStorageImpl) CreateWebhook(ctx context.Context, wh *backend.Webhook) error {
	s.l().Mth("create").C(ctx).F(kit.KV{"whId": wh.Id}).Dbg()
	err := s.pg.Instance.Create(s.toWebhookDto(wh)).Error
	if err != nil {
		return errors.ErrWhStorageCreate(ctx, err)
	}
	return nil
}

func (s *webhookStorageImpl) UpdateWebhook(ctx context.Context, wh *backend.Webhook) error {
	s.l().Mth("update").C(ctx).F(kit.KV{"whId": wh.Id}).Dbg()
	err := s.pg.Instance.Scopes(update()).Save(s.toWebhookDto(wh)).Error
	if err != nil {
		return errors.ErrWhStorageUpdate(ctx, err)
	}
	return nil
}

func (s *webhookStorageImpl) MergeWebhook(ctx context.Context, wh *backend.Webhook) error {
	s.l().Mth("merge").C(ctx).F(kit.KV{"whId": wh.Id}).Dbg()
	err := s.pg.Instance.Scopes(merge()).Create(s.toWebhookDto(wh)).Error
	if err != nil {
		return errors.ErrWhStorageMerge(ctx, err)
	}
	return nil
}

func (s *webhookStorageImpl) DeleteWebhook(ctx context.Context, whId string) error {
	s.l().Mth("update").C(ctx).F(kit.KV{"whId": whId}).Dbg()
	err := s.pg.Instance.Delete(&webhook{Id: whId}).Error
	if err != nil {
		return errors.ErrWhStorageDelete(ctx, err)
	}
	return nil
}

func (s *webhookStorageImpl) SearchWebhook(ctx context.Context, cr *backend.SearchWebhookCriteria) ([]*backend.Webhook, error) {
	s.l().Mth("search").C(ctx).Dbg()
	query := s.pg.Instance.Model(&webhook{})
	if cr.Event != "" {
		query = query.Where("? = ANY(events)", cr.Event)
	}
	// make query
	var dtos []*webhook
	if err := query.Find(&dtos).Error; err != nil {
		return nil, errors.ErrWhStorageGetDb(ctx, err)
	}
	return s.toWebhooksBackend(dtos), nil
}

func (s *webhookStorageImpl) GetWebhook(ctx context.Context, whId string) (*backend.Webhook, error) {
	s.l().Mth("get").C(ctx).Dbg()
	if whId == "" {
		return nil, nil
	}
	dto := &webhook{Id: whId}
	res := s.pg.Instance.Limit(1).Find(&dto)
	if res.Error != nil {
		return nil, errors.ErrWhStorageGetDb(ctx, res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return s.toWebhookBackend(dto), nil
}
