package sessions

import "github.com/mikhailbolshakov/ocpi/transport/http"

func GetRoutes(c Controller) []*http.Route {
	return []*http.Route{
		// sender
		http.R("/ocpi/2.2.1/sender/sessions", c.SenderGetSessions).GET().Auth(http.TokenB).OcpiLogging(),
		http.R("/ocpi/2.2.1/sender/sessions/{session_id}/preferences", c.SenderPutChargingPreferences).PUT().Auth(http.TokenB).OcpiLogging(),
		// receiver
		http.R("/ocpi/2.2.1/receiver/sessions/{country_code}/{party_id}/{session_id}", c.ReceiverGetSession).GET().Auth(http.TokenB).OcpiLogging(),
		http.R("/ocpi/2.2.1/receiver/sessions/{country_code}/{party_id}/{session_id}", c.ReceiverPostSession).PUT().Auth(http.TokenB).OcpiLogging(),
		http.R("/ocpi/2.2.1/receiver/sessions/{country_code}/{party_id}/{session_id}", c.ReceiverPatchSession).PATCH().Auth(http.TokenB).OcpiLogging(),
	}
}
