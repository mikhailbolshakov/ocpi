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
	"testing"
)

type commandUcTestSuite struct {
	kit.Suite
	uc                   usecase.CommandUc
	platformService      *mocks.PlatformService
	commandService       *mocks.CommandService
	locationService      *mocks.LocationService
	remoteCommandRep     *mocks.RemoteCommandRepository
	partyService         *mocks.PartyService
	webhook              *mocks.WebhookCallService
	localPlatformService *mocks.LocalPlatformService
	tokenUc              *mocks.TokenUc
	tokenService         *mocks.TokenService
	sessionService       *mocks.SessionService
	converter            *mocks.CommandConverter
	tokenGen             *mocks.TokenGenerator
}

func (s *commandUcTestSuite) SetupSuite() {
	s.Suite.Init(ocpi.LF())
}

func (s *commandUcTestSuite) SetupTest() {
	s.commandService = &mocks.CommandService{}
	s.locationService = &mocks.LocationService{}
	s.remoteCommandRep = &mocks.RemoteCommandRepository{}
	s.partyService = &mocks.PartyService{}
	s.webhook = &mocks.WebhookCallService{}
	s.localPlatformService = &mocks.LocalPlatformService{}
	s.tokenUc = &mocks.TokenUc{}
	s.tokenService = &mocks.TokenService{}
	s.sessionService = &mocks.SessionService{}
	s.converter = &mocks.CommandConverter{}
	s.platformService = &mocks.PlatformService{}
	s.tokenGen = &mocks.TokenGenerator{}

	s.uc = NewCommandUc(
		s.platformService,
		s.commandService,
		s.remoteCommandRep,
		s.partyService,
		s.locationService,
		s.webhook,
		s.localPlatformService,
		s.tokenUc,
		s.tokenService,
		s.sessionService,
		s.tokenGen,
	)
}

func TestCommandUcSuite(t *testing.T) {
	suite.Run(t, new(commandUcTestSuite))
}

func (s *commandUcTestSuite) Test_OnRemoteStartSession_Success() {

	rq := &model.OcpiStartSession{
		LocationId:  "location123",
		EvseId:      "evse123",
		ConnectorId: "connector123",
		Token:       &model.OcpiToken{Id: "token123"},
	}

	connector := &domain.Connector{OcpiItem: domain.OcpiItem{PlatformId: "local"}}
	s.localPlatformService.On("GetPlatformId", s.Ctx).Return("local")
	s.locationService.On("GetConnector", s.Ctx, rq.LocationId, rq.EvseId, rq.ConnectorId).Return(connector, nil)
	s.tokenUc.On("OnRemoteTokenPut", s.Ctx, "platform123", rq.Token).Return(nil)
	s.tokenService.On("GetToken", s.Ctx, "token123").Return(&domain.Token{}, nil)

	cmd := &domain.Command{
		Id:      "cmd123",
		Details: domain.CommandDetails{StartSession: &domain.StartSession{Token: &domain.Token{}}},
	}
	s.commandService.On("Create", s.Ctx, mock.Anything).Return(cmd, nil)
	s.webhook.On("OnStartSession", s.Ctx, mock.Anything).Return(nil)

	response, err := s.uc.OnRemoteStartSession(s.Ctx, "platform123", rq)
	s.NoError(err)
	s.NotNil(response)
	s.Equal(domain.CmdResponseTypeAccepted, response.Result)
}

func (s *commandUcTestSuite) Test_OnRemoteStartSession_ConnectorNotFound() {
	rq := &model.OcpiStartSession{
		LocationId:  "location123",
		EvseId:      "evse123",
		ConnectorId: "connector123",
	}

	s.locationService.On("GetConnector", s.Ctx, rq.LocationId, rq.EvseId, rq.ConnectorId).Return(nil, nil)
	_, err := s.uc.OnRemoteStartSession(s.Ctx, "platform123", rq)
	s.AssertAppErr(err, errors.ErrCodeCmdConNotFound)
}

func (s *commandUcTestSuite) Test_OnRemoteStartSession_LocationNotBelongToLocalPlatform() {
	rq := &model.OcpiStartSession{
		LocationId:  "location123",
		EvseId:      "evse123",
		ConnectorId: "connector123",
	}

	connector := &domain.Connector{OcpiItem: domain.OcpiItem{PlatformId: "remote"}}
	s.locationService.On("GetConnector", s.Ctx, rq.LocationId, rq.EvseId, rq.ConnectorId).Return(connector, nil)
	s.localPlatformService.On("GetPlatformId", s.Ctx).Return("local")

	_, err := s.uc.OnRemoteStartSession(s.Ctx, "platform123", rq)
	s.AssertAppErr(err, errors.ErrCodeLocNotBelongLocalPlatform)
}

func (s *commandUcTestSuite) Test_OnRemoteStartSession_TokenError() {
	rq := &model.OcpiStartSession{
		LocationId:  "location123",
		EvseId:      "evse123",
		ConnectorId: "connector123",
		Token:       &model.OcpiToken{Id: "token123"},
	}

	connector := &domain.Connector{OcpiItem: domain.OcpiItem{PlatformId: "local"}}
	s.localPlatformService.On("GetPlatformId", s.Ctx).Return("local")
	s.locationService.On("GetConnector", s.Ctx, rq.LocationId, rq.EvseId, rq.ConnectorId).Return(connector, nil)
	s.tokenUc.On("OnRemoteTokenPut", s.Ctx, "platform123", rq.Token).Return(errors.ErrTknNotValid(s.Ctx))

	_, err := s.uc.OnRemoteStartSession(s.Ctx, "platform123", rq)
	s.AssertAppErr(err, errors.ErrCodeTknNotValid)
}
