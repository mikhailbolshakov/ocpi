package tariffs

import (
	"github.com/mikhailbolshakov/ocpi/transport/http"
)

func GetRoutes(c Controller) []*http.Route {
	return []*http.Route{
		http.R("/backend/tariffs", c.PutTariff).POST().ApiKey(),
		http.R("/backend/tariffs/pull", c.PullTariffs).POST().ApiKey(),
		http.R("/backend/tariffs/{trfId}", c.GetTariff).GET().ApiKey(),
		http.R("/backend/tariffs/search/query", c.SearchTariffs).GET().ApiKey(),
	}
}
