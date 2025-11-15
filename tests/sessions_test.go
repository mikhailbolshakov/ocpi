//go:build dev

package tests

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type sessionsTestSuite struct {
	emulatorSuite
}

func (s *sessionsTestSuite) SetupSuite() {
	// load config
	cfg, err := ocpi.LoadConfig()
	s.NoError(err)

	s.setCfg(cfg)
	s.init()
}

func (s *sessionsTestSuite) TearDownSuite() {
	s.close()
}

func TestSessionsSuite(t *testing.T) {
	suite.Run(t, new(sessionsTestSuite))
}

func (s *sessionsTestSuite) Test_StartSessionByLocal_RemoteWebhook_RespondAccepted() {

	// create webhook server
	var receivedEvent string
	var receivedBody map[string]interface{}
	server := s.openWebhookServer("7890", "/test/webhook", func(ev string, body map[string]interface{}) {
		receivedEvent = ev
		receivedBody = body
	})
	defer server.Close()

	// register webhook
	wh := &backend.Webhook{
		Id:     kit.NewRandString(),
		ApiKey: kit.NewRandString(),
		Events: []string{backend.WhEventStartSession},
		Url:    s.cfg.Tests.WebhookUrl,
	}
	wh, err := s.emulatorSdk.CreateUpdateWebhook(s.Ctx, wh)
	s.NoError(err)
	defer func() { _ = s.emulatorSdk.DeleteWebhook(s.Ctx, wh.Id) }()

	// create and send a new location in the remote system
	party := party()
	loc := location(party.PartyId)
	s.NoError(s.emulatorSdk.PutLocation(s.Ctx, loc))

	// await the location pushed to the local platform
	loc, err = awaitLocation(s.Ctx, s.localSdk, loc.Id, nil)
	s.NoError(err)

	// send start session command from the local platform
	cmd := &backend.StartSessionRequest{
		Id:                     kit.NewId(),
		LocationId:             loc.Id,
		EvseId:                 loc.Evses[0].Id,
		ConnectorId:            loc.Evses[0].Connectors[0].Id,
		AuthorizationReference: kit.NewRandString(),
		Token: backend.Token{
			Id:           kit.NewId(),
			PartyId:      party.PartyId,
			CountryCode:  "RS",
			VisualNumber: kit.NewRandString(),
			ContractId:   kit.NewRandString(),
			Issuer:       "some",
			WhiteList:    backend.TokenWLTypeAlways,
			Valid:        kit.BoolPtr(true),
			Type:         backend.TokenTypeAppUser,
			Lang:         "en",
			LastUpdated:  kit.Now(),
		},
		PartyId:     party.PartyId,
		CountryCode: "RS",
	}
	s.NoError(s.localSdk.StartSession(s.Ctx, cmd))

	// await command created on the remote platform
	if err := <-kit.Await(func() (bool, error) {
		return receivedEvent == backend.WhEventStartSession && receivedBody != nil, nil
	}, time.Millisecond*300, time.Second*5); err != nil {
		s.Fatal(err)
	}
	s.NotEmpty(receivedBody["data"])
	remoteCmdId := receivedBody["data"].(map[string]interface{})["id"].(string)

	// set response by the remote platform
	s.NoError(s.emulatorSdk.PutCommandResponse(s.Ctx, remoteCmdId, &backend.CommandResponse{
		Status: backend.CmdResultTypeAccepted,
	}))

	// await local command is processed ok-
	if err := <-kit.Await(func() (bool, error) {
		c, err := s.localSdk.GetCommand(s.Ctx, cmd.Id)
		if err != nil {
			return false, err
		}
		return c != nil && c.Status == "ok", nil
	}, time.Millisecond*300, time.Second*5); err != nil {
		s.Fatal(err)
	}
}

func (s *sessionsTestSuite) Test_StartSessionByRemote_LocalWebhook_RespondRejected() {

	// create webhook server
	var receivedEvent string
	var receivedBody map[string]interface{}
	server := s.openWebhookServer("7890", "/test/webhook", func(ev string, body map[string]interface{}) {
		receivedEvent = ev
		receivedBody = body
	})
	defer server.Close()

	// register webhook
	wh := &backend.Webhook{
		Id:     kit.NewRandString(),
		ApiKey: kit.NewRandString(),
		Events: []string{backend.WhEventStartSession},
		Url:    s.cfg.Tests.WebhookUrl,
	}
	wh, err := s.localSdk.CreateUpdateWebhook(s.Ctx, wh)
	s.NoError(err)
	defer func() { _ = s.localSdk.DeleteWebhook(s.Ctx, wh.Id) }()

	// create and send a new location in the local platform
	party := party()
	loc := location(party.PartyId)
	s.NoError(s.localSdk.PutLocation(s.Ctx, loc))

	// await the location pushed to the remote platform
	loc, err = awaitLocation(s.Ctx, s.emulatorSdk, loc.Id, nil)
	s.NoError(err)

	// send start session command from the remote platform
	cmd := &backend.StartSessionRequest{
		Id:                     kit.NewId(),
		LocationId:             loc.Id,
		EvseId:                 loc.Evses[0].Id,
		ConnectorId:            loc.Evses[0].Connectors[0].Id,
		AuthorizationReference: kit.NewRandString(),
		Token: backend.Token{
			Id:           kit.NewId(),
			PartyId:      party.PartyId,
			CountryCode:  "RS",
			VisualNumber: kit.NewRandString(),
			ContractId:   kit.NewRandString(),
			Issuer:       "some",
			WhiteList:    backend.TokenWLTypeAlways,
			Valid:        kit.BoolPtr(true),
			Type:         backend.TokenTypeAppUser,
			Lang:         "en",
			LastUpdated:  kit.Now(),
		},
		PartyId:     party.PartyId,
		CountryCode: "RS",
	}
	s.NoError(s.emulatorSdk.StartSession(s.Ctx, cmd))

	// await command created on the local platform
	if err := <-kit.Await(func() (bool, error) {
		return receivedEvent == backend.WhEventStartSession && receivedBody != nil, nil
	}, time.Millisecond*300, time.Second*5); err != nil {
		s.Fatal(err)
	}
	s.NotEmpty(receivedBody["data"])
	locCmdId := receivedBody["data"].(map[string]interface{})["id"].(string)

	// set response by the remote platform
	s.NoError(s.localSdk.PutCommandResponse(s.Ctx, locCmdId, &backend.CommandResponse{
		Status: backend.CmdResultTypeFailed,
	}))

	// await remote command is processed ok
	if err := <-kit.Await(func() (bool, error) {
		c, err := s.emulatorSdk.GetCommand(s.Ctx, cmd.Id)
		if err != nil {
			return false, err
		}
		return c != nil && c.Status == "failed", nil
	}, time.Millisecond*300, time.Second*5); err != nil {
		s.Fatal(err)
	}
}

func (s *sessionsTestSuite) Test_StartStopSessionByLocal_RemoteSession_Completed() {

	// create webhook server
	var receivedEvent string
	var receivedBody map[string]interface{}
	server := s.openWebhookServer("7890", "/test/webhook", func(ev string, body map[string]interface{}) {
		receivedEvent = ev
		receivedBody = body
	})
	defer server.Close()

	// register webhook
	wh := &backend.Webhook{
		Id:     kit.NewRandString(),
		ApiKey: kit.NewRandString(),
		Events: []string{backend.WhEventStartSession, backend.WhEventStopSession},
		Url:    s.cfg.Tests.WebhookUrl,
	}
	wh, err := s.emulatorSdk.CreateUpdateWebhook(s.Ctx, wh)
	s.NoError(err)
	defer func() { _ = s.emulatorSdk.DeleteWebhook(s.Ctx, wh.Id) }()

	// create and send a new location in the remote system
	party := party()
	loc := location(party.PartyId)
	s.NoError(s.emulatorSdk.PutLocation(s.Ctx, loc))

	// await the location pushed to the local platform
	loc, err = awaitLocation(s.Ctx, s.localSdk, loc.Id, nil)
	s.NoError(err)

	// send start session command from the local platform
	startRq := &backend.StartSessionRequest{
		Id:                     kit.NewId(),
		LocationId:             loc.Id,
		EvseId:                 loc.Evses[0].Id,
		ConnectorId:            loc.Evses[0].Connectors[0].Id,
		AuthorizationReference: kit.NewRandString(),
		Token: backend.Token{
			Id:           kit.NewId(),
			PartyId:      party.PartyId,
			CountryCode:  "RS",
			VisualNumber: kit.NewRandString(),
			ContractId:   kit.NewRandString(),
			Issuer:       "some",
			WhiteList:    backend.TokenWLTypeAlways,
			Valid:        kit.BoolPtr(true),
			Type:         backend.TokenTypeAppUser,
			Lang:         "en",
			LastUpdated:  kit.Now(),
		},
		PartyId:     party.PartyId,
		CountryCode: "RS",
	}
	s.NoError(s.localSdk.StartSession(s.Ctx, startRq))

	// await command created on the remote platform
	if err := <-kit.Await(func() (bool, error) {
		return receivedEvent == backend.WhEventStartSession && receivedBody != nil, nil
	}, time.Millisecond*300, time.Second*5); err != nil {
		s.Fatal(err)
	}
	s.NotEmpty(receivedBody["data"])
	remoteCmdId := receivedBody["data"].(map[string]interface{})["id"].(string)

	// set response by the remote platform
	s.NoError(s.emulatorSdk.PutCommandResponse(s.Ctx, remoteCmdId, &backend.CommandResponse{
		Status: backend.CmdResultTypeAccepted,
	}))

	// await local command is processed ok-
	if err := <-kit.Await(func() (bool, error) {
		c, err := s.localSdk.GetCommand(s.Ctx, startRq.Id)
		if err != nil {
			return false, err
		}
		return c != nil && c.Status == "ok", nil
	}, time.Millisecond*300, time.Second*5); err != nil {
		s.Fatal(err)
	}

	// create session
	sess := session(party.PartyId, startRq.AuthorizationReference, loc.Id, loc.Evses[0].Id, loc.Evses[0].Connectors[0].Id)
	sess.CdrToken.Id = startRq.Token.Id
	sess.CdrToken.ContractId = startRq.Token.ContractId
	sess.CdrToken.Type = startRq.Token.Type
	sess.CdrToken.PartyId = startRq.Token.PartyId
	sess.CdrToken.CountryCode = startRq.Token.CountryCode
	s.NoError(s.emulatorSdk.PutSession(s.Ctx, sess))

	// await local session
	if err := <-kit.Await(func() (bool, error) {
		ss, err := s.localSdk.GetSession(s.Ctx, sess.Id)
		if err != nil {
			return false, err
		}
		return ss != nil, nil
	}, time.Millisecond*300, time.Second*5); err != nil {
		s.Fatal(err)
	}

	// send MV
	for i := 0; i < 10; i++ {
		kwh := kit.Float64Ptr(*sess.Kwh + float64(i)*1000.0)
		// send MV
		sessMv := &backend.Session{
			Id:            sess.Id,
			StartDateTime: sess.StartDateTime,
			Kwh:           kwh,
			AuthRef:       startRq.AuthorizationReference,
			ChargingPeriods: []*backend.ChargingPeriod{
				{
					StartDateTime: kit.Now(),
					Dimensions: []*backend.CdrDimension{
						{
							Type:   backend.DimensionTypeEnergyImport,
							Volume: float64(*kwh),
						},
					},
				},
			},
			LastUpdated: kit.Now(),
		}
		s.NoError(s.emulatorSdk.PatchSession(s.Ctx, sessMv))
		time.Sleep(time.Millisecond * 500)
	}

	// send stop session command from the local platform
	stopRq := &backend.StopSessionRequest{
		Id:          kit.NewRandString(),
		SessionId:   sess.Id,
		PartyId:     party.PartyId,
		CountryCode: "RS",
	}
	s.NoError(s.localSdk.StopSession(s.Ctx, stopRq))

	// await command created on the remote platform
	if err := <-kit.Await(func() (bool, error) {
		return receivedEvent == backend.WhEventStopSession && receivedBody != nil, nil
	}, time.Millisecond*300, time.Second*5); err != nil {
		s.Fatal(err)
	}
	s.NotEmpty(receivedBody["data"])
	remoteCmdId = receivedBody["data"].(map[string]interface{})["id"].(string)

	// set response by the remote platform
	s.NoError(s.emulatorSdk.PutCommandResponse(s.Ctx, remoteCmdId, &backend.CommandResponse{
		Status: backend.CmdResultTypeAccepted,
	}))

	// await local command is processed ok-
	if err := <-kit.Await(func() (bool, error) {
		c, err := s.localSdk.GetCommand(s.Ctx, stopRq.Id)
		if err != nil {
			return false, err
		}
		return c != nil && c.Status == "ok", nil
	}, time.Millisecond*300, time.Second*5); err != nil {
		s.Fatal(err)
	}

	// put session status
	sessMv := &backend.Session{
		Id:            sess.Id,
		StartDateTime: sess.StartDateTime,
		EndDateTime:   kit.NowPtr(),
		AuthRef:       startRq.AuthorizationReference,
		LastUpdated:   kit.Now(),
		Status:        backend.SessionStatusCompleted,
	}
	s.NoError(s.emulatorSdk.PatchSession(s.Ctx, sessMv))

	// await local session completed
	if err := <-kit.Await(func() (bool, error) {
		ss, err := s.localSdk.GetSession(s.Ctx, sess.Id)
		if err != nil {
			return false, err
		}
		return ss != nil && ss.Status == backend.SessionStatusCompleted, nil
	}, time.Millisecond*300, time.Second*5); err != nil {
		s.Fatal(err)
	}

}
