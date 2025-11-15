package credentials

import (
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/model"
)

func (c *ctrlImpl) toVersionEndpointsApi(endpoints domain.ModuleEndpoints) []model.OcpiVersionModuleEndpoint {
	var r []model.OcpiVersionModuleEndpoint
	for moduleId, roles := range endpoints {
		for role, ep := range roles {
			r = append(r, model.OcpiVersionModuleEndpoint{
				Id:   moduleId,
				Role: role,
				Url:  string(ep),
			})
		}
	}
	return r
}
