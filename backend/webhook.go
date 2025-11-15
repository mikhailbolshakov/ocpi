package backend

import (
	"context"
)

const (
	WhEventPartyChanged      = "party.changed"
	WhEventLocationChanged   = "location.changed"
	WhEventEvseChanged       = "evse.changed"
	WhEventConnectorChanged  = "connector.changed"
	WhEventTariffChanged     = "tariff.changed"
	WhEventTokenChanged      = "token.changed"
	WhEventSessionChanged    = "session.changed"
	WhEventCommandResponse   = "command.response"
	WhEventStartSession      = "command.start-session"
	WhEventStopSession       = "command.stop-session"
	WhEventCdrChanged        = "cdr.changed"
	WhEventReservation       = "command.reservation"
	WhEventCancelReservation = "command.reservation-cancel"
)

type Webhook struct {
	Id     string   // Id unique ID
	ApiKey string   // ApiKey used to call webhook
	Events []string // Events list of events to call webhook
	Url    string   // Url of webhook
}

type SearchWebhookCriteria struct {
	Event string
}

type WebhookService interface {
	// CreateUpdate registers a new webhook or update existent one
	CreateUpdate(ctx context.Context, wh *Webhook) (*Webhook, error)
	// Create registers a new webhook
	Create(ctx context.Context, wh *Webhook) (*Webhook, error)
	// Update updates an existent webhook
	Update(ctx context.Context, wh *Webhook) (*Webhook, error)
	// Delete deletes a webhook
	Delete(ctx context.Context, whId string) error
	// Search retrieves webhooks by criteria
	Search(ctx context.Context, cr *SearchWebhookCriteria) ([]*Webhook, error)
	// Get retrieves webhook by id
	Get(ctx context.Context, whId string) (*Webhook, error)
}

type WebhookCallService interface {
	// OnPartiesChanged makes a webhook call when parties changed
	OnPartiesChanged(ctx context.Context, parties ...*Party) error
	// OnLocationsChanged makes a webhook call when locations changed
	OnLocationsChanged(ctx context.Context, locs ...*Location) error
	// OnEvseChanged makes a webhook call when evses changed
	OnEvseChanged(ctx context.Context, evses ...*Evse) error
	// OnConnectorChanged makes a webhook call when connectors changed
	OnConnectorChanged(ctx context.Context, cons ...*Connector) error
	// OnTariffsChanged makes a webhook call when tariffs changed
	OnTariffsChanged(ctx context.Context, tariffs ...*Tariff) error
	// OnTokensChanged makes a webhook call when tariffs changed
	OnTokensChanged(ctx context.Context, tokens ...*Token) error
	// OnSessionsChanged makes a webhook call when sessions changed
	OnSessionsChanged(ctx context.Context, sessions ...*Session) error
	// OnCommandResponse makes a webhook call when a command response arrives
	OnCommandResponse(ctx context.Context, cmd *Command) error
	// OnStartSession makes a webhook call when a start session requested
	OnStartSession(ctx context.Context, cmd *Command) error
	// OnStopSession makes a webhook call when a stop session requested
	OnStopSession(ctx context.Context, cmd *Command) error
	// OnCdrChanged makes a webhook call when sessions changed
	OnCdrChanged(ctx context.Context, cdr *Cdr) error
	// OnReserveNow makes a webhook call when a reservation requested
	OnReserveNow(ctx context.Context, cmd *Command) error
	// OnCancelReservation makes a webhook call when a cancel reservation requested
	OnCancelReservation(ctx context.Context, cmd *Command) error
}

type WebhookRepository interface {
	// CallAsync executes webhook in async manner
	CallAsync(ctx context.Context, url, apiKey, event string, payload any)
}

type WebhookStorage interface {
	// CreateWebhook registers a new webhook
	CreateWebhook(ctx context.Context, wh *Webhook) error
	// UpdateWebhook updates an existent webhook
	UpdateWebhook(ctx context.Context, wh *Webhook) error
	// MergeWebhook merges webhook
	MergeWebhook(ctx context.Context, wh *Webhook) error
	// DeleteWebhook deletes a webhook
	DeleteWebhook(ctx context.Context, whId string) error
	// SearchWebhook searches webhooks
	SearchWebhook(ctx context.Context, cr *SearchWebhookCriteria) ([]*Webhook, error)
	// GetWebhook retrieves webhook by id
	GetWebhook(ctx context.Context, whId string) (*Webhook, error)
}
