package smalltrip

import (
	"log/slog"
	"net/http"

	"forge.capytal.company/loreddev/x/tinyssert"
)

type Option func(*router)

func WithAssertions(assertions tinyssert.Assertions) Option {
	return func(r *router) {
		r.assert = assertions
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(r *router) {
		r.log = logger
	}
}

func WithServeMux(mux *http.ServeMux) Option {
	return func(r *router) {
		r.mux = mux
	}
}
