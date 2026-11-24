package main

import (
	"log"
	"restaurant-service/internal/app"
	"restaurant-service/internal/config"
)

func main() {

	cfg := config.Load()

	bootstrap, err := app.NewApp(cfg)

	if err != nil {
		log.Fatal(err)
	}

	if err := bootstrap.Start(); err != nil {
		log.Fatal(err)
	}

}
