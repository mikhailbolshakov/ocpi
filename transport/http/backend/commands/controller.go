package commands

import (
	kitHttp "github.com/mikhailbolshakov/kit/http"
	service "github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/usecase"
	"net/http"
)

type Controller interface {
	kitHttp.Controller
	StartSession(http.ResponseWriter, *http.Request)
	StopSession(http.ResponseWriter, *http.Request)
	Reservation(http.ResponseWriter, *http.Request)
	CancelReservation(http.ResponseWriter, *http.Request)
	UnLockConnector(http.ResponseWriter, *http.Request)
	PutCommandResponse(http.ResponseWriter, *http.Request)
	GetCommand(http.ResponseWriter, *http.Request)
	SearchCommands(w http.ResponseWriter, r *http.Request)
}

type ctrlImpl struct {
	kitHttp.BaseController
	cmdUc         usecase.CommandUc
	converter     usecase.CommandConverter
	localPlatform domain.LocalPlatformService
	cmdService    domain.CommandService
}

func NewController(cmdUc usecase.CommandUc, converter usecase.CommandConverter, localPlatform domain.LocalPlatformService,
	cmdService domain.CommandService) Controller {
	return &ctrlImpl{
		BaseController: kitHttp.BaseController{Logger: service.LF()},
		cmdUc:          cmdUc,
		converter:      converter,
		localPlatform:  localPlatform,
		cmdService:     cmdService,
	}
}

// StartSession godoc
// @Summary sends "START" command to the remote platform
// @Accept json
// @Param request body backend.StartSessionRequest true "start command request"
// @Success 200
// @Failure 500 {object} http.Error
// @Router /backend/commands/sessions/start [post]
// @tags commands
func (c *ctrlImpl) StartSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[backend.StartSessionRequest](ctx, r)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	err = c.cmdUc.OnLocalStartSession(ctx, c.converter.StartSessionCommandBackendToDomain(rq, c.localPlatform.GetPlatformId(ctx)))
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, kitHttp.EmptyOkResponse)
}

// StopSession godoc
// @Summary sends "STOP" command to the remote platform
// @Accept json
// @Param request body backend.StopSessionRequest true "stop command request"
// @Success 200
// @Failure 500 {object} http.Error
// @Router /backend/commands/sessions/stop [post]
// @tags commands
func (c *ctrlImpl) StopSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[backend.StopSessionRequest](ctx, r)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	err = c.cmdUc.OnLocalStopSession(ctx, c.converter.StopSessionCommandBackendToDomain(rq, c.localPlatform.GetPlatformId(ctx)))
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, kitHttp.EmptyOkResponse)
}

// Reservation godoc
// @Summary sends "RESERVE_NOW" command to the remote platform
// @Accept json
// @Param request body backend.ReserveNowRequest true "reservation request"
// @Success 200
// @Failure 500 {object} http.Error
// @Router /backend/commands/reservations [post]
// @tags commands
func (c *ctrlImpl) Reservation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[backend.ReserveNowRequest](ctx, r)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	err = c.cmdUc.OnLocalReserve(ctx, c.converter.ReserveNowCommandBackendToDomain(rq, c.localPlatform.GetPlatformId(ctx)))
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, kitHttp.EmptyOkResponse)
}

// CancelReservation godoc
// @Summary sends "RESERVATION_CANCELLATION" command to the remote platform
// @Accept json
// @Param resId path string true "reservation ID to cancel"
// @Param request body backend.CancelReservationRequest true "reservation cancellation request"
// @Success 200
// @Failure 500 {object} http.Error
// @Router /backend/commands/reservations [delete]
// @tags commands
func (c *ctrlImpl) CancelReservation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[backend.CancelReservationRequest](ctx, r)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	err = c.cmdUc.OnLocalCancelReservation(ctx, c.converter.CancelReservationCommandBackendToDomain(rq, c.localPlatform.GetPlatformId(ctx)))
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, kitHttp.EmptyOkResponse)
}

func (c *ctrlImpl) UnLockConnector(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}

// PutCommandResponse godoc
// @Summary sends a response for the incoming command from the remote platform
// @Accept json
// @Param cmdId path string true "OCPI command ID"
// @Param request body backend.CommandResponse true "stop command request"
// @Success 200
// @Failure 500 {object} http.Error
// @Router /backend/commands/{cmdId}/response [post]
// @tags commands
func (c *ctrlImpl) PutCommandResponse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cmdId, err := c.Var(ctx, r, "cmdId", true)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	rq, err := kitHttp.DecodeRequest[backend.CommandResponse](ctx, r)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	err = c.cmdUc.OnLocalCommandSetResponse(ctx, cmdId, rq.Status, rq.ErrMsg)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, kitHttp.EmptyOkResponse)
}

// GetCommand godoc
// @Summary retrieves a command object by id
// @Accept json
// @Param cmdId path string true "OCPI command ID"
// @Success 200 {object} backend.Command
// @Failure 500 {object} http.Error
// @Router /backend/commands/{cmdId} [get]
// @tags commands
func (c *ctrlImpl) GetCommand(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cmdId, err := c.Var(ctx, r, "cmdId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	cmd, err := c.cmdService.Get(ctx, cmdId)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.converter.CommandDomainToBackend(cmd))
}

// SearchCommands godoc
// @Summary retrieves commands objects by criteria
// @Accept json
// @Param offset query string false "number of items to offset from the beginning"
// @Param limit query string false "number of items to retrieve"
// @Param dateFrom query string false "items updated after the given date"
// @Param dateTo query string false "items updated before the given date"
// @Param refId query string false "reference id"
// @Param authRef query string false "auth reference"
// @Param cmd query string false "cmd type"
// @Success 200 {object} backend.CommandSearchResponse
// @Failure 500 {object} http.Error
// @Router /backend/commands/search/query [get]
// @tags commands
func (c *ctrlImpl) SearchCommands(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error
	cr := &domain.CommandSearchCriteria{}

	cr.Offset, err = c.FormValInt(ctx, r, "offset", true)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	cr.Limit, err = c.FormValInt(ctx, r, "limit", true)
	if err != nil {
		c.RespondError(w, err)
		return
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

	cr.RefId, err = c.FormVal(ctx, r, "refId", true)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	cr.AuthRef, err = c.FormVal(ctx, r, "authRef", true)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	cr.Cmd, err = c.FormVal(ctx, r, "cmd", true)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	cr.DateFrom, err = c.FormValTime(ctx, r, "updatedFrom", true)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	partyId, err := c.FormVal(ctx, r, "partyId", true)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	countryCode, err := c.FormVal(ctx, r, "countryCode", true)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	if partyId != "" && countryCode != "" {
		cr.ExtId = &domain.PartyExtId{
			PartyId:     partyId,
			CountryCode: countryCode,
		}
	}

	rs, err := c.cmdService.SearchCommands(ctx, cr)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, &backend.CommandSearchResponse{
		PageInfo: &backend.PageResponse{
			Total: rs.Total,
			Limit: rs.Limit,
		},
		Items: c.converter.CommandsDomainToBackend(rs.Items),
	})
}
