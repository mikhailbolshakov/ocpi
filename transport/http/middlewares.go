package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mikhailbolshakov/kit"
	ocpiCfg "github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/domain"
	"github.com/mikhailbolshakov/ocpi/errors"
	"github.com/mikhailbolshakov/ocpi/model"
	"github.com/mikhailbolshakov/ocpi/transport/http/ocpi"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	HeaderXRealIp       = "x-real-ip"
	HeaderXForwarderFor = "x-forwarder-for"
	HeaderApiKey        = "x-api-key"
)

type Middleware struct {
	ocpi.Controller
	platformService domain.PlatformService
	ocpiLogging     domain.OcpiLogService
	cfg             *ocpiCfg.CfgOcpiConfig
}

func NewMiddleware(platformService domain.PlatformService, ocpiLogging domain.OcpiLogService, cfg *ocpiCfg.CfgOcpiConfig) *Middleware {
	return &Middleware{
		Controller:      ocpi.NewController(),
		platformService: platformService,
		ocpiLogging:     ocpiLogging,
		cfg:             cfg,
	}
}

func (m *Middleware) AuthAccessTokenMiddleware(next http.HandlerFunc, tokenTypes ...string) http.HandlerFunc {

	f := func(w http.ResponseWriter, r *http.Request) {

		// check if context is a request context
		ctxRq, err := kit.MustRequest(r.Context())
		if err != nil {
			m.OcpiRespondError(r, w, err)
			return
		}
		ctx := r.Context()

		// extract token
		token, err := m.ExtractToken(ctx, r)
		if err != nil {
			m.OcpiRespondError(r, w, errors.ErrAuthFailed(ctx))
			return
		}

		// try to find platform by the token
		var platform *domain.Platform
		for _, tokenType := range tokenTypes {
			switch tokenType {
			case TokenA:
				platform, err = m.platformService.GetByTokenA(ctx, domain.PlatformToken(token))
			case TokenB:
				platform, err = m.platformService.GetByTokenB(ctx, domain.PlatformToken(token))
			case TokenC:
				platform, err = m.platformService.GetByTokenC(ctx, domain.PlatformToken(token))
			}
			if err != nil {
				m.OcpiRespondError(r, w, errors.ErrAuthFailed(ctx))
				return
			}
			if platform != nil {
				break
			}
		}

		// no platform found
		if platform == nil {
			m.OcpiRespondError(r, w, errors.ErrAuthFailed(ctx))
			return
		}

		if platform.Status == domain.ConnectionStatusSuspended {
			m.OcpiRespondError(r, w, errors.ErrPlatformNotAvailable(ctx))
			return
		}

		// populate context
		r = r.WithContext(ctxRq.WithKv(model.OcpiCtxPlatform, platform.Id).ToContext(r.Context()))

		next.ServeHTTP(w, r)
	}

	return f
}

func (m *Middleware) SetContextMiddleware(next http.Handler) http.Handler {

	f := func(w http.ResponseWriter, r *http.Request) {

		// init context
		ctxRq := kit.NewRequestCtx().Rest()

		// extract OCPI standard headers
		headers := m.ExtractHeaders(r)

		// set request ID if specified
		if headers.RequestId != "" {
			ctxRq = ctxRq.WithRequestId(headers.RequestId)
		} else {
			ctxRq = ctxRq.WithNewRequestId()
		}

		// set OCPI headers to context
		if headers.CorrelationId != "" {
			ctxRq = ctxRq.WithKv(model.OcpiCtxCorrelationId, headers.CorrelationId)
		}
		if headers.FromCountryCode != "" {
			ctxRq = ctxRq.WithKv(model.OcpiCtxFromCountryCode, headers.FromCountryCode)
		}
		if headers.FromPartyId != "" {
			ctxRq = ctxRq.WithKv(model.OcpiCtxFromParty, headers.FromPartyId)
		}
		if headers.ToCountryCode != "" {
			ctxRq = ctxRq.WithKv(model.OcpiCtxToCountryCode, headers.ToCountryCode)
		}
		if headers.ToPartyId != "" {
			ctxRq = ctxRq.WithKv(model.OcpiCtxToParty, headers.ToPartyId)
		}

		// set client ip header coming from client
		clientIP := r.Header.Get(HeaderXRealIp)
		// try to get ip from x-forwarder-for
		if clientIP == "" {
			clientIP = r.Header.Get(HeaderXForwarderFor)
		}
		if clientIP != "" {
			ctxRq = ctxRq.WithClientIp(clientIP)
		}

		ctx := ctxRq.ToContext(r.Context())

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}

func (m *Middleware) ApiKeyMiddleware(next http.Handler) http.HandlerFunc {

	f := func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		if m.cfg.Local.ApiKey != "" {
			// api key header coming from client
			apiKey := r.Header.Get(HeaderApiKey)

			// check api-key
			if apiKey == "" || apiKey != m.cfg.Local.ApiKey {
				m.RespondError(w, errors.ErrAuthBackendFailed(ctx))
				return
			}
		}

		next.ServeHTTP(w, r)
	}

	return f
}

func (m *Middleware) WithTimeoutMiddleware(next http.HandlerFunc, timeoutSec int) http.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request) {
		timeoutHandler := http.TimeoutHandler(next, time.Duration(timeoutSec)*time.Second, "")
		timeoutHandler.ServeHTTP(w, r)
	}

	return f
}

type writerWrapper struct {
	http.ResponseWriter
	status int
	body   []byte
}

func (r *writerWrapper) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.body = append(r.body, b...)
	return size, err
}

func (r *writerWrapper) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode) // write status code using original http.ResponseWriter
	r.status = statusCode                    // capture status code
}

var (
	modules = []string{
		domain.ModuleIdCredentials,
		domain.ModuleIdCdrs,
		domain.ModuleIdCommands,
		domain.ModuleIdHubClientInfo,
		domain.ModuleIdLocations,
		domain.ModuleIdSessions,
		domain.ModuleIdTariffs,
		domain.ModuleIdTokens,
		"versions",
	}
)

func (m *Middleware) OcpiLoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		start := kit.Now()

		msg := &domain.LogMessage{
			Url:           fmt.Sprintf("http://%s%s", r.Host, r.RequestURI),
			Token:         r.Header.Get(model.OcpiHeaderAuth),
			RequestId:     r.Header.Get(model.OcpiHeaderRequestId),
			CorrelationId: r.Header.Get(model.OcpiHeaderCorrelationId),
			Headers:       r.Header,
			ToPlatform:    m.cfg.Local.Platform.Id,
			In:            true,
		}

		// get event
		for _, mod := range modules {
			if strings.Index(r.RequestURI, mod) > 0 {
				msg.Event = strings.ToLower(fmt.Sprintf("%s.%s", mod, r.Method))
				break
			}
		}
		if msg.Event == "" {
			msg.Event = "version-details.get"
		}

		// token
		msg.Token, _ = m.ExtractToken(ctx, r)

		// request body
		if r.Body != nil {
			body, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewReader(body))
			if body != nil {
				msg.RequestBody = make(map[string]interface{})
				_ = json.Unmarshal(body, &msg.RequestBody)

			}
		}

		// call logging
		defer m.ocpiLogging.Log(ctx, msg)

		// apply writer wrapper to gather response
		wrapper := &writerWrapper{ResponseWriter: w}

		// next
		next.ServeHTTP(wrapper, r)

		// duration
		msg.DurationMs = time.Since(start).Milliseconds()

		// response
		msg.ResponseStatus = wrapper.status
		if wrapper.body != nil {
			rsMap := make(map[string]interface{})
			_ = json.Unmarshal(wrapper.body, &rsMap)
			if rsMap != nil {
				msg.ResponseBody = rsMap
				if st, ok := rsMap["status_code"]; ok && st != nil {
					msg.OcpiStatus = int(st.(float64))
				}
			}
		}

		// platform
		if rqCtx, ok := kit.Request(ctx); ok {
			if v, ok := rqCtx.Kv[model.OcpiCtxPlatform]; ok && v != nil {
				msg.FromPlatform = v.(string)
			}
		}
	}

	return f
}
