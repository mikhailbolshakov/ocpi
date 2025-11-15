package ocpi

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/cluster"
	"github.com/mikhailbolshakov/kit/goroutine"
	service "github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/model"
	"github.com/mikhailbolshakov/ocpi/usecase"
)

type Adapter interface {
	cluster.Adapter
	usecase.RemotePlatformRepository
	usecase.RemoteHubClientInfoRepository
	usecase.RemoteLocationRepository
	usecase.RemoteTariffRepository
	usecase.RemoteTokenRepository
	usecase.RemoteSessionRepository
	usecase.RemoteCommandRepository
	usecase.RemoteCdrRepository
}

type adapterImpl struct {
	ocpiRestClient
	logService domain.OcpiLogService
	cfg        *service.CfgOcpiConfig
}

func NewAdapter(logService domain.OcpiLogService) Adapter {
	return &adapterImpl{
		logService: logService,
	}
}

func (a *adapterImpl) l() kit.CLogger {
	return service.L().Cmp("ocpi-adapter")
}

func (a *adapterImpl) Init(ctx context.Context, config interface{}) error {
	a.l().Mth("init").Dbg()
	a.cfg = config.(*service.CfgOcpiConfig)
	if a.cfg.Remote.Mock {
		a.ocpiRestClient = newMockOcpiRestClient(a.logService)
	} else {
		a.ocpiRestClient = newOcpiRestClient(a.logService)
	}
	if err := a.ocpiRestClient.Init(ctx, a.cfg.Remote); err != nil {
		return err
	}
	return nil
}

func (a *adapterImpl) Close(ctx context.Context) error {
	return a.ocpiRestClient.Close(ctx)
}

func (a *adapterImpl) GetVersions(ctx context.Context, rq *usecase.OcpiRepositoryBaseRequest) (domain.Versions, error) {
	a.l().C(ctx).Mth("get-versions").Dbg()
	rs, err := a.ocpiRestClient.GetVersions(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId)
	if err != nil {
		return nil, err
	}
	return a.toVersionsDomain(rs), nil
}

func (a *adapterImpl) GetVersionDetails(ctx context.Context, rq *usecase.OcpiRepositoryBaseRequest) (domain.ModuleEndpoints, error) {
	a.l().C(ctx).Mth("get-ver-details").Dbg()
	rs, err := a.ocpiRestClient.GetVersionDetails(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId)
	if err != nil {
		return nil, err
	}
	return a.toVersionDetailsDomain(rs), nil
}

func (a *adapterImpl) PostCredentials(ctx context.Context, rq *usecase.OcpiRepositoryRequestG[*model.OcpiCredentials]) (*model.OcpiCredentials, error) {
	a.l().C(ctx).Mth("post-cred").Dbg()
	return a.ocpiRestClient.PostCredentials(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, rq.Request)
}

func (a *adapterImpl) PutCredentials(ctx context.Context, rq *usecase.OcpiRepositoryRequestG[*model.OcpiCredentials]) (*model.OcpiCredentials, error) {
	a.l().C(ctx).Mth("put-cred").Dbg()
	return a.ocpiRestClient.PutCredentials(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, rq.Request)
}

func (a *adapterImpl) GetCredentials(ctx context.Context, rq *usecase.OcpiRepositoryBaseRequest) (*model.OcpiCredentials, error) {
	a.l().C(ctx).Mth("get-cred").Dbg()
	return a.ocpiRestClient.GetCredentials(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId)
}

func (a *adapterImpl) DeleteCredentials(ctx context.Context, rq *usecase.OcpiRepositoryBaseRequest) error {
	a.l().C(ctx).Mth("del-cred").Dbg()
	return a.ocpiRestClient.DeleteCredentials(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId)
}

func (a *adapterImpl) PutClientInfoAsync(ctx context.Context, rq *usecase.OcpiRepositoryErrHandlerRequestG[*model.OcpiClientInfo]) {
	l := a.l().C(ctx).Mth("put-client-info").Dbg()
	goroutine.New().WithLogger(l).Go(ctx, func() {
		if err := a.ocpiRestClient.PutClientInfo(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, rq.Request); err != nil {
			rq.Handler(err)
		}
	})
}

func (a *adapterImpl) PutClientInfo(ctx context.Context, rq *usecase.OcpiRepositoryRequestG[*model.OcpiClientInfo]) error {
	a.l().C(ctx).Mth("put-client-info").Dbg()
	return a.ocpiRestClient.PutClientInfo(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, rq.Request)
}

func (a *adapterImpl) GetClientInfos(ctx context.Context, rq *usecase.OcpiRepositoryPagingRequest) ([]*model.OcpiClientInfo, error) {
	a.l().C(ctx).Mth("get-hub-client-page").Dbg()
	return a.ocpiRestClient.GetHubClientInfo(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, &rq.OcpiGetPageRequest)
}

func (a *adapterImpl) PutLocationAsync(ctx context.Context, rq *usecase.OcpiRepositoryErrHandlerRequestG[*model.OcpiLocation]) {
	l := a.l().C(ctx).Mth("put-location").Dbg()
	goroutine.New().WithLogger(l).Go(ctx, func() {
		if err := a.ocpiRestClient.PutLocation(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, rq.Request); err != nil {
			rq.Handler(err)
		}
	})
}

func (a *adapterImpl) PatchLocationAsync(ctx context.Context, rq *usecase.OcpiRepositoryErrHandlerRequestG[*model.OcpiLocation]) {
	l := a.l().C(ctx).Mth("patch-location").Dbg()
	goroutine.New().WithLogger(l).Go(ctx, func() {
		if err := a.ocpiRestClient.PatchLocation(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, rq.Request); err != nil {
			rq.Handler(err)
		}
	})
}

func (a *adapterImpl) GetLocations(ctx context.Context, rq *usecase.OcpiRepositoryPagingRequest) ([]*model.OcpiLocation, error) {
	a.l().C(ctx).Mth("get-loc-page").Dbg()
	return a.ocpiRestClient.GetLocationPage(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, &rq.OcpiGetPageRequest)
}

func (a *adapterImpl) GetLocation(ctx context.Context, rq *usecase.OcpiRepositoryIdRequest) (*model.OcpiLocation, error) {
	a.l().C(ctx).Mth("get-loc").Dbg()
	return a.ocpiRestClient.GetLocation(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, rq.Id)
}

func (a *adapterImpl) PutEvseAsync(ctx context.Context, rq *usecase.OcpiRepositoryErrHandlerRequestG[*model.OcpiEvse], party *model.OcpiPartyId, locId string) {
	l := a.l().C(ctx).Mth("put-evse-async").Dbg()
	goroutine.New().WithLogger(l).Go(ctx, func() {
		if err := a.ocpiRestClient.PutEvse(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, rq.Request, party, locId); err != nil {
			rq.Handler(err)
		}
	})
}

func (a *adapterImpl) PatchEvseAsync(ctx context.Context, rq *usecase.OcpiRepositoryErrHandlerRequestG[*model.OcpiEvse], party *model.OcpiPartyId, locId string) {
	l := a.l().C(ctx).Mth("patch-evse-async").Dbg()
	goroutine.New().WithLogger(l).Go(ctx, func() {
		if err := a.ocpiRestClient.PatchEvse(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, rq.Request, party, locId); err != nil {
			rq.Handler(err)
		}
	})
}

func (a *adapterImpl) GetEvse(ctx context.Context, rq *usecase.OcpiRepositoryBaseRequest, locId, evseId string) (*model.OcpiEvse, error) {
	a.l().C(ctx).Mth("get-evse").Dbg()
	return a.ocpiRestClient.GetEvse(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, locId, evseId)
}

func (a *adapterImpl) PutConnectorAsync(ctx context.Context, rq *usecase.OcpiRepositoryErrHandlerRequestG[*model.OcpiConnector], party *model.OcpiPartyId, evseId, locId string) {
	l := a.l().C(ctx).Mth("put-con-async").Dbg()
	goroutine.New().WithLogger(l).Go(ctx, func() {
		if err := a.ocpiRestClient.PutCon(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, rq.Request, party, locId, evseId); err != nil {
			rq.Handler(err)
		}
	})
}

func (a *adapterImpl) PatchConnectorAsync(ctx context.Context, rq *usecase.OcpiRepositoryErrHandlerRequestG[*model.OcpiConnector], party *model.OcpiPartyId, evseId, locId string) {
	l := a.l().C(ctx).Mth("patch-con-async").Dbg()
	goroutine.New().WithLogger(l).Go(ctx, func() {
		if err := a.ocpiRestClient.PatchCon(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, rq.Request, party, locId, evseId); err != nil {
			rq.Handler(err)
		}
	})
}

func (a *adapterImpl) GetConnector(ctx context.Context, rq *usecase.OcpiRepositoryBaseRequest, locId, evseId, conId string) (*model.OcpiConnector, error) {
	a.l().C(ctx).Mth("get-con").Dbg()
	return a.ocpiRestClient.GetCon(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, locId, evseId, conId)
}

func (a *adapterImpl) PutTariffAsync(ctx context.Context, rq *usecase.OcpiRepositoryErrHandlerRequestG[*model.OcpiTariff]) {
	l := a.l().C(ctx).Mth("put-trf-async").Dbg()
	goroutine.New().WithLogger(l).Go(ctx, func() {
		if err := a.ocpiRestClient.PutTariff(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, rq.Request); err != nil {
			rq.Handler(err)
		}
	})
}

func (a *adapterImpl) PatchTariffAsync(ctx context.Context, rq *usecase.OcpiRepositoryErrHandlerRequestG[*model.OcpiTariff]) {
	l := a.l().C(ctx).Mth("patch-trf-async").Dbg()
	goroutine.New().WithLogger(l).Go(ctx, func() {
		if err := a.ocpiRestClient.PatchTariff(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, rq.Request); err != nil {
			rq.Handler(err)
		}
	})
}

func (a *adapterImpl) GetTariffs(ctx context.Context, rq *usecase.OcpiRepositoryPagingRequest) ([]*model.OcpiTariff, error) {
	a.l().C(ctx).Mth("get-trf-page").Dbg()
	return a.ocpiRestClient.GetTariffPage(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, &rq.OcpiGetPageRequest)
}

func (a *adapterImpl) GetTariff(ctx context.Context, rq *usecase.OcpiRepositoryIdRequest) (*model.OcpiTariff, error) {
	a.l().C(ctx).Mth("get-trf").Dbg()
	return a.ocpiRestClient.GetTariff(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, rq.Id)
}

func (a *adapterImpl) PutTokenAsync(ctx context.Context, rq *usecase.OcpiRepositoryErrHandlerRequestG[*model.OcpiToken]) {
	l := a.l().C(ctx).Mth("put-tkn-async").Dbg()
	goroutine.New().WithLogger(l).Go(ctx, func() {
		if err := a.ocpiRestClient.PutToken(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, rq.Request); err != nil {
			rq.Handler(err)
		}
	})
}

func (a *adapterImpl) PatchTokenAsync(ctx context.Context, rq *usecase.OcpiRepositoryErrHandlerRequestG[*model.OcpiToken]) {
	l := a.l().C(ctx).Mth("patch-tkn-async").Dbg()
	goroutine.New().WithLogger(l).Go(ctx, func() {
		if err := a.ocpiRestClient.PatchToken(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, rq.Request); err != nil {
			rq.Handler(err)
		}
	})
}

func (a *adapterImpl) GetTokens(ctx context.Context, rq *usecase.OcpiRepositoryPagingRequest) ([]*model.OcpiToken, error) {
	a.l().C(ctx).Mth("get-tkn-page").Dbg()
	return a.ocpiRestClient.GetTokenPage(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, &rq.OcpiGetPageRequest)
}

func (a *adapterImpl) GetToken(ctx context.Context, rq *usecase.OcpiRepositoryIdRequest) (*model.OcpiToken, error) {
	a.l().C(ctx).Mth("get-tkn").Dbg()
	return a.ocpiRestClient.GetToken(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, rq.Id)
}

func (a *adapterImpl) PutSessionAsync(ctx context.Context, rq *usecase.OcpiRepositoryErrHandlerRequestG[*model.OcpiSession]) {
	l := a.l().C(ctx).Mth("put-sess-async").Dbg()
	goroutine.New().WithLogger(l).Go(ctx, func() {
		if err := a.ocpiRestClient.PutSession(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, rq.Request); err != nil {
			rq.Handler(err)
		}
	})
}

func (a *adapterImpl) PatchSessionAsync(ctx context.Context, rq *usecase.OcpiRepositoryErrHandlerRequestG[*model.OcpiSession]) {
	l := a.l().C(ctx).Mth("patch-sess-async").Dbg()
	goroutine.New().WithLogger(l).Go(ctx, func() {
		if err := a.ocpiRestClient.PatchSession(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, rq.Request); err != nil {
			rq.Handler(err)
		}
	})
}

func (a *adapterImpl) GetSessions(ctx context.Context, rq *usecase.OcpiRepositoryPagingRequest) ([]*model.OcpiSession, error) {
	a.l().C(ctx).Mth("get-sess-page").Dbg()
	return a.ocpiRestClient.GetSessionPage(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, &rq.OcpiGetPageRequest)
}

func (a *adapterImpl) GetSession(ctx context.Context, rq *usecase.OcpiRepositoryIdRequest) (*model.OcpiSession, error) {
	a.l().C(ctx).Mth("get-sess").Dbg()
	return a.ocpiRestClient.GetSession(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, rq.Id)
}

func (a *adapterImpl) PostCommandAsync(ctx context.Context, rq *usecase.OcpiRepositoryErrHandlerRequest, cmdType string, cmd any) {
	l := a.l().C(ctx).Mth("post-cmd-async").Dbg()
	goroutine.New().WithLogger(l).Go(ctx, func() {
		if err := a.ocpiRestClient.PostCommand(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, cmdType, cmd); err != nil {
			rq.Handler(err)
		}
	})
}

func (a *adapterImpl) PostCommandResponseAsync(ctx context.Context, rq *usecase.OcpiRepositoryErrHandlerRequestG[*model.OcpiCommandResult]) {
	l := a.l().C(ctx).Mth("post-cmd-rs-async").Dbg()
	goroutine.New().WithLogger(l).Go(ctx, func() {
		if err := a.ocpiRestClient.PostCommandResponse(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, rq.Request); err != nil {
			rq.Handler(err)
		}
	})
}

func (a *adapterImpl) PostCdrAsync(ctx context.Context, rq *usecase.OcpiRepositoryErrHandlerRequestG[*model.OcpiCdr]) {
	l := a.l().C(ctx).Mth("post-cdr-async").Dbg()
	goroutine.New().WithLogger(l).Go(ctx, func() {
		if err := a.ocpiRestClient.PostCdr(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, rq.Request); err != nil {
			rq.Handler(err)
		}
	})
}

func (a *adapterImpl) GetCdrs(ctx context.Context, rq *usecase.OcpiRepositoryPagingRequest) ([]*model.OcpiCdr, error) {
	a.l().C(ctx).Mth("get-cdrs-page").Dbg()
	return a.ocpiRestClient.GetCdrsPage(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, &rq.OcpiGetPageRequest)
}

func (a *adapterImpl) GetCdr(ctx context.Context, rq *usecase.OcpiRepositoryIdRequest) (*model.OcpiCdr, error) {
	a.l().C(ctx).Mth("get-cdr").Dbg()
	return a.ocpiRestClient.GetCdr(ctx, string(rq.Endpoint), string(rq.Token), rq.FromPlatformId, rq.ToPlatformId, rq.Id)
}
