package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"restaurant-service/internal/service"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {

	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		http.Error(w, "Invalid order id", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrderByID(r.Context(), id)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
		
	}

	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {

	var request service.CreateOrderRequest

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	order, err := h.service.CreateOrderWithItems(r.Context(), request)

	if err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(order); err != nil {
		log.Fatal(err)
	}

}
