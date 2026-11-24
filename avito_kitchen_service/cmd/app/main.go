package main

import (
	"avito-kitchen-service/internal/app"
	"avito-kitchen-service/internal/config"
	"log"
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
