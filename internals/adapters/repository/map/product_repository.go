package maprepo

import (
	"context"
	"errors"

	"github.com/Vedant2104/inventory-system/internals/domain"
	"github.com/Vedant2104/inventory-system/internals/ports"
	"github.com/google/uuid"
)

type ProductRepository struct {
	products map[string]domain.Product
}

var _ ports.ProductRepository = (*ProductRepository)(nil)

func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		products: map[string]domain.Product{},
	}
}

func (p *ProductRepository) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	product.ID = uuid.New().String()
	p.products[product.ID] = *product
	return product, nil
}

func (p *ProductRepository) GetAllProduct(ctx context.Context) ([]*domain.Product, error) {
	products := []*domain.Product{}

	for _, product := range p.products {
		products = append(products, &product)
	}

	return products, nil
}

func (p *ProductRepository) GetProductById(ctx context.Context, Id string) (*domain.Product, error) {
	product, ok := p.products[Id]
	if !ok {
		return nil, errors.New("product not found")
	}

	return &product, nil
}

func (p *ProductRepository) UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {

	p.products[product.ID] = *product
	return product, nil

}

func (p *ProductRepository) DeleteProduct(ctx context.Context, Id string) error {

	delete(p.products, Id)
	return nil
}
