package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/fx"
)

func NewMongoClient(lc fx.Lifecycle, cfg *Config) (*mongo.Client, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb : %w", err)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			return client.Ping(pingCtx, nil)
		},
		OnStop: func(ctx context.Context) error {
			return client.Disconnect(ctx)
		},
	})

	return client, nil
}

func NewMongoDatabase(client *mongo.Client) *mongo.Database {
	return client.Database("inventory")
}

func NewProductCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection("products")
}

func NewCategoryCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection("product_categories")
}

var MongoModule = fx.Options(
	fx.Provide(
		NewMongoClient,
		NewMongoDatabase,
		fx.Annotate(
			NewProductCollection,
			fx.ResultTags(`name:"productCollection"`),
		),
		fx.Annotate(
			NewCategoryCollection,
			fx.ResultTags(`name:"categoryCollection"`),
		),
	),
)
