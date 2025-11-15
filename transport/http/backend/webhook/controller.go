package webhook

import (
	kitHttp "github.com/mikhailbolshakov/kit/http"
	service "github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
	"net/http"
)

type Controller interface {
	kitHttp.Controller
	CreateUpdateWebhook(http.ResponseWriter, *http.Request)
	GetWebhook(http.ResponseWriter, *http.Request)
	DeleteWebhook(http.ResponseWriter, *http.Request)
	SearchWebhooks(http.ResponseWriter, *http.Request)
}

type ctrlImpl struct {
	kitHttp.BaseController
	whService backend.WebhookService
}

func NewController(whService backend.WebhookService) Controller {
	return &ctrlImpl{
		whService:      whService,
		BaseController: kitHttp.BaseController{Logger: service.LF()},
	}
}

// CreateUpdateWebhook godoc
// @Summary creates or updates webhook object
// @Accept json
// @Param request body Webhook true "webhook object"
// @Success 200 {object} Webhook
// @Failure 500 {object} http.Error
// @Router /webhooks [post]
// @tags webhooks
func (c *ctrlImpl) CreateUpdateWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[Webhook](ctx, r)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	wh, err := c.whService.CreateUpdate(ctx, c.toWebhookRequestDomain(rq))
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.toWebhookApi(wh))
}

// GetWebhook godoc
// @Summary retrieves a webhook object
// @Accept json
// @Param whId path string true "webhook id"
// @Success 200 {object} Webhook
// @Failure 500 {object} http.Error
// @Router /webhooks/{whId} [get]
// @tags webhooks
func (c *ctrlImpl) GetWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	whId, err := c.Var(ctx, r, "whId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	wh, err := c.whService.Get(ctx, whId)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.toWebhookApi(wh))
}

// DeleteWebhook godoc
// @Summary deleted a webhook object
// @Accept json
// @Param whId path string true "webhook id"
// @Success 200
// @Failure 500 {object} http.Error
// @Router /webhooks/{whId} [delete]
// @tags webhooks
func (c *ctrlImpl) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	whId, err := c.Var(ctx, r, "whId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	err = c.whService.Delete(ctx, whId)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, kitHttp.EmptyOkResponse)
}

// SearchWebhooks godoc
// @Summary retrieves webhook objects by criteria
// @Accept json
// @Param event query string false "event type"
// @Success 200 {array} Webhook
// @Failure 500 {object} http.Error
// @Router /webhooks/search/query [get]
// @tags webhooks
func (c *ctrlImpl) SearchWebhooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq := &backend.SearchWebhookCriteria{}
	var err error
	rq.Event, err = c.FormVal(ctx, r, "event", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	whs, err := c.whService.Search(ctx, rq)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.toWebhooksApi(whs))
}
