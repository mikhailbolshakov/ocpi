package hub

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
	GetHubClientInfo(http.ResponseWriter, *http.Request)
	GetClientInfo(http.ResponseWriter, *http.Request)
	PutClientInfo(http.ResponseWriter, *http.Request)
}

type ctrlImpl struct {
	ocpi.Controller
	partyService    domain.PartyService
	localPlatform   domain.LocalPlatformService
	hubUc           usecase.HubUc
	senderSearchUrl string
}

func NewController(partyService domain.PartyService, localPlatform domain.LocalPlatformService, hubUc usecase.HubUc,
	cfg *cfg.CfgOcpiConfig) Controller {
	return &ctrlImpl{
		Controller:      ocpi.NewController(),
		partyService:    partyService,
		localPlatform:   localPlatform,
		hubUc:           hubUc,
		senderSearchUrl: fmt.Sprintf("%s/%s", cfg.Local.Url, "ocpi/2.2.1/sender/hubclientinfo"),
	}
}

func (c *ctrlImpl) GetHubClientInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq := &domain.PartySearchCriteria{
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

	rs, err := c.partyService.Search(ctx, rq)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	r = r.WithContext(c.SetResponseWithNextPageCtx(ctx, rs.PageResponse, c.senderSearchUrl))

	c.OcpiRespondOK(r, w, c.toClientInfosApi(rs))
}

func (c *ctrlImpl) GetClientInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error
	var extPartyId domain.PartyExtId
	extPartyId.PartyId, err = c.Var(ctx, r, model.OcpiQueryParamPartyId, false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	extPartyId.CountryCode, err = c.Var(ctx, r, model.OcpiQueryParamCountryCode, false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	rs, err := c.partyService.GetByExtId(ctx, extPartyId)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	if rs == nil {
		c.OcpiRespondNotFoundError(r, w)
		return
	}

	c.OcpiRespondOK(r, w, c.toClientInfoApi(rs))
}

func (c *ctrlImpl) PutClientInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[model.OcpiClientInfo](ctx, r)
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

	platformId, err := c.PlatformId(ctx)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	err = c.hubUc.OnRemoteClientInfoPut(ctx, platformId, rq)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	c.OcpiRespondOK(r, w, nil)
}
