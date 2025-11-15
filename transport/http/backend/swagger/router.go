package swagger

import (
	_ "github.com/mikhailbolshakov/ocpi/swagger"
	"github.com/mikhailbolshakov/ocpi/transport/http"
	httpSwagger "github.com/swaggo/http-swagger"
)

func GetRoutes() []*http.Route {
	return []*http.Route{
		http.R("/swagger/*", nil).Handler(httpSwagger.Handler(httpSwagger.InstanceName("ocpi"))).GET(),
	}
}
