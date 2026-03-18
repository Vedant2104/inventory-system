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

type ProductCategoryRepository struct {
	collection *mongo.Collection
}

func NewProductCategoryRepository(collection *mongo.Collection) *ProductCategoryRepository {
	return &ProductCategoryRepository{collection: collection}
}

var _ ports.ProductCategoryRepository = (*ProductCategoryRepository)(nil)

type mongoProductCategory struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
}

func (p *ProductCategoryRepository) CreateProductCategory(ctx context.Context, productCategory *domain.ProductCategory) (*domain.ProductCategory, error) {
	doc := &mongoProductCategory{
		Name:        productCategory.Name,
		Description: productCategory.Description,
	}

	result, err := p.collection.InsertOne(ctx, doc)

	if err != nil {
		return nil, err
	}

	objectId, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return nil, errors.New("error fetching insertedId")
	}
	productCategory.ID = objectId.Hex()
	return productCategory, nil
}

func (p *ProductCategoryRepository) GetAllProductCategory(ctx context.Context) ([]*domain.ProductCategory, error) {

	filter := bson.D{{}}
	opt := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})

	cur, err := p.collection.Find(ctx, filter, opt)

	if err != nil {
		return nil, err
	}

	var categories []*domain.ProductCategory

	for cur.Next(ctx) {
		var doc *mongoProductCategory
		err := cur.Decode(&doc)
		if err != nil {
			return nil, err
		}
		category := &domain.ProductCategory{
			ID:          doc.ID.Hex(),
			Name:        doc.Name,
			Description: doc.Description,
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (p *ProductCategoryRepository) GetProductCategoryById(ctx context.Context, id string) (*domain.ProductCategory, error) {
	objectId, err := bson.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: objectId}}

	var doc *mongoProductCategory
	err = p.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return nil, err
	}

	return &domain.ProductCategory{
		ID:          doc.ID.Hex(),
		Name:        doc.Name,
		Description: doc.Description,
	}, nil
}

func (p *ProductCategoryRepository) UpdateProductCategory(ctx context.Context, productCategory *domain.ProductCategory) (*domain.ProductCategory, error) {
	id := productCategory.ID
	objectId , err := bson.ObjectIDFromHex(id)
	if err != nil{
		return nil ,err
	}
	filter := bson.D{{Key :"_id" , Value :objectId}}
	update := bson.M{"$set": bson.M{"name": productCategory.Name, "description": productCategory.Description}}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedCategory *mongoProductCategory
	err  = p.collection.FindOneAndUpdate(ctx , filter , update , opt).Decode(&updatedCategory)
	if err != nil{
		return nil ,err
	}
	return &domain.ProductCategory{ID: updatedCategory.ID.Hex(), Name: updatedCategory.Name, Description: updatedCategory.Description}, nil
}

func (p *ProductCategoryRepository) DeleteProductCategory(ctx context.Context, id string) error {
	objectID , err := bson.ObjectIDFromHex(id)
	if err != nil{
		return err
	}

	filter := bson.D{{Key :"_id" , Value :objectID}}

	_  ,err = p.collection.DeleteOne(ctx, filter)

	if err != nil{
		return err
	}
	return nil
}
