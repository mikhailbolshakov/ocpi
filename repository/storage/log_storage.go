package storage

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/goroutine"
	"github.com/mikhailbolshakov/kit/storages/pg"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/errors"
	"gorm.io/gorm"
	"time"
)

const (
	flushInterval = time.Second * 3
	flushMaxItems = 100
)

type logDetails struct {
	Rq     any `json:"rq,omitempty"`
	Rs     any `json:"rs,omitempty"`
	Header any `json:"hdr,omitempty"`
}

type log struct {
	Event          string  `gorm:"column:event"`
	Url            *string `gorm:"column:url"`
	Token          *string `gorm:"column:token"`
	RequestId      *string `gorm:"column:rq_id"`
	CorrelationId  *string `gorm:"column:corr_id"`
	FromPlatformId *string `gorm:"column:from_platform_id"`
	ToPlatformId   *string `gorm:"column:to_platform_id"`
	Err            *string `gorm:"column:err"`
	Status         int     `gorm:"column:status"`
	OcpiStatus     int     `gorm:"column:ocpi_status"`
	Details        string  `gorm:"column:details"`
	DurationMs     int64   `gorm:"column:dur_ms"`
	Incoming       bool    `gorm:"column:incoming"`
	CreatedAt      *time.Time
}

type logStorageImpl struct {
	pg            *pg.Storage
	eventChan     chan *log
	cancelCtx     context.Context
	cancelFn      context.CancelFunc
	flashInterval time.Duration
	flashMaxItems int
}

func (s *logStorageImpl) l() kit.CLogger {
	return ocpi.L().Cmp("log-storage")
}

func newLogStorage(pg *pg.Storage) *logStorageImpl {
	return &logStorageImpl{
		pg:            pg,
		flashInterval: flushInterval,
		flashMaxItems: flushMaxItems,
	}
}

func (s *logStorageImpl) Save(ctx context.Context, msg *domain.LogMessage) {
	s.l().C(ctx).Mth("save").Dbg()
	s.eventChan <- s.toDto(msg)
}

func (s *logStorageImpl) SearchLog(ctx context.Context, criteria *domain.SearchLogCriteria) ([]*domain.LogMessage, error) {
	s.l().C(ctx).Mth("search").Dbg()

	// make query
	var dtos []*log

	criteria.PagingRequest.SortBy = []*kit.SortRequest{{Field: "created_at"}}
	if criteria.PagingRequest.Size <= 0 {
		criteria.PagingRequest.Size = 20
	}

	if err := s.pg.Instance.
		Scopes(s.buildSearchQuery(criteria), pg.Paging(criteria.PagingRequest)).
		Find(&dtos).Error; err != nil {
		return nil, errors.ErrSessStorageGet(ctx, err)
	}

	if len(dtos) == 0 {
		return nil, nil
	}

	return kit.Select(dtos, s.toDomain), nil
}

func (s *logStorageImpl) init(ctx context.Context, flashInterval time.Duration, flashMaxItems int) error {
	s.l().C(ctx).Mth("init").Dbg()

	// setup config
	if flashInterval > 0 {
		s.flashInterval = flashInterval
	}
	if flashMaxItems > 0 {
		s.flashMaxItems = flashMaxItems
	}

	// cancel running forcibly
	if s.cancelFn != nil {
		s.cancelFn()
	}

	// init cancellation context
	s.eventChan = make(chan *log, 999)
	s.cancelCtx, s.cancelFn = context.WithCancel(ctx)

	// run db writer in a separate goroutine
	goroutine.New().WithLogger(s.l().C(ctx).Mth("writer")).WithRetry(goroutine.Unrestricted).Go(ctx, func() { s.writer(ctx) })
	return nil
}

func (s *logStorageImpl) writer(ctx context.Context) {
	for keepGoing := true; keepGoing; {
		var eventsBatch []*log
		expire := time.After(s.flashInterval)
		for {
			select {
			case ev, ok := <-s.eventChan:
				if !ok {
					keepGoing = false
					goto flush
				}
				eventsBatch = append(eventsBatch, ev)
				if len(eventsBatch) == s.flashMaxItems {
					goto flush
				}
			// flush when timeout expires
			case <-expire:
				goto flush
			// leave when context cancelled
			case <-ctx.Done():
				keepGoing = false
				goto flush
			}
		}
	flush:
		if len(eventsBatch) > 0 {
			if err := s.createLogs(ctx, eventsBatch); err != nil {
				s.l().C(ctx).Mth("writer").E(err).Err()
			}
		}
	}
}

func (s *logStorageImpl) createLogs(ctx context.Context, dtos []*log) error {
	s.l().C(ctx).Mth("create").Dbg()
	if err := s.pg.Instance.Create(dtos).Error; err != nil {
		return errors.ErrLogStorageCreate(ctx, err)
	}
	return nil
}

func (s *logStorageImpl) close(ctx context.Context) {
	if s.cancelFn != nil {
		s.cancelFn()
	}
	close(s.eventChan)
}

func (s *logStorageImpl) buildSearchQuery(criteria *domain.SearchLogCriteria) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query := db.Model(&log{})
		// populate conditions
		if len(criteria.Events) > 0 {
			query = query.Where("event in (?)", criteria.Events)
		}
		if criteria.DateTo != nil {
			query = query.Where("created_at <= ?", *criteria.DateTo)
		}
		if criteria.DateFrom != nil {
			query = query.Where("created_at >= ?", *criteria.DateFrom)
		}
		if criteria.Incoming != nil {
			query = query.Where("incoming = ?", *criteria.Incoming)
		}
		if criteria.RequestId != "" {
			query = query.Where("rq_id = ?", criteria.RequestId)
		}
		if criteria.ToPlatform != "" {
			query = query.Where("to_platform_id = ?", criteria.ToPlatform)
		}
		if criteria.FromPlatform != "" {
			query = query.Where("from_platform_id = ?", criteria.FromPlatform)
		}
		if criteria.OcpiStatus != nil {
			query = query.Where("ocpi_status = ?", *criteria.OcpiStatus)
		}
		if criteria.HttpStatus != nil {
			query = query.Where("status = ?", *criteria.HttpStatus)
		}
		if criteria.Error != nil {
			if *criteria.Error {
				query = query.Where("err is not null")
			} else {
				query = query.Where("err is null")
			}
		}
		return query
	}
}
