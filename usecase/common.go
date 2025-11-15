package usecase

import (
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/model"
)

type ErrorHandler func(err error)

type OcpiRepositoryBaseRequest struct {
	Endpoint       domain.Endpoint
	Token          domain.PlatformToken
	FromPlatformId string
	ToPlatformId   string
}

type OcpiRepositoryRequestG[T any] struct {
	OcpiRepositoryBaseRequest
	Request T
}

type OcpiRepositoryIdRequest struct {
	OcpiRepositoryBaseRequest
	Id string
}

type OcpiRepositoryErrHandlerRequest struct {
	OcpiRepositoryBaseRequest
	Handler ErrorHandler
}

type OcpiRepositoryErrHandlerRequestG[T any] struct {
	OcpiRepositoryErrHandlerRequest
	Request T
}

type OcpiRepositoryPagingRequest struct {
	OcpiRepositoryBaseRequest
	model.OcpiGetPageRequest
}
