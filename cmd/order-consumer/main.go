package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/application"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/config"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/delivery/messaging"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/enum"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/infrastructure/repository"
	"gorm.io/gorm"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	viperConfig := config.NewViper()
	logger := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, logger)
	logger.Info("Starting worker service")

	ctx, cancel := context.WithCancel(context.Background())

	go RunOrderConsumer(logger, viperConfig, ctx, db)

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGTERM)

	stop := false
	for !stop {
		select {
		case s := <-terminateSignals:
			logger.Info("Got one of stop signals, shutting down worker gracefully, SIGNAL NAME :", s)
			cancel()
			stop = true
		}
	}

	time.Sleep(5 * time.Second) // wait for consumer to finish processing
	logger.Info("Worker service stopped")
}

func RunOrderConsumer(logger *logrus.Logger, viperConfig *viper.Viper, ctx context.Context, db *gorm.DB) {
	logger.Info("setup order consumer")

	// ===== Repository =====
	orderRepo := repository.NewOrderRepository(logger)
	outboxRepo := repository.NewOutboxRepository(logger)

	// ===== Usecase =====
	orderUsecase := application.NewOrderUseCase(
		db,
		logger,
		orderRepo,
		outboxRepo,
	)
	
	orderConsumerGroup := config.NewKafkaConsumerGroup(viperConfig, logger)
	orderHandler := messaging.NewOrderConsumer(logger, orderUsecase)

	topics := []string{
		string(enum.InventoryTopic),
		string(enum.PaymentTopic),
	}

	messaging.ConsumeTopics(ctx, orderConsumerGroup, topics, logger, orderHandler.Consume)
}