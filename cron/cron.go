package cron

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/cron"
	"github.com/mikhailbolshakov/kit/goroutine"
	service "github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/usecase"
	"time"
)

type cronImpl struct {
	cronManager cron.Manager
	commandUc   usecase.CommandUc
}

func NewCron(cronManager cron.Manager, commandUc usecase.CommandUc) cron.CronHandler {
	return &cronImpl{
		cronManager: cronManager,
		commandUc:   commandUc,
	}
}

func (c *cronImpl) l() kit.CLogger {
	return service.L().Cmp("cron-handler")
}

func (c *cronImpl) Register(ctx context.Context) {
	c.cronManager.Add(ctx, "local-cmd-deadline").
		Every(time.Minute).
		Action(c.localCmdDeadlineAsync())
	c.cronManager.Add(ctx, "remote-cmd-deadline").
		Every(time.Minute).
		Action(c.remoteCmdDeadlineAsync())
}

func (c *cronImpl) localCmdDeadlineAsync() cron.Action {
	return func(ctxFn func() context.Context) {
		ctx := ctxFn()
		goroutine.New().
			WithLogger(c.l().C(ctx).Mth("local-cmd-deadline")).
			Go(ctx, func() {
				c.commandUc.LocalCommandsDeadlineCronHandler(ctx)
			})
	}
}

func (c *cronImpl) remoteCmdDeadlineAsync() cron.Action {
	return func(ctxFn func() context.Context) {
		ctx := ctxFn()
		goroutine.New().
			WithLogger(c.l().C(ctx).Mth("remote-cmd-deadline")).
			Go(ctx, func() {
				c.commandUc.RemoteCommandsDeadlineCronHandler(ctx)
			})
	}
}
