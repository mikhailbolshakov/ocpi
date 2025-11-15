package commands

import (
	"context"
	kitHttp "github.com/mikhailbolshakov/kit/http"
	"github.com/mikhailbolshakov/ocpi/errors"
	"github.com/mikhailbolshakov/ocpi/model"
	"github.com/mikhailbolshakov/ocpi/transport/http/ocpi"
	"github.com/mikhailbolshakov/ocpi/usecase"
	"net/http"
	"sync"
)

type commandDetails struct {
	RqFunc  func() any
	CmdFunc func(context.Context, string, any) (*model.OcpiCommandResponse, error)
}

type Controller interface {
	kitHttp.Controller
	SenderSetCommandResponse(http.ResponseWriter, *http.Request)
	ReceiverExecCommand(http.ResponseWriter, *http.Request)
}

type ctrlImpl struct {
	ocpi.Controller
	once       sync.Once
	commandMap map[string]commandDetails
	commandUc  usecase.CommandUc
}

func NewController(commandUc usecase.CommandUc) Controller {
	return &ctrlImpl{
		Controller: ocpi.NewController(),
		commandUc:  commandUc,
	}
}

func (c *ctrlImpl) getCommandsCfg() map[string]commandDetails {
	c.once.Do(
		func() {
			c.commandMap = map[string]commandDetails{
				model.CommandStartSession: {
					RqFunc: func() any { return &model.OcpiStartSession{} },
					CmdFunc: func(ctx context.Context, platformId string, rq any) (*model.OcpiCommandResponse, error) {
						return c.commandUc.OnRemoteStartSession(ctx, platformId, rq.(*model.OcpiStartSession))
					},
				},
				model.CommandStopSession: {
					RqFunc: func() any { return &model.OcpiStopSession{} },
					CmdFunc: func(ctx context.Context, platformId string, rq any) (*model.OcpiCommandResponse, error) {
						return c.commandUc.OnRemoteStopSession(ctx, platformId, rq.(*model.OcpiStopSession))
					},
				},
				model.CommandReserve: {
					RqFunc: func() any { return &model.OcpiReserveNow{} },
					CmdFunc: func(ctx context.Context, platformId string, rq any) (*model.OcpiCommandResponse, error) {
						return c.commandUc.OnRemoteReserve(ctx, platformId, rq.(*model.OcpiReserveNow))
					},
				},
				model.CommandCancelReservation: {
					RqFunc: func() any { return &model.OcpiCancelReservation{} },
					CmdFunc: func(ctx context.Context, platformId string, rq any) (*model.OcpiCommandResponse, error) {
						return c.commandUc.OnRemoteCancelReservation(ctx, platformId, rq.(*model.OcpiCancelReservation))
					},
				},
				model.CommandUnlockConnector: {
					RqFunc: func() any { return &model.OcpiUnlockConnector{} },
					CmdFunc: func(ctx context.Context, platformId string, rq any) (*model.OcpiCommandResponse, error) {
						return c.commandUc.OnRemoteUnlockConnector(ctx, platformId, rq.(*model.OcpiUnlockConnector))
					},
				},
			}
		})
	return c.commandMap
}

func (c *ctrlImpl) SenderSetCommandResponse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	command, err := c.Var(ctx, r, model.OcpiQueryParamCommand, false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	// check command
	if _, ok := c.getCommandsCfg()[command]; !ok {
		c.OcpiRespondError(r, w, errors.ErrCmdNotSupported(ctx))
		return
	}

	uid, err := c.Var(ctx, r, model.OcpiQueryParamUid, false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	rq, err := kitHttp.DecodeRequest[model.OcpiCommandResult](ctx, r)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	platformId, err := c.PlatformId(ctx)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	err = c.commandUc.OnRemoteSetResponse(ctx, platformId, uid, rq)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	c.OcpiRespondOK(r, w, nil)

}

func (c *ctrlImpl) ReceiverExecCommand(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := c.EnsureCtxHeaders(ctx, model.OcpiCtxFromParty, model.OcpiCtxFromCountryCode); err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	command, err := c.Var(ctx, r, model.OcpiQueryParamCommand, false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	platformId, err := c.PlatformId(ctx)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	// check command
	cmdDet, ok := c.getCommandsCfg()[command]
	if !ok {
		c.OcpiRespondError(r, w, errors.ErrCmdNotSupported(ctx))
		return
	}
	rq := cmdDet.RqFunc()
	if err := c.DecodeRequest(ctx, r, rq); err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	rs, err := cmdDet.CmdFunc(ctx, platformId, rq)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	c.OcpiRespondOK(r, w, rs)
}
