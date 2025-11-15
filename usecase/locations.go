package usecase

import (
	"context"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/model"
	"time"
)

type LocationConverter interface {
	// LocationDomainToModel converts location domain to ocpi model
	LocationDomainToModel(loc *domain.Location) *model.OcpiLocation
	// LocationsDomainToModel converts locations domain to ocpi model
	LocationsDomainToModel(locs []*domain.Location) []*model.OcpiLocation
	// LocationModelToDomain converts location model to domain
	LocationModelToDomain(loc *model.OcpiLocation, platformId string) *domain.Location
	// LocationDomainToBackend converts location domain to backend
	LocationDomainToBackend(loc *domain.Location) *backend.Location
	// LocationsDomainToBackend converts location domain to backend
	LocationsDomainToBackend(locs []*domain.Location) []*backend.Location
	// LocationBackendToDomain converts location backend to domain
	LocationBackendToDomain(loc *backend.Location, platformId string) *domain.Location

	// EvseDomainToModel converts evse domain to ocpi model
	EvseDomainToModel(evse *domain.Evse) *model.OcpiEvse
	// EvsesDomainToModel converts evses domain to ocpi model
	EvsesDomainToModel(evses []*domain.Evse) []*model.OcpiEvse
	// EvseModelToDomain converts evse ocpi model to domain
	EvseModelToDomain(e *model.OcpiEvse, cc, partyId, platformId, locId string) *domain.Evse
	// EvseDomainToBackend converts evse domain to backend
	EvseDomainToBackend(evse *domain.Evse) *backend.Evse
	// EvseBackendToDomain converts evse backend to domain
	EvseBackendToDomain(e *backend.Evse, platformId, locId string) *domain.Evse
	// EvsesDomainToBackend converts evses domain to backend
	EvsesDomainToBackend(evses []*domain.Evse) []*backend.Evse

	// ConnectorDomainToModel converts connector domain to ocpi model
	ConnectorDomainToModel(con *domain.Connector) *model.OcpiConnector
	// ConnectorsDomainToModel converts connectors domain to ocpi model
	ConnectorsDomainToModel(cons []*domain.Connector) []*model.OcpiConnector
	// ConnectorModelToDomain converts connector model to domain
	ConnectorModelToDomain(con *model.OcpiConnector, cc, partyId, platformId, locId, evseId string) *domain.Connector
	// ConnectorDomainToBackend converts connector domain to backend
	ConnectorDomainToBackend(con *domain.Connector) *backend.Connector
	// ConnectorBackendToDomain converts connector backend to domain
	ConnectorBackendToDomain(con *backend.Connector, platformId, locId, evseId string) *domain.Connector
	// ConnectorsDomainToBackend converts connector domain to backend
	ConnectorsDomainToBackend(cons []*domain.Connector) []*backend.Connector
}

type LocationUc interface {
	// OnLocalLocationChanged handles changing location in local platform
	OnLocalLocationChanged(ctx context.Context, loc *domain.Location) error
	// OnRemoteLocationsPull handles request to pull locations from remote platforms (fired by cron)
	OnRemoteLocationsPull(ctx context.Context, from, to *time.Time) error
	// OnRemoteLocationsPullWhenPushNotSupported handles request to pull locations from remote platforms which don't support push (fired by cron)
	OnRemoteLocationsPullWhenPushNotSupported(ctx context.Context, from, to *time.Time) error
	// OnRemoteLocationPut handles put location in remote platform
	OnRemoteLocationPut(ctx context.Context, platformId string, loc *model.OcpiLocation) error
	// OnRemoteLocationPatch handles patch location in remote platform
	OnRemoteLocationPatch(ctx context.Context, platformId string, loc *model.OcpiLocation) error

	// OnLocalEvseChanged handles changing evse in local platform
	OnLocalEvseChanged(ctx context.Context, evse *domain.Evse) error
	// OnLocalEvseStatusChanged handles changing status of evse in local platform
	OnLocalEvseStatusChanged(ctx context.Context, locId, evseId, status string) error
	// OnRemoteEvsePut handles put location in remote platform
	OnRemoteEvsePut(ctx context.Context, platformId, locId, countryCode, partyId string, evse *model.OcpiEvse) error
	// OnRemoteEvsePatch handles patch location in remote platform
	OnRemoteEvsePatch(ctx context.Context, platformId, locId, countryCode, partyId string, evse *model.OcpiEvse) error

	// OnLocalConnectorChanged handles changing connector in local platform
	OnLocalConnectorChanged(ctx context.Context, con *domain.Connector) error
	// OnRemoteConnectorPut handles put evse in remote platform
	OnRemoteConnectorPut(ctx context.Context, platformId, locId, evseId, countryCode, partyId string, con *model.OcpiConnector) error
	// OnRemoteConnectorPatch handles patch connector in remote platform
	OnRemoteConnectorPatch(ctx context.Context, platformId, locId, evseId, countryCode, partyId string, con *model.OcpiConnector) error
}

type RemoteLocationRepository interface {
	// PutLocationAsync puts location
	PutLocationAsync(ctx context.Context, rq *OcpiRepositoryErrHandlerRequestG[*model.OcpiLocation])
	// PatchLocationAsync patches location
	PatchLocationAsync(ctx context.Context, rq *OcpiRepositoryErrHandlerRequestG[*model.OcpiLocation])
	// GetLocations retrieves locations
	GetLocations(ctx context.Context, rq *OcpiRepositoryPagingRequest) ([]*model.OcpiLocation, error)
	// GetLocation retrieves location by id
	GetLocation(ctx context.Context, rq *OcpiRepositoryIdRequest) (*model.OcpiLocation, error)
	// PutEvseAsync puts evse
	PutEvseAsync(ctx context.Context, rq *OcpiRepositoryErrHandlerRequestG[*model.OcpiEvse], party *model.OcpiPartyId, locId string)
	// PatchEvseAsync patches evse
	PatchEvseAsync(ctx context.Context, rq *OcpiRepositoryErrHandlerRequestG[*model.OcpiEvse], party *model.OcpiPartyId, locId string)
	// GetEvse retrieves evse
	GetEvse(ctx context.Context, rq *OcpiRepositoryBaseRequest, locId, evseId string) (*model.OcpiEvse, error)
	// PutConnectorAsync puts connector
	PutConnectorAsync(ctx context.Context, rq *OcpiRepositoryErrHandlerRequestG[*model.OcpiConnector], party *model.OcpiPartyId, evseId, locId string)
	// PatchConnectorAsync patches connector
	PatchConnectorAsync(ctx context.Context, rq *OcpiRepositoryErrHandlerRequestG[*model.OcpiConnector], party *model.OcpiPartyId, evseId, locId string)
	// GetConnector retrieves connector
	GetConnector(ctx context.Context, rq *OcpiRepositoryBaseRequest, locId, evseId, conId string) (*model.OcpiConnector, error)
}
