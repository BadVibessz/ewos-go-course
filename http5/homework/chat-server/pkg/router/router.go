package router

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Middleware = func(next http.Handler) http.Handler

func MakeRoutes(basePath string, routers map[string]chi.Router, middlewares []Middleware) *chi.Mux {
	r := chi.NewRouter()

	for _, middleware := range middlewares {
		r.Use(middleware)
	}

	for routerPath, router := range routers {
		r.Mount(fmt.Sprintf("%s%s", basePath, routerPath), router)
	}

	return r
}
