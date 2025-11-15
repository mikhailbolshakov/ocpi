package usecase

import (
	"context"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/model"
	"time"
)

type TariffConverter interface {
	// TariffDomainToModel converts tariff domain to ocpi model
	TariffDomainToModel(trf *domain.Tariff) *model.OcpiTariff
	// TariffsDomainToModel converts tariffs domain to ocpi model
	TariffsDomainToModel(ts []*domain.Tariff) []*model.OcpiTariff
	// TariffModelToDomain converts tariff model to domain
	TariffModelToDomain(trf *model.OcpiTariff, platformId string) *domain.Tariff
	// TariffsModelToDomain converts tariffs model to domain
	TariffsModelToDomain(ts []*model.OcpiTariff, platformId string) []*domain.Tariff
	// TariffDomainToBackend converts tariff domain to backend
	TariffDomainToBackend(trf *domain.Tariff) *backend.Tariff
	// TariffsDomainToBackend converts tariff domain to backend
	TariffsDomainToBackend(ts []*domain.Tariff) []*backend.Tariff
	// TariffsBackendToDomain converts tariffs backend to domain
	TariffsBackendToDomain(ts []*backend.Tariff) []*domain.Tariff
	// TariffBackendToDomain converts tariff backend to domain
	TariffBackendToDomain(trf *backend.Tariff) *domain.Tariff
}

type TariffUc interface {
	// OnLocalTariffChanged handles changing tariff in local platform
	OnLocalTariffChanged(ctx context.Context, trf *domain.Tariff) error
	// OnRemoteTariffsPull handles request to pull tariffs from remote platforms (fired by cron)
	OnRemoteTariffsPull(ctx context.Context, from, to *time.Time) error
	// OnRemoteTariffsPullWhenPushNotSupported handles request to pull tariffs from remote platforms which don't support push (fired by cron)
	OnRemoteTariffsPullWhenPushNotSupported(ctx context.Context, from, to *time.Time) error
	// OnRemoteTariffPut handles put tariff in remote platform
	OnRemoteTariffPut(ctx context.Context, platformId string, trf *model.OcpiTariff) error
	// OnRemoteTariffPatch handles patch tariff in remote platform
	OnRemoteTariffPatch(ctx context.Context, platformId string, trf *model.OcpiTariff) error
}

type RemoteTariffRepository interface {
	// PutTariffAsync puts tariff
	PutTariffAsync(ctx context.Context, rq *OcpiRepositoryErrHandlerRequestG[*model.OcpiTariff])
	// PatchTariffAsync patches tariff
	PatchTariffAsync(ctx context.Context, rq *OcpiRepositoryErrHandlerRequestG[*model.OcpiTariff])
	// GetTariffs retrieves tariffs
	GetTariffs(ctx context.Context, rq *OcpiRepositoryPagingRequest) ([]*model.OcpiTariff, error)
	// GetTariff retrieves tariff by id
	GetTariff(ctx context.Context, rq *OcpiRepositoryIdRequest) (*model.OcpiTariff, error)
}
