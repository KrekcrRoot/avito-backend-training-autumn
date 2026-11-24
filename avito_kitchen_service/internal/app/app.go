package app

import (
	"avito-kitchen-service/internal/config"
	"avito-kitchen-service/internal/infrastructure/db"
	"avito-kitchen-service/internal/infrastructure/db/repository"
	"avito-kitchen-service/internal/infrastructure/db/sqlc"
	"avito-kitchen-service/internal/integration/restaurant"
	"avito-kitchen-service/internal/service"
	"avito-kitchen-service/internal/service/auth"
	"avito-kitchen-service/internal/transport/http/handler"
	"avito-kitchen-service/internal/transport/http/router"
	"context"
	"fmt"
	"net/http"
	"time"
)

type App struct {
	Config               *config.Config
	UserService          *service.UserService
	RestaurantService    *service.RestaurantService
	OrderService         *service.OrderService
	JWTService           *auth.JWTService
	RestaurantRepository *repository.RestaurantRepository
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

	userRepository := repository.NewUserRepository(queries)
	userService := service.NewUserService(userRepository)

	httpClient := &http.Client{Timeout: 5 * time.Second}
	client := restaurant.NewClient(httpClient)

	restaurantRepository := repository.NewRestaurantRepository(queries)
	restaurantService := service.NewRestaurantService(restaurantRepository, client)

	orderRepository := repository.NewOrderRepository(pool, queries)
	orderService := service.NewOrderService(orderRepository, restaurantRepository, client)

	jwtService := auth.NewJWTService(cfg.JWTSecret)

	return &App{
		Config:               cfg,
		UserService:          userService,
		RestaurantService:    restaurantService,
		OrderService:         orderService,
		JWTService:           jwtService,
		RestaurantRepository: restaurantRepository,
	}, nil
}

func (app *App) Start() error {

	handlers := &handler.Handlers{
		AuthHandler:       handler.NewAuthHandler(app.JWTService, app.UserService),
		RestaurantHandler: handler.NewRestaurantHandler(app.RestaurantService),
		OrderHandler:      handler.NewOrderHandler(app.OrderService),
	}

	r := router.NewRouter(handlers, app.RestaurantRepository)

	addr := fmt.Sprintf(":%s", app.Config.AppPort)

	if err := http.ListenAndServe(addr, r); err != nil {
		return err
	}

	return nil

}
