package main

import (
	"fmt"
	"fp-designpattern/internal/config"
	"fp-designpattern/pkg/timezone"
)

func main() {
	app := config.NewFiber(config.NewViper())
	app.Static("/images", "./public/images")
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validate := config.NewValidator(viperConfig)
	timezone.InitTimeLocation()
	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
	})

	webPort := viperConfig.GetInt("web.port")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
