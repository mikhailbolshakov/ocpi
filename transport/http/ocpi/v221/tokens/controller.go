package tokens

import (
	"fmt"
	kitHttp "github.com/mikhailbolshakov/kit/http"
	cfg "github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/errors"
	"github.com/mikhailbolshakov/ocpi/model"
	"github.com/mikhailbolshakov/ocpi/transport/http/ocpi"
	"github.com/mikhailbolshakov/ocpi/usecase"
	"net/http"
)

type Controller interface {
	kitHttp.Controller
	SenderGetTokens(http.ResponseWriter, *http.Request)
	SenderAuthToken(http.ResponseWriter, *http.Request)

	ReceiverGetToken(http.ResponseWriter, *http.Request)
	ReceiverPutToken(http.ResponseWriter, *http.Request)
	ReceiverPatchToken(http.ResponseWriter, *http.Request)
}

type ctrlImpl struct {
	ocpi.Controller
	tokenService    domain.TokenService
	localPlatform   domain.LocalPlatformService
	converter       usecase.TokenConverter
	tokenUc         usecase.TokenUc
	senderSearchUrl string
}

func NewController(tokenService domain.TokenService, localPlatform domain.LocalPlatformService, converter usecase.TokenConverter,
	tokenUc usecase.TokenUc, cfg *cfg.CfgOcpiConfig) Controller {
	return &ctrlImpl{
		Controller:      ocpi.NewController(),
		tokenService:    tokenService,
		localPlatform:   localPlatform,
		converter:       converter,
		tokenUc:         tokenUc,
		senderSearchUrl: fmt.Sprintf("%s/%s", cfg.Local.Url, "ocpi/2.2.1/sender/tokens"),
	}
}

func (c *ctrlImpl) SenderGetTokens(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve tokens by the local platform
	rq := &domain.TokenSearchCriteria{
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

	rs, err := c.tokenService.SearchTokens(ctx, rq)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	r = r.WithContext(c.SetResponseWithNextPageCtx(ctx, rs.PageResponse, c.senderSearchUrl))

	c.OcpiRespondOK(r, w, c.converter.TokensDomainToModel(rs.Items))
}

func (c *ctrlImpl) SenderAuthToken(w http.ResponseWriter, r *http.Request) {
	c.OcpiRespondError(r, w, errors.ErrCmdNotSupported(r.Context()))
}

func (c *ctrlImpl) ReceiverGetToken(w http.ResponseWriter, r *http.Request) {
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
	tknId, err := c.Var(ctx, r, "token_id", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	tkn, err := c.tokenService.GetToken(ctx, tknId)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	if tkn == nil {
		c.OcpiRespondNotFoundError(r, w)
		return
	}

	c.OcpiRespondOK(r, w, c.converter.TokenDomainToModel(tkn))
}

func (c *ctrlImpl) ReceiverPutToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[model.OcpiToken](ctx, r)
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
	rq.Id, err = c.Var(ctx, r, "token_id", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	platformId, err := c.PlatformId(ctx)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	err = c.tokenUc.OnRemoteTokenPut(ctx, platformId, rq)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	c.OcpiRespondOK(r, w, nil)
}

func (c *ctrlImpl) ReceiverPatchToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[model.OcpiToken](ctx, r)
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
	rq.Id, err = c.Var(ctx, r, "token_id", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	platformId, err := c.PlatformId(ctx)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	err = c.tokenUc.OnRemoteTokenPatch(ctx, platformId, rq)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	c.OcpiRespondOK(r, w, nil)
}
