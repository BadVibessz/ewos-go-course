package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

func MakeRoutes(basePath string, routers []chi.Router) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)

	for _, router := range routers {
		r.Mount(basePath, router)
	}

	return r
}
