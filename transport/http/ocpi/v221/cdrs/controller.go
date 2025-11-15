package cdrs

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
	SenderGetCdrs(http.ResponseWriter, *http.Request)
	ReceiverGetCdr(http.ResponseWriter, *http.Request)
	ReceiverPostCdr(http.ResponseWriter, *http.Request)
}

type ctrlImpl struct {
	ocpi.Controller
	cdrService      domain.CdrService
	localPlatform   domain.LocalPlatformService
	converter       usecase.CdrConverter
	cdrUc           usecase.CdrUc
	senderSearchUrl string
}

func NewController(cdrService domain.CdrService, localPlatform domain.LocalPlatformService, converter usecase.CdrConverter,
	cdrUc usecase.CdrUc, cfg *cfg.CfgOcpiConfig) Controller {
	return &ctrlImpl{
		Controller:      ocpi.NewController(),
		cdrService:      cdrService,
		localPlatform:   localPlatform,
		converter:       converter,
		cdrUc:           cdrUc,
		senderSearchUrl: fmt.Sprintf("%s/%s", cfg.Local.Url, "ocpi/2.2.1/sender/cdrs"),
	}
}

func (c *ctrlImpl) SenderGetCdrs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve cdrs by the local platform
	rq := &domain.CdrSearchCriteria{
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

	rs, err := c.cdrService.SearchCdrs(ctx, rq)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	r = r.WithContext(c.SetResponseWithNextPageCtx(ctx, rs.PageResponse, c.senderSearchUrl))

	c.OcpiRespondOK(r, w, c.converter.CdrsDomainToModel(rs.Items))
}

func (c *ctrlImpl) ReceiverGetCdr(w http.ResponseWriter, r *http.Request) {
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
	sessId, err := c.Var(ctx, r, "cdr_id", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	sess, err := c.cdrService.GetCdr(ctx, sessId)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	if sess == nil {
		c.OcpiRespondNotFoundError(r, w)
		return
	}

	c.OcpiRespondOK(r, w, c.converter.CdrDomainToModel(sess))
}

func (c *ctrlImpl) ReceiverPostCdr(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[model.OcpiCdr](ctx, r)
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
	rq.Id, err = c.Var(ctx, r, "cdr_id", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	platformId, err := c.PlatformId(ctx)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	err = c.cdrUc.OnRemoteCdrPut(ctx, platformId, rq)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	c.OcpiRespondOK(r, w, nil)
}
