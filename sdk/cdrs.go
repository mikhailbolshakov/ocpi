package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	service "github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
)

func (s *Sdk) PutCdr(ctx context.Context, rq *backend.Cdr) error {
	l := service.L().C(ctx).Mth("put-cdr").Dbg()

	rqJs, err := json.Marshal(rq)
	if err != nil {
		return err
	}

	_, err = s.POST(ctx, fmt.Sprintf("%s/backend/cdrs", s.baseUrl), rqJs)
	if err != nil {
		return err
	}
	l.Dbg("ok")
	return nil
}

func (s *Sdk) GetCdr(ctx context.Context, cdrId string) (*backend.Cdr, error) {
	service.L().C(ctx).Mth("get-cdr").Dbg()

	rs, err := s.GET(ctx, fmt.Sprintf("%s/backend/cdrs/%s", s.baseUrl, cdrId))
	if err != nil {
		return nil, err
	}

	var p *backend.Cdr
	err = json.Unmarshal(rs, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Sdk) SearchCdrs(ctx context.Context, params map[string]interface{}) (*backend.CdrSearchResponse, error) {
	service.L().C(ctx).Mth("search-cdrs").Dbg()

	rs, err := s.GET(ctx, fmt.Sprintf("%s/backend/cdrs/search/query%s", s.baseUrl, s.toUrlParams(params)))
	if err != nil {
		return nil, err
	}

	var p *backend.CdrSearchResponse
	err = json.Unmarshal(rs, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
