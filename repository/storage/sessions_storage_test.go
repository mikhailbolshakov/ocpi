//go:build integration

package storage

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type sessionsTestSuite struct {
	kit.Suite
	storage domain.SessionStorage
	adapter Adapter
}

func (s *sessionsTestSuite) SetupSuite() {
	s.Suite.Init(ocpi.LF())

	// load config
	cfg, err := ocpi.LoadConfig()
	if err != nil {
		s.Fatal(err)
	}

	s.adapter = NewAdapter()
	s.NoError(s.adapter.Init(s.Ctx, cfg.Storages))

	s.storage = s.adapter
}

func (s *sessionsTestSuite) TearDownSuite() {
	_ = s.adapter.Close(s.Ctx)
}

func TestSessionSuite(t *testing.T) {
	suite.Run(t, new(sessionsTestSuite))
}

func (s *sessionsTestSuite) Test_Session_CRUD() {

	// get when no exists
	act, err := s.storage.GetSession(s.Ctx, kit.NewId(), false)
	s.NoError(err)
	s.Empty(act)

	// get when no exists
	act, err = s.storage.GetSession(s.Ctx, kit.NewId(), true)
	s.NoError(err)
	s.Empty(act)

	// create new
	sess := s.session()
	s.NoError(s.storage.MergeSession(s.Ctx, sess))

	// get
	act, err = s.storage.GetSession(s.Ctx, sess.Id, false)
	s.NoError(err)
	s.Equal(act, sess)

	// update
	sess.Details.Kwh = kit.Float64Ptr(30000)
	s.NoError(s.storage.MergeSession(s.Ctx, sess))

	// get
	act, err = s.storage.GetSession(s.Ctx, sess.Id, false)
	s.NoError(err)
	s.Equal(act, sess)

	// update
	sess.Details.Kwh = kit.Float64Ptr(40000)
	s.NoError(s.storage.UpdateSession(s.Ctx, sess))

	// get
	act, err = s.storage.GetSession(s.Ctx, sess.Id, false)
	s.NoError(err)
	s.Equal(act, sess)
}

func (s *sessionsTestSuite) Test_ChargingPeriods_CRUD() {

	// create new
	sess := s.session()
	sess.ChargingPeriods = []*domain.ChargingPeriod{s.chargingPeriod()}
	s.NoError(s.storage.MergeSession(s.Ctx, sess))
	s.NoError(s.storage.CreateChargingPeriods(s.Ctx, sess, sess.ChargingPeriods))

	// get
	act, err := s.storage.GetSession(s.Ctx, sess.Id, true)
	s.NoError(err)
	s.Equal(act, sess)

	// add charging period
	s.NoError(s.storage.CreateChargingPeriods(s.Ctx, sess, []*domain.ChargingPeriod{s.chargingPeriod()}))

	// get
	act, err = s.storage.GetSession(s.Ctx, sess.Id, true)
	s.NoError(err)
	s.Len(act.ChargingPeriods, 2)

	// update
	sess.ChargingPeriods = []*domain.ChargingPeriod{s.chargingPeriod()}
	s.NoError(s.storage.UpdateChargingPeriods(s.Ctx, sess, sess.ChargingPeriods))

	// get
	act, err = s.storage.GetSession(s.Ctx, sess.Id, true)
	s.NoError(err)
	s.Len(act.ChargingPeriods, 1)
	s.Equal(sess.ChargingPeriods, act.ChargingPeriods)

}

func (s *sessionsTestSuite) Test_Search() {
	// create new
	sess := s.session()
	sess.ChargingPeriods = []*domain.ChargingPeriod{s.chargingPeriod()}
	s.NoError(s.storage.MergeSession(s.Ctx, sess))
	s.NoError(s.storage.CreateChargingPeriods(s.Ctx, sess, sess.ChargingPeriods))

	rs, err := s.storage.SearchSessions(s.Ctx, &domain.SessionSearchCriteria{
		PageRequest: domain.PageRequest{},
		ExtId: &domain.PartyExtId{
			PartyId:     sess.ExtId.PartyId,
			CountryCode: sess.ExtId.CountryCode,
		},
		IncPlatforms:        nil,
		ExcPlatforms:        nil,
		WithChargingPeriods: true,
	})
	s.NoError(err)
	s.Len(rs.Items, 1)
	s.NotEmpty(rs.Items[0].Id)
	s.NotEmpty(rs.Items[0].ChargingPeriods)
	s.NotEmpty(rs.Total)
	s.Equal(*rs.Total, 1)
	s.Equal(*rs.Limit, domain.PageSizeDefault)

	rs, err = s.storage.SearchSessions(s.Ctx, &domain.SessionSearchCriteria{
		PageRequest: domain.PageRequest{
			DateFrom: kit.TimePtr(kit.Now().Add(-time.Hour)),
			DateTo:   kit.TimePtr(kit.Now().Add(time.Hour)),
			Limit:    kit.IntPtr(domain.PageSizeMaxLimit + 10),
		},
		IncPlatforms:        []string{sess.PlatformId},
		ExcPlatforms:        []string{kit.NewId()},
		WithChargingPeriods: false,
	})
	s.NoError(err)
	s.Len(rs.Items, 1)
	s.NotEmpty(rs.Items[0].Id)
	s.Empty(rs.Items[0].ChargingPeriods)
	s.NotEmpty(rs.Total)
	s.Equal(*rs.Total, 1)
	s.Equal(*rs.Limit, domain.PageSizeMaxLimit)

	rs, err = s.storage.SearchSessions(s.Ctx, &domain.SessionSearchCriteria{
		PageRequest: domain.PageRequest{},
		Ids:         []string{sess.Id},
	})
	s.NoError(err)
	s.Len(rs.Items, 1)

	rs, err = s.storage.SearchSessions(s.Ctx, &domain.SessionSearchCriteria{
		PageRequest: domain.PageRequest{},
		AuthRef:     sess.Details.AuthRef,
	})
	s.NoError(err)
	s.Len(rs.Items, 1)
}

func (s *sessionsTestSuite) Test_DeleteByExt() {

	sess := s.session()

	// create a session
	s.NoError(s.storage.MergeSession(s.Ctx, sess))

	// get session
	act, err := s.storage.GetSession(s.Ctx, sess.Id, false)
	s.NoError(err)
	s.Equal(act, sess)

	// delete by ext
	s.NoError(s.storage.DeleteSessionsByExtId(s.Ctx, sess.ExtId))

	// get session
	act, err = s.storage.GetSession(s.Ctx, sess.Id, false)
	s.NoError(err)
	s.Empty(act)

}

func (s *sessionsTestSuite) session() *domain.Session {
	partyId := kit.NewRandString()
	return &domain.Session{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     partyId,
				CountryCode: "RS",
			},
			PlatformId:  kit.NewRandString(),
			RefId:       kit.NewRandString(),
			LastUpdated: kit.Now(),
		},
		Id: kit.NewId(),
		Details: domain.SessionDetails{
			StartDateTime: kit.TimePtr(kit.Now()),
			EndDateTime:   nil,
			Kwh:           kit.Float64Ptr(10000),
			CdrToken: &domain.CdrToken{
				PartyExtId: domain.PartyExtId{
					PartyId:     partyId,
					CountryCode: "RS",
				},
				Id:         kit.NewId(),
				Type:       "APP_USER",
				ContractId: kit.NewId(),
			},
			AuthMethod:  domain.AuthMethodCommand,
			AuthRef:     kit.NewRandString(),
			LocationId:  kit.NewId(),
			EvseId:      kit.NewId(),
			ConnectorId: kit.NewId(),
			MeterId:     kit.NewId(),
			Currency:    "RSD",
			TotalCost: &domain.Price{
				ExclVat: 100,
			},
			Status: domain.SessionStatusActive,
		},
	}
}

func (s *sessionsTestSuite) chargingPeriod() *domain.ChargingPeriod {
	return &domain.ChargingPeriod{
		StartDateTime: kit.Now(),
		Dimensions: []*domain.CdrDimension{
			{
				Type:   domain.DimensionTypeEnergy,
				Volume: 1234.34,
			},
		},
		TariffId:              kit.NewId(),
		EncodingMethod:        "method",
		EncodingMethodVersion: kit.IntPtr(1),
		PublicKey:             kit.NewRandString(),
		SignedValues: []*domain.SignedValue{
			{
				Nature:     kit.NewRandString(),
				PlainData:  kit.NewRandString(),
				SignedData: kit.NewRandString(),
			},
		},
		Url: "https://test.com/",
	}
}
