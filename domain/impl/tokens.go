package impl

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/errors"
)

type tokenService struct {
	base
	storage domain.TokenStorage
}

func NewTokenService(storage domain.TokenStorage) domain.TokenService {
	return &tokenService{
		storage: storage,
	}
}

func (s *tokenService) l() kit.CLogger {
	return ocpi.L().Cmp("tkn-svc")
}

var (
	tokenTypeMap = map[string]struct{}{
		domain.TokenTypeAdHocUser: {},
		domain.TokenTypeAppUser:   {},
		domain.TokenTypeRfid:      {},
		domain.TokenTypeOther:     {},
	}

	tokenWLMap = map[string]struct{}{
		domain.TokenWLTypeAlways:         {},
		domain.TokenWLTypeAllowed:        {},
		domain.TokenWLTypeNever:          {},
		domain.TokenWLTypeAllowedOffline: {},
	}

	profileMap = map[string]struct{}{
		domain.ProfileTypeCheap:   {},
		domain.ProfileTypeFast:    {},
		domain.ProfileTypeGreen:   {},
		domain.ProfileTypeRegular: {},
	}
)

func (s *tokenService) PutToken(ctx context.Context, tkn *domain.Token) (*domain.Token, error) {
	l := s.l().C(ctx).Mth("put-tkn").F(kit.KV{"tknId": tkn.Id}).Dbg()

	if tkn.Id == "" {
		return nil, errors.ErrTknIdEmpty(ctx)
	}

	// search by id
	stored, err := s.storage.GetToken(ctx, tkn.Id)
	if err != nil {
		return nil, err
	}

	// check last_updated
	if stored != nil && tkn.LastUpdated.Before(stored.LastUpdated) {
		l.Warn("later changes found")
		return nil, nil
	}

	// validate
	err = s.validateAndPopulatePut(ctx, tkn, stored)
	if err != nil {
		return nil, err
	}

	err = s.storage.MergeToken(ctx, tkn)
	if err != nil {
		return nil, err
	}

	return tkn, nil
}

func (s *tokenService) MergeToken(ctx context.Context, tkn *domain.Token) (*domain.Token, error) {
	l := s.l().C(ctx).Mth("merge-tkn").F(kit.KV{"tknId": tkn.Id}).Dbg()

	if tkn.Id == "" {
		return nil, errors.ErrTknIdEmpty(ctx)
	}

	// search by id
	stored, err := s.storage.GetToken(ctx, tkn.Id)
	if err != nil {
		return nil, err
	}
	if stored == nil {
		return nil, errors.ErrTknNotFound(ctx)
	}

	// check last_updated
	if tkn.LastUpdated.Before(stored.LastUpdated) {
		l.Warn("later changes found")
		return nil, nil
	}

	// validate
	err = s.validateAndPopulateMerge(ctx, tkn, stored)
	if err != nil {
		return nil, err
	}

	err = s.storage.UpdateToken(ctx, stored)
	if err != nil {
		return nil, err
	}

	return stored, nil
}

func (s *tokenService) GetToken(ctx context.Context, tknId string) (*domain.Token, error) {
	s.l().C(ctx).Mth("get-tkn").Dbg()
	if tknId == "" {
		return nil, errors.ErrTknIdEmpty(ctx)
	}
	return s.storage.GetToken(ctx, tknId)
}

func (s *tokenService) SearchTokens(ctx context.Context, cr *domain.TokenSearchCriteria) (*domain.TokenSearchResponse, error) {
	s.l().C(ctx).Mth("search-tkn").Dbg()
	if cr.Limit == nil {
		cr.Limit = kit.IntPtr(20)
	}
	return s.storage.SearchTokens(ctx, cr)
}

func (s *tokenService) DeleteTokensByExtId(ctx context.Context, extId domain.PartyExtId) error {
	s.l().C(ctx).Mth("del-ext").F(kit.KV{"partyId": extId.PartyId, "country": extId.CountryCode}).Dbg()
	if err := s.validateExtId(ctx, extId); err != nil {
		return err
	}
	return s.storage.DeleteTokensByExtId(ctx, extId)
}

func (s *tokenService) ValidateToken(ctx context.Context, tkn *domain.Token) error {

	if err := s.validateOcpiItem(ctx, &tkn.OcpiItem); err != nil {
		return err
	}
	if tkn.Id == "" {
		return errors.ErrTknIdEmpty(ctx)
	}
	if err := s.validateId(ctx, tkn.Id, "id"); err != nil {
		return err
	}

	// valid
	if tkn.Details.Valid == nil {
		return errors.ErrTknEmptyAttr(ctx, "token", "valid")
	}

	// type
	if tkn.Details.Type == "" {
		return errors.ErrTknEmptyAttr(ctx, "token", "type")
	}
	if _, ok := tokenTypeMap[tkn.Details.Type]; !ok {
		return errors.ErrTknInvalidAttr(ctx, "token", "type")
	}

	// contract
	if tkn.Details.ContractId == "" {
		return errors.ErrTknEmptyAttr(ctx, "token", "contract_id")
	}
	if err := s.validateId(ctx, tkn.Details.ContractId, "contract_id"); err != nil {
		return err
	}

	// visual number
	if err := s.validateMaxLen(ctx, tkn.Details.VisualNumber, 64, "contract_id"); err != nil {
		return err
	}

	// group id
	if err := s.validateMaxLen(ctx, tkn.Details.GroupId, 36, "group_id"); err != nil {
		return err
	}

	// issuer
	if tkn.Details.Issuer == "" {
		return errors.ErrTknEmptyAttr(ctx, "token", "issuer")
	}
	if err := s.validateMaxLen(ctx, tkn.Details.Issuer, 64, "issuer"); err != nil {
		return err
	}

	// while list type
	if tkn.Details.WhiteList == "" {
		return errors.ErrTknEmptyAttr(ctx, "token", "whitelist")
	}
	if _, ok := tokenWLMap[tkn.Details.WhiteList]; !ok {
		return errors.ErrTknInvalidAttr(ctx, "token", "whitelist")
	}

	// lang
	if tkn.Details.Lang != "" {
		if err := s.validateLang(ctx, "token", tkn.Details.Lang); err != nil {
			return err
		}
	}

	// profile type
	if tkn.Details.DefaultProfileType != "" {
		if _, ok := profileMap[tkn.Details.DefaultProfileType]; !ok {
			return errors.ErrTknInvalidAttr(ctx, "token", "default_profile_type")
		}
	}

	// energy contract
	if tkn.Details.EnergyContract != nil {
		if tkn.Details.EnergyContract.SupplierName == "" {
			return errors.ErrTknEmptyAttr(ctx, "energy_contract", "supplier_name")
		}
		if err := s.validateMaxLen(ctx, tkn.Details.EnergyContract.SupplierName, 64, "supplier_name"); err != nil {
			return err
		}
		if err := s.validateMaxLen(ctx, tkn.Details.EnergyContract.ContractId, 64, "contract_id"); err != nil {
			return err
		}
	}

	return nil
}

func (s *tokenService) validateAndPopulatePut(ctx context.Context, tkn, stored *domain.Token) error {

	if stored != nil {
		tkn.LastSent = stored.LastSent
		if tkn.PlatformId == "" {
			tkn.PlatformId = stored.PlatformId
		}
		if tkn.RefId == "" {
			tkn.RefId = stored.RefId
		}
		if tkn.ExtId.PartyId == "" || tkn.ExtId.CountryCode == "" {
			tkn.ExtId = stored.ExtId
		}
	}

	return s.ValidateToken(ctx, tkn)
}

func (s *tokenService) validateAndPopulateMerge(ctx context.Context, tkn, stored *domain.Token) error {

	stored.LastUpdated = tkn.LastUpdated

	if tkn.PlatformId != "" {
		stored.PlatformId = tkn.PlatformId
	}
	if tkn.RefId != "" {
		stored.RefId = tkn.RefId
	}
	if tkn.ExtId.PartyId != "" && tkn.ExtId.CountryCode != "" {
		stored.ExtId = tkn.ExtId
	}
	if tkn.Details.Type != "" {
		stored.Details.Type = tkn.Details.Type
	}
	if tkn.Details.ContractId != "" {
		stored.Details.ContractId = tkn.Details.ContractId
	}
	if tkn.Details.VisualNumber != "" {
		stored.Details.VisualNumber = tkn.Details.VisualNumber
	}
	if tkn.Details.Issuer != "" {
		stored.Details.Issuer = tkn.Details.Issuer
	}
	if tkn.Details.GroupId != "" {
		stored.Details.GroupId = tkn.Details.GroupId
	}
	if tkn.Details.Valid != nil {
		stored.Details.Valid = tkn.Details.Valid
	}
	if tkn.Details.WhiteList != "" {
		stored.Details.WhiteList = tkn.Details.WhiteList
	}
	if tkn.Details.Lang != "" {
		stored.Details.Lang = tkn.Details.Lang
	}
	if tkn.Details.DefaultProfileType != "" {
		stored.Details.DefaultProfileType = tkn.Details.DefaultProfileType
	}
	if tkn.Details.EnergyContract != nil {
		stored.Details.EnergyContract = tkn.Details.EnergyContract
	}
	return s.ValidateToken(ctx, stored)
}
