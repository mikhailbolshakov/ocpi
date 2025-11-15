package impl

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/model"
	"github.com/mikhailbolshakov/ocpi/usecase"
	"github.com/stretchr/testify/suite"
	"testing"
)

type pageReaderTestSuite struct {
	kit.Suite
}

func (s *pageReaderTestSuite) SetupSuite() {
	s.Suite.Init(ocpi.LF())
}

func (s *pageReaderTestSuite) SetupTest() {
}

func (s *pageReaderTestSuite) TearDownSuite() {}

func TestPageReaderSuite(t *testing.T) {
	suite.Run(t, new(pageReaderTestSuite))
}

func (s *pageReaderTestSuite) Test_Locations() {

	test := func(rs [][]*model.OcpiLocation, expected int) {
		i := 0
		readFn := func(context.Context, *usecase.OcpiRepositoryPagingRequest) ([]*model.OcpiLocation, error) {
			for i < len(rs) {
				r := rs[i]
				i++
				return r, nil
			}
			return nil, nil
		}
		var res []*model.OcpiLocation
		pr := NewPageReader(readFn)
		ch := pr.GetPage(s.Ctx, buildOcpiRepositoryRequest("http://ep", domain.PlatformToken(kit.NewRandString()), &domain.Platform{}, &domain.Platform{}), 10, nil, nil)
		func() {
			for v := range ch {
				res = append(res, v...)
			}
		}()
		s.Len(res, expected)
	}

	test(nil, 0)
	test([][]*model.OcpiLocation{{}, {}}, 0)
	test([][]*model.OcpiLocation{{{Id: "1"}}, {{Id: "2"}}}, 2)
	test([][]*model.OcpiLocation{{{Id: "1"}, {Id: "1.1"}}, {{Id: "2"}, {Id: "2.2"}}}, 4)
	test([][]*model.OcpiLocation{{{Id: "1"}, {Id: "1.1"}, {Id: "1.3"}}, {{Id: "2"}, {Id: "2.2"}}}, 5)
}

func (s *pageReaderTestSuite) Test_Tariffs() {

	test := func(rs [][]*model.OcpiTariff, expected int) {
		i := 0
		readFn := func(context.Context, *usecase.OcpiRepositoryPagingRequest) ([]*model.OcpiTariff, error) {
			for i < len(rs) {
				r := rs[i]
				i++
				return r, nil
			}
			return nil, nil
		}
		var res []*model.OcpiTariff
		pr := NewPageReader(readFn)
		ch := pr.GetPage(s.Ctx, buildOcpiRepositoryRequest("http://ep", domain.PlatformToken(kit.NewRandString()), &domain.Platform{}, &domain.Platform{}), 10, nil, nil)
		func() {
			for v := range ch {
				res = append(res, v...)
			}
		}()
		s.Len(res, expected)
	}

	test(nil, 0)
	test([][]*model.OcpiTariff{{}, {}}, 0)
	test([][]*model.OcpiTariff{{{Id: "1"}}, {{Id: "2"}}}, 2)
	test([][]*model.OcpiTariff{{{Id: "1"}, {Id: "1.1"}}, {{Id: "2"}, {Id: "2.2"}}}, 4)
	test([][]*model.OcpiTariff{{{Id: "1"}, {Id: "1.1"}, {Id: "1.3"}}, {{Id: "2"}, {Id: "2.2"}}}, 5)
}
