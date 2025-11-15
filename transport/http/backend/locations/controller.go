package locations

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
	PutLocation(http.ResponseWriter, *http.Request)
	PullLocations(http.ResponseWriter, *http.Request)
	GetLocation(http.ResponseWriter, *http.Request)
	SearchLocations(http.ResponseWriter, *http.Request)
	PutEvse(http.ResponseWriter, *http.Request)
	SetEvseStatus(http.ResponseWriter, *http.Request)
	GetEvse(http.ResponseWriter, *http.Request)
	SearchEvses(http.ResponseWriter, *http.Request)
	PutConnector(http.ResponseWriter, *http.Request)
	GetConnector(http.ResponseWriter, *http.Request)
	SearchConnectors(http.ResponseWriter, *http.Request)
}

type ctrlImpl struct {
	kitHttp.BaseController
	locationUc      usecase.LocationUc
	converter       usecase.LocationConverter
	localPlatform   domain.LocalPlatformService
	locationService domain.LocationService
}

func NewController(locationUc usecase.LocationUc, converter usecase.LocationConverter, localPlatform domain.LocalPlatformService,
	locationService domain.LocationService) Controller {
	return &ctrlImpl{
		BaseController:  kitHttp.BaseController{Logger: service.LF()},
		locationUc:      locationUc,
		converter:       converter,
		localPlatform:   localPlatform,
		locationService: locationService,
	}
}

// PutLocation godoc
// @Summary updates location object in OCPI. It may contain evse and connectors as well
// @Accept json
// @Param request body backend.Location true "location object"
// @Success 200
// @Failure 500 {object} http.Error
// @Router /backend/locations [post]
// @tags locations
func (c *ctrlImpl) PutLocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[backend.Location](ctx, r)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	err = c.locationUc.OnLocalLocationChanged(ctx, c.converter.LocationBackendToDomain(rq, c.localPlatform.GetPlatformId(ctx)))
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, kitHttp.EmptyOkResponse)
}

// PullLocations godoc
// @Summary triggers pulling locations from the remote platforms
// @Accept json
// @Param request body backend.PullRequest true "request object"
// @Success 200
// @Failure 500 {object} http.Error
// @Router /backend/locations/pull [post]
// @tags locations
func (c *ctrlImpl) PullLocations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[backend.PullRequest](ctx, r)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	err = c.locationUc.OnRemoteLocationsPull(ctx, rq.From, rq.To)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, kitHttp.EmptyOkResponse)
}

// GetLocation godoc
// @Summary retrieves a location object by id
// @Accept json
// @Param locId path string true "OCPI location ID"
// @Success 200 {object} backend.Location
// @Failure 500 {object} http.Error
// @Router /backend/locations/{locId} [get]
// @tags locations
func (c *ctrlImpl) GetLocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	locId, err := c.Var(ctx, r, "locId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	loc, err := c.locationService.GetLocation(ctx, locId, true)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.converter.LocationDomainToBackend(loc))
}

// SearchLocations godoc
// @Summary retrieves location objects by criteria
// @Accept json
// @Param offset query string false "number of items to offset from the beginning"
// @Param limit query string false "number of items to retrieve"
// @Param dateFrom query string false "items updated after the given date"
// @Param dateTo query string false "items updated before the given date"
// @Param refId query string false "reference id"
// @Param incPlatforms query string false "comma separated platforms to include"
// @Param excPlatforms query string false "comma separated platforms to exclude"
// @Param ids query string false "comma separated list of ids"
// @Success 200 {object} backend.LocationSearchResponse
// @Failure 500 {object} http.Error
// @Router /backend/locations/search/query [get]
// @tags locations
func (c *ctrlImpl) SearchLocations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error
	cr := &domain.LocationSearchCriteria{}

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

	rs, err := c.locationService.SearchLocations(ctx, cr)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, &backend.LocationSearchResponse{
		PageInfo: &backend.PageResponse{
			Total: rs.Total,
			Limit: rs.Limit,
		},
		Items: c.converter.LocationsDomainToBackend(rs.Items),
	})
}

// PutEvse godoc
// @Summary updates evse object in OCPI. It may contain connectors as well
// @Accept json
// @Param locId path string true "OCPI location ID"
// @Param request body backend.Evse true "evse object"
// @Success 200
// @Failure 500 {object} http.Error
// @Router /backend/locations/{locId}/evses [post]
// @tags locations
func (c *ctrlImpl) PutEvse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[backend.Evse](ctx, r)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	rq.LocationId, err = c.Var(ctx, r, "locId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	err = c.locationUc.OnLocalEvseChanged(ctx, c.converter.EvseBackendToDomain(rq, c.localPlatform.GetPlatformId(ctx), rq.LocationId))
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, kitHttp.EmptyOkResponse)
}

// SetEvseStatus godoc
// @Summary updates evse status in OCPI
// @Accept json
// @Param locId path string true "OCPI location ID"
// @Param evseId path string true "OCPI evse ID"
// @Param status query string true "status"
// @Success 200
// @Failure 500 {object} http.Error
// @Router /backend/locations/{locId}/evses/{evseId}/status [post]
// @tags locations
func (c *ctrlImpl) SetEvseStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	locId, err := c.Var(ctx, r, "locId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}
	evseId, err := c.Var(ctx, r, "evseId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}
	status, err := c.FormVal(ctx, r, "status", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	err = c.locationUc.OnLocalEvseStatusChanged(ctx, locId, evseId, status)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, kitHttp.EmptyOkResponse)
}

// GetEvse godoc
// @Summary retrieves an evse object by id
// @Accept json
// @Param locId path string true "OCPI location ID"
// @Param evseId path string true "OCPI evse ID"
// @Success 200 {object} backend.Evse
// @Failure 500 {object} http.Error
// @Router /backend/locations/{locId}/evses/{evseId} [get]
// @tags locations
func (c *ctrlImpl) GetEvse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	locId, err := c.Var(ctx, r, "locId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	evseId, err := c.Var(ctx, r, "evseId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	evse, err := c.locationService.GetEvse(ctx, locId, evseId, true)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.converter.EvseDomainToBackend(evse))
}

// SearchEvses godoc
// @Summary retrieves evse objects by criteria
// @Accept json
// @Param offset query string false "number of items to offset from the beginning"
// @Param limit query string false "number of items to retrieve"
// @Param dateFrom query string false "items updated after the given date"
// @Param dateTo query string false "items updated before the given date"
// @Success 200 {object} backend.EvseSearchResponse
// @Failure 500 {object} http.Error
// @Router /backend/evses/search/query [get]
// @tags locations
func (c *ctrlImpl) SearchEvses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error
	cr := &domain.EvseSearchCriteria{}

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
	rs, err := c.locationService.SearchEvses(ctx, cr)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, &backend.EvseSearchResponse{
		PageInfo: &backend.PageResponse{
			Total: rs.Total,
			Limit: rs.Limit,
		},
		Items: c.converter.EvsesDomainToBackend(rs.Items),
	})
}

// PutConnector godoc
// @Summary updates connector object in OCPI
// @Accept json
// @Param locId path string true "OCPI location ID"
// @Param evseId path string true "OCPI evse ID"
// @Param request body backend.Connector true "connector object"
// @Success 200
// @Failure 500 {object} http.Error
// @Router /backend/locations/{locId}/evses/{evseId}/connectors [post]
// @tags locations
func (c *ctrlImpl) PutConnector(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[backend.Connector](ctx, r)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	rq.LocationId, err = c.Var(ctx, r, "locId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}
	rq.EvseId, err = c.Var(ctx, r, "evseId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	err = c.locationUc.OnLocalConnectorChanged(ctx, c.converter.ConnectorBackendToDomain(rq, c.localPlatform.GetPlatformId(ctx), rq.LocationId, rq.EvseId))
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, kitHttp.EmptyOkResponse)
}

// GetConnector godoc
// @Summary retrieves a connector object by id
// @Accept json
// @Param locId path string true "OCPI location ID"
// @Param evseId path string true "OCPI evse ID"
// @Param conId path string true "OCPI connector ID"
// @Success 200 {object} backend.Connector
// @Failure 500 {object} http.Error
// @Router /backend/locations/{locId}/evses/{evseId}/connectors/{conId} [get]
// @tags locations
func (c *ctrlImpl) GetConnector(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	locId, err := c.Var(ctx, r, "locId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	evseId, err := c.Var(ctx, r, "evseId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	conId, err := c.Var(ctx, r, "conId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	con, err := c.locationService.GetConnector(ctx, locId, evseId, conId)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.converter.ConnectorDomainToBackend(con))
}

// SearchConnectors godoc
// @Summary retrieves connector objects by criteria
// @Accept json
// @Param offset query string false "number of items to offset from the beginning"
// @Param limit query string false "number of items to retrieve"
// @Param dateFrom query string false "items updated after the given date"
// @Param dateTo query string false "items updated before the given date"
// @Success 200 {object} backend.ConnectorSearchResponse
// @Failure 500 {object} http.Error
// @Router /backend/connectors/search/query [get]
// @tags locations
func (c *ctrlImpl) SearchConnectors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error
	cr := &domain.ConnectorSearchCriteria{}

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

	rs, err := c.locationService.SearchConnectors(ctx, cr)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, &backend.ConnectorSearchResponse{
		PageInfo: &backend.PageResponse{
			Total: rs.Total,
			Limit: rs.Limit,
		},
		Items: c.converter.ConnectorsDomainToBackend(rs.Items),
	})
}
