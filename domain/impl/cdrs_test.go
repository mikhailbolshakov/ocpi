package impl

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/stretchr/testify/suite"
	"testing"
)

type cdrTestSuite struct {
	kit.Suite
}

func (s *cdrTestSuite) SetupSuite() {
	s.Suite.Init(ocpi.LF())
}

func (s *cdrTestSuite) SetupTest() {
}

func (s *cdrTestSuite) TearDownSuite() {}

func TestCdrSuite(t *testing.T) {
	suite.Run(t, new(cdrTestSuite))
}
