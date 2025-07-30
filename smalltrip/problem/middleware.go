package problem

import (
	"context"
	"net/http"

	"forge.capytal.company/loreddev/x/smalltrip/middleware"
)

func Middleware(h Handler) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), contextKey, h)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func HandlerMiddleware(fallback ...Handler) Handler {
	return func(p Problem) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler := r.Context().Value(contextKey)
			if h, ok := handler.(Handler); handler != nil && ok {
				h(p).ServeHTTP(w, r)
			} else if len(fallback) > 0 {
				fallback[0](p).ServeHTTP(w, r)
			}
		})
	}
}

var contextKey = "x-smalltrip-problems-middleware-handler"
