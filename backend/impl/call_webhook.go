package impl

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
)

type webhookCall struct {
	webhook    backend.WebhookService
	repository backend.WebhookRepository
}

func NewWebhookCallService(webhook backend.WebhookService, repository backend.WebhookRepository) backend.WebhookCallService {
	return &webhookCall{
		webhook:    webhook,
		repository: repository,
	}
}

func (w *webhookCall) l() kit.CLogger {
	return ocpi.L().Cmp("wh-call-svc")
}

func (w *webhookCall) OnPartiesChanged(ctx context.Context, parties ...*backend.Party) error {
	w.l().C(ctx).Mth("on-parties").Dbg()
	return w.callAsync(ctx, backend.WhEventPartyChanged, parties)
}

func (w *webhookCall) OnLocationsChanged(ctx context.Context, locs ...*backend.Location) error {
	w.l().C(ctx).Mth("on-locations").Dbg()
	return w.callAsync(ctx, backend.WhEventLocationChanged, locs)
}

func (w *webhookCall) OnEvseChanged(ctx context.Context, evses ...*backend.Evse) error {
	w.l().C(ctx).Mth("on-evse").Dbg()
	return w.callAsync(ctx, backend.WhEventEvseChanged, evses)
}

func (w *webhookCall) OnConnectorChanged(ctx context.Context, cons ...*backend.Connector) error {
	w.l().C(ctx).Mth("on-con").Dbg()
	return w.callAsync(ctx, backend.WhEventConnectorChanged, cons)
}

func (w *webhookCall) OnTariffsChanged(ctx context.Context, tariffs ...*backend.Tariff) error {
	w.l().C(ctx).Mth("on-tariff").Dbg()
	return w.callAsync(ctx, backend.WhEventTariffChanged, tariffs)
}

func (w *webhookCall) OnTokensChanged(ctx context.Context, tokens ...*backend.Token) error {
	w.l().C(ctx).Mth("on-token").Dbg()
	return w.callAsync(ctx, backend.WhEventTokenChanged, tokens)
}

func (w *webhookCall) OnSessionsChanged(ctx context.Context, sessions ...*backend.Session) error {
	w.l().C(ctx).Mth("on-sess").Dbg()
	return w.callAsync(ctx, backend.WhEventSessionChanged, sessions)
}

func (w *webhookCall) OnCommandResponse(ctx context.Context, cmd *backend.Command) error {
	w.l().C(ctx).Mth("on-cmd-rs").Dbg()
	return w.callAsync(ctx, backend.WhEventCommandResponse, cmd)
}

func (w *webhookCall) OnStartSession(ctx context.Context, cmd *backend.Command) error {
	w.l().C(ctx).Mth("on-start-sess-cmd").Dbg()
	return w.callAsync(ctx, backend.WhEventStartSession, cmd)
}

func (w *webhookCall) OnStopSession(ctx context.Context, cmd *backend.Command) error {
	w.l().C(ctx).Mth("on-stop-sess-cmd").Dbg()
	return w.callAsync(ctx, backend.WhEventStopSession, cmd)
}

func (w *webhookCall) OnCdrChanged(ctx context.Context, cdr *backend.Cdr) error {
	w.l().C(ctx).Mth("on-cdr").Dbg()
	return w.callAsync(ctx, backend.WhEventCdrChanged, cdr)
}

func (w *webhookCall) OnReserveNow(ctx context.Context, cmd *backend.Command) error {
	w.l().C(ctx).Mth("on-res").Dbg()
	return w.callAsync(ctx, backend.WhEventReservation, cmd)
}

func (w *webhookCall) OnCancelReservation(ctx context.Context, cmd *backend.Command) error {
	w.l().C(ctx).Mth("on-cancel-res").Dbg()
	return w.callAsync(ctx, backend.WhEventCancelReservation, cmd)
}

func (w *webhookCall) callAsync(ctx context.Context, event string, payload any) error {
	// registered webhooks
	webhooks, err := w.getByEvent(ctx, event)
	if err != nil {
		return err
	}

	// call webhooks through repository
	for _, wh := range webhooks {
		w.repository.CallAsync(ctx, wh.Url, wh.ApiKey, event, payload)
	}
	return nil
}

func (w *webhookCall) getByEvent(ctx context.Context, event string) ([]*backend.Webhook, error) {
	return w.webhook.Search(ctx, &backend.SearchWebhookCriteria{Event: event})
}
