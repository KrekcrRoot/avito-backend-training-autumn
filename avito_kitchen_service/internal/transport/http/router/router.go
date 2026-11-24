package router

import (
	"avito-kitchen-service/internal/infrastructure/db/repository"
	"avito-kitchen-service/internal/transport/http/handler"
	"avito-kitchen-service/internal/transport/http/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(
	handlers *handler.Handlers,
	restaurantRepository *repository.RestaurantRepository,
) http.Handler {
	r := chi.NewRouter()

	r.Get("/_info", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Group(func(r chi.Router) {

		r.Post("/register", handlers.AuthHandler.Register)
		r.Post("/login", handlers.AuthHandler.Login)

	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(handlers.AuthHandler.JWTService))

		r.Get("/restaurants", handlers.RestaurantHandler.GetAllRestaurants)

		r.Get(
			"/restaurants/{id}/products",
			handlers.RestaurantHandler.GetProducts,
		)

		r.Post("/orders", handlers.OrderHandler.CreateOrder)
		r.Get("/orders", handlers.OrderHandler.GetOrders)

		r.Get("/orders/{id}", handlers.OrderHandler.GetOrder)

	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.RestaurantAuthMiddleware(restaurantRepository))

		r.Post(
			"/integration/orders/status",
			handlers.OrderHandler.UpdateOrderStatus,
		)
	})

	return r
}
