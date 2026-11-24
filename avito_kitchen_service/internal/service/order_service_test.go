package service

import (
	"avito-kitchen-service/internal/transport/http/dto"
	"context"
	"testing"
)

func TestOrderService_CreateOrder_EmptyItems(t *testing.T) {
	service := &OrderService{}

	_, err := service.CreateOrder(
		context.Background(),
		1,
		dto.CreateOrderRequest{
			RestaurantID: 1,
			Items:        nil,
		},
	)

	if err == nil {
		t.Fatal("expected error for empty order")
	}
}

func TestOrderService_CreateOrder_InvalidQuantity(t *testing.T) {
	service := &OrderService{}

	_, err := service.CreateOrder(
		context.Background(),
		1,
		dto.CreateOrderRequest{
			RestaurantID: 1,
			Items: []dto.CreateOrderItemRequest{
				{
					ProductID: 1,
					Quantity:  0,
				},
			},
		},
	)

	if err == nil {
		t.Fatal("expected error for invalid quantity")
	}
}

func TestOrderService_UpdateOrderStatus_InvalidStatus(t *testing.T) {
	service := &OrderService{}

	_, err := service.UpdateOrderStatus(
		context.Background(),
		1,
		1,
		"not-exist-status",
	)

	if err == nil {
		t.Fatal("expected error for invalid status")
	}
}
