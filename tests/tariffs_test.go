//go:build dev

package tests

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/sdk"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type tariffsTestSuite struct {
	emulatorSuite
}

func (s *tariffsTestSuite) SetupSuite() {
	// load config
	cfg, err := ocpi.LoadConfig()
	s.NoError(err)

	s.setCfg(cfg)
	s.init()
}

func (s *tariffsTestSuite) TearDownSuite() {
	s.close()
}

func TestTariffsSuite(t *testing.T) {
	suite.Run(t, new(tariffsTestSuite))
}

func (s *tariffsTestSuite) Test_Tariffs_PutInLocal_PushedToEmulator() {

	// create and send a new tariff in local system
	party := party()
	trf := tariff(party.PartyId)
	s.NoError(s.localSdk.PutTariff(s.Ctx, trf))

	// await
	expected := s.mustAwaitTariffs(s.Ctx, s.emulatorSdk, trf.Id, nil)
	s.Equal(expected.Type, trf.Type)
	s.Equal(expected.PartyId, trf.PartyId)
	s.Equal(expected.CountryCode, trf.CountryCode)
	s.Equal(len(expected.Elements), len(trf.Elements))
	s.Equal(expected.Elements[0].PriceComponents[0].Price, trf.Elements[0].PriceComponents[0].Price)
	s.Equal(expected.Elements[0].PriceComponents[0].Type, trf.Elements[0].PriceComponents[0].Type)
	s.EqualValues(expected.Elements[0].Restrictions.MinKwh, trf.Elements[0].Restrictions.MinKwh)

	// update tariff
	trf.LastUpdated = kit.Now()
	trf.Elements[0].PriceComponents[0].Price = 9999.9
	s.NoError(s.localSdk.PutTariff(s.Ctx, trf))

	// await
	_ = s.mustAwaitTariffs(s.Ctx, s.emulatorSdk, trf.Id, func(t *backend.Tariff) bool {
		return t.Elements[0].PriceComponents[0].Price == trf.Elements[0].PriceComponents[0].Price
	})

	// search tariffs
	srchRs, err := s.emulatorSdk.SearchTariffs(s.Ctx, map[string]any{
		"date_from": kit.Now().Add(-10 * time.Second),
		"date_to":   kit.Now().Add(10 * time.Second),
	})
	s.NoError(err)
	s.NotEmpty(srchRs.PageInfo)
	s.NotEmpty(srchRs.PageInfo.Total)
	s.NotEmpty(kit.Filter(srchRs.Items, func(t *backend.Tariff) bool { return t.Id == trf.Id }))

}

func (s *tariffsTestSuite) Test_Tariffs_PutLocal_Webhook() {

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
		Events: []string{backend.WhEventTariffChanged},
		Url:    "http://localhost:7890/test/webhook",
	}
	wh, err := s.localSdk.CreateUpdateWebhook(s.Ctx, wh)
	s.NoError(err)
	defer func() { _ = s.localSdk.DeleteWebhook(s.Ctx, wh.Id) }()

	// create in emulator
	party := party()
	trf := tariff(party.PartyId)
	s.NoError(s.emulatorSdk.PutTariff(s.Ctx, trf))

	// await event and payload
	if err := <-kit.Await(func() (bool, error) {
		return receivedEvent == backend.WhEventTariffChanged && receivedBody != nil, nil
	}, time.Millisecond*300, time.Second*5); err != nil {
		s.Fatal(err)
	}

	s.NotEmpty(receivedBody["data"])
	tariffs := receivedBody["data"].([]any)
	s.Len(tariffs, 1)
	s.Equal(tariffs[0].(map[string]any)["id"].(string), trf.Id)

}

func (s *tariffsTestSuite) mustAwaitTariffs(ctx context.Context, sdk *sdk.Sdk, trfId string, fn func(trf *backend.Tariff) bool) *backend.Tariff {
	l, err := awaitTariffs(ctx, sdk, trfId, fn)
	s.NoError(err)
	return l
}
