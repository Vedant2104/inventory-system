package main

import (
	"log"
	"net/http"
	"os"

	httprepo "github.com/Vedant2104/inventory-system/internals/adapters/http"
	maprepo "github.com/Vedant2104/inventory-system/internals/adapters/repository/map"
	"github.com/Vedant2104/inventory-system/internals/infrastructure/logger"
	"github.com/Vedant2104/inventory-system/internals/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	port := getEnv("PORT", "8080")

	productRepo := maprepo.NewProductRepository()

	productService := service.NewProductService(productRepo)

	ProductHandler := httprepo.NewProductHandler(productService)

	router := chi.NewRouter()

	logger := logger.GetLogger()

	router.Mount("/", httprepo.RegisterRoutes(ProductHandler, logger))

	log.Println("Server Running at port", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}

}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
