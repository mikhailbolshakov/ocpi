package storage

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type totalCount struct {
	TotalCount int `gorm:"column:total_count"`
}

func pagingLimit(rqLimit *int) *int {
	if rqLimit == nil {
		return kit.IntPtr(domain.PageSizeDefault)
	}
	if *rqLimit > domain.PageSizeMaxLimit {
		return kit.IntPtr(domain.PageSizeMaxLimit)
	}
	return rqLimit
}

func paging(rq domain.PageRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if rq.Offset == nil {
			rq.Offset = kit.IntPtr(0)
		}
		return db.Limit(*pagingLimit(rq.Limit)).Offset(*rq.Offset)
	}
}

func orderByLastUpdated(desc bool) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(clause.OrderByColumn{Column: clause.Column{Name: "last_updated"}, Desc: desc})
	}
}

func merge() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Clauses(clause.OnConflict{UpdateAll: true})
	}
}

func update() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Omit("created_at")
	}
}

func nextPage(rq domain.PageRequest, total *int) *domain.PageRequest {

	if total == nil {
		return nil
	}
	if rq.Limit == nil {
		rq.Limit = pagingLimit(rq.Limit)
	}
	if rq.Offset == nil {
		rq.Offset = kit.IntPtr(0)
	}

	// last page
	if (*rq.Limit)+(*rq.Offset) >= *total {
		return nil
	}

	return &domain.PageRequest{
		Offset:   kit.IntPtr(*rq.Offset + *rq.Limit),
		Limit:    rq.Limit,
		DateFrom: rq.DateFrom,
		DateTo:   rq.DateTo,
	}

}
