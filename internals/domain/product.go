package domain

import (
	"errors"
	"strings"
)

var (
	ErrInvalidID          = errors.New("id cannot be empty")
	ErrInvalidName        = errors.New("length of name must be greater than 3")
	ErrInvalidDescription = errors.New("length of description must be greater than 5")
	ErrInvalidCategory    = errors.New("length of category must be greater than 3")
	ErrInvalidPrice       = errors.New("price cannot be negative")
	ErrInvalidBrand       = errors.New("length of brand must be greater than 3")
	ErrInvalidQuantity    = errors.New("quantity cannot be negative")
)

type Product struct {
	ID          string
	Name        string
	Description string
	Category    string
	Price       int
	Brand       string
	Quantity    int
}

func NewProduct(Name string, Description string, Category string, Price int, Brand string, Quantity int) (*Product, error) {
	Name = strings.TrimSpace(Name)
	Description = strings.TrimSpace(Description)
	Category = strings.TrimSpace(Category)
	Brand = strings.TrimSpace(Brand)

	if len(Name) < 3 {
		return nil, ErrInvalidName
	}
	if len(Description) < 5 {
		return nil, ErrInvalidDescription
	}
	if len(Category) < 3 {
		return nil, ErrInvalidCategory
	}
	if Price < 0 {
		return nil, ErrInvalidPrice
	}
	if len(Brand) < 3 {
		return nil, ErrInvalidBrand
	}
	if Quantity < 0 {
		return nil, ErrInvalidQuantity
	}

	return &Product{
		ID:          "",
		Name:        Name,
		Description: Description,
		Category:    Category,
		Price:       Price,
		Brand:       Brand,
		Quantity:    Quantity,
	}, nil
}

func (p *Product) UpdateProductValidation(name *string, description *string, category *string, price *int, brand *string, quantity *int) error {

	if name != nil {
		*name = strings.TrimSpace(*name)
		if len(*name) < 3 {
			return ErrInvalidName
		}
		p.Name = *name
	}

	if description != nil {
		*description = strings.TrimSpace(*description)
		if len(*description) < 5 {
			return ErrInvalidDescription
		}
		p.Description = *description
	}

	if category != nil {
		*category = strings.TrimSpace(*category)
		if len(*category) < 3 {
			return ErrInvalidCategory
		}
		p.Category = *category
	}

	if price != nil {
		if *price < 0 {
			return ErrInvalidPrice
		}
		p.Price = *price
	}

	if brand != nil {
		*brand = strings.TrimSpace(*brand)
		if len(*brand) < 3 {
			return ErrInvalidBrand
		}
		p.Brand = *brand
	}

	if quantity != nil {
		if *quantity < 0 {
			return ErrInvalidQuantity
		}
		p.Quantity = *quantity
	}

	return nil

}
