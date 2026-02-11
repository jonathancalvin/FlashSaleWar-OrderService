package test_usecase

import (
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/application"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/infrastructure/repository"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/test"
	"github.com/sirupsen/logrus"
)

func setupOrderUseCase() (application.OrderUseCase) {
    db := test.SetupTestDB()
    log := logrus.New()
    
    orderRepo := repository.NewOrderRepository(log)
    outboxRepo := repository.NewOutboxRepository(log)
    uc := application.NewOrderUseCase(db, log, orderRepo, outboxRepo)
    
    return uc
}