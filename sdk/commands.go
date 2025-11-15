package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	service "github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
)

func (s *Sdk) StartSession(ctx context.Context, rq *backend.StartSessionRequest) error {
	l := service.L().C(ctx).Mth("start-sess").Dbg()

	rqJs, err := json.Marshal(rq)
	if err != nil {
		return err
	}

	_, err = s.POST(ctx, fmt.Sprintf("%s/backend/commands/sessions/start", s.baseUrl), rqJs)
	if err != nil {
		return err
	}
	l.Dbg("ok")
	return nil
}

func (s *Sdk) StopSession(ctx context.Context, rq *backend.StopSessionRequest) error {
	l := service.L().C(ctx).Mth("stop-sess").Dbg()

	rqJs, err := json.Marshal(rq)
	if err != nil {
		return err
	}

	_, err = s.POST(ctx, fmt.Sprintf("%s/backend/commands/sessions/stop", s.baseUrl), rqJs)
	if err != nil {
		return err
	}
	l.Dbg("ok")
	return nil
}

func (s *Sdk) Reservation(ctx context.Context, rq *backend.ReserveNowRequest) error {
	l := service.L().C(ctx).Mth("reservation").Dbg()

	rqJs, err := json.Marshal(rq)
	if err != nil {
		return err
	}

	_, err = s.POST(ctx, fmt.Sprintf("%s/backend/commands/reservations", s.baseUrl), rqJs)
	if err != nil {
		return err
	}
	l.Dbg("ok")
	return nil
}

func (s *Sdk) CancelReservations(ctx context.Context, rq *backend.CancelReservationRequest) error {
	l := service.L().C(ctx).Mth("reservation-cancel").Dbg()

	rqJs, err := json.Marshal(rq)
	if err != nil {
		return err
	}

	_, err = s.DELETE(ctx, fmt.Sprintf("%s/backend/commands/reservations", s.baseUrl), rqJs)
	if err != nil {
		return err
	}
	l.Dbg("ok")
	return nil
}

func (s *Sdk) PutCommandResponse(ctx context.Context, cmdId string, rq *backend.CommandResponse) error {
	l := service.L().C(ctx).Mth("put-cmd-rs").Dbg()

	rqJs, err := json.Marshal(rq)
	if err != nil {
		return err
	}

	_, err = s.POST(ctx, fmt.Sprintf("%s/backend/commands/%s/response", s.baseUrl, cmdId), rqJs)
	if err != nil {
		return err
	}
	l.Dbg("ok")
	return nil
}

func (s *Sdk) GetCommand(ctx context.Context, cmdId string) (*backend.Command, error) {
	service.L().C(ctx).Mth("get-cmd").Dbg()

	rs, err := s.GET(ctx, fmt.Sprintf("%s/backend/commands/%s", s.baseUrl, cmdId))
	if err != nil {
		return nil, err
	}

	var p *backend.Command
	err = json.Unmarshal(rs, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Sdk) SearchCommands(ctx context.Context, params map[string]interface{}) (*backend.CommandSearchResponse, error) {
	service.L().C(ctx).Mth("search-cmd").Dbg()

	rs, err := s.GET(ctx, fmt.Sprintf("%s/backend/commands/search/query%s", s.baseUrl, s.toUrlParams(params)))
	if err != nil {
		return nil, err
	}

	var p *backend.CommandSearchResponse
	err = json.Unmarshal(rs, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
