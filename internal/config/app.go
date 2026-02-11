package config

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/application"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/delivery/http"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/delivery/http/route"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/infrastructure/repository"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *gin.Engine
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(cfg *BootstrapConfig) {
	// ===== Repository =====
	orderRepo := repository.NewOrderRepository(cfg.Log)
	outboxRepo := repository.NewOutboxRepository(cfg.Log)

	// ===== Usecase =====
	orderUsecase := application.NewOrderUseCase(
		cfg.DB,
		cfg.Log,
		orderRepo,
		outboxRepo,
	)

	// ===== Controller =====
	orderController := http.NewOrderController(orderUsecase, cfg.Log, cfg.Validate)

	// ===== Routes =====
	routeConfig := route.RouteConfig{
		App:          	 cfg.App,
		Log:             cfg.Log,
		OrderController: orderController,
	}
	routeConfig.Setup()
}