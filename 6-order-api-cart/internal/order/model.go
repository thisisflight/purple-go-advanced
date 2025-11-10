package order

import (
	"fmt"
	"purple/links/internal/product"
	"purple/links/internal/user"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type OrderStatus uint8

type OrderStatusData struct {
	StatusID uint8  `json:"status_id"`
	Name     string `json:"name"`
}

const (
	OrderStatusPending    OrderStatus = 0
	OrderStatusConfirmed  OrderStatus = 1
	OrderStatusProcessing OrderStatus = 2

	OrderStatusShipped   OrderStatus = 10
	OrderStatusDelivered OrderStatus = 11
	OrderStatusCompleted OrderStatus = 12

	OrderStatusCancelled OrderStatus = 20
	OrderStatusRefunded  OrderStatus = 21
	OrderStatusFailed    OrderStatus = 22
)

type Order struct {
	gorm.Model
	UserID     uint        `json:"user_id" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;index"`
	Status     OrderStatus `json:"status" gorm:"default:0"`
	Notes      string      `json:"notes" gorm:"type:text"`
	OrderItems []OrderItem `json:"order_items"`

	User user.User `json:"user" gorm:"foreignKey:UserID;references:ID"`
}

func (Order) TableName() string {
	return "order"
}

type OrderItem struct {
	gorm.Model
	OrderID   uint            `json:"order_id" gorm:"not null;index"`
	ProductID uint            `json:"product_id" gorm:"not null;index"`
	Quantity  int             `json:"quantity" gorm:"not null;default:1"`
	Price     decimal.Decimal `json:"price" gorm:"type:decimal(17,2);default:0"`
	Discount  decimal.Decimal `json:"discount" gorm:"type:decimal(17,2);default:0"`

	Order   Order           `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Product product.Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

func (OrderItem) TableName() string {
	return "order_item"
}

func (o *Order) GetOrderStatusData() OrderStatusData {
	var statusName string
	switch o.Status {
	case OrderStatusPending:
		statusName = "pending"
	case OrderStatusConfirmed:
		statusName = "confirmed"
	case OrderStatusProcessing:
		statusName = "processing"
	case OrderStatusShipped:
		statusName = "shipped"
	case OrderStatusDelivered:
		statusName = "delivered"
	case OrderStatusCompleted:
		statusName = "completed"
	case OrderStatusCancelled:
		statusName = "cancelled"
	case OrderStatusRefunded:
		statusName = "refunded"
	case OrderStatusFailed:
		statusName = "failed"
	default:
		statusName = "unknown"
	}
	return OrderStatusData{
		StatusID: uint8(o.Status),
		Name:     statusName,
	}
}

func (oi *OrderItem) GetSubtotal() decimal.Decimal {
	subtotal := oi.Price.Mul(decimal.NewFromInt(int64(oi.Quantity)))
	return subtotal.Sub(oi.Discount)
}

func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if oi.Price.IsZero() {
		var product product.Product
		if err := tx.First(&product, oi.ProductID).Error; err != nil {
			return fmt.Errorf("failed to get product price: %w", err)
		}
		oi.Price = product.Price
	}
	return nil
}
