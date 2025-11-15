package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mikhailbolshakov/kit"
	service "github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
)

func (s *Sdk) PutLocation(ctx context.Context, rq *backend.Location) error {
	l := service.L().C(ctx).Mth("put-loc").Dbg()

	rqJs, err := json.Marshal(rq)
	if err != nil {
		return err
	}

	_, err = s.POST(ctx, fmt.Sprintf("%s/backend/locations", s.baseUrl), rqJs)
	if err != nil {
		return err
	}
	l.Dbg("ok")
	return nil
}

func (s *Sdk) PullLocations(ctx context.Context, rq *backend.PullRequest) error {
	l := service.L().C(ctx).Mth("pull-loc").Dbg()

	rqJs, err := json.Marshal(rq)
	if err != nil {
		return err
	}

	_, err = s.POST(ctx, fmt.Sprintf("%s/backend/locations/pull", s.baseUrl), rqJs)
	if err != nil {
		return err
	}
	l.Dbg("ok")
	return nil
}

func (s *Sdk) GetLocation(ctx context.Context, locId string) (*backend.Location, error) {
	service.L().C(ctx).Mth("get-loc").Dbg()

	rs, err := s.GET(ctx, fmt.Sprintf("%s/backend/locations/%s", s.baseUrl, locId))
	if err != nil {
		return nil, err
	}

	var p *backend.Location
	err = json.Unmarshal(rs, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Sdk) PutEvse(ctx context.Context, rq *backend.Evse) error {
	l := service.L().C(ctx).Mth("put-evse").Dbg()

	rqJs, err := json.Marshal(rq)
	if err != nil {
		return err
	}

	_, err = s.POST(ctx, fmt.Sprintf("%s/backend/locations/%s/evses", s.baseUrl, rq.LocationId), rqJs)
	if err != nil {
		return err
	}
	l.Dbg("ok")
	return nil
}

func (s *Sdk) SetEvseStatus(ctx context.Context, locId, evseId, status string) error {
	l := service.L().C(ctx).Mth("set-evse-status").Dbg()
	_, err := s.POST(ctx, fmt.Sprintf("%s/backend/locations/%s/evses/%s/status?status=%s", s.baseUrl, locId, evseId, status), nil)
	if err != nil {
		return err
	}
	l.Dbg("ok")
	return nil
}

func (s *Sdk) GetEvse(ctx context.Context, locId, evseId string) (*backend.Evse, error) {
	service.L().C(ctx).Mth("get-evse").Dbg()

	rs, err := s.GET(ctx, fmt.Sprintf("%s/backend/locations/%s/evses/%s", s.baseUrl, locId, evseId))
	if err != nil {
		return nil, err
	}

	var p *backend.Evse
	err = json.Unmarshal(rs, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Sdk) PutConnector(ctx context.Context, rq *backend.Connector) error {
	l := service.L().C(ctx).Mth("put-connector").Dbg()

	rqJs, err := json.Marshal(rq)
	if err != nil {
		return err
	}

	_, err = s.POST(ctx, fmt.Sprintf("%s/backend/locations/%s/evses/%s/connectors", s.baseUrl, rq.LocationId, rq.EvseId), rqJs)
	if err != nil {
		return err
	}
	l.Dbg("ok")
	return nil
}

func (s *Sdk) UpdateConnector(ctx context.Context, rq *backend.Connector) (*backend.Connector, error) {
	l := service.L().C(ctx).Mth("update-evse").Dbg()

	rqJs, err := json.Marshal(rq)
	if err != nil {
		return nil, err
	}

	rs, err := s.PUT(ctx, fmt.Sprintf("%s/backend/locations/%s/evses/%s/connectors/%s", s.baseUrl, rq.LocationId, rq.EvseId, rq.Id), rqJs)
	if err != nil {
		return nil, err
	}

	var p *backend.Connector
	err = json.Unmarshal(rs, &p)
	if err != nil {
		return nil, err
	}
	l.F(kit.KV{"id": p.Id}).Dbg("ok")
	return p, nil
}

func (s *Sdk) GetConnector(ctx context.Context, locId, evseId, conId string) (*backend.Connector, error) {
	service.L().C(ctx).Mth("get-evse").Dbg()

	rs, err := s.GET(ctx, fmt.Sprintf("%s/backend/locations/%s/evses/%s/connectors/%s", s.baseUrl, locId, evseId, conId))
	if err != nil {
		return nil, err
	}

	var p *backend.Connector
	err = json.Unmarshal(rs, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Sdk) SearchLocations(ctx context.Context, params map[string]interface{}) (*backend.LocationSearchResponse, error) {
	service.L().C(ctx).Mth("search-locations").Dbg()

	rs, err := s.GET(ctx, fmt.Sprintf("%s/backend/locations/search/query%s", s.baseUrl, s.toUrlParams(params)))
	if err != nil {
		return nil, err
	}

	var p *backend.LocationSearchResponse
	err = json.Unmarshal(rs, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Sdk) SearchEvses(ctx context.Context, params map[string]interface{}) (*backend.EvseSearchResponse, error) {
	service.L().C(ctx).Mth("search-evses").Dbg()

	rs, err := s.GET(ctx, fmt.Sprintf("%s/backend/evses/search/query%s", s.baseUrl, s.toUrlParams(params)))
	if err != nil {
		return nil, err
	}

	var p *backend.EvseSearchResponse
	err = json.Unmarshal(rs, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Sdk) SearchConnectors(ctx context.Context, params map[string]interface{}) (*backend.ConnectorSearchResponse, error) {
	service.L().C(ctx).Mth("search-connectors").Dbg()

	rs, err := s.GET(ctx, fmt.Sprintf("%s/backend/connectors/search/query%s", s.baseUrl, s.toUrlParams(params)))
	if err != nil {
		return nil, err
	}

	var p *backend.ConnectorSearchResponse
	err = json.Unmarshal(rs, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
