package middleware

import (
	"avito-kitchen-service/internal/infrastructure/db/repository"
	"avito-kitchen-service/internal/transport/http/middleware/constants"
	"context"
	"net/http"
)

func RestaurantAuthMiddleware(
	restaurantRepository *repository.RestaurantRepository,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			apiKey := r.Header.Get("X-API-Key")
			if apiKey == "" {
				http.Error(w, "missing api key", http.StatusUnauthorized)
				return
			}

			restaurant, err := restaurantRepository.GetByAPIKey(
				r.Context(),
				apiKey,
			)
			if err != nil {
				http.Error(w, "invalid api key", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(
				r.Context(),
				constants.RestaurantIDKey,
				restaurant.ID,
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
