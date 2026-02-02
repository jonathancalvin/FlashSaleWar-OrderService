package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/enum"
)

type Order struct {
	OrderID          string        		`gorm:"column:order_id;type:varchar(36);primaryKey"`
	UserID           string        		`gorm:"column:user_id;type:varchar(36);index;not null"`
	IdempotencyKey   string        		`gorm:"column:idempotency_key;type:varchar(64);uniqueIndex;not null"`
	Status           enum.OrderStatus   `gorm:"column:status;type:varchar(20);not null"`
	Currency  		 string     		`gorm:"column:currency;type:varchar(10)"`
	TotalAmount 	 float64   			`gorm:"column:total_amount;type:decimal(10,2)"`
	
	ExpiredAt 		 time.Time     		`gorm:"column:expired_at"`
	CreatedAt 		 time.Time
	UpdatedAt 		 time.Time

	OrderItems 		 []OrderItem 		`gorm:"foreignKey:OrderID;references:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func NewOrder(
	userID string,
	idempotencyKey string,
	status enum.OrderStatus,
	expiredAt time.Time,
	currency string,
	totalAmount float64,
) *Order {

	now := time.Now().UTC()

	return &Order{
		OrderID:        uuid.NewString(),
		UserID:         userID,
		IdempotencyKey: idempotencyKey,
		Status:         status,
		Currency:       currency,
		TotalAmount:    totalAmount,

		ExpiredAt: expiredAt,
		CreatedAt: now,
		UpdatedAt: now,

		OrderItems: make([]OrderItem, 0),
	}
}

func (o *Order) AddItem(
	skuID string,
	qty int,
	price float64,
	currency string,
) {
	item := OrderItem{
		OrderItemID: uuid.NewString(),
		OrderID:     o.OrderID,
		SkuID:       skuID,
		Quantity:    qty,
		Price:       price,
		Currency:    currency,
	}
	o.OrderItems = append(o.OrderItems, item)
}

func (o *Order) TableName() string {
	return "orders"
}