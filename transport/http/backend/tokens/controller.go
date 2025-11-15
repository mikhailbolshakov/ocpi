package tokens

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
	PutToken(http.ResponseWriter, *http.Request)
	PullTokens(http.ResponseWriter, *http.Request)
	GetToken(http.ResponseWriter, *http.Request)
	SearchTokens(http.ResponseWriter, *http.Request)
}

type ctrlImpl struct {
	kitHttp.BaseController
	tknUc         usecase.TokenUc
	converter     usecase.TokenConverter
	localPlatform domain.LocalPlatformService
	trfService    domain.TokenService
}

func NewController(trfUc usecase.TokenUc, converter usecase.TokenConverter, localPlatform domain.LocalPlatformService,
	trfService domain.TokenService) Controller {
	return &ctrlImpl{
		BaseController: kitHttp.BaseController{Logger: service.LF()},
		tknUc:          trfUc,
		converter:      converter,
		localPlatform:  localPlatform,
		trfService:     trfService,
	}
}

// PutToken godoc
// @Summary updates token object in OCPI
// @Accept json
// @Param request body backend.Token true "token object"
// @Success 200
// @Failure 500 {object} http.Error
// @Router /backend/tokens [post]
// @tags tokens
func (c *ctrlImpl) PutToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[backend.Token](ctx, r)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	err = c.tknUc.OnLocalTokenChanged(ctx, c.converter.TokenBackendToDomain(rq, c.localPlatform.GetPlatformId(ctx)))
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, kitHttp.EmptyOkResponse)
}

// PullTokens godoc
// @Summary triggers pulling tokens from the remote platforms
// @Accept json
// @Param request body backend.PullRequest true "request object"
// @Success 200
// @Failure 500 {object} http.Error
// @Router /backend/tokens/pull [post]
// @tags tokens
func (c *ctrlImpl) PullTokens(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[backend.PullRequest](ctx, r)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	err = c.tknUc.OnRemoteTokensPull(ctx, rq.From, rq.To)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, kitHttp.EmptyOkResponse)
}

// GetToken godoc
// @Summary retrieves a token object by id
// @Accept json
// @Param trfId path string true "OCPI token ID"
// @Success 200 {object} backend.Tariff
// @Failure 500 {object} http.Error
// @Router /backend/tokens/{tknId} [get]
// @tags tokens
func (c *ctrlImpl) GetToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	trfId, err := c.Var(ctx, r, "tknId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	trf, err := c.trfService.GetToken(ctx, trfId)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.converter.TokenDomainToBackend(trf))
}

// SearchTokens godoc
// @Summary retrieves token objects by criteria
// @Accept json
// @Param offset query string false "number of items to offset from the beginning"
// @Param limit query string false "number of items to retrieve"
// @Param dateFrom query string false "items updated after the given date"
// @Param dateTo query string false "items updated before the given date"
// @Param refId query string false "reference id"
// @Param incPlatforms query string false "comma separated platforms to include"
// @Param excPlatforms query string false "comma separated platforms to exclude"
// @Param ids query string false "comma separated list of ids"
// @Success 200 {object} backend.TokenSearchResponse
// @Failure 500 {object} http.Error
// @Router /backend/tokens/search/query [get]
// @tags tokens
func (c *ctrlImpl) SearchTokens(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error
	cr := &domain.TokenSearchCriteria{}

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

	rs, err := c.trfService.SearchTokens(ctx, cr)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, &backend.TokenSearchResponse{
		PageInfo: &backend.PageResponse{
			Total: rs.Total,
			Limit: rs.Limit,
		},
		Items: c.converter.TokensDomainToBackend(rs.Items),
	})
}
