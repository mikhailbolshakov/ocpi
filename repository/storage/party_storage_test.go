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

type partyTestSuite struct {
	kit.Suite
	storage domain.PartyStorage
	adapter Adapter
}

func (s *partyTestSuite) SetupSuite() {
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

func (s *partyTestSuite) TearDownSuite() {
	_ = s.adapter.Close(s.Ctx)
}

func TestPartySuite(t *testing.T) {
	suite.Run(t, new(partyTestSuite))
}

func (s *partyTestSuite) Test_GetWhenEmpty() {
	// get by id when not exists
	act, err := s.storage.GetParty(s.Ctx, kit.NewId())
	s.NoError(err)
	s.Empty(act)

	// get by ref when not exists
	act, err = s.storage.GetByRefId(s.Ctx, kit.NewId())
	s.NoError(err)
	s.Empty(act)

	// get by ref when not exists
	act, err = s.storage.GetPartyByExtId(s.Ctx, domain.PartyExtId{PartyId: kit.NewId(), CountryCode: "RS"})
	s.NoError(err)
	s.Empty(act)

	// get by ref when not exists
	acts, err := s.storage.GetPartiesByPlatform(s.Ctx, kit.NewId())
	s.NoError(err)
	s.Empty(acts)
}

func (s *partyTestSuite) Test_CRUD() {
	// create new
	p := s.party()
	s.NoError(s.storage.CreateParty(s.Ctx, p))

	// get by id when not exists
	act, err := s.storage.GetParty(s.Ctx, p.Id)
	s.NoError(err)
	s.Equal(p, act)

	// get by ref when not exists
	act, err = s.storage.GetByRefId(s.Ctx, p.RefId)
	s.NoError(err)
	s.Equal(p, act)

	// get by ref when not exists
	act, err = s.storage.GetPartyByExtId(s.Ctx, domain.PartyExtId{PartyId: p.ExtId.PartyId, CountryCode: p.ExtId.CountryCode})
	s.NoError(err)
	s.Equal(p, act)

	// get by ref when not exists
	acts, err := s.storage.GetPartiesByPlatform(s.Ctx, p.PlatformId)
	s.NoError(err)
	s.NotEmpty(acts)
	s.Equal(p, acts[0])

	// update
	p.BusinessDetails.Name = "another"
	s.NoError(s.storage.UpdateParty(s.Ctx, p))

	// get by id when not exists
	act, err = s.storage.GetParty(s.Ctx, p.Id)
	s.NoError(err)
	s.Equal(p, act)

	// get by ref when not exists
	act, err = s.storage.GetByRefId(s.Ctx, p.RefId)
	s.NoError(err)
	s.Equal(p, act)

	// get by ref when not exists
	act, err = s.storage.GetPartyByExtId(s.Ctx, domain.PartyExtId{PartyId: p.ExtId.PartyId, CountryCode: p.ExtId.CountryCode})
	s.NoError(err)
	s.Equal(p, act)

	// get by ref when not exists
	acts, err = s.storage.GetPartiesByPlatform(s.Ctx, p.PlatformId)
	s.NoError(err)
	s.NotEmpty(acts)
	s.Equal(p, acts[0])
}

func (s *partyTestSuite) Test_Search() {
	// create new
	p := s.party()
	s.NoError(s.storage.CreateParty(s.Ctx, p))

	rs, err := s.storage.Search(s.Ctx, &domain.PartySearchCriteria{
		PageRequest:  domain.PageRequest{},
		IncRoles:     []string{p.Roles[0]},
		ExcRoles:     []string{domain.RoleOTHER},
		IncPlatforms: []string{p.PlatformId},
		ExcPlatforms: []string{kit.NewId()},
	})
	s.NoError(err)
	s.Len(rs.Items, 1)

	rs, err = s.storage.Search(s.Ctx, &domain.PartySearchCriteria{
		PageRequest: domain.PageRequest{
			DateFrom: kit.TimePtr(kit.Now().Add(-time.Hour)),
			DateTo:   kit.TimePtr(kit.Now().Add(time.Hour)),
		},
		ExtId: &domain.PartyExtId{
			PartyId:     p.ExtId.PartyId,
			CountryCode: p.ExtId.CountryCode,
		},
	})
	s.NoError(err)
	s.Len(rs.Items, 1)
	s.NotEmpty(rs.Items[0].Id)
	s.NotEmpty(rs.Total)
	s.Equal(*rs.Total, 1)
	s.Equal(*rs.Limit, domain.PageSizeDefault)

	rs, err = s.storage.Search(s.Ctx, &domain.PartySearchCriteria{
		PageRequest: domain.PageRequest{
			Limit: kit.IntPtr(domain.PageSizeMaxLimit + 10),
		},
		RefId: p.RefId,
	})
	s.NoError(err)
	s.Len(rs.Items, 1)
	s.NotEmpty(rs.Items[0].Id)
	s.NotEmpty(rs.Total)
	s.Equal(*rs.Total, 1)
	s.Equal(*rs.Limit, domain.PageSizeMaxLimit)

	rs, err = s.storage.Search(s.Ctx, &domain.PartySearchCriteria{
		PageRequest: domain.PageRequest{},
		Ids:         []string{p.Id},
	})
	s.NoError(err)
	s.Len(rs.Items, 1)
}

func (s *partyTestSuite) Test_DeleteByExt() {

	p := s.party()

	// create a party
	s.NoError(s.storage.CreateParty(s.Ctx, p))

	// get party
	act, err := s.storage.GetParty(s.Ctx, p.Id)
	s.NoError(err)
	s.Equal(act, p)

	// delete by ext
	s.NoError(s.storage.DeletePartyByExtId(s.Ctx, p.ExtId))

	// get party
	act, err = s.storage.GetParty(s.Ctx, p.Id)
	s.NoError(err)
	s.Empty(act)

}

func (s *partyTestSuite) party() *domain.Party {
	return &domain.Party{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     kit.NewId(),
				CountryCode: "RS",
			},
			PlatformId:  kit.NewRandString(),
			RefId:       kit.NewRandString(),
			LastUpdated: kit.Now(),
		},
		Id:    kit.NewId(),
		Roles: []string{domain.RoleCPO, domain.RoleEMSP},
		BusinessDetails: &domain.BusinessDetails{
			Name:    "name",
			Website: "https://test.com/site",
			Logo: &domain.Image{
				Url:       "https://test.com/logo",
				Thumbnail: "https://test.com/th",
				Category:  domain.ImageCategoryOperator,
				Type:      "jpg",
				Width:     100,
				Height:    100,
			},
		},
		Status: domain.ConnectionStatusConnected,
	}
}
