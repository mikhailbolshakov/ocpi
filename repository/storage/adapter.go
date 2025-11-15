package storage

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	kitCluster "github.com/mikhailbolshakov/kit/cluster"
	"github.com/mikhailbolshakov/kit/storages/pg"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/domain"
)

// Adapter provides a contract to access a remote service
type Adapter interface {
	kitCluster.Adapter
	domain.PlatformStorage
	domain.PartyStorage
	domain.OcpiLogStorage
	domain.LocationStorage
	domain.TariffStorage
	domain.TokenStorage
	domain.SessionStorage
	domain.CommandStorage
	domain.CdrStorage
	backend.WebhookStorage
}

// adapterImpl implements storage adapter
type adapterImpl struct {
	*platformStorageImpl
	*partyStorageImpl
	*logStorageImpl
	*locationStorageImpl
	*tariffStorageImpl
	*tokenStorageImpl
	*webhookStorageImpl
	*sessionStorageImpl
	*commandStorageImpl
	*cdrStorageImpl
	pg *pg.Storage
}

// NewAdapter creates a new instance of the adapter
func NewAdapter() Adapter {
	return &adapterImpl{}
}

func (a *adapterImpl) l() kit.CLogger {
	return ocpi.L().Cmp("adapter")
}

func (a *adapterImpl) Init(ctx context.Context, cfg interface{}) error {
	a.l().C(ctx).Mth("init").Dbg()

	config := cfg.(*ocpi.CfgStorages)

	// init postgres
	var err error
	a.pg, err = pg.Open(config.Database.Master, ocpi.LF())
	if err != nil {
		return err
	}
	// applying migrations
	if config.Database.MigPath != "" {
		db, _ := a.pg.Instance.DB()
		m := pg.NewMigration(db, config.Database.MigPath, ocpi.LF())
		if err := m.Up(); err != nil {
			return err
		}
	}

	// init storages
	a.platformStorageImpl = newPlatformStorage(a.pg)
	a.partyStorageImpl = newPartyStorage(a.pg)
	a.logStorageImpl = newLogStorage(a.pg)
	if err := a.logStorageImpl.init(ctx, 0, 0); err != nil {
		return err
	}
	a.locationStorageImpl = newLocationStorage(a.pg)
	a.webhookStorageImpl = newWebhookStorage(a.pg)
	a.tariffStorageImpl = newTariffStorage(a.pg)
	a.tokenStorageImpl = newTokenStorage(a.pg)
	a.sessionStorageImpl = newSessionStorage(a.pg)
	a.commandStorageImpl = newCommandStorage(a.pg)
	a.cdrStorageImpl = newCdrStorage(a.pg)

	return nil
}

func (a *adapterImpl) Close(ctx context.Context) error {
	if a.pg != nil {
		a.pg.Close()
	}
	a.logStorageImpl.close(ctx)
	return nil
}
