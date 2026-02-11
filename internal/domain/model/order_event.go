package model

import (
	"time"

	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/enum"
)

type OrderCreatedPayload struct {
	OrderID     string             `json:"order_id"`
	UserID      string             `json:"user_id"`
	TotalAmount float64            `json:"total_amount"`
	Currency    string             `json:"currency"`
	Items       []OrderItemPayload `json:"items"`
	CreatedAt   time.Time          `json:"created_at"`
}

type OrderItemPayload struct {
	SkuID    string  `json:"sku_id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func (p OrderCreatedPayload) GetResourceID() string {
    return p.OrderID
}

// topic: order.paid
type OrderPaidPayload struct {
	OrderID string    `json:"order_id"`
	PaidAt  time.Time `json:"paid_at"`
}

func (p OrderPaidPayload) GetResourceID() string {
    return p.OrderID
}

// topic: order.cancelled
type OrderCancelledPayload struct {
	OrderID string `json:"order_id"`
	Reason  enum.CancelReason `json:"reason"`
}

func (p OrderCancelledPayload) GetResourceID() string {
	return p.OrderID
}

// topic: order.expired
type OrderExpiredPayload struct {
	OrderID string `json:"order_id"`
	Reason  enum.ExpireReason `json:"reason"`
}

func (p OrderExpiredPayload) GetResourceID() string {
	return p.OrderID
}