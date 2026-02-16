package httprepo

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handler *ProductHandler) http.Handler {
	r := chi.NewRouter()

	r.Route("/product", func(r chi.Router) {
		r.Post("/", handler.CreateProduct)
		r.Get("/", handler.GetAllProduct)
		r.Get("/{id}", handler.GetProductById)
		r.Delete("/{id}", handler.DeleteProduct)
		r.Patch("/{id}", handler.UpdateProduct)
	})

	return r
}
