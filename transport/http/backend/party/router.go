package party

import (
	"github.com/mikhailbolshakov/ocpi/transport/http"
)

func GetRoutes(c Controller) []*http.Route {
	return []*http.Route{
		// backend parties
		http.R("/backend/parties", c.PutParty).POST().ApiKey(),
		http.R("/backend/parties/pull", c.PullParties).POST().ApiKey(),
		http.R("/backend/parties/{partyId}", c.GetParty).GET().ApiKey(),
		http.R("/backend/parties/search/query", c.SearchParties).GET().ApiKey(),
	}
}
