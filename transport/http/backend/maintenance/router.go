package maintenance

import (
	"github.com/mikhailbolshakov/ocpi/transport/http"
)

func GetRoutes(c Controller) []*http.Route {
	return []*http.Route{
		http.R("/maintenance/parties/{partyId}/{countryCode}", c.DeleteParty).DELETE().ApiKey(),
		http.R("/maintenance/logs", c.SearchLogEntries).GET().ApiKey(),
	}
}
