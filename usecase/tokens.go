package usecase

import (
	"context"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/model"
	"time"
)

type TokenConverter interface {
	// TokenDomainToModel converts token domain to ocpi model
	TokenDomainToModel(tkn *domain.Token) *model.OcpiToken
	// TokensDomainToModel converts tokens domain to ocpi model
	TokensDomainToModel(ts []*domain.Token) []*model.OcpiToken
	// TokenModelToDomain converts token model to domain
	TokenModelToDomain(tkn *model.OcpiToken, platformId string) *domain.Token
	// TokenDomainToBackend converts token domain to backend
	TokenDomainToBackend(tkn *domain.Token) *backend.Token
	// TokensDomainToBackend converts token domain to backend
	TokensDomainToBackend(ts []*domain.Token) []*backend.Token
	// TokenBackendToDomain converts token backend to domain
	TokenBackendToDomain(tkn *backend.Token, platformId string) *domain.Token
}

type TokenUc interface {
	// OnLocalTokenChanged handles changing token in local platform
	OnLocalTokenChanged(ctx context.Context, tkn *domain.Token) error
	// OnRemoteTokensPull handles request to pull tokens from remote platforms (fired by cron)
	OnRemoteTokensPull(ctx context.Context, from, to *time.Time) error
	// OnRemoteTokensPullWhenPushNotSupported handles request to pull tokens from remote platforms which don't support push (fired by cron)
	OnRemoteTokensPullWhenPushNotSupported(ctx context.Context, from, to *time.Time) error
	// OnRemoteTokenPut handles put token in remote platform
	OnRemoteTokenPut(ctx context.Context, platformId string, tkn *model.OcpiToken) error
	// OnRemoteTokenPatch handles patch token in remote platform
	OnRemoteTokenPatch(ctx context.Context, platformId string, tkn *model.OcpiToken) error

	// GetOrCreateLocalToken first tries to find token by id, is not exists, create a new one local token with default attr
	GetOrCreateLocalToken(ctx context.Context, tkn *domain.Token) (*domain.Token, error)
}

type RemoteTokenRepository interface {
	// PutTokenAsync puts token
	PutTokenAsync(ctx context.Context, rq *OcpiRepositoryErrHandlerRequestG[*model.OcpiToken])
	// PatchTokenAsync patches token
	PatchTokenAsync(ctx context.Context, rq *OcpiRepositoryErrHandlerRequestG[*model.OcpiToken])
	// GetTokens retrieves tokens
	GetTokens(ctx context.Context, rq *OcpiRepositoryPagingRequest) ([]*model.OcpiToken, error)
	// GetToken retrieves token by id
	GetToken(ctx context.Context, rq *OcpiRepositoryIdRequest) (*model.OcpiToken, error)
}
