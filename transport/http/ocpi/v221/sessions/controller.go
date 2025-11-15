package sessions

import (
	"fmt"
	kitHttp "github.com/mikhailbolshakov/kit/http"
	cfg "github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/model"
	"github.com/mikhailbolshakov/ocpi/transport/http/ocpi"
	"github.com/mikhailbolshakov/ocpi/usecase"
	"net/http"
)

type Controller interface {
	kitHttp.Controller
	SenderGetSessions(http.ResponseWriter, *http.Request)
	ReceiverGetSession(http.ResponseWriter, *http.Request)
	ReceiverPostSession(http.ResponseWriter, *http.Request)
	ReceiverPatchSession(http.ResponseWriter, *http.Request)
	SenderPutChargingPreferences(http.ResponseWriter, *http.Request)
}

type ctrlImpl struct {
	ocpi.Controller
	sessionService  domain.SessionService
	localPlatform   domain.LocalPlatformService
	converter       usecase.SessionConverter
	sessionUc       usecase.SessionUc
	senderSearchUrl string
}

func NewController(sessionService domain.SessionService, localPlatform domain.LocalPlatformService, converter usecase.SessionConverter,
	sessionUc usecase.SessionUc, cfg *cfg.CfgOcpiConfig) Controller {
	return &ctrlImpl{
		Controller:      ocpi.NewController(),
		sessionService:  sessionService,
		localPlatform:   localPlatform,
		converter:       converter,
		sessionUc:       sessionUc,
		senderSearchUrl: fmt.Sprintf("%s/%s", cfg.Local.Url, "ocpi/2.2.1/sender/sessions"),
	}
}

func (c *ctrlImpl) SenderGetSessions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve sessions by the local platform
	rq := &domain.SessionSearchCriteria{
		IncPlatforms: []string{c.localPlatform.GetPlatformId(ctx)},
	}

	var err error
	rq.DateFrom, err = c.FormValTime(ctx, r, model.OcpiQueryParamDateFrom, true)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	rq.DateTo, err = c.FormValTime(ctx, r, model.OcpiQueryParamDateTo, true)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	rq.Offset, err = c.FormValInt(ctx, r, model.OcpiQueryParamOffset, true)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	rq.Limit, err = c.FormValInt(ctx, r, model.OcpiQueryParamLimit, true)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	rs, err := c.sessionService.SearchSessions(ctx, rq)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	r = r.WithContext(c.SetResponseWithNextPageCtx(ctx, rs.PageResponse, c.senderSearchUrl))

	c.OcpiRespondOK(r, w, c.converter.SessionsDomainToModel(rs.Items))
}

func (c *ctrlImpl) ReceiverGetSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, err := c.Var(ctx, r, model.OcpiQueryParamPartyId, false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	_, err = c.Var(ctx, r, model.OcpiQueryParamCountryCode, false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	sessId, err := c.Var(ctx, r, "session_id", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	sess, err := c.sessionService.GetSessionWithPeriods(ctx, sessId)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	if sess == nil {
		c.OcpiRespondNotFoundError(r, w)
		return
	}

	c.OcpiRespondOK(r, w, c.converter.SessionDomainToModel(sess))
}

func (c *ctrlImpl) ReceiverPostSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[model.OcpiSession](ctx, r)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	rq.PartyId, err = c.Var(ctx, r, model.OcpiQueryParamPartyId, false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	rq.CountryCode, err = c.Var(ctx, r, model.OcpiQueryParamCountryCode, false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	rq.Id, err = c.Var(ctx, r, "session_id", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	platformId, err := c.PlatformId(ctx)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	err = c.sessionUc.OnRemoteSessionPut(ctx, platformId, rq)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	c.OcpiRespondOK(r, w, nil)
}

func (c *ctrlImpl) ReceiverPatchSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[model.OcpiSession](ctx, r)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	rq.PartyId, err = c.Var(ctx, r, model.OcpiQueryParamPartyId, false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	rq.CountryCode, err = c.Var(ctx, r, model.OcpiQueryParamCountryCode, false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	rq.Id, err = c.Var(ctx, r, "session_id", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	platformId, err := c.PlatformId(ctx)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	err = c.sessionUc.OnRemoteSessionPatch(ctx, platformId, rq)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	c.OcpiRespondOK(r, w, nil)
}

func (c *ctrlImpl) SenderPutChargingPreferences(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}
