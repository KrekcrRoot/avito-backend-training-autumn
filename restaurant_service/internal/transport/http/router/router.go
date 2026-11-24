package router

import (
	"net/http"
	"restaurant-service/internal/config"
	"restaurant-service/internal/transport/http/handler"
	"restaurant-service/internal/transport/http/middleware"

	"github.com/go-chi/chi/v5"
)

func NewRouter(handlers *handler.Handlers, cfg *config.Config) http.Handler {

	r := chi.NewRouter()

	r.Get("/_info", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.APIKeyMiddle(cfg.ApiKey))
		r.Get("/products", handlers.ProductHandler.GetProducts)

		r.Get("/orders/{id}", handlers.OrderHandler.GetOrder)
		r.Post("/orders", handlers.OrderHandler.CreateOrder)
	})

	return r

}
