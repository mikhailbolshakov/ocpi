package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	service "github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
)

func (s *Sdk) PutTariff(ctx context.Context, rq *backend.Tariff) error {
	l := service.L().C(ctx).Mth("put-trf").Dbg()

	rqJs, err := json.Marshal(rq)
	if err != nil {
		return err
	}

	_, err = s.POST(ctx, fmt.Sprintf("%s/backend/tariffs", s.baseUrl), rqJs)
	if err != nil {
		return err
	}
	l.Dbg("ok")
	return nil
}

func (s *Sdk) PullTariffs(ctx context.Context, rq *backend.PullRequest) error {
	l := service.L().C(ctx).Mth("pull-trfs").Dbg()

	rqJs, err := json.Marshal(rq)
	if err != nil {
		return err
	}

	_, err = s.POST(ctx, fmt.Sprintf("%s/backend/tariffs/pull", s.baseUrl), rqJs)
	if err != nil {
		return err
	}
	l.Dbg("ok")
	return nil
}

func (s *Sdk) GetTariff(ctx context.Context, trfId string) (*backend.Tariff, error) {
	service.L().C(ctx).Mth("get-trf").Dbg()

	rs, err := s.GET(ctx, fmt.Sprintf("%s/backend/tariffs/%s", s.baseUrl, trfId))
	if err != nil {
		return nil, err
	}

	var p *backend.Tariff
	err = json.Unmarshal(rs, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Sdk) SearchTariffs(ctx context.Context, params map[string]interface{}) (*backend.TariffSearchResponse, error) {
	service.L().C(ctx).Mth("search-trfs").Dbg()

	rs, err := s.GET(ctx, fmt.Sprintf("%s/backend/tariffs/search/query%s", s.baseUrl, s.toUrlParams(params)))
	if err != nil {
		return nil, err
	}

	var p *backend.TariffSearchResponse
	err = json.Unmarshal(rs, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
