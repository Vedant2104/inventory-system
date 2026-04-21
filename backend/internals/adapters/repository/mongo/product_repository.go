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
	Category    bson.ObjectID `bson:"category"`
	Price       int           `bson:"price"`
	Brand       string        `bson:"brand"`
	Quantity    int           `bson:"quantity"`
}

type mongoResultDao struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
	Category    struct {
		ID          bson.ObjectID `bson:"_id,omitempty"`
		Name        string        `bson:"name"`
		Description string        `bson:"description"`
	} `bson:"category"`
	Price    int    `bson:"price"`
	Brand    string `bson:"brand"`
	Quantity int    `bson:"quantity"`
}

type ProductCountByCategoryDao struct {
	CategoryId bson.ObjectID `bson:"_id"`
	Category   string        `bson:"category"`
	Count      int           `bson:"count"`
}

type PriceSegmentationDao struct {
	CategoryId bson.ObjectID `bson:"category_id"`
	Category   string        `bson:"category"`
	Budget     int           `bson:"budget"`
	MidRange   int           `bson:"midRange"`
	Premium    int           `bson:"premium"`
}

func (p *ProductRepository) mapToDomain(product *mongoResultDao) *domain.Product {
	return &domain.Product{
		ID:          product.ID.Hex(),
		Name:        product.Name,
		Description: product.Description,
		Category: &domain.ProductCategory{
			ID:          product.Category.ID.Hex(),
			Name:        product.Category.Name,
			Description: product.Category.Description,
		},
		Price:    product.Price,
		Brand:    product.Brand,
		Quantity: product.Quantity,
	}
}

func (p *ProductRepository) mapFromDomain(product *domain.Product) (*mongoProduct, error) {
	catId, err := bson.ObjectIDFromHex(product.Category.ID)
	if err != nil {
		return nil, err
	}
	return &mongoProduct{
		Name:        product.Name,
		Description: product.Description,
		Category:    catId,
		Price:       product.Price,
		Brand:       product.Brand,
		Quantity:    product.Quantity,
	}, nil
}

func (p *ProductRepository) buildProductAggregation(filter bson.D) mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "product_categories"},
			{Key: "localField", Value: "category"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "category"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$category"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
	}
}

func (p *ProductRepository) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {

	doc, err := p.mapFromDomain(product)
	if err != nil {
		return nil, err
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

func (p *ProductRepository) GetAllProduct(ctx context.Context, category string) ([]*domain.Product, error) {

	filter := bson.D{{}}
	if category != "" {
		catId, err := bson.ObjectIDFromHex(category)
		if err != nil {
			return nil, err
		}
		filter = bson.D{{Key: "category", Value: catId}}
	}
	// opt := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})
	pipeline := p.buildProductAggregation(filter)
	cur, err := p.collection.Aggregate(ctx, pipeline)

	if err != nil {
		return nil, err
	}

	var products []*domain.Product

	for cur.Next(ctx) {
		var product mongoResultDao
		if err := cur.Decode(&product); err != nil {
			return nil, err
		}
		var doc = p.mapToDomain(&product)
		products = append(products, doc)
	}
	defer cur.Close(ctx)

	return products, nil
}

func (p *ProductRepository) GetProductById(ctx context.Context, id string) (*domain.Product, error) {

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid product Id")
	}

	filter := bson.D{{Key: "_id", Value: objectID}}
	pipeline := p.buildProductAggregation(filter)
	cur, err := p.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var product mongoResultDao
	if cur.Next(ctx) {
		if err := cur.Decode(&product); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("product not found")
	}

	return p.mapToDomain(&product), nil
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

func (p *ProductRepository) UpdateProduct(ctx context.Context, product *domain.Product) error {
	id := product.ID
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: objectId}}
	catId, err := bson.ObjectIDFromHex(product.Category.ID)
	if err != nil {
		return err
	}
	update := bson.M{"$set": bson.M{"name": product.Name, "description": product.Description, "category": catId, "price": product.Price, "brand": product.Brand, "quantity": product.Quantity}}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedProduct mongoProduct

	err = p.collection.FindOneAndUpdate(ctx, filter, update, opt).Decode(&updatedProduct)

	if err != nil {
		return err
	}
	return nil
}

func (p *ProductRepository) BulkCreate(ctx context.Context, products []domain.Product) error {
	var docs []any

	for _, product := range products {
		doc, err := p.mapFromDomain(&product)
		if err != nil {
			return err
		}
		docs = append(docs, doc)
	}

	_, err := p.collection.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductRepository) ReportLowStockedProducts(ctx context.Context, threshold int) ([]*domain.LowStockProducts, error) {
	filter := bson.D{{Key: "quantity", Value: bson.D{{Key: "$lt", Value: threshold}}}}

	pipeline := p.buildProductAggregation(filter)

	cur, err := p.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var products []*domain.LowStockProducts

	for cur.Next(ctx) {
		var mongoProduct mongoResultDao
		if err := cur.Decode(&mongoProduct); err != nil {
			return nil, err
		}
		product := &domain.LowStockProducts{
			Id:       mongoProduct.ID.Hex(),
			Product:  mongoProduct.Name,
			Brand:    mongoProduct.Brand,
			Quantity: mongoProduct.Quantity,
		}
		products = append(products, product)

	}

	return products, nil

}

func (p *ProductRepository) ReportProductCountByCategory(ctx context.Context, minValue int, maxValue int) ([]domain.ProductCountByCategory, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$category"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},

		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "product_categories"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "category_info"},
		}}},

		{{Key: "$unwind", Value: "$category_info"}},

		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: "$_id"},
			{Key: "category", Value: "$category_info.name"},
			{Key: "count", Value: 1},
		}}},
		{{Key: "$match", Value: bson.D{
			{Key: "count", Value: bson.D{
				{Key: "$gte", Value: minValue},
				{Key: "$lte", Value: maxValue},
			}},
		}}},
	}

	cur, err := p.collection.Aggregate(ctx, pipeline)

	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var result []domain.ProductCountByCategory

	for cur.Next(ctx) {
		var doc ProductCountByCategoryDao
		err := cur.Decode(&doc)
		if err != nil {
			return nil, err
		}
		result = append(result, domain.ProductCountByCategory{
			CategoryId: doc.CategoryId.Hex(),
			Category:   doc.Category,
			Count:      doc.Count,
		})
	}

	return result, nil
}

func (p *ProductRepository) ReportPriceSegmentation(ctx context.Context) ([]domain.PriceSegmentation, error) {
	pipeline := mongo.Pipeline{
		// Stage 1: Lookup
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "product_categories"},
			{Key: "localField", Value: "category"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "cat_info"},
		}}},

		// Stage 2: Unwind
		bson.D{{Key: "$unwind", Value: "$cat_info"}},

		// Stage 3: Group
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$cat_info._id"},
			{Key: "name", Value: bson.D{{Key: "$first", Value: "$cat_info.name"}}},
			{Key: "budget", Value: bson.D{{Key: "$sum", Value: bson.D{
				{Key: "$cond", Value: bson.A{
					bson.D{{Key: "$lt", Value: bson.A{"$price", 2000}}},
					1, 0,
				}},
			}}}},
			{Key: "midrange", Value: bson.D{{Key: "$sum", Value: bson.D{
				{Key: "$cond", Value: bson.A{
					bson.D{{Key: "$and", Value: bson.A{
						bson.D{{Key: "$gt", Value: bson.A{"$price", 2000}}},
						bson.D{{Key: "$lt", Value: bson.A{"$price", 5000}}},
					}}},
					1, 0,
				}},
			}}}},
			{Key: "premium", Value: bson.D{{Key: "$sum", Value: bson.D{
				{Key: "$cond", Value: bson.A{
					bson.D{{Key: "$gt", Value: bson.A{"$price", 5000}}},
					1, 0,
				}},
			}}}},
			{Key: "total", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},

		// Stage 4: Project
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "category_id", Value: "$_id"},
			{Key: "category", Value: "$name"},
			{Key: "budget", Value: 1},
			{Key: "midRange", Value: "$midrange"},
			{Key: "premium", Value: 1},
			{Key: "_id", Value: 0},
		}}},
	}

	cur, err := p.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var docs []PriceSegmentationDao

	err = cur.All(ctx, &docs)

	if err != nil {
		return nil, err
	}

	var result []domain.PriceSegmentation
	for _, doc := range docs {
		result = append(result, domain.PriceSegmentation{
			CategoryId: doc.CategoryId.Hex(),
			Category:   doc.Category,
			Budget:     doc.Budget,
			MidRange:   doc.MidRange,
			Premium:    doc.Premium,
		})
	}
	return result, nil

}
