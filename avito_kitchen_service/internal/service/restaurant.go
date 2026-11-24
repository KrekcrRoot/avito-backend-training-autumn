package service

import (
	"avito-kitchen-service/internal/infrastructure/db/repository"
	"avito-kitchen-service/internal/infrastructure/db/sqlc"
	"avito-kitchen-service/internal/integration/restaurant"
	"avito-kitchen-service/internal/transport/http/dto"
	"context"
)

type RestaurantService struct {
	restaurantRepository *repository.RestaurantRepository
	client               *restaurant.Client
}

func NewRestaurantService(
	restaurantRepository *repository.RestaurantRepository,
	client *restaurant.Client,
) *RestaurantService {
	return &RestaurantService{
		restaurantRepository: restaurantRepository,
		client:               client,
	}
}

func (s *RestaurantService) GetAll(ctx context.Context) ([]sqlc.Restaurant, error) {
	return s.restaurantRepository.GetAll(ctx)
}

func (s *RestaurantService) GetProducts(ctx context.Context, restaurantID int64) ([]dto.Product, error) {

	dbRestaurant, err := s.restaurantRepository.GetByID(ctx, restaurantID)

	if err != nil {
		return nil, err
	}

	return s.client.GetProducts(ctx, dbRestaurant)

}
