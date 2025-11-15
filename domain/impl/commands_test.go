package impl

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/stretchr/testify/suite"
	"testing"
)

type commandTestSuite struct {
	kit.Suite
}

func (s *commandTestSuite) SetupSuite() {
	s.Suite.Init(ocpi.LF())
}

func (s *commandTestSuite) SetupTest() {
}

func (s *commandTestSuite) TearDownSuite() {}

func TestCommandSuite(t *testing.T) {
	suite.Run(t, new(commandTestSuite))
}
