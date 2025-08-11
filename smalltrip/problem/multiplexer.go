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

package problem

import (
	"net/http"
	"strings"

	"forge.capytal.company/loreddev/x/smalltrip/multiplexer"
)

func Multiplexer(m multiplexer.Multiplexer, opts ...MultiplexerOption) multiplexer.Multiplexer {
	mux := &mux{Multiplexer: m, methodList: DefaultMethods}

	for _, opt := range opts {
		opt(mux)
	}

	return mux
}

type mux struct {
	notFound         http.Handler
	methodNotAllowed http.Handler
	methodList       []string

	multiplexer.Multiplexer
}

func (m *mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i := &interceptor{
		notFound:         m.notFound,
		methodNotAllowed: m.methodNotAllowed,
		methodList:       m.methodList,

		mux: m,

		w: w,
		r: r,
	}
	m.Multiplexer.ServeHTTP(i, r)
}

type interceptor struct {
	notFound         http.Handler
	methodNotAllowed http.Handler
	methodList       []string

	mux multiplexer.Multiplexer

	statusCode int
	intercept  bool

	w http.ResponseWriter
	r *http.Request
}

var _ http.ResponseWriter = (*interceptor)(nil)

func (e *interceptor) Header() http.Header {
	return e.w.Header()
}

func (e *interceptor) WriteHeader(statusCode int) {
	if statusCode > 399 && strings.Contains(e.w.Header().Get("Content-Type"), "text/plain") {
		e.w.Header().Del("Content-Type")
		e.intercept = true
		e.statusCode = statusCode
		return
	}
	e.w.WriteHeader(statusCode)
}

func (e *interceptor) Write(data []byte) (int, error) {
	if e.intercept && e.statusCode == http.StatusMethodNotAllowed {
		method := e.r.Method
		_, current := e.mux.Handler(e.r)

		var allowed []string
		for _, m := range e.methodList {
			e.r.Method = m
			if _, p := e.mux.Handler(e.r); p != current {
				allowed = append(allowed, m)
			}
		}

		e.r.Method = method

		if e.methodNotAllowed != nil {
			e.w.Header().Set("Allow", strings.Join(allowed, ", "))
			e.methodNotAllowed.ServeHTTP(e.w, e.r)
		} else {
			NewMethodNotAllowed(allowed).ServeHTTP(e.w, e.r)
		}
		return len(data), nil

	} else if e.intercept && e.statusCode == http.StatusNotFound && e.notFound != nil {
		e.notFound.ServeHTTP(e.w, e.r)
		return len(data), nil

	} else if e.intercept {
		NewDetailed(e.statusCode, strings.TrimSpace(string(data))).ServeHTTP(e.w, e.r)
		return len(data), nil
	}

	return e.w.Write(data)
}

var DefaultMethods = multiplexer.DefaultMethods

type MultiplexerOption func(*mux)

func WithNotFound(h http.Handler) MultiplexerOption {
	return func(m *mux) {
		m.notFound = h
	}
}

func WithMethodNotAllowed(h http.Handler) MultiplexerOption {
	return func(m *mux) {
		m.methodNotAllowed = h
	}
}

func WithMethod(method string) MultiplexerOption {
	return func(m *mux) {
		m.methodList = append(m.methodList, method)
	}
}

func WithMethodList(list []string) MultiplexerOption {
	return func(m *mux) {
		m.methodList = list
	}
}
