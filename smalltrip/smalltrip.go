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
	mws []middleware.Middleware
	log *slog.Logger
}

var _ Router = (*router)(nil)

func NewRouter(options ...Option) Router {
	r := &router{
		mux: http.NewServeMux(),
		mws: []middleware.Middleware{},
		log: slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})).WithGroup("smalltrip-router"),
	}

	for _, option := range options {
		option(r)
	}

	return r
}

func (router *router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	router.Handle(pattern, http.HandlerFunc(handler))
}

func (router *router) Handle(pattern string, handler http.Handler) {
	log := router.log.With(slog.String("pattern", pattern), slog.String("handler", fmt.Sprintf("%T", handler)))
	log.Info("Adding route")

	for _, m := range router.mws {
		log.Debug("Wrapping with middleware", slog.String("middleware", fmt.Sprintf("%T", m)))
		handler = m(handler)
	}

	router.mux.Handle(pattern, handler)
}

func (router *router) Use(m middleware.Middleware) {
	router.log.Info("Middleware added", slog.String("middleware", fmt.Sprintf("%T", m)))

	if router.mws == nil {
		router.mws = []middleware.Middleware{}
	}
	router.mws = append(router.mws, m)
}

func (router *router) Handler(r *http.Request) (http.Handler, string) {
	return router.mux.Handler(r)
}

func (router *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}
