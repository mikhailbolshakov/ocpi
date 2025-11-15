package tariffs

import "github.com/mikhailbolshakov/ocpi/transport/http"

func GetRoutes(c Controller) []*http.Route {
	return []*http.Route{
		// sender
		http.R("/ocpi/2.2.1/sender/tariffs", c.SenderGetTariffs).GET().Auth(http.TokenB).OcpiLogging(),
		// receiver
		http.R("/ocpi/2.2.1/receiver/tariffs/{country_code}/{party_id}/{tariff_id}", c.ReceiverGetTariff).GET().Auth(http.TokenB).OcpiLogging(),
		http.R("/ocpi/2.2.1/receiver/tariffs/{country_code}/{party_id}/{tariff_id}", c.ReceiverPutTariff).PUT().Auth(http.TokenB).OcpiLogging(),
	}
}
