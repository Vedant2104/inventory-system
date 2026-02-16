package service

import (
	"context"
	"errors"

	"github.com/Vedant2104/inventory-system/internals/domain"
	"github.com/Vedant2104/inventory-system/internals/ports"
)

type ProductService struct {
	repo ports.ProductRepository
}

func NewProductService(productRepository ports.ProductRepository) *ProductService {
	return &ProductService{repo: productRepository}
}

func (s *ProductService) CreateProduct(ctx context.Context, name string, description string, category string, price int, brand string, quantity int) (*domain.Product, error) {

	product, err := domain.NewProduct(name, description, category, price, brand, quantity)

	if err != nil {
		return nil, err
	}

	return s.repo.CreateProduct(ctx, product)
}

func (s *ProductService) GetAllProduct(ctx context.Context) ([]*domain.Product, error) {
	return s.repo.GetAllProduct(ctx)
}

func (s *ProductService) GetProductById(ctx context.Context, ID string) (*domain.Product, error) {
	return s.repo.GetProductById(ctx, ID)
}

func (s *ProductService) DeleteProduct(ctx context.Context, ID string) error {
	existing_product, _ := s.repo.GetProductById(ctx, ID)
	if existing_product == nil {
		return errors.New("product not found")
	}
	return s.repo.DeleteProduct(ctx, ID)
}

func (s *ProductService) UpdateProduct(ctx context.Context, id string, name string, description string, category string, price int, brand string, quantity int) (*domain.Product, error) {

	existing_product, _ := s.repo.GetProductById(ctx, id)
	if existing_product == nil {
		return nil, errors.New("product not found")
	}
	copy_product := *existing_product
	err := copy_product.UpdateProductValidation(&name, &description, &category, &price, &brand, &quantity)
	if err != nil {
		return nil, err
	}
	return s.repo.UpdateProduct(ctx, &copy_product)
}
