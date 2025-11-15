package impl

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/goroutine"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/errors"
	"github.com/mikhailbolshakov/ocpi/usecase"
)

type maintenanceUc struct {
	ucBase
	platformService      domain.PlatformService
	localPlatformService domain.LocalPlatformService
	partyService         domain.PartyService
	locService           domain.LocationService
	cmdService           domain.CommandService
	sessService          domain.SessionService
	cdrsService          domain.CdrService
	trfService           domain.TariffService
	tknService           domain.TokenService
}

func NewMaintenanceUc(platformService domain.PlatformService,
	localPlatformService domain.LocalPlatformService,
	partyService domain.PartyService,
	locService domain.LocationService,
	cmdService domain.CommandService,
	sessService domain.SessionService,
	cdrsService domain.CdrService,
	trfService domain.TariffService,
	tknService domain.TokenService,
	tokenGen domain.TokenGenerator,
) usecase.MaintenanceUc {
	return &maintenanceUc{
		ucBase:               newBase(platformService, partyService, tokenGen),
		platformService:      platformService,
		localPlatformService: localPlatformService,
		partyService:         partyService,
		locService:           locService,
		cmdService:           cmdService,
		sessService:          sessService,
		cdrsService:          cdrsService,
		trfService:           trfService,
		tknService:           tknService,
	}
}

func (s *maintenanceUc) l() kit.CLogger {
	return ocpi.L().Cmp("maintenance-uc")
}

func (s *maintenanceUc) DeleteLocalPartyByExt(ctx context.Context, extId domain.PartyExtId) error {
	l := s.l().C(ctx).Mth("delete-party").F(kit.KV{"partyId": extId.PartyId, "country": extId.CountryCode}).Dbg()

	// get local platform
	localPlatform, err := s.localPlatformService.Get(ctx)
	if err != nil {
		return err
	}

	// get by ext
	stored, err := s.partyService.GetByExtId(ctx, extId)
	if err != nil {
		return err
	}
	if stored == nil {
		return errors.ErrPartyNotFoundByExt(ctx)
	}
	if stored != nil && stored.PlatformId != localPlatform.Id {
		return errors.ErrPartyNotBelongLocalPlatform(ctx)
	}

	eg := goroutine.NewGroup(ctx).WithLogger(l)

	eg.Go(func() error {
		return s.partyService.DeletePartyByExtId(ctx, extId)
	})
	eg.Go(func() error {
		return s.locService.DeleteLocationsByExtId(ctx, extId)
	})
	eg.Go(func() error {
		return s.trfService.DeleteTariffsByExtId(ctx, extId)
	})
	eg.Go(func() error {
		return s.tknService.DeleteTokensByExtId(ctx, extId)
	})
	eg.Go(func() error {
		return s.cmdService.DeleteCommandsByExt(ctx, extId)
	})
	eg.Go(func() error {
		return s.sessService.DeleteSessionsByExtId(ctx, extId)
	})
	eg.Go(func() error {
		return s.cdrsService.DeleteCdrsByExtId(ctx, extId)
	})

	return eg.Wait()
}
