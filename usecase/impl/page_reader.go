package impl

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/goroutine"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/model"
	"github.com/mikhailbolshakov/ocpi/usecase"
	"time"
)

const (
	maxNumberOfPages = 100
)

type pageReader[R usecase.PageResponse] struct {
	pageFn func(context.Context, *usecase.OcpiRepositoryPagingRequest) (R, error)
}

func NewPageReader[R usecase.PageResponse](fn func(context.Context, *usecase.OcpiRepositoryPagingRequest) (R, error)) usecase.PageReader[R] {
	return &pageReader[R]{
		pageFn: fn,
	}
}

func (p *pageReader[R]) l() kit.CLogger {
	return ocpi.L().Cmp("page-reader")
}

func (p *pageReader[R]) GetPage(ctx context.Context, rq *usecase.OcpiRepositoryBaseRequest, pageSize int, dateFrom, dateTo *time.Time) chan R {
	l := p.l().C(ctx).Mth("get-page").Dbg()

	res := make(chan R, 10)

	goroutine.New().WithLogger(l).Go(ctx, func() {
		offset := 0
		pageIndex := 0

		defer close(res)

		for {
			select {

			case <-ctx.Done():
				return

			default:

				rs, err := p.pageFn(ctx, &usecase.OcpiRepositoryPagingRequest{
					OcpiRepositoryBaseRequest: *rq,
					OcpiGetPageRequest: model.OcpiGetPageRequest{
						DateFrom: dateFrom,
						DateTo:   dateTo,
						Offset:   &offset,
						Limit:    &pageSize,
					},
				})
				if err != nil {
					l.E(err).St().Err()
					return
				}

				if len(rs) == 0 {
					return
				}

				// we have to ensure we won't stick if a remote platform doesn't implement paging correctly
				if pageIndex > maxNumberOfPages {
					l.Dbg("max number of pages reached")
					return
				}

				offset += pageSize
				pageIndex++

				res <- rs
			}
		}
	},
	)
	return res
}
