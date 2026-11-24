package handler

import (
	"avito-kitchen-service/internal/service"
	"avito-kitchen-service/internal/transport/http/dto"
	"avito-kitchen-service/internal/transport/http/middleware/constants"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var request dto.CreateOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(constants.UserIDKey).(int64)

	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	order, err := h.orderService.CreateOrder(
		r.Context(),
		userID,
		request,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	WriteJSON(w, http.StatusCreated, order)

}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(constants.UserIDKey).(int64)

	if !ok {
		WriteError(w, "UNAUTHORIZED", "unauthorized", http.StatusUnauthorized)
		return
	}

	orders, err := h.orderService.GetOrdersByUserID(
		r.Context(),
		userID,
	)

	if err != nil {
		WriteError(w, "INTERNAL_ERROR", "failed to get orders", http.StatusInternalServerError)
		return
	}

	WriteJSON(w, http.StatusOK, orders)

}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(constants.UserIDKey).(int64)

	if !ok {
		WriteError(w, "UNAUTHORIZED", "unauthorized", http.StatusUnauthorized)
		return
	}

	orderID, err := strconv.ParseInt(
		chi.URLParam(r, "id"),
		10,
		64,
	)
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	order, err := h.orderService.GetOrder(
		r.Context(),
		userID,
		orderID,
	)

	if err != nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	WriteJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	var request dto.UpdateOrderStatusRequest

	restaurantID, ok := r.Context().Value(constants.RestaurantIDKey).(int64)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	order, err := h.orderService.UpdateOrderStatus(
		r.Context(),
		restaurantID,
		request.KitchenOrderID,
		request.Status,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	WriteJSON(w, http.StatusOK, order)

}
