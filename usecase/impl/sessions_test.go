package impl

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type sessionUcTestSuite struct {
	kit.Suite
	uc                   *sessionUc
	platformService      *mocks.PlatformService
	sessionService       *mocks.SessionService
	remoteSessionRep     *mocks.RemoteSessionRepository
	partyService         *mocks.PartyService
	webhook              *mocks.WebhookCallService
	cmdService           *mocks.CommandService
	localPlatformService *mocks.LocalPlatformService
	tokenService         *mocks.TokenService
}

func (s *sessionUcTestSuite) SetupSuite() {
	s.Suite.Init(ocpi.LF())
	s.platformService = &mocks.PlatformService{}
	s.sessionService = &mocks.SessionService{}
	s.remoteSessionRep = &mocks.RemoteSessionRepository{}
	s.partyService = &mocks.PartyService{}
	s.webhook = &mocks.WebhookCallService{}
	s.cmdService = &mocks.CommandService{}
	s.localPlatformService = &mocks.LocalPlatformService{}
	s.tokenService = &mocks.TokenService{}
	s.uc = NewSessionUc(s.platformService, s.sessionService, s.remoteSessionRep, s.partyService, s.webhook, s.cmdService, s.localPlatformService, s.tokenService, nil).(*sessionUc)
}

func (s *sessionUcTestSuite) SetupTest() {
}

func (s *sessionUcTestSuite) TearDownSuite() {}

func TestSessionUcSuite(t *testing.T) {
	suite.Run(t, new(sessionUcTestSuite))
}

func (s *sessionUcTestSuite) Test_OnLocalSessionChanged_UpdateAndRemotePut() {
	s.localPlatformService.On("Get", mock.Anything).Return(&domain.Platform{}, nil)
	sess := &domain.Session{Id: kit.NewId(), Details: domain.SessionDetails{CdrToken: &domain.CdrToken{Id: kit.NewId()}}}
	tkn := &domain.Token{OcpiItem: domain.OcpiItem{PlatformId: "remote", ExtId: domain.PartyExtId{PartyId: kit.NewRandString()}}}
	s.tokenService.On("GetToken", s.Ctx, sess.Details.CdrToken.Id).Return(tkn, nil)
	s.sessionService.On("PutSession", s.Ctx, sess).Return(sess, nil)
	s.platformService.On("RoleEndpoint", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(domain.Endpoint("url"))
	platform := &domain.Platform{TokenC: domain.PlatformToken(kit.NewRandString()), Protocol: &domain.ProtocolDetails{PushSupport: domain.PushSupport{Sessions: true}}, Status: domain.ConnectionStatusConnected}
	s.platformService.On("Get", s.Ctx, tkn.PlatformId).Return(platform, nil)
	s.remoteSessionRep.On("PutSessionAsync", s.Ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	s.NoError(s.uc.OnLocalSessionChanged(s.Ctx, sess))
	s.AssertNumberOfCalls(&s.remoteSessionRep.Mock, "PutSessionAsync", 1)
}

func (s *sessionUcTestSuite) Test_OnLocalSessionPatched_UpdateAndRemotePatch() {
	s.localPlatformService.On("Get", mock.Anything).Return(&domain.Platform{}, nil)
	sess := &domain.Session{Id: kit.NewId(), Details: domain.SessionDetails{CdrToken: &domain.CdrToken{Id: kit.NewId()}}}
	tkn := &domain.Token{OcpiItem: domain.OcpiItem{PlatformId: "remote", ExtId: domain.PartyExtId{PartyId: kit.NewRandString()}}}
	s.tokenService.On("GetToken", s.Ctx, sess.Details.CdrToken.Id).Return(tkn, nil)
	s.sessionService.On("MergeSession", s.Ctx, sess).Return(sess, nil)
	s.platformService.On("RoleEndpoint", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(domain.Endpoint("url"))
	platform := &domain.Platform{TokenC: domain.PlatformToken(kit.NewRandString()), Protocol: &domain.ProtocolDetails{PushSupport: domain.PushSupport{Sessions: true}}, Status: domain.ConnectionStatusConnected}
	s.platformService.On("Get", s.Ctx, tkn.PlatformId).Return(platform, nil)
	s.remoteSessionRep.On("PatchSessionAsync", s.Ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	s.NoError(s.uc.OnLocalSessionPatched(s.Ctx, sess))
	s.AssertNumberOfCalls(&s.remoteSessionRep.Mock, "PatchSessionAsync", 1)
}
