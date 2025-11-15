package usecase

import (
	"context"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/model"
	"time"
)

type SessionConverter interface {
	// SessionDomainToModel converts session domain to ocpi model
	SessionDomainToModel(sess *domain.Session) *model.OcpiSession
	// SessionsDomainToModel converts sessions domain to ocpi model
	SessionsDomainToModel(ts []*domain.Session) []*model.OcpiSession
	// SessionModelToDomain converts session model to domain
	SessionModelToDomain(sess *model.OcpiSession, platformId string) *domain.Session
	// SessionDomainToBackend converts session domain to backend
	SessionDomainToBackend(sess *domain.Session) *backend.Session
	// SessionsDomainToBackend converts session domain to backend
	SessionsDomainToBackend(ts []*domain.Session) []*backend.Session
	// SessionBackendToDomain converts session backend to domain
	SessionBackendToDomain(sess *backend.Session, platformId string) *domain.Session
	// TokenToCdrTokenDomain converts token object to cdr token
	TokenToCdrTokenDomain(tkn *domain.Token) *domain.CdrToken
}

type SessionUc interface {
	// OnLocalSessionChanged handles changing session in local platform
	OnLocalSessionChanged(ctx context.Context, sess *domain.Session) error
	// OnLocalSessionPatched handles patching session
	OnLocalSessionPatched(ctx context.Context, sess *domain.Session) error
	// OnRemoteSessionsPull handles request to pull sessions from remote platforms (fired by cron)
	OnRemoteSessionsPull(ctx context.Context, from, to *time.Time) error
	// OnRemoteSessionsPullWhenPushNotSupported handles request to pull sessions from remote platforms which don't support push (fired by cron)
	OnRemoteSessionsPullWhenPushNotSupported(ctx context.Context, from, to *time.Time) error
	// OnRemoteSessionPut handles put session in remote platform
	OnRemoteSessionPut(ctx context.Context, platformId string, sess *model.OcpiSession) error
	// OnRemoteSessionPatch handles patch session in remote platform
	OnRemoteSessionPatch(ctx context.Context, platformId string, sess *model.OcpiSession) error
}

type RemoteSessionRepository interface {
	// PutSessionAsync puts session
	PutSessionAsync(ctx context.Context, rq *OcpiRepositoryErrHandlerRequestG[*model.OcpiSession])
	// PatchSessionAsync patches session
	PatchSessionAsync(ctx context.Context, rq *OcpiRepositoryErrHandlerRequestG[*model.OcpiSession])
	// GetSessions retrieves sessions
	GetSessions(ctx context.Context, rq *OcpiRepositoryPagingRequest) ([]*model.OcpiSession, error)
	// GetSession retrieves session by id
	GetSession(ctx context.Context, rq *OcpiRepositoryIdRequest) (*model.OcpiSession, error)
}
