package service

import (
	"avito-kitchen-service/internal/infrastructure/db/repository"
	"avito-kitchen-service/internal/infrastructure/db/sqlc"
	"avito-kitchen-service/internal/integration/restaurant"
	"avito-kitchen-service/internal/transport/http/dto"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type OrderService struct {
	orderRepository      *repository.OrderRepository
	restaurantRepository *repository.RestaurantRepository
	restaurantClient     *restaurant.Client
}

func NewOrderService(
	orderRepository *repository.OrderRepository,
	restaurantRepository *repository.RestaurantRepository,
	restaurantClient *restaurant.Client,
) *OrderService {
	return &OrderService{
		orderRepository:      orderRepository,
		restaurantRepository: restaurantRepository,
		restaurantClient:     restaurantClient,
	}
}

func (s *OrderService) CreateOrder(
	ctx context.Context,
	userID int64,
	request dto.CreateOrderRequest,
) (sqlc.Order, error) {

	if len(request.Items) == 0 {
		return sqlc.Order{}, fmt.Errorf("order must contain at least one item")
	}

	for _, item := range request.Items {
		if item.Quantity <= 0 {
			return sqlc.Order{}, fmt.Errorf(
				"product %d: quantity must be greater than zero",
				item.ProductID,
			)
		}
	}

	dbRestaurant, err := s.restaurantRepository.GetByID(
		ctx,
		request.RestaurantID,
	)
	if err != nil {
		return sqlc.Order{}, err
	}

	products, err := s.restaurantClient.GetProducts(ctx, dbRestaurant)
	if err != nil {
		return sqlc.Order{}, err
	}

	productMap := make(map[int64]dto.Product)

	for _, product := range products {
		productMap[product.ID] = product
	}

	var total float64

	snapshotItems := make([]repository.CreateOrderItem, 0, len(request.Items))

	for _, item := range request.Items {

		product, ok := productMap[item.ProductID]
		if !ok {
			return sqlc.Order{}, fmt.Errorf(
				"product %d not found",
				item.ProductID,
			)
		}

		if !product.IsAvailable {
			return sqlc.Order{}, fmt.Errorf(
				"product %d is unavailable",
				item.ProductID,
			)
		}

		total += product.Price * float64(item.Quantity)

		snapshotItems = append(snapshotItems, repository.CreateOrderItem{
			RestaurantProductID: product.ID,
			Name:                product.Name,
			Price:               product.Price,
			Quantity:            item.Quantity,
		})
	}

	totalNumeric := pgtype.Numeric{}

	err = totalNumeric.Scan(fmt.Sprintf("%.2f", total))
	if err != nil {
		return sqlc.Order{}, err
	}

	order, err := s.orderRepository.CreateWithItems(
		ctx,
		sqlc.CreateOrderParams{
			RestaurantID: request.RestaurantID,
			UserID:       userID,
			TotalPrice:   totalNumeric,
		},
		snapshotItems,
	)

	if err != nil {
		return sqlc.Order{}, err
	}

	restaurantItems := make([]restaurant.CreateOrderItemRequest, 0, len(request.Items))

	for _, item := range request.Items {

		restaurantItems = append(
			restaurantItems,
			restaurant.CreateOrderItemRequest{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
			},
		)
	}

	restaurantOrder, err := s.restaurantClient.CreateOrder(
		ctx,
		dbRestaurant,
		order.ID,
		restaurantItems,
	)
	if err != nil {
		return sqlc.Order{}, err
	}

	order, err = s.orderRepository.SetRestaurantOrderId(ctx, order.ID, restaurantOrder.ID)

	if err != nil {
		return sqlc.Order{}, err
	}

	return order, nil
}

func (s *OrderService) GetOrdersByUserID(
	ctx context.Context,
	userID int64,
) ([]sqlc.Order, error) {
	return s.orderRepository.GetOrdersByUserId(ctx, userID)
}

func (s *OrderService) GetOrder(
	ctx context.Context,
	userID int64,
	orderID int64,
) (dto.OrderResponse, error) {

	order, err := s.orderRepository.GetOrderByID(ctx, orderID)
	if err != nil {
		return dto.OrderResponse{}, err
	}

	if order.UserID != userID {
		return dto.OrderResponse{}, fmt.Errorf("order not found")
	}

	items, err := s.orderRepository.GetOrderItemsByOrderID(ctx, orderID)
	if err != nil {
		return dto.OrderResponse{}, err
	}

	response := dto.OrderResponse{
		ID:                order.ID,
		RestaurantID:      order.RestaurantID,
		UserID:            order.UserID,
		RestaurantOrderID: order.RestaurantOrderID,
		Status:            order.Status,
		TotalPrice:        order.TotalPrice,
		CreatedAt:         order.CreatedAt,
		Items:             make([]dto.OrderItemResponse, 0, len(items)),
	}

	for _, item := range items {
		response.Items = append(response.Items, dto.OrderItemResponse{
			RestaurantProductID: item.RestaurantProductID,
			Name:                item.Name,
			Price:               item.Price,
			Quantity:            item.Quantity,
		})
	}

	return response, nil
}

func (s *OrderService) UpdateOrderStatus(
	ctx context.Context,
	restaurantID int64,
	orderID int64,
	status string,
) (sqlc.Order, error) {

	allowedStatuses := map[string]bool{
		"new":      true,
		"accepted": true,
		"cooking":  true,
		"ready":    true,
		"rejected": true,
		"canceled": true,
	}

	if !allowedStatuses[status] {
		return sqlc.Order{}, fmt.Errorf("invalid order status: %s", status)
	}

	return s.orderRepository.UpdateOrderStatus(ctx, restaurantID, orderID, status)
}
