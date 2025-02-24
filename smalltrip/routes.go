package smalltrip

import (
	"fmt"
	"net/http"
	"strings"
)

type Route interface {
	http.Handler
	fmt.Stringer
}

type RouteGroup interface {
	Routes() []Route
}

type route struct {
	method string
	host   string
	path   string

	handler http.Handler
}

func newRoute(method, host, path string, handler http.Handler) Route {
	return &route{method, host, path, handler}
}

func (r *route) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.handler.ServeHTTP(w, req)
}

func (r *route) String() string {
	path := r.host + r.path

	if r.method != "" {
		path = r.method + " " + path
	}

	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	if strings.HasSuffix(path, "...}/") {
		path = strings.TrimSuffix(path, "/")
	}

	return path
}
