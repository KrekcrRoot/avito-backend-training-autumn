package repository

import (
	"avito-kitchen-service/internal/infrastructure/db/sqlc"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func NewOrderRepository(pool *pgxpool.Pool, queries *sqlc.Queries) *OrderRepository {
	return &OrderRepository{
		pool:    pool,
		queries: queries,
	}
}

func (r *OrderRepository) GetOrderByID(ctx context.Context, id int64) (sqlc.Order, error) {
	return r.queries.GetOrderByID(ctx, id)
}

func (r *OrderRepository) GetOrdersByUserId(ctx context.Context, id int64) ([]sqlc.Order, error) {
	return r.queries.GetOrdersByUserID(ctx, id)
}

func (r *OrderRepository) CreateOrder(ctx context.Context, args sqlc.CreateOrderParams) (sqlc.Order, error) {
	return r.queries.CreateOrder(ctx, args)
}

func (r *OrderRepository) AddOrderSnapshotItem(ctx context.Context, args sqlc.AddOrderSnapshotItemParams) (sqlc.OrderItem, error) {
	return r.queries.AddOrderSnapshotItem(ctx, args)
}

type CreateOrderItem struct {
	RestaurantProductID int64
	Name                string
	Price               float64
	Quantity            int32
}

func (r *OrderRepository) CreateWithItems(
	ctx context.Context,
	orderParams sqlc.CreateOrderParams,
	items []CreateOrderItem,
) (sqlc.Order, error) {

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return sqlc.Order{}, err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Printf("failed to rollback transaction: %v", err)
		}
	}()

	q := r.queries.WithTx(tx)

	order, err := q.CreateOrder(ctx, orderParams)
	if err != nil {
		return sqlc.Order{}, err
	}

	for _, item := range items {

		price := pgtype.Numeric{}

		err = price.Scan(fmt.Sprintf("%.2f", item.Price))
		if err != nil {
			return sqlc.Order{}, err
		}

		_, err := q.AddOrderSnapshotItem(ctx, sqlc.AddOrderSnapshotItemParams{
			OrderID:             order.ID,
			RestaurantProductID: item.RestaurantProductID,
			Name:                item.Name,
			Price:               price,
			Quantity:            item.Quantity,
		})
		if err != nil {
			return sqlc.Order{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return sqlc.Order{}, err
	}

	return order, nil
}

func (r *OrderRepository) SetRestaurantOrderId(
	ctx context.Context,
	id int64,
	restaurantOrderId int64,
) (sqlc.Order, error) {
	return r.queries.SetRestaurantOrderID(ctx, sqlc.SetRestaurantOrderIDParams{
		ID: id,
		RestaurantOrderID: pgtype.Int8{
			Int64: restaurantOrderId,
			Valid: true,
		},
	})
}

func (r *OrderRepository) GetOrderItemsByOrderID(ctx context.Context, orderID int64) ([]sqlc.OrderItem, error) {
	return r.queries.GetOrderItemsByOrderID(ctx, orderID)
}

func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, restaurantID int64, orderID int64, status string) (sqlc.Order, error) {
	return r.queries.UpdateOrderStatus(ctx, sqlc.UpdateOrderStatusParams{
		RestaurantID: restaurantID,
		RestaurantOrderID: pgtype.Int8{
			Int64: orderID,
			Valid: true,
		},
		Status: status,
	})
}
