package httprepo

import (
	"net/http"

	"github.com/rs/zerolog"
)

func RegisterProductHandler(mux *http.ServeMux, handler *ProductHandler, logger zerolog.Logger) *http.Handler {

	mux.HandleFunc("GET /product", handler.GetAllProduct)
	mux.HandleFunc("POST /product", handler.CreateProduct)
	mux.HandleFunc("GET /product/{id}", handler.GetProductById)
	mux.HandleFunc("PATCH /product/{id}", handler.UpdateProduct)
	mux.HandleFunc("DELETE /product/{id}", handler.DeleteProduct)

	middleware := RequestLogger(logger)

	wrappedMux := middleware(mux)
	return &wrappedMux
}
