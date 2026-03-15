//go:build integration
// +build integration

package service

import (
	"context"
	"log"
	"os"
	"testing"

	mongorepo "github.com/Vedant2104/inventory-system/internals/adapters/repository/mongo"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var testClient *mongo.Client
var testDb *mongo.Database

func TestMain(m *testing.M) {

	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found", err)
	}

	url := os.Getenv("MONGO_URI")
	ctx := context.Background()

	client, err := mongo.Connect(options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal("Failed to connect to mongodb", err)
	}
	testClient = client
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to connect to mongodb", err)
	}

	testDb = client.Database("inventory_test")

	testDb.Drop(ctx)

	code := m.Run()

	testDb.Drop(ctx)
	testClient.Disconnect(ctx)
	os.Exit(code)
}

func TestIntegation_CreateProduct(t *testing.T) {
	productColl := testDb.Collection("products")
	catColl := testDb.Collection("product_categories")
	productRepo := mongorepo.NewProductRepository(productColl)
	catRepo := mongorepo.NewProductCategoryRepository(catColl)
	CatService := NewProductCategoryService(catRepo)
	prodService := NewProductService(productRepo, CatService)

	t.Run("should create product and save to mongo", func(t *testing.T) {
		ctx := context.Background()
		cat, _ := CatService.CreateProductCategory(ctx, "test_Category", "test_Description")
		prod, err := prodService.CreateProduct(ctx, "test_product", "test_description", cat.ID, 100, "test_brand", 10)
		if err != nil {
			t.Fatalf("expected no error but got %v", err.Error())
		}
		if prod == nil {
			t.Fatalf("expected product but got nil")
		}
		t.Cleanup(func() {
			productColl.DeleteMany(ctx, bson.M{})
			catColl.DeleteMany(ctx, bson.M{})
		})
		var dbResult bson.M
		objectId, err := bson.ObjectIDFromHex(prod.ID)
		if err != nil {
			log.Fatalf("failed to convert id to object id")
		}
		err = productColl.FindOne(ctx, bson.M{"_id": objectId}).Decode(&dbResult)

		if err != nil {
			t.Fatalf("expected no error but got %v", err.Error())
		}

		if dbResult["name"] != "test_product" {
			t.Fatalf("expected product name to be test_product but got %v", dbResult["name"])
		}
	})
}

func TestIntegration_GetAllProducts(t *testing.T) {
	productColl := testDb.Collection("products")
	catColl := testDb.Collection("product_categories")
	productRepo := mongorepo.NewProductRepository(productColl)
	catRepo := mongorepo.NewProductCategoryRepository(catColl)
	CatService := NewProductCategoryService(catRepo)
	prodService := NewProductService(productRepo, CatService)

	t.Run("should return all products with right category", func(t *testing.T) {
		ctx := context.Background()
		cat1, _ := CatService.CreateProductCategory(ctx, "test_Category", "test_Description")
		cat2, _ := CatService.CreateProductCategory(ctx, "test_Category2", "test_Description2")

		prodService.CreateProduct(ctx, "test_product", "test_description", cat1.ID, 100, "test_brand", 10)
		prodService.CreateProduct(ctx, "test_product2", "test_description2", cat2.ID, 100, "test_brand2", 10)
		prodService.CreateProduct(ctx, "test_product3", "test_description3", cat1.ID, 100, "test_brand3", 10)

		products, err := prodService.GetAllProduct(ctx, cat1.ID)
		if err != nil {
			t.Fatalf("expected no error but got %v", err.Error())
		}
		if products == nil {
			t.Fatalf("expected products but got nil")
		}

		for _, product := range products {
			if product.Category.ID != cat1.ID {
				t.Fatalf("expected product category id to be %v but got %v", cat1.ID, product.Category.ID)
			}
			if product.Name != "test_product" && product.Name != "test_product3" {
				t.Fatalf("expected product name to be test_product or test_product3 but got %v", product.Name)
			}
		}

		t.Cleanup(func() {
			productColl.DeleteMany(ctx, bson.M{})
			catColl.DeleteMany(ctx, bson.M{})
		})
	})
}

func TestIntegation_UpdateProduct(t *testing.T) {
	catColl := testDb.Collection("product_categories")
	catrepo := mongorepo.NewProductCategoryRepository(catColl)
	catService := NewProductCategoryService(catrepo)
	prodColl := testDb.Collection("products")
	prodrepo := mongorepo.NewProductRepository(prodColl)
	prodService := NewProductService(prodrepo, catService)

	t.Run("should update product and only desired fields should be updated and save to mongo", func(t *testing.T) {
		ctx := context.Background()

		cat, _ := catService.CreateProductCategory(ctx, "test_Category", "test_Description")
		prod, _ := prodService.CreateProduct(ctx, "test_product", "test_description", cat.ID, 100, "test_brand", 10)
		newName := "test_product2"
		updatedProd ,err  := prodService.UpdateProduct(ctx, prod.ID, &newName, nil, nil, nil, nil, nil)
		if err != nil{
			log.Fatal(err)
		}

		if updatedProd == nil{
			t.Fatalf("expected product but got nil")
		}
		if updatedProd.Name != newName{
			t.Fatalf("expected product name to be %v but got %v", newName, updatedProd.Name)
		}
		if updatedProd.Description != prod.Description{
			t.Fatalf("expected product description to be %v but got %v", prod.Description, updatedProd.Description)
		}

		var dbResult bson.M
		objectId, err := bson.ObjectIDFromHex(prod.ID)
		if err != nil {
			log.Fatalf("failed to convert id to object id")
		}
		err = prodColl.FindOne(ctx, bson.M{"_id": objectId}).Decode(&dbResult)
		if err != nil {
			t.Fatalf("expected no error but got %v", err.Error())
		}
		if dbResult["name"] != newName {
			t.Fatalf("expected product name to be %v but got %v", newName, dbResult["name"])
		}

		t.Cleanup(func() {
			prodColl.DeleteMany(ctx, bson.M{})
			catColl.DeleteMany(ctx, bson.M{})
		})
	})
}

func TestIntegration_CreateProductCategory(t *testing.T){
	catColl := testDb.Collection("product_categories")
	catrepo := mongorepo.NewProductCategoryRepository(catColl)
	catService := NewProductCategoryService(catrepo)

	t.Run("should create product category and save to mongo", func(t *testing.T) {
		ctx := context.Background()

		cat , err := catService.CreateProductCategory(ctx , "test_Category" , "test_Description")

		if err != nil{
			t.Fatalf("expected no error but got %v", err.Error())
		}

		if cat == nil{
			t.Fatalf("expected product category but got nil")
		}

		var dbResult bson.M
		objectId, err := bson.ObjectIDFromHex(cat.ID)
		if err != nil {
			log.Fatalf("failed to convert id to object id")
		}
		err = catColl.FindOne(ctx, bson.M{"_id": objectId}).Decode(&dbResult)
		
		t.Cleanup(func() {
			catColl.DeleteMany(ctx, bson.M{})
		})

	})
}

func TestIntegration_UpdateProductCategory(t *testing.T){
	catColl := testDb.Collection("product_categories")
	catrepo := mongorepo.NewProductCategoryRepository(catColl)
	catService := NewProductCategoryService(catrepo)

	t.Run("should update product category and only desired fields should be updated and save to mongo", func(t *testing.T) {
		ctx := context.Background()

		cat , _ := catService.CreateProductCategory(ctx , "test_Category" , "test_Description")

		newName := "test_Category2"
		updatedCat , err := catService.UpdateProductCategory(ctx , cat.ID , &newName , nil)
		if err != nil{
			log.Fatal(err)
		}
		if updatedCat == nil{
			t.Fatalf("expected product category but got nil")
		}
		if updatedCat.Name != newName{
			t.Fatalf("expected product category name to be %v but got %v", newName, updatedCat.Name)
		}

		var dbResult bson.M
		objectId, err := bson.ObjectIDFromHex(cat.ID)
		if err != nil {
			log.Fatalf("failed to convert id to object id")
		}
		err = catColl.FindOne(ctx, bson.M{"_id": objectId}).Decode(&dbResult)
		if err != nil {
			t.Fatalf("expected no error but got %v", err.Error())
		}
		if dbResult["name"] != newName {
			t.Fatalf("expected product category name to be %v but got %v", newName, dbResult["name"])
		}

		t.Cleanup(func() {
			catColl.DeleteMany(ctx, bson.M{})
		})
	})

}