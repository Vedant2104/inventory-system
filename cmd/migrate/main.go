package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

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

	log.Println("Starting Migration...")

	var uniqueNames []string
	if err  = product_collection.Distinct(ctx, "category", bson.D{}).Decode(&uniqueNames) ; err != nil{
		log.Fatal(err)
	}

	
	for _ , catName := range uniqueNames{
		if catName == ""{
			continue
		}
		filter := bson.D{{Key: "name", Value: catName}}
		update := bson.M{
			"$setOnInsert" : bson.M{
				"name" : catName,
				"description" : catName,
			},
		}
		opts := options.UpdateOne().SetUpsert(true)
		_ ,err := category_collection.UpdateOne(ctx, filter, update, opts)
		if err != nil{
			log.Fatal(err)
		}

		var category struct{
			ID bson.ObjectID `bson:"_id"`
		}
		if err = category_collection.FindOne(ctx, filter).Decode(&category);err != nil{
			log.Fatal(err)
		}

		update = bson.M{
			"$set" : bson.M{
				"category" : category.ID,
			},
		}
		filter = bson.D{{Key: "category", Value: catName}}
		result , err := product_collection.UpdateMany(ctx, filter, update)
		if err != nil{
			log.Fatal(err)
		}
		log.Printf("Updated %d documents", result.ModifiedCount)
	}

	log.Println("Migration Completed")
}	
