package repository

import (
	"context"
	"log"
	"restaurant-service/internal/infrastructure/db/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

type CreateOrderItem struct {
	ProductID int64
	Quantity  int32
	Price     pgtype.Numeric
}

func NewOrderRepository(pool *pgxpool.Pool, queries *sqlc.Queries) *OrderRepository {
	return &OrderRepository{
		pool:    pool,
		queries: queries,
	}
}

func (r *OrderRepository) GetAll(ctx context.Context) ([]sqlc.Order, error) {
	return r.queries.GetOrders(ctx)
}

func (r *OrderRepository) GetByID(ctx context.Context, id int64) (sqlc.Order, error) {
	return r.queries.GetOrder(ctx, id)
}

func (r *OrderRepository) GetOrderItems(ctx context.Context, id int64) ([]sqlc.OrderItem, error) {
	return r.queries.GetOrderItems(ctx, id)
}

func (r *OrderRepository) Create(ctx context.Context, kitchenOrderID int64) (sqlc.Order, error) {
	return r.queries.CreateOrder(ctx, kitchenOrderID)
}

func (r *OrderRepository) AddOrderItem(ctx context.Context, params sqlc.AddOrderItemParams) (sqlc.OrderItem, error) {
	return r.queries.AddOrderItem(ctx, params)
}

func (r *OrderRepository) CreateWithItems(
	ctx context.Context,
	kitchenOrderID int64,
	orderItems []CreateOrderItem,
) (sqlc.Order, error) {

	tx, err := r.pool.Begin(ctx)

	if err != nil {
		return sqlc.Order{}, err
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Printf("failed to rollback transaction: %v", err)
		}
	}(tx, ctx)

	queries := r.queries.WithTx(tx)

	order, err := queries.CreateOrder(ctx, kitchenOrderID)

	if err != nil {
		return sqlc.Order{}, err
	}

	for _, item := range orderItems {
		_, err := queries.AddOrderItem(ctx, sqlc.AddOrderItemParams{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})

		if err != nil {
			return sqlc.Order{}, nil
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return sqlc.Order{}, err
	}

	return order, nil

}

func (r *OrderRepository) UpdateStatus(
	ctx context.Context,
	orderID int64,
	status string,
) (sqlc.Order, error) {
	return r.queries.UpdateOrderStatus(ctx, sqlc.UpdateOrderStatusParams{
		ID:     orderID,
		Status: status,
	})
}
