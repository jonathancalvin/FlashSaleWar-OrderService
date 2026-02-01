package entity

import (
	"time"

	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/enum"
)

type Order struct {
	OrderID          string        		`gorm:"column:order_id;type:varchar(36);primaryKey"`
	UserID           string        		`gorm:"column:user_id;type:varchar(36);index;not null"`
	IdempotencyKey   string        		`gorm:"column:idempotency_key;type:varchar(64);uniqueIndex;not null"`
	Status           enum.OrderStatus   `gorm:"column:status;type:varchar(20);not null"`
	Currency  		 string     		`gorm:"column:currency;type:varchar(10)"`
	TotalAmount 	 float64   			`gorm:"column:total_amount;type:decimal(10,2)"`
	
	ExpiresAt 		 time.Time     		`gorm:"column:expires_at"`
	CreatedAt 		 time.Time
	UpdatedAt 		 time.Time

	OrderItems 		 []OrderItem 		`gorm:"foreignKey:OrderID;references:OrderID"`
}

func (o *Order) TableName() string {
	return "orders"
}