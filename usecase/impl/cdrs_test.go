package impl

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/stretchr/testify/suite"
	"testing"
)

type cdrUcTestSuite struct {
	kit.Suite
}

func (s *cdrUcTestSuite) SetupSuite() {
	s.Suite.Init(ocpi.LF())
}

func (s *cdrUcTestSuite) SetupTest() {
}

func (s *cdrUcTestSuite) TearDownSuite() {}

func TestCdrUcSuite(t *testing.T) {
	suite.Run(t, new(cdrUcTestSuite))
}
