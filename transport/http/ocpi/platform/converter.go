package platform

import (
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/model"
)

func (c *ctrlImpl) toPlatformVersionsApi(versions domain.Versions) []*model.OcpiVersion {
	var r []*model.OcpiVersion
	for v, ep := range versions {
		r = append(r, &model.OcpiVersion{
			Version: v,
			Url:     string(ep),
		})
	}
	return r
}
