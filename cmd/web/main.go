package main

import (
	"fmt"

	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/config"
)

func main() {
	// config
	viperConfig := config.NewViper()

	// logger
	log := config.NewLogger(viperConfig)

	// infra
	db := config.NewDatabase(viperConfig, log)
	validate := config.NewValidator()
	router := config.NewGin(viperConfig)

	// bootstrap
	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		Router:   router,
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
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal(err)
	}
}
