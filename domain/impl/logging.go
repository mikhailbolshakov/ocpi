package impl

import (
	"context"
	"github.com/mikhailbolshakov/ocpi/domain"
)

type ocpiLogImpl struct {
	storage domain.OcpiLogStorage
}

func NewOcpiLogService(storage domain.OcpiLogStorage) domain.OcpiLogService {
	return &ocpiLogImpl{
		storage: storage,
	}
}

func (l *ocpiLogImpl) Init(severity string) {
}

func (l *ocpiLogImpl) Log(ctx context.Context, msg *domain.LogMessage) {
	l.storage.Save(ctx, msg)
}

func (l *ocpiLogImpl) Search(ctx context.Context, criteria *domain.SearchLogCriteria) ([]*domain.LogMessage, error) {
	if criteria.Size <= 0 {
		criteria.Size = 100
	}
	return l.storage.SearchLog(ctx, criteria)
}
