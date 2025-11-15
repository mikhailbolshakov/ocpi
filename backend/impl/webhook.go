package impl

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/errors"
)

type webhookImpl struct {
	storage backend.WebhookStorage
}

func NewWebhookService(storage backend.WebhookStorage) backend.WebhookService {
	return &webhookImpl{
		storage: storage,
	}
}

func (w *webhookImpl) l() kit.CLogger {
	return ocpi.L().Cmp("wh-svc")
}

func (w *webhookImpl) CreateUpdate(ctx context.Context, wh *backend.Webhook) (*backend.Webhook, error) {
	w.l().C(ctx).Mth("create-update").Dbg()

	if wh.Id == "" {
		wh.Id = kit.NewId()
	}

	err := w.validate(ctx, wh)
	if err != nil {
		return nil, err
	}

	err = w.storage.MergeWebhook(ctx, wh)
	if err != nil {
		return nil, err
	}

	return wh, nil
}

func (w *webhookImpl) Create(ctx context.Context, wh *backend.Webhook) (*backend.Webhook, error) {
	w.l().C(ctx).Mth("create").Dbg()

	if wh.Id == "" {
		wh.Id = kit.NewId()
	} else {
		stored, err := w.storage.GetWebhook(ctx, wh.Id)
		if err != nil {
			return nil, err
		}
		if stored != nil {
			return nil, errors.ErrWhAlreadyExists(ctx)
		}
	}

	err := w.validate(ctx, wh)
	if err != nil {
		return nil, err
	}

	err = w.storage.CreateWebhook(ctx, wh)
	if err != nil {
		return nil, err
	}

	return wh, nil
}

func (w *webhookImpl) Update(ctx context.Context, wh *backend.Webhook) (*backend.Webhook, error) {
	w.l().C(ctx).Mth("update").Dbg()

	err := w.validate(ctx, wh)
	if err != nil {
		return nil, err
	}

	stored, err := w.storage.GetWebhook(ctx, wh.Id)
	if err != nil {
		return nil, err
	}
	if stored == nil {
		return nil, errors.ErrWhNotFound(ctx)
	}

	err = w.storage.UpdateWebhook(ctx, wh)
	if err != nil {
		return nil, err
	}

	return wh, nil
}

func (w *webhookImpl) Delete(ctx context.Context, whId string) error {
	w.l().C(ctx).Mth("delete").Dbg()

	if whId == "" {
		return errors.ErrWhIdEmpty(ctx)
	}

	stored, err := w.storage.GetWebhook(ctx, whId)
	if err != nil {
		return err
	}
	if stored == nil {
		return errors.ErrWhNotFound(ctx)
	}

	err = w.storage.DeleteWebhook(ctx, whId)
	return err
}

func (w *webhookImpl) Search(ctx context.Context, cr *backend.SearchWebhookCriteria) ([]*backend.Webhook, error) {
	w.l().C(ctx).Mth("search").Dbg()
	return w.storage.SearchWebhook(ctx, cr)
}

func (w *webhookImpl) Get(ctx context.Context, whId string) (*backend.Webhook, error) {
	w.l().C(ctx).Mth("get").Dbg()

	if whId == "" {
		return nil, errors.ErrWhIdEmpty(ctx)
	}

	return w.storage.GetWebhook(ctx, whId)
}

func (w *webhookImpl) validate(ctx context.Context, wh *backend.Webhook) error {
	if wh.Id == "" {
		return errors.ErrWhIdEmpty(ctx)
	}
	if wh.ApiKey == "" {
		return errors.ErrWhApiKeyEmpty(ctx)
	}
	if len(wh.Events) == 0 {
		return errors.ErrWhEventsEmpty(ctx)
	}
	if wh.Url == "" || !kit.IsUrlValid(wh.Url) {
		return errors.ErrWhUrlInvalid(ctx)
	}
	return nil
}
