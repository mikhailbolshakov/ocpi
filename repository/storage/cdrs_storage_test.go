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

type cdrsTestSuite struct {
	kit.Suite
	storage domain.CdrStorage
	adapter Adapter
}

func (s *cdrsTestSuite) SetupSuite() {
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

func (s *cdrsTestSuite) TearDownSuite() {
	_ = s.adapter.Close(s.Ctx)
}

func TestCdrSuite(t *testing.T) {
	suite.Run(t, new(cdrsTestSuite))
}

func (s *cdrsTestSuite) Test_CRUD() {
	// get when no exists
	act, err := s.storage.GetCdr(s.Ctx, kit.NewId())
	s.NoError(err)
	s.Empty(act)

	// create new
	cdr := s.cdr()
	s.NoError(s.storage.MergeCdr(s.Ctx, cdr))

	// get
	act, err = s.storage.GetCdr(s.Ctx, cdr.Id)
	s.NoError(err)
	s.Equal(act, cdr)

	// update
	cdr.Details.TotalEnergy = 30000
	s.NoError(s.storage.MergeCdr(s.Ctx, cdr))

	// get
	act, err = s.storage.GetCdr(s.Ctx, cdr.Id)
	s.NoError(err)
	s.Equal(act, cdr)

	// update
	cdr.Details.TotalEnergy = 40000
	s.NoError(s.storage.UpdateCdr(s.Ctx, cdr))

	// get
	act, err = s.storage.GetCdr(s.Ctx, cdr.Id)
	s.NoError(err)
	s.Equal(act, cdr)
}

func (s *cdrsTestSuite) Test_Search() {
	// create new
	cdr := s.cdr()
	s.NoError(s.storage.MergeCdr(s.Ctx, cdr))

	rs, err := s.storage.SearchCdrs(s.Ctx, &domain.CdrSearchCriteria{
		PageRequest: domain.PageRequest{},
		ExtId: &domain.PartyExtId{
			PartyId:     cdr.ExtId.PartyId,
			CountryCode: cdr.ExtId.CountryCode,
		},
		IncPlatforms: nil,
		ExcPlatforms: nil,
	})
	s.NoError(err)
	s.Len(rs.Items, 1)
	s.NotEmpty(rs.Items[0].Id)
	s.NotEmpty(rs.Total)
	s.Equal(*rs.Total, 1)
	s.Equal(*rs.Limit, domain.PageSizeDefault)

	rs, err = s.storage.SearchCdrs(s.Ctx, &domain.CdrSearchCriteria{
		PageRequest: domain.PageRequest{
			DateFrom: kit.TimePtr(kit.Now().Add(-time.Hour)),
			DateTo:   kit.TimePtr(kit.Now().Add(time.Hour)),
			Limit:    kit.IntPtr(domain.PageSizeMaxLimit + 10),
		},
		IncPlatforms: []string{cdr.PlatformId},
		ExcPlatforms: []string{kit.NewId()},
	})
	s.NoError(err)
	s.Len(rs.Items, 1)
	s.NotEmpty(rs.Items[0].Id)
	s.NotEmpty(rs.Total)
	s.Equal(*rs.Total, 1)
	s.Equal(*rs.Limit, domain.PageSizeMaxLimit)

	rs, err = s.storage.SearchCdrs(s.Ctx, &domain.CdrSearchCriteria{
		PageRequest: domain.PageRequest{},
		Ids:         []string{cdr.Id},
	})
	s.NoError(err)
	s.Len(rs.Items, 1)
}

func (s *cdrsTestSuite) Test_DeleteByExt() {

	sess := s.cdr()

	// create a cdr
	s.NoError(s.storage.MergeCdr(s.Ctx, sess))

	// get cdr
	act, err := s.storage.GetCdr(s.Ctx, sess.Id)
	s.NoError(err)
	s.Equal(act, sess)

	// delete by ext
	s.NoError(s.storage.DeleteCdrsByExtId(s.Ctx, sess.ExtId))

	// get cdr
	act, err = s.storage.GetCdr(s.Ctx, sess.Id)
	s.NoError(err)
	s.Empty(act)

}

func (s *cdrsTestSuite) cdr() *domain.Cdr {
	partyId := kit.NewRandString()
	return &domain.Cdr{
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
		Details: domain.CdrDetails{
			StartDateTime: kit.Now(),
			EndDateTime:   kit.Now().Add(time.Hour),
			SessionId:     kit.NewId(),
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
			CdrLocation: domain.CdrLocation{},
			MeterId:     kit.NewId(),
			Currency:    "RSD",
			Tariffs:     nil,
			ChargingPeriods: []*domain.ChargingPeriod{
				{
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
				},
			},
			SignedData:               nil,
			TotalCost:                domain.Price{ExclVat: 100},
			TotalFixedCost:           &domain.Price{ExclVat: 100},
			TotalEnergy:              1000,
			TotalEnergyCost:          &domain.Price{ExclVat: 100},
			TotalTime:                10,
			TotalTimeCost:            &domain.Price{ExclVat: 100},
			TotalParkingTime:         kit.Float64Ptr(10.0),
			TotalParkingCost:         &domain.Price{ExclVat: 100},
			TotalReservationCost:     &domain.Price{ExclVat: 100},
			Remark:                   "remark",
			InvoiceReferenceId:       kit.NewId(),
			Credit:                   true,
			CreditReferenceId:        kit.NewId(),
			HomeChargingCompensation: false,
		},
	}
}
