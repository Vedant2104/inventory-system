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
	collection := db.Collection("products")
	// productRepo := maprepo.NewProductRepository()
	productRepo := mongorepo.NewProductRepository(collection)
	productService := service.NewProductService(productRepo)

	ProductHandler := httprepo.NewProductHandler(productService)

	// router := chi.NewRouter()
	router := http.ServeMux{}
	logger := logger.GetLogger()

	// router.Mount("/", httprepo.RegisterRoutes(ProductHandler, logger))
	productHandler := httprepo.RegisterProductHandler(&router, ProductHandler, logger)
	log.Println("Server Running at port", port)

	if err := http.ListenAndServe(":"+port, *productHandler); err != nil {
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
