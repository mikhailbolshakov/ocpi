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

type commandsTestSuite struct {
	kit.Suite
	storage domain.CommandStorage
	adapter Adapter
}

func (s *commandsTestSuite) SetupSuite() {
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

func (s *commandsTestSuite) TearDownSuite() {
	_ = s.adapter.Close(s.Ctx)
}

func TestCommandsSuite(t *testing.T) {
	suite.Run(t, new(commandsTestSuite))
}

func (s *commandsTestSuite) Test_CRUD() {
	// get command
	act, err := s.storage.GetCommand(s.Ctx, kit.NewId())
	s.NoError(err)
	s.Empty(act)

	cmd := s.command()

	// create a command
	s.NoError(s.storage.CreateCommand(s.Ctx, cmd))

	// get command
	act, err = s.storage.GetCommand(s.Ctx, cmd.Id)
	s.NoError(err)
	s.Equal(act, cmd)

	// update command
	cmd.Status = domain.CmdStatusRequestProcessedOk
	cmd.LastUpdated = kit.Now()
	s.NoError(s.storage.UpdateCommand(s.Ctx, cmd))

	// get command
	act, err = s.storage.GetCommand(s.Ctx, cmd.Id)
	s.NoError(err)
	s.Equal(act, cmd)

	// get command by auth ref
	act, err = s.storage.GetCommandByAuthRef(s.Ctx, cmd.AuthRef)
	s.NoError(err)
	s.Equal(act, cmd)

	// delete a command
	s.NoError(s.storage.DeleteCommand(s.Ctx, cmd.Id))

	// get command
	act, err = s.storage.GetCommand(s.Ctx, cmd.Id)
	s.NoError(err)
	s.Empty(act)

}

func (s *commandsTestSuite) Test_Search() {
	// create new
	cmd := s.command()
	s.NoError(s.storage.CreateCommand(s.Ctx, cmd))
	defer func() { _ = s.storage.DeleteCommand(s.Ctx, cmd.Id) }()

	rs, err := s.storage.SearchCommands(s.Ctx, &domain.CommandSearchCriteria{
		PageRequest: domain.PageRequest{},
		ExtId: &domain.PartyExtId{
			PartyId:     cmd.ExtId.PartyId,
			CountryCode: cmd.ExtId.CountryCode,
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

	rs, err = s.storage.SearchCommands(s.Ctx, &domain.CommandSearchCriteria{
		PageRequest: domain.PageRequest{
			DateFrom: kit.TimePtr(kit.Now().Add(-time.Hour)),
			DateTo:   kit.TimePtr(kit.Now().Add(time.Hour)),
			Limit:    kit.IntPtr(domain.PageSizeMaxLimit + 10),
		},
		IncPlatforms: []string{cmd.PlatformId},
		ExcPlatforms: []string{kit.NewId()},
	})
	s.NoError(err)
	s.Len(rs.Items, 1)
	s.NotEmpty(rs.Items[0].Id)
	s.NotEmpty(rs.Total)
	s.Equal(*rs.Total, 1)
	s.Equal(*rs.Limit, domain.PageSizeMaxLimit)

	rs, err = s.storage.SearchCommands(s.Ctx, &domain.CommandSearchCriteria{
		PageRequest: domain.PageRequest{
			Limit: kit.IntPtr(1),
		},
		AuthRef: cmd.AuthRef,
		Cmd:     cmd.Cmd,
	})
	s.NoError(err)
	s.Len(rs.Items, 1)

	rs, err = s.storage.SearchCommands(s.Ctx, &domain.CommandSearchCriteria{
		PageRequest: domain.PageRequest{
			Limit: kit.IntPtr(1),
		},
		Ids: []string{cmd.Id},
	})
	s.NoError(err)
	s.Len(rs.Items, 1)

	// create new
	cmd1 := s.command()
	cmd1.Status = domain.CmdStatusRequestAccepted
	cmd1.Deadline = kit.Now().Add(-time.Hour)
	s.NoError(s.storage.CreateCommand(s.Ctx, cmd1))
	defer func() { _ = s.storage.DeleteCommand(s.Ctx, cmd1.Id) }()

	cmd2 := s.command()
	cmd2.Status = domain.CmdStatusRequestAccepted
	cmd2.Deadline = kit.Now().Add(-time.Hour)
	s.NoError(s.storage.CreateCommand(s.Ctx, cmd2))
	defer func() { _ = s.storage.DeleteCommand(s.Ctx, cmd2.Id) }()

	rs, err = s.storage.SearchCommands(s.Ctx, &domain.CommandSearchCriteria{
		PageRequest: domain.PageRequest{
			Limit: kit.IntPtr(1),
		},
		Ids:         []string{cmd1.Id, cmd2.Id},
		RetrieveAll: true,
		DeadlineLE:  kit.NowPtr(),
		Statuses:    []string{domain.CmdStatusRequestAccepted},
	})
	s.NoError(err)
	s.Len(rs.Items, 2)

}

func (s *commandsTestSuite) Test_SearchByReservationId() {
	// create new
	cmd := s.command()
	cmd.Details.StartSession = nil
	cmd.Cmd = domain.CmdReserve
	cmd.Details.Reserve = &domain.ReserveNow{
		ExpireDate:    kit.Now(),
		ReservationId: kit.NewRandString(),
		LocationId:    kit.NewId(),
		EvseId:        kit.NewId(),
		ConnectorId:   kit.NewId(),
	}

	s.NoError(s.storage.CreateCommand(s.Ctx, cmd))
	defer func() { _ = s.storage.DeleteCommand(s.Ctx, cmd.Id) }()

	rs, err := s.storage.SearchCommands(s.Ctx, &domain.CommandSearchCriteria{
		PageRequest:   domain.PageRequest{},
		ReservationId: cmd.Details.Reserve.ReservationId,
	})
	s.NoError(err)
	s.Len(rs.Items, 1)
	s.NotEmpty(rs.Items[0].Id)
	s.NotEmpty(rs.Total)
	s.Equal(*rs.Total, 1)
	s.Equal(*rs.Limit, domain.PageSizeDefault)
}

func (s *commandsTestSuite) Test_DeleteByExt() {

	cmd := s.command()

	// create a command
	s.NoError(s.storage.CreateCommand(s.Ctx, cmd))

	// get command
	act, err := s.storage.GetCommand(s.Ctx, cmd.Id)
	s.NoError(err)
	s.Equal(act, cmd)

	// delete by ext
	s.NoError(s.storage.DeleteCommandsByExt(s.Ctx, cmd.ExtId))

	// get command
	act, err = s.storage.GetCommand(s.Ctx, cmd.Id)
	s.NoError(err)
	s.Empty(act)

}

func (s *commandsTestSuite) command() *domain.Command {
	return &domain.Command{
		OcpiItem: domain.OcpiItem{
			ExtId: domain.PartyExtId{
				PartyId:     kit.NewRandString(),
				CountryCode: "RS",
			},
			PlatformId:  kit.NewRandString(),
			RefId:       kit.NewRandString(),
			LastUpdated: kit.Now(),
		},
		Id:       kit.NewId(),
		Status:   domain.CmdStatusRequestAccepted,
		Cmd:      domain.CmdStartSession,
		Deadline: kit.Now().Add(time.Hour * 24 * 365),
		Details: domain.CommandDetails{
			ResponseUrl: "https://test-url",
			Processing: domain.Processing{
				Status: domain.CmdResultTypeFailed,
				ErrMsg: "Error",
			},
			StartSession: &domain.StartSession{
				LocationId:  kit.NewId(),
				EvseId:      kit.NewId(),
				ConnectorId: kit.NewId(),
			},
		},
		AuthRef: kit.NewRandString(),
	}
}
