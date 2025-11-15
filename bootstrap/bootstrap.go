package bootstrap

import (
	"context"
	"github.com/mikhailbolshakov/kit/cluster"
	"github.com/mikhailbolshakov/kit/cron"
	kitHttp "github.com/mikhailbolshakov/kit/http"
	"github.com/mikhailbolshakov/kit/monitoring"
	"github.com/mikhailbolshakov/kit/profile"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
	impl3 "github.com/mikhailbolshakov/ocpi/backend/impl"
	ocpiCron "github.com/mikhailbolshakov/ocpi/cron"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/domain/impl"
	ocpiRep "github.com/mikhailbolshakov/ocpi/repository/adapters/ocpi"
	"github.com/mikhailbolshakov/ocpi/repository/adapters/webhook"
	"github.com/mikhailbolshakov/ocpi/repository/storage"
	"github.com/mikhailbolshakov/ocpi/transport/grpc"
	"github.com/mikhailbolshakov/ocpi/transport/http"
	bkndCdrs "github.com/mikhailbolshakov/ocpi/transport/http/backend/cdrs"
	bkndCmd "github.com/mikhailbolshakov/ocpi/transport/http/backend/commands"
	bkndLoc "github.com/mikhailbolshakov/ocpi/transport/http/backend/locations"
	bkndMnt "github.com/mikhailbolshakov/ocpi/transport/http/backend/maintenance"
	bkndParty "github.com/mikhailbolshakov/ocpi/transport/http/backend/party"
	bkndPlatform "github.com/mikhailbolshakov/ocpi/transport/http/backend/platform"
	bkndSess "github.com/mikhailbolshakov/ocpi/transport/http/backend/sessions"
	bkndSwg "github.com/mikhailbolshakov/ocpi/transport/http/backend/swagger"
	bkndTrf "github.com/mikhailbolshakov/ocpi/transport/http/backend/tariffs"
	bkndTkn "github.com/mikhailbolshakov/ocpi/transport/http/backend/tokens"
	bkndWebhook "github.com/mikhailbolshakov/ocpi/transport/http/backend/webhook"
	"github.com/mikhailbolshakov/ocpi/transport/http/ocpi/platform"
	"github.com/mikhailbolshakov/ocpi/transport/http/ocpi/v221/cdrs"
	"github.com/mikhailbolshakov/ocpi/transport/http/ocpi/v221/commands"
	"github.com/mikhailbolshakov/ocpi/transport/http/ocpi/v221/credentials"
	"github.com/mikhailbolshakov/ocpi/transport/http/ocpi/v221/hub"
	"github.com/mikhailbolshakov/ocpi/transport/http/ocpi/v221/locations"
	"github.com/mikhailbolshakov/ocpi/transport/http/ocpi/v221/sessions"
	"github.com/mikhailbolshakov/ocpi/transport/http/ocpi/v221/tariffs"
	"github.com/mikhailbolshakov/ocpi/transport/http/ocpi/v221/tokens"
	"github.com/mikhailbolshakov/ocpi/usecase"
	impl2 "github.com/mikhailbolshakov/ocpi/usecase/impl"
)

// ServiceImpl implements a service bootstrapping
// all dependencies between layers must be specified here
type ServiceImpl struct {
	cfg                  *ocpi.Config
	loadCfgFn            func() (*ocpi.Config, error)
	monitoring           monitoring.MetricsServer
	profile              profile.Server
	storageAdapter       storage.Adapter
	http                 *kitHttp.Server
	grpc                 *grpc.Server
	platformService      domain.PlatformService
	localPlatformService domain.LocalPlatformService
	tokenGen             domain.TokenGenerator
	logService           domain.OcpiLogService
	credentialsUc        usecase.CredentialsUc
	credConverter        usecase.CredentialsConverter
	partyService         domain.PartyService
	locationService      domain.LocationService
	locationUc           usecase.LocationUc
	locConverter         usecase.LocationConverter
	hubUc                usecase.HubUc
	ocpiAdapter          ocpiRep.Adapter
	trfService           domain.TariffService
	trfUc                usecase.TariffUc
	trfConverter         usecase.TariffConverter
	tknService           domain.TokenService
	tknUc                usecase.TokenUc
	tknConverter         usecase.TokenConverter
	sessService          domain.SessionService
	sessUc               usecase.SessionUc
	sessConverter        usecase.SessionConverter
	cmdService           domain.CommandService
	cmdUc                usecase.CommandUc
	cmdConverter         usecase.CommandConverter
	cdrService           domain.CdrService
	cdrUc                usecase.CdrUc
	cdrConverter         usecase.CdrConverter
	webhookService       backend.WebhookService
	webhookCallService   backend.WebhookCallService
	webhookAdapter       webhook.Adapter
	maintenanceUc        usecase.MaintenanceUc
	cronManager          cron.Manager
}

// New creates a new instance of the service
func New() cluster.Service {
	s := &ServiceImpl{
		loadCfgFn:   ocpi.LoadConfig,
		monitoring:  monitoring.NewMetricsServer(ocpi.LF()),
		profile:     profile.New(ocpi.LF()),
		cronManager: cron.NewManager(ocpi.LF()),
	}
	s.locConverter = impl2.NewLocationConverter()
	s.credConverter = impl2.NewCredentialsConverter()
	s.storageAdapter = storage.NewAdapter()
	s.logService = impl.NewOcpiLogService(s.storageAdapter)
	s.webhookAdapter = webhook.NewAdapter(s.logService)
	s.webhookService = impl3.NewWebhookService(s.storageAdapter)
	s.webhookCallService = impl3.NewWebhookCallService(s.webhookService, s.webhookAdapter)
	s.tokenGen = impl.NewTokenGenerator()
	s.ocpiAdapter = ocpiRep.NewAdapter(s.logService)
	s.partyService = impl.NewPartyService(s.storageAdapter)
	s.platformService = impl.NewPlatformService(s.storageAdapter, s.tokenGen, s.partyService)
	s.localPlatformService = impl.NewLocalPlatformService(s.platformService, s.partyService)
	s.hubUc = impl2.NewHubUc(s.platformService, s.ocpiAdapter, s.partyService, s.webhookCallService, s.localPlatformService, s.tokenGen)
	s.credentialsUc = impl2.NewCredentialsUc(s.platformService, s.localPlatformService, s.tokenGen, s.ocpiAdapter, s.partyService, s.webhookCallService, s.hubUc)
	s.locationService = impl.NewLocationService(s.storageAdapter)
	s.locationUc = impl2.NewLocationUc(s.platformService, s.locationService, s.ocpiAdapter, s.partyService, s.webhookCallService, s.localPlatformService, s.tokenGen)
	s.trfConverter = impl2.NewTariffConverter()
	s.trfService = impl.NewTariffService(s.storageAdapter)
	s.trfUc = impl2.NewTariffUc(s.platformService, s.trfService, s.ocpiAdapter, s.partyService, s.webhookCallService, s.localPlatformService, s.tokenGen)
	s.tknConverter = impl2.NewTokenConverter()
	s.tknService = impl.NewTokenService(s.storageAdapter)
	s.tknUc = impl2.NewTokenUc(s.platformService, s.tknService, s.ocpiAdapter, s.partyService, s.webhookCallService, s.localPlatformService, s.tokenGen)
	s.cmdService = impl.NewCmdService(s.tknService, s.storageAdapter)
	s.sessConverter = impl2.NewSessionConverter()
	s.sessService = impl.NewSessionService(s.storageAdapter)
	s.sessUc = impl2.NewSessionUc(s.platformService, s.sessService, s.ocpiAdapter, s.partyService, s.webhookCallService, s.cmdService, s.localPlatformService, s.tknService, s.tokenGen)
	s.cdrConverter = impl2.NewCdrConverter(s.trfConverter)
	s.cdrService = impl.NewCdrService(s.storageAdapter, s.trfService)
	s.cdrUc = impl2.NewCdrUc(s.platformService, s.cdrService, s.ocpiAdapter, s.partyService, s.webhookCallService,
		s.sessService, s.localPlatformService, s.locationService, s.trfService, s.tknService, s.tokenGen)
	s.cmdConverter = impl2.NewCommandConverter(s.tknConverter)
	s.cmdUc = impl2.NewCommandUc(s.platformService, s.cmdService, s.ocpiAdapter, s.partyService, s.locationService, s.webhookCallService,
		s.localPlatformService, s.tknUc, s.tknService, s.sessService, s.tokenGen)
	s.maintenanceUc = impl2.NewMaintenanceUc(s.platformService, s.localPlatformService, s.partyService, s.locationService, s.cmdService,
		s.sessService, s.cdrService, s.trfService, s.tknService, s.tokenGen)
	return s
}

func (s *ServiceImpl) SetConfigLoadFn(fn func() (*ocpi.Config, error)) {
	s.loadCfgFn = fn
}

func (s *ServiceImpl) GetCode() string {
	return ocpi.Meta.ServiceCode()
}

func (s *ServiceImpl) initHttpServer(ctx context.Context) error {
	// create HTTP server
	s.http = kitHttp.NewHttpServer(s.cfg.Http, ocpi.LF())

	// create and set middlewares
	mdw := http.NewMiddleware(s.platformService, s.logService, s.cfg.Ocpi)
	s.http.RootRouter.Use(mdw.SetContextMiddleware)

	// ocpi routing
	routeBuilder := http.NewRouteBuilder(s.http, mdw)
	routeBuilder.SetRoutes(platform.GetRoutes(platform.NewController(s.localPlatformService)))
	routeBuilder.SetRoutes(credentials.GetRoutes(credentials.NewController(s.credentialsUc, s.platformService, s.localPlatformService)))
	routeBuilder.SetRoutes(hub.GetRoutes(hub.NewController(s.partyService, s.localPlatformService, s.hubUc, s.cfg.Ocpi)))
	routeBuilder.SetRoutes(locations.GetRoutes(locations.NewController(s.locationUc, s.locationService, s.localPlatformService, s.locConverter, s.cfg.Ocpi)))
	routeBuilder.SetRoutes(tariffs.GetRoutes(tariffs.NewController(s.trfService, s.localPlatformService, s.trfConverter, s.trfUc, s.cfg.Ocpi)))
	routeBuilder.SetRoutes(tokens.GetRoutes(tokens.NewController(s.tknService, s.localPlatformService, s.tknConverter, s.tknUc, s.cfg.Ocpi)))
	routeBuilder.SetRoutes(sessions.GetRoutes(sessions.NewController(s.sessService, s.localPlatformService, s.sessConverter, s.sessUc, s.cfg.Ocpi)))
	routeBuilder.SetRoutes(cdrs.GetRoutes(cdrs.NewController(s.cdrService, s.localPlatformService, s.cdrConverter, s.cdrUc, s.cfg.Ocpi)))
	routeBuilder.SetRoutes(commands.GetRoutes(commands.NewController(s.cmdUc)))

	// backend routing
	routeBuilder.SetRoutes(bkndPlatform.GetRoutes(bkndPlatform.NewController(s.platformService, s.credentialsUc, s.credConverter, s.tokenGen)))
	routeBuilder.SetRoutes(bkndWebhook.GetRoutes(bkndWebhook.NewController(s.webhookService)))
	routeBuilder.SetRoutes(bkndParty.GetRoutes(bkndParty.NewController(s.credentialsUc, s.hubUc, s.credConverter, s.localPlatformService, s.partyService)))
	routeBuilder.SetRoutes(bkndLoc.GetRoutes(bkndLoc.NewController(s.locationUc, s.locConverter, s.localPlatformService, s.locationService)))
	routeBuilder.SetRoutes(bkndTrf.GetRoutes(bkndTrf.NewController(s.trfUc, s.trfConverter, s.localPlatformService, s.trfService)))
	routeBuilder.SetRoutes(bkndTkn.GetRoutes(bkndTkn.NewController(s.tknUc, s.tknConverter, s.localPlatformService, s.tknService)))
	routeBuilder.SetRoutes(bkndSess.GetRoutes(bkndSess.NewController(s.sessUc, s.sessConverter, s.localPlatformService, s.sessService)))
	routeBuilder.SetRoutes(bkndCdrs.GetRoutes(bkndCdrs.NewController(s.cdrUc, s.cdrConverter, s.localPlatformService, s.cdrService)))
	routeBuilder.SetRoutes(bkndCmd.GetRoutes(bkndCmd.NewController(s.cmdUc, s.cmdConverter, s.localPlatformService, s.cmdService)))
	routeBuilder.SetRoutes(bkndMnt.GetRoutes(bkndMnt.NewController(s.maintenanceUc, s.logService)))
	routeBuilder.SetRoutes(bkndSwg.GetRoutes())

	return routeBuilder.Build()
}

// Init does all initializations
func (s *ServiceImpl) Init(ctx context.Context) error {

	// load config
	var err error
	s.cfg, err = s.loadCfgFn()
	if err != nil {
		return err
	}

	// set log config
	ocpi.Logger.Init(s.cfg.Log)

	// init storage
	if err := s.storageAdapter.Init(ctx, s.cfg.Storages); err != nil {
		return err
	}

	// init adapters
	if err := s.ocpiAdapter.Init(ctx, s.cfg.Ocpi); err != nil {
		return err
	}
	if err := s.webhookAdapter.Init(ctx, s.cfg.Ocpi.Local.Webhook); err != nil {
		return err
	}

	// init services
	if err := s.platformService.Init(ctx, s.cfg.Ocpi); err != nil {
		return err
	}
	if err := s.localPlatformService.Init(ctx, s.cfg.Ocpi); err != nil {
		return err
	}

	// init http server
	if err := s.initHttpServer(ctx); err != nil {
		return err
	}

	// init grpc server
	s.grpc = grpc.New(s.credentialsUc, s.locationUc, s.hubUc, s.trfUc, s.tknUc, s.sessUc, s.cdrUc)
	if err = s.grpc.Init(s.cfg.Grpc); err != nil {
		return err
	}

	// init monitoring
	if s.cfg.Monitoring.Enabled {
		// enable error monitoring
		errMonitor := monitoring.NewErrorMonitoring()
		ocpi.Logger.SetErrorHook(errMonitor)
		// register metrics collectors
		if err := s.monitoring.Init(s.cfg.Monitoring, errMonitor); err != nil {
			return err
		}
	}

	// profiling
	if s.cfg.Profile.Enabled {
		// init profiling server
		if err := s.profile.Init(s.cfg.Profile); err != nil {
			return err
		}
	}

	// register cron
	ocpiCron.NewCron(s.cronManager, s.cmdUc).Register(ctx)

	return nil
}

func (s *ServiceImpl) Start(ctx context.Context) error {

	s.grpc.ListenAsync()

	// listen HTTP connections
	s.http.Listen()

	// listen for scraping metrics
	if s.cfg.Monitoring.Enabled {
		s.monitoring.Listen()
	}

	// profiling listener
	if s.cfg.Profile.Enabled {
		s.profile.Listen()
	}

	// initialize home ocpi platform
	if err := s.localPlatformService.InitializePlatform(ctx); err != nil {
		return err
	}

	//start cron manager
	s.cronManager.Start(ctx)

	return nil
}

func (s *ServiceImpl) Close(ctx context.Context) {
	s.cronManager.Stop(ctx)
	_ = s.storageAdapter.Close(ctx)
	_ = s.ocpiAdapter.Close(ctx)
	_ = s.webhookAdapter.Close(ctx)
	if s.cfg.Monitoring.Enabled {
		s.monitoring.Close()
	}
	if s.cfg.Profile.Enabled {
		s.profile.Close()
	}
	s.http.Close()
	s.grpc.Close()
}
