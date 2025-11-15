package platform

import (
	kitHttp "github.com/mikhailbolshakov/kit/http"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/transport/http/ocpi"
	"net/http"
)

type Controller interface {
	kitHttp.Controller
	Health(http.ResponseWriter, *http.Request)
	GetHomeVersions(http.ResponseWriter, *http.Request)
}

type ctrlImpl struct {
	ocpi.Controller
	localPlatformService domain.LocalPlatformService
}

func NewController(localPlatformService domain.LocalPlatformService) Controller {
	return &ctrlImpl{
		localPlatformService: localPlatformService,
		Controller:           ocpi.NewController(),
	}
}

func (c *ctrlImpl) Health(w http.ResponseWriter, r *http.Request) {
	c.RespondOK(w, kitHttp.EmptyOkResponse)
}

func (c *ctrlImpl) GetHomeVersions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	platform, err := c.localPlatformService.Get(ctx)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	c.OcpiRespondOK(r, w, c.toPlatformVersionsApi(platform.VersionInfo.Available))
}
