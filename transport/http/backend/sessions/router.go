package sessions

import (
	"github.com/mikhailbolshakov/ocpi/transport/http"
)

func GetRoutes(c Controller) []*http.Route {
	return []*http.Route{
		http.R("/backend/sessions", c.PutSession).POST().ApiKey(),
		http.R("/backend/sessions/{sessId}", c.PatchSession).PATCH().ApiKey(),
		http.R("/backend/sessions/{sessId}", c.GetSession).GET().ApiKey(),
		http.R("/backend/sessions/search/query", c.SearchSessions).GET().ApiKey(),
	}
}
