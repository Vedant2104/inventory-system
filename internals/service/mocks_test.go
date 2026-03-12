package service

import (
	"context"

	"github.com/Vedant2104/inventory-system/internals/domain"
)

type mockProductRepository struct {
	CreateProductFn  func(ctx context.Context, product *domain.Product) (*domain.Product, error)
	GetAllProductFn  func(ctx context.Context, category string) ([]*domain.Product, error)
	GetProductByIdFn func(ctx context.Context, ID string) (*domain.Product, error)
	UpdateProductFn  func(ctx context.Context, product *domain.Product) error
	DeleteProductFn  func(ctx context.Context, ID string) error
	// BulkCreate     func(ctx context.Context, products []domain.Product) error
}

func (m *mockProductRepository) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	return m.CreateProductFn(ctx, product)
}

func (m *mockProductRepository) GetAllProduct(ctx context.Context, category string) ([]*domain.Product, error) {
	return m.GetAllProductFn(ctx, category)
}

func (m *mockProductRepository) GetProductById(ctx context.Context, ID string) (*domain.Product, error) {
	return m.GetProductByIdFn(ctx, ID)
}

func (m *mockProductRepository) UpdateProduct(ctx context.Context, product *domain.Product) error {
	return m.UpdateProductFn(ctx, product)
}

func (m *mockProductRepository) DeleteProduct(ctx context.Context, ID string) error {
	return m.DeleteProductFn(ctx, ID)
}

func (m *mockProductRepository) BulkCreate(ctx context.Context, products []domain.Product) error {
	return nil
}

type mockProductCategoryRepository struct {
	GetProductCategoryByIdFn func(ctx context.Context, ID string) (*domain.ProductCategory, error)
	GetAllProductCategoryFn  func(ctx context.Context) ([]*domain.ProductCategory, error)
	CreateProductCategoryFn  func(ctx context.Context, productCategory *domain.ProductCategory) (*domain.ProductCategory, error)
	UpdateProductCategoryFn  func(ctx context.Context, productCategory *domain.ProductCategory) (*domain.ProductCategory, error)
	DeleteProductCategoryFn  func(ctx context.Context, ID string) error
}

func (m *mockProductCategoryRepository) GetProductCategoryById(ctx context.Context, ID string) (*domain.ProductCategory, error) {
	return m.GetProductCategoryByIdFn(ctx, ID)
}

func (m *mockProductCategoryRepository) GetAllProductCategory(ctx context.Context) ([]*domain.ProductCategory, error) {
	return m.GetAllProductCategoryFn(ctx)
}

func (m *mockProductCategoryRepository) CreateProductCategory(ctx context.Context, productCategory *domain.ProductCategory) (*domain.ProductCategory, error) {
	return m.CreateProductCategoryFn(ctx, productCategory)
}

func (m *mockProductCategoryRepository) UpdateProductCategory(ctx context.Context, productCategory *domain.ProductCategory) (*domain.ProductCategory, error) {
	return m.UpdateProductCategoryFn(ctx, productCategory)
}

func (m *mockProductCategoryRepository) DeleteProductCategory(ctx context.Context, ID string) error {
	return m.DeleteProductCategoryFn(ctx, ID)
}
