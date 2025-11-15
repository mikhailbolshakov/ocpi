package cdrs

import "github.com/mikhailbolshakov/ocpi/transport/http"

func GetRoutes(c Controller) []*http.Route {
	return []*http.Route{
		// sender
		http.R("/ocpi/2.2.1/sender/cdrs", c.SenderGetCdrs).GET().Auth(http.TokenB).OcpiLogging(),
		// receiver
		http.R("/ocpi/2.2.1/receiver/cdrs/{country_code}/{party_id}/{cdr_id}", c.ReceiverGetCdr).GET().Auth(http.TokenB).OcpiLogging(),
		http.R("/ocpi/2.2.1/receiver/cdrs/{country_code}/{party_id}/{cdr_id}", c.ReceiverPostCdr).POST().Auth(http.TokenB).OcpiLogging(),
	}
}
