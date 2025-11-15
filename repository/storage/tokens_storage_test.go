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

type tokensTestSuite struct {
	kit.Suite
	storage domain.TokenStorage
	adapter Adapter
}

func (s *tokensTestSuite) SetupSuite() {
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

func (s *tokensTestSuite) TearDownSuite() {
	_ = s.adapter.Close(s.Ctx)
}

func TestTokensSuite(t *testing.T) {
	suite.Run(t, new(tokensTestSuite))
}

func (s *tokensTestSuite) Test_CRUD() {
	// get when no exists
	act, err := s.storage.GetToken(s.Ctx, kit.NewId())
	s.NoError(err)
	s.Empty(act)

	// create new
	trf := s.token()
	s.NoError(s.storage.MergeToken(s.Ctx, trf))

	// get
	act, err = s.storage.GetToken(s.Ctx, trf.Id)
	s.NoError(err)
	s.Equal(act, trf)

	// update
	trf.Details.VisualNumber = "another"
	s.NoError(s.storage.MergeToken(s.Ctx, trf))

	// get
	act, err = s.storage.GetToken(s.Ctx, trf.Id)
	s.NoError(err)
	s.Equal(act, trf)

	// update
	trf.Details.VisualNumber = "another-2"
	s.NoError(s.storage.UpdateToken(s.Ctx, trf))

	// get
	act, err = s.storage.GetToken(s.Ctx, trf.Id)
	s.NoError(err)
	s.Equal(act, trf)
}

func (s *tokensTestSuite) Test_Search() {
	// create new
	tkn := s.token()
	s.NoError(s.storage.MergeToken(s.Ctx, tkn))

	rs, err := s.storage.SearchTokens(s.Ctx, &domain.TokenSearchCriteria{
		PageRequest: domain.PageRequest{},
		ExtId: &domain.PartyExtId{
			PartyId:     tkn.ExtId.PartyId,
			CountryCode: tkn.ExtId.CountryCode,
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

	rs, err = s.storage.SearchTokens(s.Ctx, &domain.TokenSearchCriteria{
		PageRequest: domain.PageRequest{
			DateFrom: kit.TimePtr(kit.Now().Add(-time.Hour)),
			DateTo:   kit.TimePtr(kit.Now().Add(time.Hour)),
			Limit:    kit.IntPtr(domain.PageSizeMaxLimit + 10),
		},
		IncPlatforms: []string{tkn.PlatformId},
		ExcPlatforms: []string{kit.NewId()},
	})
	s.NoError(err)
	s.Len(rs.Items, 1)
	s.NotEmpty(rs.Items[0].Id)
	s.NotEmpty(rs.Total)
	s.Equal(*rs.Total, 1)
	s.Equal(*rs.Limit, domain.PageSizeMaxLimit)

	rs, err = s.storage.SearchTokens(s.Ctx, &domain.TokenSearchCriteria{
		PageRequest: domain.PageRequest{},
		Ids:         []string{tkn.Id},
	})
	s.NoError(err)
	s.Len(rs.Items, 1)
}

func (s *tokensTestSuite) Test_DeleteByExt() {

	tkn := s.token()

	// create a token
	s.NoError(s.storage.MergeToken(s.Ctx, tkn))

	// get token
	act, err := s.storage.GetToken(s.Ctx, tkn.Id)
	s.NoError(err)
	s.Equal(act, tkn)

	// delete by ext
	s.NoError(s.storage.DeleteTokensByExtId(s.Ctx, tkn.ExtId))

	// get token
	act, err = s.storage.GetToken(s.Ctx, tkn.Id)
	s.NoError(err)
	s.Empty(act)

}

func (s *tokensTestSuite) token() *domain.Token {
	partyId := kit.NewRandString()
	return &domain.Token{
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
		Details: domain.TokenDetails{
			Type:               domain.TokenTypeRfid,
			ContractId:         kit.NewRandString(),
			VisualNumber:       "num",
			Issuer:             "issuer",
			GroupId:            kit.NewRandString(),
			Valid:              kit.BoolPtr(true),
			WhiteList:          domain.TokenWLTypeNever,
			Lang:               "en",
			DefaultProfileType: domain.ProfileTypeCheap,
			EnergyContract:     nil,
		},
	}
}
