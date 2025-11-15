package sessions

import (
	"github.com/mikhailbolshakov/kit"
	kitHttp "github.com/mikhailbolshakov/kit/http"
	service "github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/usecase"
	"net/http"
	"strings"
)

type Controller interface {
	kitHttp.Controller
	PutSession(http.ResponseWriter, *http.Request)
	PatchSession(http.ResponseWriter, *http.Request)
	GetSession(http.ResponseWriter, *http.Request)
	SearchSessions(http.ResponseWriter, *http.Request)
}

type ctrlImpl struct {
	kitHttp.BaseController
	sessUc        usecase.SessionUc
	converter     usecase.SessionConverter
	localPlatform domain.LocalPlatformService
	sessService   domain.SessionService
}

func NewController(sessUc usecase.SessionUc, converter usecase.SessionConverter, localPlatform domain.LocalPlatformService,
	sessService domain.SessionService) Controller {
	return &ctrlImpl{
		BaseController: kitHttp.BaseController{Logger: service.LF()},
		sessUc:         sessUc,
		converter:      converter,
		localPlatform:  localPlatform,
		sessService:    sessService,
	}
}

// PutSession godoc
// @Summary updates session object in OCPI
// @Accept json
// @Param request body backend.Session true "session object"
// @Success 200
// @Failure 500 {object} http.Error
// @Router /backend/sessions [post]
// @tags sessions
func (c *ctrlImpl) PutSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[backend.Session](ctx, r)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	err = c.sessUc.OnLocalSessionChanged(ctx, c.converter.SessionBackendToDomain(rq, c.localPlatform.GetPlatformId(ctx)))
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, kitHttp.EmptyOkResponse)
}

// PatchSession godoc
// @Summary patches (partly updates) session object in OCPI
// @Accept json
// @Param request body backend.Session true "session object"
// @Success 200
// @Failure 500 {object} http.Error
// @Router /backend/sessions [patch]
// @tags sessions
func (c *ctrlImpl) PatchSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[backend.Session](ctx, r)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	err = c.sessUc.OnLocalSessionPatched(ctx, c.converter.SessionBackendToDomain(rq, c.localPlatform.GetPlatformId(ctx)))
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, kitHttp.EmptyOkResponse)
}

// GetSession godoc
// @Summary retrieves a session object by id
// @Accept json
// @Param sessId path string true "OCPI session ID"
// @Param withChargingPeriods query bool false "if true charging periods are retrieved"
// @Success 200 {object} backend.Session
// @Failure 500 {object} http.Error
// @Router /backend/sessions/{sessId} [get]
// @tags sessions
func (c *ctrlImpl) GetSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessId, err := c.Var(ctx, r, "sessId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	withChargingPeriods, err := c.FormValBool(ctx, r, "withChargingPeriods", true)
	if err != nil {
		c.RespondError(w, err)
		return
	}
	if withChargingPeriods == nil {
		withChargingPeriods = kit.BoolPtr(false)
	}

	var sess *domain.Session
	if *withChargingPeriods {
		sess, err = c.sessService.GetSessionWithPeriods(ctx, sessId)
	} else {
		sess, err = c.sessService.GetSession(ctx, sessId)
	}
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.converter.SessionDomainToBackend(sess))
}

// SearchSessions godoc
// @Summary retrieves session objects by criteria
// @Accept json
// @Param offset query string false "number of items to offset from the beginning"
// @Param limit query string false "number of items to retrieve"
// @Param dateFrom query string false "items updated after the given date"
// @Param dateTo query string false "items updated before the given date"
// @Param refId query string false "reference id"
// @Param partyId query string false "OCPI party id"
// @Param countryCode query string false "OCPI country code"
// @Param withChargingPeriods query string false "if true charging periods are retrieved"
// @Param incPlatforms query string false "comma separated platforms to include"
// @Param excPlatforms query string false "comma separated platforms to exclude"
// @Param ids query string false "comma separated list of ids"
// @Param authRef query string false "auth reference"
// @Success 200 {object} backend.SessionSearchResponse
// @Failure 500 {object} http.Error
// @Router /backend/sessions/search/query [get]
// @tags sessions
func (c *ctrlImpl) SearchSessions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error
	cr := &domain.SessionSearchCriteria{}

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

	incPlatforms, err := c.FormVal(ctx, r, "incPlatforms", true)
	if err != nil {
		c.RespondError(w, err)
		return
	}
	if incPlatforms != "" {
		cr.IncPlatforms = strings.Split(incPlatforms, ",")
	}

	excPlatforms, err := c.FormVal(ctx, r, "excPlatforms", true)
	if err != nil {
		c.RespondError(w, err)
		return
	}
	if excPlatforms != "" {
		cr.ExcPlatforms = strings.Split(excPlatforms, ",")
	}

	ids, err := c.FormVal(ctx, r, "ids", true)
	if err != nil {
		c.RespondError(w, err)
		return
	}
	if ids != "" {
		cr.Ids = strings.Split(ids, ",")
	}

	withChargingPeriods, err := c.FormValBool(ctx, r, "withChargingPeriods", true)
	if err != nil {
		c.RespondError(w, err)
		return
	}
	if withChargingPeriods != nil {
		cr.WithChargingPeriods = *withChargingPeriods
	}

	rs, err := c.sessService.SearchSessions(ctx, cr)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, &backend.SessionSearchResponse{
		PageInfo: &backend.PageResponse{
			Total: rs.Total,
			Limit: rs.Limit,
		},
		Items: c.converter.SessionsDomainToBackend(rs.Items),
	})
}
