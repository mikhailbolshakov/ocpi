package storage

import (
	"context"
	"github.com/jackc/pgtype"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/goroutine"
	"github.com/mikhailbolshakov/kit/storages/pg"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/errors"
	"gorm.io/gorm"
	"time"
)

type location struct {
	pg.GormDto
	Id          string        `gorm:"column:id;primaryKey"`
	Details     *pgtype.JSONB `gorm:"column:details"`
	PartyId     string        `gorm:"column:party_id"`
	CountryCode string        `gorm:"column:country_code"`
	PlatformId  string        `gorm:"column:platform_id"`
	RefId       *string       `gorm:"column:ref_id"`
	LastUpdated time.Time     `gorm:"column:last_updated"`
	LastSent    *time.Time    `gorm:"column:last_sent"`
}

type locationRead struct {
	Location   location   `gorm:"embedded"`
	TotalCount totalCount `gorm:"embedded"`
}

type evse struct {
	pg.GormDto
	Id          string        `gorm:"column:id;primaryKey"`
	LocationId  string        `gorm:"column:location_id"`
	Status      string        `gorm:"column:status"`
	Details     *pgtype.JSONB `gorm:"column:details"`
	PartyId     string        `gorm:"column:party_id"`
	CountryCode string        `gorm:"column:country_code"`
	PlatformId  string        `gorm:"column:platform_id"`
	RefId       *string       `gorm:"column:ref_id"`
	LastUpdated time.Time     `gorm:"column:last_updated"`
	LastSent    *time.Time    `gorm:"column:last_sent"`
}

type evseRead struct {
	Evse       evse       `gorm:"embedded"`
	TotalCount totalCount `gorm:"embedded"`
}

type connector struct {
	pg.GormDto
	Id          string        `gorm:"column:id;primaryKey"`
	LocationId  string        `gorm:"column:location_id"`
	EvseId      string        `gorm:"column:evse_id"`
	Details     *pgtype.JSONB `gorm:"column:details"`
	PartyId     string        `gorm:"column:party_id"`
	CountryCode string        `gorm:"column:country_code"`
	PlatformId  string        `gorm:"column:platform_id"`
	RefId       *string       `gorm:"column:ref_id"`
	LastUpdated time.Time     `gorm:"column:last_updated"`
	LastSent    *time.Time    `gorm:"column:last_sent"`
}

type connectorRead struct {
	Connector  connector  `gorm:"embedded"`
	TotalCount totalCount `gorm:"embedded"`
}

type locationStorageImpl struct {
	pg *pg.Storage
}

func (s *locationStorageImpl) l() kit.CLogger {
	return ocpi.L().Cmp("loc-storage")
}

func newLocationStorage(pg *pg.Storage) *locationStorageImpl {
	return &locationStorageImpl{
		pg: pg,
	}
}

func (s *locationStorageImpl) GetLocation(ctx context.Context, id string, withEvse bool) (*domain.Location, error) {
	l := s.l().C(ctx).Mth("get-loc").F(kit.KV{"locId": id}).Dbg()

	if id == "" {
		return nil, nil
	}

	eg := goroutine.NewGroup(ctx).WithLogger(l)

	locDto := &location{}
	var evseDtos []*evse
	var conDtos []*connector

	eg.Go(func() error {
		res := s.pg.Instance.Where("id = ?", id).Limit(1).Find(&locDto)
		if res.Error != nil {
			return errors.ErrLocStorageGet(ctx, res.Error)
		}
		if res.RowsAffected == 0 {
			locDto = nil
		}
		return nil
	})

	if withEvse {
		// retrieve evses
		eg.Go(func() error {
			if err := s.pg.Instance.Where("location_id = ?", id).Find(&evseDtos).Error; err != nil {
				return errors.ErrEvseStorageGet(ctx, err)
			}
			return nil
		})
		// retrieve connectors
		eg.Go(func() error {
			if err := s.pg.Instance.Where("location_id = ?", id).Find(&conDtos).Error; err != nil {
				return errors.ErrConStorageGet(ctx, err)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return s.toLocationDomain(locDto, evseDtos, conDtos), nil
}

func (s *locationStorageImpl) MergeLocation(ctx context.Context, loc *domain.Location) error {
	l := s.l().C(ctx).Mth("merge-loc").F(kit.KV{"locId": loc.Id}).Dbg()
	eg := goroutine.NewGroup(ctx).WithLogger(l)

	// transaction
	tx := s.pg.Instance.Begin()
	if tx.Error != nil {
		return errors.ErrLocStorageTx(ctx, tx.Error)
	}

	// merge location
	eg.Go(func() error {
		return s.pg.Instance.Scopes(merge()).Create(s.toLocationDto(loc)).Error
	})

	// merge evses
	for _, evse := range loc.Evses {
		evse := evse
		eg.Go(func() error {
			return s.pg.Instance.Scopes(merge()).Create(s.toEvseDto(evse)).Error
		})
	}

	// merge connectors
	for _, evse := range loc.Evses {
		for _, con := range evse.Connectors {
			con := con
			eg.Go(func() error {
				return s.pg.Instance.Scopes(merge()).Create(s.toConnectorDto(con)).Error
			})
		}
	}

	// rollback if error
	if err := eg.Wait(); err != nil {
		tx.Rollback()
		return errors.ErrLocStorageMerge(ctx, err)
	}
	tx.Commit()
	return nil
}

func (s *locationStorageImpl) UpdateLocation(ctx context.Context, loc *domain.Location) error {
	s.l().C(ctx).Mth("update-loc").F(kit.KV{"locId": loc.Id}).Dbg()
	if err := s.pg.Instance.Scopes(update()).Save(s.toLocationDto(loc)).Error; err != nil {
		return errors.ErrLocStorageUpdate(ctx, err)
	}
	return nil
}

func (s *locationStorageImpl) DeleteLocationsByExtId(ctx context.Context, extId domain.PartyExtId) error {
	l := s.l().C(ctx).Mth("delete-ext").F(kit.KV{"partyId": extId.PartyId, "country": extId.CountryCode}).Dbg()
	if extId.PartyId == "" || extId.CountryCode == "" {
		return nil
	}
	eg := goroutine.NewGroup(ctx).WithLogger(l)

	// transaction
	tx := s.pg.Instance.Begin()
	if tx.Error != nil {
		return errors.ErrEvseStorageTx(ctx, tx.Error)
	}

	// delete locations
	eg.Go(func() error {
		return s.pg.Instance.Where(`party_id = ? and country_code = ?`, extId.PartyId, extId.CountryCode).
			Delete(&location{}).Error
	})

	// delete evses
	eg.Go(func() error {
		return s.pg.Instance.Where(`party_id = ? and country_code = ?`, extId.PartyId, extId.CountryCode).
			Delete(&evse{}).Error
	})

	// delete connectors
	eg.Go(func() error {
		return s.pg.Instance.Where(`party_id = ? and country_code = ?`, extId.PartyId, extId.CountryCode).
			Delete(&connector{}).Error
	})

	// rollback if error
	if err := eg.Wait(); err != nil {
		tx.Rollback()
		return errors.ErrConStorageDelete(ctx, err)
	}
	tx.Commit()
	return nil
}

func (s *locationStorageImpl) SearchLocations(ctx context.Context, cr *domain.LocationSearchCriteria) (*domain.LocationSearchResponse, error) {
	l := s.l().Mth("search-loc").C(ctx).Dbg()

	rs := &domain.LocationSearchResponse{
		PageResponse: domain.PageResponse{
			Limit: pagingLimit(cr.PageRequest.Limit),
			Total: kit.IntPtr(0),
		},
	}

	// make query
	var locDtosRead []*locationRead
	var evseDtos []*evse
	var conDtos []*connector

	if err := s.pg.Instance.
		Scopes(s.buildLocSearchQuery(cr), paging(cr.PageRequest), orderByLastUpdated(true)).
		Find(&locDtosRead).Error; err != nil {
		return nil, errors.ErrLocStorageGetDb(ctx, err)
	}

	if len(locDtosRead) == 0 {
		return rs, nil
	}

	// gather location Ids
	locIds := make([]string, 0, len(locDtosRead))
	locDtos := make([]*location, 0, len(locDtosRead))
	for _, loc := range locDtosRead {
		locIds = append(locIds, loc.Location.Id)
		locDtos = append(locDtos, &loc.Location)
	}

	if len(locIds) > 0 {
		eg := goroutine.NewGroup(ctx).WithLogger(l)

		// retrieve evses
		eg.Go(func() error {
			if err := s.pg.Instance.Where("location_id in (?)", locIds).Find(&evseDtos).Error; err != nil {
				return errors.ErrEvseStorageGet(ctx, err)
			}
			return nil
		})

		// retrieve connectors
		eg.Go(func() error {
			if err := s.pg.Instance.Where("location_id in (?)", locIds).Find(&conDtos).Error; err != nil {
				return errors.ErrConStorageGet(ctx, err)
			}
			return nil
		})

		if err := eg.Wait(); err != nil {
			return nil, err
		}

	}

	// build response
	rs.Items = s.toLocationsDomain(locDtos, evseDtos, conDtos)
	rs.Total = &locDtosRead[0].TotalCount.TotalCount
	rs.NextPage = nextPage(cr.PageRequest, rs.Total)

	return rs, nil
}

func (s *locationStorageImpl) GetEvse(ctx context.Context, locId, evseId string, withConnector bool) (*domain.Evse, error) {
	l := s.l().C(ctx).Mth("get-evse").F(kit.KV{"evseId": evseId}).Dbg()

	if locId == "" || evseId == "" {
		return nil, nil
	}

	eg := goroutine.NewGroup(ctx).WithLogger(l)

	evseDto := &evse{}
	var conDtos []*connector

	eg.Go(func() error {
		res := s.pg.Instance.Where("id = ? and location_id = ?", evseId, locId).Limit(1).Find(&evseDto)
		if res.Error != nil {
			return errors.ErrEvseStorageGet(ctx, res.Error)
		}
		if res.RowsAffected == 0 {
			evseDto = nil
		}
		return nil
	})

	// retrieve connectors
	if withConnector {
		eg.Go(func() error {
			if err := s.pg.Instance.Where("evse_id = ? and location_id = ?", evseId, locId).Find(&conDtos).Error; err != nil {
				return errors.ErrConStorageGet(ctx, err)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return s.toEvseDomain(evseDto, conDtos), nil
}

func (s *locationStorageImpl) SearchEvses(ctx context.Context, cr *domain.EvseSearchCriteria) (*domain.EvseSearchResponse, error) {
	s.l().Mth("search-evse").C(ctx).Dbg()

	rs := &domain.EvseSearchResponse{
		PageResponse: domain.PageResponse{
			Limit: pagingLimit(cr.PageRequest.Limit),
			Total: kit.IntPtr(0),
		},
	}

	// make query
	var evseDtosRead []*evseRead
	var conDtos []*connector

	if err := s.pg.Instance.
		Scopes(s.buildEvseSearchQuery(cr), paging(cr.PageRequest), orderByLastUpdated(true)).
		Find(&evseDtosRead).Error; err != nil {
		return nil, errors.ErrEvseStorageGetDb(ctx, err)
	}

	if len(evseDtosRead) == 0 {
		return rs, nil
	}

	// gather evse Ids
	evseIds := make([]string, 0, len(evseDtosRead))
	evseDtos := make([]*evse, 0, len(evseDtosRead))
	for _, evse := range evseDtosRead {
		evseIds = append(evseIds, evse.Evse.Id)
		evseDtos = append(evseDtos, &evse.Evse)
	}

	if len(evseIds) > 0 {
		if err := s.pg.Instance.Where("evse_id in (?)", evseIds).Find(&conDtos).Error; err != nil {
			return nil, errors.ErrConStorageGet(ctx, err)
		}
	}

	conMap := make(map[string][]*connector)
	for _, con := range conDtos {
		conMap[con.EvseId] = append(conMap[con.EvseId], con)
	}

	// build response
	rs.Items = s.toEvsesDomain(evseDtos, conMap)
	rs.Total = &evseDtosRead[0].TotalCount.TotalCount
	return rs, nil
}

func (s *locationStorageImpl) UpdateEvse(ctx context.Context, evse *domain.Evse) error {
	l := s.l().C(ctx).Mth("update-evse").F(kit.KV{"evseId": evse.Id}).Dbg()

	eg := goroutine.NewGroup(ctx).WithLogger(l)

	// transaction
	tx := s.pg.Instance.Begin()
	if tx.Error != nil {
		return errors.ErrEvseStorageTx(ctx, tx.Error)
	}

	// update evse
	eg.Go(func() error {
		return s.pg.Instance.Scopes(update()).Save(s.toEvseDto(evse)).Error
	})

	// update last_updated
	eg.Go(func() error {
		return s.pg.Instance.Model(&location{Id: evse.LocationId}).Scopes(update()).
			Update("last_updated", evse.LastUpdated).
			Error
	})

	// update evse
	eg.Go(func() error {
		return s.pg.Instance.Scopes(update()).Save(s.toEvseDto(evse)).Error
	})

	// rollback if error
	if err := eg.Wait(); err != nil {
		tx.Rollback()
		return errors.ErrEvseStorageUpdate(ctx, err)
	}
	tx.Commit()

	return nil
}

func (s *locationStorageImpl) MergeEvse(ctx context.Context, evse *domain.Evse) error {
	l := s.l().C(ctx).Mth("merge-evse").F(kit.KV{"evseId": evse.Id}).Dbg()
	eg := goroutine.NewGroup(ctx).WithLogger(l)

	// transaction
	tx := s.pg.Instance.Begin()
	if tx.Error != nil {
		return errors.ErrEvseStorageTx(ctx, tx.Error)
	}

	// update last_updated
	eg.Go(func() error {
		return s.pg.Instance.Scopes(update()).Model(&location{Id: evse.LocationId}).
			Update("last_updated", evse.LastUpdated).
			Error
	})

	// merge evse
	eg.Go(func() error {
		return s.pg.Instance.Scopes(merge()).Create(s.toEvseDto(evse)).Error
	})

	// merge evses
	for _, con := range evse.Connectors {
		con := con
		eg.Go(func() error {
			return s.pg.Instance.Scopes(merge()).Create(s.toConnectorDto(con)).Error
		})
	}

	// rollback if error
	if err := eg.Wait(); err != nil {
		tx.Rollback()
		return errors.ErrEvseStorageMerge(ctx, err)
	}
	tx.Commit()
	return nil
}

func (s *locationStorageImpl) MergeConnector(ctx context.Context, con *domain.Connector) error {
	l := s.l().C(ctx).Mth("merge-con").F(kit.KV{"conId": con.Id}).Dbg()
	eg := goroutine.NewGroup(ctx).WithLogger(l)

	// transaction
	tx := s.pg.Instance.Begin()
	if tx.Error != nil {
		return errors.ErrEvseStorageTx(ctx, tx.Error)
	}

	// update last_updated
	eg.Go(func() error {
		return s.pg.Instance.Scopes(update()).Model(&location{Id: con.LocationId}).
			Update("last_updated", con.LastUpdated).
			Error
	})

	// update last_updated
	eg.Go(func() error {
		return s.pg.Instance.Scopes(update()).Model(&evse{Id: con.EvseId}).
			Update("last_updated", con.LastUpdated).
			Error
	})

	// merge connector
	eg.Go(func() error {
		return s.pg.Instance.Scopes(merge()).Create(s.toConnectorDto(con)).Error
	})

	// rollback if error
	if err := eg.Wait(); err != nil {
		tx.Rollback()
		return errors.ErrConStorageMerge(ctx, err)
	}
	tx.Commit()
	return nil

}

func (s *locationStorageImpl) UpdateConnector(ctx context.Context, con *domain.Connector) error {
	l := s.l().C(ctx).Mth("update-con").F(kit.KV{"conId": con.Id}).Dbg()
	eg := goroutine.NewGroup(ctx).WithLogger(l)

	// transaction
	tx := s.pg.Instance.Begin()
	if tx.Error != nil {
		return errors.ErrEvseStorageTx(ctx, tx.Error)
	}

	// update last_updated
	eg.Go(func() error {
		return s.pg.Instance.Scopes(update()).Model(&location{Id: con.LocationId}).
			Update("last_updated", con.LastUpdated).
			Error
	})

	// update last_updated
	eg.Go(func() error {
		return s.pg.Instance.Scopes(update()).Model(&evse{Id: con.EvseId}).
			Update("last_updated", con.LastUpdated).
			Error
	})

	// update connector
	eg.Go(func() error {
		return s.pg.Instance.Scopes(update()).Save(s.toConnectorDto(con)).Error
	})

	// rollback if error
	if err := eg.Wait(); err != nil {
		tx.Rollback()
		return errors.ErrConStorageUpdate(ctx, err)
	}
	tx.Commit()
	return nil
}

func (s *locationStorageImpl) SearchConnectors(ctx context.Context, cr *domain.ConnectorSearchCriteria) (*domain.ConnectorSearchResponse, error) {
	s.l().Mth("search-con").C(ctx).Dbg()

	rs := &domain.ConnectorSearchResponse{
		PageResponse: domain.PageResponse{
			Limit: pagingLimit(cr.PageRequest.Limit),
			Total: kit.IntPtr(0),
		},
	}

	// make query
	var conDtosRead []*connectorRead

	if err := s.pg.Instance.
		Scopes(s.buildConSearchQuery(cr), paging(cr.PageRequest), orderByLastUpdated(true)).
		Find(&conDtosRead).Error; err != nil {
		return nil, errors.ErrConStorageGetDb(ctx, err)
	}

	if len(conDtosRead) == 0 {
		return rs, nil
	}

	conDtos := make([]*connector, 0, len(conDtosRead))
	for _, con := range conDtosRead {
		conDtos = append(conDtos, &con.Connector)
	}

	// build response
	rs.Items = s.toConnectorsDomain(conDtos)
	rs.Total = &conDtosRead[0].TotalCount.TotalCount
	return rs, nil
}

func (s *locationStorageImpl) GetConnector(ctx context.Context, locId, evseId, conId string) (*domain.Connector, error) {
	s.l().C(ctx).Mth("get-con").F(kit.KV{"evseId": evseId, "conId": conId}).Dbg()
	if locId == "" || evseId == "" || conId == "" {
		return nil, nil
	}
	conDto := &connector{}
	res := s.pg.Instance.Where("id = ? and evse_id = ? and location_id = ?", conId, evseId, locId).Find(&conDto)
	if err := res.Error; err != nil {
		return nil, errors.ErrConStorageGet(ctx, err)
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return s.toConnectorDomain(conDto), nil
}

func (s *locationStorageImpl) buildLocSearchQuery(criteria *domain.LocationSearchCriteria) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query := db.Table("locations").Select("locations.*, count(*) over() total_count")
		// populate conditions
		if criteria.ExtId != nil {
			query = query.Where("party_id = ? and country_code = ?", criteria.ExtId.PartyId, criteria.ExtId.CountryCode)
		}
		if criteria.DateTo != nil {
			query = query.Where("last_updated <= ?", *criteria.DateTo)
		}
		if criteria.DateFrom != nil {
			query = query.Where("last_updated >= ?", *criteria.DateFrom)
		}
		if len(criteria.IncPlatforms) > 0 {
			query = query.Where("platform_id in (?)", criteria.IncPlatforms)
		}
		if len(criteria.ExcPlatforms) > 0 {
			query = query.Where("platform_id not in (?)", criteria.ExcPlatforms)
		}
		if len(criteria.Ids) > 0 {
			query = query.Where("id in (?)", criteria.Ids)
		}
		if criteria.RefId != "" {
			query = query.Where("ref_id = ?", criteria.RefId)
		}
		return query
	}
}

func (s *locationStorageImpl) buildConSearchQuery(criteria *domain.ConnectorSearchCriteria) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query := db.Table("connectors").Select("connectors.*, count(*) over() total_count")
		// populate conditions
		if criteria.ExtId != nil {
			query = query.Where("party_id = ? and country_code = ?", criteria.ExtId.PartyId, criteria.ExtId.CountryCode)
		}
		if criteria.DateTo != nil {
			query = query.Where("last_updated <= ?", *criteria.DateTo)
		}
		if criteria.DateFrom != nil {
			query = query.Where("last_updated >= ?", *criteria.DateFrom)
		}
		if criteria.RefId != "" {
			query = query.Where("ref_id = ?", criteria.RefId)
		}
		return query
	}
}

func (s *locationStorageImpl) buildEvseSearchQuery(criteria *domain.EvseSearchCriteria) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query := db.Table("evses").Select("evses.*, count(*) over() total_count")
		// populate conditions
		if criteria.ExtId != nil {
			query = query.Where("party_id = ? and country_code = ?", criteria.ExtId.PartyId, criteria.ExtId.CountryCode)
		}
		if criteria.DateTo != nil {
			query = query.Where("last_updated <= ?", *criteria.DateTo)
		}
		if criteria.DateFrom != nil {
			query = query.Where("last_updated >= ?", *criteria.DateFrom)
		}
		if criteria.RefId != "" {
			query = query.Where("ref_id = ?", criteria.RefId)
		}
		return query
	}
}
