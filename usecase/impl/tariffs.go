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
	trfWorkersNum = 4
	trfPageSize   = 100
)

type tariffUc struct {
	ucBase
	localPlatform   domain.LocalPlatformService
	tariffService   domain.TariffService
	remoteTariffRep usecase.RemoteTariffRepository
	partyService    domain.PartyService
	webhook         backend.WebhookCallService
	converter       usecase.TariffConverter
}

func NewTariffUc(platformService domain.PlatformService, tariffService domain.TariffService, remoteTariffRep usecase.RemoteTariffRepository,
	partyService domain.PartyService, webhook backend.WebhookCallService, localPlatform domain.LocalPlatformService, tokenGen domain.TokenGenerator) usecase.TariffUc {
	return &tariffUc{
		ucBase:          newBase(platformService, partyService, tokenGen),
		tariffService:   tariffService,
		remoteTariffRep: remoteTariffRep,
		partyService:    partyService,
		webhook:         webhook,
		converter:       NewTariffConverter(),
		localPlatform:   localPlatform,
	}
}

func (t *tariffUc) l() kit.CLogger {
	return ocpi.L().Cmp("trf-uc")
}

func (t *tariffUc) OnLocalTariffChanged(ctx context.Context, trf *domain.Tariff) error {
	l := t.l().C(ctx).Mth("on-trf-changed-loc").F(kit.KV{"locId": trf.Id}).Dbg()

	localPlatform, err := t.localPlatform.Get(ctx)
	if err != nil {
		return err
	}

	// check tariff is of the local platform
	stored, err := t.tariffService.GetTariff(ctx, trf.Id)
	if err != nil {
		return err
	}
	if stored != nil && stored.PlatformId != t.localPlatform.GetPlatformId(ctx) {
		return errors.ErrTrfNotBelongLocalPlatform(ctx)
	}

	// merge tariff to local platform
	trf, err = t.tariffService.PutTariff(ctx, trf)
	if err != nil {
		return err
	}

	// no changes applied
	if trf == nil {
		l.Warn("no changes applied")
		return nil
	}

	// get platforms to push tariff
	platforms, err := t.getPlatformsToPush(ctx, trf.PlatformId)
	if err != nil {
		return err
	}

	l.DbgF("%d platforms to push", len(platforms))

	// for each platform
	ocpiTrf := t.converter.TariffDomainToModel(trf)
	for _, platform := range platforms {
		platform := platform
		// check if tariff receiver is supported by the remote platform
		ep := t.platformService.RoleEndpoint(ctx, platform, model.ModuleIdTariffs, model.OcpiReceiver)
		if ep != "" && (platform.Protocol == nil || platform.Protocol.PushSupport.Tariffs) {
			// push tariff to a remote platform
			rq := buildOcpiRepositoryErrHandlerRequestG(ep, t.tokenC(platform), localPlatform, platform, ocpiTrf, l)
			t.remoteTariffRep.PutTariffAsync(ctx, rq)
		} else {
			l.F(kit.KV{"platform": platform.Id}).Dbg("push not supported")
		}
	}

	return nil
}

func (t *tariffUc) OnRemoteTariffPut(ctx context.Context, platformId string, trf *model.OcpiTariff) error {
	t.l().C(ctx).Mth("on-trf-put-rem").F(kit.KV{"platformId": platformId, "locId": trf.Id}).Dbg()
	return t.modifyTariff(ctx, platformId, trf, t.tariffService.PutTariff)
}

func (t *tariffUc) OnRemoteTariffPatch(ctx context.Context, platformId string, trf *model.OcpiTariff) error {
	t.l().C(ctx).Mth("on-trf-patch-rem").F(kit.KV{"platformId": platformId, "locId": trf.Id}).Dbg()
	return t.modifyTariff(ctx, platformId, trf, t.tariffService.MergeTariff)
}

func (t *tariffUc) OnRemoteTariffsPull(ctx context.Context, from, to *time.Time) error {
	// get platforms to pull from
	platforms, err := t.getPlatformsToPull(ctx)
	if err != nil {
		return err
	}

	return t.remoteTariffsPull(ctx, from, to, platforms)
}

func (t *tariffUc) OnRemoteTariffsPullWhenPushNotSupported(ctx context.Context, from, to *time.Time) error {
	// get platforms to pull from
	platforms, err := t.getPlatformsToPull(ctx)
	if err != nil {
		return err
	}

	// retrieve only those which don't support tariffs push
	platforms = kit.Filter(platforms, func(p *domain.Platform) bool { return p.Protocol != nil && !p.Protocol.PushSupport.Tariffs })

	return t.remoteTariffsPull(ctx, from, to, platforms)
}

func (t *tariffUc) getPlatformsToPush(ctx context.Context, originalPlatformId string) ([]*domain.Platform, error) {
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

func (t *tariffUc) getPlatformsToPull(ctx context.Context) ([]*domain.Platform, error) {
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

func (t *tariffUc) modifyTariff(ctx context.Context, platformId string, trf *model.OcpiTariff, modifyFunc func(context.Context, *domain.Tariff) (*domain.Tariff, error)) error {
	l := t.l().C(ctx).Mth("trf-modify").F(kit.KV{"platformId": platformId, "trfId": trf.Id}).Dbg()

	// get and check platform
	_, err := t.getConnectedPlatform(ctx, platformId)
	if err != nil {
		return err
	}

	// check location is of the remote platform
	stored, err := t.tariffService.GetTariff(ctx, trf.Id)
	if err != nil {
		return err
	}
	if stored != nil && stored.PlatformId == t.localPlatform.GetPlatformId(ctx) {
		return errors.ErrTrfNotBelongRemotePlatform(ctx)
	}

	// get or create party
	_, err = t.getCreateParty(ctx, platformId, trf.PartyId, trf.CountryCode)
	if err != nil {
		return err
	}

	// merge tariff to local platform
	trfDom, err := modifyFunc(ctx, t.converter.TariffModelToDomain(trf, platformId))
	if err != nil {
		return err
	}

	// no changes applied
	if trfDom == nil {
		l.Warn("no changes applied")
		return nil
	}

	// call webhook
	err = t.webhook.OnTariffsChanged(ctx, t.converter.TariffDomainToBackend(trfDom))
	return err
}

func (t *tariffUc) remoteTariffsPull(ctx context.Context, from, to *time.Time, platforms []*domain.Platform) error {
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
	ch := make(chan channelData, trfWorkersNum)
	for i := 0; i < trfWorkersNum; i++ {
		goroutine.New().WithLogger(l).Go(ctx, func() {
			for v := range ch {
				trf := v.data.(*model.OcpiTariff)
				err := t.OnRemoteTariffPut(ctx, v.platformId, trf)
				if err != nil {
					t.l().C(ctx).Mth("remote-put").
						F(kit.KV{"platformId": v.platformId, "trfId": trf.Id}).
						E(err).St().Err()
				}
			}
		})
	}

	// for each platform execute pagers and put retrieved items to worker routines
	pgReader := NewPageReader(t.remoteTariffRep.GetTariffs)
	eg := goroutine.NewGroup(ctx).WithLogger(l)
	for _, platform := range platforms {
		platform := platform
		// check if sender isn't supported by the remote platform
		ep := t.platformService.RoleEndpoint(ctx, platform, model.ModuleIdTariffs, model.OcpiSender)
		if ep == "" {
			continue
		}
		eg.Go(func() error {
			for tariffs := range pgReader.GetPage(ctx, buildOcpiRepositoryRequest(ep, t.tokenC(platform), localPlatform, platform), trfPageSize, from, to) {
				for _, trf := range tariffs {
					ch <- channelData{platformId: platform.Id, data: trf}
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
