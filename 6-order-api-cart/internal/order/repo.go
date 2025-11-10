package order

import (
	"context"
	"purple/links/internal/user"
	"purple/links/pkg/db"
)

type OrderRepository struct {
	db *db.DB
}

type OrderItemRepository struct {
	db *db.DB
}

func NewOrderRepository(db *db.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func NewOrderItemRepository(db *db.DB) *OrderItemRepository {
	return &OrderItemRepository{
		db: db,
	}
}

func (or *OrderRepository) Create(order *Order) (*Order, error) {
	result := or.db.Create(order)
	if result.Error != nil {
		return nil, result.Error
	}
	return order, nil
}

func (oir *OrderItemRepository) Create(order_items []OrderItem) ([]OrderItem, error) {
	result := oir.db.Create(&order_items)
	if result.Error != nil {
		return nil, result.Error
	}
	return order_items, nil
}

func (or *OrderRepository) CreateWithItems(
	ctx context.Context,
	req *OrderCreateRequest,
	user *user.User) (*Order, error) {

	order := &Order{
		Notes:  req.Notes,
		UserID: user.ID,
	}

	tx := or.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	orderItems := make([]OrderItem, len(req.Items))
	for i, itemReq := range req.Items {
		orderItems[i] = OrderItem{
			OrderID:   order.ID,
			ProductID: uint(itemReq.ProductID),
			Quantity:  itemReq.Quantity,
		}
	}

	if err := tx.Create(&orderItems).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return order, nil
}

func (or *OrderRepository) GetOrder(orderID uint64) (*Order, error) {
	var order Order
	err := or.db.
		Preload("OrderItems").
		Preload("OrderItems.Product").
		First(&order, orderID).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (or *OrderRepository) GetListByUserID(userID uint) ([]Order, error) {
	orders := make([]Order, 20)

	err := or.db.
		Where("user_id = ?", userID).
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Order("created_at DESC").
		Find(&orders).Error

	if err != nil {
		return nil, err
	}
	return orders, nil
}
