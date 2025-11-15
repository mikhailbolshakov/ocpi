package usecase

import (
	"context"
	"github.com/mikhailbolshakov/ocpi/domain"
)

type MaintenanceUc interface {
	// DeleteLocalPartyByExt deletes local party and all related entities
	// !!!! Be careful!! Use for testing purposes only
	DeleteLocalPartyByExt(ctx context.Context, extId domain.PartyExtId) error
}
