package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Vedant2104/inventory-system/internals/domain"
)

func TestProductCategoryService_CreateProductCategory(t *testing.T) {
	mockProductCatRepo := &mockProductCategoryRepository{}

	catService := NewProductCategoryService(mockProductCatRepo)

	tests := []struct {
		test_name     string
		name          string
		description   string
		setupMocks    func()
		wantErr       bool
		expectedError string
	}{
		{
			test_name:   "Success Scenario",
			name:        "testCategory",
			description: "testDescription",
			setupMocks: func() {
				mockProductCatRepo.CreateProductCategoryFn = func(ctx context.Context, p *domain.ProductCategory) (*domain.ProductCategory, error) {
					p.ID = "testId"
					return p, nil
				}
			},
			wantErr: false,
		}, {
			test_name:     "Invalid Name",
			name:          "",
			wantErr:       true,
			expectedError: "length of name should be greater than 3",
		}, {
			test_name:     "Invalid Description",
			name:          "testCategory",
			description:   "",
			wantErr:       true,
			expectedError: "length of description should be greater than 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			if tt.setupMocks != nil {
				tt.setupMocks()
			}

			_, err := catService.CreateProductCategory(context.Background(), tt.name, tt.description)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
					return
				}

				if err.Error() != tt.expectedError {
					t.Errorf("expected error %v got %v", tt.expectedError, err.Error())
					return
				}
			} else if err != nil {
				t.Errorf("expected no error but got %v", err.Error())
			}

		})
	}
}

func TestProductCategoryService_GetProductCategoryById(t *testing.T) {
	mockProductCatRepo := &mockProductCategoryRepository{}
	catService := NewProductCategoryService(mockProductCatRepo)

	t.Run("Success Scenario", func(t *testing.T) {
		mockProductCatRepo.GetProductCategoryByIdFn = func(ctx context.Context, ID string) (*domain.ProductCategory, error) {
			return &domain.ProductCategory{ID: ID}, nil
		}

		_, err := catService.GetProductCategoryById(context.Background(), "id")
		if err != nil {
			t.Errorf("expected no error but got %v", err.Error())
		}
	})

	t.Run("Invalid Id", func(t *testing.T) {
		id := ""

		_, err := catService.GetProductCategoryById(context.Background(), id)

		if err == nil {
			t.Errorf("expected error but got nil")
			return
		}
	})
}

func TestProductCategoryService_GetAllProductCategory(t *testing.T) {
	mockProductCatRepo := &mockProductCategoryRepository{}
	catService := NewProductCategoryService(mockProductCatRepo)

	t.Run("Success Scenario", func(t *testing.T) {
		mockProductCatRepo.GetAllProductCategoryFn = func(ctx context.Context) ([]*domain.ProductCategory, error) {
			return []*domain.ProductCategory{}, nil
		}

		_, err := catService.GetAllProductCategory(context.Background())
		if err != nil {
			t.Errorf("expected no error but got %v", err.Error())
		}
	})
}

func TestProductCategoryService_DeleteProductCategory(t *testing.T) {
	mockProductCatRepo := &mockProductCategoryRepository{}
	catService := NewProductCategoryService(mockProductCatRepo)

	tests := []struct {
		test_name     string
		id            string
		setupMocks    func()
		wantErr       bool
		expectedError string
	}{
		{
			test_name: "Success Scenario",
			id:        "testId",
			setupMocks: func() {
				mockProductCatRepo.GetProductCategoryByIdFn = func(ctx context.Context, ID string) (*domain.ProductCategory, error) {
					return &domain.ProductCategory{ID: ID}, nil
				}
				mockProductCatRepo.DeleteProductCategoryFn = func(ctx context.Context, ID string) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			test_name:     "Invalid Id",
			id:            "",
			wantErr:       true,
			expectedError: "invalid id",
		},
		{
			test_name: "Product Category Not Found",
			id:        "testId",
			setupMocks: func() {
				mockProductCatRepo.GetProductCategoryByIdFn = func(ctx context.Context, ID string) (*domain.ProductCategory, error) {
					return nil, errors.New("product category not found")
				}
			},
			wantErr:       true,
			expectedError: "product category not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			if tt.setupMocks != nil {
				tt.setupMocks()
			}

			err := catService.DeleteProductCategory(context.Background(), tt.id)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
					return
				}

				if err.Error() != tt.expectedError {
					t.Errorf("expected error %v got %v", tt.expectedError, err.Error())
					return
				}
			} else if err != nil {
				t.Errorf("expected no error but got %v", err.Error())
			}
		})
	}
}

func TestProductCategoryService_UpdateProductCategory(t *testing.T) {
	mockProductCatRepo := &mockProductCategoryRepository{}
	catService := NewProductCategoryService(mockProductCatRepo)

	validName := "testName"
	validDescription := "testDescription"
	invalidName := "te"
	invalidDescription := "de"
	tests := []struct {
		test_name     string
		id            string
		name          *string
		description   *string
		setupMocks    func()
		wantErr       bool
		expectedError string
	}{
		{
			test_name:   "Success Scenario",
			id:          "testId",
			name:        &validName,
			description: &validDescription,
			setupMocks: func() {
				mockProductCatRepo.GetProductCategoryByIdFn = func(ctx context.Context, ID string) (*domain.ProductCategory, error) {
					return &domain.ProductCategory{ID: ID}, nil
				}
				mockProductCatRepo.UpdateProductCategoryFn = func(ctx context.Context, productCategory *domain.ProductCategory) (*domain.ProductCategory, error) {
					return nil, nil
				}
			},
			wantErr: false,
		},
		{
			test_name:     "Invalid Id",
			id:            "",
			wantErr:       true,
			expectedError: "invalid id",
		}, {
			test_name: "Product Category Not Found",
			id:        "testId",
			setupMocks: func() {
				mockProductCatRepo.GetProductCategoryByIdFn = func(ctx context.Context, ID string) (*domain.ProductCategory, error) {
					return nil, errors.New("product category not found")
				}
			},
			wantErr:       true,
			expectedError: "product category not found",
		},
		{
			test_name: "Invalid Name",
			id:        "testId",
			name:      &invalidName,
			setupMocks: func() {
				mockProductCatRepo.GetProductCategoryByIdFn = func(ctx context.Context, ID string) (*domain.ProductCategory, error) {
					return &domain.ProductCategory{ID: ID}, nil
				}
			},
			wantErr:       true,
			expectedError: "length of name should be greater than 3",
		},
		{
			test_name:   "Invalid Description",
			id:          "testId",
			description: &invalidDescription,
			setupMocks: func() {
				mockProductCatRepo.GetProductCategoryByIdFn = func(ctx context.Context, ID string) (*domain.ProductCategory, error) {
					return &domain.ProductCategory{ID: ID}, nil
				}
			},
			wantErr:       true,
			expectedError: "length of description should be greater than 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			if tt.setupMocks != nil {
				tt.setupMocks()
			}

			_, err := catService.UpdateProductCategory(context.Background(), tt.id, tt.name, tt.description)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
					return
				}

				if err.Error() != tt.expectedError {
					t.Errorf("expected error %v got %v", tt.expectedError, err.Error())
					return
				}
			} else if err != nil {
				t.Errorf("expected no error but got %v", err.Error())
			}
		})
	}
}
