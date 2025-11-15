package tariffs

import (
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
	PutTariff(http.ResponseWriter, *http.Request)
	PullTariffs(http.ResponseWriter, *http.Request)
	GetTariff(http.ResponseWriter, *http.Request)
	SearchTariffs(http.ResponseWriter, *http.Request)
}

type ctrlImpl struct {
	kitHttp.BaseController
	trfUc         usecase.TariffUc
	converter     usecase.TariffConverter
	localPlatform domain.LocalPlatformService
	trfService    domain.TariffService
}

func NewController(trfUc usecase.TariffUc, converter usecase.TariffConverter, localPlatform domain.LocalPlatformService,
	trfService domain.TariffService) Controller {
	return &ctrlImpl{
		BaseController: kitHttp.BaseController{Logger: service.LF()},
		trfUc:          trfUc,
		converter:      converter,
		localPlatform:  localPlatform,
		trfService:     trfService,
	}
}

// PutTariff godoc
// @Summary updates tariff object in OCPI
// @Accept json
// @Param request body backend.Tariff true "tariff object"
// @Success 200
// @Failure 500 {object} http.Error
// @Router /backend/tariffs [post]
// @tags tariffs
func (c *ctrlImpl) PutTariff(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[backend.Tariff](ctx, r)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	rq.PlatformId = c.localPlatform.GetPlatformId(ctx)
	err = c.trfUc.OnLocalTariffChanged(ctx, c.converter.TariffBackendToDomain(rq))
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, kitHttp.EmptyOkResponse)
}

// PullTariffs godoc
// @Summary triggers pulling tariffs from the remote platforms
// @Accept json
// @Param request body backend.PullRequest true "request object"
// @Success 200
// @Failure 500 {object} http.Error
// @Router /backend/tariffs/pull [post]
// @tags tariffs
func (c *ctrlImpl) PullTariffs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[backend.PullRequest](ctx, r)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	err = c.trfUc.OnRemoteTariffsPull(ctx, rq.From, rq.To)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, kitHttp.EmptyOkResponse)
}

// GetTariff godoc
// @Summary retrieves a tariff object by id
// @Accept json
// @Param trfId path string true "OCPI tariffs ID"
// @Success 200 {object} backend.Tariff
// @Failure 500 {object} http.Error
// @Router /backend/tariffs/{trfId} [get]
// @tags tariffs
func (c *ctrlImpl) GetTariff(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	trfId, err := c.Var(ctx, r, "trfId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	trf, err := c.trfService.GetTariff(ctx, trfId)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.converter.TariffDomainToBackend(trf))
}

// SearchTariffs godoc
// @Summary retrieves tariff objects by criteria
// @Accept json
// @Param offset query string false "number of items to offset from the beginning"
// @Param limit query string false "number of items to retrieve"
// @Param dateFrom query string false "items updated after the given date"
// @Param dateTo query string false "items updated before the given date"
// @Param refId query string false "reference id"
// @Param incPlatforms query string false "comma separated platforms to include"
// @Param excPlatforms query string false "comma separated platforms to exclude"
// @Param ids query string false "comma separated list of ids"
// @Success 200 {object} backend.TariffSearchResponse
// @Failure 500 {object} http.Error
// @Router /backend/tariffs/search/query [get]
// @tags tariffs
func (c *ctrlImpl) SearchTariffs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error
	cr := &domain.TariffSearchCriteria{}

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

	rs, err := c.trfService.SearchTariffs(ctx, cr)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, &backend.TariffSearchResponse{
		PageInfo: &backend.PageResponse{
			Total: rs.Total,
			Limit: rs.Limit,
		},
		Items: c.converter.TariffsDomainToBackend(rs.Items),
	})
}
