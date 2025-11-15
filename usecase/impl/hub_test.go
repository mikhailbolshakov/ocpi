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

type hubUcTestSuite struct {
	kit.Suite
	uc              usecase.HubUc
	platformService *mocks.PlatformService
	partyService    *mocks.PartyService
	remoteRep       *mocks.RemoteHubClientInfoRepository
	webhook         *mocks.WebhookCallService
	localPlatform   *mocks.LocalPlatformService
}

func (s *hubUcTestSuite) SetupSuite() {
	s.Suite.Init(ocpi.LF())
}

func (s *hubUcTestSuite) SetupTest() {
	s.platformService = &mocks.PlatformService{}
	s.remoteRep = &mocks.RemoteHubClientInfoRepository{}
	s.partyService = &mocks.PartyService{}
	s.webhook = &mocks.WebhookCallService{}
	s.localPlatform = &mocks.LocalPlatformService{}
	s.uc = NewHubUc(s.platformService, s.remoteRep, s.partyService, s.webhook, s.localPlatform, nil)
}

func (s *hubUcTestSuite) TearDownSuite() {}

func TestHubUcSuite(t *testing.T) {
	suite.Run(t, new(hubUcTestSuite))
}

func (s *hubUcTestSuite) Test_OnLocalClientInfoChanged_WhenNotOfLocalPlatform() {
	locPlatform := &domain.Platform{}
	s.localPlatform.On("Get", s.Ctx).Return(locPlatform, nil)
	ci := &domain.Party{Id: kit.NewId()}
	ci.PlatformId = "remote"
	s.partyService.On("Get", s.Ctx, ci.Id).Return(ci, nil)
	s.localPlatform.On("GetPlatformId", s.Ctx).Return("local")
	s.AssertAppErr(s.uc.OnLocalClientInfoChanged(s.Ctx, ci), errors.ErrCodePartyNotBelongLocalPlatform)
}

func (s *hubUcTestSuite) Test_OnLocalClientInfoChanged_WhenPushNotConfigured() {
	locPlatform := &domain.Platform{}
	s.localPlatform.On("Get", s.Ctx).Return(locPlatform, nil)
	ci := &domain.Party{Id: kit.NewId(), Roles: []string{domain.RoleCPO, domain.RoleEMSP}}
	ci.PlatformId = "local"
	s.partyService.On("Get", s.Ctx, ci.Id).Return(ci, nil)
	s.platformService.On("RoleEndpoint", mock.Anything, mock.Anything, model.ModuleIdHubClientInfo, model.OcpiReceiver).Return(domain.Endpoint("url"))
	platforms := []*domain.Platform{
		{TokenC: domain.PlatformToken(kit.NewRandString()), Protocol: &domain.ProtocolDetails{PushSupport: domain.PushSupport{HubClientInfo: false}}},
		{TokenC: domain.PlatformToken(kit.NewRandString()), Protocol: &domain.ProtocolDetails{PushSupport: domain.PushSupport{HubClientInfo: false}}},
	}
	s.platformService.On("Search", s.Ctx, mock.Anything).Return(platforms, nil)
	s.localPlatform.On("GetPlatformId", s.Ctx).Return("local")
	s.NoError(s.uc.OnLocalClientInfoChanged(s.Ctx, ci))
	s.AssertNumberOfCalls(&s.remoteRep.Mock, "PutClientInfoAsync", 0)
}

func (s *hubUcTestSuite) Test_OnLocalClientInfoChanged_UpdateAndRemotePut() {
	locPlatform := &domain.Platform{}
	s.localPlatform.On("Get", s.Ctx).Return(locPlatform, nil)

	ci := &domain.Party{Id: kit.NewId(), Roles: []string{domain.RoleCPO, domain.RoleEMSP}}
	ci.PlatformId = "local"
	s.partyService.On("Get", s.Ctx, ci.Id).Return(ci, nil)
	s.platformService.On("RoleEndpoint", mock.Anything, mock.Anything, model.ModuleIdHubClientInfo, model.OcpiReceiver).Return(domain.Endpoint("url"))
	platforms := []*domain.Platform{
		{TokenC: domain.PlatformToken(kit.NewRandString()), Protocol: &domain.ProtocolDetails{PushSupport: domain.PushSupport{HubClientInfo: true}}},
		{TokenC: domain.PlatformToken(kit.NewRandString()), Protocol: &domain.ProtocolDetails{PushSupport: domain.PushSupport{HubClientInfo: true}}},
	}
	s.platformService.On("Search", s.Ctx, mock.Anything).Return(platforms, nil)
	s.localPlatform.On("GetPlatformId", s.Ctx).Return("local")
	s.remoteRep.On("PutClientInfo", s.Ctx, mock.Anything).Return(nil)
	s.NoError(s.uc.OnLocalClientInfoChanged(s.Ctx, ci))
	s.AssertNumberOfCalls(&s.remoteRep.Mock, "PutClientInfo", 4)
}

func (s *hubUcTestSuite) Test_OnRemoteClientInfoPut_WhenNotOfRemotePlatform() {
	platform := &domain.Platform{Id: kit.NewId(), Status: domain.ConnectionStatusConnected}
	ci := &domain.Party{
		Id:    kit.NewId(),
		Roles: []string{domain.RoleCPO, domain.RoleEMSP},
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{PartyId: kit.NewId(), CountryCode: "RS"},
		},
	}
	ci.PlatformId = "local"
	ciOcpi := &model.OcpiClientInfo{
		OcpiPartyId: model.OcpiPartyId{
			PartyId:     ci.OcpiItem.ExtId.PartyId,
			CountryCode: ci.OcpiItem.ExtId.CountryCode,
		},
	}
	s.platformService.On("Get", s.Ctx, platform.Id).Return(platform, nil)
	s.partyService.On("GetByExtId", s.Ctx, ci.OcpiItem.ExtId).Return(ci, nil)
	s.localPlatform.On("GetPlatformId", s.Ctx).Return("local")
	s.AssertAppErr(s.uc.OnRemoteClientInfoPut(s.Ctx, platform.Id, ciOcpi), errors.ErrCodePartyNotBelongRemotePlatform)
}

func (s *hubUcTestSuite) Test_OnRemoteClientInfoPut_Ok() {
	platform := &domain.Platform{Id: kit.NewId(), Status: domain.ConnectionStatusConnected}
	ci := &domain.Party{
		Id:    kit.NewId(),
		Roles: []string{domain.RoleCPO, domain.RoleEMSP},
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{PartyId: kit.NewId(), CountryCode: "RS"},
		},
	}
	ci.PlatformId = "remote"
	ciOcpi := &model.OcpiClientInfo{
		OcpiPartyId: model.OcpiPartyId{
			PartyId:     ci.OcpiItem.ExtId.PartyId,
			CountryCode: ci.OcpiItem.ExtId.CountryCode,
		},
	}
	s.platformService.On("Get", s.Ctx, platform.Id).Return(platform, nil)
	s.partyService.On("GetByExtId", s.Ctx, ci.OcpiItem.ExtId).Return(ci, nil)
	s.localPlatform.On("GetPlatformId", s.Ctx).Return("local")
	s.webhook.On("OnPartiesChanged", mock.Anything, mock.Anything).Return(nil)
	s.partyService.On("Merge", s.Ctx, mock.Anything).Return(&domain.Party{}, nil)
	s.NoError(s.uc.OnRemoteClientInfoPut(s.Ctx, platform.Id, ciOcpi))
	s.AssertCalled(&s.partyService.Mock, "Merge", s.Ctx, mock.Anything)
	s.AssertCalled(&s.webhook.Mock, "OnPartiesChanged", mock.Anything, mock.Anything)
}

func (s *hubUcTestSuite) Test_OnRemoteClientInfoPull() {
	locPlatform := &domain.Platform{}
	s.localPlatform.On("Get", s.Ctx).Return(locPlatform, nil)

	platforms := []*domain.Platform{
		{Id: kit.NewId(), TokenC: domain.PlatformToken(kit.NewRandString())},
		{Id: kit.NewId(), TokenC: domain.PlatformToken(kit.NewRandString())},
	}
	s.localPlatform.On("GetPlatformId", mock.Anything).Return("local")
	s.platformService.On("Search", mock.Anything, mock.Anything).Return(platforms, nil)
	s.platformService.On("RoleEndpoint", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(domain.Endpoint("url"))
	s.platformService.On("Get", mock.Anything, mock.Anything).Return(&domain.Platform{Status: domain.ConnectionStatusConnected}, nil)
	s.partyService.On("GetByExtId", mock.Anything, mock.Anything).Return(&domain.Party{}, nil)
	s.partyService.On("Merge", mock.Anything, mock.Anything).Return(&domain.Party{}, nil)
	cnt := atomic.NewInt32(0)
	s.webhook.On("OnPartiesChanged", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			cnt.Inc()
		}).
		Return(nil)

	cis := []*model.OcpiClientInfo{
		{
			OcpiPartyId: model.OcpiPartyId{
				PartyId:     kit.NewRandString(),
				CountryCode: "RS",
			},
		},
		{
			OcpiPartyId: model.OcpiPartyId{
				PartyId:     kit.NewRandString(),
				CountryCode: "RS",
			},
		},
	}
	for _, p := range platforms {
		tknC := p.TokenC
		s.remoteRep.On("GetClientInfos", mock.Anything, mock.MatchedBy(func(r *usecase.OcpiRepositoryPagingRequest) bool {
			return r.Token == tknC && *r.Offset >= 500
		})).Return(nil, nil)
		s.remoteRep.On("GetClientInfos", mock.Anything, mock.MatchedBy(func(r *usecase.OcpiRepositoryPagingRequest) bool {
			return r.Token == tknC && *r.Offset < 500
		})).Return(cis, nil)
	}

	s.NoError(s.uc.OnRemoteClientInfosPull(s.Ctx, nil, nil))
	if err := <-kit.Await(func() (bool, error) {
		return cnt.Load() == 20, nil
	}, time.Millisecond*50, time.Second*50); err != nil {
		s.Fatal(err)
	}
}

func (s *hubUcTestSuite) Test_OnRemoteClientInfoPull_WhenNoClientInfo() {
	locPlatform := &domain.Platform{}
	s.localPlatform.On("Get", s.Ctx).Return(locPlatform, nil)

	platforms := []*domain.Platform{
		{Id: kit.NewId(), TokenC: domain.PlatformToken(kit.NewRandString())},
		{Id: kit.NewId(), TokenC: domain.PlatformToken(kit.NewRandString())},
	}
	s.localPlatform.On("GetPlatformId", s.Ctx).Return("local")
	s.platformService.On("Search", s.Ctx, mock.Anything).Return(platforms, nil)
	s.platformService.On("RoleEndpoint", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(domain.Endpoint("url"))
	s.platformService.On("Get", mock.Anything, mock.Anything).Return(&domain.Platform{Status: domain.ConnectionStatusConnected}, nil)
	s.partyService.On("GetByExtId", mock.Anything, mock.Anything).Return(&domain.Party{}, nil)
	cnt := atomic.NewInt32(0)
	s.remoteRep.On("GetClientInfos", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			cnt.Inc()
		}).
		Return(nil, nil)

	s.NoError(s.uc.OnRemoteClientInfosPull(s.Ctx, nil, nil))
	if err := <-kit.Await(func() (bool, error) {
		return cnt.Load() == 2, nil
	}, time.Millisecond*50, time.Second*5); err != nil {
		s.Fatal(err)
	}
}
