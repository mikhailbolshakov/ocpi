package impl

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/errors"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type baseTestSuite struct {
	kit.Suite
	svc *base
}

func (s *baseTestSuite) SetupSuite() {
	s.Suite.Init(ocpi.LF())
}

func (s *baseTestSuite) SetupTest() {
	s.svc = &base{}
}

func (s *baseTestSuite) TearDownSuite() {}

func TestBaseSuite(t *testing.T) {
	suite.Run(t, new(baseTestSuite))
}

func (s *baseTestSuite) Test_ValidateOcpiItem() {
	// valid
	oi := s.ocpiItem()
	s.NoError(s.svc.validateOcpiItem(s.Ctx, oi))

	// empty platform
	oi = s.ocpiItem()
	oi.PlatformId = ""
	s.AssertAppErr(s.svc.validateOcpiItem(s.Ctx, oi), errors.ErrCodePlatformIdEmpty)

	// empty country code
	oi = s.ocpiItem()
	oi.ExtId.CountryCode = ""
	s.AssertAppErr(s.svc.validateOcpiItem(s.Ctx, oi), errors.ErrCodeExtIdInvalid)

	// invalid country code
	oi = s.ocpiItem()
	oi.ExtId.CountryCode = "XXX"
	s.AssertAppErr(s.svc.validateOcpiItem(s.Ctx, oi), errors.ErrCodeCountryCodeInvalid)

	// empty party
	oi = s.ocpiItem()
	oi.ExtId.PartyId = ""
	s.AssertAppErr(s.svc.validateOcpiItem(s.Ctx, oi), errors.ErrCodeExtIdInvalid)

	// invalid last_updated
	oi = s.ocpiItem()
	oi.LastUpdated = time.Time{}
	s.AssertAppErr(s.svc.validateOcpiItem(s.Ctx, oi), errors.ErrCodeLastUpdatedInvalid)

}

func (s *baseTestSuite) Test_ValidateDisplayText() {
	// valid
	dt := s.displayText()
	s.NoError(s.svc.validateDisplayText(s.Ctx, "", dt))

	// invalid lang
	dt = s.displayText()
	dt.Language = "invalid"
	s.Error(s.svc.validateDisplayText(s.Ctx, "", dt))

	// empty lang
	dt = s.displayText()
	dt.Language = ""
	s.Error(s.svc.validateDisplayText(s.Ctx, "", dt))

	// empty text
	dt = s.displayText()
	dt.Text = ""
	s.Error(s.svc.validateDisplayText(s.Ctx, "", dt))
}

func (s *baseTestSuite) Test_ValidateImage() {

	// valid
	s.Run("valid", func() {
		img := s.image()
		s.NoError(s.svc.validateImage(s.Ctx, "", img))
	})

	// invalid url
	s.Run("invalid url", func() {
		img := s.image()
		img.Url = "invalid"
		s.Error(s.svc.validateImage(s.Ctx, "", img))
	})

	// invalid url
	s.Run("invalid thumbnail", func() {
		img := s.image()
		img.Thumbnail = "invalid"
		s.Error(s.svc.validateImage(s.Ctx, "", img))
	})

	// invalid category
	s.Run("invalid category", func() {
		img := s.image()
		img.Category = "invalid"
		s.Error(s.svc.validateImage(s.Ctx, "", img))
	})

	s.Run("invalid size", func() {
		img := s.image()
		img.Height = -100
		s.Error(s.svc.validateImage(s.Ctx, "", img))
	})

}

func (s *baseTestSuite) ocpiItem() *domain.OcpiItem {
	return &domain.OcpiItem{
		ExtId: domain.PartyExtId{
			PartyId:     "PPP",
			CountryCode: "RS",
		},
		PlatformId:  kit.NewRandString(),
		RefId:       "",
		LastUpdated: kit.Now(),
		LastSent:    nil,
	}
}

func (s *baseTestSuite) displayText() *domain.DisplayText {
	return &domain.DisplayText{
		Language: "en",
		Text:     "text",
	}
}

func (s *baseTestSuite) image() *domain.Image {
	return &domain.Image{
		Url:       "https://test.com/image",
		Thumbnail: "https://test.com/tn",
		Category:  domain.ImageCategoryOperator,
		Type:      "jpg",
		Width:     100,
		Height:    100,
	}
}
