package order

import (
	"fmt"
	"purple/links/pkg/db"
)

type OrderValidator struct {
	db *db.DB
}

func NewOrderValidator(db *db.DB) *OrderValidator {
	return &OrderValidator{db: db}
}

func (v *OrderValidator) ValidateCreateRequest(req *OrderCreateRequest) error {
	if err := v.validateProductsExist(req); err != nil {
		return fmt.Errorf("products validation failed: %w", err)
	}

	return nil
}

func (v *OrderValidator) validateProductsExist(req *OrderCreateRequest) error {
	productIDs := make([]int64, len(req.Items))
	for i, item := range req.Items {
		productIDs[i] = item.ProductID
	}

	var count int64
	err := v.db.Table("product").
		Where("id IN (?) AND deleted_at is null", productIDs).
		Count(&count).Error

	if err != nil {
		return err
	}

	if int(count) != len(productIDs) {
		return fmt.Errorf("some products not found or not active")
	}

	return nil
}
