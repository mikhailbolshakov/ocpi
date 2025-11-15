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

type tariffsTestSuite struct {
	kit.Suite
	storage domain.TariffStorage
	adapter Adapter
}

func (s *tariffsTestSuite) SetupSuite() {
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

func (s *tariffsTestSuite) TearDownSuite() {
	_ = s.adapter.Close(s.Ctx)
}

func TestTariffsSuite(t *testing.T) {
	suite.Run(t, new(tariffsTestSuite))
}

func (s *tariffsTestSuite) Test_CRUD() {
	// get when no exists
	act, err := s.storage.GetTariff(s.Ctx, kit.NewId())
	s.NoError(err)
	s.Empty(act)

	// create new
	trf := s.tariff()
	s.NoError(s.storage.MergeTariff(s.Ctx, trf))

	// get
	act, err = s.storage.GetTariff(s.Ctx, trf.Id)
	s.NoError(err)
	s.Equal(act, trf)

	// update
	trf.Details.MinPrice.InclVat = kit.Float64Ptr(30000)
	s.NoError(s.storage.MergeTariff(s.Ctx, trf))

	// get
	act, err = s.storage.GetTariff(s.Ctx, trf.Id)
	s.NoError(err)
	s.Equal(act, trf)

	// update
	trf.Details.MinPrice.InclVat = kit.Float64Ptr(40000)
	s.NoError(s.storage.UpdateTariff(s.Ctx, trf))

	// get
	act, err = s.storage.GetTariff(s.Ctx, trf.Id)
	s.NoError(err)
	s.Equal(act, trf)
}

func (s *tariffsTestSuite) Test_Search() {
	// create new
	trf := s.tariff()
	s.NoError(s.storage.MergeTariff(s.Ctx, trf))

	rs, err := s.storage.SearchTariffs(s.Ctx, &domain.TariffSearchCriteria{
		PageRequest: domain.PageRequest{},
		ExtId: &domain.PartyExtId{
			PartyId:     trf.ExtId.PartyId,
			CountryCode: trf.ExtId.CountryCode,
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

	rs, err = s.storage.SearchTariffs(s.Ctx, &domain.TariffSearchCriteria{
		PageRequest: domain.PageRequest{
			DateFrom: kit.TimePtr(kit.Now().Add(-time.Hour)),
			DateTo:   kit.TimePtr(kit.Now().Add(time.Hour)),
			Limit:    kit.IntPtr(domain.PageSizeMaxLimit + 10),
		},
		IncPlatforms: []string{trf.PlatformId},
		ExcPlatforms: []string{kit.NewId()},
	})
	s.NoError(err)
	s.Len(rs.Items, 1)
	s.NotEmpty(rs.Items[0].Id)
	s.NotEmpty(rs.Total)
	s.Equal(*rs.Total, 1)
	s.Equal(*rs.Limit, domain.PageSizeMaxLimit)

	rs, err = s.storage.SearchTariffs(s.Ctx, &domain.TariffSearchCriteria{
		PageRequest: domain.PageRequest{},
		Ids:         []string{trf.Id},
	})
	s.NoError(err)
	s.Len(rs.Items, 1)
}

func (s *tariffsTestSuite) Test_DeleteByExt() {

	trf := s.tariff()

	// create a tariff
	s.NoError(s.storage.MergeTariff(s.Ctx, trf))

	// get tariff
	act, err := s.storage.GetTariff(s.Ctx, trf.Id)
	s.NoError(err)
	s.Equal(act, trf)

	// delete by ext
	s.NoError(s.storage.DeleteTariffsByExtId(s.Ctx, trf.ExtId))

	// get tariff
	act, err = s.storage.GetTariff(s.Ctx, trf.Id)
	s.NoError(err)
	s.Empty(act)

}

func (s *tariffsTestSuite) tariff() *domain.Tariff {
	partyId := kit.NewRandString()
	return &domain.Tariff{
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
		Details: domain.TariffDetails{
			Currency: "RSD",
			Type:     domain.TariffTypeReg,
			TariffAltText: []*domain.DisplayText{
				{
					Language: "en",
					Text:     "some text",
				},
			},
			TariffAltUrl: "http://test.com/",
			MinPrice: &domain.Price{
				ExclVat: 1000,
				InclVat: kit.Float64Ptr(800),
			},
			MaxPrice: &domain.Price{
				ExclVat: 1000,
				InclVat: kit.Float64Ptr(800),
			},
			Elements: []*domain.TariffElement{
				{
					PriceComponents: []*domain.PriceComponent{
						{
							Type:     domain.TariffDimEnergy,
							Price:    100,
							Vat:      kit.Float64Ptr(10),
							StepSize: 1,
						},
					},
					Restrictions: &domain.TariffRestrictions{
						StartTime:   "10:00",
						EndTime:     "20:00",
						StartDate:   kit.NowPtr(),
						EndDate:     kit.NowPtr(),
						MinKwh:      kit.Float64Ptr(10),
						MaxKwh:      kit.Float64Ptr(100),
						MinCurrent:  kit.Float64Ptr(10),
						MaxCurrent:  kit.Float64Ptr(100),
						MinPower:    kit.Float64Ptr(10),
						MaxPower:    kit.Float64Ptr(100),
						MinDuration: kit.Float64Ptr(10),
						MaxDuration: kit.Float64Ptr(100),
						DayOfWeek:   []string{domain.DayMon, domain.DayFri},
						Reservation: domain.Reservation,
					},
				},
			},
			StartDateTime: kit.TimePtr(kit.Now()),
			EndDateTime:   kit.TimePtr(kit.Now().Add(time.Hour * 24)),
			EnergyMix: &domain.EnergyMix{
				SupplierName:      "aaa",
				EnergyProductName: "aaa",
			},
		},
	}
}
