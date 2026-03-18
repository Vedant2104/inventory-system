package httprepo

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Vedant2104/inventory-system/internals/service"
)

type ProductCategoryHandler struct {
	service *service.ProductCategoryService
}

func NewProductCategoryHandler(service *service.ProductCategoryService) *ProductCategoryHandler {
	return &ProductCategoryHandler{
		service: service,
	}
}

func (h *ProductCategoryHandler) GetAllProductCategory(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	products, err := h.service.GetAllProductCategory(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (h *ProductCategoryHandler) CreateProductCategory(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	category, err := h.service.CreateProductCategory(ctx, input.Name, input.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}


func (h *ProductCategoryHandler) GetProductById(w http.ResponseWriter, r *http.Request){
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id := r.PathValue("id")

	category, err := h.service.GetProductCategoryById(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(category)
}

func (h *ProductCategoryHandler) DeleteProductCategory (w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	id := r.PathValue("id")

	err := h.service.DeleteProductCategory(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode("Product Category deleted successfully")
}	

func (h *ProductCategoryHandler) UpdateProductCategory(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id := r.PathValue("id")
	var doc struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}
	err := json.NewDecoder(r.Body).Decode(&doc)
	if err != nil{
	http.Error(w, err.Error(), http.StatusBadRequest)
	}

	updatedCategory , err := h.service.UpdateProductCategory(ctx, id, doc.Name, doc.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedCategory)
}