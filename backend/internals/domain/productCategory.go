package domain

import (
	"errors"
	"strings"
)

type ProductCategory struct {
	ID          string
	Name        string
	Description string
}

func NewProductCategory(name string, description string) (*ProductCategory, error) {
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)

	if len(name) < 3 {
		return nil, errors.New("length of name should be greater than 3")
	}
	if len(description) < 3 {
		return nil, errors.New("length of description should be greater than 3")
	}

	return &ProductCategory{
		ID:          "",
		Name:        name,
		Description: description,
	}, nil
}

func (p *ProductCategory) UpdateProductCategoryValidation(name *string, description *string) error {

	if name != nil {
		*name = strings.TrimSpace(*name)
		if len(*name) < 3 {
			return errors.New("length of name should be greater than 3")
		}
		p.Name = *name
	}

	if description != nil {
		*description = strings.TrimSpace(*description)
		if len(*description) < 3 {
			return errors.New("length of description should be greater than 3")
		}
		p.Description = *description
	}
	return nil
}
