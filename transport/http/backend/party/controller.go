package party

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
	PutParty(http.ResponseWriter, *http.Request)
	GetParty(http.ResponseWriter, *http.Request)
	SearchParties(http.ResponseWriter, *http.Request)
	PullParties(http.ResponseWriter, *http.Request)
}

type ctrlImpl struct {
	kitHttp.BaseController
	partyService         domain.PartyService
	credentialUc         usecase.CredentialsUc
	credentialsConverter usecase.CredentialsConverter
	localPlatform        domain.LocalPlatformService
	hubUc                usecase.HubUc
}

func NewController(credentialUc usecase.CredentialsUc, hubUc usecase.HubUc, credentialsConverter usecase.CredentialsConverter,
	localPlatform domain.LocalPlatformService, partyService domain.PartyService) Controller {
	return &ctrlImpl{
		BaseController:       kitHttp.BaseController{Logger: service.LF()},
		credentialUc:         credentialUc,
		hubUc:                hubUc,
		credentialsConverter: credentialsConverter,
		localPlatform:        localPlatform,
		partyService:         partyService,
	}
}

// PutParty godoc
// @Summary updates party object in OCPI
// @Accept json
// @Param request body backend.Party true "party object"
// @Success 200
// @Failure 500 {object} http.Error
// @Router /backend/parties [post]
// @tags parties
func (c *ctrlImpl) PutParty(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[backend.Party](ctx, r)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	err = c.credentialUc.OnLocalPartyChanged(ctx, c.credentialsConverter.PartyBackendToDomain(c.localPlatform.GetPlatformId(ctx), rq))
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, kitHttp.EmptyOkResponse)
}

// GetParty godoc
// @Summary retrieves a party object by id
// @Accept json
// @Param partyId path string true "party ID"
// @Success 200 {object} backend.Party
// @Failure 500 {object} http.Error
// @Router /backend/parties/{partyId} [get]
// @tags parties
func (c *ctrlImpl) GetParty(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	partyId, err := c.Var(ctx, r, "partyId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	party, err := c.partyService.Get(ctx, partyId)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.credentialsConverter.PartyDomainToBackend(party))
}

// SearchParties godoc
// @Summary retrieves parties objects by criteria
// @Accept json
// @Param offset query string false "number of items to offset from the beginning"
// @Param limit query string false "number of items to retrieve"
// @Param refId query string false "party reference id"
// @Param dateFrom query string false "items updated after the given date"
// @Param dateTo query string false "items updated before the given date"
// @Param partyId query string false "OCPI party id"
// @Param countryCode query string false "OCPI country code"
// @Param incPlatforms query string false "comma separated platforms to include"
// @Param excPlatforms query string false "comma separated platforms to exclude"
// @Param ids query string false "comma separated list of ids"
// @Success 200 {object} backend.PartySearchResponse
// @Failure 500 {object} http.Error
// @Router /backend/parties/search/query [get]
// @tags parties
func (c *ctrlImpl) SearchParties(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error
	cr := &domain.PartySearchCriteria{}

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

	cr.RefId, err = c.FormVal(ctx, r, "refId", true)
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

	rs, err := c.partyService.Search(ctx, cr)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, &backend.PartySearchResponse{
		PageInfo: &backend.PageResponse{
			Total: rs.Total,
			Limit: rs.Limit,
		},
		Items: c.credentialsConverter.PartiesDomainToBackend(rs.Items),
	})

}

// PullParties godoc
// @Summary triggers pulling parties from the remote platforms
// @Accept json
// @Success 200
// @Failure 500 {object} http.Error
// @Router /backend/parties/pull [post]
// @tags parties
func (c *ctrlImpl) PullParties(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := c.credentialUc.OnRemotePartyPull(ctx)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, kitHttp.EmptyOkResponse)
}
