package hub

import "github.com/mikhailbolshakov/ocpi/transport/http"

func GetRoutes(c Controller) []*http.Route {
	return []*http.Route{
		// sender
		http.R("/ocpi/2.2.1/sender/hubclientinfo", c.GetHubClientInfo).GET().Auth(http.TokenB).OcpiLogging(),
		// receiver
		http.R("/ocpi/2.2.1/receiver/hubclientinfo/{country_code}/{party_id}", c.GetClientInfo).GET().Auth(http.TokenB).OcpiLogging(),
		http.R("/ocpi/2.2.1/receiver/hubclientinfo/{country_code}/{party_id}", c.PutClientInfo).PUT().Auth(http.TokenB).OcpiLogging(),
	}
}
