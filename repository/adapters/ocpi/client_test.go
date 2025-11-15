//go:build debug

package ocpi

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/mocks"
	"github.com/mikhailbolshakov/ocpi/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type clientTestSuite struct {
	kit.Suite
	clSvc  *clientImpl
	logSvc *mocks.OcpiLogService
}

func (s *clientTestSuite) SetupSuite() {
	s.Suite.Init(ocpi.LF())
}

func (s *clientTestSuite) SetupTest() {
	s.logSvc = &mocks.OcpiLogService{}
	s.logSvc.On("Log", mock.Anything, mock.Anything)
	s.clSvc = newOcpiRestClient(s.logSvc).(*clientImpl)
	s.clSvc.Init(s.Ctx, &ocpi.CfgOcpiRemote{
		Mock:    false,
		Timeout: kit.IntPtr(20),
	})
}

func (s *clientTestSuite) TearDownSuite() {}

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(clientTestSuite))
}
