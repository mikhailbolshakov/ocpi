package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mikhailbolshakov/kit"
	service "github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
)

func (s *Sdk) PostPlatform(ctx context.Context, rq *backend.PlatformRequest) (*backend.Platform, error) {
	l := service.L().C(ctx).Mth("post-platform").Dbg()

	rqJs, err := json.Marshal(rq)
	if err != nil {
		return nil, err
	}

	rs, err := s.POST(ctx, fmt.Sprintf("%s/platforms", s.baseUrl), rqJs)
	if err != nil {
		return nil, err
	}

	var p *backend.Platform
	err = json.Unmarshal(rs, &p)
	if err != nil {
		return nil, err
	}
	l.Dbg("ok")
	return p, nil
}

func (s *Sdk) SetPlatformStatus(ctx context.Context, platformId, status string) (*backend.Platform, error) {
	l := service.L().C(ctx).Mth("set-platform-status").Dbg()

	rs, err := s.PUT(ctx, fmt.Sprintf("%s/platforms/%s/status?status=%s", s.baseUrl, platformId, status), nil)
	if err != nil {
		return nil, err
	}

	var p *backend.Platform
	err = json.Unmarshal(rs, &p)
	if err != nil {
		return nil, err
	}
	l.Dbg("ok")
	return p, nil
}

func (s *Sdk) GetPlatform(ctx context.Context, platformId string) (*backend.Platform, error) {
	l := service.L().C(ctx).Mth("get-platform").F(kit.KV{"platformId": platformId}).Dbg()

	rs, err := s.GET(ctx, fmt.Sprintf("%s/platforms/%s", s.baseUrl, platformId))
	if err != nil {
		return nil, err
	}

	var p *backend.Platform
	err = json.Unmarshal(rs, &p)
	if err != nil {
		return nil, err
	}
	l.Dbg("ok")
	return p, nil
}

func (s *Sdk) ConnectPlatform(ctx context.Context, platformId string) (*backend.Platform, error) {
	l := service.L().C(ctx).Mth("connect-platform").Dbg()

	rs, err := s.POST(ctx, fmt.Sprintf("%s/platforms/%s/connections", s.baseUrl, platformId), nil)
	if err != nil {
		return nil, err
	}

	var p *backend.Platform
	err = json.Unmarshal(rs, &p)
	if err != nil {
		return nil, err
	}
	l.Dbg("ok")
	return p, nil
}

func (s *Sdk) UpdateConnectionPlatform(ctx context.Context, platformId string) (*backend.Platform, error) {
	l := service.L().C(ctx).Mth("update-connection-platform").Dbg()

	rs, err := s.PUT(ctx, fmt.Sprintf("%s/platforms/%s/connections", s.baseUrl, platformId), nil)
	if err != nil {
		return nil, err
	}

	var p *backend.Platform
	err = json.Unmarshal(rs, &p)
	if err != nil {
		return nil, err
	}
	l.Dbg("ok")
	return p, nil
}
