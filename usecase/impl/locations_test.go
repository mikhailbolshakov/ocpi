package impl

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/errors"
	"github.com/mikhailbolshakov/ocpi/mocks"
	"github.com/mikhailbolshakov/ocpi/model"
	"github.com/mikhailbolshakov/ocpi/usecase"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/atomic"
	"testing"
	"time"
)

type locationUcTestSuite struct {
	kit.Suite
	uc                usecase.LocationUc
	platformService   *mocks.PlatformService
	locationService   *mocks.LocationService
	remoteLocationRep *mocks.RemoteLocationRepository
	partyService      *mocks.PartyService
	webhook           *mocks.WebhookCallService
	localPlatform     *mocks.LocalPlatformService
}

func (s *locationUcTestSuite) SetupSuite() {
	s.Suite.Init(ocpi.LF())
}

func (s *locationUcTestSuite) SetupTest() {
	s.platformService = &mocks.PlatformService{}
	s.locationService = &mocks.LocationService{}
	s.remoteLocationRep = &mocks.RemoteLocationRepository{}
	s.partyService = &mocks.PartyService{}
	s.webhook = &mocks.WebhookCallService{}
	s.localPlatform = &mocks.LocalPlatformService{}
	s.uc = NewLocationUc(s.platformService, s.locationService, s.remoteLocationRep, s.partyService, s.webhook, s.localPlatform, nil)
}

func (s *locationUcTestSuite) TearDownSuite() {}

func TestLocationUcSuite(t *testing.T) {
	suite.Run(t, new(locationUcTestSuite))
}

func (s *locationUcTestSuite) Test_OnLocalLocationChanged_WhenNotOfLocalPlatform() {
	s.localPlatform.On("Get", mock.Anything).Return(&domain.Platform{}, nil)
	loc := &domain.Location{Id: kit.NewId()}
	loc.PlatformId = "remote"
	s.locationService.On("GetLocation", s.Ctx, loc.Id, false).Return(loc, nil)
	s.localPlatform.On("GetPlatformId", s.Ctx).Return("local")
	s.AssertAppErr(s.uc.OnLocalLocationChanged(s.Ctx, loc), errors.ErrCodeLocNotBelongLocalPlatform)
}

func (s *locationUcTestSuite) Test_OnLocalLocationChanged_WhenNoChanges() {
	s.localPlatform.On("Get", mock.Anything).Return(&domain.Platform{}, nil)
	loc := &domain.Location{Id: kit.NewId()}
	loc.PlatformId = "local"
	s.locationService.On("GetLocation", s.Ctx, loc.Id, false).Return(loc, nil)
	s.locationService.On("PutLocation", s.Ctx, loc).Return(nil, nil)
	s.localPlatform.On("GetPlatformId", s.Ctx).Return("local")
	s.NoError(s.uc.OnLocalLocationChanged(s.Ctx, loc))
}

func (s *locationUcTestSuite) Test_OnLocalLocationChanged_UpdateAndRemotePut() {
	s.localPlatform.On("Get", mock.Anything).Return(&domain.Platform{}, nil)
	loc := &domain.Location{Id: kit.NewId()}
	loc.PlatformId = "local"
	s.locationService.On("GetLocation", s.Ctx, loc.Id, false).Return(loc, nil)
	s.locationService.On("PutLocation", s.Ctx, loc).Return(loc, nil)
	s.platformService.On("RoleEndpoint", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(domain.Endpoint("url"))
	platforms := []*domain.Platform{
		{TokenC: domain.PlatformToken(kit.NewRandString()), Protocol: &domain.ProtocolDetails{PushSupport: domain.PushSupport{Locations: true}}},
		{TokenC: domain.PlatformToken(kit.NewRandString()), Protocol: &domain.ProtocolDetails{PushSupport: domain.PushSupport{Locations: true}}},
	}
	s.platformService.On("Search", s.Ctx, mock.Anything).Return(platforms, nil)
	s.localPlatform.On("GetPlatformId", s.Ctx).Return("local")
	s.remoteLocationRep.On("PutLocationAsync", s.Ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	s.NoError(s.uc.OnLocalLocationChanged(s.Ctx, loc))
	s.AssertNumberOfCalls(&s.remoteLocationRep.Mock, "PutLocationAsync", 2)
}

func (s *locationUcTestSuite) Test_OnRemoteLocationPut_WhenNotOfRemotePlatform() {
	s.localPlatform.On("Get", mock.Anything).Return(&domain.Platform{}, nil)
	platform := &domain.Platform{Id: kit.NewId(), Status: domain.ConnectionStatusConnected}
	loc := &domain.Location{Id: kit.NewId()}
	loc.PlatformId = "local"
	locOcpi := &model.OcpiLocation{Id: loc.Id}
	s.platformService.On("Get", s.Ctx, platform.Id).Return(platform, nil)
	s.locationService.On("GetLocation", s.Ctx, loc.Id, false).Return(loc, nil)
	s.localPlatform.On("GetPlatformId", s.Ctx).Return("local")
	s.AssertAppErr(s.uc.OnRemoteLocationPut(s.Ctx, platform.Id, locOcpi), errors.ErrCodeLocNotBelongRemotePlatform)
}

func (s *locationUcTestSuite) Test_OnRemoteLocationPatch_Ok() {
	s.localPlatform.On("Get", mock.Anything).Return(&domain.Platform{}, nil)
	platform := &domain.Platform{Id: kit.NewId(), Status: domain.ConnectionStatusConnected}
	locOcpi := &model.OcpiLocation{Id: kit.NewId()}
	s.platformService.On("Get", s.Ctx, platform.Id).Return(platform, nil)
	s.locationService.On("GetLocation", s.Ctx, locOcpi.Id, false).Return(nil, nil)
	s.partyService.On("GetByExtId", mock.Anything, mock.Anything).Return(&domain.Party{}, nil)
	s.webhook.On("OnLocationsChanged", mock.Anything, mock.Anything).Return(nil)
	s.locationService.On("MergeLocation", s.Ctx, mock.Anything).Return(&domain.Location{}, nil)
	s.NoError(s.uc.OnRemoteLocationPatch(s.Ctx, platform.Id, locOcpi))
	s.AssertCalled(&s.locationService.Mock, "MergeLocation", s.Ctx, mock.Anything)
	s.AssertCalled(&s.webhook.Mock, "OnLocationsChanged", mock.Anything, mock.Anything)
}

func (s *locationUcTestSuite) Test_OnRemoteLocationPull() {
	s.localPlatform.On("Get", mock.Anything).Return(&domain.Platform{}, nil)

	platforms := []*domain.Platform{
		{Id: kit.NewId(), TokenC: domain.PlatformToken(kit.NewRandString())},
		{Id: kit.NewId(), TokenC: domain.PlatformToken(kit.NewRandString())},
	}
	s.platformService.On("Search", s.Ctx, mock.Anything).Return(platforms, nil)
	s.platformService.On("RoleEndpoint", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(domain.Endpoint("url"))
	s.platformService.On("Get", mock.Anything, mock.Anything).Return(&domain.Platform{Status: domain.ConnectionStatusConnected}, nil)
	s.partyService.On("GetByExtId", mock.Anything, mock.Anything).Return(&domain.Party{}, nil)
	s.locationService.On("PutLocation", mock.Anything, mock.Anything).Return(&domain.Location{}, nil)
	s.locationService.On("GetLocation", mock.Anything, mock.Anything, false).Return(nil, nil)
	cnt := atomic.NewInt32(0)
	s.webhook.On("OnLocationsChanged", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			cnt.Inc()
		}).
		Return(nil)

	locations := []*model.OcpiLocation{
		{
			Id: kit.NewId(),
		},
		{
			Id: kit.NewId(),
		},
	}
	for _, p := range platforms {
		tknC := p.TokenC
		s.remoteLocationRep.On("GetLocations", mock.Anything, mock.MatchedBy(func(rq *usecase.OcpiRepositoryPagingRequest) bool {
			return rq.Token == tknC && *rq.Offset == 500
		})).Return(nil, nil)
		s.remoteLocationRep.On("GetLocations", mock.Anything, mock.MatchedBy(func(rq *usecase.OcpiRepositoryPagingRequest) bool {
			return rq.Token == tknC && *rq.Offset < 500
		})).Return(locations, nil)
	}

	s.NoError(s.uc.OnRemoteLocationsPull(s.Ctx, nil, nil))
	if err := <-kit.Await(func() (bool, error) {
		return cnt.Load() == 20, nil
	}, time.Millisecond*50, time.Second*5); err != nil {
		s.Fatal(err)
	}
}

func (s *locationUcTestSuite) Test_OnRemoteLocationPull_WhenNoLocations() {
	s.localPlatform.On("Get", mock.Anything).Return(&domain.Platform{}, nil)

	platforms := []*domain.Platform{
		{Id: kit.NewId(), TokenC: domain.PlatformToken(kit.NewRandString())},
		{Id: kit.NewId(), TokenC: domain.PlatformToken(kit.NewRandString())},
	}
	s.platformService.On("Search", s.Ctx, mock.Anything).Return(platforms, nil)
	s.platformService.On("RoleEndpoint", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(domain.Endpoint("url"))
	s.platformService.On("Get", mock.Anything, mock.Anything).Return(&domain.Platform{Status: domain.ConnectionStatusConnected}, nil)
	s.partyService.On("GetByExtId", mock.Anything, mock.Anything).Return(&domain.Party{}, nil)
	s.locationService.On("PutLocation", mock.Anything, mock.Anything).Return(&domain.Location{}, nil)
	s.webhook.On("OnLocationsChanged", mock.Anything, mock.Anything).Return(nil)

	cnt := atomic.NewInt32(0)
	s.remoteLocationRep.On("GetLocations", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			cnt.Inc()
		}).
		Return(nil, nil)

	s.NoError(s.uc.OnRemoteLocationsPull(s.Ctx, nil, nil))
	if err := <-kit.Await(func() (bool, error) {
		return cnt.Load() == 2, nil
	}, time.Millisecond*50, time.Second*5); err != nil {
		s.Fatal(err)
	}
}

func (s *locationUcTestSuite) Test_OnLocalEvseChanged_WhenNotOfLocalPlatform() {
	s.localPlatform.On("Get", mock.Anything).Return(&domain.Platform{}, nil)
	evse := &domain.Evse{Id: kit.NewId(), LocationId: kit.NewId()}
	evse.PlatformId = "remote"
	s.locationService.On("GetEvse", s.Ctx, evse.LocationId, evse.Id, false).Return(evse, nil)
	s.localPlatform.On("GetPlatformId", s.Ctx).Return("local")
	s.AssertAppErr(s.uc.OnLocalEvseChanged(s.Ctx, evse), errors.ErrCodeLocNotBelongLocalPlatform)
}

func (s *locationUcTestSuite) Test_OnLocalEvseChanged_WhenNoChanges() {
	s.localPlatform.On("Get", mock.Anything).Return(&domain.Platform{}, nil)
	evse := &domain.Evse{Id: kit.NewId(), LocationId: kit.NewId()}
	evse.PlatformId = "local"
	s.locationService.On("GetEvse", s.Ctx, evse.LocationId, evse.Id, false).Return(evse, nil)
	s.locationService.On("PutEvse", s.Ctx, evse).Return(nil, nil)
	s.localPlatform.On("GetPlatformId", s.Ctx).Return("local")
	s.NoError(s.uc.OnLocalEvseChanged(s.Ctx, evse))
}

func (s *locationUcTestSuite) Test_OnLocalEvseChanged_WhenNoPush() {
	s.localPlatform.On("Get", mock.Anything).Return(&domain.Platform{}, nil)
	evse := &domain.Evse{Id: kit.NewId(), LocationId: kit.NewId()}
	evse.PlatformId = "local"
	s.locationService.On("GetEvse", s.Ctx, evse.LocationId, evse.Id, false).Return(evse, nil)
	s.locationService.On("PutEvse", s.Ctx, evse).Return(evse, nil)
	s.platformService.On("RoleEndpoint", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(domain.Endpoint("url"))
	platforms := []*domain.Platform{
		{TokenC: domain.PlatformToken(kit.NewRandString()), Protocol: &domain.ProtocolDetails{PushSupport: domain.PushSupport{Locations: false}}},
		{TokenC: domain.PlatformToken(kit.NewRandString()), Protocol: &domain.ProtocolDetails{PushSupport: domain.PushSupport{Locations: false}}},
	}
	s.platformService.On("Search", s.Ctx, mock.Anything).Return(platforms, nil)
	s.localPlatform.On("GetPlatformId", s.Ctx).Return("local")
	s.NoError(s.uc.OnLocalEvseChanged(s.Ctx, evse))
}

func (s *locationUcTestSuite) Test_OnLocalEvseChanged_UpdateAndRemotePut() {
	s.localPlatform.On("Get", mock.Anything).Return(&domain.Platform{}, nil)
	evse := &domain.Evse{Id: kit.NewId(), LocationId: kit.NewId()}
	evse.PlatformId = "local"
	s.locationService.On("GetEvse", s.Ctx, evse.LocationId, evse.Id, false).Return(evse, nil)
	s.locationService.On("PutEvse", s.Ctx, evse).Return(evse, nil)
	s.platformService.On("RoleEndpoint", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(domain.Endpoint("url"))
	platforms := []*domain.Platform{
		{TokenC: domain.PlatformToken(kit.NewRandString()), Protocol: &domain.ProtocolDetails{PushSupport: domain.PushSupport{Locations: true}}},
		{TokenC: domain.PlatformToken(kit.NewRandString()), Protocol: &domain.ProtocolDetails{PushSupport: domain.PushSupport{Locations: true}}},
	}
	s.platformService.On("Search", s.Ctx, mock.Anything).Return(platforms, nil)
	s.localPlatform.On("GetPlatformId", s.Ctx).Return("local")
	s.remoteLocationRep.On("PutEvseAsync", s.Ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	s.NoError(s.uc.OnLocalEvseChanged(s.Ctx, evse))
	s.AssertNumberOfCalls(&s.remoteLocationRep.Mock, "PutEvseAsync", 2)
}

func (s *locationUcTestSuite) Test_OnLocalEvseStatusChanged_UpdateAndRemotePut() {
	s.localPlatform.On("Get", mock.Anything).Return(&domain.Platform{}, nil)
	evse := &domain.Evse{Id: kit.NewId(), LocationId: kit.NewId()}
	evse.PlatformId = "local"
	s.locationService.On("GetEvse", s.Ctx, evse.LocationId, evse.Id, false).Return(evse, nil)
	s.locationService.On("MergeEvse", s.Ctx, evse).Return(evse, nil)
	s.platformService.On("RoleEndpoint", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(domain.Endpoint("url"))
	platforms := []*domain.Platform{
		{TokenC: domain.PlatformToken(kit.NewRandString()), Protocol: &domain.ProtocolDetails{PushSupport: domain.PushSupport{Locations: true}}},
		{TokenC: domain.PlatformToken(kit.NewRandString()), Protocol: &domain.ProtocolDetails{PushSupport: domain.PushSupport{Locations: true}}},
	}
	s.platformService.On("Search", s.Ctx, mock.Anything).Return(platforms, nil)
	s.localPlatform.On("GetPlatformId", s.Ctx).Return("local")
	s.remoteLocationRep.On("PatchEvseAsync", s.Ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	s.NoError(s.uc.OnLocalEvseStatusChanged(s.Ctx, evse.LocationId, evse.Id, domain.EvseStatusAvailable))
	s.AssertNumberOfCalls(&s.remoteLocationRep.Mock, "PatchEvseAsync", 2)
}

func (s *locationUcTestSuite) Test_OnRemoteEvsePut_WhenNotOfRemotePlatform() {
	s.localPlatform.On("Get", mock.Anything).Return(&domain.Platform{}, nil)
	platform := &domain.Platform{Id: kit.NewId(), Status: domain.ConnectionStatusConnected}
	evse := &domain.Evse{Id: kit.NewId(), LocationId: kit.NewId()}
	evse.PlatformId = "local"
	evseOcpi := &model.OcpiEvse{Uid: evse.Id}
	s.platformService.On("Get", s.Ctx, platform.Id).Return(platform, nil)
	s.locationService.On("GetEvse", s.Ctx, evse.LocationId, evse.Id, false).Return(evse, nil)
	s.localPlatform.On("GetPlatformId", s.Ctx).Return("local")
	s.AssertAppErr(s.uc.OnRemoteEvsePut(s.Ctx, platform.Id, evse.LocationId, evse.ExtId.CountryCode, evse.ExtId.PartyId, evseOcpi), errors.ErrCodeLocNotBelongRemotePlatform)
}

func (s *locationUcTestSuite) Test_OnRemoteEvsePatch_Ok() {
	s.localPlatform.On("Get", mock.Anything).Return(&domain.Platform{}, nil)
	platform := &domain.Platform{Id: kit.NewId(), Status: domain.ConnectionStatusConnected}
	evseOcpi := &model.OcpiEvse{Uid: kit.NewId()}
	locId := kit.NewId()
	s.platformService.On("Get", s.Ctx, platform.Id).Return(platform, nil)
	s.locationService.On("GetEvse", s.Ctx, locId, evseOcpi.Uid, false).Return(nil, nil)
	s.partyService.On("GetByExtId", mock.Anything, mock.Anything).Return(&domain.Party{}, nil)
	s.webhook.On("OnEvseChanged", mock.Anything, mock.Anything).Return(nil)
	s.locationService.On("MergeEvse", s.Ctx, mock.Anything).Return(&domain.Evse{}, nil)
	s.NoError(s.uc.OnRemoteEvsePatch(s.Ctx, platform.Id, locId, "", "", evseOcpi))
	s.AssertCalled(&s.locationService.Mock, "MergeEvse", s.Ctx, mock.Anything)
	s.AssertCalled(&s.webhook.Mock, "OnEvseChanged", mock.Anything, mock.Anything)
}

func (s *locationUcTestSuite) Test_OnLocalConChanged_WhenNotOfLocalPlatform() {
	s.localPlatform.On("Get", mock.Anything).Return(&domain.Platform{}, nil)
	con := &domain.Connector{Id: kit.NewId(), LocationId: kit.NewId(), EvseId: kit.NewId()}
	con.PlatformId = "remote"
	s.locationService.On("GetConnector", s.Ctx, con.LocationId, con.EvseId, con.Id).Return(con, nil)
	s.localPlatform.On("GetPlatformId", s.Ctx).Return("local")
	s.AssertAppErr(s.uc.OnLocalConnectorChanged(s.Ctx, con), errors.ErrCodeLocNotBelongLocalPlatform)
}

func (s *locationUcTestSuite) Test_OnLocalConChanged_WhenNoChanges() {
	s.localPlatform.On("Get", mock.Anything).Return(&domain.Platform{}, nil)
	con := &domain.Connector{Id: kit.NewId(), LocationId: kit.NewId(), EvseId: kit.NewId()}
	con.PlatformId = "remote"
	s.locationService.On("GetConnector", s.Ctx, con.LocationId, con.EvseId, con.Id).Return(nil, nil)
	s.locationService.On("PutConnector", s.Ctx, con).Return(nil, nil)
	s.localPlatform.On("GetPlatformId", s.Ctx).Return("local")
	s.NoError(s.uc.OnLocalConnectorChanged(s.Ctx, con))
}

func (s *locationUcTestSuite) Test_OnLocalConChanged_WhenNoPush() {
	s.localPlatform.On("Get", mock.Anything).Return(&domain.Platform{}, nil)
	con := &domain.Connector{Id: kit.NewId(), LocationId: kit.NewId(), EvseId: kit.NewId()}
	con.PlatformId = "remote"
	s.locationService.On("GetConnector", s.Ctx, con.LocationId, con.EvseId, con.Id).Return(nil, nil)
	s.locationService.On("PutConnector", s.Ctx, con).Return(con, nil)
	s.platformService.On("RoleEndpoint", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(domain.Endpoint("url"))
	platforms := []*domain.Platform{
		{TokenC: domain.PlatformToken(kit.NewRandString()), Protocol: &domain.ProtocolDetails{PushSupport: domain.PushSupport{Locations: false}}},
		{TokenC: domain.PlatformToken(kit.NewRandString()), Protocol: &domain.ProtocolDetails{PushSupport: domain.PushSupport{Locations: false}}},
	}
	s.platformService.On("Search", s.Ctx, mock.Anything).Return(platforms, nil)
	s.localPlatform.On("GetPlatformId", s.Ctx).Return("local")
	s.NoError(s.uc.OnLocalConnectorChanged(s.Ctx, con))
}

func (s *locationUcTestSuite) Test_OnLocalConChanged_UpdateAndRemotePut() {
	s.localPlatform.On("Get", mock.Anything).Return(&domain.Platform{}, nil)
	con := &domain.Connector{Id: kit.NewId(), LocationId: kit.NewId(), EvseId: kit.NewId()}
	con.PlatformId = "remote"
	s.locationService.On("GetConnector", s.Ctx, con.LocationId, con.EvseId, con.Id).Return(nil, nil)
	s.locationService.On("PutConnector", s.Ctx, con).Return(con, nil)
	s.platformService.On("RoleEndpoint", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(domain.Endpoint("url"))
	platforms := []*domain.Platform{
		{TokenC: domain.PlatformToken(kit.NewRandString()), Protocol: &domain.ProtocolDetails{PushSupport: domain.PushSupport{Locations: true}}},
		{TokenC: domain.PlatformToken(kit.NewRandString()), Protocol: &domain.ProtocolDetails{PushSupport: domain.PushSupport{Locations: true}}},
	}
	s.platformService.On("Search", s.Ctx, mock.Anything).Return(platforms, nil)
	s.localPlatform.On("GetPlatformId", s.Ctx).Return("local")
	s.remoteLocationRep.On("PutConnectorAsync", s.Ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	s.NoError(s.uc.OnLocalConnectorChanged(s.Ctx, con))
	s.AssertNumberOfCalls(&s.remoteLocationRep.Mock, "PutConnectorAsync", 2)
}

func (s *locationUcTestSuite) Test_OnRemoteConPut_WhenNotOfRemotePlatform() {
	platform := &domain.Platform{Id: kit.NewId(), Status: domain.ConnectionStatusConnected}
	con := &domain.Connector{Id: kit.NewId(), LocationId: kit.NewId()}
	con.PlatformId = "local"
	conOcpi := &model.OcpiConnector{Id: con.Id}
	s.platformService.On("Get", s.Ctx, platform.Id).Return(platform, nil)
	s.locationService.On("GetConnector", s.Ctx, con.LocationId, con.EvseId, con.Id).Return(con, nil)
	s.localPlatform.On("GetPlatformId", s.Ctx).Return("local")
	s.AssertAppErr(s.uc.OnRemoteConnectorPut(s.Ctx, platform.Id, con.LocationId, con.EvseId, con.ExtId.CountryCode, con.ExtId.PartyId, conOcpi), errors.ErrCodeLocNotBelongRemotePlatform)
}

func (s *locationUcTestSuite) Test_OnRemoteConPatch_Ok() {
	platform := &domain.Platform{Id: kit.NewId(), Status: domain.ConnectionStatusConnected}
	conOcpi := &model.OcpiConnector{Id: kit.NewId()}
	locId, evseId := kit.NewId(), kit.NewId()
	s.platformService.On("Get", s.Ctx, platform.Id).Return(platform, nil)
	s.locationService.On("GetConnector", s.Ctx, locId, evseId, conOcpi.Id).Return(nil, nil)
	s.partyService.On("GetByExtId", mock.Anything, mock.Anything).Return(&domain.Party{}, nil)
	s.webhook.On("OnConnectorChanged", mock.Anything, mock.Anything).Return(nil)
	s.locationService.On("MergeConnector", s.Ctx, mock.Anything).Return(&domain.Connector{}, nil)
	s.NoError(s.uc.OnRemoteConnectorPatch(s.Ctx, platform.Id, locId, evseId, "", "", conOcpi))
	s.AssertCalled(&s.locationService.Mock, "MergeConnector", s.Ctx, mock.Anything)
	s.AssertCalled(&s.webhook.Mock, "OnConnectorChanged", mock.Anything, mock.Anything)
}
