package impl

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/errors"
)

type platformService struct {
	storage      domain.PlatformStorage
	partyService domain.PartyService
	tokenGen     domain.TokenGenerator
	cfg          *ocpi.CfgOcpiConfig
}

func NewPlatformService(storage domain.PlatformStorage, tokenGen domain.TokenGenerator, partyService domain.PartyService) domain.PlatformService {
	return &platformService{
		storage:      storage,
		tokenGen:     tokenGen,
		partyService: partyService,
	}
}

func (p *platformService) l() kit.CLogger {
	return ocpi.L().Cmp("platform-svc")
}

func (p *platformService) Init(ctx context.Context, cfg *ocpi.CfgOcpiConfig) error {
	p.cfg = cfg
	return nil
}

func (p *platformService) Merge(ctx context.Context, platform *domain.Platform) (*domain.Platform, error) {
	p.l().C(ctx).Mth("merge").Dbg()

	if platform.Id == "" {
		return nil, errors.ErrPlatformIdEmpty(ctx)
	}

	// prepare a new platform
	if platform.Status == "" {
		platform.Status = domain.ConnectionStatusPlanned
	}

	// request stored platform
	stored, err := p.storage.GetPlatform(ctx, platform.Id)
	if err != nil {
		return nil, err
	}

	if stored != nil {
		platform.Status = stored.Status
		platform.Remote = stored.Remote
		if platform.VersionInfo.Current == "" {
			platform.VersionInfo.Current = stored.VersionInfo.Current
		}
		if platform.VersionInfo.Available == nil {
			platform.VersionInfo.Available = stored.VersionInfo.Available
		}
		if platform.VersionInfo.VersionEp == "" {
			platform.VersionInfo.VersionEp = stored.VersionInfo.VersionEp
		}
		if platform.Role == "" {
			platform.Role = stored.Role
		}
		if platform.Name == "" {
			platform.Name = stored.Name
		}
		if platform.TokenA == "" {
			platform.TokenA = stored.TokenA
		}
		if platform.TokenB == "" {
			platform.TokenB = stored.TokenB
		}
		if platform.TokenC == "" {
			platform.TokenC = stored.TokenC
		}
		if platform.TokenC == "" {
			platform.TokenC = stored.TokenC
		}
		if platform.Endpoints == nil {
			platform.Endpoints = stored.Endpoints
		}
		if platform.TokenBase64 == nil {
			platform.TokenBase64 = stored.TokenBase64
		}
		if platform.Protocol == nil {
			platform.Protocol = stored.Protocol
		}
	}

	if platform.Protocol == nil {
		platform.Protocol = &domain.ProtocolDetails{}
	}

	// validation
	err = p.validatePlatform(ctx, platform)
	if err != nil {
		return nil, err
	}

	// update storage
	if stored != nil {
		err = p.storage.UpdatePlatform(ctx, platform)
	} else {
		err = p.storage.CreatePlatform(ctx, platform)
	}
	if err != nil {
		return nil, err
	}

	return platform, nil
}

func (p *platformService) SetStatus(ctx context.Context, platformId string, status string) (*domain.Platform, error) {
	p.l().C(ctx).Mth("status").Dbg()

	// get stored
	platform, err := p.mustGet(ctx, platformId)
	if err != nil {
		return nil, err
	}

	platform.Status = status

	// validation
	err = p.validatePlatform(ctx, platform)
	if err != nil {
		return nil, err
	}

	// if connected status
	if platform.Status == domain.ConnectionStatusConnected {
		// for remote platform tokens must be populated
		if platform.Remote &&
			(platform.TokenB == "" || platform.TokenC == "") {
			return nil, errors.ErrNoTokensSpecifiedForConnection(ctx)
		}
		// list of versions and currently used version
		if len(platform.VersionInfo.Current) == 0 {
			return nil, errors.ErrNoVersionsSpecifiedForConnection(ctx)
		}
		if _, ok := platform.VersionInfo.Available[platform.VersionInfo.Current]; !ok {
			return nil, errors.ErrCurrentVersionInvalid(ctx)
		}
	}

	// put to storage
	err = p.storage.UpdatePlatform(ctx, platform)
	if err != nil {
		return nil, err
	}

	return platform, nil
}

func (p *platformService) Get(ctx context.Context, platformId string) (*domain.Platform, error) {
	p.l().C(ctx).Mth("get").Dbg()
	if platformId == "" {
		return nil, errors.ErrPlatformIdEmpty(ctx)
	}
	return p.storage.GetPlatform(ctx, platformId)
}

func (p *platformService) GetByTokenA(ctx context.Context, token domain.PlatformToken) (*domain.Platform, error) {
	p.l().C(ctx).Mth("get-by-token-a").Dbg()
	if token == "" {
		return nil, errors.ErrPlatformTokenAEmpty(ctx)
	}
	return p.storage.GetPlatformByTokenA(ctx, p.tryBase64Token(token))
}

func (p *platformService) GetByTokenB(ctx context.Context, token domain.PlatformToken) (*domain.Platform, error) {
	p.l().C(ctx).Mth("get-by-token-b").Dbg()
	if token == "" {
		return nil, errors.ErrPlatformTokenEmpty(ctx)
	}
	return p.storage.GetPlatformByTokenB(ctx, p.tryBase64Token(token))
}

func (p *platformService) GetByTokenC(ctx context.Context, token domain.PlatformToken) (*domain.Platform, error) {
	p.l().C(ctx).Mth("get-by-token-с").Dbg()
	if token == "" {
		return nil, errors.ErrPlatformTokenEmpty(ctx)
	}
	return p.storage.GetPlatformByTokenC(ctx, p.tryBase64Token(token))
}

func (p *platformService) Search(ctx context.Context, cr *domain.PlatformSearchCriteria) ([]*domain.Platform, error) {
	p.l().C(ctx).Mth("search").Dbg()
	return p.storage.SearchPlatforms(ctx, cr)
}

func (p *platformService) RoleEndpoint(ctx context.Context, platform *domain.Platform, module, role string) domain.Endpoint {
	// check if push location is supported
	locModule, ok := platform.Endpoints[module]
	if !ok {
		// module not supported
		return ""
	}
	ep, ok := locModule[role]
	if !ok {
		// receiver interface isn't supported
		return ""
	}
	return ep
}

// tryBase64Token tries to decode from base64, otherwise returns a source token
func (p *platformService) tryBase64Token(token domain.PlatformToken) domain.PlatformToken {
	decoded, ok := p.tokenGen.TryBase64Decode(token)
	if ok {
		return decoded
	}
	return token
}

func (p *platformService) validatePlatform(ctx context.Context, platform *domain.Platform) error {

	if platform == nil {
		return errors.ErrPlatformEmpty(ctx)
	}
	if platform.Name == "" {
		return errors.ErrPlatformNameEmpty(ctx)
	}
	if platform.TokenA == "" {
		return errors.ErrPlatformTokenAEmpty(ctx)
	}
	if platform.VersionInfo.VersionEp == "" {
		return errors.ErrPlatformVersionEpEmpty(ctx)
	}

	// url
	if !kit.IsUrlValid(string(platform.VersionInfo.VersionEp)) {
		return errors.ErrPlatformVersionEpNotValid(ctx)
	}

	// role
	supported, ok := domain.RoleMap[platform.Role]
	if !ok {
		return errors.ErrPlatformRoleInvalid(ctx)
	}
	if !supported {
		return errors.ErrPlatformRoleNotSupported(ctx)
	}

	// if connected status
	if platform.Status == domain.ConnectionStatusConnected {
		// for remote platform tokens must be populated
		if platform.Remote &&
			(platform.TokenB == "" || platform.TokenC == "") {
			return errors.ErrNoTokensSpecifiedForConnection(ctx)
		}
		// list of versions and currently used version
		if len(platform.VersionInfo.Current) == 0 {
			return errors.ErrNoVersionsSpecifiedForConnection(ctx)
		}
		if _, ok := platform.VersionInfo.Available[platform.VersionInfo.Current]; !ok {
			return errors.ErrCurrentVersionInvalid(ctx)
		}
	}
	return nil
}

func (p *platformService) mustGet(ctx context.Context, platformId string) (*domain.Platform, error) {
	if platformId == "" {
		return nil, errors.ErrPlatformIdEmpty(ctx)
	}
	// get stored
	platform, err := p.storage.GetPlatform(ctx, platformId)
	if err != nil {
		return nil, err
	}
	if platform == nil {
		return nil, errors.ErrPlatformNotFound(ctx, platformId)
	}
	return platform, nil
}
