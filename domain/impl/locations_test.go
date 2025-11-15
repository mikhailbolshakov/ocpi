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

type locationTestSuite struct {
	kit.Suite
	svc     *locationService
	storage *mocks.LocationStorage
}

func (s *locationTestSuite) SetupSuite() {
	s.Suite.Init(ocpi.LF())
}

func (s *locationTestSuite) SetupTest() {
	s.storage = &mocks.LocationStorage{}
	s.svc = NewLocationService(s.storage).(*locationService)
}

func (s *locationTestSuite) TearDownSuite() {}

func TestLocationSuite(t *testing.T) {
	suite.Run(t, new(locationTestSuite))
}

func (s *locationTestSuite) Test_PutLocation_WhenLaterExists_Skip() {
	loc := s.location()
	stored := s.location()
	stored.Id = loc.Id
	stored.OcpiItem = loc.OcpiItem
	loc.LastUpdated = kit.Now().Add(-time.Hour)
	s.storage.On("GetLocation", s.Ctx, loc.Id, true).Return(stored, nil)
	r, err := s.svc.PutLocation(s.Ctx, loc)
	s.NoError(err)
	s.Nil(r)
}

func (s *locationTestSuite) Test_PutLocation_Ok() {
	loc := s.location()
	s.storage.On("GetLocation", s.Ctx, loc.Id, true).Return(nil, nil)
	s.storage.On("MergeLocation", s.Ctx, loc).Return(nil)
	_, err := s.svc.PutLocation(s.Ctx, loc)
	s.NoError(err)
	s.AssertCalled(&s.storage.Mock, "MergeLocation", s.Ctx, loc)
}

func (s *locationTestSuite) Test_MergeLocation_WhenLastUpdatedLater_Skip() {
	stored := s.location()
	s.storage.On("GetLocation", s.Ctx, stored.Id, false).Return(stored, nil)
	loc := &domain.Location{Id: stored.Id}
	loc.LastUpdated = kit.Now().Add(-2 * time.Hour)
	loc.Details.Coordinates.Latitude = "67.234567"
	r, err := s.svc.MergeLocation(s.Ctx, loc)
	s.NoError(err)
	s.Nil(r)
}

func (s *locationTestSuite) Test_MergeLocation_WhenEmptyLastUpdated_Fail() {
	stored := s.location()
	s.storage.On("GetLocation", s.Ctx, stored.Id, false).Return(stored, nil)
	loc := &domain.Location{Id: stored.Id}
	_, err := s.svc.MergeLocation(s.Ctx, loc)
	s.Error(err)
}

func (s *locationTestSuite) Test_MergeLocation_WhenChanged() {
	stored := s.location()
	s.storage.On("GetLocation", s.Ctx, stored.Id, false).Return(stored, nil)
	s.storage.On("MergeLocation", s.Ctx, stored).Return(nil)
	loc := &domain.Location{Id: stored.Id}
	stored.LastUpdated = kit.Now().Add(-time.Hour)
	loc.LastUpdated = kit.Now()
	loc.Details.Coordinates.Latitude = "67.234567"
	r, err := s.svc.MergeLocation(s.Ctx, loc)
	s.NoError(err)
	s.NotEmpty(r)
	s.Equal(stored.ExtId, r.ExtId)
	s.Equal(loc.Details.Coordinates.Latitude, r.Details.Coordinates.Latitude)
}

func (s *locationTestSuite) Test_ValidateAddress() {
	// valid
	loc := s.location()
	s.NoError(s.svc.validateAddress(s.Ctx, &loc.Details))

	// empty address
	loc = s.location()
	loc.Details.Address = ""
	s.Error(s.svc.validateAddress(s.Ctx, &loc.Details))

	// empty city
	loc = s.location()
	loc.Details.City = ""
	s.Error(s.svc.validateAddress(s.Ctx, &loc.Details))

	// empty country
	loc = s.location()
	loc.Details.Country = ""
	s.Error(s.svc.validateAddress(s.Ctx, &loc.Details))

	// invalid country
	loc = s.location()
	loc.Details.Country = "XXX"
	s.Error(s.svc.validateAddress(s.Ctx, &loc.Details))
}

func (s *locationTestSuite) Test_ValidateCoordinates() {
	// valid
	loc := s.location()
	s.NoError(s.svc.validateCoordinates(s.Ctx, &loc.Details))

	// lat empty
	loc = s.location()
	loc.Details.Coordinates.Latitude = ""
	s.Error(s.svc.validateCoordinates(s.Ctx, &loc.Details))

	// lat invalid
	loc = s.location()
	loc.Details.Coordinates.Latitude = "123.678"
	s.Error(s.svc.validateCoordinates(s.Ctx, &loc.Details))

	// rel lat invalid
	loc = s.location()
	loc.Details.RelatedLocations[0].Latitude = "123.678"
	s.Error(s.svc.validateCoordinates(s.Ctx, &loc.Details))
}

func (s *locationTestSuite) Test_ValidateHours() {
	// valid
	loc := s.location()
	s.NoError(s.svc.validateHours(s.Ctx, "", loc.Details.OpeningTimes))

	// invalid regular
	loc = s.location()
	loc.Details.OpeningTimes.RegularHours = nil
	s.Error(s.svc.validateHours(s.Ctx, "", loc.Details.OpeningTimes))

	// invalid weekday
	loc = s.location()
	loc.Details.OpeningTimes.RegularHours[0].Weekday = 10
	s.Error(s.svc.validateHours(s.Ctx, "", loc.Details.OpeningTimes))

	// invalid begin time
	loc = s.location()
	loc.Details.OpeningTimes.RegularHours[0].PeriodBegin = "25:25"
	s.Error(s.svc.validateHours(s.Ctx, "", loc.Details.OpeningTimes))

	// invalid end time
	loc = s.location()
	loc.Details.OpeningTimes.RegularHours[0].PeriodEnd = "AA"
	s.Error(s.svc.validateHours(s.Ctx, "", loc.Details.OpeningTimes))

	// invalid openings
	loc = s.location()
	loc.Details.OpeningTimes.ExceptionalOpenings[0] = &domain.ExceptionalPeriod{
		PeriodBegin: time.Time{},
		PeriodEnd:   time.Time{},
	}
	s.Error(s.svc.validateHours(s.Ctx, "", loc.Details.OpeningTimes))

	loc = s.location()
	loc.Details.OpeningTimes.ExceptionalOpenings[0] = &domain.ExceptionalPeriod{
		PeriodBegin: kit.Now().Add(time.Hour),
		PeriodEnd:   kit.Now(),
	}
	s.Error(s.svc.validateHours(s.Ctx, "", loc.Details.OpeningTimes))
}

func (s *locationTestSuite) Test_ValidateTZ() {
	// valid
	loc := s.location()
	s.NoError(s.svc.validateTZ(s.Ctx, &loc.Details))

	// empty
	loc = s.location()
	loc.Details.TimeZone = ""
	s.Error(s.svc.validateTZ(s.Ctx, &loc.Details))

	// empty
	loc = s.location()
	loc.Details.TimeZone = "242142"
	s.Error(s.svc.validateTZ(s.Ctx, &loc.Details))
}

func (s *locationTestSuite) Test_ValidateLocation() {
	// valid
	loc := s.location()
	s.NoError(s.svc.validateLocation(s.Ctx, loc))

	// invalid parking type
	loc = s.location()
	loc.Details.ParkingType = "invalid"
	s.Error(s.svc.validateLocation(s.Ctx, loc))

	// invalid facilities
	loc = s.location()
	loc.Details.Facilities = []string{"invalid"}
	s.Error(s.svc.validateLocation(s.Ctx, loc))
}

func (s *locationTestSuite) Test_ValidateEvse() {
	// valid
	evse := s.evse()
	s.NoError(s.svc.validateEvse(s.Ctx, evse))

	// empty id
	evse = s.evse()
	evse.Id = ""
	s.Error(s.svc.validateEvse(s.Ctx, evse))

	// empty location
	evse = s.evse()
	evse.LocationId = ""
	s.Error(s.svc.validateEvse(s.Ctx, evse))

	// empty status
	evse = s.evse()
	evse.Status = ""
	s.Error(s.svc.validateEvse(s.Ctx, evse))

	// invalid status
	evse = s.evse()
	evse.Status = "invalid"
	s.Error(s.svc.validateEvse(s.Ctx, evse))

	// invalid status schedule
	evse = s.evse()
	evse.Details.StatusSchedule[0].Status = "invalid"
	s.Error(s.svc.validateEvse(s.Ctx, evse))
	evse = s.evse()
	evse.Details.StatusSchedule[0].PeriodEnd = kit.TimePtr(kit.Now().Add(-time.Hour))
	s.Error(s.svc.validateEvse(s.Ctx, evse))

	// invalid capabilities
	evse = s.evse()
	evse.Details.Capabilities = []string{"invalid"}
	s.Error(s.svc.validateEvse(s.Ctx, evse))

	// invalid coordinates
	evse = s.evse()
	evse.Details.Coordinates.Latitude = "invalid"
	s.Error(s.svc.validateEvse(s.Ctx, evse))

	// invalid parking restriction
	evse = s.evse()
	evse.Details.ParkingRestrictions = []string{"invalid"}
	s.Error(s.svc.validateEvse(s.Ctx, evse))

}

func (s *locationTestSuite) Test_PutEvse_WhenLaterExists_Skip() {
	evse := s.evse()
	stored := s.evse()
	stored.Id = evse.Id
	stored.OcpiItem = evse.OcpiItem
	evse.LastUpdated = kit.Now().Add(-time.Hour)
	s.storage.On("GetEvse", s.Ctx, evse.LocationId, evse.Id, true).Return(stored, nil)
	r, err := s.svc.PutEvse(s.Ctx, evse)
	s.NoError(err)
	s.Nil(r)
}

func (s *locationTestSuite) Test_PutEvse_Ok() {
	evse := s.evse()
	s.storage.On("GetEvse", s.Ctx, evse.LocationId, evse.Id, true).Return(nil, nil)
	s.storage.On("MergeEvse", s.Ctx, evse).Return(nil)
	_, err := s.svc.PutEvse(s.Ctx, evse)
	s.NoError(err)
	s.AssertCalled(&s.storage.Mock, "MergeEvse", s.Ctx, evse)
}

func (s *locationTestSuite) Test_MergeEvse_WhenLastUpdatedLater_Skip() {
	stored := s.evse()
	s.storage.On("GetEvse", s.Ctx, stored.LocationId, stored.Id, false).Return(stored, nil)
	evse := &domain.Evse{Id: stored.Id, LocationId: stored.LocationId}
	evse.LastUpdated = kit.Now().Add(-2 * time.Hour)
	evse.Details.Coordinates = &domain.GeoLocation{
		Latitude:  "67.234567",
		Longitude: "67.234567",
	}
	r, err := s.svc.MergeEvse(s.Ctx, evse)
	s.NoError(err)
	s.Nil(r)
}

func (s *locationTestSuite) Test_MergeEvse_WhenEmptyLastUpdated_Fail() {
	stored := s.evse()
	s.storage.On("GetEvse", s.Ctx, stored.LocationId, stored.Id).Return(stored, nil)
	evse := &domain.Evse{Id: stored.Id, LocationId: stored.LocationId}
	_, err := s.svc.MergeEvse(s.Ctx, evse)
	s.Error(err)
}

func (s *locationTestSuite) Test_MergeEvse_WhenChanged() {
	stored := s.evse()
	s.storage.On("GetEvse", s.Ctx, stored.LocationId, stored.Id, false).Return(stored, nil)
	s.storage.On("MergeEvse", s.Ctx, stored).Return(nil)
	evse := &domain.Evse{Id: stored.Id, LocationId: stored.LocationId}
	stored.LastUpdated = kit.Now().Add(-time.Hour)
	evse.LastUpdated = kit.Now()
	evse.Details.Coordinates = &domain.GeoLocation{
		Latitude:  "67.234567",
		Longitude: "67.234567",
	}
	r, err := s.svc.MergeEvse(s.Ctx, evse)
	s.NoError(err)
	s.NotEmpty(r)
	s.Equal(stored.ExtId, r.ExtId)
	s.Equal(evse.Details.Coordinates.Latitude, r.Details.Coordinates.Latitude)
}

func (s *locationTestSuite) Test_ValidateConnector() {
	// valid
	con := s.connector()
	s.NoError(s.svc.validateConnector(s.Ctx, con))

	// empty id
	con = s.connector()
	con.Id = ""
	s.Error(s.svc.validateConnector(s.Ctx, con))

	// empty loc id
	con = s.connector()
	con.LocationId = ""
	s.Error(s.svc.validateConnector(s.Ctx, con))

	// empty evse id
	con = s.connector()
	con.EvseId = ""
	s.Error(s.svc.validateConnector(s.Ctx, con))

	// empty standard
	con = s.connector()
	con.Details.Standard = ""
	s.Error(s.svc.validateConnector(s.Ctx, con))

	// invalid standard
	con = s.connector()
	con.Details.Standard = "invalid"
	s.Error(s.svc.validateConnector(s.Ctx, con))

	// empty format
	con = s.connector()
	con.Details.Format = ""
	s.Error(s.svc.validateConnector(s.Ctx, con))

	// invalid format
	con = s.connector()
	con.Details.Format = "invalid"
	s.Error(s.svc.validateConnector(s.Ctx, con))

	// empty power type
	con = s.connector()
	con.Details.PowerType = ""
	s.Error(s.svc.validateConnector(s.Ctx, con))

	// invalid power type
	con = s.connector()
	con.Details.PowerType = "invalid"
	s.Error(s.svc.validateConnector(s.Ctx, con))

	// invalid voltage
	con = s.connector()
	con.Details.MaxVoltage = -10
	s.Error(s.svc.validateConnector(s.Ctx, con))
	con.Details.MaxVoltage = 1000000000
	s.Error(s.svc.validateConnector(s.Ctx, con))

	// invalid amperage
	con = s.connector()
	con.Details.MaxAmperage = -10
	s.Error(s.svc.validateConnector(s.Ctx, con))
	con.Details.MaxAmperage = 1000000000
	s.Error(s.svc.validateConnector(s.Ctx, con))

	// invalid power
	con = s.connector()
	con.Details.MaxElectricPower = kit.Float64Ptr(-10)
	s.Error(s.svc.validateConnector(s.Ctx, con))
	con.Details.MaxElectricPower = kit.Float64Ptr(1000000000)
	s.Error(s.svc.validateConnector(s.Ctx, con))

	// empty tariffs
	//con = s.connector()
	//con.Details.TariffIds = nil
	//s.Error(s.svc.validateConnector(s.Ctx, con))

	// invalid terms
	con = s.connector()
	con.Details.TermsAndConditions = "invalid"
	s.Error(s.svc.validateConnector(s.Ctx, con))

}

func (s *locationTestSuite) Test_PutConnector_WhenLaterExists_Skip() {
	connector := s.connector()
	stored := s.connector()
	stored.Id = connector.Id
	stored.OcpiItem = connector.OcpiItem
	connector.LastUpdated = kit.Now().Add(-time.Hour)
	s.storage.On("GetConnector", s.Ctx, connector.LocationId, connector.EvseId, connector.Id).Return(stored, nil)
	r, err := s.svc.PutConnector(s.Ctx, connector)
	s.NoError(err)
	s.Nil(r)
}

func (s *locationTestSuite) Test_PutConnector_Ok() {
	connector := s.connector()
	s.storage.On("GetConnector", s.Ctx, connector.LocationId, connector.EvseId, connector.Id).Return(nil, nil)
	s.storage.On("MergeConnector", s.Ctx, connector).Return(nil)
	_, err := s.svc.PutConnector(s.Ctx, connector)
	s.NoError(err)
	s.AssertCalled(&s.storage.Mock, "MergeConnector", s.Ctx, connector)
}

func (s *locationTestSuite) Test_PutLocationEvseConnector_Ok() {
	loc := s.location()
	evse := s.evse()
	evse.Connectors = []*domain.Connector{s.connector()}
	loc.Evses = []*domain.Evse{evse}
	s.storage.On("GetLocation", s.Ctx, loc.Id, true).Return(nil, nil)
	s.storage.On("MergeLocation", s.Ctx, loc).Return(nil)
	act, err := s.svc.PutLocation(s.Ctx, loc)
	s.NoError(err)
	s.Equal(act.ExtId, act.Evses[0].ExtId)
	s.Equal(act.LastUpdated, act.Evses[0].LastUpdated)
	s.Equal(act.Id, act.Evses[0].LocationId)
	s.Equal(act.Evses[0].ExtId, act.Evses[0].Connectors[0].ExtId)
	s.Equal(act.Evses[0].LastUpdated, act.Evses[0].Connectors[0].LastUpdated)
	s.Equal(act.Evses[0].Id, act.Evses[0].Connectors[0].EvseId)
	s.Equal(act.Evses[0].LocationId, act.Evses[0].Connectors[0].LocationId)
}

func (s *locationTestSuite) Test_MergeConnector_WhenLastUpdatedLater_Skip() {
	stored := s.connector()
	s.storage.On("GetConnector", s.Ctx, stored.LocationId, stored.EvseId, stored.Id).Return(stored, nil)
	connector := &domain.Connector{Id: stored.Id, LocationId: stored.LocationId, EvseId: stored.EvseId}
	connector.LastUpdated = kit.Now().Add(-time.Hour)
	connector.Details.Standard = domain.ConnectorTypeDomesticA
	r, err := s.svc.MergeConnector(s.Ctx, connector)
	s.NoError(err)
	s.Nil(r)
}

func (s *locationTestSuite) Test_MergeConnector_WhenEmptyLastUpdated_Fail() {
	stored := s.connector()
	s.storage.On("GetConnector", s.Ctx, stored.LocationId, stored.EvseId, stored.Id).Return(stored, nil)
	connector := &domain.Connector{Id: stored.Id, LocationId: stored.LocationId, EvseId: stored.EvseId}
	_, err := s.svc.MergeConnector(s.Ctx, connector)
	s.Error(err)
}

func (s *locationTestSuite) Test_MergeConnector_WhenChanged() {
	stored := s.connector()
	s.storage.On("GetConnector", s.Ctx, stored.LocationId, stored.EvseId, stored.Id).Return(stored, nil)
	s.storage.On("MergeConnector", s.Ctx, stored).Return(nil)
	connector := &domain.Connector{Id: stored.Id, LocationId: stored.LocationId, EvseId: stored.EvseId}
	stored.LastUpdated = kit.Now().Add(-time.Hour)
	connector.LastUpdated = kit.Now()
	connector.Details.Standard = domain.ConnectorTypeDomesticA
	r, err := s.svc.MergeConnector(s.Ctx, connector)
	s.NoError(err)
	s.NotEmpty(r)
	s.Equal(stored.ExtId, r.ExtId)
	s.Equal(connector.Details.Standard, r.Details.Standard)
}

func (s *locationTestSuite) location() *domain.Location {
	return &domain.Location{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     "EVG",
				CountryCode: "RS",
			},
			PlatformId:  "test",
			LastUpdated: kit.Now(),
			LastSent:    nil,
		},
		Id: kit.NewId(),
		Details: domain.LocationDetails{
			Publish:          nil,
			PublishAllowedTo: nil,
			Name:             "name",
			Address:          "address",
			City:             "city",
			PostalCode:       "11111",
			State:            "state",
			Country:          "SRB",
			Coordinates: domain.GeoLocation{
				Latitude:  "57.5464646",
				Longitude: "37.3543636",
			},
			RelatedLocations: []*domain.AdditionalGeoLocation{
				{
					GeoLocation: domain.GeoLocation{
						Latitude:  "57.5464646",
						Longitude: "37.3543636",
					},
					Name: &domain.DisplayText{
						Language: "ru",
						Text:     "text",
					},
				},
			},
			ParkingType: domain.ParkingTypeDriveway,
			Directions:  nil,
			Operator: &domain.BusinessDetails{
				Name:    "name",
				Website: "https://test.com/page",
				Logo:    nil,
			},
			Facilities: []string{domain.FacilityHotel, domain.FacilityAirport},
			TimeZone:   "Europe/Belgrade",
			OpeningTimes: &domain.Hours{
				TwentyFourSeven: false,
				RegularHours: []*domain.RegularHours{
					{
						Weekday:     1,
						PeriodBegin: "10:00",
						PeriodEnd:   "20:00",
					},
				},
				ExceptionalOpenings: []*domain.ExceptionalPeriod{
					{
						PeriodBegin: kit.Now().Add(-time.Hour),
						PeriodEnd:   kit.Now(),
					},
				},
				ExceptionalClosings: nil,
			},
			ChargingWhenClosed: nil,
			Images:             nil,
			EnergyMix:          nil,
		},
		Evses: nil,
	}
}

func (s *locationTestSuite) evse() *domain.Evse {
	return &domain.Evse{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     "EVG",
				CountryCode: "RS",
			},
			PlatformId:  "test",
			LastUpdated: kit.Now(),
			LastSent:    nil,
		},
		Id:         kit.NewId(),
		LocationId: kit.NewId(),
		Status:     domain.EvseStatusAvailable,
		Details: domain.EvseDetails{
			EvseId: kit.NewId(),
			StatusSchedule: []*domain.StatusSchedule{
				{
					PeriodBegin: kit.Now(),
					PeriodEnd:   nil,
					Status:      domain.EvseStatusCharging,
				},
			},
			Capabilities: []string{domain.CapabilityChargingProfile},
			Coordinates: &domain.GeoLocation{
				Latitude:  "57.5464646",
				Longitude: "37.3543636",
			},
			ParkingRestrictions: []string{domain.ParkingRestrictionEvOnly},
			Images:              nil,
		},
		Connectors: nil,
	}
}

func (s *locationTestSuite) connector() *domain.Connector {
	return &domain.Connector{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     "EVG",
				CountryCode: "RS",
			},
			PlatformId:  "test",
			LastUpdated: kit.Now(),
			LastSent:    nil,
		},
		Id:         kit.NewId(),
		LocationId: kit.NewId(),
		EvseId:     kit.NewId(),
		Details: domain.ConnectorDetails{
			Standard:           domain.ConnectorTypeChademo,
			Format:             domain.FormatSocket,
			PowerType:          domain.PowerTypeDc,
			MaxVoltage:         10,
			MaxAmperage:        10,
			MaxElectricPower:   kit.Float64Ptr(10.0),
			TariffIds:          []string{kit.NewId()},
			TermsAndConditions: "https://test.com/terms",
		},
	}
}
