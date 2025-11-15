package webhook

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/cluster"
	service "github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/domain"
)

type Adapter interface {
	cluster.Adapter
	backend.WebhookRepository
}

type adapterImpl struct {
	webhookRestClient
	logService domain.OcpiLogService
	cfg        *service.CfgWebHook
}

func NewAdapter(logService domain.OcpiLogService) Adapter {
	return &adapterImpl{
		logService: logService,
	}
}

func (a *adapterImpl) l() kit.CLogger {
	return service.L().Cmp("wh-adapter")
}

func (a *adapterImpl) Init(ctx context.Context, config interface{}) error {
	a.l().Mth("init").Dbg()
	a.cfg = config.(*service.CfgWebHook)
	if !a.cfg.Mock {
		a.webhookRestClient = newWhRestClient()
	} else {
		a.webhookRestClient = newMockWhRestClient()
	}
	if err := a.webhookRestClient.Init(ctx, a.cfg); err != nil {
		return err
	}
	return nil
}

func (a *adapterImpl) Close(ctx context.Context) error {
	return a.webhookRestClient.Close(ctx)
}
