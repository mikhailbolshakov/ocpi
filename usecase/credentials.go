package usecase

import (
	"context"
	_ "github.com/eapache/channels"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/model"
)

type CredentialsConverter interface {
	// PartyDomainToBackend converts party domain to backend
	PartyDomainToBackend(party *domain.Party) *backend.Party
	// PartiesDomainToBackend converts parties domain to backend
	PartiesDomainToBackend(parties []*domain.Party) []*backend.Party
	// PartiesModelToDomain converts ocpi model to domain
	PartiesModelToDomain(platformId, status string, roles ...*model.OcpiCredentialRole) []*domain.Party
	// PartiesDomainToModel converts domain to ocpi model
	PartiesDomainToModel(parties ...*domain.Party) []*model.OcpiCredentialRole
	// PartyBackendToDomain converts party backend to domain
	PartyBackendToDomain(platformId string, p *backend.Party) *domain.Party
	// PlatformBackendToDomain converts platform backend to domain
	PlatformBackendToDomain(rq *backend.PlatformRequest) *domain.Platform
	// PlatformDomainToBackend converts platform domain to Backend
	PlatformDomainToBackend(p *domain.Platform) *backend.Platform
}

type CredentialsUc interface {
	// EstablishConnection establishes connection with the remote platform
	EstablishConnection(ctx context.Context, receiverPlatformId string) (*domain.Platform, error)
	// UpdateConnection updates existent connection with the remote platform
	UpdateConnection(ctx context.Context, receiverPlatformId string) (*domain.Platform, error)
	// AcceptConnection accepts connection from the remote platform
	AcceptConnection(ctx context.Context, senderPlatformId string, rq *model.OcpiCredentials) (*model.OcpiCredentials, error)
	// OnRemoteGetCredentials remote platform requests local credentials
	OnRemoteGetCredentials(ctx context.Context, platformId string) (*model.OcpiCredentials, error)
	// OnRemoteDeleteCredentials remote platform requests deletes credentials
	OnRemoteDeleteCredentials(ctx context.Context, platformId string) error
	// OnRemotePartyPull initiates by cron, goes through all the connected platforms, retrieves and updates current list of parties
	OnRemotePartyPull(ctx context.Context) error
	// OnLocalPartyChanged party changed in local platform
	OnLocalPartyChanged(ctx context.Context, party *domain.Party) error
}

type RemotePlatformRepository interface {
	// GetVersions requests versions of the remote platform
	GetVersions(ctx context.Context, rq *OcpiRepositoryBaseRequest) (domain.Versions, error)
	// GetVersionDetails requests version details  of the remote platform
	GetVersionDetails(ctx context.Context, rq *OcpiRepositoryBaseRequest) (domain.ModuleEndpoints, error)
	// PostCredentials posts credentials to establish connection
	PostCredentials(ctx context.Context, rq *OcpiRepositoryRequestG[*model.OcpiCredentials]) (*model.OcpiCredentials, error)
	// PutCredentials updates credentials for existent connection
	PutCredentials(ctx context.Context, rq *OcpiRepositoryRequestG[*model.OcpiCredentials]) (*model.OcpiCredentials, error)
	// GetCredentials requests credentials for existent connection
	GetCredentials(ctx context.Context, rq *OcpiRepositoryBaseRequest) (*model.OcpiCredentials, error)
	// DeleteCredentials requests credentials deletion
	DeleteCredentials(ctx context.Context, rq *OcpiRepositoryBaseRequest) error
}
