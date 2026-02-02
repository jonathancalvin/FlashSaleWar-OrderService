package test

import (
	"log"

	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(
		&entity.Order{},
		&entity.OrderItem{},
	); err != nil {
		log.Fatal(err)
	}

	return db
}