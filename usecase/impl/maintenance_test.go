package impl

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/stretchr/testify/suite"
	"testing"
)

type maintenanceTestSuite struct {
	kit.Suite
}

func (s *maintenanceTestSuite) SetupSuite() {
	s.Suite.Init(ocpi.LF())
}

func (s *maintenanceTestSuite) SetupTest() {
}

func (s *maintenanceTestSuite) TearDownSuite() {}

func TestMaintenanceSuite(t *testing.T) {
	suite.Run(t, new(maintenanceTestSuite))
}
