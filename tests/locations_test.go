//go:build dev

package tests

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/sdk"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type locationsTestSuite struct {
	emulatorSuite
}

func (s *locationsTestSuite) SetupSuite() {
	// load config
	cfg, err := ocpi.LoadConfig()
	s.NoError(err)

	s.setCfg(cfg)
	s.init()
}

func (s *locationsTestSuite) TearDownSuite() {
	s.close()
}

func TestLocationsSuite(t *testing.T) {
	suite.Run(t, new(locationsTestSuite))
}

func (s *locationsTestSuite) Test_Location_PutLocInLocal_PushedToEmulator() {

	// create and send a new location in local system
	party := party()
	loc := location(party.PartyId)
	s.NoError(s.localSdk.PutLocation(s.Ctx, loc))

	// await
	expected := s.awaitLocation(s.emulatorSdk, loc.Id, nil)
	s.Equal(expected.Name, loc.Name)
	s.Equal(expected.PartyId, loc.PartyId)
	s.Equal(expected.CountryCode, loc.CountryCode)
	s.Equal(len(expected.Evses), len(loc.Evses))
	s.Equal(expected.Evses[0].Status, loc.Evses[0].Status)
	s.Equal(expected.Owner.Inn, loc.Owner.Inn)

	// update
	loc.Evses[0].Status = "INOPERATIVE"
	loc.Name = "Another name"
	s.NoError(s.localSdk.PutLocation(s.Ctx, loc))
	expected = s.awaitLocation(s.emulatorSdk, loc.Id, func(location *backend.Location) bool {
		return location.Evses[0].Status == loc.Evses[0].Status
	})
	s.Equal(expected.Name, loc.Name)

}

func (s *locationsTestSuite) Test_Location_PutLocEvseConnectorBySeparateCallsInLocal_PushedToEmulator() {

	// create and send a new location without evse in local system
	party := party()
	loc := location(party.PartyId)
	evse := loc.Evses[0]
	loc.Evses = nil
	s.NoError(s.localSdk.PutLocation(s.Ctx, loc))

	// await
	s.NotEmpty(s.awaitLocation(s.emulatorSdk, loc.Id, nil))

	// put evse without connector
	con := evse.Connectors[0]
	evse.Connectors = nil
	s.NoError(s.localSdk.PutEvse(s.Ctx, evse))
	// await
	s.NotEmpty(s.awaitLocation(s.emulatorSdk, loc.Id, func(loc *backend.Location) bool {
		return len(loc.Evses) > 0
	}))
	// get evse
	s.NotEmpty(s.emulatorSdk.GetEvse(s.Ctx, loc.Id, evse.Id))

	// put connector
	s.NoError(s.localSdk.PutConnector(s.Ctx, con))
	// await
	s.NotEmpty(s.awaitLocation(s.emulatorSdk, loc.Id, func(loc *backend.Location) bool {
		return len(loc.Evses) > 0 && len(loc.Evses[0].Connectors) > 0
	}))
	// get connector
	s.NotEmpty(s.emulatorSdk.GetConnector(s.Ctx, loc.Id, evse.Id, con.Id))

	// update evse only
	evse.Status = "INOPERATIVE"
	s.NoError(s.localSdk.PutEvse(s.Ctx, evse))
	// await
	s.NotEmpty(s.awaitLocation(s.emulatorSdk, loc.Id, func(loc *backend.Location) bool {
		return len(loc.Evses) > 0 && loc.Evses[0].Status == "INOPERATIVE"
	}))

	// update connector only
	con.Standard = "GBT_AC"
	s.NoError(s.localSdk.PutConnector(s.Ctx, con))
	// await
	s.NotEmpty(s.awaitLocation(s.emulatorSdk, loc.Id, func(loc *backend.Location) bool {
		return len(loc.Evses) > 0 && len(loc.Evses[0].Connectors) > 0 && loc.Evses[0].Connectors[0].Standard == "GBT_AC"
	}))

}

func (s *locationsTestSuite) Test_Location_PutLocLocal_Webhook() {

	var receivedEvent string
	var receivedBody map[string]interface{}
	// create webhook server
	server := s.openWebhookServer("7890", "/test/webhook", func(ev string, body map[string]interface{}) {
		receivedEvent = ev
		receivedBody = body
	})
	defer server.Close()

	// register webhook
	wh := &backend.Webhook{
		Id:     kit.NewRandString(),
		ApiKey: kit.NewRandString(),
		Events: []string{backend.WhEventLocationChanged},
		Url:    s.cfg.Tests.WebhookUrl,
	}
	wh, err := s.localSdk.CreateUpdateWebhook(s.Ctx, wh)
	s.NoError(err)
	defer func() { _ = s.localSdk.DeleteWebhook(s.Ctx, wh.Id) }()

	// create and send a new party in local system
	party := party()
	loc := location(party.PartyId)
	s.NoError(s.emulatorSdk.PutLocation(s.Ctx, loc))

	// await event and payload
	if err := <-kit.Await(func() (bool, error) {
		return receivedEvent == backend.WhEventLocationChanged && receivedBody != nil, nil
	}, time.Millisecond*300, time.Second*5); err != nil {
		s.Fatal(err)
	}

	s.NotEmpty(receivedBody["data"])
	locations := receivedBody["data"].([]any)
	s.Len(locations, 1)
	s.Equal(locations[0].(map[string]any)["id"].(string), loc.Id)
	s.NotEmpty(loc.Owner)
	s.NotEmpty(loc.Owner.Inn)

}

func (s *locationsTestSuite) awaitLocation(sdk *sdk.Sdk, locId string, fn func(location *backend.Location) bool) *backend.Location {
	l, err := awaitLocation(s.Ctx, sdk, locId, fn)
	if err != nil {
		s.Fatal(err)
	}
	return l
}
