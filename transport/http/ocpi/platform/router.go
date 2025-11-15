package platform

import (
	"github.com/mikhailbolshakov/ocpi/transport/http"
)

func GetRoutes(c Controller) []*http.Route {
	return []*http.Route{
		http.R("/ocpi/health", c.Health).GET().NoAuth(),
		http.R("/ocpi/versions", c.GetHomeVersions).Auth(http.TokenA, http.TokenB).GET(),
	}
}
