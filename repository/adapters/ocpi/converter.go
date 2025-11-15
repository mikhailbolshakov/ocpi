package ocpi

import (
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/model"
)

func (a *adapterImpl) toVersionsDomain(versions []*model.OcpiVersion) domain.Versions {
	res := make(domain.Versions)
	for _, v := range versions {
		res[v.Version] = domain.Endpoint(v.Url)
	}
	return res
}

func (a *adapterImpl) toVersionDetailsDomain(versionDetails *model.OcpiVersionDetails) domain.ModuleEndpoints {
	if versionDetails == nil {
		return nil
	}
	r := make(domain.ModuleEndpoints)
	for _, vd := range versionDetails.Endpoints {
		if r[vd.Id] == nil {
			r[vd.Id] = make(map[string]domain.Endpoint)
		}
		r[vd.Id][vd.Role] = domain.Endpoint(vd.Url)
	}
	return r
}
