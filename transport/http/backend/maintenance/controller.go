package maintenance

import (
	kitHttp "github.com/mikhailbolshakov/kit/http"
	service "github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/usecase"
	"net/http"
	"strings"
)

type Controller interface {
	kitHttp.Controller
	DeleteParty(http.ResponseWriter, *http.Request)
	SearchLogEntries(http.ResponseWriter, *http.Request)
}

type ctrlImpl struct {
	kitHttp.BaseController
	maintenanceUc usecase.MaintenanceUc
	logService    domain.OcpiLogService
}

func NewController(maintenanceUc usecase.MaintenanceUc, logService domain.OcpiLogService) Controller {
	return &ctrlImpl{
		BaseController: kitHttp.BaseController{Logger: service.LF()},
		maintenanceUc:  maintenanceUc,
		logService:     logService,
	}
}

// DeleteParty godoc
// @Summary delete party of the local platform by external party id (partyId + country code) and all the related objects
// @Accept json
// @Param partyId path string true "party ID"
// @Success 200
// @Failure 500 {object} http.Error
// @Router /maintenance/parties/{partyId}/{countryCode} [delete]
// @tags maintenance
func (c *ctrlImpl) DeleteParty(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	partyId, err := c.Var(ctx, r, "partyId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}
	country, err := c.Var(ctx, r, "countryCode", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	err = c.maintenanceUc.DeleteLocalPartyByExt(ctx, domain.PartyExtId{
		PartyId:     partyId,
		CountryCode: country,
	})
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, kitHttp.EmptyOkResponse)
}

// SearchLogEntries godoc
// @Summary allows retrieving log entries
// @Accept json
// @Param requestId path string false "request id"
// @Param fromPlatform path string false "source platform id"
// @Param toPlatform path string false "target platform id"
// @Param ocpiStatus path int false "ocpi status"
// @Param httpStatus path int false "http status"
// @Param incoming path bool false "if true, only incoming events"
// @Param error path bool false "if true, only events with errors"
// @Param events query string false "comma separated events"
// @Param dateFrom query string false "logs created after the given date"
// @Param dateTo query string false "logs created before the given date"
// @Param size query int false "size of items"
// @Param index query int false "index"
// @Param short path bool false "if true, short data is retrieved"
// @Success 200 array LogMessage
// @Failure 500 {object} http.Error
// @Router /maintenance/logs [get]
// @tags maintenance
func (c *ctrlImpl) SearchLogEntries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cr := &domain.SearchLogCriteria{}
	size, index, err := c.FormPaging(ctx, r, nil)
	if err != nil {
		c.RespondError(w, err)
		return
	}
	if size != nil {
		cr.PagingRequest.Size = *size
	}
	if index != nil {
		cr.PagingRequest.Index = *index
	}

	cr.DateFrom, err = c.FormValTime(ctx, r, "dateFrom", true)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	cr.DateTo, err = c.FormValTime(ctx, r, "dateTo", true)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	events, err := c.FormVal(ctx, r, "events", true)
	if err != nil {
		c.RespondError(w, err)
		return
	}
	if events != "" {
		cr.Events = strings.Split(events, ",")
	}

	cr.RequestId, err = c.FormVal(ctx, r, "requestId", true)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	cr.FromPlatform, err = c.FormVal(ctx, r, "fromPlatform", true)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	cr.ToPlatform, err = c.FormVal(ctx, r, "toPlatform", true)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	cr.OcpiStatus, err = c.FormValInt(ctx, r, "ocpiStatus", true)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	cr.HttpStatus, err = c.FormValInt(ctx, r, "httpStatus", true)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	cr.Incoming, err = c.FormValBool(ctx, r, "incoming", true)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	cr.Error, err = c.FormValBool(ctx, r, "error", true)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	short, err := c.FormValBool(ctx, r, "short", true)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	messages, err := c.logService.Search(ctx, cr)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.toLogMessages(messages, short))
}
