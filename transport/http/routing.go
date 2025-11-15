package http

import (
	"github.com/gorilla/mux"
	"github.com/mikhailbolshakov/kit"
	kitHttp "github.com/mikhailbolshakov/kit/http"
	"github.com/mikhailbolshakov/ocpi/errors"
	"net/http"
)

const (
	TokenA = "A"
	TokenB = "B"
	TokenC = "C"
)

type RouteBuilder struct {
	http   *kitHttp.Server
	mdw    *Middleware
	routes []*Route
}

func NewRouteBuilder(http *kitHttp.Server, mdw *Middleware) *RouteBuilder {
	return &RouteBuilder{
		http: http,
		mdw:  mdw,
	}
}

func (r *RouteBuilder) SetRoutes(routes []*Route) {
	r.routes = append(r.routes, routes...)
}

type Router struct {
}

type Route struct {
	id          string
	url         string
	urlPrefix   string
	verbs       []string
	handleFn    http.HandlerFunc
	handler     http.Handler
	middlewares []mux.MiddlewareFunc
	authTokens  []string
	subRouter   bool
	ocpiLogging bool
	apiKey      bool
}

func (r *RouteBuilder) Build() error {
	for _, route := range r.routes {
		// validate route
		err := route.validate()
		if err != nil {
			return err
		}
		// setup http routing
		httpRouter := r.http.RootRouter
		if route.subRouter {
			httpRouter = r.http.RootRouter.PathPrefix(route.urlPrefix).Handler(route.handler).Subrouter()
			// set middlewares is passed
			if len(route.middlewares) > 0 {
				httpRouter.Use(route.middlewares...)
			}
		}
		if route.handleFn != nil {
			handleFn := route.handleFn
			// if authentication, apply special middleware
			if r.mdw != nil {
				if len(route.authTokens) > 0 {
					handleFn = r.mdw.AuthAccessTokenMiddleware(handleFn, route.authTokens...)
				}
				if route.ocpiLogging {
					handleFn = r.mdw.OcpiLoggingMiddleware(handleFn)
				}
				if route.apiKey {
					handleFn = r.mdw.ApiKeyMiddleware(handleFn)
				}
			}
			httpRouter.HandleFunc(route.url, handleFn).Methods(route.verbs...)
		} else if route.handler != nil {
			// if handler specified, it means all processing done by it
			httpRouter.PathPrefix(route.urlPrefix).Handler(route.handler)
		}
	}
	return nil
}

// R starts building a new route with url and handle function
func R(url string, f func(http.ResponseWriter, *http.Request)) *Route {
	return &Route{
		id:       kit.NewRandString(),
		url:      url,
		handleFn: f,
	}
}

// SubRouter allows specifying a new area of routes with its own set of middlewares
func (r *Route) SubRouter(urlPrefix string) *Route {
	r.subRouter = true
	r.urlPrefix = urlPrefix
	return r
}

// Url specifies route's URL
func (r *Route) Url(url string) *Route {
	r.url = url
	return r
}

// PathPrefix specifies URL prefix
func (r *Route) PathPrefix(urlPrefix string) *Route {
	r.urlPrefix = urlPrefix
	return r
}

// POST applies post verb
func (r *Route) POST() *Route {
	r.verbs = append(r.verbs, "POST")
	return r
}

// PUT applies put verb
func (r *Route) PUT() *Route {
	r.verbs = append(r.verbs, "PUT")
	return r
}

// PATCH applies patch verb
func (r *Route) PATCH() *Route {
	r.verbs = append(r.verbs, "PATCH")
	return r
}

// GET applies get verb
func (r *Route) GET() *Route {
	r.verbs = append(r.verbs, "GET")
	return r
}

// DELETE applies delete verb
func (r *Route) DELETE() *Route {
	r.verbs = append(r.verbs, "DELETE")
	return r
}

// Auth marks route as authorized by token with types (A, B, C)
func (r *Route) Auth(tokenTypes ...string) *Route {
	r.authTokens = tokenTypes
	return r
}

// NoAuth marks route as not authorized
func (r *Route) NoAuth() *Route {
	r.authTokens = nil
	return r
}

// HandleFn specifies a handle function for route
func (r *Route) HandleFn(f func(http.ResponseWriter, *http.Request)) *Route {
	r.handleFn = f
	return r
}

// Handler allows specifying a handler which is applied to URL prefix
func (r *Route) Handler(h http.Handler) *Route {
	r.handler = h
	return r
}

// OcpiLogging allows logging all the OCPI requests
func (r *Route) OcpiLogging() *Route {
	r.ocpiLogging = true
	return r
}

// ApiKey applies auth by api key
func (r *Route) ApiKey() *Route {
	r.apiKey = true
	return r
}

// Middlewares allows specifying special middlewares applied to the route
// Note! It's applied only to SubRoute
func (r *Route) Middlewares(mdws ...mux.MiddlewareFunc) *Route {
	r.middlewares = append(r.middlewares, mdws...)
	return r
}

func (r *Route) validate() error {
	if r.url == "" && r.urlPrefix == "" {
		return errors.ErrRouteNotValid(r.url)
	}
	if len(r.verbs) == 0 && r.handleFn != nil {
		return errors.ErrRouteNotValid(r.url)
	}
	if r.handler == nil && r.handleFn == nil {
		return errors.ErrRouteNotValid(r.url)
	}
	if len(r.middlewares) > 0 && !r.subRouter {
		return errors.ErrRouteNotValid(r.url)
	}
	return nil
}
