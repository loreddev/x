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

package middleware

import (
	"fmt"
	"log/slog"
	"math/rand"
	"net"
	"net/http"
)

func Logger(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lw := &loggerResponseWriter{w, 0}

			addr := loggerGetAddr(r)
			if net.ParseIP(addr) == nil {
				addr = fmt.Sprintf("INVALID %s", addr)
			}

			log := logger.With(
				slog.String("id", randHash(5)),
				slog.String("method", fmt.Sprintf("%4s", r.Method)),
				slog.String("addr", addr),
				slog.String("path", r.URL.Path),
			)

			log.Debug("NEW REQUEST", slog.String("status", "000"))

			next.ServeHTTP(lw, r)

			log = log.With(slog.String("status", fmt.Sprintf("%3d", lw.statusCode)))

			switch {
			case lw.statusCode >= 500:
				log.Warn("ERR REQUEST")
			case lw.statusCode >= 400:
				log.Info("INV REQUEST")
			case lw.statusCode >= 200:
				log.Debug("END REQUEST")
			default:
				log.Debug("MSC REQUEST")
			}
		})
	}
}

func loggerGetAddr(r *http.Request) string {
	if i := r.Header.Get("CF-Connecting-IP"); i != "" {
		return i
	}
	if i := r.Header.Get("X-Forwarded-For"); i != "" {
		return i
	}
	if i := r.Header.Get("X-Real-IP"); i != "" {
		return i
	}
	return r.RemoteAddr
}

type loggerResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *loggerResponseWriter) WriteHeader(s int) {
	w.statusCode = s
	w.ResponseWriter.WriteHeader(s)
}

const hashChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// This is not the most performant function, as a TODO we could
// improve based on this Stackoberflow thread:
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func randHash(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = hashChars[rand.Int63()%int64(len(hashChars))]
	}
	return string(b)
}
