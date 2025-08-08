package middleware

import (
	"net/http"
	"strings"
)

func FormMethod(key string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if v := r.FormValue(key); v != "" {
				r.Method = strings.ToUpper(v)
			}
			next.ServeHTTP(w, r)
		})
	}
}
