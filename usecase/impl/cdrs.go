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
	cdrWorkersNum = 4
	cdrPageSize   = 20
)

type cdrUc struct {
	ucBase
	cdrService           domain.CdrService
	remoteCdrRep         usecase.RemoteCdrRepository
	partyService         domain.PartyService
	webhook              backend.WebhookCallService
	converter            usecase.CdrConverter
	sessService          domain.SessionService
	locService           domain.LocationService
	tariffService        domain.TariffService
	localPlatformService domain.LocalPlatformService
	tokenService         domain.TokenService
}

func NewCdrUc(platformService domain.PlatformService, cdrService domain.CdrService, remoteCdrRep usecase.RemoteCdrRepository,
	partyService domain.PartyService, webhook backend.WebhookCallService, sessService domain.SessionService, localPlatformService domain.LocalPlatformService,
	locService domain.LocationService, tariffService domain.TariffService, tokenService domain.TokenService, tokenGen domain.TokenGenerator) usecase.CdrUc {
	return &cdrUc{
		ucBase:               newBase(platformService, partyService, tokenGen),
		cdrService:           cdrService,
		remoteCdrRep:         remoteCdrRep,
		partyService:         partyService,
		webhook:              webhook,
		sessService:          sessService,
		localPlatformService: localPlatformService,
		locService:           locService,
		tariffService:        tariffService,
		tokenService:         tokenService,
		converter:            NewCdrConverter(NewTariffConverter()),
	}
}

func (s *cdrUc) l() kit.CLogger {
	return ocpi.L().Cmp("cdr-uc")
}

func (s *cdrUc) OnLocalCdrChanged(ctx context.Context, cdr *backend.Cdr) error {
	l := s.l().C(ctx).Mth("on-cdr-changed-loc").F(kit.KV{"cdrId": cdr.Id}).Dbg()

	// local platform
	localPlatform, err := s.localPlatformService.Get(ctx)
	if err != nil {
		return err
	}

	// get valid session
	sess, err := s.mustGetLocalSession(ctx, cdr.SessionId)
	if err != nil {
		return err
	}

	// check and get token
	if sess.Details.CdrToken == nil {
		return errors.ErrCdrTokenEmpty(ctx)
	}
	tkn, err := s.mustGetToken(ctx, sess.Details.CdrToken.Id)
	if err != nil {
		return err
	}

	// get platform
	platform, err := s.getConnectedPlatform(ctx, tkn.PlatformId)
	if err != nil {
		return err
	}

	// get location-evse-connector
	loc, evse, con, err := s.mustGetLocalConnector(ctx, sess.Details.LocationId, sess.Details.EvseId, sess.Details.ConnectorId)
	if err != nil {
		return err
	}

	// convert to domain
	cdrDom := s.converter.CdrBackendToDomain(cdr, sess, loc, evse, con)

	// merge cdr to local platform
	cdrDom, err = s.cdrService.PutCdr(ctx, cdrDom)
	if err != nil {
		return err
	}

	// no changes applied
	if cdrDom == nil {
		l.Warn("no changes applied")
		return nil
	}

	// set header to route message
	s.setFromPartyCtx(ctx, sess.ExtId)
	s.setToPartyCtx(ctx, tkn.ExtId)

	ep := s.platformService.RoleEndpoint(ctx, platform, model.ModuleIdCdrs, model.OcpiReceiver)
	if ep != "" && (platform.Protocol == nil || platform.Protocol.PushSupport.Cdrs) {
		// push cdr to a remote platform
		ocpiRq := buildOcpiRepositoryErrHandlerRequestG(ep, s.tokenC(platform), localPlatform, platform, s.converter.CdrDomainToModel(cdrDom), l)
		s.remoteCdrRep.PostCdrAsync(ctx, ocpiRq)
	} else {
		l.F(kit.KV{"platform": platform.Id}).Dbg("push not supported")
	}

	return nil
}

func (s *cdrUc) OnRemoteCdrsPull(ctx context.Context, from, to *time.Time) error {

	// get platforms to pull from
	platforms, err := s.getPlatformsToPull(ctx)
	if err != nil {
		return err
	}

	return s.remoteCdrsPull(ctx, from, to, platforms)
}

func (s *cdrUc) OnRemoteCdrsPullWhenPushNotSupported(ctx context.Context, from, to *time.Time) error {
	// get platforms to pull from
	platforms, err := s.getPlatformsToPull(ctx)
	if err != nil {
		return err
	}

	// retrieve only those which don't support cdrs push
	platforms = kit.Filter(platforms, func(p *domain.Platform) bool { return p.Protocol != nil && !p.Protocol.PushSupport.Cdrs })

	return s.remoteCdrsPull(ctx, from, to, platforms)
}

func (s *cdrUc) OnRemoteCdrPut(ctx context.Context, platformId string, cdr *model.OcpiCdr) error {
	l := s.l().C(ctx).Mth("on-cdr-put-rem").F(kit.KV{"platformId": platformId, "locId": cdr.Id}).Dbg()

	// get and check platform
	_, err := s.getConnectedPlatform(ctx, platformId)
	if err != nil {
		return err
	}

	// check cdr is of the remote platform
	stored, err := s.cdrService.GetCdr(ctx, cdr.Id)
	if err != nil {
		return err
	}
	if stored != nil && stored.PlatformId == s.localPlatformService.GetPlatformId(ctx) {
		return errors.ErrCdrInvalidPlatform(ctx)
	}

	// get session
	_, err = s.mustGetRemoteSession(ctx, cdr.SessionId)
	if err != nil {
		return err
	}

	// merge evse to the local platform
	cdrDom, err := s.cdrService.PutCdr(ctx, s.converter.CdrModelToDomain(cdr, platformId))
	if err != nil {
		return err
	}

	// no changes
	if cdrDom == nil {
		l.Warn("no changes applied")
		return nil
	}

	// call webhook
	return s.webhook.OnCdrChanged(ctx, s.converter.CdrDomainToBackend(cdrDom))
}

func (s *cdrUc) getPlatformsToPull(ctx context.Context) ([]*domain.Platform, error) {
	platforms, err := s.platformService.Search(ctx, &domain.PlatformSearchCriteria{
		Statuses: []string{domain.ConnectionStatusConnected}, // connected platforms
		Remote:   kit.BoolPtr(true),                          // remote platforms
		ExcRoles: []string{domain.RoleCPO},                   // don't pull from EMSP platforms
	})
	if err != nil {
		return nil, err
	}
	return platforms, nil
}

func (s *cdrUc) mustGetSession(ctx context.Context, sessionId string) (*domain.Session, error) {
	// check if sessionId exists
	if sessionId == "" {
		return nil, errors.ErrCdrSessionIdEmpty(ctx)
	}

	// get session & check
	sess, err := s.sessService.GetSession(ctx, sessionId)
	if err != nil {
		return nil, err
	}
	if sess == nil {
		return nil, errors.ErrSessNotFound(ctx)
	}

	return sess, nil
}

func (s *cdrUc) mustGetConnector(ctx context.Context, locationId, evseId, connectorId string) (*domain.Location, *domain.Evse, *domain.Connector, error) {

	if locationId == "" {
		return nil, nil, nil, errors.ErrLocIdEmpty(ctx)
	}
	if evseId == "" {
		return nil, nil, nil, errors.ErrEvseIdEmpty(ctx)
	}
	if connectorId == "" {
		return nil, nil, nil, errors.ErrConIdEmpty(ctx)
	}

	// get session & check
	loc, err := s.locService.GetLocation(ctx, locationId, true)
	if err != nil {
		return nil, nil, nil, err
	}
	if loc == nil {
		return nil, nil, nil, errors.ErrLocationNotFound(ctx)
	}

	var evse *domain.Evse
	var con *domain.Connector
	for _, e := range loc.Evses {
		if e.Id == evseId {
			evse = e
			for _, c := range evse.Connectors {
				if c.Id == connectorId {
					con = c
					break
				}
			}
			break
		}
	}

	if evse == nil {
		return nil, nil, nil, errors.ErrEvseNotFound(ctx)
	}
	if con == nil {
		return nil, nil, nil, errors.ErrConNotFound(ctx)
	}

	return loc, evse, con, nil
}

func (s *cdrUc) mustGetRemoteSession(ctx context.Context, sessionId string) (*domain.Session, error) {

	sess, err := s.mustGetSession(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	// check if a session belongs to the remote platform
	if sess.PlatformId == s.localPlatformService.GetPlatformId(ctx) {
		return nil, errors.ErrCdrSessInvalidPlatform(ctx)
	}
	return sess, nil
}

func (s *cdrUc) mustGetLocalSession(ctx context.Context, sessionId string) (*domain.Session, error) {

	sess, err := s.mustGetSession(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	// check if a session belongs to the remote platform
	if sess.PlatformId != s.localPlatformService.GetPlatformId(ctx) {
		return nil, errors.ErrCdrSessInvalidPlatform(ctx)
	}
	return sess, nil
}

func (s *cdrUc) mustGetLocalConnector(ctx context.Context, locationId, evseId, connectorId string) (*domain.Location, *domain.Evse, *domain.Connector, error) {

	loc, evse, con, err := s.mustGetConnector(ctx, locationId, evseId, connectorId)
	if err != nil {
		return nil, nil, nil, err
	}

	// check if a session belongs to the remote platform
	if loc.PlatformId != s.localPlatformService.GetPlatformId(ctx) {
		return nil, nil, nil, errors.ErrCdrLocInvalidPlatform(ctx)
	}
	return loc, evse, con, nil
}

func (s *cdrUc) mustGetRemoteConnector(ctx context.Context, locationId, evseId, connectorId string) (*domain.Location, *domain.Evse, *domain.Connector, error) {

	loc, evse, con, err := s.mustGetConnector(ctx, locationId, evseId, connectorId)
	if err != nil {
		return nil, nil, nil, err
	}

	// check if a session belongs to the remote platform
	if loc.PlatformId == s.localPlatformService.GetPlatformId(ctx) {
		return nil, nil, nil, errors.ErrCdrLocInvalidPlatform(ctx)
	}
	return loc, evse, con, nil
}

func (s *cdrUc) mustGetToken(ctx context.Context, id string) (*domain.Token, error) {
	if id == "" {
		return nil, errors.ErrTknIdEmpty(ctx)
	}

	// get token
	tkn, err := s.tokenService.GetToken(ctx, id)
	if err != nil {
		return nil, err
	}
	if tkn == nil {
		return nil, errors.ErrTknNotFound(ctx)
	}

	return tkn, nil
}

func (s *cdrUc) remoteCdrsPull(ctx context.Context, from, to *time.Time, platforms []*domain.Platform) error {
	l := s.l().C(ctx).Mth("remote-pull").Dbg()

	// get local platform
	localPlatform, err := s.localPlatformService.Get(ctx)
	if err != nil {
		return err
	}

	ctx = kit.Copy(ctx)

	// execute workers
	ch := make(chan channelData, cdrWorkersNum)
	for i := 0; i < cdrWorkersNum; i++ {
		goroutine.New().WithLogger(l).Go(ctx, func() {
			for v := range ch {
				cdr := v.data.(*model.OcpiCdr)
				err := s.OnRemoteCdrPut(ctx, v.platformId, cdr)
				if err != nil {
					s.l().C(ctx).Mth("remote-put").
						F(kit.KV{"platformId": v.platformId, "cdrId": cdr.Id}).
						E(err).St().Err()
				}
			}
		})
	}

	// for each platform execute pagers and put retrieved items to worker routines
	pgReader := NewPageReader(s.remoteCdrRep.GetCdrs)
	eg := goroutine.NewGroup(ctx).WithLogger(l)
	for _, platform := range platforms {
		platform := platform
		// check if sender isn't supported by the remote platform
		ep := s.platformService.RoleEndpoint(ctx, platform, model.ModuleIdCdrs, model.OcpiSender)
		if ep == "" {
			continue
		}
		eg.Go(func() error {
			for cdrs := range pgReader.GetPage(ctx, buildOcpiRepositoryRequest(ep, s.tokenC(platform), localPlatform, platform), cdrPageSize, from, to) {
				for _, cdr := range cdrs {
					ch <- channelData{platformId: platform.Id, data: cdr}
				}
			}
			return nil
		})
	}

	// close channels when done
	goroutine.New().WithLogger(l).Go(ctx, func() {
		defer close(ch)
		_ = eg.Wait()
	})

	return nil
}
