package multiplexer

import "net/http"

func WithFormMethod(mux Multiplexer, key string) Multiplexer {
	return &formMethodMux{key: key, Multiplexer: mux}
}

type formMethodMux struct {
	key string
	Multiplexer
}

func (mux *formMethodMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if v := r.FormValue(mux.key); v != "" {
		r.Method = v
	}
	mux.Multiplexer.ServeHTTP(w, r)
}
