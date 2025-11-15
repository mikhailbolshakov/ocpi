package tokens

import (
	"github.com/mikhailbolshakov/ocpi/transport/http"
)

func GetRoutes(c Controller) []*http.Route {
	return []*http.Route{
		http.R("/backend/tokens", c.PutToken).POST().ApiKey(),
		http.R("/backend/tokens/pull", c.PullTokens).POST().ApiKey(),
		http.R("/backend/tokens/{tknId}", c.GetToken).GET().ApiKey(),
		http.R("/backend/tokens/search/query", c.SearchTokens).GET().ApiKey(),
	}
}
