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
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"runtime"

	"forge.capytal.company/loreddev/x/smalltrip/middleware"
	"forge.capytal.company/loreddev/x/smalltrip/multiplexer"
)

type Router interface {
	multiplexer.Multiplexer
	Use(middleware.Middleware)
}

type router struct {
	mux multiplexer.Multiplexer
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
	log := router.log.With(slog.String("pattern", pattern), slog.String("handler", getValueType(handler)))
	log.Info("Adding route")

	var hf http.Handler = http.HandlerFunc(handler)

	for _, m := range router.mws {
		log.Debug("Wrapping with middleware", slog.String("middleware", getValueType(m)))
		hf = m(hf)
	}

	router.mux.Handle(pattern, hf)
}

func (router *router) Handle(pattern string, handler http.Handler) {
	log := router.log.With(slog.String("pattern", pattern), slog.String("handler", getValueType(handler)))
	log.Info("Adding route")

	for _, m := range router.mws {
		log.Debug("Wrapping with middleware", slog.String("middleware", getValueType(m)))
		handler = m(handler)
	}

	router.mux.Handle(pattern, handler)
}

func (router *router) Use(m middleware.Middleware) {
	router.log.Info("Middleware added", slog.String("middleware", getValueType(m)))

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

func getValueType[T any](value T) (name string) {
	defer func() {
		if rc := recover(); rc != nil {
			name = fmt.Sprintf("%T", value)
		}
	}()

	v := reflect.ValueOf(value)

	if v.Kind() == reflect.Pointer {
		return getValueType(v.Elem().Interface())
	}

	if v.Kind() == reflect.Func {
		fc := runtime.FuncForPC(v.Pointer())
		if fc != nil {
			return fc.Name()
		}
	}

	if p, n := v.Type().PkgPath(), v.Type().Name(); p != "" && n != "" {
		return fmt.Sprintf("%s.%s", p, n)
	}

	return fmt.Sprintf("%T", value)
}
