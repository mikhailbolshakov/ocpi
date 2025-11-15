package usecase

import (
	"context"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/model"
	"time"
)

type HubUc interface {
	// OnRemoteClientInfosPull handles request to pull client info from remote platforms (fired by cron)
	OnRemoteClientInfosPull(ctx context.Context, from, to *time.Time) error
	// OnRemoteClientInfosPullWhenPushNotSupported handles request to pull client info from remote platforms which don't support push (fired by cron)
	OnRemoteClientInfosPullWhenPushNotSupported(ctx context.Context, from, to *time.Time) error
	// OnRemoteClientInfoPut executed when a remote platform sends put request
	OnRemoteClientInfoPut(ctx context.Context, platformId string, loc *model.OcpiClientInfo) error
	// OnLocalClientInfoChanged executed when local client info changed
	OnLocalClientInfoChanged(ctx context.Context, party *domain.Party) error
}

type RemoteHubClientInfoRepository interface {
	// PutClientInfo puts client info to remote platform
	PutClientInfo(ctx context.Context, rq *OcpiRepositoryRequestG[*model.OcpiClientInfo]) error
	// PutClientInfoAsync puts client info to remote platform
	PutClientInfoAsync(ctx context.Context, rq *OcpiRepositoryErrHandlerRequestG[*model.OcpiClientInfo])
	// GetClientInfos retrieves client info from remote platforms
	GetClientInfos(ctx context.Context, rq *OcpiRepositoryPagingRequest) ([]*model.OcpiClientInfo, error)
}
