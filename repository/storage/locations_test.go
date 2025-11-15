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

type locationsTestSuite struct {
	kit.Suite
	storage domain.LocationStorage
	adapter Adapter
}

func (s *locationsTestSuite) SetupSuite() {
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

func (s *locationsTestSuite) TearDownSuite() {
	_ = s.adapter.Close(s.Ctx)
}

func TestLocationsSuite(t *testing.T) {
	suite.Run(t, new(locationsTestSuite))
}

func (s *locationsTestSuite) Test_Location_CRUD() {

	// get when not exists with evse
	act, err := s.storage.GetLocation(s.Ctx, kit.NewId(), true)
	s.NoError(err)
	s.Empty(act)

	// get when not exists without evse
	act, err = s.storage.GetLocation(s.Ctx, kit.NewId(), false)
	s.NoError(err)
	s.Empty(act)

	// merge when not exists
	loc := s.location()
	s.NoError(s.storage.MergeLocation(s.Ctx, loc))

	// get with evse
	act, err = s.storage.GetLocation(s.Ctx, loc.Id, true)
	s.NoError(err)
	s.Equal(act, loc)

	// get without evse
	act, err = s.storage.GetLocation(s.Ctx, loc.Id, false)
	s.NoError(err)
	s.NotEmpty(act)
	s.Empty(act.Evses)

	// merge when exists
	loc.Details.Name = "another"
	loc.Evses[0].Status = domain.EvseStatusCharging
	loc.Evses[0].Connectors[0].Details.PowerType = domain.PowerTypeDc
	s.NoError(s.storage.MergeLocation(s.Ctx, loc))

	// get with evse
	act, err = s.storage.GetLocation(s.Ctx, loc.Id, true)
	s.NoError(err)
	s.Equal(act, loc)

	// get without evse
	act, err = s.storage.GetLocation(s.Ctx, loc.Id, false)
	s.NoError(err)
	s.NotEmpty(act)
	s.Empty(act.Evses)

	// update
	loc.Details.Name = "another-2"
	s.NoError(s.storage.UpdateLocation(s.Ctx, loc))

	// get without evse
	act, err = s.storage.GetLocation(s.Ctx, loc.Id, false)
	s.NoError(err)
	s.NotEmpty(act)
	s.Empty(act.Evses)

}

func (s *locationsTestSuite) Test_Location_Search() {
	// merge when not exists
	loc := s.location()
	s.NoError(s.storage.MergeLocation(s.Ctx, loc))

	rs, err := s.storage.SearchLocations(s.Ctx, &domain.LocationSearchCriteria{
		PageRequest:  domain.PageRequest{},
		ExtId:        nil,
		IncPlatforms: []string{loc.PlatformId},
		ExcPlatforms: []string{kit.NewId()},
	})
	s.NoError(err)
	s.Len(rs.Items, 1)
	s.NotEmpty(rs.Items[0].Id)
	s.NotEmpty(rs.Total)
	s.Equal(*rs.Total, 1)
	s.Equal(*rs.Limit, domain.PageSizeDefault)

	rs, err = s.storage.SearchLocations(s.Ctx, &domain.LocationSearchCriteria{
		PageRequest: domain.PageRequest{
			DateFrom: kit.TimePtr(kit.Now().Add(-time.Hour)),
			DateTo:   kit.TimePtr(kit.Now().Add(time.Hour)),
			Limit:    kit.IntPtr(domain.PageSizeMaxLimit + 10),
		},
		ExtId: &domain.PartyExtId{
			PartyId:     loc.ExtId.PartyId,
			CountryCode: loc.ExtId.CountryCode,
		},
	})
	s.NoError(err)
	s.Len(rs.Items, 1)
	s.NotEmpty(rs.Items[0].Id)
	s.NotEmpty(rs.Total)
	s.Equal(*rs.Total, 1)
	s.Equal(*rs.Limit, domain.PageSizeMaxLimit)

	rs, err = s.storage.SearchLocations(s.Ctx, &domain.LocationSearchCriteria{
		PageRequest: domain.PageRequest{},
		Ids:         []string{loc.Id},
	})
	s.NoError(err)
	s.Len(rs.Items, 1)

}

func (s *locationsTestSuite) Test_DeleteLocationByExt() {

	// merge when not exists
	loc := s.location()
	s.NoError(s.storage.MergeLocation(s.Ctx, loc))

	// get when not exists without evse
	act, err := s.storage.GetLocation(s.Ctx, loc.Id, true)
	s.NoError(err)
	s.NotEmpty(act)

	// delete by ext
	s.NoError(s.storage.DeleteLocationsByExtId(s.Ctx, loc.ExtId))

	// check
	act, err = s.storage.GetLocation(s.Ctx, loc.Id, true)
	s.NoError(err)
	s.Empty(act)

	evse, err := s.storage.GetEvse(s.Ctx, loc.Id, loc.Evses[0].Id, true)
	s.NoError(err)
	s.Empty(evse)

	con, err := s.storage.GetConnector(s.Ctx, loc.Id, loc.Evses[0].Id, loc.Evses[0].Connectors[0].Id)
	s.NoError(err)
	s.Empty(con)
}

func (s *locationsTestSuite) Test_Evse_CRUD() {

	// get
	act, err := s.storage.GetEvse(s.Ctx, kit.NewId(), kit.NewId(), true)
	s.NoError(err)
	s.Empty(act)

	// get without connectors
	act, err = s.storage.GetEvse(s.Ctx, kit.NewId(), kit.NewId(), false)
	s.NoError(err)
	s.Empty(act)

	// merge when not exists
	evse := s.evse(kit.NewId(), kit.NewId())
	s.NoError(s.storage.MergeEvse(s.Ctx, evse))

	// get
	act, err = s.storage.GetEvse(s.Ctx, evse.LocationId, evse.Id, true)
	s.NoError(err)
	s.Equal(act, evse)

	// get without connectors
	act, err = s.storage.GetEvse(s.Ctx, evse.LocationId, evse.Id, false)
	s.NoError(err)
	s.NotEmpty(act)
	s.Equal(act.Details, evse.Details)
	s.Empty(act.Connectors)

	// merge when exists
	evse.Status = domain.EvseStatusCharging
	evse.Connectors[0].Details.PowerType = domain.PowerTypeDc
	s.NoError(s.storage.MergeEvse(s.Ctx, evse))

	// get
	act, err = s.storage.GetEvse(s.Ctx, evse.LocationId, evse.Id, true)
	s.NoError(err)
	s.Equal(act, evse)

	// update
	evse.Status = domain.EvseStatusBlocked
	s.NoError(s.storage.UpdateEvse(s.Ctx, evse))

	// get
	act, err = s.storage.GetEvse(s.Ctx, evse.LocationId, evse.Id, true)
	s.NoError(err)
	s.Equal(act, evse)

}

func (s *locationsTestSuite) Test_Evse_Search() {
	// merge when not exists
	evse := s.evse(kit.NewId(), kit.NewId())
	s.NoError(s.storage.MergeEvse(s.Ctx, evse))

	rs, err := s.storage.SearchEvses(s.Ctx, &domain.EvseSearchCriteria{
		PageRequest: domain.PageRequest{},
		ExtId: &domain.PartyExtId{
			PartyId:     evse.ExtId.PartyId,
			CountryCode: evse.ExtId.CountryCode,
		},
	})
	s.NoError(err)
	s.Len(rs.Items, 1)
	s.NotEmpty(rs.Items[0].Id)
	s.NotEmpty(rs.Total)
	s.Equal(*rs.Total, 1)
	s.Equal(*rs.Limit, domain.PageSizeDefault)

	rs, err = s.storage.SearchEvses(s.Ctx, &domain.EvseSearchCriteria{
		PageRequest: domain.PageRequest{
			DateFrom: kit.TimePtr(kit.Now().Add(-time.Hour)),
			DateTo:   kit.TimePtr(kit.Now().Add(time.Hour)),
			Limit:    kit.IntPtr(domain.PageSizeMaxLimit + 10),
		},
		ExtId: &domain.PartyExtId{
			PartyId:     evse.ExtId.PartyId,
			CountryCode: evse.ExtId.CountryCode,
		},
	})
	s.NoError(err)
	s.Len(rs.Items, 1)
	s.NotEmpty(rs.Items[0].Id)
	s.NotEmpty(rs.Total)
	s.Equal(*rs.Total, 1)
	s.Equal(*rs.Limit, domain.PageSizeMaxLimit)
}

func (s *locationsTestSuite) Test_Connector_CRUD() {

	// get
	act, err := s.storage.GetConnector(s.Ctx, kit.NewId(), kit.NewId(), kit.NewId())
	s.NoError(err)
	s.Empty(act)

	// merge when not exists
	con := s.connector(kit.NewId(), kit.NewId(), kit.NewId())
	s.NoError(s.storage.MergeConnector(s.Ctx, con))

	// get
	act, err = s.storage.GetConnector(s.Ctx, con.LocationId, con.EvseId, con.Id)
	s.NoError(err)
	s.Equal(act, con)

	// merge when exists
	con.Details.PowerType = domain.PowerTypeDc
	s.NoError(s.storage.MergeConnector(s.Ctx, con))

	// get
	act, err = s.storage.GetConnector(s.Ctx, con.LocationId, con.EvseId, con.Id)
	s.NoError(err)
	s.Equal(act, con)

	// update
	con.Details.PowerType = domain.PowerTypeAc3Phase
	s.NoError(s.storage.UpdateConnector(s.Ctx, con))

	// get
	act, err = s.storage.GetConnector(s.Ctx, con.LocationId, con.EvseId, con.Id)
	s.NoError(err)
	s.Equal(act, con)

}

func (s *locationsTestSuite) Test_Connector_Search() {
	// merge when not exists
	evse := s.evse(kit.NewId(), kit.NewId())
	s.NoError(s.storage.MergeEvse(s.Ctx, evse))

	rs, err := s.storage.SearchEvses(s.Ctx, &domain.EvseSearchCriteria{
		PageRequest: domain.PageRequest{},
		ExtId: &domain.PartyExtId{
			PartyId:     evse.ExtId.PartyId,
			CountryCode: evse.ExtId.CountryCode,
		},
	})
	s.NoError(err)
	s.Len(rs.Items, 1)
	s.NotEmpty(rs.Items[0].Id)
	s.NotEmpty(rs.Total)
	s.Equal(*rs.Total, 1)
	s.Equal(*rs.Limit, domain.PageSizeDefault)

	rs, err = s.storage.SearchEvses(s.Ctx, &domain.EvseSearchCriteria{
		PageRequest: domain.PageRequest{
			DateFrom: kit.TimePtr(kit.Now().Add(-time.Hour)),
			DateTo:   kit.TimePtr(kit.Now().Add(time.Hour)),
			Limit:    kit.IntPtr(domain.PageSizeMaxLimit + 10),
		},
		ExtId: &domain.PartyExtId{
			PartyId:     evse.ExtId.PartyId,
			CountryCode: evse.ExtId.CountryCode,
		},
	})
	s.NoError(err)
	s.Len(rs.Items, 1)
	s.NotEmpty(rs.Items[0].Id)
	s.NotEmpty(rs.Total)
	s.Equal(*rs.Total, 1)
	s.Equal(*rs.Limit, domain.PageSizeMaxLimit)
}

func (s *locationsTestSuite) location() *domain.Location {
	l := &domain.Location{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     kit.NewRandString(),
				CountryCode: "RS",
			},
			PlatformId:  kit.NewRandString(),
			RefId:       kit.NewId(),
			LastUpdated: kit.Now(),
		},
		Id: kit.NewId(),
		Details: domain.LocationDetails{
			Publish:          kit.BoolPtr(true),
			PublishAllowedTo: nil,
			Name:             "test location",
			Address:          "Beograd St. Despota Stefana",
			City:             "Beograd",
			PostalCode:       "567890",
			State:            "Beograd",
			Country:          "SRB",
			Coordinates: domain.GeoLocation{
				Latitude:  "53.4563463",
				Longitude: "43.2352352",
			},
			RelatedLocations: []*domain.AdditionalGeoLocation{
				{
					GeoLocation: domain.GeoLocation{
						Latitude:  "53.4563463",
						Longitude: "43.2352352",
					},
					Name: &domain.DisplayText{
						Language: "en",
						Text:     "Parking",
					},
				},
			},
			ParkingType: "ALONG_MOTORWAY",
			Directions: []*domain.DisplayText{
				{
					Language: "en",
					Text:     "On the right",
				},
			},
			Operator: &domain.BusinessDetails{
				Name: "operator",
			},
			SubOperator: &domain.BusinessDetails{
				Name: "sub operator",
			},
			Owner: &domain.BusinessDetails{
				Name: "owner",
			},
			Facilities: []string{"HOTEL", "CAFE"},
			TimeZone:   "Europe/Beograd",
			OpeningTimes: &domain.Hours{
				TwentyFourSeven: true,
				RegularHours: []*domain.RegularHours{
					{
						Weekday:     1,
						PeriodBegin: "10:00",
						PeriodEnd:   "20:00",
					},
				},
				ExceptionalOpenings: []*domain.ExceptionalPeriod{
					{
						PeriodBegin: kit.Now(),
						PeriodEnd:   kit.Now(),
					},
				},
			},
			ChargingWhenClosed: kit.BoolPtr(true),
			Images: []*domain.Image{
				{
					Url:      "https://chargers.com/image",
					Category: "CHARGER",
					Type:     "gif",
				},
			},
			EnergyMix: &domain.EnergyMix{
				IsGreenEnergy: true,
				EnergySources: []*domain.EnergySource{
					{
						Source:     "NUCLEAR",
						Percentage: 99,
					},
				},
				EnvironImpact: []*domain.EnvironmentalImpact{
					{
						Category: "CARBON_DIOXIDE",
						Amount:   50.5,
					},
				},
				SupplierName:      "noname",
				EnergyProductName: "product",
			},
		},
	}
	l.Evses = []*domain.Evse{s.evse(l.ExtId.PartyId, l.Id)}
	return l
}

func (s *locationsTestSuite) evse(partyId, locId string) *domain.Evse {
	evse := &domain.Evse{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     partyId,
				CountryCode: "RS",
			},
			PlatformId:  kit.NewRandString(),
			RefId:       kit.NewId(),
			LastUpdated: kit.Now(),
		},
		Id:         kit.NewId(),
		LocationId: locId,
		Status:     domain.EvseStatusAvailable,
		Details: domain.EvseDetails{
			EvseId: kit.NewRandString(),
			StatusSchedule: []*domain.StatusSchedule{
				{
					PeriodBegin: kit.Now(),
					PeriodEnd:   kit.TimePtr(kit.Now().Add(time.Hour)),
					Status:      "INOPERATIVE",
				},
			},
			Capabilities: []string{"CHARGING_PROFILE_CAPABLE"},
			FloorLevel:   "4",
			Coordinates: &domain.GeoLocation{
				Latitude:  "53.4563463",
				Longitude: "43.2352352",
			},
			PhysicalReference: "On the right",
			Directions: []*domain.DisplayText{
				{
					Language: "en",
					Text:     "On the right",
				},
			},
			ParkingRestrictions: []string{"PLUGGED"},
			Images: []*domain.Image{
				{
					Url:      "https://chargers.com/image",
					Category: "CHARGER",
					Type:     "gif",
				},
			},
		},
	}
	evse.Connectors = []*domain.Connector{s.connector(partyId, locId, evse.Id)}
	return evse
}

func (s *locationsTestSuite) connector(partyId, locId, evseId string) *domain.Connector {
	return &domain.Connector{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     partyId,
				CountryCode: "RS",
			},
			PlatformId:  kit.NewRandString(),
			RefId:       kit.NewId(),
			LastUpdated: kit.Now(),
		},
		Id:         kit.NewId(),
		LocationId: locId,
		EvseId:     evseId,
		Details: domain.ConnectorDetails{
			Standard:           "CHADEMO",
			Format:             "SOCKET",
			PowerType:          "AC_1_PHASE",
			MaxVoltage:         10,
			MaxAmperage:        10,
			MaxElectricPower:   kit.Float64Ptr(150000),
			TariffIds:          []string{kit.NewRandString()},
			TermsAndConditions: "https://terms.com",
		},
	}
}
