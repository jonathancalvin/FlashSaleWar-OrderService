package model

import (
	"time"

	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/enum"
)
type InventoryItem struct {
	ItemID   string `json:"item_id"`
	Quantity int    `json:"quantity"`
}

// InventoryReserved represents the 'INVENTORY_RESERVED' event
type InventoryReservedPayload struct {
	OrderID string          `json:"order_id"`
	Items   []InventoryItem `json:"items"`
	ReservationExpiresAt time.Time `json:"reservation_expires_at"`
}

// InventoryReservationFailed represents the 'INVENTORY_RESERVATION_FAILED' event
type InventoryReservationFailedPayload struct {
	OrderID string          `json:"order_id"`
	FailedItems []InventoryItem `json:"failed_items"`
	Reason  string          `json:"reason"`
}

// InventoryReleased represents the 'INVENTORY_RELEASED' event
type InventoryReleasedPayload struct {
	OrderID string          `json:"order_id"`
	Items   []InventoryItem `json:"items"`
	Reason  enum.InventoryReleaseReason `json:"reason"`
}