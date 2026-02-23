package mongorepo

import (
	"context"
	"errors"

	"github.com/Vedant2104/inventory-system/internals/domain"
	"github.com/Vedant2104/inventory-system/internals/ports"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ProductRepository struct {
	collection *mongo.Collection
}

var _ ports.ProductRepository = (*ProductRepository)(nil)

func NewProductRepository(collection *mongo.Collection) *ProductRepository {
	return &ProductRepository{collection: collection}
}

type mongoProduct struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
	Category    string        `bson:"category"`
	Price       int           `bson:"price"`
	Brand       string        `bson:"brand"`
	Quantity    int           `bson:"s"`
}

func (p *ProductRepository) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	doc := mongoProduct{
		Name:        product.Name,
		Description: product.Description,
		Category:    product.Category,
		Price:       product.Price,
		Brand:       product.Brand,
		Quantity:    product.Quantity,
	}

	result, err := p.collection.InsertOne(ctx, doc)

	if err != nil {
		return nil, err
	}

	ObjectId, ok := result.InsertedID.(bson.ObjectID)

	if !ok {
		return nil, errors.New("failed to convert to inserted id")
	}

	product.ID = ObjectId.Hex()
	return product, nil
}

func (p *ProductRepository) GetAllProduct(ctx context.Context) ([]*domain.Product, error) {

	filter := bson.D{{}}

	opt := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})

	cur, err := p.collection.Find(ctx, filter, opt)

	if err != nil {
		return nil, err
	}

	var products []*domain.Product

	for cur.Next(ctx) {
		var product mongoProduct
		if err := cur.Decode(&product); err != nil {
			return nil, err
		}
		var doc = &domain.Product{
			ID:          product.ID.Hex(),
			Name:        product.Name,
			Description: product.Description,
			Category:    product.Category,
			Price:       product.Price,
			Brand:       product.Brand,
			Quantity:    product.Quantity,
		}
		products = append(products, doc)
	}

	return products, nil
}

func (p *ProductRepository) GetProductById(ctx context.Context, id string) (*domain.Product, error) {

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid product Id")
	}

	filter := bson.D{{Key: "_id", Value: objectID}}
	var product mongoProduct
	err = p.collection.FindOne(ctx, filter).Decode(&product)

	if err != nil {
		return nil, err
	}

	return &domain.Product{
		ID:          product.ID.Hex(),
		Name:        product.Name,
		Description: product.Description,
		Category:    product.Category,
		Price:       product.Price,
		Brand:       product.Brand,
		Quantity:    product.Quantity,
	}, nil
}

func (p *ProductRepository) DeleteProduct(ctx context.Context, id string) error {
	objectId, err := bson.ObjectIDFromHex(id)

	if err != nil {
		return errors.New("invalid product Id")
	}

	filter := bson.D{{Key: "_id", Value: objectId}}

	_, err = p.collection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}
	return nil
}

func (p *ProductRepository) UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	id := product.ID
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: objectId}}
	update := bson.M{"$set": bson.M{"name": product.Name, "description": product.Description, "category": product.Category, "price": product.Price, "brand": product.Brand, "quantity": product.Quantity}}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedProduct mongoProduct

	err = p.collection.FindOneAndUpdate(ctx, filter, update, opt).Decode(&updatedProduct)

	if err != nil {
		return nil, err
	}
	return &domain.Product{
		ID:          updatedProduct.ID.Hex(),
		Name:        updatedProduct.Name,
		Description: updatedProduct.Description,
		Category:    updatedProduct.Category,
		Price:       updatedProduct.Price,
		Brand:       updatedProduct.Brand,
		Quantity:    updatedProduct.Quantity,
	}, nil
}
