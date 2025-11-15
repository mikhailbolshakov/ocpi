package usecase

import (
	"context"
	"github.com/mikhailbolshakov/ocpi/model"
	"time"
)

// PageReader is a generic reader allowing read from remote repository
type PageReader[R PageResponse] interface {
	// GetPage retrieves one page
	GetPage(ctx context.Context, rq *OcpiRepositoryBaseRequest, pageSize int, dateFrom, dateTo *time.Time) chan R
}

type PageResponse interface {
	[]*model.OcpiLocation | []*model.OcpiTariff | []*model.OcpiToken | []*model.OcpiSession | []*model.OcpiClientInfo | []*model.OcpiCdr
}
