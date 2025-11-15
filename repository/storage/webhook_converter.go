package storage

import (
	"github.com/mikhailbolshakov/ocpi/backend"
)

func (s *webhookStorageImpl) toWebhookDto(wh *backend.Webhook) *webhook {
	if wh == nil {
		return nil
	}
	return &webhook{
		Id:     wh.Id,
		ApiKey: wh.ApiKey,
		Url:    wh.Url,
		Events: wh.Events,
	}
}

func (s *webhookStorageImpl) toWebhookBackend(dto *webhook) *backend.Webhook {
	if dto == nil {
		return nil
	}
	return &backend.Webhook{
		Id:     dto.Id,
		ApiKey: dto.ApiKey,
		Events: dto.Events,
		Url:    dto.Url,
	}
}

func (s *webhookStorageImpl) toWebhooksBackend(dtos []*webhook) []*backend.Webhook {
	var r []*backend.Webhook
	for _, dto := range dtos {
		r = append(r, s.toWebhookBackend(dto))
	}
	return r
}
