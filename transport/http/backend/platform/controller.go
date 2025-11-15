package platform

import (
	kitHttp "github.com/mikhailbolshakov/kit/http"
	service "github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/backend"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/usecase"
	"net/http"
)

type Controller interface {
	kitHttp.Controller
	PostPlatform(http.ResponseWriter, *http.Request)
	UpdatePlatformStatus(http.ResponseWriter, *http.Request)
	GetPlatform(http.ResponseWriter, *http.Request)
	EstablishConnection(http.ResponseWriter, *http.Request)
	UpdateConnection(http.ResponseWriter, *http.Request)
	GenToken(http.ResponseWriter, *http.Request)
}

type ctrlImpl struct {
	kitHttp.BaseController
	platformService domain.PlatformService
	credentialsUc   usecase.CredentialsUc
	converter       usecase.CredentialsConverter
	genToken        domain.TokenGenerator
}

func NewController(platformService domain.PlatformService, credentialsUc usecase.CredentialsUc, converter usecase.CredentialsConverter, genToken domain.TokenGenerator) Controller {
	return &ctrlImpl{
		platformService: platformService,
		credentialsUc:   credentialsUc,
		converter:       converter,
		genToken:        genToken,
		BaseController:  kitHttp.BaseController{Logger: service.LF()},
	}
}

// PostPlatform godoc
// @Summary creates or updates platform details
// @Param request body backend.PlatformRequest true "platform request"
// @Accept json
// @Success 200 {object} backend.Platform
// @Failure 500 {object} http.Error
// @Router /platforms [post]
// @tags platform
func (c *ctrlImpl) PostPlatform(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq, err := kitHttp.DecodeRequest[backend.PlatformRequest](ctx, r)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	platform, err := c.platformService.Merge(ctx, c.converter.PlatformBackendToDomain(rq))
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.converter.PlatformDomainToBackend(platform))
}

// UpdatePlatformStatus godoc
// @Summary updates platform status
// @Param platformId path string true "platform ID"
// @Param status path string true "platform status"
// @Accept json
// @Success 200 {object} backend.Platform
// @Failure 500 {object} http.Error
// @Router /platforms/{platformId}/status [post]
// @tags platform
func (c *ctrlImpl) UpdatePlatformStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	platformId, err := c.Var(ctx, r, "platformId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	status, err := c.FormVal(ctx, r, "status", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	platform, err := c.platformService.SetStatus(ctx, platformId, status)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.converter.PlatformDomainToBackend(platform))
}

// GetPlatform godoc
// @Summary retrieves a platform object by id
// @Accept json
// @Param platformId path string true "platform ID"
// @Success 200 {object} backend.Platform
// @Failure 500 {object} http.Error
// @Router /platforms/{platformId} [get]
// @tags platform
func (c *ctrlImpl) GetPlatform(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	platformId, err := c.Var(ctx, r, "platformId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	platform, err := c.platformService.Get(ctx, platformId)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.converter.PlatformDomainToBackend(platform))
}

// EstablishConnection godoc
// @Summary establishes a connection with a remote platform as a sender
// @Param platformId path string true "platform ID to connect to"
// @Accept json
// @Success 200 {object} backend.Platform
// @Failure 500 {object} http.Error
// @Router /platforms/{platformId}/connections [post]
// @tags platform
func (c *ctrlImpl) EstablishConnection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	platformId, err := c.Var(ctx, r, "platformId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	platform, err := c.credentialsUc.EstablishConnection(ctx, platformId)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.converter.PlatformDomainToBackend(platform))
}

// UpdateConnection godoc
// @Summary updates an existent connection with a remote platform
// @Param platformId path string true "platform ID to connect to"
// @Accept json
// @Success 200 {object} backend.Platform
// @Failure 500 {object} http.Error
// @Router /platforms/{platformId}/connections [put]
// @tags platform
func (c *ctrlImpl) UpdateConnection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	platformId, err := c.Var(ctx, r, "platformId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	platform, err := c.credentialsUc.UpdateConnection(ctx, platformId)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.converter.PlatformDomainToBackend(platform))
}

// GenToken godoc
// @Summary generates a new platform token
// @Accept json
// @Success 200
// @Failure 500 {object} http.Error
// @Router /platforms/tokens/generate [get]
// @tags platform
func (c *ctrlImpl) GenToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token, err := c.genToken.Generate(ctx)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, token)
}
