package router

import (
	"net/http"

	"forge.capytal.company/loreddev/x/groute/middleware"
)

var DefaultRouter = NewRouter()

func Handle(pattern string, handler http.Handler) {
	DefaultRouter.Handle(pattern, handler)
}

func HandleFunc(pattern string, handler http.HandlerFunc) {
	DefaultRouter.HandleFunc(pattern, handler)
}

func Use(m middleware.Middleware) {
	DefaultRouter.Use(m)
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	DefaultRouter.ServeHTTP(w, r)
}
