package impl

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/goroutine"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/errors"
	"github.com/mikhailbolshakov/ocpi/model"
	"github.com/mikhailbolshakov/ocpi/usecase"
	"time"
)

const (
	locWorkersNum = 16
	locPageSize   = 100
)

type locationUc struct {
	ucBase
	localPlatform     domain.LocalPlatformService
	locationService   domain.LocationService
	remoteLocationRep usecase.RemoteLocationRepository
	partyService      domain.PartyService
	webhook           backend.WebhookCallService
	converter         usecase.LocationConverter
}

func NewLocationUc(platformService domain.PlatformService, locationService domain.LocationService,
	remoteLocationRep usecase.RemoteLocationRepository, partyService domain.PartyService,
	webhook backend.WebhookCallService, localPlatform domain.LocalPlatformService, tokenGen domain.TokenGenerator) usecase.LocationUc {
	return &locationUc{
		ucBase:            newBase(platformService, partyService, tokenGen),
		locationService:   locationService,
		remoteLocationRep: remoteLocationRep,
		partyService:      partyService,
		webhook:           webhook,
		localPlatform:     localPlatform,
		converter:         NewLocationConverter(),
	}
}

func (l *locationUc) l() kit.CLogger {
	return ocpi.L().Cmp("loc-uc")
}

func (l *locationUc) OnLocalLocationChanged(ctx context.Context, loc *domain.Location) error {
	lg := l.l().C(ctx).Mth("on-loc-changed-loc").F(kit.KV{"locId": loc.Id}).Dbg()

	// get local platform
	localPlatform, err := l.localPlatform.Get(ctx)
	if err != nil {
		return err
	}

	// check location is of the local platform
	stored, err := l.locationService.GetLocation(ctx, loc.Id, false)
	if err != nil {
		return err
	}
	if stored != nil && stored.PlatformId != l.localPlatform.GetPlatformId(ctx) {
		return errors.ErrLocNotBelongLocalPlatform(ctx)
	}

	// put location to local platform
	loc, err = l.locationService.PutLocation(ctx, loc)
	if err != nil {
		return err
	}

	// no changes applied
	if loc == nil {
		lg.Warn("no changes applied")
		return nil
	}

	// get platforms to push location
	platforms, err := l.getPlatformsToPush(ctx, loc.PlatformId)
	if err != nil {
		return err
	}
	lg.DbgF("%d platforms to push", len(platforms))

	// for each platform
	ocpiLoc := l.converter.LocationDomainToModel(loc)
	for _, platform := range platforms {
		platform := platform
		// check if location receiver is supported by the remote platform
		ep := l.platformService.RoleEndpoint(ctx, platform, model.ModuleIdLocations, model.OcpiReceiver)
		if ep != "" && (platform.Protocol == nil || platform.Protocol.PushSupport.Locations) {
			// push location to a remote platform
			rq := buildOcpiRepositoryErrHandlerRequestG(ep, l.tokenC(platform), localPlatform, platform, ocpiLoc, lg)
			l.remoteLocationRep.PutLocationAsync(ctx, rq)
		} else {
			lg.F(kit.KV{"platform": platform.Id}).Dbg("push not supported")
		}
	}

	return nil
}

func (l *locationUc) OnRemoteLocationPut(ctx context.Context, platformId string, loc *model.OcpiLocation) error {
	l.l().C(ctx).Mth("on-loc-put-rem").F(kit.KV{"platformId": platformId, "locId": loc.Id}).Dbg()
	return l.modifyLocation(ctx, platformId, loc, l.locationService.PutLocation)
}

func (l *locationUc) OnRemoteLocationPatch(ctx context.Context, platformId string, loc *model.OcpiLocation) error {
	l.l().C(ctx).Mth("on-loc-patch-rem").F(kit.KV{"platformId": platformId, "locId": loc.Id}).Dbg()
	return l.modifyLocation(ctx, platformId, loc, l.locationService.MergeLocation)
}

func (l *locationUc) OnRemoteLocationsPull(ctx context.Context, from, to *time.Time) error {

	// get platforms to pull from
	platforms, err := l.getPlatformsToPull(ctx)
	if err != nil {
		return err
	}

	return l.remoteLocationsPull(ctx, from, to, platforms)
}

func (l *locationUc) OnRemoteLocationsPullWhenPushNotSupported(ctx context.Context, from, to *time.Time) error {

	// get platforms to pull from
	platforms, err := l.getPlatformsToPull(ctx)
	if err != nil {
		return err
	}

	// retrieve only those which don't support locations push
	platforms = kit.Filter(platforms, func(p *domain.Platform) bool { return p.Protocol != nil && !p.Protocol.PushSupport.Locations })

	return l.remoteLocationsPull(ctx, from, to, platforms)
}

func (l *locationUc) OnLocalEvseChanged(ctx context.Context, evse *domain.Evse) error {
	lg := l.l().C(ctx).Mth("on-evse-changed-loc").F(kit.KV{"evseId": evse.Id}).Dbg()

	// get local platform
	localPlatform, err := l.localPlatform.Get(ctx)
	if err != nil {
		return err
	}

	// check evse is of the local platform
	stored, err := l.locationService.GetEvse(ctx, evse.LocationId, evse.Id, false)
	if err != nil {
		return err
	}
	if stored != nil && stored.PlatformId != l.localPlatform.GetPlatformId(ctx) {
		return errors.ErrLocNotBelongLocalPlatform(ctx)
	}

	// merge evse to local platform
	evse, err = l.locationService.PutEvse(ctx, evse)
	if err != nil {
		return err
	}

	// no changes applied
	if evse == nil {
		lg.Warn("no changes applied")
		return nil
	}

	// get platforms to push evse
	platforms, err := l.getPlatformsToPush(ctx, evse.PlatformId)
	if err != nil {
		return err
	}

	// for each platform
	ocpiEvse := l.converter.EvseDomainToModel(evse)
	ocpiParty := &model.OcpiPartyId{
		PartyId:     evse.ExtId.PartyId,
		CountryCode: evse.ExtId.CountryCode,
	}
	for _, platform := range platforms {
		platform := platform
		// check if location receiver is supported by the remote platform
		ep := l.platformService.RoleEndpoint(ctx, platform, model.ModuleIdLocations, model.OcpiReceiver)
		if ep != "" && (platform.Protocol == nil || platform.Protocol.PushSupport.Locations) {
			// push evse to a remote platform
			rq := buildOcpiRepositoryErrHandlerRequestG(ep, l.tokenC(platform), localPlatform, platform, ocpiEvse, lg)
			l.remoteLocationRep.PutEvseAsync(ctx, rq, ocpiParty, evse.LocationId)
		} else {
			lg.F(kit.KV{"platform": platform.Id}).Dbg("push not supported")
		}
	}

	return nil
}

func (l *locationUc) OnLocalEvseStatusChanged(ctx context.Context, locId, evseId, status string) error {
	lg := l.l().C(ctx).Mth("on-evse-changed-loc").F(kit.KV{"evseId": evseId, "status": status}).Dbg()

	// get local platform
	localPlatform, err := l.localPlatform.Get(ctx)
	if err != nil {
		return err
	}

	// check location is of the local platform
	evse, err := l.locationService.GetEvse(ctx, locId, evseId, false)
	if err != nil {
		return err
	}
	if evse == nil {
		return errors.ErrEvseNotFound(ctx)
	}
	if evse.PlatformId != l.localPlatform.GetPlatformId(ctx) {
		return errors.ErrLocNotBelongLocalPlatform(ctx)
	}

	// merge evse
	evse.Status = status
	evse.LastUpdated = kit.Now()
	evse, err = l.locationService.MergeEvse(ctx, evse)
	if err != nil {
		return err
	}

	// no changes applied
	if evse == nil {
		lg.Warn("no changes applied")
		return nil
	}

	// get platforms to push evse
	platforms, err := l.getPlatformsToPush(ctx, evse.PlatformId)
	if err != nil {
		return err
	}

	// for each platform
	ocpiEvse := &model.OcpiEvse{
		Uid:         evse.Id,
		Status:      evse.Status,
		LastUpdated: evse.LastUpdated,
	}
	ocpiParty := &model.OcpiPartyId{
		PartyId:     evse.ExtId.PartyId,
		CountryCode: evse.ExtId.CountryCode,
	}
	for _, platform := range platforms {
		platform := platform
		// check if location receiver is supported by the remote platform
		ep := l.platformService.RoleEndpoint(ctx, platform, model.ModuleIdLocations, model.OcpiReceiver)
		if ep != "" && (platform.Protocol == nil || platform.Protocol.PushSupport.Locations) {
			// push evse to a remote platform
			rq := buildOcpiRepositoryErrHandlerRequestG(ep, l.tokenC(platform), localPlatform, platform, ocpiEvse, lg)
			l.remoteLocationRep.PatchEvseAsync(ctx, rq, ocpiParty, evse.LocationId)

		} else {
			lg.F(kit.KV{"platform": platform.Id}).Dbg("push not supported")
		}
	}

	return nil
}

func (l *locationUc) OnRemoteEvsePut(ctx context.Context, platformId, locId, countryCode, partyId string, evse *model.OcpiEvse) error {
	l.l().C(ctx).Mth("on-evse-put-rem").F(kit.KV{"platformId": platformId, "evseId": evse.Uid}).Dbg()
	return l.modifyEvse(ctx, platformId, locId, countryCode, partyId, evse, l.locationService.PutEvse)
}

func (l *locationUc) OnRemoteEvsePatch(ctx context.Context, platformId, locId, countryCode, partyId string, evse *model.OcpiEvse) error {
	l.l().C(ctx).Mth("on-evse-patch-rem").F(kit.KV{"platformId": platformId, "evseId": evse.Uid}).Dbg()
	return l.modifyEvse(ctx, platformId, locId, countryCode, partyId, evse, l.locationService.MergeEvse)
}

func (l *locationUc) OnLocalConnectorChanged(ctx context.Context, con *domain.Connector) error {
	lg := l.l().C(ctx).Mth("on-con-changed-loc").F(kit.KV{"evseId": con.EvseId, "conId": con.Id}).Dbg()

	// get local platform
	localPlatform, err := l.localPlatform.Get(ctx)
	if err != nil {
		return err
	}

	// check connector is of the local platform
	stored, err := l.locationService.GetConnector(ctx, con.LocationId, con.EvseId, con.Id)
	if err != nil {
		return err
	}
	if stored != nil && stored.PlatformId != l.localPlatform.GetPlatformId(ctx) {
		return errors.ErrLocNotBelongLocalPlatform(ctx)
	}

	// merge connector to local platform
	con, err = l.locationService.PutConnector(ctx, con)
	if err != nil {
		return err
	}

	// no changes applied
	if con == nil {
		lg.Warn("no changes applied")
		return nil
	}

	// get platforms to push connectors
	platforms, err := l.getPlatformsToPush(ctx, con.PlatformId)
	if err != nil {
		return err
	}

	// for each platform
	ocpiCon := l.converter.ConnectorDomainToModel(con)
	ocpiParty := &model.OcpiPartyId{
		PartyId:     con.ExtId.PartyId,
		CountryCode: con.ExtId.CountryCode,
	}
	for _, platform := range platforms {
		platform := platform
		// check if location receiver is supported by the remote platform
		ep := l.platformService.RoleEndpoint(ctx, platform, model.ModuleIdLocations, model.OcpiReceiver)
		if ep != "" && (platform.Protocol == nil || platform.Protocol.PushSupport.Locations) {
			// push connector to a remote platform
			rq := buildOcpiRepositoryErrHandlerRequestG(ep, l.tokenC(platform), localPlatform, platform, ocpiCon, lg)
			l.remoteLocationRep.PutConnectorAsync(ctx, rq, ocpiParty, con.LocationId, con.EvseId)
		} else {
			lg.F(kit.KV{"platform": platform.Id}).Dbg("push not supported")
		}
	}

	return nil
}

func (l *locationUc) OnRemoteConnectorPut(ctx context.Context, platformId, locId, evseId, countryCode, partyId string, con *model.OcpiConnector) error {
	l.l().C(ctx).Mth("on-con-put-rem").F(kit.KV{"platformId": platformId, "evseId": evseId, "conId": con.Id}).Dbg()
	return l.modifyConnector(ctx, platformId, locId, evseId, countryCode, partyId, con, l.locationService.PutConnector)
}

func (l *locationUc) OnRemoteConnectorPatch(ctx context.Context, platformId, locId, evseId, countryCode, partyId string, con *model.OcpiConnector) error {
	l.l().C(ctx).Mth("on-con-patch-rem").F(kit.KV{"platformId": platformId, "evseId": evseId}).Dbg()
	return l.modifyConnector(ctx, platformId, locId, evseId, countryCode, partyId, con, l.locationService.MergeConnector)
}

func (l *locationUc) getPlatformsToPush(ctx context.Context, originalPlatformId string) ([]*domain.Platform, error) {
	platforms, err := l.platformService.Search(ctx, &domain.PlatformSearchCriteria{
		Statuses: []string{domain.ConnectionStatusConnected}, // connected platforms
		ExcIds:   []string{originalPlatformId},               // exclude source platform
		Remote:   kit.BoolPtr(true),                          // remote platforms
		ExcRoles: []string{domain.RoleCPO},                   // don't push to CPO platforms
	})
	if err != nil {
		return nil, err
	}
	return platforms, nil
}

func (l *locationUc) getPlatformsToPull(ctx context.Context) ([]*domain.Platform, error) {
	platforms, err := l.platformService.Search(ctx, &domain.PlatformSearchCriteria{
		Statuses: []string{domain.ConnectionStatusConnected}, // connected platforms
		Remote:   kit.BoolPtr(true),                          // remote platforms
		ExcRoles: []string{domain.RoleEMSP},                  // don't pull from EMSP platforms
	})
	if err != nil {
		return nil, err
	}
	return platforms, nil
}

func (l *locationUc) modifyLocation(ctx context.Context, platformId string, loc *model.OcpiLocation, modifyFunc func(context.Context, *domain.Location) (*domain.Location, error)) error {
	lg := l.l().C(ctx).Mth("loc-modify").F(kit.KV{"platformId": platformId, "locId": loc.Id}).Dbg()

	// get and check platform
	_, err := l.getConnectedPlatform(ctx, platformId)
	if err != nil {
		return err
	}

	// check location is of the remote platform
	stored, err := l.locationService.GetLocation(ctx, loc.Id, false)
	if err != nil {
		return err
	}
	if stored != nil && stored.PlatformId == l.localPlatform.GetPlatformId(ctx) {
		return errors.ErrLocNotBelongRemotePlatform(ctx)
	}

	// get or create party
	_, err = l.getCreateParty(ctx, platformId, loc.PartyId, loc.CountryCode)
	if err != nil {
		return err
	}

	// modify location to local platform
	locDom, err := modifyFunc(ctx, l.converter.LocationModelToDomain(loc, platformId))
	if err != nil {
		return err
	}

	// no changes applied
	if locDom == nil {
		lg.Warn("no changes applied")
		return nil
	}

	// call webhook
	return l.webhook.OnLocationsChanged(ctx, l.converter.LocationDomainToBackend(locDom))
}

func (l *locationUc) modifyEvse(ctx context.Context, platformId, locId, countryCode, partyId string, evse *model.OcpiEvse, modifyFunc func(context.Context, *domain.Evse) (*domain.Evse, error)) error {
	lg := l.l().C(ctx).Mth("evse-modify").F(kit.KV{"platformId": platformId, "evseId": evse.Uid}).Dbg()

	// get and check platform
	_, err := l.getConnectedPlatform(ctx, platformId)
	if err != nil {
		return err
	}

	// check evse is of the remote platform
	stored, err := l.locationService.GetEvse(ctx, locId, evse.Uid, false)
	if err != nil {
		return err
	}
	if stored != nil && stored.PlatformId == l.localPlatform.GetPlatformId(ctx) {
		return errors.ErrLocNotBelongRemotePlatform(ctx)
	}

	// get or create party
	_, err = l.getCreateParty(ctx, platformId, partyId, countryCode)
	if err != nil {
		return err
	}

	// merge evse to the local platform
	evseDom, err := modifyFunc(ctx, l.converter.EvseModelToDomain(evse, countryCode, partyId, platformId, locId))
	if err != nil {
		return err
	}

	// no changes applied
	if evseDom == nil {
		lg.Warn("no changes applied")
		return nil
	}

	// call webhook
	return l.webhook.OnEvseChanged(ctx, l.converter.EvseDomainToBackend(evseDom))
}

func (l *locationUc) modifyConnector(ctx context.Context, platformId, locId, evseId, countryCode, partyId string, con *model.OcpiConnector, modifyFunc func(context.Context, *domain.Connector) (*domain.Connector, error)) error {
	lg := l.l().C(ctx).Mth("con-modify").F(kit.KV{"platformId": platformId, "evseId": evseId}).Dbg()

	// get and check platform
	_, err := l.getConnectedPlatform(ctx, platformId)
	if err != nil {
		return err
	}

	// check con is of the remote platform
	stored, err := l.locationService.GetConnector(ctx, locId, evseId, con.Id)
	if err != nil {
		return err
	}
	if stored != nil && stored.PlatformId == l.localPlatform.GetPlatformId(ctx) {
		return errors.ErrLocNotBelongRemotePlatform(ctx)
	}

	// get or create party
	_, err = l.getCreateParty(ctx, platformId, partyId, countryCode)
	if err != nil {
		return err
	}

	// patch connector to local platform
	conDom, err := modifyFunc(ctx, l.converter.ConnectorModelToDomain(con, countryCode, partyId, platformId, locId, evseId))
	if err != nil {
		return err
	}

	// no changes applied
	if conDom == nil {
		lg.Warn("no changes applied")
		return nil
	}

	// call webhook
	return l.webhook.OnConnectorChanged(ctx, l.converter.ConnectorDomainToBackend(conDom))
}

func (l *locationUc) remoteLocationsPull(ctx context.Context, from, to *time.Time, platforms []*domain.Platform) error {
	lg := l.l().C(ctx).Mth("remote-pull").Dbg()

	if len(platforms) == 0 {
		return nil
	}

	// get local platform
	localPlatform, err := l.localPlatform.Get(ctx)
	if err != nil {
		return err
	}

	ctx = kit.Copy(ctx)

	// execute workers
	ch := make(chan channelData, locWorkersNum)
	for i := 0; i < locWorkersNum; i++ {
		goroutine.New().WithLogger(lg).Go(ctx, func() {
			for v := range ch {
				loc := v.data.(*model.OcpiLocation)
				err := l.OnRemoteLocationPut(ctx, v.platformId, loc)
				if err != nil {
					l.l().C(ctx).Mth("remote-put").
						F(kit.KV{"platformId": v.platformId, "locId": loc.Id}).
						E(err).St().Err()
				}
			}
		})
	}

	// for each platform execute pagers and put retrieved items to worker routines
	pgReader := NewPageReader(l.remoteLocationRep.GetLocations)
	eg := goroutine.NewGroup(ctx).WithLogger(lg)
	for _, platform := range platforms {
		platform := platform
		// check if sender isn't supported by the remote platform
		ep := l.platformService.RoleEndpoint(ctx, platform, model.ModuleIdLocations, model.OcpiSender)
		if ep == "" {
			continue
		}
		eg.Go(func() error {
			for locs := range pgReader.GetPage(ctx, buildOcpiRepositoryRequest(ep, l.tokenC(platform), localPlatform, platform), locPageSize, from, to) {
				for _, loc := range locs {
					ch <- channelData{platformId: platform.Id, data: loc}
				}
			}
			return nil
		})
	}

	// close channels when done
	goroutine.New().WithLogger(lg).Go(ctx, func() {
		defer close(ch)
		_ = eg.Wait()
	})

	return nil
}
