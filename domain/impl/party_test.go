package impl

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/stretchr/testify/suite"
	"testing"
)

type partyTestSuite struct {
	kit.Suite
}

func (s *partyTestSuite) SetupSuite() {
	s.Suite.Init(ocpi.LF())
}

func (s *partyTestSuite) SetupTest() {
}

func (s *partyTestSuite) TearDownSuite() {}

func TestPartySuite(t *testing.T) {
	suite.Run(t, new(partyTestSuite))
}
