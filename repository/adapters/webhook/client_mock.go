package webhook

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	service "github.com/mikhailbolshakov/ocpi"
)

type mockClientImpl struct{}

func newMockWhRestClient() webhookRestClient {
	return &mockClientImpl{}
}

func (s *mockClientImpl) l() kit.CLogger {
	return service.L().Cmp("mock-wh-rest")
}

func (s *mockClientImpl) Init(ctx context.Context, config *service.CfgWebHook) error {
	s.l().Mth("init").Dbg()
	return nil
}

func (s *mockClientImpl) Close(ctx context.Context) error {
	s.l().Mth("close").Dbg()
	return nil
}

func (s *mockClientImpl) CallAsync(ctx context.Context, url, apiKey, event string, payload any) {
	s.l().Mth("call-async").F(kit.KV{"url": url, "event": event}).Dbg("ok")
}
