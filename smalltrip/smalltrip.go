// Copyright 2025-present Gustavo "Guz" L. de Mello
// Copyright 2025-present The Lored.dev Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package smalltrip

import (
	"io"
	"log/slog"
	"net/http"
	"path"
	"strings"

	"forge.capytal.company/loreddev/x/smalltrip/middleware"
	"forge.capytal.company/loreddev/x/tinyssert"
)

type Router interface {
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))

	Use(middleware middleware.Middleware)

	http.Handler
}

type router struct {
	mux         *http.ServeMux
	routes      map[string]Route
	middlewares []middleware.Middleware

	assert tinyssert.Assertions
	log    *slog.Logger
}

var (
	_ Router     = (*router)(nil)
	_ RouteGroup = (*router)(nil)
)

func NewRouter(options ...Option) Router {
	r := &router{
		mux:         http.NewServeMux(),
		routes:      map[string]Route{},
		middlewares: []middleware.Middleware{},

		assert: tinyssert.NewDisabledAssertions(),
		log:    slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})).WithGroup("smalltrip-router"),
	}

	for _, option := range options {
		option(r)
	}

	return r
}

func (r *router) Handle(pattern string, handler http.Handler) {
	r.assert.NotNil(handler, "Handler should not be nil, invalid state.")
	r.assert.NotZero(pattern, "Path should not be empty, invalid state.")
	r.assert.NotNil(r.log)

	log := r.log.With(slog.String("pattern", pattern))
	log.Info("Adding route")

	if router, ok := handler.(RouteGroup); ok {
		r.log.Debug("Route has nested router as handler, handling router's routes")

		r.handleGroup(pattern, router)
		return
	}

	method, host, p := parsePattern(pattern)
	r.assert.NotZero(p)

	log.Debug("Parsed route pattern",
		slog.String("method", method), slog.String("host", host), slog.String("path", p))

	route := newRoute(method, host, p, handler)
	r.handleRoute(route)
}

func (r *router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	r.Handle(pattern, http.HandlerFunc(handler))
}

func (r *router) Use(m middleware.Middleware) {
	r.assert.NotNil(m, "Middleware should not be nil value, invalid state")
	r.assert.NotNil(r.middlewares)
	r.assert.NotNil(r.log)

	r.log.Info("Adding middleware", slog.Any("middleware", m))

	r.middlewares = append(r.middlewares, m)
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *router) Routes() []Route {
	r.assert.NotNil(r.routes)

	rs := make([]Route, len(r.routes))

	var i int
	for _, v := range r.routes {
		rs[i] = v
		i++
	}

	return rs
}

func (r *router) handleGroup(pattern string, group RouteGroup) {
	r.assert.NotNil(group, "Router should not be nil, invalid state.")
	r.assert.NotZero(pattern, "Pattern should not be empty, invalid state.")
	r.assert.NotNil(r.mux)
	r.assert.NotNil(r.log)

	log := r.log.With(slog.String("pattern", pattern))

	method, host, p := parsePattern(pattern)
	r.assert.NotZero(p)

	log.Debug("Parsed route pattern",
		slog.String("method", method), slog.String("host", host), slog.String("path", p))

	for _, route := range group.Routes() {
		log := log.With("route-pattern", route.String())
		log.Debug("Adding group's route to parent")

		rMethod, rHost, rPath := parsePattern(route.String())

		log.Debug("Parsed route pattern",
			slog.String("route-method", method), slog.String("route-host", host), slog.String("route-path", p))

		if method != "" && rMethod != "" {
			r.assert.Equal(method, rMethod, "Nested group's route has incompatible method in route %q", pattern)
		}
		if host != "" && rHost != "" {
			r.assert.Equal(method, rMethod, "Nested group's route has incompatible method in route %q", pattern)
		}

		if method == "" {
			log.Debug("Parent method is empty, using route's method")
			method = rMethod
		}
		if host == "" {
			log.Debug("Parent host is empty, using route's host")
			host = rHost
		}

		route = newRoute(method, host, path.Join(p, rPath), route)

		log.Debug("Adding final route", slog.String("final-pattern", route.String()))

		r.handleRoute(route)
	}
}

func (r *router) handleRoute(route Route) {
	r.assert.NotNil(route, "Route should not be nil, invalid state.")
	r.assert.NotZero(route.String(), "Route pattern should not be empty, invalid state.")
	r.assert.NotNil(r.routes)
	r.assert.NotNil(r.mux)
	r.assert.NotNil(r.log)

	if len(r.middlewares) == 0 {
		pattern := route.String()
		r.routes[pattern] = route
		r.mux.Handle(pattern, route)

		return
	}

	log := r.log.With("pattern", route.String())

	handler := route.(http.Handler)

	for _, m := range r.middlewares {
		log.Debug("Wrapping route handler with middleware", slog.Any("middleware", m))

		handler = m(route)
	}

	method, host, p := parsePattern(route.String())
	r.assert.NotZero(p)

	route = newRoute(method, host, p, handler)

	pattern := route.String()
	r.routes[pattern] = route
	r.mux.Handle(pattern, route)
}

func parsePattern(pattern string) (method, host, p string) {
	pattern = strings.TrimSpace(pattern)

	// ServerMux patterns are "[METHOD ][HOST]/[PATH]", so to parsing it, we must
	// first split it between "[METHOD ][HOST]" and "[PATH]"
	ps := strings.Split(pattern, "/")

	p = path.Join("/", strings.Join(ps[1:], "/"))

	// If "[METHOD ][HOST]" is empty, we just have the path and can send it back
	if ps[0] == "" {
		return "", "", p
	}

	// Split string again, if method is not defined, this will end up being just []string{"[HOST]"}
	// since there isn't a space before the host. If there is a method defined, this will end up as
	// []string{"[METHOD]","[HOST]"}, with "[HOST]" being possibly a empty string.
	mh := strings.Split(ps[0], " ")

	// If slice is of length 1, this means it is []string{"[HOST]"}
	if len(mh) == 1 {
		return "", host, p
	}

	return mh[0], mh[1], p
}
