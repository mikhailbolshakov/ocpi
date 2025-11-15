package platform

import (
	"github.com/mikhailbolshakov/ocpi/transport/http"
)

func GetRoutes(c Controller) []*http.Route {
	return []*http.Route{
		http.R("/platforms", c.PostPlatform).POST().ApiKey(),
		http.R("/platforms/{platformId}/status", c.UpdatePlatformStatus).POST().ApiKey(),
		http.R("/platforms/{platformId}", c.GetPlatform).GET().ApiKey(),
		http.R("/platforms/{platformId}/connections", c.EstablishConnection).POST().ApiKey(),
		http.R("/platforms/{platformId}/connections", c.UpdateConnection).PUT().ApiKey(),
		http.R("/platforms/tokens/generate", c.GenToken).GET().ApiKey(),
	}
}
