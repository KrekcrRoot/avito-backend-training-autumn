package app

import (
	"context"
	"fmt"
	"net/http"
	"restaurant-service/internal/config"
	"restaurant-service/internal/infrastructure/db"
	"restaurant-service/internal/infrastructure/db/repository"
	"restaurant-service/internal/infrastructure/db/sqlc"
	"restaurant-service/internal/integration/kitchen"
	"restaurant-service/internal/service"
	"restaurant-service/internal/transport/http/handler"
	"restaurant-service/internal/transport/http/router"
	"time"

	"github.com/go-chi/chi/v5"
)

type App struct {
	Config         *config.Config
	ProductService *service.ProductService
	OrderService   *service.OrderService
	Server         *chi.Mux
}

func NewApp(cfg *config.Config) (*App, error) {

	ctx := context.Background()

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	pool, err := db.NewPool(ctx, connString)

	if err != nil {
		return nil, err
	}

	queries := sqlc.New(pool)

	productRepository := repository.NewProductRepository(queries)
	productService := service.NewProductService(productRepository)

	httpClient := &http.Client{Timeout: 5 * time.Second}
	kitchenClient := kitchen.NewClient(httpClient)

	orderRepository := repository.NewOrderRepository(pool, queries)
	orderService := service.NewOrderService(
		orderRepository,
		productRepository,
		kitchenClient,
		cfg,
	)

	app := &App{
		Config:         cfg,
		ProductService: productService,
		OrderService:   orderService,
		Server:         chi.NewRouter(),
	}

	return app, nil

}

func (app *App) Start() error {

	handlers := &handler.Handlers{
		ProductHandler: handler.NewProductHandler(
			app.ProductService,
		),
		OrderHandler: handler.NewOrderHandler(
			app.OrderService,
		),
	}

	r := router.NewRouter(handlers, app.Config)

	addr := fmt.Sprintf(":%s", app.Config.AppPort)

	if err := http.ListenAndServe(addr, r); err != nil {
		return err
	}

	return nil

}
