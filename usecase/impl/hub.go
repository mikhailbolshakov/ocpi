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
	ciWorkersNum = 16
	ciPageSize   = 100
)

type hubUc struct {
	ucBase
	localPlatform       domain.LocalPlatformService
	remoteClientInfoRep usecase.RemoteHubClientInfoRepository
	partyService        domain.PartyService
	webhook             backend.WebhookCallService
}

func NewHubUc(platformService domain.PlatformService, remoteClientInfoRep usecase.RemoteHubClientInfoRepository, partyService domain.PartyService,
	webhook backend.WebhookCallService, localPlatform domain.LocalPlatformService, tokenGen domain.TokenGenerator) usecase.HubUc {
	return &hubUc{
		ucBase:              newBase(platformService, partyService, tokenGen),
		partyService:        partyService,
		remoteClientInfoRep: remoteClientInfoRep,
		webhook:             webhook,
		localPlatform:       localPlatform,
	}
}

func (h *hubUc) l() kit.CLogger {
	return ocpi.L().Cmp("hub-uc")
}

func (h *hubUc) OnRemoteClientInfosPull(ctx context.Context, from, to *time.Time) error {
	// get platforms to pull from
	platforms, err := h.getPlatformsToPull(ctx)
	if err != nil {
		return err
	}

	return h.remoteClientInfosPull(ctx, from, to, platforms)
}

func (h *hubUc) OnRemoteClientInfosPullWhenPushNotSupported(ctx context.Context, from, to *time.Time) error {
	// get platforms to pull from
	platforms, err := h.getPlatformsToPull(ctx)
	if err != nil {
		return err
	}

	// retrieve only those which don't support client info push
	platforms = kit.Filter(platforms, func(p *domain.Platform) bool { return p.Protocol != nil && !p.Protocol.PushSupport.HubClientInfo })

	return h.remoteClientInfosPull(ctx, from, to, platforms)
}

func (h *hubUc) OnRemoteClientInfoPut(ctx context.Context, platformId string, ci *model.OcpiClientInfo) error {
	l := h.l().C(ctx).Mth("on-ci-changed-rem").F(kit.KV{"partyId": ci.PartyId}).Dbg()

	// get and check platform
	_, err := h.getConnectedPlatform(ctx, platformId)
	if err != nil {
		return err
	}

	// check party is of the remote platform
	stored, err := h.partyService.GetByExtId(ctx, domain.PartyExtId{
		PartyId:     ci.OcpiPartyId.PartyId,
		CountryCode: ci.OcpiPartyId.CountryCode,
	})
	if err != nil {
		return err
	}
	if stored != nil && stored.PlatformId == h.localPlatform.GetPlatformId(ctx) {
		return errors.ErrPartyNotBelongRemotePlatform(ctx)
	}

	// merge location to local platform
	ciDom, err := h.partyService.Merge(ctx, h.partyModelToDomain(platformId, ci, stored))
	if err != nil {
		return err
	}

	// no changes applied
	if ciDom == nil {
		l.Warn("no changes applied")
		return nil
	}

	// call webhook
	err = h.webhook.OnPartiesChanged(ctx, h.partyDomainToBackend(ciDom))
	if err != nil {
		return err
	}

	return err
}

func (h *hubUc) OnLocalClientInfoChanged(ctx context.Context, party *domain.Party) error {
	l := h.l().C(ctx).Mth("on-ci-changed-loc").Dbg()

	localPlatform, err := h.localPlatform.Get(ctx)
	if err != nil {
		return err
	}

	// check location is of the local platform
	stored, err := h.partyService.Get(ctx, party.Id)
	if err != nil {
		return err
	}
	if stored != nil && stored.PlatformId != h.localPlatform.GetPlatformId(ctx) {
		return errors.ErrPartyNotBelongLocalPlatform(ctx)
	}

	// get remote platforms to push
	remotePlatforms, err := h.getPlatformsToPush(ctx, party.PlatformId)
	if err != nil {
		return err
	}

	// for each platform
	cis := h.partyDomainToModel(party)
	for _, platform := range remotePlatforms {
		// check if client info receiver is supported by the remote platform
		ep := h.platformService.RoleEndpoint(ctx, platform, model.ModuleIdHubClientInfo, model.OcpiReceiver)
		if ep != "" && (platform.Protocol == nil || platform.Protocol.PushSupport.HubClientInfo) {
			// push client info to the remote platform
			for _, ci := range cis {
				if err := h.remoteClientInfoRep.PutClientInfo(ctx, buildOcpiRepositoryRequestG(ep, h.tokenC(platform), localPlatform, platform, ci)); err != nil {
					l.F(kit.KV{"platform": platform.Id}).E(err).St().Err()
				}
			}
		}
	}

	return nil
}

func (h *hubUc) getPlatformsToPush(ctx context.Context, partyPlatformId string) ([]*domain.Platform, error) {
	platforms, err := h.platformService.Search(ctx, &domain.PlatformSearchCriteria{
		Statuses: []string{domain.ConnectionStatusConnected}, // connected platforms
		ExcIds:   []string{partyPlatformId},                  // exclude source platform
		Remote:   kit.BoolPtr(true),                          // remote platforms
	})
	if err != nil {
		return nil, err
	}
	return platforms, nil
}

func (h *hubUc) getPlatformsToPull(ctx context.Context) ([]*domain.Platform, error) {
	platforms, err := h.platformService.Search(ctx, &domain.PlatformSearchCriteria{
		Statuses: []string{domain.ConnectionStatusConnected}, // connected platforms
		Remote:   kit.BoolPtr(true),                          // remote platforms
	})
	if err != nil {
		return nil, err
	}
	return platforms, nil
}

func (h *hubUc) remoteClientInfosPull(ctx context.Context, from, to *time.Time, platforms []*domain.Platform) error {
	l := h.l().C(ctx).Mth("remote-pull").Dbg()

	if len(platforms) == 0 {
		return nil
	}

	localPlatform, err := h.localPlatform.Get(ctx)
	if err != nil {
		return err
	}

	ctx = kit.Copy(ctx)

	// execute workers
	ch := make(chan channelData, ciWorkersNum)
	for i := 0; i < ciWorkersNum; i++ {
		goroutine.New().WithLogger(l).Go(ctx, func() {
			for v := range ch {
				ci := v.data.(*model.OcpiClientInfo)
				err := h.OnRemoteClientInfoPut(ctx, v.platformId, ci)
				if err != nil {
					h.l().C(ctx).Mth("remote-put").
						F(kit.KV{"platformId": v.platformId, "partyId": ci.PartyId}).
						E(err).St().Err()
				}
			}
		})
	}

	// for each platform execute pagers and put retrieved items to worker routines
	pgReader := NewPageReader(h.remoteClientInfoRep.GetClientInfos)
	eg := goroutine.NewGroup(ctx).WithLogger(l)
	for _, platform := range platforms {
		platform := platform
		// check if sender isn't supported by the remote platform
		ep := h.platformService.RoleEndpoint(ctx, platform, model.ModuleIdHubClientInfo, model.OcpiSender)
		if ep == "" {
			continue
		}
		eg.Go(func() error {
			for cis := range pgReader.GetPage(ctx, buildOcpiRepositoryRequest(ep, h.tokenC(platform), localPlatform, platform), ciPageSize, from, to) {
				for _, ci := range cis {
					ch <- channelData{platformId: platform.Id, data: ci}
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
