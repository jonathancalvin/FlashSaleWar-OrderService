package model

import (
	"time"

	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/enum"
)

// PaymentIntentCreated represents the 'payment.intent_created' event
type PaymentIntentCreatedPayload struct {
	PaymentIntentID string    `json:"payment_intent_id"`
	OrderID         string    `json:"order_id"`
	Amount          int64     `json:"amount"`
	Currency        string    `json:"currency"`
	ExpiresAt       time.Time `json:"expires_at"`
}

// PaymentSucceeded represents the 'payment.succeeded' event
type PaymentSucceededPayload struct {
	PaymentIntentID string    `json:"payment_intent_id"`
	OrderID         string    `json:"order_id"`
	GatewayEventID  string    `json:"gateway_event_id"`
	PaidAt          time.Time `json:"paid_at"`
}

// PaymentFailed represents the 'payment.failed' event
type PaymentFailedPayload struct {
	PaymentIntentID string `json:"payment_intent_id"`
	OrderID         string `json:"order_id"`
	GatewayEventID  string `json:"gateway_event_id"`
	Reason          enum.PaymentFailureReason `json:"reason"`
}