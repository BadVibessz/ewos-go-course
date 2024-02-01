package router

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Middleware = func(next http.Handler) http.Handler

func MakeRoutes(basePath string, routers map[string]chi.Router, middlewares []Middleware) chi.Router {
	r := chi.NewRouter()

	for _, middleware := range middlewares {
		r.Use(middleware)
	}

	for routerPath, router := range routers {
		r.Mount(fmt.Sprintf("%s%s", basePath, routerPath), router)
	}

	return r
}
