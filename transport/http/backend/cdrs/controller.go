package cdrs

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
	PostCdr(http.ResponseWriter, *http.Request)
	GetCdr(http.ResponseWriter, *http.Request)
	SearchCdrs(http.ResponseWriter, *http.Request)
}

type ctrlImpl struct {
	kitHttp.BaseController
	cdrUc         usecase.CdrUc
	converter     usecase.CdrConverter
	localPlatform domain.LocalPlatformService
	cdrService    domain.CdrService
}

func NewController(cdrUc usecase.CdrUc, converter usecase.CdrConverter, localPlatform domain.LocalPlatformService,
	cdrService domain.CdrService) Controller {
	return &ctrlImpl{
		BaseController: kitHttp.BaseController{Logger: service.LF()},
		cdrUc:          cdrUc,
		converter:      converter,
		localPlatform:  localPlatform,
		cdrService:     cdrService,
	}
}

// PostCdr godoc
// @Summary creates cdr object in OCPI
// @Accept json
// @Param request body backend.Cdr true "cdr object"
// @Success 200
// @Failure 500 {object} http.Error
// @Router /backend/cdrs [post]
// @tags cdrs
func (c *ctrlImpl) PostCdr(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[backend.Cdr](ctx, r)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	rq.PlatformId = c.localPlatform.GetPlatformId(ctx)
	err = c.cdrUc.OnLocalCdrChanged(ctx, rq)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, kitHttp.EmptyOkResponse)
}

// GetCdr godoc
// @Summary retrieves a cdr object by id
// @Accept json
// @Param sessId path string true "OCPI cdr ID"
// @Success 200 {object} backend.Cdr
// @Failure 500 {object} http.Error
// @Router /backend/cdrs/{cdrId} [get]
// @tags cdrs
func (c *ctrlImpl) GetCdr(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cdrId, err := c.Var(ctx, r, "cdrId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	cdr, err := c.cdrService.GetCdr(ctx, cdrId)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.converter.CdrDomainToBackend(cdr))
}

// SearchCdrs godoc
// @Summary retrieves cdr objects by criteria
// @Accept json
// @Param offset query string false "number of items to offset from the beginning"
// @Param limit query string false "number of items to retrieve"
// @Param dateFrom query string false "items updated after the given date"
// @Param dateTo query string false "items updated before the given date"
// @Param refId query string false "reference id"
// @Param partyId query string false "OCPI party id"
// @Param countryCode query string false "OCPI country code"
// @Param incPlatforms query string false "comma separated platforms to include"
// @Param excPlatforms query string false "comma separated platforms to exclude"
// @Param ids query string false "comma separated list of ids"
// @Success 200 {object} backend.CdrSearchResponse
// @Failure 500 {object} http.Error
// @Router /backend/cdrs/search/query [get]
// @tags cdrs
func (c *ctrlImpl) SearchCdrs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error
	cr := &domain.CdrSearchCriteria{}

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

	rs, err := c.cdrService.SearchCdrs(ctx, cr)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, &backend.CdrSearchResponse{
		PageInfo: &backend.PageResponse{
			Total: rs.Total,
			Limit: rs.Limit,
		},
		Items: c.converter.CdrsDomainToBackend(rs.Items),
	})
}
