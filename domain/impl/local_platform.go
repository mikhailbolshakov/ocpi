package impl

import (
	"context"
	"fmt"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/errors"
	"github.com/mikhailbolshakov/ocpi/model"
	"strings"
)

type localPlatformService struct {
	platformService domain.PlatformService
	partyService    domain.PartyService
	cfg             *ocpi.CfgOcpiConfig
}

func NewLocalPlatformService(platformService domain.PlatformService, partyService domain.PartyService) domain.LocalPlatformService {
	return &localPlatformService{
		platformService: platformService,
		partyService:    partyService,
	}
}

func (p *localPlatformService) l() kit.CLogger {
	return ocpi.L().Cmp("local-platform-svc")
}

var (
	Modules = map[string]struct {
		SenderOnly   bool
		ReceiverOnly bool
		NotSupported bool
	}{
		domain.ModuleIdCredentials:   {SenderOnly: true},
		domain.ModuleIdCdrs:          {},
		domain.ModuleIdCommands:      {},
		domain.ModuleIdHubClientInfo: {},
		domain.ModuleIdLocations:     {},
		domain.ModuleIdSessions:      {},
		domain.ModuleIdTariffs:       {},
		domain.ModuleIdTokens:        {},
	}
)

func (p *localPlatformService) Init(ctx context.Context, cfg *ocpi.CfgOcpiConfig) error {
	p.cfg = cfg
	return nil
}

func (p *localPlatformService) InitializePlatform(ctx context.Context) error {
	p.l().C(ctx).Mth("init").Dbg()
	// get home platform
	platform, err := p.platformService.Get(ctx, p.cfg.Local.Platform.Id)
	if err != nil {
		return err
	}
	found := platform != nil
	if !found {
		platform = &domain.Platform{}
	}
	platform.Id = p.cfg.Local.Platform.Id
	platform.Role = p.cfg.Local.Platform.Role
	platform.Name = p.cfg.Local.Platform.Name
	platform.TokenA = domain.PlatformToken(p.cfg.Local.Platform.TokenA)
	platform.Status = domain.ConnectionStatusConnected
	platform.VersionInfo.Available = make(domain.Versions)
	platform.TokenBase64 = kit.BoolPtr(false)

	platform.Protocol = &domain.ProtocolDetails{}
	platform.Protocol.PushSupport.Tokens = true
	platform.Protocol.PushSupport.Cdrs = true
	platform.Protocol.PushSupport.Commands = true
	platform.Protocol.PushSupport.Credentials = true
	platform.Protocol.PushSupport.HubClientInfo = true
	platform.Protocol.PushSupport.Locations = true
	platform.Protocol.PushSupport.Sessions = true
	platform.Protocol.PushSupport.Tariffs = true

	for _, v := range p.cfg.Local.Platform.Versions {
		platform.VersionInfo.Available[v] = domain.Endpoint(fmt.Sprintf("%s/ocpi/%s", p.cfg.Local.Url, v))
		platform.VersionInfo.Current = v
	}
	platform.VersionInfo.VersionEp = domain.Endpoint(fmt.Sprintf("%s/ocpi/versions", p.cfg.Local.Url))
	platform.Endpoints = p.GetEndpoints(ctx, platform.VersionInfo.Current)

	// create or update platform
	_, err = p.platformService.Merge(ctx, platform)
	if err != nil {
		return err
	}

	// merge party
	party := &domain.Party{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     p.cfg.Local.Party.PartyId,
				CountryCode: p.cfg.Local.Party.CountryCode,
			},
			PlatformId:  p.cfg.Local.Platform.Id,
			LastUpdated: kit.Now(),
		},
		Roles: strings.Split(p.cfg.Local.Party.Roles, ","),
		BusinessDetails: &domain.BusinessDetails{
			Name: p.cfg.Local.Platform.Name,
		},
		Status: domain.ConnectionStatusConnected,
	}
	_, err = p.partyService.Merge(ctx, party)
	if err != nil {
		return err
	}

	return nil
}

func (p *localPlatformService) Get(ctx context.Context) (*domain.Platform, error) {
	p.l().C(ctx).Mth("get").Dbg()
	// get home platform
	platform, err := p.platformService.Get(ctx, p.cfg.Local.Platform.Id)
	if err != nil {
		return nil, err
	}
	if platform == nil {
		return nil, errors.ErrPlatformNotFound(ctx, p.cfg.Local.Platform.Id)
	}
	return platform, nil
}

func (p *localPlatformService) GetPlatformId(ctx context.Context) string {
	return p.cfg.Local.Platform.Id
}

func (p *localPlatformService) GetEndpoints(ctx context.Context, version string) domain.ModuleEndpoints {
	p.l().C(ctx).Mth("get-endpoints").Dbg()
	r := make(domain.ModuleEndpoints)
	for moduleId, d := range Modules {
		if d.NotSupported {
			continue
		}
		if d.SenderOnly {
			r[moduleId] = domain.RoleEndpoint{
				model.OcpiSender: domain.Endpoint(fmt.Sprintf("%s/ocpi/%s/%s", p.cfg.Local.Url, version, moduleId)),
			}
		} else if d.ReceiverOnly {
			r[moduleId] = domain.RoleEndpoint{
				model.OcpiReceiver: domain.Endpoint(fmt.Sprintf("%s/ocpi/%s/receiver/%s", p.cfg.Local.Url, version, moduleId)),
			}
		} else {
			r[moduleId] = domain.RoleEndpoint{
				model.OcpiSender:   domain.Endpoint(fmt.Sprintf("%s/ocpi/%s/sender/%s", p.cfg.Local.Url, version, moduleId)),
				model.OcpiReceiver: domain.Endpoint(fmt.Sprintf("%s/ocpi/%s/receiver/%s", p.cfg.Local.Url, version, moduleId)),
			}
		}
	}
	return r
}
