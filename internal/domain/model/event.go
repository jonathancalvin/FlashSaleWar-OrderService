package model

import "time"

type IdentifiablePayload interface {
    GetResourceID() string
}

type EventEnvelope struct {
	EventID       string    `json:"event_id"`
	EventType     string    `json:"event_type"`
	OccurredAt    time.Time `json:"occurred_at"`
	SchemaVersion int       `json:"schema_version"`
	Payload       any       `json:"payload"`
}