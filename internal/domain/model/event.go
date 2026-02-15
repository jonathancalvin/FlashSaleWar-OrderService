package model

import (
	"time"

	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/enum"
)

type IdentifiablePayload interface {
    GetResourceID() string
}

type EventEnvelope struct {
	EventID       string    `json:"event_id"`
	EventType     enum.EventType `json:"event_type"`
	OccurredAt    time.Time `json:"occurred_at"`
	SchemaVersion int       `json:"schema_version"`
	Payload       any       `json:"payload"`
}