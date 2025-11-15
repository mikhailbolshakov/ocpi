package impl

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/stretchr/testify/suite"
	"testing"
)

type tokenTestSuite struct {
	kit.Suite
}

func (s *tokenTestSuite) SetupSuite() {
	s.Suite.Init(ocpi.LF())
}

func (s *tokenTestSuite) SetupTest() {
}

func (s *tokenTestSuite) TearDownSuite() {}

func TestTokenSuite(t *testing.T) {
	suite.Run(t, new(tokenTestSuite))
}
