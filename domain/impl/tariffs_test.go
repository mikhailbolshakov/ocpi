package impl

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type tariffTestSuite struct {
	kit.Suite
	svc     *tariffService
	storage *mocks.TariffStorage
}

func (s *tariffTestSuite) SetupSuite() {
	s.Suite.Init(ocpi.LF())
}

func (s *tariffTestSuite) SetupTest() {
	s.storage = &mocks.TariffStorage{}
	s.svc = NewTariffService(s.storage).(*tariffService)
}

func (s *tariffTestSuite) TearDownSuite() {}

func TestTariffSuite(t *testing.T) {
	suite.Run(t, new(tariffTestSuite))
}

func (s *tariffTestSuite) Test_PutTariff_WhenLaterExists_Skip() {
	trf := s.tariff()
	stored := s.tariff()
	stored.Id = trf.Id
	stored.OcpiItem = trf.OcpiItem
	trf.LastUpdated = kit.Now().Add(-time.Hour)
	s.storage.On("GetTariff", s.Ctx, trf.Id).Return(stored, nil)
	r, err := s.svc.PutTariff(s.Ctx, trf)
	s.NoError(err)
	s.Nil(r)
}

func (s *tariffTestSuite) Test_PutTariff_Ok() {
	trf := s.tariff()
	s.storage.On("GetTariff", s.Ctx, trf.Id).Return(nil, nil)
	s.storage.On("MergeTariff", s.Ctx, trf).Return(nil)
	_, err := s.svc.PutTariff(s.Ctx, trf)
	s.NoError(err)
	s.AssertCalled(&s.storage.Mock, "MergeTariff", s.Ctx, trf)
}

func (s *tariffTestSuite) Test_MergeTariff_WhenLastUpdatedLater_Skip() {
	stored := s.tariff()
	s.storage.On("GetTariff", s.Ctx, stored.Id).Return(stored, nil)
	trf := &domain.Tariff{Id: stored.Id}
	trf.LastUpdated = kit.Now().Add(-2 * time.Hour)
	trf.Details.Currency = "USD"
	r, err := s.svc.MergeTariff(s.Ctx, trf)
	s.NoError(err)
	s.Nil(r)
}

func (s *tariffTestSuite) Test_MergeTariff_WhenEmptyLastUpdated_Fail() {
	stored := s.tariff()
	s.storage.On("GetTariff", s.Ctx, stored.Id).Return(stored, nil)
	trf := &domain.Tariff{Id: stored.Id}
	_, err := s.svc.MergeTariff(s.Ctx, trf)
	s.Error(err)
}

func (s *tariffTestSuite) Test_MergeTariff_WhenChanged() {
	stored := s.tariff()
	s.storage.On("GetTariff", s.Ctx, stored.Id).Return(stored, nil)
	s.storage.On("UpdateTariff", s.Ctx, stored).Return(nil)
	trf := &domain.Tariff{Id: stored.Id}
	stored.LastUpdated = kit.Now().Add(-time.Hour)
	trf.LastUpdated = kit.Now()
	trf.Details.Currency = "USD"
	r, err := s.svc.MergeTariff(s.Ctx, trf)
	s.NoError(err)
	s.NotEmpty(r)
	s.Equal(stored.ExtId, r.ExtId)
	s.Equal(trf.Details.Currency, r.Details.Currency)
}

func (s *tariffTestSuite) Test_ValidatePrice() {
	// valid
	trf := s.tariff()
	s.NoError(s.svc.validatePrice(s.Ctx, "", "", trf.Details.MinPrice))

	// excl vat less zero
	trf = s.tariff()
	trf.Details.MinPrice.ExclVat = -1
	s.Error(s.svc.validatePrice(s.Ctx, "", "", trf.Details.MinPrice))

	// incl vat less zero
	trf = s.tariff()
	trf.Details.MinPrice.InclVat = kit.Float64Ptr(-10)
	s.Error(s.svc.validatePrice(s.Ctx, "", "", trf.Details.MinPrice))

}

func (s *tariffTestSuite) Test_ValidatePriceComponent() {
	// valid
	pc := s.tariff().Details.Elements[0].PriceComponents[0]
	s.NoError(s.svc.validatePriceComponent(s.Ctx, pc))

	// empty type
	pc = s.tariff().Details.Elements[0].PriceComponents[0]
	pc.Type = ""
	s.Error(s.svc.validatePriceComponent(s.Ctx, pc))

	// invalid type
	pc = s.tariff().Details.Elements[0].PriceComponents[0]
	pc.Type = "invalid"
	s.Error(s.svc.validatePriceComponent(s.Ctx, pc))

	// invalid step size
	pc = s.tariff().Details.Elements[0].PriceComponents[0]
	pc.StepSize = -1
	s.Error(s.svc.validatePriceComponent(s.Ctx, pc))

	// invalid vat
	pc = s.tariff().Details.Elements[0].PriceComponents[0]
	pc.Vat = kit.Float64Ptr(-1)
	s.Error(s.svc.validatePriceComponent(s.Ctx, pc))

	// invalid price
	pc = s.tariff().Details.Elements[0].PriceComponents[0]
	pc.Price = -1
	s.Error(s.svc.validatePriceComponent(s.Ctx, pc))
}

func (s *tariffTestSuite) Test_ValidateRestrictions() {
	// valid
	r := s.tariff().Details.Elements[0].Restrictions
	s.NoError(s.svc.validateRestriction(s.Ctx, r))

	// invalid start time
	r = s.tariff().Details.Elements[0].Restrictions
	r.StartTime = "invalid"
	s.Error(s.svc.validateRestriction(s.Ctx, r))

	// invalid start time
	r = s.tariff().Details.Elements[0].Restrictions
	r.StartTime = "55:55"
	s.Error(s.svc.validateRestriction(s.Ctx, r))

	// invalid end time
	r = s.tariff().Details.Elements[0].Restrictions
	r.EndTime = "invalid"
	s.Error(s.svc.validateRestriction(s.Ctx, r))

	// invalid end time
	r = s.tariff().Details.Elements[0].Restrictions
	r.EndTime = "55:55"
	s.Error(s.svc.validateRestriction(s.Ctx, r))

	// invalid min kwh
	r = s.tariff().Details.Elements[0].Restrictions
	r.MinKwh = kit.Float64Ptr(-1)
	s.Error(s.svc.validateRestriction(s.Ctx, r))

	// invalid max kwh
	r = s.tariff().Details.Elements[0].Restrictions
	r.MaxKwh = kit.Float64Ptr(-1)
	s.Error(s.svc.validateRestriction(s.Ctx, r))

	// invalid day of week
	r = s.tariff().Details.Elements[0].Restrictions
	r.DayOfWeek = []string{"invalid"}
	s.Error(s.svc.validateRestriction(s.Ctx, r))

	// invalid reservation
	r = s.tariff().Details.Elements[0].Restrictions
	r.Reservation = "invalid"
	s.Error(s.svc.validateRestriction(s.Ctx, r))
}

func (s *tariffTestSuite) Test_Validate() {
	// valid
	trf := s.tariff()
	s.NoError(s.svc.Validate(s.Ctx, trf))

	// invalid currency
	trf.Details.Currency = "invalid"
	s.Error(s.svc.Validate(s.Ctx, trf))

	// invalid type
	trf.Details.Type = "invalid"
	s.Error(s.svc.Validate(s.Ctx, trf))

	// invalid start/end date
	trf.Details.StartDateTime = kit.TimePtr(kit.Now().Add(time.Hour))
	trf.Details.EndDateTime = kit.TimePtr(kit.Now().Add(-time.Hour))
	s.Error(s.svc.Validate(s.Ctx, trf))
}

func (s *tariffTestSuite) tariff() *domain.Tariff {
	return &domain.Tariff{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     "PPP",
				CountryCode: "RS",
			},
			PlatformId:  kit.NewRandString(),
			RefId:       kit.NewRandString(),
			LastUpdated: kit.Now(),
			LastSent:    nil,
		},
		Id: kit.NewRandString(),
		Details: domain.TariffDetails{
			Currency: "RSD",
			Type:     domain.TariffTypeReg,
			TariffAltText: []*domain.DisplayText{
				{
					Language: "en",
					Text:     "text",
				},
			},
			TariffAltUrl: "https://example.com/url",
			MinPrice: &domain.Price{
				ExclVat: 10.0,
				InclVat: nil,
			},
			MaxPrice: &domain.Price{
				ExclVat: 10000.0,
				InclVat: nil,
			},
			Elements:      []*domain.TariffElement{s.element()},
			StartDateTime: kit.TimePtr(kit.Now().Add(-time.Hour * 24 * 10)),
			EndDateTime:   kit.TimePtr(kit.Now().Add(time.Hour * 24 * 10)),
			EnergyMix:     nil,
		},
	}
}

func (s *tariffTestSuite) element() *domain.TariffElement {
	return &domain.TariffElement{
		PriceComponents: []*domain.PriceComponent{s.component()},
		Restrictions:    s.restriction(),
	}
}

func (s *tariffTestSuite) component() *domain.PriceComponent {
	return &domain.PriceComponent{
		Type:     domain.DimensionTypeEnergy,
		Price:    20,
		Vat:      kit.Float64Ptr(2.0),
		StepSize: 1,
	}
}

func (s *tariffTestSuite) restriction() *domain.TariffRestrictions {
	return &domain.TariffRestrictions{
		StartTime:   "10:00",
		EndTime:     "20:00",
		StartDate:   kit.TimePtr(kit.Now().Add(-time.Hour)),
		EndDate:     kit.TimePtr(kit.Now().Add(time.Hour)),
		MinKwh:      kit.Float64Ptr(1),
		MaxKwh:      kit.Float64Ptr(100),
		MinCurrent:  kit.Float64Ptr(1),
		MaxCurrent:  kit.Float64Ptr(100),
		MinPower:    kit.Float64Ptr(1),
		MaxPower:    kit.Float64Ptr(100),
		MinDuration: kit.Float64Ptr(1),
		MaxDuration: kit.Float64Ptr(100),
		DayOfWeek:   []string{domain.DayMon, domain.DayFri},
		Reservation: "",
	}
}
