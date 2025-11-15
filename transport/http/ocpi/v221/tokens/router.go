package tokens

import "github.com/mikhailbolshakov/ocpi/transport/http"

func GetRoutes(c Controller) []*http.Route {
	return []*http.Route{
		// sender
		http.R("/ocpi/2.2.1/sender/tokens", c.SenderGetTokens).GET().Auth(http.TokenB).OcpiLogging(),
		http.R("/ocpi/2.2.1/sender/tokens/{token_id}/authorize", c.SenderAuthToken).POST().Auth(http.TokenB).OcpiLogging(),
		// receiver
		http.R("/ocpi/2.2.1/receiver/tokens/{country_code}/{party_id}/{token_id}", c.ReceiverGetToken).GET().Auth(http.TokenB).OcpiLogging(),
		http.R("/ocpi/2.2.1/receiver/tokens/{country_code}/{party_id}/{token_id}", c.ReceiverPutToken).PUT().Auth(http.TokenB).OcpiLogging(),
		http.R("/ocpi/2.2.1/receiver/tokens/{country_code}/{party_id}/{token_id}", c.ReceiverPatchToken).PATCH().Auth(http.TokenB).OcpiLogging(),
	}
}
