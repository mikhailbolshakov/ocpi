//go:build integration

package storage

import (
	"fmt"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/storages/pg"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type loggerTestSuite struct {
	kit.Suite
	pg *pg.Storage
}

func (s *loggerTestSuite) SetupSuite() {
	s.Suite.Init(ocpi.LF())

	// load config
	cfg, err := ocpi.LoadConfig()
	if err != nil {
		s.Fatal(err)
	}

	s.pg, err = pg.Open(cfg.Storages.Database.Master, ocpi.LF())
	if err != nil {
		s.Fatal(err)
	}
}

func (s *loggerTestSuite) TearDownSuite() {
	s.pg.Close()
}

func TestLoggerSuite(t *testing.T) {
	suite.Run(t, new(loggerTestSuite))
}

func (s *loggerTestSuite) Test_WhenFlushByTimeout() {
	dbLogger := newLogStorage(s.pg)
	err := dbLogger.init(s.Ctx, time.Millisecond*200, 10)
	if err != nil {
		s.Fatal(err)
	}
	defer dbLogger.close(s.Ctx)
	// put5 entries
	rqId := kit.NewRandString()
	for i := 0; i < 5; i++ {
		msg := s.logMsg()
		msg.RequestId = rqId
		dbLogger.Save(s.Ctx, msg)
	}
	// await
	var res []*log
	if err := <-kit.Await(func() (bool, error) {
		if err := s.pg.Instance.Where("rq_id = ?", rqId).Find(&res).Error; err != nil {
			return false, err
		}
		return len(res) == 5, nil
	}, time.Millisecond*300, time.Second*2); err != nil {
		s.Fatal(err)
	}
	s.Len(res, 5)
}

func (s *loggerTestSuite) Test_WhenFlushByNumberOfItems() {
	dbLogger := newLogStorage(s.pg)
	err := dbLogger.init(s.Ctx, time.Second*5, 5)
	if err != nil {
		s.Fatal(err)
	}
	defer dbLogger.close(s.Ctx)
	rqId := kit.NewRandString()
	for i := 0; i < 5; i++ {
		msg := s.logMsg()
		msg.RequestId = rqId
		dbLogger.Save(s.Ctx, msg)
	}
	// await
	var res []*log
	if err := <-kit.Await(func() (bool, error) {
		if err := s.pg.Instance.Where("rq_id = ?", rqId).Find(&res).Error; err != nil {
			return false, err
		}
		return len(res) == 5, nil
	}, time.Millisecond*300, time.Second*2); err != nil {
		s.Fatal(err)
	}
	s.Len(res, 5)
}

func (s *loggerTestSuite) Test_Search() {
	dbLogger := newLogStorage(s.pg)

	msg1 := s.logMsg()
	msg1.In = true
	msg1.OcpiStatus = 2000
	msg1.ResponseStatus = 400
	msg2 := s.logMsg()
	msg2.Err = nil
	s.NoError(dbLogger.createLogs(s.Ctx, []*log{dbLogger.toDto(msg1), dbLogger.toDto(msg2)}))

	r, err := dbLogger.SearchLog(s.Ctx, &domain.SearchLogCriteria{RequestId: msg1.RequestId, FromPlatform: msg1.FromPlatform})
	s.NoError(err)
	s.Len(r, 1)

	r, err = dbLogger.SearchLog(s.Ctx, &domain.SearchLogCriteria{RequestId: msg1.RequestId, OcpiStatus: &msg1.OcpiStatus, Incoming: kit.BoolPtr(true)})
	s.NoError(err)
	s.Len(r, 1)

	r, err = dbLogger.SearchLog(s.Ctx, &domain.SearchLogCriteria{RequestId: msg1.RequestId, ToPlatform: msg1.ToPlatform, Error: kit.BoolPtr(true)})
	s.NoError(err)
	s.Len(r, 1)

	r, err = dbLogger.SearchLog(s.Ctx, &domain.SearchLogCriteria{RequestId: msg2.RequestId, Incoming: kit.BoolPtr(false), Error: kit.BoolPtr(false)})
	s.NoError(err)
	s.Len(r, 1)

	r, err = dbLogger.SearchLog(s.Ctx, &domain.SearchLogCriteria{RequestId: msg2.RequestId, DateFrom: kit.TimePtr(kit.Now().Add(-time.Minute)), DateTo: kit.TimePtr(kit.Now().Add(time.Minute))})
	s.NoError(err)
	s.Len(r, 1)

}

func (s *loggerTestSuite) logMsg() *domain.LogMessage {
	return &domain.LogMessage{
		Event:         domain.LogEventPutLocation,
		Url:           "https://ocpi.test/locations",
		Token:         kit.NewRandString(),
		RequestId:     kit.NewRandString(),
		CorrelationId: kit.NewRandString(),
		FromPlatform:  "platform-1",
		ToPlatform:    "platform-2",
		RequestBody: struct {
			Val string `json:"val"`
		}{
			Val: "value",
		},
		ResponseBody: struct {
			Val string `json:"val"`
		}{
			Val: "value",
		},
		Headers: map[string]string{
			"key": "val",
		},
		ResponseStatus: 200,
		OcpiStatus:     1000,
		Err:            fmt.Errorf("error"),
		In:             false,
		DurationMs:     0,
	}
}
