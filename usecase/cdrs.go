package usecase

import (
	"context"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/model"
	"time"
)

type CdrConverter interface {
	// CdrDomainToModel converts cdr domain to ocpi model
	CdrDomainToModel(cdr *domain.Cdr) *model.OcpiCdr
	// CdrsDomainToModel converts cdrs domain to ocpi model
	CdrsDomainToModel(ts []*domain.Cdr) []*model.OcpiCdr
	// CdrModelToDomain converts cdr model to domain
	CdrModelToDomain(cdr *model.OcpiCdr, platformId string) *domain.Cdr
	// CdrDomainToBackend converts cdr domain to backend
	CdrDomainToBackend(cdr *domain.Cdr) *backend.Cdr
	// CdrsDomainToBackend converts cdr domain to backend
	CdrsDomainToBackend(ts []*domain.Cdr) []*backend.Cdr
	// CdrBackendToDomain converts cdr backend to domain
	CdrBackendToDomain(cdr *backend.Cdr, sess *domain.Session, loc *domain.Location, evse *domain.Evse, con *domain.Connector) *domain.Cdr
}

type CdrUc interface {
	// OnLocalCdrChanged handles changing cdr in local platform
	OnLocalCdrChanged(ctx context.Context, cdr *backend.Cdr) error
	// OnRemoteCdrsPull handles request to pull cdrs from remote platforms (fired by cron)
	OnRemoteCdrsPull(ctx context.Context, from, to *time.Time) error
	// OnRemoteCdrsPullWhenPushNotSupported handles request to pull cdrs from remote platforms which don't support push (fired by cron)
	OnRemoteCdrsPullWhenPushNotSupported(ctx context.Context, from, to *time.Time) error
	// OnRemoteCdrPut handles put cdr in remote platform
	OnRemoteCdrPut(ctx context.Context, platformId string, cdr *model.OcpiCdr) error
}

type RemoteCdrRepository interface {
	// PostCdrAsync posts cdr
	PostCdrAsync(ctx context.Context, rq *OcpiRepositoryErrHandlerRequestG[*model.OcpiCdr])
	// GetCdrs retrieves cdrs
	GetCdrs(ctx context.Context, rq *OcpiRepositoryPagingRequest) ([]*model.OcpiCdr, error)
	// GetCdr retrieves cdr by id
	GetCdr(ctx context.Context, rq *OcpiRepositoryIdRequest) (*model.OcpiCdr, error)
}
