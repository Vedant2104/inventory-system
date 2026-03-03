package httprepo

import (
	"net/http"
)

func RegisterProductHandler(mux *http.ServeMux, handler *ProductHandler) {

	mux.HandleFunc("GET /product", handler.GetAllProduct)
	mux.HandleFunc("POST /product", handler.CreateProduct)
	mux.HandleFunc("GET /product/{id}", handler.GetProductById)
	mux.HandleFunc("PATCH /product/{id}", handler.UpdateProduct)
	mux.HandleFunc("DELETE /product/{id}", handler.DeleteProduct)
	mux.HandleFunc("POST /product/bulk", handler.BulkCreateFromCSV)
}

func RegisterProductCategoryHandler(mux *http.ServeMux, handler *ProductCategoryHandler) {

	mux.HandleFunc("GET /category", handler.GetAllProductCategory)
	mux.HandleFunc("POST /category", handler.CreateProductCategory)
	mux.HandleFunc("GET /category/{id}", handler.GetProductById)
	mux.HandleFunc("PATCH /category/{id}", handler.UpdateProductCategory)
	mux.HandleFunc("DELETE /category/{id}", handler.DeleteProductCategory)

}
