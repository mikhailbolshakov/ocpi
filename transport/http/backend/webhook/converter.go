package webhook

import "github.com/mikhailbolshakov/ocpi/backend"

func (c *ctrlImpl) toWebhookRequestDomain(wh *Webhook) *backend.Webhook {
	if wh == nil {
		return nil
	}
	return &backend.Webhook{
		Id:     wh.Id,
		ApiKey: wh.ApiKey,
		Events: wh.Events,
		Url:    wh.Url,
	}
}

func (c *ctrlImpl) toWebhookApi(wh *backend.Webhook) *Webhook {
	if wh == nil {
		return nil
	}
	return &Webhook{
		Id:     wh.Id,
		ApiKey: wh.ApiKey,
		Events: wh.Events,
		Url:    wh.Url,
	}
}

func (c *ctrlImpl) toWebhooksApi(whs []*backend.Webhook) []*Webhook {
	var r []*Webhook
	for _, wh := range whs {
		r = append(r, c.toWebhookApi(wh))
	}
	return r
}
