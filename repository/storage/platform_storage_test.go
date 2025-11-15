//go:build integration

package storage

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/stretchr/testify/suite"
	"testing"
)

type platformStorageTestSuite struct {
	kit.Suite
	storage domain.PlatformStorage
	adapter Adapter
}

func (s *platformStorageTestSuite) SetupSuite() {
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

func (s *platformStorageTestSuite) TearDownSuite() {
	_ = s.adapter.Close(s.Ctx)
}

func TestPlatformStorageSuite(t *testing.T) {
	suite.Run(t, new(platformStorageTestSuite))
}

func (s *platformStorageTestSuite) platform() *domain.Platform {
	return &domain.Platform{
		Id:          kit.NewRandString(),
		TokenA:      domain.PlatformToken(kit.NewRandString()),
		TokenB:      domain.PlatformToken(kit.NewRandString()),
		TokenC:      domain.PlatformToken(kit.NewRandString()),
		TokenBase64: kit.BoolPtr(true),
		Name:        "test",
		Role:        domain.RoleOTHER,
		VersionInfo: domain.VersionInfo{
			Current: "2.2.1",
			Available: map[string]domain.Endpoint{
				"2.2.1": "https://test.dev/ocpi/2.2.1",
			},
			VersionEp: "https://test.dev/ocpi/2.2.1/versions",
		},
		Endpoints: map[string]domain.RoleEndpoint{
			"credentials": map[string]domain.Endpoint{
				"SENDER": "https://test.dev/ocpi/2.2.1/sender/credentials",
			},
		},
		Status: domain.ConnectionStatusSuspended,
		Remote: true,
		Protocol: &domain.ProtocolDetails{
			PushSupport: domain.PushSupport{
				Locations: true,
			},
		},
	}
}

func (s *platformStorageTestSuite) Test_CRUD() {
	// get when not exists
	act, err := s.storage.GetPlatform(s.Ctx, kit.NewId())
	s.NoError(err)
	s.Empty(act)

	// get by token when not exists
	act, err = s.storage.GetPlatformByTokenA(s.Ctx, domain.PlatformToken(kit.NewRandString()))
	s.NoError(err)
	s.Empty(act)

	// create
	p := s.platform()
	s.NoError(s.storage.CreatePlatform(s.Ctx, p))

	// get
	act, err = s.storage.GetPlatform(s.Ctx, p.Id)
	s.NoError(err)
	s.Equal(act, p)

	// get by token
	act, err = s.storage.GetPlatformByTokenA(s.Ctx, p.TokenA)
	s.NoError(err)
	s.Equal(act, p)

	// update
	p.Name = "another name"
	s.NoError(s.storage.UpdatePlatform(s.Ctx, p))

	// get
	act, err = s.storage.GetPlatform(s.Ctx, p.Id)
	s.NoError(err)
	s.Equal(act, p)

	// delete
	s.NoError(s.storage.DeletePlatform(s.Ctx, p))

	// get
	act, err = s.storage.GetPlatform(s.Ctx, p.Id)
	s.NoError(err)
	s.Empty(act)

	// get by token when not exists
	act, err = s.storage.GetPlatformByTokenA(s.Ctx, domain.PlatformToken(kit.NewRandString()))
	s.NoError(err)
	s.Empty(act)
}

func (s *platformStorageTestSuite) Test_Search() {
	// create
	p1 := s.platform()
	s.NoError(s.storage.CreatePlatform(s.Ctx, p1))
	p2 := s.platform()
	s.NoError(s.storage.CreatePlatform(s.Ctx, p2))

	rs, err := s.storage.SearchPlatforms(s.Ctx, &domain.PlatformSearchCriteria{
		Roles:    []string{p1.Role, p2.Role},
		Statuses: []string{p1.Status},
		IncIds:   []string{p1.Id, p2.Id},
	})
	s.NoError(err)
	s.Len(rs, 2)

	rs, err = s.storage.SearchPlatforms(s.Ctx, &domain.PlatformSearchCriteria{
		IncIds: []string{p1.Id, p2.Id},
		Remote: kit.BoolPtr(true),
	})
	s.NoError(err)
	s.Len(rs, 2)
}
