package impl

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/stretchr/testify/suite"
	"testing"
)

type platformTestSuite struct {
	kit.Suite
}

func (s *platformTestSuite) SetupSuite() {
	s.Suite.Init(ocpi.LF())
}

func (s *platformTestSuite) SetupTest() {
}

func (s *platformTestSuite) TearDownSuite() {}

func TestPlatformSuite(t *testing.T) {
	suite.Run(t, new(platformTestSuite))
}
