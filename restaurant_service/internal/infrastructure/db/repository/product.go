package repository

import (
	"context"
	"restaurant-service/internal/infrastructure/db/sqlc"
)

type ProductRepository struct {
	queries *sqlc.Queries
}

func NewProductRepository(queries *sqlc.Queries) *ProductRepository {
	return &ProductRepository{
		queries: queries,
	}
}

func (r *ProductRepository) GetAll(ctx context.Context) ([]sqlc.Product, error) {
	return r.queries.GetProducts(ctx)
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (sqlc.Product, error) {
	return r.queries.GetProduct(ctx, id)
}
