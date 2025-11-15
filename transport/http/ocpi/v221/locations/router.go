package locations

import "github.com/mikhailbolshakov/ocpi/transport/http"

func GetRoutes(c Controller) []*http.Route {
	return []*http.Route{
		// sender
		http.R("/ocpi/2.2.1/sender/locations", c.SenderGetLocations).GET().Auth(http.TokenB).OcpiLogging(),
		http.R("/ocpi/2.2.1/sender/locations/{location_id}", c.SenderGetLocation).GET().Auth(http.TokenB).OcpiLogging(),
		http.R("/ocpi/2.2.1/sender/locations/{location_id}/{evse_uid}", c.SenderGetEvse).GET().Auth(http.TokenB).OcpiLogging(),
		http.R("/ocpi/2.2.1/sender/locations/{location_id}/{evse_uid}/{connector_id}", c.SenderGetConnector).GET().Auth(http.TokenB).OcpiLogging(),

		// receiver
		http.R("/ocpi/2.2.1/receiver/locations/{country_code}/{party_id}/{location_id}", c.ReceiverGetLocation).GET().Auth(http.TokenB).OcpiLogging(),
		http.R("/ocpi/2.2.1/receiver/locations/{country_code}/{party_id}/{location_id}/{evse_uid}", c.ReceiverGetEvse).GET().Auth(http.TokenB).OcpiLogging(),
		http.R("/ocpi/2.2.1/receiver/locations/{country_code}/{party_id}/{location_id}/{evse_uid}/{connector_id}", c.ReceiverGetConnector).GET().Auth(http.TokenB).OcpiLogging(),
		http.R("/ocpi/2.2.1/receiver/locations/{country_code}/{party_id}/{location_id}", c.ReceiverPutLocation).PUT().Auth(http.TokenB).OcpiLogging(),
		http.R("/ocpi/2.2.1/receiver/locations/{country_code}/{party_id}/{location_id}/{evse_uid}", c.ReceiverPutEvse).PUT().Auth(http.TokenB).OcpiLogging(),
		http.R("/ocpi/2.2.1/receiver/locations/{country_code}/{party_id}/{location_id}/{evse_uid}/{connector_id}", c.ReceiverPutConnector).PUT().Auth(http.TokenB).OcpiLogging(),
		http.R("/ocpi/2.2.1/receiver/locations/{country_code}/{party_id}/{location_id}", c.ReceiverPatchLocation).PATCH().Auth(http.TokenB).OcpiLogging(),
		http.R("/ocpi/2.2.1/receiver/locations/{country_code}/{party_id}/{location_id}/{evse_uid}", c.ReceiverPatchEvse).PATCH().Auth(http.TokenB).OcpiLogging(),
		http.R("/ocpi/2.2.1/receiver/locations/{country_code}/{party_id}/{location_id}/{evse_uid}/{connector_id}", c.ReceiverPatchConnector).PATCH().Auth(http.TokenB).OcpiLogging(),
	}
}
