package credentials

import (
	"github.com/mikhailbolshakov/ocpi/transport/http"
)

func GetRoutes(c Controller) []*http.Route {
	return []*http.Route{
		http.R("/ocpi/2.2.1", c.GetVersionDetails).GET().Auth(http.TokenA, http.TokenB).OcpiLogging(),
		http.R("/ocpi/2.2.1/credentials", c.PostCredentials).POST().Auth(http.TokenA).OcpiLogging(),
		http.R("/ocpi/2.2.1/credentials", c.PutCredentials).PUT().Auth(http.TokenB).OcpiLogging(),
		http.R("/ocpi/2.2.1/credentials", c.GetCredentials).GET().Auth(http.TokenB).OcpiLogging(),
		http.R("/ocpi/2.2.1/credentials", c.DeleteCredentials).DELETE().Auth(http.TokenB).OcpiLogging(),
	}
}
