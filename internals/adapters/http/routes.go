package httprepo

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

func RegisterRoutes(handler *ProductHandler, logger zerolog.Logger) http.Handler {
	r := chi.NewRouter()
	r.Use(RequestLogger(logger))
	r.Route("/product", func(r chi.Router) {
		r.Post("/", handler.CreateProduct)
		r.Get("/", handler.GetAllProduct)
		r.Get("/{id}", handler.GetProductById)
		r.Delete("/{id}", handler.DeleteProduct)
		r.Patch("/{id}", handler.UpdateProduct)
	})

	return r
}
