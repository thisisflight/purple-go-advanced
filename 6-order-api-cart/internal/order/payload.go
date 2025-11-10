package order

import (
	"time"

	"github.com/shopspring/decimal"
)

type OrderCreateRequest struct {
	Items  []OrderItemRequest `json:"items" validate:"required,min=1,dive"`
	Notes  string             `json:"notes"`
	UserID *uint              `json:"-"`
}

type OrderCreateResponse struct {
	ID        uint            `json:"id"`
	UserID    uint            `json:"user_id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Status    OrderStatusData `json:"status"`
	Notes     string          `json:"notes"`
}

type OrderDetailResponse struct {
	ID        uint              `json:"id"`
	UserID    uint              `json:"user_id"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Status    OrderStatusData   `json:"status"`
	Notes     string            `json:"notes"`
	Items     []OrderItemDetail `json:"items"`
}

type OrderListResponse struct {
	Total int64                 `json:"total"`
	Items []OrderDetailResponse `json:"items"`
}

type OrderItemRequest struct {
	ProductID int64 `json:"product_id" validate:"required,gt=0"`
	Quantity  int   `json:"quantity" validate:"omitempty,gt=0"`
	OrderID   *uint `json:"-"`
}

type OrderItemDetail struct {
	ID       uint               `json:"id"`
	Quantity int                `json:"quantity" validate:"omitempty,gt=0"`
	Price    decimal.Decimal    `json:"price"`
	Discount decimal.Decimal    `json:"discount"`
	Product  OrderProductDetail `json:"product"`
}

type OrderProductDetail struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
