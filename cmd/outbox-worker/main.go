package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/config"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/infrastructure/messaging"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/infrastructure/repository"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)

	log.Info("starting outbox worker")

	// setup dependencies
	outboxRepo := repository.NewOutboxRepository(log)

	producer, err := config.NewKafkaProducer(
		viperConfig,
		log,
	)
	
	if err != nil {
		log.Fatal(err)
	}

	worker := messaging.NewOutboxWorker(db, outboxRepo, producer, log)

	ctx, cancel := context.WithCancel(context.Background())
	go worker.Run(ctx)

	// graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	log.Info("shutdown signal received")

	cancel()
	_ = producer.Close()
}
