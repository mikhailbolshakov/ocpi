package ocpi

import (
	"context"
	"fmt"
	"github.com/mikhailbolshakov/kit"
	kitHttp "github.com/mikhailbolshakov/kit/http"
	service "github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/errors"
	"github.com/mikhailbolshakov/ocpi/model"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Controller struct {
	kitHttp.BaseController
}

func NewController() Controller {
	return Controller{
		BaseController: kitHttp.BaseController{Logger: service.LF()},
	}
}

func (c *Controller) ExtractToken(ctx context.Context, r *http.Request) (string, error) {
	// check and extract Authorization data
	authHeader := r.Header.Get(model.OcpiHeaderAuth)
	if authHeader == "" {
		return "", errors.ErrAuthFailed(ctx)
	}
	splitToken := strings.Split(authHeader, " ")
	if len(splitToken) < 2 || splitToken[1] == "" {
		return "", errors.ErrAuthFailed(ctx)
	}
	return splitToken[1], nil
}

func (c *Controller) ExtractHeaders(r *http.Request) *model.OcpiRequestHeader {
	return &model.OcpiRequestHeader{
		RequestId:       r.Header.Get(model.OcpiHeaderRequestId),
		CorrelationId:   r.Header.Get(model.OcpiHeaderCorrelationId),
		FromPartyId:     r.Header.Get(model.OcpiHeaderFromPartyId),
		FromCountryCode: r.Header.Get(model.OcpiHeaderFromCountryCode),
		ToPartyId:       r.Header.Get(model.OcpiHeaderToPartyId),
		ToCountryCode:   r.Header.Get(model.OcpiHeaderToCountryCode),
	}
}

func (c *Controller) EnsureCtxHeaders(ctx context.Context, attrs ...string) error {

	rqCtx, ok := kit.Request(ctx)
	if !ok || rqCtx.Kv == nil {
		return errors.ErrInvalidContext(ctx)
	}

	for _, a := range attrs {
		if _, ok := rqCtx.Kv[a]; !ok {
			return errors.ErrHeaderParamEmpty(ctx, a)
		}
	}

	return nil
}

func (c *Controller) PlatformId(ctx context.Context) (string, error) {
	appCtx, ok := kit.Request(ctx)
	if !ok {
		return "", errors.ErrAppCtxPlatformIdEmpty(ctx)
	}
	if appCtx.Kv != nil && appCtx.Kv[model.OcpiCtxPlatform] != nil {
		return appCtx.Kv[model.OcpiCtxPlatform].(string), nil
	}
	return "", errors.ErrAppCtxPlatformIdEmpty(ctx)
}

func (c *Controller) OcpiRespondOK(r *http.Request, w http.ResponseWriter, data any) {
	// set correlation header
	if r.Header.Get(model.OcpiHeaderCorrelationId) != "" {
		w.Header().Set(model.OcpiHeaderCorrelationId, r.Header.Get(model.OcpiHeaderCorrelationId))
	}
	// set total/limit headers
	appCtx, ok := kit.Request(r.Context())
	if ok {
		if v, ok := appCtx.Kv[model.OcpiHeaderTotalCount]; ok && v != nil {
			w.Header().Set(model.OcpiHeaderTotalCount, fmt.Sprintf("%d", v.(int)))
		}
		if v, ok := appCtx.Kv[model.OcpiHeaderLimit]; ok && v != nil {
			w.Header().Set(model.OcpiHeaderLimit, fmt.Sprintf("%d", v.(int)))
		}
		if v, ok := appCtx.Kv[model.OcpiHeaderLink]; ok && v != nil {
			w.Header().Set(model.OcpiHeaderLink, fmt.Sprintf("%s", v.(string)))
		}
	}
	// respond
	c.RespondJson(w, http.StatusOK, &model.OcpiResponseAny{
		OcpiResponse: model.OcpiResponse{
			StatusCode:    model.OcpiStatusCodeOk,
			StatusMessage: model.OcpiStatusMessageSuccess,
			Timestamp:     kit.Now().Format(time.RFC3339),
		},
		Data: data,
	})
}

func (c *Controller) OcpiRespondError(r *http.Request, w http.ResponseWriter, err error) {
	// set correlation
	if r.Header.Get(model.OcpiHeaderCorrelationId) != "" {
		w.Header().Set(model.OcpiHeaderCorrelationId, r.Header.Get(model.OcpiHeaderCorrelationId))
	}
	// build OCPI response
	ocpiErr := &model.OcpiResponse{
		Timestamp:     kit.Now().Format(time.RFC3339),
		StatusCode:    model.OcpiStatusGenServerError,
		StatusMessage: err.Error(),
	}
	httpStatus := http.StatusInternalServerError
	// check if this is an app error
	if appErr, ok := kit.IsAppErr(err); ok {
		if ocpiStatus, ok := appErr.Fields()[model.OcpiStatusField]; ok && ocpiStatus != nil {
			ocpiErr.StatusCode = ocpiStatus.(int)
		}
		ocpiErr.StatusMessage = appErr.Message()
		if httpSt := appErr.HttpStatus(); httpSt != nil {
			httpStatus = int(*httpSt)
		}
	}
	if c.Logger != nil {
		c.Logger().Cmp("ocpi").Pr("rest").E(err).St().Err()
	}
	c.RespondJson(w, httpStatus, ocpiErr)
}

func (c *Controller) OcpiRespondNotFoundError(r *http.Request, w http.ResponseWriter) {
	// set correlation header
	if r.Header.Get(model.OcpiHeaderCorrelationId) != "" {
		w.Header().Set(model.OcpiHeaderCorrelationId, r.Header.Get(model.OcpiHeaderCorrelationId))
	}
	// respond
	c.RespondJson(w, http.StatusNotFound, &model.OcpiResponseAny{
		OcpiResponse: model.OcpiResponse{
			StatusCode:    model.OcpiStatusCodeOk,
			StatusMessage: model.OcpiStatusMessageSuccess,
			Timestamp:     kit.Now().Format(time.RFC3339),
		},
	})
}

func (c *Controller) SetResponseCtx(ctx context.Context, rs domain.PageResponse) context.Context {
	appCtx, ok := kit.Request(ctx)
	if !ok {
		return ctx
	}
	if rs.Total != nil {
		appCtx.WithKv(model.OcpiHeaderTotalCount, *rs.Total)
	}
	if rs.Limit != nil {
		appCtx.WithKv(model.OcpiHeaderLimit, *rs.Limit)
	}
	return appCtx.ToContext(ctx)
}

func (c *Controller) SetResponseWithNextPageCtx(ctx context.Context, rs domain.PageResponse, prefixUrl string) context.Context {
	appCtx, ok := kit.Request(ctx)
	if !ok {
		return ctx
	}
	if rs.Total != nil {
		appCtx.WithKv(model.OcpiHeaderTotalCount, *rs.Total)
	}
	if rs.Limit != nil {
		appCtx.WithKv(model.OcpiHeaderLimit, *rs.Limit)
	}
	if rs.NextPage != nil {
		appCtx.WithKv(model.OcpiHeaderLink, c.buildNextPageUrl(rs.NextPage, prefixUrl))
	}
	return appCtx.ToContext(ctx)
}

func (c *Controller) buildNextPageUrl(nextPage *domain.PageRequest, prefixUrl string) string {

	v := url.Values{}

	if nextPage.DateFrom != nil {
		v.Add(model.OcpiQueryParamDateFrom, (*nextPage.DateFrom).Format(time.RFC3339))
	}
	if nextPage.DateTo != nil {
		v.Add(model.OcpiQueryParamDateTo, (*nextPage.DateTo).Format(time.RFC3339))
	}
	if nextPage.Limit != nil {
		v.Add(model.OcpiQueryParamLimit, strconv.Itoa(*nextPage.Limit))
	}
	if nextPage.Offset != nil {
		v.Add(model.OcpiQueryParamOffset, strconv.Itoa(*nextPage.Offset))
	}
	values := v.Encode()
	if values == "" {
		return prefixUrl
	}

	return fmt.Sprintf("%s?%s", prefixUrl, values)
}
