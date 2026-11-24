package dto

import (
	"avito-kitchen-service/internal/infrastructure/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type Order struct {
	ID             int64  `json:"id"`
	KitchenOrderID int64  `json:"kitchen_order_id"`
	Status         string `json:"status"`
}

type OrderResponse struct {
	ID                int64               `json:"id"`
	RestaurantID      int64               `json:"restaurant_id"`
	UserID            int64               `json:"user_id"`
	RestaurantOrderID pgtype.Int8         `json:"restaurant_order_id"`
	Status            string              `json:"status"`
	TotalPrice        pgtype.Numeric      `json:"total_price"`
	CreatedAt         pgtype.Timestamptz  `json:"created_at"`
	Items             []OrderItemResponse `json:"items"`
}

type OrderItemResponse struct {
	RestaurantProductID int64          `json:"restaurant_product_id"`
	Name                string         `json:"name"`
	Price               pgtype.Numeric `json:"price"`
	Quantity            int32          `json:"quantity"`
}

type CreateOrderRequest struct {
	RestaurantID int64                    `json:"restaurant_id"`
	Items        []CreateOrderItemRequest `json:"items"`
}

type CreateOrderItemRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

type CreateOrderResult struct {
	Order sqlc.Order
}

type UpdateOrderStatusRequest struct {
	KitchenOrderID int64  `json:"kitchen_order_id"`
	Status         string `json:"status"`
}
