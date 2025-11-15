package webhook

import (
	"github.com/mikhailbolshakov/ocpi/transport/http"
)

func GetRoutes(c Controller) []*http.Route {
	return []*http.Route{
		http.R("/webhooks", c.CreateUpdateWebhook).POST().ApiKey(),
		http.R("/webhooks/{whId}", c.GetWebhook).GET().ApiKey(),
		http.R("/webhooks/{whId}", c.DeleteWebhook).DELETE().ApiKey(),
		http.R("/webhooks/search/query", c.SearchWebhooks).GET().ApiKey(),
	}
}
