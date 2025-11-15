//go:build dev

package tests

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
	sdk2 "github.com/mikhailbolshakov/ocpi/sdk"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type credentialsTestSuite struct {
	emulatorSuite
}

func (s *credentialsTestSuite) SetupSuite() {
	// load config
	cfg, err := ocpi.LoadConfig()
	s.NoError(err)

	s.setCfg(cfg)
	s.init()

}

func (s *credentialsTestSuite) TearDownSuite() {
	s.close()
}

func TestCredentialsSuite(t *testing.T) {
	suite.Run(t, new(credentialsTestSuite))
}

func (s *credentialsTestSuite) Test_Hub_NewPartyInLocal_PushedToEmulator() {

	// create and send a new party in local system
	party := party()
	party.PartyId = "TS1"
	err := s.localSdk.PutParty(s.Ctx, party)
	s.NoError(err)

	// await
	expected := s.expectParty(s.emulatorSdk, party.CountryCode, party.PartyId, nil)
	s.Equal(expected.CountryCode, party.CountryCode)
	s.Equal(expected.Status, party.Status)

	// delete party
	s.NoError(s.localSdk.DeletePartyByExt(s.Ctx, party.PartyId, party.CountryCode))

}

func (s *credentialsTestSuite) Test_Hub_PartyStatusChangedInLocal_PushedToEmulator() {

	// create and send a new party in local system
	party := party()
	party.PartyId = "TS2"
	err := s.localSdk.PutParty(s.Ctx, party)
	s.NoError(err)

	// await
	expected := s.expectParty(s.emulatorSdk, party.CountryCode, party.PartyId, func(p *backend.Party) bool {
		return p.Status == backend.ConnectionStatusConnected
	})
	s.Equal(expected.CountryCode, party.CountryCode)
	s.Equal(expected.Status, party.Status)

	party.Status = backend.ConnectionStatusSuspended
	err = s.localSdk.PutParty(s.Ctx, party)
	s.NoError(err)

	// await
	expected = s.expectParty(s.emulatorSdk, party.CountryCode, party.PartyId, func(p *backend.Party) bool {
		return p.Status == backend.ConnectionStatusSuspended
	})
	s.Equal(expected.CountryCode, party.CountryCode)
	s.Equal(expected.Status, party.Status)

	// delete party
	s.NoError(s.localSdk.DeletePartyByExt(s.Ctx, party.PartyId, party.CountryCode))

}

func (s *credentialsTestSuite) Test_PartyChangedInLocal_PullInEmulator() {

	// create and send a new party in local system
	party := party()
	party.PartyId = "TS3"
	err := s.localSdk.PutParty(s.Ctx, party)
	s.NoError(err)

	// call pull parties command
	s.NoError(s.emulatorSdk.PullParties(s.Ctx))

	// await
	expected := s.expectParty(s.emulatorSdk, party.CountryCode, party.PartyId, nil)
	s.Equal(expected.CountryCode, party.CountryCode)
	s.Equal(expected.Status, party.Status)

	// delete party
	s.NoError(s.localSdk.DeletePartyByExt(s.Ctx, party.PartyId, party.CountryCode))

}

func (s *credentialsTestSuite) Test_Delete_CantOfRemotePlatform() {

	// create and send a new party in remote system
	party := party()
	party.PartyId = "TS3"
	err := s.localSdk.PutParty(s.Ctx, party)
	s.NoError(err)

	// call pull parties command
	s.NoError(s.emulatorSdk.PullParties(s.Ctx))

	// await
	expected := s.expectParty(s.emulatorSdk, party.CountryCode, party.PartyId, nil)
	s.Equal(expected.CountryCode, party.CountryCode)
	s.Equal(expected.Status, party.Status)

	// delete party
	s.Error(s.emulatorSdk.DeletePartyByExt(s.Ctx, party.PartyId, party.CountryCode))
	s.NoError(s.localSdk.DeletePartyByExt(s.Ctx, party.PartyId, party.CountryCode))

}

func (s *credentialsTestSuite) expectParty(sdk *sdk2.Sdk, countryCode, partyId string, fn func(party *backend.Party) bool) *backend.Party {
	var party *backend.Party
	if err := <-kit.Await(func() (bool, error) {
		var err error
		rs, err := sdk.SearchParties(s.Ctx, map[string]interface{}{"countryCode": countryCode, "partyId": partyId})
		if err != nil {
			return false, err
		}
		if len(rs.Items) > 0 {
			party = rs.Items[0]
		} else {
			return false, nil
		}
		if fn != nil {
			return fn(party), nil
		}
		return true, nil
	}, time.Millisecond*300, time.Second*3); err != nil {
		s.Fatal(err)
	}
	return party
}
