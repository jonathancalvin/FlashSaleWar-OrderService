package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type OutboxEvent struct {
    ID          uuid.UUID      `gorm:"column:id;type:uuid;primaryKey"` 
    AggregateID uuid.UUID      `gorm:"column:aggregate_id;type:uuid;not null;index:idx_outbox_aggregate"`
    EventType   string         `gorm:"column:event_type;type:varchar(100);not null"`
    Payload     datatypes.JSON `gorm:"column:payload;not null"`

    Status     string  `gorm:"column:status;type:varchar(20);not null;index:idx_outbox_pending,priority:1"`
    RetryCount int     `gorm:"column:retry_count;not null;default:0"`
    LastError  *string `gorm:"column:last_error;type:text"`
    
    CreatedAt   time.Time  `gorm:"column:created_at;not null;autoCreateTime;index:idx_outbox_pending,priority:2"`
    ProcessedAt *time.Time `gorm:"column:processed_at"`
}

func (OutboxEvent) TableName() string {
	return "outbox_events"
}

func NewOutboxEvent(
	aggregateID uuid.UUID,
	eventType string,
	payload datatypes.JSON,
	status string,
) *OutboxEvent {
	now := time.Now().UTC()

	return &OutboxEvent{
		ID:          uuid.New(),
		AggregateID: aggregateID,
		EventType:   eventType,
		Payload:     payload,
		Status:      status,
		RetryCount:  0,
		CreatedAt:   now,
	}
}