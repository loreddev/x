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
	"context"
	"fmt"
	"net/http"

	"forge.capytal.company/loreddev/x/smalltrip/middleware"
)

type ContextKey string

var DefaultContextKey ContextKey = "x-smalltrip-problems-middleware-handler"

// TODO?: BufferedMiddleware, a middleware which can respond or redirect to
// a error page even after the first Write

func PanicMiddleware() middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rv := recover(); rv != nil {
					err := fmt.Errorf("panic recovered: %+v", rv)
					NewInternalServerError(err).ServeHTTP(w, r)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func Middleware(h Handler) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), DefaultContextKey, h)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func HandlerMiddleware(fallback ...Handler) Handler {
	return func(p Problem) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler := r.Context().Value(DefaultContextKey)
			if h, ok := handler.(Handler); handler != nil && ok {
				h(p).ServeHTTP(w, r)
			} else if len(fallback) > 0 {
				fallback[0](p).ServeHTTP(w, r)
			}
		})
	}
}
