package locations

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
	SenderGetLocations(http.ResponseWriter, *http.Request)
	SenderGetLocation(http.ResponseWriter, *http.Request)
	SenderGetEvse(http.ResponseWriter, *http.Request)
	SenderGetConnector(http.ResponseWriter, *http.Request)

	ReceiverGetLocation(http.ResponseWriter, *http.Request)
	ReceiverGetEvse(http.ResponseWriter, *http.Request)
	ReceiverGetConnector(http.ResponseWriter, *http.Request)
	ReceiverPutLocation(http.ResponseWriter, *http.Request)
	ReceiverPutEvse(http.ResponseWriter, *http.Request)
	ReceiverPutConnector(http.ResponseWriter, *http.Request)
	ReceiverPatchLocation(http.ResponseWriter, *http.Request)
	ReceiverPatchEvse(http.ResponseWriter, *http.Request)
	ReceiverPatchConnector(http.ResponseWriter, *http.Request)
}

type ctrlImpl struct {
	ocpi.Controller
	locationUc      usecase.LocationUc
	locationService domain.LocationService
	localPlatform   domain.LocalPlatformService
	converter       usecase.LocationConverter
	senderSearchUrl string
}

func NewController(locationUc usecase.LocationUc, locationService domain.LocationService, localPlatform domain.LocalPlatformService,
	converter usecase.LocationConverter, cfg *cfg.CfgOcpiConfig) Controller {
	return &ctrlImpl{
		locationUc:      locationUc,
		localPlatform:   localPlatform,
		locationService: locationService,
		Controller:      ocpi.NewController(),
		converter:       converter,
		senderSearchUrl: fmt.Sprintf("%s/%s", cfg.Local.Url, "ocpi/2.2.1/sender/locations"),
	}
}

func (c *ctrlImpl) SenderGetLocations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve locations by the local platform
	rq := &domain.LocationSearchCriteria{
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

	rs, err := c.locationService.SearchLocations(ctx, rq)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	r = r.WithContext(c.SetResponseWithNextPageCtx(ctx, rs.PageResponse, c.senderSearchUrl))

	c.OcpiRespondOK(r, w, c.converter.LocationsDomainToModel(rs.Items))
}

func (c *ctrlImpl) SenderGetLocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	locationId, err := c.Var(ctx, r, "location_id", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	loc, err := c.locationService.GetLocation(ctx, locationId, true)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	if loc == nil {
		c.OcpiRespondNotFoundError(r, w)
		return
	}

	c.OcpiRespondOK(r, w, c.converter.LocationDomainToModel(loc))

}

func (c *ctrlImpl) SenderGetEvse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	locationId, err := c.Var(ctx, r, "location_id", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	evseId, err := c.Var(ctx, r, "evse_uid", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	evse, err := c.locationService.GetEvse(ctx, locationId, evseId, true)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	if evse == nil {
		c.OcpiRespondNotFoundError(r, w)
		return
	}

	c.OcpiRespondOK(r, w, c.converter.EvseDomainToModel(evse))
}

func (c *ctrlImpl) SenderGetConnector(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	locationId, err := c.Var(ctx, r, "location_id", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	evseId, err := c.Var(ctx, r, "evse_uid", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	conId, err := c.Var(ctx, r, "connector_id", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	con, err := c.locationService.GetConnector(ctx, locationId, evseId, conId)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	if con == nil {
		c.OcpiRespondNotFoundError(r, w)
		return
	}

	c.OcpiRespondOK(r, w, c.converter.ConnectorDomainToModel(con))
}

func (c *ctrlImpl) ReceiverGetLocation(w http.ResponseWriter, r *http.Request) {
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
	locationId, err := c.Var(ctx, r, "location_id", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	loc, err := c.locationService.GetLocation(ctx, locationId, true)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	if loc == nil {
		c.OcpiRespondNotFoundError(r, w)
		return
	}

	c.OcpiRespondOK(r, w, c.converter.LocationDomainToModel(loc))
}

func (c *ctrlImpl) ReceiverGetEvse(w http.ResponseWriter, r *http.Request) {
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
	locationId, err := c.Var(ctx, r, "location_id", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	evseId, err := c.Var(ctx, r, "evse_uid", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	evse, err := c.locationService.GetEvse(ctx, locationId, evseId, true)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	if evse == nil {
		c.OcpiRespondNotFoundError(r, w)
		return
	}

	c.OcpiRespondOK(r, w, c.converter.EvseDomainToModel(evse))
}

func (c *ctrlImpl) ReceiverGetConnector(w http.ResponseWriter, r *http.Request) {
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
	locationId, err := c.Var(ctx, r, "location_id", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	evseId, err := c.Var(ctx, r, "evse_uid", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	conId, err := c.Var(ctx, r, "connector_id", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	con, err := c.locationService.GetConnector(ctx, locationId, evseId, conId)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	if con == nil {
		c.OcpiRespondNotFoundError(r, w)
		return
	}

	c.OcpiRespondOK(r, w, c.converter.ConnectorDomainToModel(con))
}

func (c *ctrlImpl) ReceiverPutLocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[model.OcpiLocation](ctx, r)
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
	rq.Id, err = c.Var(ctx, r, "location_id", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	platformId, err := c.PlatformId(ctx)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	err = c.locationUc.OnRemoteLocationPut(ctx, platformId, rq)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	c.OcpiRespondOK(r, w, nil)
}

func (c *ctrlImpl) ReceiverPutEvse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[model.OcpiEvse](ctx, r)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	partyId, err := c.Var(ctx, r, model.OcpiQueryParamPartyId, false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	countryCode, err := c.Var(ctx, r, model.OcpiQueryParamCountryCode, false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	locationId, err := c.Var(ctx, r, "location_id", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	rq.Uid, err = c.Var(ctx, r, "evse_uid", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	platformId, err := c.PlatformId(ctx)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	err = c.locationUc.OnRemoteEvsePut(ctx, platformId, locationId, countryCode, partyId, rq)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	c.OcpiRespondOK(r, w, nil)
}

func (c *ctrlImpl) ReceiverPutConnector(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[model.OcpiConnector](ctx, r)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	partyId, err := c.Var(ctx, r, model.OcpiQueryParamPartyId, false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	countryCode, err := c.Var(ctx, r, model.OcpiQueryParamCountryCode, false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	locationId, err := c.Var(ctx, r, "location_id", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	evseId, err := c.Var(ctx, r, "evse_uid", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	rq.Id, err = c.Var(ctx, r, "connector_id", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	platformId, err := c.PlatformId(ctx)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	err = c.locationUc.OnRemoteConnectorPut(ctx, platformId, locationId, evseId, countryCode, partyId, rq)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	c.OcpiRespondOK(r, w, nil)
}

func (c *ctrlImpl) ReceiverPatchLocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[model.OcpiLocation](ctx, r)
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
	rq.Id, err = c.Var(ctx, r, "location_id", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	platformId, err := c.PlatformId(ctx)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	err = c.locationUc.OnRemoteLocationPatch(ctx, platformId, rq)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	c.OcpiRespondOK(r, w, nil)
}

func (c *ctrlImpl) ReceiverPatchEvse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[model.OcpiEvse](ctx, r)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	partyId, err := c.Var(ctx, r, model.OcpiQueryParamPartyId, false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	countryCode, err := c.Var(ctx, r, model.OcpiQueryParamCountryCode, false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	locationId, err := c.Var(ctx, r, "location_id", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	rq.Uid, err = c.Var(ctx, r, "evse_uid", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	platformId, err := c.PlatformId(ctx)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	err = c.locationUc.OnRemoteEvsePatch(ctx, platformId, locationId, countryCode, partyId, rq)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	c.OcpiRespondOK(r, w, nil)
}

func (c *ctrlImpl) ReceiverPatchConnector(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[model.OcpiConnector](ctx, r)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	partyId, err := c.Var(ctx, r, model.OcpiQueryParamPartyId, false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	countryCode, err := c.Var(ctx, r, model.OcpiQueryParamCountryCode, false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	locationId, err := c.Var(ctx, r, "location_id", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	evseId, err := c.Var(ctx, r, "evse_uid", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}
	rq.Id, err = c.Var(ctx, r, "connector_id", false)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	platformId, err := c.PlatformId(ctx)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	err = c.locationUc.OnRemoteConnectorPatch(ctx, platformId, locationId, evseId, countryCode, partyId, rq)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	c.OcpiRespondOK(r, w, nil)
}
