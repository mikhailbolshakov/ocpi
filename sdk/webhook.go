package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	service "github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
)

func (s *Sdk) CreateUpdateWebhook(ctx context.Context, rq *backend.Webhook) (*backend.Webhook, error) {
	l := service.L().C(ctx).Mth("create-update-wh").Dbg()

	rqJs, err := json.Marshal(rq)
	if err != nil {
		return nil, err
	}

	rs, err := s.POST(ctx, fmt.Sprintf("%s/webhooks", s.baseUrl), rqJs)
	if err != nil {
		return nil, err
	}

	var p *backend.Webhook
	err = json.Unmarshal(rs, &p)
	if err != nil {
		return nil, err
	}

	l.Dbg("ok")
	return p, nil
}

func (s *Sdk) GetWebhook(ctx context.Context, whId string) (*backend.Webhook, error) {
	l := service.L().C(ctx).Mth("get-wh").Dbg()

	rs, err := s.GET(ctx, fmt.Sprintf("%s/webhooks/%s", s.baseUrl, whId))
	if err != nil {
		return nil, err
	}

	var p *backend.Webhook
	err = json.Unmarshal(rs, &p)
	if err != nil {
		return nil, err
	}
	l.Dbg("ok")

	return p, nil
}

func (s *Sdk) DeleteWebhook(ctx context.Context, whId string) error {
	l := service.L().C(ctx).Mth("del-wh").Dbg()

	_, err := s.DELETE(ctx, fmt.Sprintf("%s/webhooks/%s", s.baseUrl, whId), nil)
	if err != nil {
		return err
	}
	l.Dbg("ok")
	return nil
}
