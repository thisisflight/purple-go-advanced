package order

import (
	"net/http"
	"purple/links/configs"
	"purple/links/pkg/constants"
	"purple/links/pkg/middleware"
	"purple/links/pkg/req"
	pkg "purple/links/pkg/res"
	"strconv"
)

type OrderHandlerDeps struct {
	*OrderService
	*configs.Config
}

type OrderHandler struct {
	*OrderService
}

func NewOrderHandler(router *http.ServeMux, deps OrderHandlerDeps) {
	handler := &OrderHandler{
		OrderService: deps.OrderService,
	}
	router.Handle("POST /order", middleware.AuthMiddleware(handler.Create(), deps.Config))
	router.Handle("GET /order/{id}", middleware.AuthMiddleware(handler.Read(), deps.Config))
	router.Handle("GET /my-orders", middleware.AuthMiddleware(handler.GetAll(), deps.Config))
}

func (h *OrderHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := req.HandleBody[OrderCreateRequest](w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		order, err := h.OrderService.CreateWithItems(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		pkg.Json(w, OrderCreateResponse{
			ID:        order.ID,
			UserID:    order.UserID,
			CreatedAt: order.CreatedAt,
			UpdatedAt: order.UpdatedAt,
			Status:    order.GetOrderStatusData(),
			Notes:     order.Notes,
		}, http.StatusCreated)
	}
}

func (h *OrderHandler) Read() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		order, err := h.OrderService.GetByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		phoneNumber, ok := r.Context().Value(constants.ContextPhoneNumber).(string)
		if !ok || phoneNumber == "" {
			http.Error(w, "phone number not found in context", http.StatusBadRequest)
			return
		}

		user, err := h.OrderService.UserRepository.GetByPhoneNumber(phoneNumber)
		if err != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		if order.UserID != user.ID {
			http.Error(w, "cant see this order", http.StatusForbidden)
			return
		}

		resp, err := h.OrderService.GetDetail(order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		pkg.Json(w, resp, http.StatusOK)
	}
}

func (h *OrderHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			http.Error(w, "Invalid offset", http.StatusBadRequest)
			return
		}

		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}

		phoneNumber, ok := r.Context().Value(constants.ContextPhoneNumber).(string)
		if !ok || phoneNumber == "" {
			http.Error(w, "phone number not found in context", http.StatusBadRequest)
			return
		}

		user, err := h.OrderService.UserRepository.GetByPhoneNumber(phoneNumber)
		if err != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		items, err := h.OrderService.GetListByUserID(user.ID, offset, limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		count := h.OrderService.Count(user.ID)

		pkg.Json(w, OrderListResponse{
			Total: count,
			Items: items,
		}, http.StatusOK)
	}
}
