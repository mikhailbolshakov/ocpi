package ocpi

import (
	"fmt"
	"github.com/mikhailbolshakov/ocpi/model"
	"net/http"
	"time"
)

type rqBuilder struct {
	rq        *restRequest
	urlParams []string
}

func newRq(url, token, ev string) *rqBuilder {
	return &rqBuilder{
		rq: &restRequest{
			Url:      url,
			Verb:     http.MethodGet,
			Token:    token,
			LogEvent: ev,
		},
	}
}

func (r *rqBuilder) Body(b any) *rqBuilder {
	r.rq.Body = b
	return r
}

func (r *rqBuilder) Verb(v string) *rqBuilder {
	r.rq.Verb = v
	return r
}

func (r *rqBuilder) ResponseModel(m any) *rqBuilder {
	r.rq.RespModel = m
	return r
}

func (r *rqBuilder) Header(k, v string) *rqBuilder {
	if r.rq.Header == nil {
		r.rq.Header = map[string]string{}
	}
	r.rq.Header[k] = v
	return r
}

func (r *rqBuilder) QueryParam(k, v string) *rqBuilder {
	if r.rq.QueryParams == nil {
		r.rq.QueryParams = map[string]string{}
	}
	r.rq.QueryParams[k] = v
	return r
}

func (r *rqBuilder) UrlParams(p ...string) *rqBuilder {
	r.urlParams = append(r.urlParams, p...)
	return r
}

func (r *rqBuilder) Page(rq *model.OcpiGetPageRequest) *rqBuilder {
	if rq.Limit != nil {
		r.QueryParam(model.OcpiQueryParamLimit, fmt.Sprintf("%d", *rq.Limit))
	}
	if rq.Offset != nil {
		r.QueryParam(model.OcpiQueryParamOffset, fmt.Sprintf("%d", *rq.Offset))
	}
	if rq.DateFrom != nil {
		r.QueryParam(model.OcpiQueryParamDateFrom, rq.DateFrom.Format(time.RFC3339))
	}
	if rq.DateTo != nil {
		r.QueryParam(model.OcpiQueryParamDateTo, rq.DateTo.Format(time.RFC3339))
	}
	return r
}

func (r *rqBuilder) B() *restRequest {
	for _, p := range r.urlParams {
		r.rq.Url += "/" + p
	}
	return r.rq
}
