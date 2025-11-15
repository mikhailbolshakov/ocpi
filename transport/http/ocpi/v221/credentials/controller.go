package credentials

import (
	kitHttp "github.com/mikhailbolshakov/kit/http"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/model"
	"github.com/mikhailbolshakov/ocpi/transport/http/ocpi"
	"github.com/mikhailbolshakov/ocpi/usecase"
	"net/http"
)

type Controller interface {
	kitHttp.Controller
	GetVersionDetails(http.ResponseWriter, *http.Request)
	PostCredentials(http.ResponseWriter, *http.Request)
	PutCredentials(http.ResponseWriter, *http.Request)
	DeleteCredentials(http.ResponseWriter, *http.Request)
	GetCredentials(http.ResponseWriter, *http.Request)
}

type ctrlImpl struct {
	ocpi.Controller
	credentialUc    usecase.CredentialsUc
	localPlatform   domain.LocalPlatformService
	platformService domain.PlatformService
}

func NewController(credentialUc usecase.CredentialsUc, platformService domain.PlatformService, localPlatform domain.LocalPlatformService) Controller {
	return &ctrlImpl{
		Controller:      ocpi.NewController(),
		credentialUc:    credentialUc,
		platformService: platformService,
		localPlatform:   localPlatform,
	}
}

func (c *ctrlImpl) GetVersionDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rs := &model.OcpiVersionDetails{
		Version:   "2.2.1",
		Endpoints: c.toVersionEndpointsApi(c.localPlatform.GetEndpoints(ctx, "2.2.1")),
	}
	c.OcpiRespondOK(r, w, rs)
}

func (c *ctrlImpl) PostCredentials(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	platformId, err := c.PlatformId(ctx)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	rq, err := kitHttp.DecodeRequest[model.OcpiCredentials](ctx, r)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	rs, err := c.credentialUc.AcceptConnection(ctx, platformId, rq)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	c.OcpiRespondOK(r, w, rs)
}

func (c *ctrlImpl) PutCredentials(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	platformId, err := c.PlatformId(ctx)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	rq, err := kitHttp.DecodeRequest[model.OcpiCredentials](ctx, r)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	rs, err := c.credentialUc.AcceptConnection(ctx, platformId, rq)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	c.OcpiRespondOK(r, w, rs)
}

func (c *ctrlImpl) DeleteCredentials(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	platformId, err := c.PlatformId(ctx)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	err = c.credentialUc.OnRemoteDeleteCredentials(ctx, platformId)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	c.OcpiRespondOK(r, w, nil)
}

func (c *ctrlImpl) GetCredentials(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	platformId, err := c.PlatformId(ctx)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	rs, err := c.credentialUc.OnRemoteGetCredentials(ctx, platformId)
	if err != nil {
		c.OcpiRespondError(r, w, err)
		return
	}

	c.OcpiRespondOK(r, w, rs)
}
