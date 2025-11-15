package impl

import (
	"context"
	"fmt"
	iso639_3 "github.com/barbashov/iso639-3"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/errors"
	"time"
)

const (
	idMaxLen      = 36
	partyIdMaxLen = 3
	minYear       = 2020
)

var (
	authMethodMap = map[string]struct{}{
		domain.AuthMethodRequest:   {},
		domain.AuthMethodCommand:   {},
		domain.AuthMethodWhitelist: {},
	}
)

type base struct{}

func (b *base) validateMaxLen(ctx context.Context, s string, maxLen int, attr string) error {
	if len([]rune(s)) > maxLen {
		return errors.ErrMaxLenExceeded(ctx, attr)
	}
	return nil
}

func (b *base) validateId(ctx context.Context, id, attr string) error {
	return b.validateMaxLen(ctx, id, idMaxLen, attr)
}

func (b *base) validateExtId(ctx context.Context, extId domain.PartyExtId) error {
	if extId.PartyId == "" || extId.CountryCode == "" {
		return errors.ErrExtIdInvalid(ctx)
	}
	if !kit.Alfa2Valid(extId.CountryCode) {
		return errors.ErrCountryCodeInvalid(ctx)
	}
	if len([]rune(extId.PartyId)) > partyIdMaxLen {
		return errors.ErrPartyIdLen(ctx)
	}
	return nil
}

func (b *base) validateLastUpdated(ctx context.Context, lastUpdated time.Time) error {
	if lastUpdated.Year() < minYear {
		return errors.ErrLastUpdatedInvalid(ctx)
	}
	return nil
}

func (b *base) validateOcpiItem(ctx context.Context, i *domain.OcpiItem) error {
	if i.PlatformId == "" {
		return errors.ErrPlatformIdEmpty(ctx)
	}
	if err := b.validateExtId(ctx, i.ExtId); err != nil {
		return err
	}
	return b.validateLastUpdated(ctx, i.LastUpdated)
}

func (b *base) validateLang(ctx context.Context, entity string, lang string) error {
	if iso639_3.FromPart1Code(lang) == nil {
		return errors.ErrLangInvalid(ctx, entity)
	}
	return nil
}

func (b *base) validateDisplayText(ctx context.Context, entity string, dt *domain.DisplayText) error {
	if dt == nil {
		return nil
	}
	if err := b.validateLang(ctx, entity, dt.Language); err != nil {
		return errors.ErrDisplayTextInvalidAttr(ctx, entity, "language")
	}
	if dt.Text == "" {
		return errors.ErrDisplayTextInvalidAttr(ctx, entity, "text")
	}
	if err := b.validateMaxLen(ctx, dt.Text, 512, entity); err != nil {
		return errors.ErrDisplayTextInvalidAttr(ctx, entity, "language")
	}
	return nil
}

func (b *base) validateImage(ctx context.Context, entity string, img *domain.Image) error {
	if img == nil {
		return nil
	}
	if _, ok := imgCatMap[img.Category]; !ok {
		return errors.ErrImageInvalidAttr(ctx, entity, "category")
	}
	if !kit.IsUrlValid(img.Url) {
		return errors.ErrImageInvalidAttr(ctx, entity, "url")
	}
	if img.Thumbnail != "" && !kit.IsUrlValid(img.Thumbnail) {
		return errors.ErrImageInvalidAttr(ctx, entity, "thumbnail")
	}
	if img.Height < 0 || img.Width < 0 {
		return errors.ErrImageInvalidAttr(ctx, entity, "size")
	}
	return nil
}

func (b *base) validateBusinessDetails(ctx context.Context, entity string, bd *domain.BusinessDetails) error {
	if bd == nil {
		return nil
	}
	if bd.Name == "" {
		return errors.ErrBusinessDetailsEmptyAttr(ctx, entity, "name")
	}
	if err := b.validateMaxLen(ctx, bd.Name, 100, entity); err != nil {
		return errors.ErrDisplayTextInvalidAttr(ctx, entity, "language")
	}
	if bd.Website != "" && !kit.IsUrlValid(bd.Website) {
		return errors.ErrBusinessDetailsInvalidAttr(ctx, entity, "website")
	}
	return b.validateImage(ctx, entity, bd.Logo)
}

func (b *base) validateTimePeriod(ctx context.Context, entity, attr, p string) error {
	if p == "" {
		return nil
	}
	if _, err := (kit.HourMinTime{}).Parse(p); err != nil {
		return errors.ErrTimePeriodInvalid(ctx, entity, attr)
	}
	return nil
}

func (b *base) validateCdrToken(ctx context.Context, entity string, token *domain.CdrToken) error {
	if err := b.validateExtId(ctx, token.PartyExtId); err != nil {
		return errors.ErrCdrTokenInvalidAttr(ctx, entity, "ext_id")
	}
	if token.Id == "" {
		return errors.ErrCdrTokenEmptyAttr(ctx, entity, "uid")
	}
	if err := b.validateId(ctx, token.Id, "token.uid"); err != nil {
		return err
	}
	if token.Type == "" {
		return errors.ErrCdrTokenEmptyAttr(ctx, entity, "type")
	}
	if _, ok := tokenTypeMap[token.Type]; !ok {
		return errors.ErrCdrTokenInvalidAttr(ctx, entity, "type")
	}
	if token.ContractId == "" {
		return errors.ErrCdrTokenEmptyAttr(ctx, entity, "contract_id")
	}
	if err := b.validateId(ctx, token.ContractId, "token.contract_id"); err != nil {
		return err
	}
	return nil
}

func (b *base) validateAuth(ctx context.Context, entity, authMethod, authRef string) error {
	// auth method
	if authMethod == "" {
		return errors.ErrAuthEmptyAttr(ctx, entity, "auth_method")
	}
	if _, ok := authMethodMap[authMethod]; !ok {
		return errors.ErrAuthInvalidAttr(ctx, entity, "auth_method")
	}
	if err := b.validateMaxLen(ctx, authRef, idMaxLen, "authorization_reference"); err != nil {
		return err
	}
	return nil
}

func (b *base) validateChargingPeriods(ctx context.Context, entity string, chargingPeriods []*domain.ChargingPeriod) error {
	for _, cp := range chargingPeriods {
		if err := b.validateId(ctx, cp.TariffId, "charging_period.tariff_id"); err != nil {
			return err
		}
		if len(cp.Dimensions) == 0 {
			return errors.ErrChargingPeriodEmptyAttr(ctx, entity, "dimensions")
		}
		for _, d := range cp.Dimensions {
			if d.Type == "" {
				return errors.ErrChargingPeriodEmptyAttr(ctx, entity, "dimension.type")
			}
			if _, ok := dimensionMap[d.Type]; !ok {
				return errors.ErrChargingPeriodInvalidAttr(ctx, entity, "dimension.type")
			}
		}
		// TODO: signed data related attributes isn't used and not validated
		if cp.Url != "" && !kit.IsUrlValid(cp.Url) {
			return errors.ErrChargingPeriodInvalidAttr(ctx, entity, "url")
		}
	}
	return nil
}

func (b *base) validatePrice(ctx context.Context, entity, attr string, price *domain.Price) error {
	if price == nil {
		return nil
	}
	if price.ExclVat < 0 {
		return errors.ErrPriceInvalidAttr(ctx, entity, fmt.Sprintf("%s.%s", attr, "excl_vat"))
	}
	if price.InclVat != nil && *price.InclVat < 0 {
		return errors.ErrPriceInvalidAttr(ctx, entity, fmt.Sprintf("%s.%s", attr, "incl_vat"))
	}
	return nil
}

func (b *base) validateAmount(ctx context.Context, entity, attr string, v *float64) error {
	if v == nil {
		return nil
	}
	if *v < 0 {
		return errors.ErrInvalidAmountAttr(ctx, entity, attr)
	}
	return nil
}
