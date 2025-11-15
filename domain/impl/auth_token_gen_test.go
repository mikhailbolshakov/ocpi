package impl

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/stretchr/testify/suite"
	"testing"
)

type authTokenGenTestSuite struct {
	kit.Suite
	gen domain.TokenGenerator
}

func (s *authTokenGenTestSuite) SetupSuite() {
	s.Suite.Init(ocpi.LF())
}

func (s *authTokenGenTestSuite) SetupTest() {
	s.gen = NewTokenGenerator()
}

func (s *authTokenGenTestSuite) TearDownSuite() {}

func TestAuthTokenGenSuite(t *testing.T) {
	suite.Run(t, new(authTokenGenTestSuite))
}

func (s *authTokenGenTestSuite) Test_Base64() {
	v, err := s.gen.Generate(s.Ctx)
	s.NoError(err)
	s.NotEmpty(v)

	encoded := s.gen.Base64Encode(v)
	s.NotEmpty(encoded)

	decoded, ok := s.gen.TryBase64Decode(encoded)
	s.Equal(decoded, v)
	s.True(ok)

}
