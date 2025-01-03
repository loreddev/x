package rerrors

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"forge.capytal.company/loreddev/x/groute/middleware"
	"github.com/a-h/templ"
)

const (
	ERROR_MIDDLEWARE_HEADER = "XX-Error-Middleware"
	ERROR_VALUE_HEADER      = "X-Error-Value"
)

type RouteError struct {
	StatusCode int            `json:"status_code"`
	Err        string         `json:"error"`
	Info       map[string]any `json:"info"`
	Endpoint   string
}

func NewRouteError(status int, error string, info ...map[string]any) RouteError {
	rerr := RouteError{StatusCode: status, Err: error}
	if len(info) > 0 {
		rerr.Info = info[0]
	} else {
		rerr.Info = map[string]any{}
	}
	return rerr
}

func (rerr RouteError) Error() string {
	return fmt.Sprintf("route error %d %s: %v", rerr.StatusCode, rerr.Endpoint, rerr.Info)
}

func (rerr RouteError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if rerr.StatusCode == 0 {
		rerr.StatusCode = http.StatusNotImplemented
	}

	if rerr.Err == "" {
		rerr.Err = "MISSING ERROR DESCRIPTION"
	}

	if rerr.Info == nil {
		rerr.Info = map[string]any{}
	}

	j, err := json.Marshal(rerr)
	if err != nil {
		j, _ = json.Marshal(RouteError{
			StatusCode: http.StatusInternalServerError,
			Err:        "Failed to marshal error message to JSON",
			Info: map[string]any{
				"source_value": fmt.Sprintf("%#v", rerr),
				"error":        err.Error(),
			},
		})
	}

	if r.Header.Get(ERROR_MIDDLEWARE_HEADER) == "enable" && prefersHtml(r.Header) {
		q := r.URL.Query()
		q.Set("error", base64.URLEncoding.EncodeToString(j))
		r.URL.RawQuery = q.Encode()

		http.Redirect(w, r, r.URL.String(), http.StatusTemporaryRedirect)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(rerr.StatusCode)
	if _, err = w.Write(j); err != nil {
		_, _ = w.Write([]byte("Failed to write error JSON string to body"))
	}
}

type ErrorMiddlewarePage func(err RouteError) templ.Component

type ErrorDisplayer struct {
	log  *slog.Logger
	page ErrorMiddlewarePage
}

func (h ErrorDisplayer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e, err := base64.URLEncoding.DecodeString(r.URL.Query().Get("error"))
	if err != nil {
		h.log.Error("Failed to decode \"error\" parameter from error redirect",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", 0),
			slog.String("data", string(e)),
		)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(
			fmt.Sprintf("Data %s\nError %s", string(e), err.Error()),
		))
		return
	}

	var rerr RouteError
	if err := json.Unmarshal(e, &rerr); err != nil {
		h.log.Error("Failed to decode \"error\" parameter from error redirect",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", 0),
			slog.String("data", string(e)),
		)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(
			fmt.Sprintf("Data %s\nError %s", string(e), err.Error()),
		))
		return
	}

	if rerr.Endpoint == "" {
		q := r.URL.Query()
		q.Del("error")
		r.URL.RawQuery = q.Encode()

		rerr.Endpoint = r.URL.String()
	}

	w.WriteHeader(rerr.StatusCode)
	if err := h.page(rerr).Render(r.Context(), w); err != nil {
		_, _ = w.Write(e)
	}
}

func NewErrorMiddleware(
	p ErrorMiddlewarePage,
	l *slog.Logger,
	notfound ...ErrorMiddlewarePage,
) middleware.Middleware {
	var nf ErrorMiddlewarePage
	if len(notfound) > 0 {
		nf = notfound[0]
	} else {
		nf = p
	}

	l = l.WithGroup("error_middleware")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Header.Set(ERROR_MIDDLEWARE_HEADER, "enable")

			if uerr := r.URL.Query().Get("error"); uerr != "" && prefersHtml(r.Header) {
				ErrorDisplayer{l, nf}.ServeHTTP(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func prefersHtml(h http.Header) bool {
	if h.Get("Accept") == "" {
		return false
	}
	return (strings.Contains(h.Get("Accept"), "text/html") ||
		strings.Contains(h.Get("Accept"), "application/xhtml+xml") ||
		strings.Contains(h.Get("Accept"), "application/xml")) &&
		!strings.Contains(h.Get("Accept"), "application/json")
}
