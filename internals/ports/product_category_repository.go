package ports

import (
	"context"

	"github.com/Vedant2104/inventory-system/internals/domain"
)

type ProductCategoryRepository interface {
	GetAllProductCategory(ctx context.Context) ([] *domain.ProductCategory , error)
	GetProductCategoryById(ctx context.Context , id string) (*domain.ProductCategory , error)
	CreateProductCategory(ctx context.Context , productCategory *domain.ProductCategory) (*domain.ProductCategory , error)
	UpdateProductCategory(ctx context.Context , productCategory *domain.ProductCategory) (*domain.ProductCategory , error)
	DeleteProductCategory(ctx context.Context , id string) error
}