package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/Vedant2104/inventory-system/internals/domain"
	"github.com/Vedant2104/inventory-system/internals/ports"
)

type ProductService struct {
	repo            ports.ProductRepository
	categoryService *ProductCategoryService
}

func NewProductService(productRepository ports.ProductRepository, categoryService *ProductCategoryService) *ProductService {
	return &ProductService{repo: productRepository, categoryService: categoryService}
}

func (s *ProductService) CreateProduct(ctx context.Context, name string, description string, category string, price int, brand string, quantity int) (*domain.Product, error) {

	cat, err := s.categoryService.GetProductCategoryById(ctx, category)
	if err != nil {
		return nil, err
	}
	if cat == nil{
		return nil , errors.New("category not found")
	}

	product, err := domain.NewProduct(name, description, cat, price, brand, quantity)

	if err != nil {
		return nil, err
	}

	return s.repo.CreateProduct(ctx, product)
}

func (s *ProductService) GetAllProduct(ctx context.Context , category string) ([]*domain.Product, error) {
	return s.repo.GetAllProduct(ctx , category)
}

func (s *ProductService) GetProductById(ctx context.Context, ID string) (*domain.Product, error) {
	if ID == "" {
		return nil, errors.New("invalid id")
	}
	return s.repo.GetProductById(ctx, ID)
}

func (s *ProductService) DeleteProduct(ctx context.Context, ID string) error {
	if ID == "" {
		return errors.New("invalid id")
	}
	existing_product, _ := s.repo.GetProductById(ctx, ID)
	if existing_product == nil {
		return errors.New("product not found")
	}
	return s.repo.DeleteProduct(ctx, ID)
}

func (s *ProductService) UpdateProduct(ctx context.Context, id string, name *string, description *string, category *string, price *int, brand *string, quantity *int) (*domain.Product, error) {
	if id == "" {
		return nil, errors.New("invalid id")
	}
	existing_product, _ := s.repo.GetProductById(ctx, id)
	if existing_product == nil {
		return nil, errors.New("product not found")
	}
	copy_product := *existing_product
	var cat *domain.ProductCategory
	if category != nil {
		cat, _ = s.categoryService.GetProductCategoryById(ctx, *category)
		if cat == nil {
			return nil, errors.New("category not found")
		}
	}
	err := copy_product.UpdateProductValidation(name, description, cat, price, brand, quantity)
	if err != nil {
		return nil, err
	}
	if s.repo.UpdateProduct(ctx, &copy_product) != nil {
		return nil, err
	}
	return &copy_product, nil
}

func (s *ProductService) BulkCreate(ctx context.Context, records [][]string) error {
	var products []domain.Product

	for i, row := range records {
		if i == 0 {
			continue
		}
		if len(row) != 6 {
			return errors.New("invalid record")
		}
		price, _ := strconv.Atoi(row[3])
		quantity, _ := strconv.Atoi(row[5])
		category, _ := s.categoryService.GetProductCategoryById(ctx, row[2])
		doc := domain.Product{
			Name:        row[0],
			Description: row[1],
			Category:    category,
			Price:       price,
			Brand:       row[4],
			Quantity:    quantity,
		}
		products = append(products, doc)
	}

	return s.repo.BulkCreate(ctx, products)
}
