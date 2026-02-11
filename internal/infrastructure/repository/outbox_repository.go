package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OutboxRepository interface {
	BaseRepository[entity.OutboxEvent]

	FindPending(tx *gorm.DB, limit int) ([]entity.OutboxEvent, error)
	MarkSent(tx *gorm.DB, id uuid.UUID) error
	MarkFailed(tx *gorm.DB, id uuid.UUID, reason string) error
	MarkProcessing(tx *gorm.DB, id uuid.UUID) error
}

type outboxRepository struct {
	Repository[entity.OutboxEvent]
	Log *logrus.Logger
}

func NewOutboxRepository(log *logrus.Logger) OutboxRepository {
	return &outboxRepository{
		Log: log,
	}
}

func (r *outboxRepository) FindPending(tx *gorm.DB, limit int) ([]entity.OutboxEvent, error) {
	var events []entity.OutboxEvent

	err := tx.
		Raw(`
			SELECT *
			FROM outbox_events
			WHERE status = 'PENDING'
			ORDER BY created_at
			LIMIT ?
			FOR UPDATE SKIP LOCKED
		`, limit).
		Scan(&events).Error

	return events, err
}

func (r *outboxRepository) MarkProcessing(tx *gorm.DB, id uuid.UUID) error {
	now := time.Now().UTC()
	return tx.
		Model(&entity.OutboxEvent{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":       "PROCESSED",
			"processed_at": &now,
		}).Error
}

func (r *outboxRepository) MarkSent(tx *gorm.DB, id uuid.UUID) error {
	now := time.Now().UTC()

	return tx.
		Model(&entity.OutboxEvent{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":       "SENT",
			"processed_at": &now,
		}).Error
}

func (r *outboxRepository) MarkFailed(tx *gorm.DB, id uuid.UUID, reason string) error {
	return tx.
		Model(&entity.OutboxEvent{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":      "FAILED",
			"retry_count": gorm.Expr("retry_count + 1"),
			"last_error":  reason,
		}).Error
}
