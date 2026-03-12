package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Vedant2104/inventory-system/internals/domain"
)

func TestProductService_CreateProduct(t *testing.T) {
	mockProductRepo := &mockProductRepository{}
	mockProductCategoryRepo := &mockProductCategoryRepository{}
	mockProductCategoryService := NewProductCategoryService(mockProductCategoryRepo)
	ProductService := NewProductService(mockProductRepo, mockProductCategoryService)

	tests := []struct {
		test_name     string
		name          string
		description   string
		category      string
		price         int
		brand         string
		quantity      int
		setupMocks    func()
		wantErr       bool
		expectedError string
	}{
		{
			test_name:   "Success Scenario",
			name:        "laptop",
			description: "laptop description",
			category:    "laptop",
			price:       1000,
			brand:       "dell",
			quantity:    10,
			setupMocks: func() {
				mockProductRepo.CreateProductFn = func(ctx context.Context, p *domain.Product) (*domain.Product, error) {
					p.ID = "generatedId"
					return p, nil
				}
				mockProductCategoryRepo.GetProductCategoryByIdFn = func(ctx context.Context, ID string) (*domain.ProductCategory, error) {
					return &domain.ProductCategory{ID: ID}, nil
				}
			},
			wantErr: false,
		},
		{
			test_name:     "Negative Price Error",
			name:          "laptop",
			description:   "laptop description",
			category:      "laptop",
			price:         -1000,
			brand:         "dell",
			quantity:      10,
			setupMocks:    nil,
			wantErr:       true,
			expectedError: "price cannot be negative",
		},
		{
			test_name:     "Invalid product Name",
			name:          "la",
			category:      "laptop",
			setupMocks:    nil,
			wantErr:       true,
			expectedError: "length of name must be greater than 3",
		},
		{
			test_name:   "Invalid Category",
			name:        "laptop",
			description: "laptop description",
			category:    "laptop",
			price:       1000,
			brand:       "dell",
			quantity:    10,
			setupMocks: func() {
				mockProductCategoryRepo.GetProductCategoryByIdFn = func(ctx context.Context, ID string) (*domain.ProductCategory, error) {
					return nil, errors.New("category not found")
				}
			},
			wantErr:       true,
			expectedError: "category not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			if tt.setupMocks != nil {
				tt.setupMocks()
			}

			_, err := ProductService.CreateProduct(context.Background(), tt.name, tt.description, tt.category, tt.price, tt.brand, tt.quantity)

			if tt.wantErr {
				if err == nil {
					t.Errorf("%s: expected error but got nil", tt.test_name)
					return
				}
				if err.Error() != tt.expectedError {
					t.Errorf("%s: expected error [%v] but got [%v]", tt.test_name, tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("%s: expected no error but got [%v]", tt.test_name, err.Error())
				}
			}
		})
	}
}

func TestProductService_GetAllProducts(t *testing.T){
	mockProductRepo := &mockProductRepository{}
	mockCategotyRepo := &mockProductCategoryRepository{}
	catService := NewProductCategoryService(mockCategotyRepo)

	ProductService := NewProductService(mockProductRepo, catService)

	t.Run("Success Scenario", func(t *testing.T) {
		mockProductRepo.GetAllProductFn = func(ctx context.Context , category string) ([]*domain.Product, error) {
			return []*domain.Product{}, nil
		}

		_, err := ProductService.GetAllProduct(context.Background(), "category")
		if err != nil {
			t.Errorf("expected no error but got %v", err.Error())
		}
	})
}

func TestProductService_GetProductById(t *testing.T) {
	mockProductRepo := &mockProductRepository{}
	mockCategotyRepo := &mockProductCategoryRepository{}
	catService := NewProductCategoryService(mockCategotyRepo)

	ProductService := NewProductService(mockProductRepo, catService)

	t.Run("Success Scenario", func(t *testing.T) {
		mockProductRepo.GetProductByIdFn = func(ctx context.Context, ID string) (*domain.Product, error) {
			return &domain.Product{ID: ID}, nil
		}

		_, err := ProductService.GetProductById(context.Background(), "id")
		if err != nil {
			t.Errorf("expected no error but got %v", err.Error())
		}

	})
	t.Run("Invalid Id", func(t *testing.T) {
		id := ""

		_, err := ProductService.GetProductById(context.Background(), id)

		if err == nil {
			t.Errorf("expected error but got nil")
			return
		}

		if err.Error() != "invalid id" {
			t.Errorf("expected error %v got %v", "invalid product Id", err.Error())
		}

	})
}

func TestProductService_UpdateProduct(t *testing.T) {
	mockProductRepo := &mockProductRepository{}
	mockCategoryRepo := &mockProductCategoryRepository{}
	catService := NewProductCategoryService(mockCategoryRepo)

	ProductService := NewProductService(mockProductRepo, catService)
	mockCategoryId := "catid"
	tests := []struct {
		test_name     string
		id            string
		name          *string
		description   *string
		category      *string
		price         *int
		brand         *string
		quantity      *int
		setupMocks    func()
		wantErr       bool
		expectedError string
	}{
		{
			test_name: "Success Scenario",
			id:        "generatedId",
			category:  &mockCategoryId,
			setupMocks: func() {
				mockProductRepo.GetProductByIdFn = func(ctx context.Context, ID string) (*domain.Product, error) {
					return &domain.Product{ID: ID}, nil
				}
				mockCategoryRepo.GetProductCategoryByIdFn = func(ctx context.Context, ID string) (*domain.ProductCategory, error) {
					return &domain.ProductCategory{ID: ID}, nil
				}
				mockProductRepo.UpdateProductFn = func(ctx context.Context, p *domain.Product) error {
					return nil
				}
			},
		}, {
			test_name:     "Invalid Id",
			id:            "",
			wantErr:       true,
			expectedError: "invalid id",
		}, {
			test_name: "Invalid Category Id",
			id:        "generatedId",
			category:  &mockCategoryId,
			setupMocks: func() {
				mockProductRepo.GetProductByIdFn = func(ctx context.Context, ID string) (*domain.Product, error) {
					return &domain.Product{ID: ID}, nil
				}
				mockCategoryRepo.GetProductCategoryByIdFn = func(ctx context.Context, ID string) (*domain.ProductCategory, error) {
					return nil, nil
				}
			},
			wantErr:       true,
			expectedError: "category not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			if tt.setupMocks != nil {
				tt.setupMocks()
			}

			_, err := ProductService.UpdateProduct(context.Background(), tt.id, tt.name, tt.description, tt.category, tt.price, tt.brand, tt.quantity)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
					return
				}

				if err.Error() != tt.expectedError {
					t.Errorf("expected error %v got %v", tt.expectedError, err.Error())
				}
				return
			} else {
				if err != nil {
					t.Errorf("expected no error but got %v", err.Error())
				}
			}
		})
	}
}

func TestProductService_DeleteProduct(t *testing.T) {
	mockProductCategoryRepository := &mockProductCategoryRepository{}
	mockProductRepository := &mockProductRepository{}

	catService := NewProductCategoryService(mockProductCategoryRepository)
	ProductService := NewProductService(mockProductRepository, catService)

	tests := []struct {
		test_name     string
		id            string
		setupMocks    func()
		wantErr       bool
		expectedError string
	}{
		{
			test_name: "Succeess Scenario",
			id:        "testId",
			setupMocks: func() {
				mockProductRepository.GetProductByIdFn = func(ctx context.Context, ID string) (*domain.Product, error) {
					return &domain.Product{ID: ID}, nil
				}
				mockProductRepository.DeleteProductFn = func(ctx context.Context, ID string) error {
					return nil
				}
			},
			wantErr: false,
		}, {
			test_name:     "Invalid Id",
			id:            "",
			wantErr:       true,
			expectedError: "invalid id",
		}, {
			test_name: "Product Not Found",
			id:        "testId",
			setupMocks: func() {
				mockProductRepository.GetProductByIdFn = func(ctx context.Context, ID string) (*domain.Product, error) {
					return nil, nil
				}
			},
			wantErr:       true,
			expectedError: "product not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			if tt.setupMocks != nil {
				tt.setupMocks()
			}

			err := ProductService.DeleteProduct(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
					return
				}

				if err.Error() != tt.expectedError {
					t.Errorf("expected error %v got %v", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got %v", err.Error())
				}
			}
		})
	}
}
