package product

import (
	"fmt"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type ProductCreateRequest struct {
	Name        string          `json:"name" validate:"required"`
	Description string          `json:"description"`
	Images      pq.StringArray  `json:"images"`
	IsActive    bool            `json:"is_active"`
	Price       decimal.Decimal `json:"price"`
}

type ProductUpdateRequest struct {
	Name        *string          `json:"name"`
	Description *string          `json:"description"`
	Images      *pq.StringArray  `json:"images"`
	IsActive    *bool            `json:"is_active"`
	Price       *decimal.Decimal `json:"price"`
}

func (r *ProductCreateRequest) Validate() error {
	if len(r.Name) > 255 {
		return fmt.Errorf("name cannot be longer than 255 characters")
	}

	if r.Price.LessThan(decimal.Zero) {
		return fmt.Errorf("price must be greater than or equal to 0")
	}

	if len(r.Description) > 5000 {
		return fmt.Errorf("description cannot be longer than 5000 characters")
	}

	if len(r.Images) > 10 {
		return fmt.Errorf("cannot have more than 10 images")
	}
	for i, img := range r.Images {
		if img == "" {
			return fmt.Errorf("image URL cannot be empty at position %d", i)
		}
	}

	return nil
}

func (r *ProductUpdateRequest) Validate() error {
	if r.Name != nil {
		if *r.Name == "" {
			return fmt.Errorf("name cannot be empty")
		}
		if len(*r.Name) > 255 {
			return fmt.Errorf("name cannot be longer than 255 characters")
		}
	}

	if r.Price != nil {
		if r.Price.LessThan(decimal.Zero) {
			return fmt.Errorf("price must be greater than or equal to 0")
		}
	}

	if r.Description != nil && len(*r.Description) > 5000 {
		return fmt.Errorf("description is too long")
	}

	if r.Images != nil {
		if len(*r.Images) > 10 {
			return fmt.Errorf("cannot have more than 10 images")
		}
		for i, img := range *r.Images {
			if img == "" {
				return fmt.Errorf("image URL cannot be empty at position %d", i)
			}
		}
	}

	return nil
}

func (r *ProductUpdateRequest) ApplyUpdates(product *Product) {
	if r.Name != nil {
		product.Name = *r.Name
	}
	if r.Description != nil {
		product.Description = *r.Description
	}
	if r.Images != nil {
		product.Images = *r.Images
	}
	if r.IsActive != nil {
		product.IsActive = *r.IsActive
	}
	if r.Price != nil {
		product.Price = *r.Price
	}
}
