package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	httprepo "github.com/Vedant2104/inventory-system/internals/adapters/http"
	mongorepo "github.com/Vedant2104/inventory-system/internals/adapters/repository/mongo"
	"github.com/Vedant2104/inventory-system/internals/infrastructure/logger"
	"github.com/Vedant2104/inventory-system/internals/service"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	port := getEnv("PORT", "8080")
	mongoURI := getEnv("MONGO_URI", "mongodb://admin:password@localhost:27018")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to mongodb", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to connect to mongodb", err)
	}
	db := client.Database("inventory")
	product_collection := db.Collection("products")
	category_collection := db.Collection("product_categories")

	categoryRepo := mongorepo.NewProductCategoryRepository(category_collection)
	categoryService := service.NewProductCategoryService(categoryRepo)
	categoryHandler := httprepo.NewProductCategoryHandler(categoryService)
	
	
	// productRepo := maprepo.NewProductRepository()
	productRepo := mongorepo.NewProductRepository(product_collection)
	productService := service.NewProductService(productRepo , categoryService)
	productHandler := httprepo.NewProductHandler(productService)

	// router := chi.NewRouter()
	router := http.NewServeMux()
	logger := logger.GetLogger()

	// router.Mount("/", httprepo.RegisterRoutes(ProductHandler, logger))
	httprepo.RegisterProductHandler(router, productHandler)
	httprepo.RegisterProductCategoryHandler(router, categoryHandler)
	log.Println("Server Running at port", port)

	middleware := httprepo.RequestLogger(logger)
	handler := middleware(router)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
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
