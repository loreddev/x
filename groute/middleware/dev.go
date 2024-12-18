package middleware

import (
	"log/slog"
	"math/rand"
	"net/http"
)

func DevMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

type loggerReponse struct {
	http.ResponseWriter
	status int
}

func (lr *loggerReponse) WriteHeader(s int) {
	lr.status = s
	lr.ResponseWriter.WriteHeader(s)
}

func NewLoggerMiddleware(l *slog.Logger) Middleware {
	l = l.WithGroup("logger_middleware")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := randHash(5)

			l.Info("NEW REQUEST",
				slog.String("id", id),
				slog.String("status", "xxx"),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)

			lw := &loggerReponse{w, http.StatusOK}
			next.ServeHTTP(lw, r)

			if lw.status >= 400 {
				l.Warn("ERR REQUEST",
					slog.String("id", id),
					slog.Int("status", lw.status),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
				)
				return
			}

			l.Info("END REQUEST",
				slog.String("id", id),
				slog.Int("status", lw.status),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)
		})
	}
}

const HASH_CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// This is not the most performant function, as a TODO we could
// improve based on this Stackoberflow thread:
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func randHash(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = HASH_CHARS[rand.Int63()%int64(len(HASH_CHARS))]
	}
	return string(b)
}
