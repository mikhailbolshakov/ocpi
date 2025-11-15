package impl

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/errors"
	"github.com/mikhailbolshakov/ocpi/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type webhookTestSuite struct {
	kit.Suite
	svc     *webhookImpl
	storage *mocks.WebhookStorage
}

func (s *webhookTestSuite) SetupSuite() {
	s.Suite.Init(ocpi.LF())
}

func (s *webhookTestSuite) SetupTest() {
	s.Suite.Init(ocpi.LF())
	s.storage = &mocks.WebhookStorage{}
	s.svc = NewWebhookService(s.storage).(*webhookImpl)
}

func (s *webhookTestSuite) TearDownSuite() {}

func TestWebhookServiceSuite(t *testing.T) {
	suite.Run(t, new(webhookTestSuite))
}

func (s *webhookTestSuite) webhook() *backend.Webhook {
	return &backend.Webhook{
		Id:     kit.NewId(),
		ApiKey: kit.NewRandString(),
		Events: []string{"ev1", "ev2"},
		Url:    "https://webhook.com/wh",
	}
}

func (s *webhookTestSuite) Test_Validate() {
	// id empty
	wh := s.webhook()
	wh.Id = ""
	s.AssertAppErr(s.svc.validate(s.Ctx, wh), errors.ErrCodeWhIdEmpty)

	// apikey empty
	wh = s.webhook()
	wh.ApiKey = ""
	s.AssertAppErr(s.svc.validate(s.Ctx, wh), errors.ErrCodeWhApiKeyEmpty)

	// url empty
	wh = s.webhook()
	wh.Url = ""
	s.AssertAppErr(s.svc.validate(s.Ctx, wh), errors.ErrCodeWhUrlInvalid)

	// url not valid
	wh = s.webhook()
	wh.Url = "wrywuiry"
	s.AssertAppErr(s.svc.validate(s.Ctx, wh), errors.ErrCodeWhUrlInvalid)

	// no events
	wh = s.webhook()
	wh.Events = nil
	s.AssertAppErr(s.svc.validate(s.Ctx, wh), errors.ErrCodeWhEventsEmpty)

	// valid
	s.NoError(s.svc.validate(s.Ctx, s.webhook()))

}

func (s *webhookTestSuite) Test_CreateUpdate() {
	wh := s.webhook()
	s.storage.On("MergeWebhook", s.Ctx, wh).Return(nil)
	act, err := s.svc.CreateUpdate(s.Ctx, wh)
	s.NoError(err)
	s.NotEmpty(act)
	s.AssertCalled(&s.storage.Mock, "MergeWebhook", s.Ctx, wh)
}

func (s *webhookTestSuite) Test_Create_WhenExists_Fail() {
	wh := s.webhook()
	s.storage.On("GetWebhook", s.Ctx, wh.Id).Return(wh, nil)
	_, err := s.svc.Create(s.Ctx, wh)
	s.AssertAppErr(err, errors.ErrCodeWhAlreadyExists)
}

func (s *webhookTestSuite) Test_Create_WhenNotExists_Ok() {
	wh := s.webhook()
	s.storage.On("GetWebhook", s.Ctx, wh.Id).Return(nil, nil)
	s.storage.On("CreateWebhook", s.Ctx, wh).Return(nil)
	act, err := s.svc.Create(s.Ctx, wh)
	s.NoError(err)
	s.NotEmpty(act)
	s.AssertCalled(&s.storage.Mock, "CreateWebhook", s.Ctx, wh)
}

func (s *webhookTestSuite) Test_Update_WhenNotExists_Fail() {
	wh := s.webhook()
	s.storage.On("GetWebhook", s.Ctx, wh.Id).Return(nil, nil)
	_, err := s.svc.Update(s.Ctx, wh)
	s.AssertAppErr(err, errors.ErrCodeWhNotFound)
}

func (s *webhookTestSuite) Test_Update_Ok() {
	wh := s.webhook()
	s.storage.On("GetWebhook", s.Ctx, wh.Id).Return(wh, nil)
	s.storage.On("UpdateWebhook", s.Ctx, wh).Return(nil)
	act, err := s.svc.Update(s.Ctx, wh)
	s.NoError(err)
	s.NotEmpty(act)
	s.AssertCalled(&s.storage.Mock, "UpdateWebhook", s.Ctx, wh)
}

func (s *webhookTestSuite) Test_Delete_WhenNotExists_Fail() {
	wh := s.webhook()
	s.storage.On("GetWebhook", s.Ctx, wh.Id).Return(nil, nil)
	s.AssertAppErr(s.svc.Delete(s.Ctx, wh.Id), errors.ErrCodeWhNotFound)
}

func (s *webhookTestSuite) Test_Delete_Ok() {
	wh := s.webhook()
	s.storage.On("GetWebhook", s.Ctx, wh.Id).Return(wh, nil)
	s.storage.On("DeleteWebhook", s.Ctx, wh.Id).Return(nil)
	s.NoError(s.svc.Delete(s.Ctx, wh.Id))
	s.AssertCalled(&s.storage.Mock, "DeleteWebhook", s.Ctx, wh.Id)
}
