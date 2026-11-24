package service

import (
	"context"
	"errors"
	"restaurant-service/internal/config"
	"restaurant-service/internal/infrastructure/db/repository"
	"restaurant-service/internal/infrastructure/db/sqlc"
	"restaurant-service/internal/integration/kitchen"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrProductNotFound    = errors.New("product not found")
	ErrProductUnavailable = errors.New("product unavailable")
)

type CreateOrderItem struct {
	ProductID int64 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

type CreateOrderRequest struct {
	KitchenOrderID int64             `json:"kitchen_order_id"`
	Items          []CreateOrderItem `json:"items"`
}

type OrderItemResponse struct {
	ProductID int64          `json:"product_id"`
	Quantity  int32          `json:"quantity"`
	Price     pgtype.Numeric `json:"price"`
}

type OrderResponse struct {
	ID             int64               `json:"id"`
	KitchenOrderID int64               `json:"kitchen_order_id"`
	Status         string              `json:"status"`
	CreatedAt      pgtype.Timestamptz  `json:"created_at"`
	UpdatedAt      pgtype.Timestamptz  `json:"updated_at"`
	Items          []OrderItemResponse `json:"items"`
}

type OrderService struct {
	orderRepository   *repository.OrderRepository
	productRepository *repository.ProductRepository
	client            *kitchen.Client
	config            *config.Config
}

func NewOrderService(
	orderRepository *repository.OrderRepository,
	productRepository *repository.ProductRepository,
	client *kitchen.Client,
	config *config.Config,
) *OrderService {
	return &OrderService{
		orderRepository:   orderRepository,
		productRepository: productRepository,
		client:            client,
		config:            config,
	}
}

func (s *OrderService) GetOrderByID(ctx context.Context, id int64) (OrderResponse, error) {
	order, err := s.orderRepository.GetByID(ctx, id)
	if err != nil {
		return OrderResponse{}, err
	}

	items, err := s.orderRepository.GetOrderItems(ctx, id)
	if err != nil {
		return OrderResponse{}, err
	}

	responseItems := make([]OrderItemResponse, 0, len(items))

	for _, item := range items {
		responseItems = append(responseItems, OrderItemResponse{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	return OrderResponse{
		ID:             order.ID,
		KitchenOrderID: order.KitchenOrderID,
		Status:         order.Status,
		CreatedAt:      order.CreatedAt,
		UpdatedAt:      order.UpdatedAt,
		Items:          responseItems,
	}, nil
}

func (s *OrderService) CreateOrderWithItems(ctx context.Context, request CreateOrderRequest) (sqlc.Order, error) {

	items := make([]repository.CreateOrderItem, 0, len(request.Items))

	for _, item := range request.Items {

		product, err := s.productRepository.GetByID(ctx, item.ProductID)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return sqlc.Order{}, ErrProductNotFound
			}

			return sqlc.Order{}, err
		}

		if !product.IsAvailable {
			return sqlc.Order{}, ErrProductUnavailable
		}

		items = append(items, repository.CreateOrderItem{
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})

	}

	order, err := s.orderRepository.CreateWithItems(
		ctx,
		request.KitchenOrderID,
		items,
	)

	if err != nil {
		return sqlc.Order{}, err
	}

	s.ProcessOrder(order)

	return order, nil
}

func (s *OrderService) changeStatus(
	ctx context.Context,
	order sqlc.Order,
	status string,
) error {

	_, err := s.orderRepository.UpdateStatus(
		ctx,
		order.ID,
		status,
	)
	if err != nil {
		return err
	}

	return s.client.UpdateOrderStatus(
		ctx,
		s.config.KitchenURL,
		s.config.ApiKey,
		order.KitchenOrderID,
		status,
	)
}

func (s *OrderService) ProcessOrder(order sqlc.Order) {
	go func() {
		ctx := context.Background()

		time.Sleep(3 * time.Second)

		if err := s.changeStatus(ctx, order, "accepted"); err != nil {
			return
		}

		time.Sleep(5 * time.Second)

		if err := s.changeStatus(ctx, order, "cooking"); err != nil {
			return
		}

		time.Sleep(5 * time.Second)

		if err := s.changeStatus(ctx, order, "ready"); err != nil {
			return
		}
	}()
}
