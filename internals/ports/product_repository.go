package ports

import (
	"context"

	"github.com/Vedant2104/inventory-system/internals/domain"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	GetAllProduct(ctx context.Context) ([]*domain.Product, error)
	GetProductById(ctx context.Context, ID string) (*domain.Product, error)
	UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	DeleteProduct(ctx context.Context, ID string) error
}
