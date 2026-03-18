package service

import (
	"context"
	"errors"

	"github.com/Vedant2104/inventory-system/internals/domain"
	"github.com/Vedant2104/inventory-system/internals/ports"
)

type ProductCategoryService struct {
	repo ports.ProductCategoryRepository
}

func NewProductCategoryService(repo ports.ProductCategoryRepository) *ProductCategoryService {
	return &ProductCategoryService{repo: repo}
}

func (s *ProductCategoryService) CreateProductCategory(ctx context.Context, name string, description string) (*domain.ProductCategory, error) {
	category, err := domain.NewProductCategory(name, description)
	if err != nil {
		return nil, err
	}
	return s.repo.CreateProductCategory(ctx, category)
}

func (s *ProductCategoryService) GetProductCategoryById(ctx context.Context, ID string) (*domain.ProductCategory, error) {
	if ID == "" {
		return nil, errors.New("invalid id")
	}
	return s.repo.GetProductCategoryById(ctx, ID)
}

func (s *ProductCategoryService) GetAllProductCategory(ctx context.Context) ([]*domain.ProductCategory, error) {
	return s.repo.GetAllProductCategory(ctx)
}

func (s *ProductCategoryService) DeleteProductCategory(ctx context.Context, ID string) error {
	if ID == "" {
		return errors.New("invalid id")
	}
	existing_product, _ := s.repo.GetProductCategoryById(ctx, ID)
	if existing_product == nil {
		return errors.New("product category not found")
	}
	return s.repo.DeleteProductCategory(ctx, ID)
}

func (s *ProductCategoryService) UpdateProductCategory(ctx context.Context, id string, name *string, description *string) (*domain.ProductCategory, error) {
	if id == "" {
		return nil, errors.New("invalid id")
	}
	existing_product, _ := s.repo.GetProductCategoryById(ctx, id)
	if existing_product == nil {
		return nil, errors.New("product category not found")
	}
	copy_categoty := *existing_product
	err := copy_categoty.UpdateProductCategoryValidation(name, description)
	if err != nil {
		return nil, err
	}

	return s.repo.UpdateProductCategory(ctx, &copy_categoty)
}
