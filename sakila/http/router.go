package http

import (
	"compress/flate"
	"sakila/sakila-film-service/sakila"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// NewRouter returns a new http router.
func NewRouter(logger sakila.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.RequestID)
	r.Use(middleware.NewCompressor(flate.DefaultCompression).Handler)
	r.Use(RequestLogger(logger))
	return r
}
