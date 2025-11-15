package impl

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/stretchr/testify/suite"
	"testing"
)

type tariffUcTestSuite struct {
	kit.Suite
}

func (s *tariffUcTestSuite) SetupSuite() {
	s.Suite.Init(ocpi.LF())
}

func (s *tariffUcTestSuite) SetupTest() {
}

func (s *tariffUcTestSuite) TearDownSuite() {}

func TestTariffUcSuite(t *testing.T) {
	suite.Run(t, new(tariffUcTestSuite))
}
