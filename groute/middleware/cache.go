package middleware

import (
	"net/http"
)

func CacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=604800, stale-while-revalidate=86400, public")
		next.ServeHTTP(w, r)
	})
}
