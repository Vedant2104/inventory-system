package ports

import (
	"context"

	"github.com/Vedant2104/inventory-system/internals/domain"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	GetAllProduct(ctx context.Context, category string) ([]*domain.Product, error)
	GetProductById(ctx context.Context, ID string) (*domain.Product, error)
	UpdateProduct(ctx context.Context, product *domain.Product) error
	DeleteProduct(ctx context.Context, ID string) error
	BulkCreate(ctx context.Context, products []domain.Product) error
	ReportLowStockedProducts(ctx context.Context, threshold int) ([]*domain.LowStockProducts, error)
	ReportProductCountByCategory(ctx context.Context, minValue int, maxValue int) ([]domain.ProductCountByCategory, error)
	ReportPriceSegmentation(ctx context.Context) ([]domain.PriceSegmentation, error)
}
