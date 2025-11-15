//go:build integration

package storage

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/stretchr/testify/suite"
	"testing"
)

type whStorageTestSuite struct {
	kit.Suite
	storage backend.WebhookStorage
	adapter Adapter
}

func (s *whStorageTestSuite) SetupSuite() {
	s.Suite.Init(ocpi.LF())

	// load config
	cfg, err := ocpi.LoadConfig()
	if err != nil {
		s.Fatal(err)
	}

	s.adapter = NewAdapter()
	s.NoError(s.adapter.Init(s.Ctx, cfg.Storages))

	s.storage = s.adapter
}

func (s *whStorageTestSuite) TearDownSuite() {
	_ = s.adapter.Close(s.Ctx)
}

func TestWebhookStorageSuite(t *testing.T) {
	suite.Run(t, new(whStorageTestSuite))
}

func (s *whStorageTestSuite) webhook() *backend.Webhook {
	return &backend.Webhook{
		Id:     kit.NewId(),
		ApiKey: kit.NewRandString(),
		Events: []string{kit.NewRandString(), kit.NewRandString()},
		Url:    "https://webhook.test",
	}
}

func (s *whStorageTestSuite) Test_CRUD() {
	// get when no exists
	act, err := s.storage.GetWebhook(s.Ctx, kit.NewId())
	s.NoError(err)
	s.Empty(act)

	// create
	wh := s.webhook()
	s.NoError(s.storage.CreateWebhook(s.Ctx, wh))

	// get
	act, err = s.storage.GetWebhook(s.Ctx, wh.Id)
	s.NoError(err)
	s.Equal(act, wh)

	// update
	wh.Events = append(wh.Events, kit.NewRandString())
	wh.ApiKey = kit.NewRandString()
	s.NoError(s.storage.UpdateWebhook(s.Ctx, wh))

	// get
	act, err = s.storage.GetWebhook(s.Ctx, wh.Id)
	s.NoError(err)
	s.Equal(act, wh)

	// delete
	s.NoError(s.storage.DeleteWebhook(s.Ctx, wh.Id))

	// get
	act, err = s.storage.GetWebhook(s.Ctx, wh.Id)
	s.NoError(err)
	s.Empty(act)

}

func (s *whStorageTestSuite) Test_Merge() {

	// create
	wh := s.webhook()
	s.NoError(s.storage.MergeWebhook(s.Ctx, wh))

	// get
	act, err := s.storage.GetWebhook(s.Ctx, wh.Id)
	s.NoError(err)
	s.Equal(act, wh)

	// update
	wh.Events = append(wh.Events, kit.NewRandString())
	wh.ApiKey = kit.NewRandString()
	s.NoError(s.storage.MergeWebhook(s.Ctx, wh))

	// get
	act, err = s.storage.GetWebhook(s.Ctx, wh.Id)
	s.NoError(err)
	s.Equal(act, wh)

	// delete
	s.NoError(s.storage.DeleteWebhook(s.Ctx, wh.Id))

}

func (s *whStorageTestSuite) Test_Search() {

	// create
	wh1 := s.webhook()
	s.NoError(s.storage.MergeWebhook(s.Ctx, wh1))
	wh2 := s.webhook()
	wh2.Events = append(wh2.Events, wh1.Events...)
	s.NoError(s.storage.MergeWebhook(s.Ctx, wh2))

	// search
	act, err := s.storage.SearchWebhook(s.Ctx, &backend.SearchWebhookCriteria{Event: wh1.Events[0]})
	s.NoError(err)
	s.Len(act, 2)

	// search
	act, err = s.storage.SearchWebhook(s.Ctx, &backend.SearchWebhookCriteria{Event: wh2.Events[0]})
	s.NoError(err)
	s.Len(act, 1)

	// search
	act, err = s.storage.SearchWebhook(s.Ctx, &backend.SearchWebhookCriteria{Event: kit.NewRandString()})
	s.NoError(err)
	s.Empty(act)

	// delete
	s.NoError(s.storage.DeleteWebhook(s.Ctx, wh1.Id))
	s.NoError(s.storage.DeleteWebhook(s.Ctx, wh2.Id))

}
