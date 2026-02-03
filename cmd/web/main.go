package main

import (
	"fmt"

	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/config"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validate := config.NewValidator()
	app := config.NewGin(viperConfig)

	// bootstrap
	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
	})

	// run server
	port := viperConfig.GetInt("web.port")
	if port == 0 {
		port = 8080
	}

	log.Infof("HTTP server running on port %d", port)
	if err := app.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal(err)
	}
}
