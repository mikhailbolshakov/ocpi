package impl

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/stretchr/testify/suite"
	"testing"
)

type tokenUcTestSuite struct {
	kit.Suite
}

func (s *tokenUcTestSuite) SetupSuite() {
	s.Suite.Init(ocpi.LF())
}

func (s *tokenUcTestSuite) SetupTest() {
}

func (s *tokenUcTestSuite) TearDownSuite() {}

func TestTokenUcSuite(t *testing.T) {
	suite.Run(t, new(tokenUcTestSuite))
}
