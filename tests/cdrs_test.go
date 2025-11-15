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

type cdrsTestSuite struct {
	emulatorSuite
}

func (s *cdrsTestSuite) SetupSuite() {
	// load config
	cfg, err := ocpi.LoadConfig()
	s.NoError(err)

	s.setCfg(cfg)
	s.init()

}

func (s *cdrsTestSuite) TearDownSuite() {
	s.close()
}

func TestCdrsSuite(t *testing.T) {
	suite.Run(t, new(cdrsTestSuite))
}

func (s *cdrsTestSuite) Test_Cdr_SentByLocal_AfterSession() {

	// create webhook server
	var receivedEvent string
	var receivedBody map[string]interface{}
	server := s.openWebhookServer("7890", "/test/webhook", func(ev string, body map[string]interface{}) {
		receivedEvent = ev
		receivedBody = body
	})
	defer server.Close()

	// register webhook on remote platform
	whR := &backend.Webhook{
		Id:     kit.NewRandString(),
		ApiKey: kit.NewRandString(),
		Events: []string{backend.WhEventStartSession, backend.WhEventStopSession},
		Url:    s.cfg.Tests.WebhookUrl,
	}
	whR, err := s.emulatorSdk.CreateUpdateWebhook(s.Ctx, whR)
	s.NoError(err)
	defer func() { _ = s.emulatorSdk.DeleteWebhook(s.Ctx, whR.Id) }()

	// register webhook on local platform
	whL := &backend.Webhook{
		Id:     kit.NewRandString(),
		ApiKey: kit.NewRandString(),
		Events: []string{backend.WhEventCdrChanged},
		Url:    s.cfg.Tests.WebhookUrl,
	}
	whL, err = s.localSdk.CreateUpdateWebhook(s.Ctx, whL)
	s.NoError(err)
	defer func() { _ = s.emulatorSdk.DeleteWebhook(s.Ctx, whL.Id) }()

	party := party()

	// create tariff
	trf := tariff(party.PartyId)
	s.NoError(s.emulatorSdk.PutTariff(s.Ctx, trf))

	// await tariff in emulator
	_, err = awaitTariffs(s.Ctx, s.localSdk, trf.Id, nil)
	s.NoError(err)

	// create and send a new location in the remote system
	loc := location(party.PartyId)
	loc.Evses[0].Connectors[0].TariffIds = []string{trf.Id}
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

	// await local command is processed ok
	_, err = awaitCommand(s.Ctx, s.localSdk, startRq.Id, "ok")

	// send session
	sess := session(party.PartyId, startRq.AuthorizationReference, loc.Id, loc.Evses[0].Id, loc.Evses[0].Connectors[0].Id)
	sess.CdrToken.Id = startRq.Token.Id
	sess.CdrToken.ContractId = startRq.Token.ContractId
	sess.CdrToken.Type = startRq.Token.Type
	sess.CdrToken.PartyId = startRq.Token.PartyId
	sess.CdrToken.CountryCode = startRq.Token.CountryCode
	s.NoError(s.emulatorSdk.PutSession(s.Ctx, sess))

	// await local session
	_, err = awaitSession(s.Ctx, s.localSdk, sess.Id)

	// send MV
	kwh := kit.Float64Ptr(*sess.Kwh + 1000.0)
	sess.ChargingPeriods = []*backend.ChargingPeriod{
		{
			StartDateTime: kit.Now(),
			Dimensions: []*backend.CdrDimension{
				{
					Type:   backend.DimensionTypeEnergyImport,
					Volume: *kwh,
				},
			},
			TariffId: trf.Id,
		},
	}
	sess.LastUpdated = kit.Now()
	s.NoError(s.emulatorSdk.PatchSession(s.Ctx, sess))

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
	_, err = awaitCommand(s.Ctx, s.localSdk, stopRq.Id, "ok")

	// put session status
	sess.EndDateTime = kit.NowPtr()
	sess.Status = backend.SessionStatusCompleted
	sess.LastUpdated = kit.Now()
	s.NoError(s.emulatorSdk.PatchSession(s.Ctx, sess))

	// put cdr
	cdr := cdr(party.PartyId, sess, []*backend.Tariff{trf})
	s.NoError(s.emulatorSdk.PutCdr(s.Ctx, cdr))

	// await webhook
	if err := <-kit.Await(func() (bool, error) {
		return receivedEvent == backend.WhEventCdrChanged && receivedBody != nil, nil
	}, time.Millisecond*300, time.Second*5); err != nil {
		s.Fatal(err)
	}
	s.NotEmpty(receivedBody["data"])
	cdrIdWh := receivedBody["data"].(map[string]interface{})["id"].(string)

	// await cdr on local platform
	locCdr, err := awaitCdr(s.Ctx, s.localSdk, cdr.Id)
	s.NoError(err)
	s.NotEmpty(locCdr)
	s.Equal(locCdr.Id, cdr.Id)
	s.Equal(locCdr.Id, cdrIdWh)
	s.Equal(locCdr.PartyId, cdr.PartyId)
	s.Equal(sess.ChargingPeriods, cdr.ChargingPeriods)

}
