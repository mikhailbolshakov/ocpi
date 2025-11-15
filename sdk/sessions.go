package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	service "github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
)

func (s *Sdk) PutSession(ctx context.Context, rq *backend.Session) error {
	l := service.L().C(ctx).Mth("put-sess").Dbg()

	rqJs, err := json.Marshal(rq)
	if err != nil {
		return err
	}

	_, err = s.POST(ctx, fmt.Sprintf("%s/backend/sessions", s.baseUrl), rqJs)
	if err != nil {
		return err
	}
	l.Dbg("ok")
	return nil
}

func (s *Sdk) PatchSession(ctx context.Context, rq *backend.Session) error {
	l := service.L().C(ctx).Mth("patch-sess").Dbg()

	rqJs, err := json.Marshal(rq)
	if err != nil {
		return err
	}

	_, err = s.PATCH(ctx, fmt.Sprintf("%s/backend/sessions/%s", s.baseUrl, rq.Id), rqJs)
	if err != nil {
		return err
	}
	l.Dbg("ok")
	return nil
}

func (s *Sdk) GetSession(ctx context.Context, sessId string) (*backend.Session, error) {
	service.L().C(ctx).Mth("get-sess").Dbg()

	rs, err := s.GET(ctx, fmt.Sprintf("%s/backend/sessions/%s", s.baseUrl, sessId))
	if err != nil {
		return nil, err
	}

	var p *backend.Session
	err = json.Unmarshal(rs, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Sdk) SearchSessions(ctx context.Context, params map[string]interface{}) (*backend.SessionSearchResponse, error) {
	service.L().C(ctx).Mth("search-sess").Dbg()

	rs, err := s.GET(ctx, fmt.Sprintf("%s/backend/sessions/search/query%s", s.baseUrl, s.toUrlParams(params)))
	if err != nil {
		return nil, err
	}

	var p *backend.SessionSearchResponse
	err = json.Unmarshal(rs, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
