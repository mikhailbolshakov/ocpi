package impl

import (
	"encoding/base64"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/mocks"
	"github.com/mikhailbolshakov/ocpi/model"
	"github.com/mikhailbolshakov/ocpi/usecase"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/atomic"
	"testing"
	"time"
)

type credentialsUcTestSuite struct {
	kit.Suite
	uc                *credentialsUc
	platformSvc       *mocks.PlatformService
	localPlatformSvc  *mocks.LocalPlatformService
	tokenGen          *mocks.TokenGenerator
	remotePlatformRep *mocks.RemotePlatformRepository
	partyService      *mocks.PartyService
	webhook           *mocks.WebhookCallService
	hubUc             *mocks.HubUc
}

func (s *credentialsUcTestSuite) SetupSuite() {
	s.Suite.Init(ocpi.LF())
}

func (s *credentialsUcTestSuite) SetupTest() {
	s.platformSvc = &mocks.PlatformService{}
	s.localPlatformSvc = &mocks.LocalPlatformService{}
	s.tokenGen = &mocks.TokenGenerator{}
	s.tokenGen.On("Generate", mock.Anything).Return(domain.PlatformToken(kit.NewRandString()), nil)
	s.remotePlatformRep = &mocks.RemotePlatformRepository{}
	s.partyService = &mocks.PartyService{}
	s.webhook = &mocks.WebhookCallService{}
	s.hubUc = &mocks.HubUc{}
	s.uc = NewCredentialsUc(s.platformSvc, s.localPlatformSvc, s.tokenGen, s.remotePlatformRep, s.partyService, s.webhook, s.hubUc).(*credentialsUc)
}

func (s *credentialsUcTestSuite) TearDownSuite() {}

func TestCredentialsUcSuite(t *testing.T) {
	suite.Run(t, new(credentialsUcTestSuite))
}

func (s *credentialsUcTestSuite) platform(role string) *domain.Platform {
	return &domain.Platform{
		Id:     kit.NewRandString(),
		TokenA: domain.PlatformToken(kit.NewRandString()),
		Role:   role,
		VersionInfo: domain.VersionInfo{
			Current: "2.2.1",
			Available: map[string]domain.Endpoint{
				"2.2.1": domain.Endpoint("http://details/" + kit.NewRandString()),
			},
			VersionEp: domain.Endpoint("http://versions/" + kit.NewRandString()),
		},
		Status: domain.ConnectionStatusConnected,
	}
}

func (s *credentialsUcTestSuite) Test_Reconnect() {
	remPlatform := s.platform(domain.RoleCPO)
	locPlatform := s.platform(domain.RoleHUB)
	remEps := domain.ModuleEndpoints{
		model.ModuleIdCredentials: map[string]domain.Endpoint{
			model.OcpiSender: "http://credentials",
		},
	}
	recCred := &model.OcpiCredentials{
		Token: kit.NewRandString(),
		Url:   kit.NewRandString(),
		Roles: []*model.OcpiCredentialRole{
			{
				OcpiPartyId: model.OcpiPartyId{
					PartyId:     kit.NewRandString(),
					CountryCode: "RS",
				},
				Role: domain.RoleCPO,
			},
		},
	}
	remVer := domain.Versions{
		"2.2.1": "http://remote",
	}
	locParty := &domain.Party{Roles: []string{domain.RoleCPO}, OcpiItem: domain.OcpiItem{ExtId: domain.PartyExtId{PartyId: kit.NewId()}}}

	var postCred *model.OcpiCredentials
	var mergedParty *domain.Party
	s.localPlatformSvc.On("Get", mock.Anything).Return(locPlatform, nil)
	s.partyService.On("Search", s.Ctx, mock.AnythingOfType("*domain.PartySearchCriteria")).Return(&domain.PartySearchResponse{Items: []*domain.Party{locParty}}, nil)
	s.remotePlatformRep.On("GetVersions", s.Ctx, mock.MatchedBy(func(r *usecase.OcpiRepositoryBaseRequest) bool {
		return r.Endpoint == remPlatform.VersionInfo.VersionEp &&
			r.Token == remPlatform.TokenA &&
			r.FromPlatformId == locPlatform.Id &&
			r.ToPlatformId == remPlatform.Id
	})).Return(remVer, nil)
	s.remotePlatformRep.On("GetVersionDetails", s.Ctx, mock.MatchedBy(func(r *usecase.OcpiRepositoryBaseRequest) bool {
		return r.Token == remPlatform.TokenA
	})).Return(remEps, nil)
	s.remotePlatformRep.On("PostCredentials", s.Ctx, mock.MatchedBy(func(r *usecase.OcpiRepositoryRequestG[*model.OcpiCredentials]) bool {
		return r.Token == remPlatform.TokenA
	})).Run(func(args mock.Arguments) {
		postCred = args.Get(1).(*usecase.OcpiRepositoryRequestG[*model.OcpiCredentials]).Request
	}).
		Return(recCred, nil)
	s.partyService.On("MergeMany", s.Ctx, mock.Anything).
		Run(func(args mock.Arguments) {
			mergedParty = args.Get(1).(*domain.Party)
		}).
		Return(nil)
	s.platformSvc.On("Merge", s.Ctx, remPlatform).Return(remPlatform, nil)
	s.platformSvc.On("SetStatus", s.Ctx, remPlatform.Id, domain.ConnectionStatusConnected).Return(remPlatform, nil)
	s.webhook.On("OnPartiesChanged", s.Ctx, mock.Anything).Return(nil)

	_, err := s.uc.reconnect(s.Ctx, remPlatform)
	s.NoError(err)
	s.NotEmpty(postCred)
	s.NotEmpty(postCred.Token)
	s.Equal(postCred.Url, string(locPlatform.VersionInfo.VersionEp))
	s.NotEmpty(postCred.Roles)
	s.NotEmpty(postCred.Roles[0].PartyId, locParty.ExtId.PartyId)

	s.NotEmpty(mergedParty)
	s.Equal(mergedParty.ExtId.PartyId, recCred.Roles[0].PartyId)

}

func (s *credentialsUcTestSuite) Test_FindProperVersion() {
	type scenario struct {
		loc domain.Versions
		rem domain.Versions
		res string
	}

	for _, cs := range []*scenario{
		{
			loc: domain.Versions{},
			rem: domain.Versions{},
			res: "",
		},
		{
			loc: nil,
			rem: nil,
			res: "",
		},
		{
			loc: domain.Versions{"1.1.0": ""},
			rem: domain.Versions{"1.1.0": ""},
			res: "1.1.0",
		},
		{
			loc: domain.Versions{"1.1.0": "", "1.2.0": ""},
			rem: domain.Versions{"1.1.0": ""},
			res: "1.1.0",
		},
		{
			loc: domain.Versions{"1.1.0": "", "1.2.0": ""},
			rem: domain.Versions{"1.1": ""},
			res: "1.1.0",
		},
		{
			loc: domain.Versions{"1.1.0": "", "1.2.1": ""},
			rem: domain.Versions{"1.5": "", "1.2.1": "", "2.0": ""},
			res: "1.2.1",
		},
		{
			loc: domain.Versions{"1.1.0": "", "1.2": "", "0.0.0": ""},
			rem: domain.Versions{"1.5": "", "1.2": "", "2.0": ""},
			res: "1.2",
		},
	} {
		s.Equal(cs.res, s.uc.findProperVersion(cs.rem, cs.loc))
	}
}

func (s *credentialsUcTestSuite) Test_AcceptConnection() {
	remPlatform := s.platform(domain.RoleCPO)
	locPlatform := s.platform(domain.RoleHUB)
	remEps := domain.ModuleEndpoints{
		model.ModuleIdCredentials: map[string]domain.Endpoint{
			model.OcpiSender: "http://credentials",
		},
	}
	recCred := &model.OcpiCredentials{
		Token: kit.NewRandString(),
		Url:   kit.NewRandString(),
		Roles: []*model.OcpiCredentialRole{
			{
				OcpiPartyId: model.OcpiPartyId{
					PartyId:     kit.NewRandString(),
					CountryCode: "RS",
				},
				Role: domain.RoleCPO,
			},
		},
	}
	remVer := domain.Versions{
		"2.2.1": "http://remote/2.2.1",
		"2.0":   "http://remote/2.0",
	}
	locParty := &domain.Party{Roles: []string{domain.RoleCPO}, OcpiItem: domain.OcpiItem{ExtId: domain.PartyExtId{PartyId: kit.NewId()}}}

	var mergedParty *domain.Party
	s.localPlatformSvc.On("Get", mock.Anything).Return(locPlatform, nil)
	s.platformSvc.On("Get", s.Ctx, remPlatform.Id).Return(remPlatform, nil)
	s.partyService.On("Search", s.Ctx, mock.AnythingOfType("*domain.PartySearchCriteria")).Return(&domain.PartySearchResponse{Items: []*domain.Party{locParty}}, nil)
	s.remotePlatformRep.On("GetVersions", s.Ctx, mock.MatchedBy(func(r *usecase.OcpiRepositoryBaseRequest) bool {
		return r.Endpoint == domain.Endpoint(recCred.Url) &&
			r.Token == domain.PlatformToken(recCred.Token) &&
			r.FromPlatformId == locPlatform.Id &&
			r.ToPlatformId == remPlatform.Id
	})).Return(remVer, nil)
	s.remotePlatformRep.On("GetVersionDetails", s.Ctx, mock.MatchedBy(func(r *usecase.OcpiRepositoryBaseRequest) bool { return r.Token == domain.PlatformToken(recCred.Token) })).Return(remEps, nil)
	s.partyService.On("MergeMany", s.Ctx, mock.Anything).
		Run(func(args mock.Arguments) {
			mergedParty = args.Get(1).(*domain.Party)
		}).
		Return(nil)
	s.platformSvc.On("Merge", s.Ctx, remPlatform).Return(remPlatform, nil)
	s.platformSvc.On("SetStatus", s.Ctx, remPlatform.Id, domain.ConnectionStatusConnected).Return(remPlatform, nil)
	s.webhook.On("OnPartiesChanged", s.Ctx, mock.Anything).Return(nil)

	_, err := s.uc.AcceptConnection(s.Ctx, remPlatform.Id, recCred)
	s.NoError(err)

	s.NotEmpty(mergedParty)
	s.Equal(mergedParty.ExtId.PartyId, recCred.Roles[0].PartyId)

}

func (s *credentialsUcTestSuite) Test_OnRemotePartyPull() {

	locPlatform := s.platform(domain.RoleHUB)
	s.localPlatformSvc.On("Get", s.Ctx).Return(locPlatform, nil)

	platforms := []*domain.Platform{
		{
			Id:     kit.NewRandString(),
			TokenC: domain.PlatformToken(kit.NewRandString()),
			Status: domain.ConnectionStatusConnected,
		},
		{
			Id:     kit.NewRandString(),
			TokenC: domain.PlatformToken(kit.NewRandString()),
			Status: domain.ConnectionStatusConnected,
		},
	}
	cred1 := &model.OcpiCredentials{
		Roles: []*model.OcpiCredentialRole{
			{
				OcpiPartyId: model.OcpiPartyId{
					PartyId:     kit.NewRandString(),
					CountryCode: "RS",
				},
				Role: model.OcpiRoleCPO,
				BusinessDetails: &model.OcpiBusinessDetails{
					Name: "name",
				},
			},
		},
	}
	cred2 := &model.OcpiCredentials{
		Roles: []*model.OcpiCredentialRole{
			{
				OcpiPartyId: model.OcpiPartyId{
					PartyId:     kit.NewRandString(),
					CountryCode: "RS",
				},
				Role: model.OcpiRoleCPO,
				BusinessDetails: &model.OcpiBusinessDetails{
					Name: "name",
				},
			},
		},
	}
	var parties []*domain.Party
	s.platformSvc.On("Search", s.Ctx, mock.Anything).Return(platforms, nil)
	s.platformSvc.On("RoleEndpoint", mock.Anything, mock.Anything, model.ModuleIdCredentials, model.OcpiSender).Return(domain.Endpoint(kit.NewRandString()), nil)
	s.remotePlatformRep.On("GetCredentials", mock.Anything, mock.MatchedBy(func(r *usecase.OcpiRepositoryBaseRequest) bool { return r.Token == platforms[0].TokenC })).Return(cred1, nil)
	s.remotePlatformRep.On("GetCredentials", mock.Anything, mock.MatchedBy(func(r *usecase.OcpiRepositoryBaseRequest) bool { return r.Token == platforms[1].TokenC })).Return(cred2, nil)
	s.partyService.On("MergeMany", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			parties = append(parties, args.Get(1).(*domain.Party))
		}).
		Return(nil)
	called := 0
	s.webhook.On("OnPartiesChanged", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			called++
		}).
		Return(nil)
	s.NoError(s.uc.OnRemotePartyPull(s.Ctx))

	s.NoError(<-kit.Await(func() (bool, error) {
		return len(parties) == 2 && called == 2, nil
	}, time.Millisecond*300, time.Second))

	for _, p := range parties {
		s.NotEmpty(p.ExtId.PartyId)
		s.NotEmpty(p.ExtId.CountryCode)
		s.NotEmpty(p.Roles)
		s.NotEmpty(p.BusinessDetails)
		s.NotEmpty(p.BusinessDetails.Name)
		s.NotEmpty(p.Status)
		s.NotEmpty(p.PlatformId)
	}
}

func (s *credentialsUcTestSuite) Test_OnRemotePartyPull_WhenNotFound() {

	locPlatform := s.platform(domain.RoleHUB)
	s.localPlatformSvc.On("Get", s.Ctx).Return(locPlatform, nil)

	platforms := []*domain.Platform{
		{
			Id:     kit.NewRandString(),
			TokenC: domain.PlatformToken(kit.NewRandString()),
			Status: domain.ConnectionStatusConnected,
		},
		{
			Id:     kit.NewRandString(),
			TokenC: domain.PlatformToken(kit.NewRandString()),
			Status: domain.ConnectionStatusConnected,
		},
	}
	cred1 := &model.OcpiCredentials{}
	cred2 := &model.OcpiCredentials{}
	s.platformSvc.On("Search", s.Ctx, mock.Anything).Return(platforms, nil)
	s.platformSvc.On("RoleEndpoint", mock.Anything, mock.Anything, model.ModuleIdCredentials, model.OcpiSender).Return(domain.Endpoint(kit.NewRandString()), nil)
	v := atomic.NewInt32(0)
	s.remotePlatformRep.On("GetCredentials", mock.Anything, mock.MatchedBy(func(r *usecase.OcpiRepositoryBaseRequest) bool { return r.Token == platforms[0].TokenC })).
		Run(func(args mock.Arguments) {
			v.Inc()
		}).
		Return(cred1, nil)
	s.remotePlatformRep.On("GetCredentials", mock.Anything, mock.MatchedBy(func(r *usecase.OcpiRepositoryBaseRequest) bool { return r.Token == platforms[1].TokenC })).
		Run(func(args mock.Arguments) {
			v.Inc()
		}).
		Return(cred2, nil)
	s.NoError(s.uc.OnRemotePartyPull(s.Ctx))

	s.NoError(<-kit.Await(func() (bool, error) {
		return v.Load() == 2, nil
	}, time.Millisecond*300, time.Second))
}

func (s *credentialsUcTestSuite) base64Encode(tkn string) string {
	return base64.StdEncoding.EncodeToString([]byte(tkn))
}

func (s *credentialsUcTestSuite) base64Decode(tkn string) (string, error) {
	v, err := base64.StdEncoding.DecodeString(string(tkn))
	if err != nil {
		return "", nil
	}
	return string(v), nil
}
