package tariffs

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
	SenderGetTariffs(http.ResponseWriter, *http.Request)
	ReceiverGetTariff(http.ResponseWriter, *http.Request)
	ReceiverPutTariff(http.ResponseWriter, *http.Request)
}

type ctrlImpl struct {
	ocpi.Controller
	tariffService   domain.TariffService
	localPlatform   domain.LocalPlatformService
	converter       usecase.TariffConverter
	tariffUc        usecase.TariffUc
	senderSearchUrl string
}

func NewController(tariffService domain.TariffService, localPlatform domain.LocalPlatformService, converter usecase.TariffConverter,
	tariffUc usecase.TariffUc, cfg *cfg.CfgOcpiConfig) Controller {
	return &ctrlImpl{
		Controller:      ocpi.NewController(),
		tariffService:   tariffService,
		localPlatform:   localPlatform,
		converter:       converter,
		tariffUc:        tariffUc,
		senderSearchUrl: fmt.Sprintf("%s/%s", cfg.Local.Url, "ocpi/2.2.1/sender/tariffs"),
	}
}

func (c *ctrlImpl) SenderGetTariffs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve tariffs by the local platform
	rq := &domain.TariffSearchCriteria{
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

	rs, err := c.tariffService.SearchTariffs(ctx, rq)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	r = r.WithContext(c.SetResponseWithNextPageCtx(ctx, rs.PageResponse, c.senderSearchUrl))

	c.OcpiRespondOK(r, w, c.converter.TariffsDomainToModel(rs.Items))
}

func (c *ctrlImpl) ReceiverGetTariff(w http.ResponseWriter, r *http.Request) {
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
	trfId, err := c.Var(ctx, r, "tariff_id", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	trf, err := c.tariffService.GetTariff(ctx, trfId)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	if trf == nil {
		c.OcpiRespondNotFoundError(r, w)
		return
	}

	c.OcpiRespondOK(r, w, c.converter.TariffDomainToModel(trf))
}

func (c *ctrlImpl) ReceiverPutTariff(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[model.OcpiTariff](ctx, r)
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
	rq.Id, err = c.Var(ctx, r, "tariff_id", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	platformId, err := c.PlatformId(ctx)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	err = c.tariffUc.OnRemoteTariffPut(ctx, platformId, rq)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	c.OcpiRespondOK(r, w, nil)
}
