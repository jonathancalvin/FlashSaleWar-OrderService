package test_usecase

import (
	"github.com/go-playground/validator/v10"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/application"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/infrastructure/repository"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/test"
	"github.com/sirupsen/logrus"
)

func setupOrderUseCase() (application.OrderUseCase) {
    db := test.SetupTestDB()
    log := logrus.New()
    v := validator.New()
    
    orderRepo := repository.NewOrderRepository(log)
    uc := application.NewOrderUseCase(db, log, v, orderRepo)
    
    return uc
}