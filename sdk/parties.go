package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	service "github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
)

func (s *Sdk) PutParty(ctx context.Context, rq *backend.Party) error {
	l := service.L().C(ctx).Mth("put-party").Dbg()

	rqJs, err := json.Marshal(rq)
	if err != nil {
		return err
	}

	_, err = s.POST(ctx, fmt.Sprintf("%s/backend/parties", s.baseUrl), rqJs)
	if err != nil {
		return err
	}

	l.Dbg("ok")
	return nil
}

func (s *Sdk) GetParty(ctx context.Context, partyId string) (*backend.Party, error) {
	l := service.L().C(ctx).Mth("get-party").Dbg()

	rs, err := s.GET(ctx, fmt.Sprintf("%s/backend/parties/%s", s.baseUrl, partyId))
	if err != nil {
		return nil, err
	}

	var p *backend.Party
	err = json.Unmarshal(rs, &p)
	if err != nil {
		return nil, err
	}
	l.Dbg("ok")
	return p, nil
}

func (s *Sdk) SearchParties(ctx context.Context, params map[string]interface{}) (*backend.PartySearchResponse, error) {
	service.L().C(ctx).Mth("search-parties").Dbg()

	rs, err := s.GET(ctx, fmt.Sprintf("%s/backend/parties/search/query%s", s.baseUrl, s.toUrlParams(params)))
	if err != nil {
		return nil, err
	}

	var p *backend.PartySearchResponse
	err = json.Unmarshal(rs, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Sdk) PullParties(ctx context.Context) error {
	l := service.L().C(ctx).Mth("pull-parties").Dbg()

	_, err := s.POST(ctx, fmt.Sprintf("%s/backend/parties/pull", s.baseUrl), nil)
	if err != nil {
		return err
	}

	l.Dbg("ok")
	return nil
}
