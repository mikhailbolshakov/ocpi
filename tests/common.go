package tests

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/sdk"
	"time"
)

func awaitLocation(ctx context.Context, sdk *sdk.Sdk, locId string, fn func(location *backend.Location) bool) (*backend.Location, error) {
	var loc *backend.Location
	if err := <-kit.Await(func() (bool, error) {
		var err error
		loc, err = sdk.GetLocation(ctx, locId)
		if err != nil {
			return false, err
		}
		if loc == nil {
			return false, nil
		}
		if fn != nil {
			return fn(loc), nil
		}
		return true, nil
	}, time.Millisecond*300, time.Second*3); err != nil {
		return nil, err
	}
	return loc, nil
}

func location(partyId string) *backend.Location {
	locId := kit.NewId()
	return &backend.Location{
		Id:               locId,
		Publish:          kit.BoolPtr(true),
		PublishAllowedTo: nil,
		Name:             "test location",
		Address:          "Beograd St. Despota Stefana",
		City:             "Beograd",
		PostalCode:       "567890",
		State:            "Beograd",
		Country:          "SRB",
		Coordinates: backend.GeoLocation{
			Latitude:  "53.4563463",
			Longitude: "43.2352352",
		},
		RelatedLocations: []*backend.AdditionalGeoLocation{
			{
				Latitude:  "53.4563463",
				Longitude: "43.2352352",
				Name: &backend.DisplayText{
					Language: "en",
					Text:     "Parking",
				},
			},
		},
		ParkingType: "ALONG_MOTORWAY",
		Directions: []*backend.DisplayText{
			{
				Language: "en",
				Text:     "On the right",
			},
		},
		Operator: &backend.BusinessDetails{
			Name: "operator",
		},
		SubOperator: &backend.BusinessDetails{
			Name: "sub operator",
		},
		Owner: &backend.BusinessDetails{
			Name: "owner",
			Inn:  "123456789",
		},
		Facilities: []string{"HOTEL", "CAFE"},
		TimeZone:   "Europe/Beograd",
		OpeningTimes: &backend.Hours{
			TwentyFourSeven: true,
			RegularHours: []*backend.RegularHours{
				{
					Weekday:     1,
					PeriodBegin: "10:00",
					PeriodEnd:   "20:00",
				},
			},
			ExceptionalOpenings: []*backend.ExceptionalPeriod{
				{
					PeriodBegin: kit.Now(),
					PeriodEnd:   kit.Now(),
				},
			},
		},
		ChargingWhenClosed: kit.BoolPtr(true),
		Images: []*backend.Image{
			{
				Url:      "https://chargers.com/image",
				Category: "CHARGER",
				Type:     "gif",
			},
		},
		EnergyMix: &backend.EnergyMix{
			IsGreenEnergy: true,
			EnergySources: []*backend.EnergySource{
				{
					Source:     "NUCLEAR",
					Percentage: 99,
				},
			},
			EnvironImpact: []*backend.EnvironmentalImpact{
				{
					Category: "CARBON_DIOXIDE",
					Amount:   50.5,
				},
			},
			SupplierName:      "noname",
			EnergyProductName: "product",
		},
		Evses:       []*backend.Evse{evse(partyId, locId)},
		PartyId:     partyId,
		CountryCode: "RS",
		RefId:       kit.NewId(),
		LastUpdated: kit.Now(),
	}
}

func evse(partyId, locId string) *backend.Evse {
	evseId := kit.NewId()
	return &backend.Evse{
		Id:         evseId,
		LocationId: locId,
		Status:     "AVAILABLE",
		EvseId:     "EVSE:1",
		StatusSchedule: []*backend.StatusSchedule{
			{
				PeriodBegin: kit.Now(),
				PeriodEnd:   kit.TimePtr(kit.Now().Add(time.Hour)),
				Status:      "INOPERATIVE",
			},
		},
		Capabilities: []string{"CHARGING_PROFILE_CAPABLE"},
		FloorLevel:   "4",
		Coordinates: &backend.GeoLocation{
			Latitude:  "53.4563463",
			Longitude: "43.2352352",
		},
		PhysicalReference: "On the right",
		Directions: []*backend.DisplayText{
			{
				Language: "en",
				Text:     "On the right",
			},
		},
		ParkingRestrictions: []string{"PLUGGED"},
		Images: []*backend.Image{
			{
				Url:      "https://chargers.com/image",
				Category: "CHARGER",
				Type:     "gif",
			},
		},
		Connectors:  []*backend.Connector{connector(partyId, locId, evseId)},
		PartyId:     partyId,
		CountryCode: "RS",
		RefId:       kit.NewId(),
		LastUpdated: kit.Now(),
	}
}

func connector(partyId, locId, evseId string) *backend.Connector {
	return &backend.Connector{
		Id:                 kit.NewId(),
		LocationId:         locId,
		EvseId:             evseId,
		Standard:           "CHADEMO",
		Format:             "SOCKET",
		PowerType:          "AC_1_PHASE",
		MaxVoltage:         10,
		MaxAmperage:        10,
		MaxElectricPower:   kit.Float64Ptr(150000),
		TariffIds:          []string{kit.NewRandString()},
		TermsAndConditions: "https://terms.com",
		PartyId:            partyId,
		CountryCode:        "RS",
		RefId:              kit.NewId(),
		LastUpdated:        kit.Now(),
	}
}

func session(partyId, authRef, locId, evseId, conId string) *backend.Session {
	return &backend.Session{
		Id:            kit.NewId(),
		StartDateTime: kit.TimePtr(kit.Now()),
		EndDateTime:   nil,
		Kwh:           kit.Float64Ptr(10000),
		CdrToken: &backend.CdrToken{
			PartyId:     partyId,
			CountryCode: "RS",
			Id:          kit.NewId(),
			Type:        "APP_USER",
			ContractId:  kit.NewId(),
		},
		AuthMethod:      backend.AuthMethodCommand,
		AuthRef:         authRef,
		LocationId:      locId,
		EvseId:          evseId,
		ConnectorId:     conId,
		MeterId:         kit.NewId(),
		Currency:        "RSD",
		ChargingPeriods: nil,
		TotalCost: &backend.Price{
			ExclVat: 100,
		},
		Status:      backend.SessionStatusActive,
		LastUpdated: kit.Now(),
		RefId:       kit.NewId(),
		PartyId:     partyId,
		CountryCode: "RS",
	}
}

func party() *backend.Party {
	return &backend.Party{
		Id:              kit.NewId(),
		Roles:           []string{backend.RoleCPO},
		PartyId:         "TST",
		CountryCode:     "RS",
		BusinessDetails: &backend.BusinessDetails{Name: "Company"},
		RefId:           kit.NewId(),
		Status:          backend.ConnectionStatusConnected,
		LastUpdated:     kit.Now(),
	}
}

func tariff(partyId string) *backend.Tariff {
	return &backend.Tariff{
		Id:       kit.NewId(),
		Currency: "RSD",
		Type:     "REGULAR",
		TariffAltText: []*backend.DisplayText{
			{
				Language: "en",
				Text:     "some tariff",
			},
		},
		TariffAltUrl: "https://ocpi.test/tariffs",
		MinPrice: &backend.Price{
			ExclVat: 1,
			InclVat: nil,
		},
		MaxPrice: &backend.Price{
			ExclVat: 99999999,
			InclVat: nil,
		},
		Elements: []*backend.TariffElement{
			{
				PriceComponents: []*backend.PriceComponent{
					{
						Type:     "ENERGY",
						Price:    100.0,
						Vat:      kit.Float64Ptr(10.0),
						StepSize: 1,
					},
				},
				Restrictions: &backend.TariffRestrictions{
					StartTime:   "07:00",
					EndTime:     "22:00",
					StartDate:   nil,
					EndDate:     nil,
					MinKwh:      kit.Float64Ptr(1),
					MaxKwh:      kit.Float64Ptr(9999999),
					MinCurrent:  kit.Float64Ptr(1),
					MaxCurrent:  kit.Float64Ptr(999),
					MinPower:    kit.Float64Ptr(1),
					MaxPower:    kit.Float64Ptr(99999999),
					MinDuration: kit.Float64Ptr(1),
					MaxDuration: kit.Float64Ptr(99999),
					DayOfWeek:   []string{"MONDAY"},
					Reservation: "",
				},
			},
		},
		StartDateTime: kit.TimePtr(kit.Now().Add(-time.Hour * 24)),
		EndDateTime:   kit.TimePtr(kit.Now().Add(time.Hour * 24)),
		EnergyMix:     nil,
		LastUpdated:   kit.Now(),
		RefId:         kit.NewId(),
		PartyId:       partyId,
		CountryCode:   "RS",
	}
}

func cdr(partyId string, sess *backend.Session, trf []*backend.Tariff) *backend.Cdr {
	return &backend.Cdr{
		Id:                       kit.NewId(),
		StartDateTime:            *sess.StartDateTime,
		EndDateTime:              *sess.EndDateTime,
		SessionId:                sess.Id,
		MeterId:                  sess.MeterId,
		Currency:                 sess.Currency,
		Tariffs:                  trf,
		ChargingPeriods:          sess.ChargingPeriods,
		TotalCost:                *sess.TotalCost,
		TotalEnergy:              *sess.Kwh,
		Remark:                   "remark",
		InvoiceReferenceId:       kit.NewId(),
		Credit:                   false,
		CreditReferenceId:        "",
		HomeChargingCompensation: false,
		LastUpdated:              kit.Now(),
		RefId:                    kit.NewId(),
		PartyId:                  partyId,
		CountryCode:              "RS",
	}
}

func awaitTariffs(ctx context.Context, sdk *sdk.Sdk, trfId string, fn func(trt *backend.Tariff) bool) (*backend.Tariff, error) {
	var trf *backend.Tariff
	if err := <-kit.Await(func() (bool, error) {
		var err error
		trf, err = sdk.GetTariff(ctx, trfId)
		if err != nil {
			return false, err
		}
		if trf == nil {
			return false, nil
		}
		if fn != nil {
			return fn(trf), nil
		}
		return true, nil
	}, time.Millisecond*300, time.Second*5); err != nil {
		return nil, err
	}
	return trf, nil
}

func awaitCommand(ctx context.Context, sdk *sdk.Sdk, cmdId, status string) (*backend.Command, error) {
	var cmd *backend.Command
	var err error
	if err := <-kit.Await(func() (bool, error) {
		cmd, err = sdk.GetCommand(ctx, cmdId)
		if err != nil {
			return false, err
		}
		return cmd != nil && cmd.Status == status, nil
	}, time.Millisecond*300, time.Second*5); err != nil {
		return nil, err
	}
	return cmd, nil
}

func awaitSession(ctx context.Context, sdk *sdk.Sdk, sessId string) (*backend.Session, error) {
	var sess *backend.Session
	var err error
	if err := <-kit.Await(func() (bool, error) {
		sess, err = sdk.GetSession(ctx, sessId)
		if err != nil {
			return false, err
		}
		return sess != nil, nil
	}, time.Millisecond*300, time.Second*5); err != nil {
		return nil, err
	}
	return sess, nil
}

func awaitCdr(ctx context.Context, sdk *sdk.Sdk, cdrId string) (*backend.Cdr, error) {
	var cdr *backend.Cdr
	var err error
	if err := <-kit.Await(func() (bool, error) {
		cdr, err = sdk.GetCdr(ctx, cdrId)
		if err != nil {
			return false, err
		}
		return cdr != nil, nil
	}, time.Millisecond*300, time.Second*5); err != nil {
		return nil, err
	}
	return cdr, nil
}
