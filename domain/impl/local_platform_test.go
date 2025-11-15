package impl

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/stretchr/testify/suite"
	"testing"
)

type localPlatformTestSuite struct {
	kit.Suite
}

func (s *localPlatformTestSuite) SetupSuite() {
	s.Suite.Init(ocpi.LF())
}

func (s *localPlatformTestSuite) SetupTest() {
}

func (s *localPlatformTestSuite) TearDownSuite() {}

func TestLocalPlatformSuite(t *testing.T) {
	suite.Run(t, new(localPlatformTestSuite))
}
