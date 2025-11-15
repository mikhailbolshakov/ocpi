package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	service "github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
)

func (s *Sdk) PutToken(ctx context.Context, rq *backend.Token) error {
	l := service.L().C(ctx).Mth("put-tkn").Dbg()

	rqJs, err := json.Marshal(rq)
	if err != nil {
		return err
	}

	_, err = s.POST(ctx, fmt.Sprintf("%s/backend/tokens", s.baseUrl), rqJs)
	if err != nil {
		return err
	}
	l.Dbg("ok")
	return nil
}

func (s *Sdk) GetToken(ctx context.Context, tknId string) (*backend.Token, error) {
	service.L().C(ctx).Mth("get-tkn").Dbg()

	rs, err := s.GET(ctx, fmt.Sprintf("%s/backend/tokens/%s", s.baseUrl, tknId))
	if err != nil {
		return nil, err
	}

	var p *backend.Token
	err = json.Unmarshal(rs, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Sdk) SearchTokens(ctx context.Context, params map[string]interface{}) (*backend.TokenSearchResponse, error) {
	service.L().C(ctx).Mth("search-tkn").Dbg()

	rs, err := s.GET(ctx, fmt.Sprintf("%s/backend/tokens/search/query%s", s.baseUrl, s.toUrlParams(params)))
	if err != nil {
		return nil, err
	}

	var p *backend.TokenSearchResponse
	err = json.Unmarshal(rs, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
