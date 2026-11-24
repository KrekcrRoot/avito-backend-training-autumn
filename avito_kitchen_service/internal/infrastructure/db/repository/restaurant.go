package repository

import (
	"avito-kitchen-service/internal/infrastructure/db/sqlc"
	"context"
)

type RestaurantRepository struct {
	queries *sqlc.Queries
}

func NewRestaurantRepository(queries *sqlc.Queries) *RestaurantRepository {
	return &RestaurantRepository{
		queries: queries,
	}
}

func (r *RestaurantRepository) GetAll(ctx context.Context) ([]sqlc.Restaurant, error) {
	return r.queries.GetRestaurants(ctx)
}

func (r *RestaurantRepository) GetByID(ctx context.Context, id int64) (sqlc.Restaurant, error) {
	return r.queries.GetRestaurantByID(ctx, id)
}

func (r *RestaurantRepository) GetByAPIKey(ctx context.Context, apiKey string) (sqlc.Restaurant, error) {
	return r.queries.GetRestaurantByAPIKey(ctx, apiKey)
}
