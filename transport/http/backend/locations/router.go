package locations

import (
	"github.com/mikhailbolshakov/ocpi/transport/http"
)

func GetRoutes(c Controller) []*http.Route {
	return []*http.Route{
		http.R("/backend/locations", c.PutLocation).POST().ApiKey(),
		http.R("/backend/locations/{locId}", c.GetLocation).GET().ApiKey(),
		http.R("/backend/locations/search/query", c.SearchLocations).GET().ApiKey(),
		http.R("/backend/locations/pull", c.PullLocations).POST().ApiKey(),

		http.R("/backend/locations/{locId}/evses", c.PutEvse).POST().ApiKey(),
		http.R("/backend/locations/{locId}/evses/{evseId}/status", c.SetEvseStatus).POST().ApiKey(),
		http.R("/backend/locations/{locId}/evses/{evseId}", c.GetEvse).GET().ApiKey(),
		http.R("/backend/evses/search/query", c.SearchEvses).GET().ApiKey(),

		http.R("/backend/locations/{locId}/evses/{evseId}/connectors", c.PutConnector).POST().ApiKey(),
		http.R("/backend/locations/{locId}/evses/{evseId}/connectors/{conId}", c.GetConnector).GET().ApiKey(),
		http.R("/backend/connectors/search/query", c.SearchConnectors).GET().ApiKey(),
	}
}
