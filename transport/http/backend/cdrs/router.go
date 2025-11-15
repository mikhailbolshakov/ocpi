package cdrs

import (
	"github.com/mikhailbolshakov/ocpi/transport/http"
)

func GetRoutes(c Controller) []*http.Route {
	return []*http.Route{
		http.R("/backend/cdrs", c.PostCdr).POST().ApiKey(),
		http.R("/backend/cdrs/{cdrId}", c.GetCdr).GET().ApiKey(),
		http.R("/backend/cdrs/search/query", c.SearchCdrs).GET().ApiKey(),
	}
}
