package router

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func MakeRoutes(basePath string, routers map[string]chi.Router) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)

	for routerPath, router := range routers {
		r.Mount(fmt.Sprintf("%s%s", basePath, routerPath), router)
	}

	return r
}
