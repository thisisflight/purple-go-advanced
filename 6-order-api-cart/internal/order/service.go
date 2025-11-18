package order

import (
	"context"
	"errors"
	"purple/links/internal/user"
	"purple/links/pkg/constants"
)

type OrderServiceDeps struct {
	OrderRepository IOrderRepository
	OrderValidator  IOrderValidator
	UserRepository  user.IUserRepository
}

type OrderService struct {
	OrderRepository IOrderRepository
	OrderValidator  IOrderValidator
	UserRepository  user.IUserRepository
}

func NewOrderService(deps OrderServiceDeps) *OrderService {
	return &OrderService{
		OrderRepository: deps.OrderRepository,
		OrderValidator:  deps.OrderValidator,
		UserRepository:  deps.UserRepository,
	}
}

func (s *OrderService) CreateWithItems(ctx context.Context, req *OrderCreateRequest) (*Order, error) {
	err := s.OrderValidator.ValidateProductsExist(req)
	if err != nil {
		return nil, err
	}

	phoneNumber, ok := ctx.Value(constants.ContextPhoneNumber).(string)
	if !ok || phoneNumber == "" {
		return nil, errors.New("phone number not found in context")
	}

	user, err := s.UserRepository.GetByPhoneNumber(phoneNumber)
	if err != nil {
		return nil, err
	}

	order, err := s.OrderRepository.CreateWithItems(ctx, req, user)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) GetByID(orderID uint64) (*Order, error) {
	order, err := s.OrderRepository.GetOrder(orderID)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) GetDetail(order *Order) (*OrderDetailResponse, error) {
	response := GetOrderDetailResponseData(order)
	return response, nil
}

func (s *OrderService) GetListByUserID(userID uint, offset, limit int) ([]OrderDetailResponse, error) {
	orders, err := s.OrderRepository.GetListByUserID(userID, offset, limit)
	if err != nil {
		return nil, err
	}

	items := make([]OrderDetailResponse, len(orders))
	for i, order := range orders {
		items[i] = *GetOrderDetailResponseData(&order)
	}

	return items, nil
}

func GetOrderDetailResponseData(order *Order) *OrderDetailResponse {
	response := &OrderDetailResponse{
		ID:        order.ID,
		UserID:    order.UserID,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
		Status:    order.GetOrderStatusData(),
		Notes:     order.Notes,
		Items:     make([]OrderItemDetail, len(order.OrderItems)),
	}

	for i, item := range order.OrderItems {
		response.Items[i] = OrderItemDetail{
			ID:       item.ID,
			Quantity: item.Quantity,
			Price:    item.Price,
			Discount: item.Discount,
			Product: OrderProductDetail{
				ID:   item.ProductID,
				Name: item.Product.Name,
			},
		}
	}
	return response
}
