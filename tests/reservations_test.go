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

type reservationTestSuite struct {
	emulatorSuite
}

func (s *reservationTestSuite) SetupSuite() {
	// load config
	cfg, err := ocpi.LoadConfig()
	s.NoError(err)

	s.setCfg(cfg)
	s.init()
}

func (s *reservationTestSuite) TearDownSuite() {
	s.close()
}

func TestReservationsSuite(t *testing.T) {
	suite.Run(t, new(reservationTestSuite))
}

func (s *reservationTestSuite) Test_ReserveByRemote_LocalWebhook_RespondRejected() {

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
		Events: []string{backend.WhEventReservation},
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

	// send reservation command from the remote platform
	cmd := &backend.ReserveNowRequest{
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
		PartyId:       party.PartyId,
		CountryCode:   "RS",
		ReservationId: kit.NewId(),
		ExpireDate:    kit.Now(),
	}
	s.NoError(s.emulatorSdk.Reservation(s.Ctx, cmd))

	// await command created on the local platform
	if err := <-kit.Await(func() (bool, error) {
		return receivedEvent == backend.WhEventReservation && receivedBody != nil, nil
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

func (s *reservationTestSuite) Test_ReserveByRemote_Accepted_Cancelled() {

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
		Events: []string{backend.WhEventReservation, backend.WhEventCancelReservation},
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

	// send reservation command from the remote platform
	resCmd := &backend.ReserveNowRequest{
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
		PartyId:       party.PartyId,
		CountryCode:   "RS",
		ReservationId: kit.NewId(),
		ExpireDate:    kit.Now(),
	}
	s.NoError(s.emulatorSdk.Reservation(s.Ctx, resCmd))

	// await command created on the local platform
	if err := <-kit.Await(func() (bool, error) {
		return receivedEvent == backend.WhEventReservation && receivedBody != nil, nil
	}, time.Millisecond*300, time.Second*5); err != nil {
		s.Fatal(err)
	}
	s.NotEmpty(receivedBody["data"])
	locCmdId := receivedBody["data"].(map[string]interface{})["id"].(string)

	// set response by the remote platform
	s.NoError(s.localSdk.PutCommandResponse(s.Ctx, locCmdId, &backend.CommandResponse{
		Status: backend.CmdResultTypeAccepted,
	}))

	// await remote command is processed ok
	if err := <-kit.Await(func() (bool, error) {
		c, err := s.emulatorSdk.GetCommand(s.Ctx, resCmd.Id)
		if err != nil {
			return false, err
		}
		return c != nil && c.Status == "ok", nil
	}, time.Millisecond*300, time.Second*5); err != nil {
		s.Fatal(err)
	}

	// send cancel reservation command from the remote platform
	cancelCmd := &backend.CancelReservationRequest{
		Id:            kit.NewId(),
		PartyId:       party.PartyId,
		CountryCode:   "RS",
		ReservationId: resCmd.ReservationId,
	}
	s.NoError(s.emulatorSdk.CancelReservations(s.Ctx, cancelCmd))

	// await command created on the local platform
	if err := <-kit.Await(func() (bool, error) {
		return receivedEvent == backend.WhEventCancelReservation && receivedBody != nil, nil
	}, time.Millisecond*300, time.Second*5); err != nil {
		s.Fatal(err)
	}
	s.NotEmpty(receivedBody["data"])
	locCmdId = receivedBody["data"].(map[string]interface{})["id"].(string)

	// set response by the remote platform
	s.NoError(s.localSdk.PutCommandResponse(s.Ctx, locCmdId, &backend.CommandResponse{
		Status: backend.CmdResultTypeAccepted,
	}))

	// await remote command is processed ok
	if err := <-kit.Await(func() (bool, error) {
		c, err := s.emulatorSdk.GetCommand(s.Ctx, cancelCmd.Id)
		if err != nil {
			return false, err
		}
		return c != nil && c.Status == "ok", nil
	}, time.Millisecond*300, time.Second*5); err != nil {
		s.Fatal(err)
	}

}
