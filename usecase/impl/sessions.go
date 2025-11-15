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
	sessWorkersNum = 4
	sessPageSize   = 20
)

type sessionUc struct {
	ucBase
	sessionService       domain.SessionService
	remoteSessionRep     usecase.RemoteSessionRepository
	partyService         domain.PartyService
	webhook              backend.WebhookCallService
	converter            usecase.SessionConverter
	cmdService           domain.CommandService
	localPlatformService domain.LocalPlatformService
	tokenService         domain.TokenService
}

func NewSessionUc(platformService domain.PlatformService, sessionService domain.SessionService, remoteSessionRep usecase.RemoteSessionRepository,
	partyService domain.PartyService, webhook backend.WebhookCallService, cmdService domain.CommandService, localPlatformService domain.LocalPlatformService,
	tokenService domain.TokenService, tokenGen domain.TokenGenerator) usecase.SessionUc {
	return &sessionUc{
		ucBase:               newBase(platformService, partyService, tokenGen),
		sessionService:       sessionService,
		remoteSessionRep:     remoteSessionRep,
		partyService:         partyService,
		webhook:              webhook,
		cmdService:           cmdService,
		localPlatformService: localPlatformService,
		tokenService:         tokenService,
		converter:            NewSessionConverter(),
	}
}

func (s *sessionUc) l() kit.CLogger {
	return ocpi.L().Cmp("sess-uc")
}

func (s *sessionUc) OnLocalSessionChanged(ctx context.Context, sess *domain.Session) error {
	l := s.l().C(ctx).Mth("on-sess-changed-loc").F(kit.KV{"sessId": sess.Id}).Dbg()

	// check token params
	if sess.Details.CdrToken == nil || sess.Details.CdrToken.Id == "" {
		return errors.ErrCdrTokenEmpty(ctx)
	}

	// get local platform
	localPlatform, err := s.localPlatformService.Get(ctx)
	if err != nil {
		return err
	}

	// get token
	tkn, err := s.mustGetToken(ctx, sess.Details.CdrToken.Id)
	if err != nil {
		return err
	}

	// set cdr token
	sess.Details.CdrToken = s.tokenToCdrToken(tkn)

	// merge session to local platform
	sess, err = s.sessionService.PutSession(ctx, sess)
	if err != nil {
		return err
	}

	// no changes applied
	if sess == nil {
		l.Warn("no changes applied")
		return nil
	}

	// get platform
	platform, err := s.getConnectedPlatform(ctx, tkn.PlatformId)
	if err != nil {
		return err
	}

	// set header to route message
	s.setFromPartyCtx(ctx, sess.ExtId)
	s.setToPartyCtx(ctx, tkn.ExtId)

	ep := s.platformService.RoleEndpoint(ctx, platform, model.ModuleIdSessions, model.OcpiReceiver)
	if ep != "" && (platform.Protocol == nil || platform.Protocol.PushSupport.Sessions) {
		// push session to a remote platform
		rq := buildOcpiRepositoryErrHandlerRequestG(ep, s.tokenC(platform), localPlatform, platform, s.converter.SessionDomainToModel(sess), l)
		s.remoteSessionRep.PutSessionAsync(ctx, rq)
	} else {
		l.F(kit.KV{"platform": platform.Id}).Dbg("push not supported")
	}

	return nil
}

func (s *sessionUc) OnLocalSessionPatched(ctx context.Context, sess *domain.Session) error {
	l := s.l().C(ctx).Mth("on-sess-patched-loc").F(kit.KV{"sessId": sess.Id}).Dbg()

	// get local platform
	localPlatform, err := s.localPlatformService.Get(ctx)
	if err != nil {
		return err
	}

	// get token id from CDR token
	tokenId, err := s.mustGetTokenId(ctx, sess)
	if err != nil {
		return err
	}

	// get token
	tkn, err := s.tokenService.GetToken(ctx, tokenId)
	if err != nil {
		return err
	}

	// if cdr token specified, populate it
	if sess.Details.CdrToken != nil {
		sess.Details.CdrToken = s.tokenToCdrToken(tkn)
	}

	// merge session to local platform
	sess, err = s.sessionService.MergeSession(ctx, sess)
	if err != nil {
		return err
	}

	// no changes applied
	if sess == nil {
		l.Warn("no changes applied")
		return nil
	}

	// get platform to push session
	platform, err := s.getConnectedPlatform(ctx, tkn.PlatformId)
	if err != nil {
		return err
	}

	// set header to route message
	s.setFromPartyCtx(ctx, sess.ExtId)
	s.setToPartyCtx(ctx, tkn.ExtId)

	ep := s.platformService.RoleEndpoint(ctx, platform, model.ModuleIdSessions, model.OcpiReceiver)
	if ep != "" && (platform.Protocol == nil || platform.Protocol.PushSupport.Sessions) {
		// push session to a remote platform
		// we are using PUT here to update the full session to avoid race condition with multiple partial requests
		rq := buildOcpiRepositoryErrHandlerRequestG(ep, s.tokenC(platform), localPlatform, platform, s.converter.SessionDomainToModel(sess), l)
		s.remoteSessionRep.PatchSessionAsync(ctx, rq)
	} else {
		l.F(kit.KV{"platform": platform.Id}).Dbg("push not supported")
	}

	return nil
}

func (s *sessionUc) OnRemoteSessionsPull(ctx context.Context, from, to *time.Time) error {
	// get platforms to pull from
	platforms, err := s.getPlatformsToPull(ctx)
	if err != nil {
		return err
	}

	return s.remoteSessionsPull(ctx, from, to, platforms)
}

func (s *sessionUc) OnRemoteSessionsPullWhenPushNotSupported(ctx context.Context, from, to *time.Time) error {
	// get platforms to pull from
	platforms, err := s.getPlatformsToPull(ctx)
	if err != nil {
		return err
	}

	// retrieve only those which don't support sessions push
	platforms = kit.Filter(platforms, func(p *domain.Platform) bool { return p.Protocol != nil && !p.Protocol.PushSupport.Sessions })

	return s.remoteSessionsPull(ctx, from, to, platforms)
}

func (s *sessionUc) OnRemoteSessionPut(ctx context.Context, platformId string, sess *model.OcpiSession) error {
	s.l().C(ctx).Mth("on-sess-put-rem").F(kit.KV{"platformId": platformId, "locId": sess.Id}).Dbg()
	return s.modifySession(ctx, platformId, sess, s.sessionService.PutSession)
}

func (s *sessionUc) OnRemoteSessionPatch(ctx context.Context, platformId string, sess *model.OcpiSession) error {
	s.l().C(ctx).Mth("on-sess-patch-rem").F(kit.KV{"platformId": platformId, "locId": sess.Id}).Dbg()
	return s.modifySession(ctx, platformId, sess, s.sessionService.MergeSession)
}

func (s *sessionUc) getPlatformsToPull(ctx context.Context) ([]*domain.Platform, error) {
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

func (s *sessionUc) modifySession(ctx context.Context, platformId string, sess *model.OcpiSession, modifyFunc func(context.Context, *domain.Session) (*domain.Session, error)) error {
	l := s.l().C(ctx).Mth("sess-modify").F(kit.KV{"platformId": platformId, "sessId": sess.Id}).Dbg()

	// get and check platform
	_, err := s.getConnectedPlatform(ctx, platformId)
	if err != nil {
		return err
	}

	// check session is of the remote platform
	stored, err := s.sessionService.GetSession(ctx, sess.Id)
	if err != nil {
		return err
	}
	if stored != nil && stored.PlatformId == s.localPlatformService.GetPlatformId(ctx) {
		return errors.ErrSessCmdInvalidPlatform(ctx)
	}

	// merge evse to the local platform
	sessDom, err := modifyFunc(ctx, s.converter.SessionModelToDomain(sess, platformId))
	if err != nil {
		return err
	}

	// no changes
	if sessDom == nil {
		l.Warn("no changes applied")
		return nil
	}

	// call webhook
	return s.webhook.OnSessionsChanged(ctx, s.converter.SessionDomainToBackend(sessDom))
}

func (s *sessionUc) mustGetToken(ctx context.Context, id string) (*domain.Token, error) {
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

func (s *sessionUc) mustGetTokenId(ctx context.Context, sess *domain.Session) (string, error) {
	if sess.Details.CdrToken != nil && sess.Details.CdrToken.Id != "" {
		return sess.Details.CdrToken.Id, nil
	}
	fullSess, err := s.sessionService.GetSession(ctx, sess.Id)
	if err != nil {
		return "", err
	}
	return fullSess.Details.CdrToken.Id, nil
}

func (s *sessionUc) remoteSessionsPull(ctx context.Context, from, to *time.Time, platforms []*domain.Platform) error {
	l := s.l().C(ctx).Mth("remote-pull").Dbg()

	if len(platforms) == 0 {
		return nil
	}

	// get local platform
	localPlatform, err := s.localPlatformService.Get(ctx)
	if err != nil {
		return err
	}

	ctx = kit.Copy(ctx)

	// execute workers
	ch := make(chan channelData, sessWorkersNum)
	for i := 0; i < sessWorkersNum; i++ {
		goroutine.New().WithLogger(l).Go(ctx, func() {
			for v := range ch {
				sess := v.data.(*model.OcpiSession)
				err := s.OnRemoteSessionPut(ctx, v.platformId, sess)
				if err != nil {
					s.l().C(ctx).Mth("remote-put").
						F(kit.KV{"platformId": v.platformId, "sessId": sess.Id}).
						E(err).St().Err()
				}
			}
		})
	}

	// for each platform execute pagers and put retrieved items to worker routines
	pgReader := NewPageReader(s.remoteSessionRep.GetSessions)
	eg := goroutine.NewGroup(ctx).WithLogger(l)
	for _, platform := range platforms {
		platform := platform
		// check if sender isn't supported by the remote platform
		ep := s.platformService.RoleEndpoint(ctx, platform, model.ModuleIdSessions, model.OcpiSender)
		if ep == "" {
			continue
		}
		eg.Go(func() error {
			for sessions := range pgReader.GetPage(ctx, buildOcpiRepositoryRequest(ep, s.tokenC(platform), localPlatform, platform), sessPageSize, from, to) {
				for _, sess := range sessions {
					ch <- channelData{platformId: platform.Id, data: sess}
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
