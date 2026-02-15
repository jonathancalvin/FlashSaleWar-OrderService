package model

import (
	"time"

	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/enum"
)

// InventoryReserved represents the 'inventory.reserved' event
type InventoryReserved struct {
	OrderID              string    `json:"order_id"`
	ItemID               string    `json:"item_id"`
	Quantity             int       `json:"quantity"`
	ReservationExpiresAt time.Time `json:"reservation_expires_at"`
}

// InventoryReservationFailed represents the 'inventory.reservation_failed' event
type InventoryReservationFailed struct {
	OrderID string `json:"order_id"`
	ItemID  string `json:"item_id"`
}

// InventoryReleased represents the 'inventory.released' event
type InventoryReleased struct {
	OrderID string `json:"order_id"`
	ItemID  string `json:"item_id"`
	Reason  enum.InventoryReleaseReason `json:"reason"`
}