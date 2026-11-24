package dto

import "avito-kitchen-service/internal/infrastructure/db/sqlc"

type RestaurantDto struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type RestaurantsDto struct {
	Restaurants []RestaurantDto `json:"restaurants"`
}

func ToDTORestaurant(r sqlc.Restaurant) RestaurantDto {
	return RestaurantDto{
		ID:   r.ID,
		Name: r.Name,
	}
}

func ToDTORestaurants(restaurants []sqlc.Restaurant) []RestaurantDto {
	result := make([]RestaurantDto, len(restaurants))

	for i, r := range restaurants {
		result[i] = ToDTORestaurant(r)
	}

	return result
}
