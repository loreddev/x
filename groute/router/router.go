package router

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"forge.capytal.company/loreddev/x/groute/middleware"
)

type Router interface {
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler http.HandlerFunc)

	Use(middleware middleware.Middleware)

	http.Handler
}

type RouterWithRoutes interface {
	Router
	Routes() []Route
}

type RouterWithMiddlewares interface {
	RouterWithRoutes
	Middlewares() []middleware.Middleware
}

type RouterWithMiddlewaresWrapper interface {
	RouterWithMiddlewares
	WrapMiddlewares(ms []middleware.Middleware, h http.Handler) http.Handler
}

type Route struct {
	Path    string
	Method  string
	Host    string
	Handler http.Handler
}

func NewRouter(mux ...*http.ServeMux) Router {
	var m *http.ServeMux
	if len(mux) > 0 {
		m = mux[0]
	} else {
		m = http.NewServeMux()
	}

	return &defaultRouter{
		m,
		[]middleware.Middleware{},
		map[string]Route{},
	}
}

type defaultRouter struct {
	mux         *http.ServeMux
	middlewares []middleware.Middleware
	routes      map[string]Route
}

func (r *defaultRouter) Handle(pattern string, h http.Handler) {
	if sr, ok := h.(Router); ok {
		r.handleRouter(pattern, sr)
	} else {
		r.handle(pattern, h)
	}
}

func (r *defaultRouter) HandleFunc(pattern string, hf http.HandlerFunc) {
	r.handle(pattern, hf)
}

func (r *defaultRouter) Use(m middleware.Middleware) {
	r.middlewares = append(r.middlewares, m)
}

func (r *defaultRouter) Routes() []Route {
	rs := make([]Route, len(r.routes))
	i := 0
	for _, r := range r.routes {
		rs[i] = r
		i++
	}
	return rs
}

func (r *defaultRouter) Middlewares() []middleware.Middleware {
	return r.middlewares
}

func (r defaultRouter) WrapMiddlewares(ms []middleware.Middleware, h http.Handler) http.Handler {
	hf := h
	for _, m := range ms {
		hf = m(hf)
	}
	return hf
}

func (r *defaultRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r defaultRouter) handle(pattern string, hf http.Handler) {
	m, h, p := r.parsePattern(pattern)
	rt := Route{
		Method:  m,
		Host:    h,
		Path:    p,
		Handler: hf,
	}
	r.handleRoute(rt)
}

func (r defaultRouter) handleRouter(pattern string, rr Router) {
	m, h, p := r.parsePattern(pattern)

	rs, ok := rr.(RouterWithRoutes)
	if !ok {
		r.handle(p, rr)
	}

	routes := rs.Routes()
	middlewares := []middleware.Middleware{}
	if rw, ok := rs.(RouterWithMiddlewares); ok {
		middlewares = rw.Middlewares()
	}

	wrap := r.WrapMiddlewares
	if rw, ok := rs.(RouterWithMiddlewaresWrapper); ok {
		wrap = rw.WrapMiddlewares
	}

	for _, route := range routes {
		route.Handler = wrap(middlewares, route.Handler)
		route.Path = path.Join(p, route.Path)

		if m != "" && route.Method != "" && m != route.Method {
			panic(
				fmt.Sprintf(
					"Nested router's route has incompatible method than defined in path %q."+
						"Router's route method is %q, while path's is %q",
					p, route.Method, m,
				),
			)
		}
		if h != "" && route.Host != "" && h != route.Host {
			panic(
				fmt.Sprintf(
					"Nested router's route has incompatible host than defined in path %q."+
						"Router's route host is %q, while path's is %q",
					p, route.Host, h,
				),
			)
		}

		r.handleRoute(route)
	}
}

func (r defaultRouter) handleRoute(rt Route) {
	if len(r.middlewares) > 0 {
		rt.Handler = r.WrapMiddlewares(r.middlewares, rt.Handler)
	}

	if rt.Path == "" || !strings.HasPrefix(rt.Path, "/") {
		panic(
			fmt.Sprintf(
				"INVALID STATE: Path of route (%#v) does not start with a leading slash",
				rt,
			),
		)
	}

	p := rt.Path
	if rt.Host != "" {
		p = fmt.Sprintf("%s%s", rt.Host, p)
	}
	if rt.Method != "" {
		p = fmt.Sprintf("%s %s", rt.Method, p)
	}

	if !strings.HasSuffix(p, "/") {
		p = fmt.Sprintf("%s/", p)
	}

	if strings.HasSuffix(p, "...}/") {
		p = strings.TrimSuffix(p, "/")
	}

	r.routes[p] = rt
	r.mux.Handle(p, rt.Handler)
}

func (r *defaultRouter) parsePattern(pattern string) (method, host, p string) {
	pattern = strings.TrimSpace(pattern)

	// ServerMux patterns are "[METHOD ][HOST]/[PATH]", so to parsing it, we must
	// first split it between "[METHOD ][HOST]" and "[PATH]"
	ps := strings.Split(pattern, "/")

	p = path.Join("/", strings.Join(ps[1:], "/"))

	// If "[METHOD ][HOST]" is empty, we just have the path and can send it back
	if ps[0] == "" {
		return "", "", p
	}

	// Split string again, if method is not defined, this will end up being just []string{"[HOST]"}
	// since there isn't a space before the host. If there is a method defined, this will end up as
	// []string{"[METHOD]","[HOST]"}, with "[HOST]" being possibly a empty string.
	mh := strings.Split(ps[0], " ")

	// If slice is of length 1, this means it is []string{"[HOST]"}
	if len(mh) == 1 {
		return "", host, p
	}

	return mh[0], mh[1], p
}
