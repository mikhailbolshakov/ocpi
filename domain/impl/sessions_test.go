package impl

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/stretchr/testify/suite"
	"testing"
)

type sessionTestSuite struct {
	kit.Suite
}

func (s *sessionTestSuite) SetupSuite() {
	s.Suite.Init(ocpi.LF())
}

func (s *sessionTestSuite) SetupTest() {
}

func (s *sessionTestSuite) TearDownSuite() {}

func TestSessionSuite(t *testing.T) {
	suite.Run(t, new(sessionTestSuite))
}
