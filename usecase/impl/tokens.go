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
	tknWorkersNum = 4
	tknPageSize   = 100
)

type tokenUc struct {
	ucBase
	localPlatform  domain.LocalPlatformService
	tokenService   domain.TokenService
	remoteTokenRep usecase.RemoteTokenRepository
	partyService   domain.PartyService
	webhook        backend.WebhookCallService
	converter      usecase.TokenConverter
}

func NewTokenUc(platformService domain.PlatformService, tokenService domain.TokenService, remoteTokenRep usecase.RemoteTokenRepository,
	partyService domain.PartyService, webhook backend.WebhookCallService, localPlatform domain.LocalPlatformService, tokenGen domain.TokenGenerator) usecase.TokenUc {
	return &tokenUc{
		ucBase:         newBase(platformService, partyService, tokenGen),
		tokenService:   tokenService,
		remoteTokenRep: remoteTokenRep,
		partyService:   partyService,
		localPlatform:  localPlatform,
		webhook:        webhook,
		converter:      NewTokenConverter(),
	}
}

func (t *tokenUc) l() kit.CLogger {
	return ocpi.L().Cmp("tkn-uc")
}

func (t *tokenUc) OnLocalTokenChanged(ctx context.Context, tkn *domain.Token) error {
	l := t.l().C(ctx).Mth("on-tkn-changed-loc").F(kit.KV{"tknId": tkn.Id}).Dbg()

	localPlatform, err := t.localPlatform.Get(ctx)
	if err != nil {
		return err
	}

	// check token is of the local platform
	stored, err := t.tokenService.GetToken(ctx, tkn.Id)
	if err != nil {
		return err
	}
	if stored != nil && stored.PlatformId != t.localPlatform.GetPlatformId(ctx) {
		return errors.ErrTrfNotBelongLocalPlatform(ctx)
	}

	// merge token to local platform
	tkn, err = t.tokenService.PutToken(ctx, tkn)
	if err != nil {
		return err
	}

	// no changes applied
	if tkn == nil {
		l.Warn("no changes applied")
		return nil
	}

	// get platforms to push token
	platforms, err := t.getPlatformsToPush(ctx, tkn.PlatformId)
	if err != nil {
		return err
	}

	l.DbgF("%d platforms to push", len(platforms))

	// for each platform
	ocpiTkn := t.converter.TokenDomainToModel(tkn)
	for _, platform := range platforms {
		platform := platform
		// check if token receiver is supported by the remote platform
		ep := t.platformService.RoleEndpoint(ctx, platform, model.ModuleIdTokens, model.OcpiReceiver)
		if ep != "" && (platform.Protocol == nil || platform.Protocol.PushSupport.Tokens) {
			// push token to a remote platform
			rq := buildOcpiRepositoryErrHandlerRequestG(ep, t.tokenC(platform), localPlatform, platform, ocpiTkn, l)
			t.remoteTokenRep.PutTokenAsync(ctx, rq)

		} else {
			l.F(kit.KV{"platform": platform.Id}).Dbg("push not supported")
		}
	}

	return nil
}

func (t *tokenUc) OnRemoteTokensPull(ctx context.Context, from, to *time.Time) error {
	// get platforms to pull from
	platforms, err := t.getPlatformsToPull(ctx)
	if err != nil {
		return err
	}

	return t.remoteTokensPull(ctx, from, to, platforms)
}

func (t *tokenUc) OnRemoteTokensPullWhenPushNotSupported(ctx context.Context, from, to *time.Time) error {
	// get platforms to pull from
	platforms, err := t.getPlatformsToPull(ctx)
	if err != nil {
		return err
	}

	// retrieve only those which don't support tokens push
	platforms = kit.Filter(platforms, func(p *domain.Platform) bool { return p.Protocol != nil && !p.Protocol.PushSupport.Tokens })

	return t.remoteTokensPull(ctx, from, to, platforms)
}

func (t *tokenUc) OnRemoteTokenPut(ctx context.Context, platformId string, tkn *model.OcpiToken) error {
	t.l().C(ctx).Mth("on-tkn-put-rem").F(kit.KV{"platformId": platformId, "tknId": tkn.Id}).Dbg()
	return t.modifyToken(ctx, platformId, tkn, t.tokenService.PutToken)
}

func (t *tokenUc) OnRemoteTokenPatch(ctx context.Context, platformId string, tkn *model.OcpiToken) error {
	t.l().C(ctx).Mth("on-tkn-patch-rem").F(kit.KV{"platformId": platformId, "tknId": tkn.Id}).Dbg()
	return t.modifyToken(ctx, platformId, tkn, t.tokenService.MergeToken)
}

func (t *tokenUc) GetOrCreateLocalToken(ctx context.Context, tkn *domain.Token) (*domain.Token, error) {
	t.l().C(ctx).Mth("get-create-tkn-loc").F(kit.KV{"tknId": tkn.Id}).Dbg()

	// try to get valid stored token
	stored, err := t.tokenService.GetToken(ctx, tkn.Id)
	if err != nil {
		return nil, err
	}
	if stored != nil && stored.Details.Valid != nil && !*stored.Details.Valid {
		return nil, errors.ErrTknNotValid(ctx)
	}

	// otherwise create a new APP_USER token with default attr
	return t.tokenService.PutToken(ctx, &domain.Token{
		OcpiItem: domain.OcpiItem{
			ExtId:       tkn.ExtId,
			PlatformId:  tkn.PlatformId,
			RefId:       tkn.RefId,
			LastUpdated: kit.Now(),
		},
		Id: tkn.Id,
		Details: domain.TokenDetails{
			Type:         domain.TokenTypeAppUser,
			ContractId:   "unknown",
			VisualNumber: tkn.Details.VisualNumber,
			Issuer:       "unknown",
			Valid:        kit.BoolPtr(true),
			WhiteList:    domain.TokenWLTypeNever,
		},
	})
}

func (t *tokenUc) getPlatformsToPush(ctx context.Context, originalPlatformId string) ([]*domain.Platform, error) {
	platforms, err := t.platformService.Search(ctx, &domain.PlatformSearchCriteria{
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

func (t *tokenUc) getPlatformsToPull(ctx context.Context) ([]*domain.Platform, error) {
	platforms, err := t.platformService.Search(ctx, &domain.PlatformSearchCriteria{
		Statuses: []string{domain.ConnectionStatusConnected}, // connected platforms
		Remote:   kit.BoolPtr(true),                          // remote platforms
		ExcRoles: []string{domain.RoleEMSP},                  // don't pull from EMSP platforms
	})
	if err != nil {
		return nil, err
	}
	return platforms, nil
}

func (t *tokenUc) modifyToken(ctx context.Context, platformId string, tkn *model.OcpiToken, modifyFunc func(context.Context, *domain.Token) (*domain.Token, error)) error {
	l := t.l().C(ctx).Mth("tkn-modify").F(kit.KV{"platformId": platformId, "tknId": tkn.Id}).Dbg()

	// get and check platform
	_, err := t.getConnectedPlatform(ctx, platformId)
	if err != nil {
		return err
	}

	// check location is of the remote platform
	stored, err := t.tokenService.GetToken(ctx, tkn.Id)
	if err != nil {
		return err
	}
	if stored != nil && stored.PlatformId == t.localPlatform.GetPlatformId(ctx) {
		return errors.ErrTrfNotBelongRemotePlatform(ctx)
	}

	// get or create party
	_, err = t.getCreateParty(ctx, platformId, tkn.PartyId, tkn.CountryCode)
	if err != nil {
		return err
	}

	// merge token to local platform
	tknDom, err := modifyFunc(ctx, t.converter.TokenModelToDomain(tkn, platformId))
	if err != nil {
		return err
	}

	// no changes applied
	if tknDom == nil {
		l.Warn("no changes applied")
		return nil
	}

	// call webhook
	err = t.webhook.OnTokensChanged(ctx, t.converter.TokenDomainToBackend(tknDom))
	return err
}

func (t *tokenUc) remoteTokensPull(ctx context.Context, from, to *time.Time, platforms []*domain.Platform) error {
	l := t.l().C(ctx).Mth("remote-pull").Dbg()

	if len(platforms) == 0 {
		return nil
	}

	localPlatform, err := t.localPlatform.Get(ctx)
	if err != nil {
		return err
	}

	ctx = kit.Copy(ctx)

	// execute workers
	ch := make(chan channelData, tknWorkersNum)
	for i := 0; i < tknWorkersNum; i++ {
		goroutine.New().WithLogger(l).Go(ctx, func() {
			for v := range ch {
				tkn := v.data.(*model.OcpiToken)
				err := t.OnRemoteTokenPut(ctx, v.platformId, tkn)
				if err != nil {
					t.l().C(ctx).Mth("remote-put").
						F(kit.KV{"platformId": v.platformId, "tknId": tkn.Id}).
						E(err).St().Err()
				}
			}
		})
	}

	// for each platform execute pagers and put retrieved items to worker routines
	pgReader := NewPageReader(t.remoteTokenRep.GetTokens)
	eg := goroutine.NewGroup(ctx).WithLogger(l)
	for _, platform := range platforms {
		platform := platform
		// check if sender isn't supported by the remote platform
		ep := t.platformService.RoleEndpoint(ctx, platform, model.ModuleIdTokens, model.OcpiSender)
		if ep == "" {
			continue
		}
		eg.Go(func() error {
			for tokens := range pgReader.GetPage(ctx, buildOcpiRepositoryRequest(ep, t.tokenC(platform), localPlatform, platform), tknPageSize, from, to) {
				for _, tkn := range tokens {
					ch <- channelData{platformId: platform.Id, data: tkn}
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
