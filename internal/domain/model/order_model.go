package model

import (
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/enum"
)

type OrderResponse struct {
	ID            string           `json:"id,omitempty"`
	UserID        string           `json:"user_id,omitempty"`
	Status        enum.OrderStatus `json:"status,omitempty"`
	TotalAmount   float64            `json:"total_amount,omitempty"`
	Currency      string           `json:"currency,omitempty"`
	CreatedAt     int64            `json:"created_at,omitempty"`
	UpdatedAt     int64            `json:"updated_at,omitempty"`
	ExpiredAt     int64            `json:"expired_at,omitempty"`
}

type CreateOrderRequest struct {
	UserID   	   string               `json:"user_id" validate:"required,max=100"`
	IdempotencyKey string          		`json:"idempotency_key" validate:"required,max=100"`
	Currency       string               `json:"currency" validate:"required,len=3"`
	TotalAmount    float64            	`json:"total_amount" validate:"required,gt=0"`
	Items    	   []CreateOrderItem    `json:"items" validate:"required,min=1,dive"`
}

type CreateOrderItem struct {
	SkuID     string `json:"sku_id" validate:"required,max=100"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
	Price     float64  `json:"price" validate:"required,min=1"`
}

type CancelOrderRequest struct {
	OrderID string `json:"-" validate:"required,max=100"`
	UserID  string `json:"user_id" validate:"required,max=100"`
}

type GetOrderRequest struct {
	OrderID string `json:"-" validate:"required,max=100"`
}

type ListOrderRequest struct {
	UserID string `json:"user_id" validate:"required,max=100"`
	Limit  int    `json:"limit,omitempty" validate:"max=100"`
	Offset int    `json:"offset,omitempty"`
}

type ReserveOrderRequest struct {
	OrderID string `json:"-" validate:"required,max=100"`
}

type MarkOrderPaidRequest struct {
	OrderID       string `json:"-" validate:"required,max=100"`
	PaymentRefID  string `json:"payment_ref_id" validate:"required,max=100"`
}