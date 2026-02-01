package entity

import (
	"time"

)

type OrderItem struct {
	OrderItemID string            `gorm:"column:order_item_id;type:varchar(36);primaryKey"`
	OrderID     string            `gorm:"column:order_id;type:varchar(36);index;not null"`
	SkuID       string            `gorm:"column:sku_id;type:varchar(64);not null"`
	Quantity    int               `gorm:"column:quantity;not null;check:quantity > 0"`
	Price       float64  		  `gorm:"column:price;type:decimal(10,2);not null"`
	Currency    string            `gorm:"column:currency;type:char(3);not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time        
}

func (OrderItem) TableName() string {
	return "order_items"
}
