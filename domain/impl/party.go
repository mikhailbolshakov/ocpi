package impl

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/goroutine"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/errors"
)

type partyService struct {
	base
	storage domain.PartyStorage
}

func NewPartyService(storage domain.PartyStorage) domain.PartyService {
	return &partyService{
		storage: storage,
	}
}

func (s *partyService) l() kit.CLogger {
	return ocpi.L().Cmp("party-svc")
}

func (s *partyService) GetByExtId(ctx context.Context, extId domain.PartyExtId) (*domain.Party, error) {
	s.l().C(ctx).Mth("get-by-ext-id").Dbg()
	if extId.PartyId == "" || extId.CountryCode == "" {
		return nil, errors.ErrPartyIdEmpty(ctx)
	}
	return s.storage.GetPartyByExtId(ctx, extId)
}

func (s *partyService) GetByPlatform(ctx context.Context, platformId string) ([]*domain.Party, error) {
	s.l().C(ctx).Mth("get-by-platform").Dbg()
	if platformId == "" {
		return nil, errors.ErrPlatformIdEmpty(ctx)
	}
	return s.storage.GetPartiesByPlatform(ctx, platformId)
}

func (s *partyService) Get(ctx context.Context, id string) (*domain.Party, error) {
	s.l().C(ctx).Mth("get").Dbg()
	if id == "" {
		return nil, errors.ErrPartyIdEmpty(ctx)
	}
	return s.storage.GetParty(ctx, id)
}

func (s *partyService) Search(ctx context.Context, cr *domain.PartySearchCriteria) (*domain.PartySearchResponse, error) {
	s.l().C(ctx).Mth("search").Dbg()
	if cr.Limit == nil {
		cr.Limit = kit.IntPtr(20)
	}
	if cr.Offset == nil {
		cr.Offset = kit.IntPtr(0)
	}
	return s.storage.Search(ctx, cr)
}

func (s *partyService) Merge(ctx context.Context, party *domain.Party) (*domain.Party, error) {
	l := s.l().C(ctx).Mth("merge").F(kit.KV{"partyId": party.Id}).Dbg()

	// validate
	err := s.validate(ctx, party)
	if err != nil {
		return nil, err
	}

	var stored *domain.Party

	// search by id
	if party.Id != "" {
		stored, err = s.storage.GetParty(ctx, party.Id)
		if err != nil {
			return nil, err
		}
	}

	// search by party
	if stored == nil && party.PlatformId != "" {
		stored, err = s.storage.GetPartyByExtId(ctx, party.ExtId)
		if err != nil {
			return nil, err
		}
	}

	if stored != nil {
		// check last_updated
		if party.LastUpdated.Before(stored.LastUpdated) {
			l.Warn("later changes found")
			return nil, nil
		}
		party.Id = stored.Id
		party.LastSent = stored.LastSent
		party.PlatformId = stored.PlatformId
		if party.RefId == "" {
			party.RefId = stored.RefId
		}
		if party.BusinessDetails == nil {
			party.BusinessDetails = stored.BusinessDetails
		}
	}

	if party.BusinessDetails == nil {
		party.BusinessDetails = &domain.BusinessDetails{}
	}
	if party.BusinessDetails.Name == "" {
		party.BusinessDetails.Name = "None"
	}

	// create/update
	if stored == nil {
		party.Id = kit.NewId()
		err = s.storage.CreateParty(ctx, party)
		if err != nil {
			return nil, err
		}
	} else {
		err = s.storage.UpdateParty(ctx, party)
		if err != nil {
			return nil, err
		}
	}

	return party, nil
}

func (s *partyService) MergeMany(ctx context.Context, parties ...*domain.Party) error {
	l := s.l().C(ctx).Mth("merge-many").Dbg()

	// check if not empty
	if len(parties) == 0 {
		return nil
	}

	ch := make(chan *domain.Party, 1)

	numWorkers := len(parties)%10 + 1
	eg := goroutine.NewGroup(ctx).WithLogger(l)

	// publisher
	eg.Go(func() error {
		defer close(ch)
		for _, party := range parties {
			ch <- party
		}
		return nil
	})

	// workers
	for i := 0; i < numWorkers; i++ {
		eg.Go(func() error {
			for party := range ch {
				_, err := s.Merge(ctx, party)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}

	return eg.Wait()
}

func (s *partyService) MarkSent(ctx context.Context, partyIds ...string) error {
	s.l().C(ctx).Mth("mark-sent").Dbg()
	if len(partyIds) == 0 {
		return nil
	}
	return s.storage.MarkSentParties(ctx, kit.Now(), partyIds...)
}

func (s *partyService) DeletePartyByExtId(ctx context.Context, extId domain.PartyExtId) error {
	s.l().C(ctx).Mth("del-ext").F(kit.KV{"partyId": extId.PartyId, "country": extId.CountryCode}).Dbg()
	if err := s.validateExtId(ctx, extId); err != nil {
		return err
	}
	return s.storage.DeletePartyByExtId(ctx, extId)
}

func (s *partyService) validate(ctx context.Context, party *domain.Party) error {

	if err := s.validateOcpiItem(ctx, &party.OcpiItem); err != nil {
		return err
	}
	if len(party.Roles) == 0 {
		return errors.ErrPartyRolesEmpty(ctx)
	}
	if err := s.validateBusinessDetails(ctx, "business_details", party.BusinessDetails); err != nil {
		return err
	}

	// roles
	party.Roles = kit.Strings(party.Roles).Distinct()
	for _, r := range party.Roles {
		if _, ok := domain.RoleMap[r]; !ok {
			return errors.ErrPartyRoleInvalid(ctx)
		}
	}

	return nil
}
