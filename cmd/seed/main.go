package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Vedant2104/inventory-system/internals/domain"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type mongoProduct struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
	Category    bson.ObjectID `bson:"category"`
	Price       int           `bson:"price"`
	Brand       string        `bson:"brand"`
	Quantity    int           `bson:"quantity"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found", err)
	}

	mongoURI := os.Getenv("MONGO_URI")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))

	if err != nil {
		log.Fatal("Failed to connect to mongodb", err)
	}

	defer client.Disconnect(ctx)
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to connect to mongodb", err)
	}

	db := client.Database("inventory")
	product_collection := db.Collection("products")
	category_collection := db.Collection("product_categories")

	log.Println("Seeding Categories...")

	categories := []domain.ProductCategory{
		{
			Name:        "Electronics",
			Description: "Tech and Gadgets",
		},
		{
			Name:        "Clothing",
			Description: "Men's and Women's Clothing",
		},
		{
			Name:        "Furniture",
			Description: "Living Room, Dining Room, Bedroom, Kitchen",
		},
		{
			Name:        "Home Decor",
			Description: "Decorative Items for Home",
		},
		{
			Name:        "Toys",
			Description: "Play and Fun",
		},
		{
			Name:        "Books",
			Description: "Read and Learn",
		},
		{
			Name:        "Sports",
			Description: "Play and Fun",
		},
		{
			Name:        "Beauty and Personal Care",
			Description: "Beauty and Personal Care",
		},
		{
			Name:        "Health and Wellness",
			Description: "Health and Wellness",
		},
		{
			Name:        "Jewelry",
			Description: "Jewelry and Accessories",
		},
	}

	catMap := make(map[string]bson.ObjectID)
	for _, category := range categories {
		filter := bson.D{{Key: "name", Value: category.Name}}
		update := bson.D{{Key: "$setOnInsert", Value: bson.D{
			{Key: "name", Value: category.Name},
			{Key: "description", Value: category.Description},
		}}}
		var doc struct {
			ID          bson.ObjectID `bson:"_id"`
			Name        string        `bson:"name"`
			Description string        `bson:"description"`
		}
		opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
		err := category_collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&doc)
		if err != nil {
			log.Fatal(err)
		}
		catMap[category.Name] = doc.ID
	}

	log.Println("Seeding Products...")

	products := []mongoProduct{
		{
			Name:        "iPhone 13",
			Description: "Smartphone",
			Category:    catMap["Electronics"],
			Price:       999,
			Brand:       "Apple",
			Quantity:    10,
		},
		{
			Name:        "Macbook Pro",
			Description: "Laptop",
			Category:    catMap["Electronics"],
			Price:       1999,
			Brand:       "Apple",
			Quantity:    5,
		},
		{
			Name:        "AirPods",
			Description: "Wireless Earbuds",
			Category:    catMap["Electronics"],
			Price:       249,
			Brand:       "Apple",
			Quantity:    15,
		},
		{
			Name:        "T-Shirt",
			Description: "Cotton T-Shirt",
			Category:    catMap["Clothing"],
			Price:       19,
			Brand:       "H&M",
			Quantity:    20,
		},
		{
			Name:        "Jeans",
			Description: "Denim Jeans",
			Category:    catMap["Clothing"],
			Price:       49,
			Brand:       "Levis",
			Quantity:    10,
		},
		{
			Name:        "Sofa",
			Description: "Leather Sofa",
			Category:    catMap["Furniture"],
			Price:       999,
			Brand:       "IKEA",
			Quantity:    5,
		},
		{
			Name:        "Chair",
			Description: "Wooden Chair",
			Category:    catMap["Furniture"],
			Price:       49,
			Brand:       "IKEA",
			Quantity:    15,
		},
		{
			Name:        "Pillow",
			Description: "Soft Pillow",
			Category:    catMap["Furniture"],
			Price:       19,
			Brand:       "IKEA",
			Quantity:    20,
		},
		{
			Name:        "Game of Thrones",
			Description: "Never going to be finished",
			Category:    catMap["Books"],
			Price:       9,
			Brand:       "Penguin",
			Quantity:    10,
		},
		{
			Name:        "Basketball",
			Description: "Play and Fun",
			Category:    catMap["Sports"],
			Price:       49,
			Brand:       "Nike",
			Quantity:    5,
		},
		{
			Name:        "Tennis Racket",
			Description: "Play and Fun",
			Category:    catMap["Sports"],
			Price:       99,
			Brand:       "Adidas",
			Quantity:    15,
		},
		{
			Name:        "Gaming Console",
			Description: "Play and Fun",
			Category:    catMap["Electronics"],
			Price:       499,
			Brand:       "Sony",
			Quantity:    10,
		},
		{
			Name:        "Smart Watch",
			Description: "Track and Monitor",
			Category:    catMap["Electronics"],
			Price:       199,
			Brand:       "Fitbit",
			Quantity:    5,
		},
		{
			Name:        "Laptop",
			Description: "Work and Play",
			Category:    catMap["Electronics"],
			Price:       999,
			Brand:       "Dell",
			Quantity:    10,
		},
		{
			Name:        "Camera",
			Description: "Capture Moments",
			Category:    catMap["Electronics"],
			Price:       499,
			Brand:       "Canon",
			Quantity:    5,
		},
		{
			Name:        "Sunscreen",
			Description: "Protect and Hydrate",
			Category:    catMap["Beauty and Personal Care"],
			Price:       49,
			Brand:       "Neutrogena",
			Quantity:    10,
		},
		{
			Name:        "Shampoo",
			Description: "Clean and Refresh",
			Category:    catMap["Beauty and Personal Care"],
			Price:       19,
			Brand:       "Pantene",
			Quantity:    5,
		},
	}

	_, err = product_collection.InsertMany(ctx, products)

	if err != nil {
		log.Fatal(err)
	}
}
