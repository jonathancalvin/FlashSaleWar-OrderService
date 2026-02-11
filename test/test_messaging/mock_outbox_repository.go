package test_messaging

import (
	"github.com/google/uuid"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/entity"
	"gorm.io/gorm"
)

type MockOutboxRepository struct {
	FindPendingFn func(tx *gorm.DB, limit int) ([]entity.OutboxEvent, error)
	MarkSentFn    func(tx *gorm.DB, id uuid.UUID) error
	MarkFailedFn  func(tx *gorm.DB, id uuid.UUID, reason string) error
	MarkProcessingFn  func(tx *gorm.DB, id uuid.UUID) error

	MarkSentCalled   int
	MarkFailedCalled int
	MarkProcessingCalled int
}

// Delete implements repository.OutboxRepository.
func (m *MockOutboxRepository) Delete(tx *gorm.DB, entity *entity.OutboxEvent) error {
	panic("unimplemented")
}

// Update implements repository.OutboxRepository.
func (m *MockOutboxRepository) Update(tx *gorm.DB, entity *entity.OutboxEvent) error {
	panic("unimplemented")
}

func (m *MockOutboxRepository) Create(tx *gorm.DB, e *entity.OutboxEvent) error {
	return nil
}

func (m *MockOutboxRepository) FindPending(tx *gorm.DB, limit int) ([]entity.OutboxEvent, error) {
	return m.FindPendingFn(tx, limit)
}

func (m *MockOutboxRepository) MarkSent(tx *gorm.DB, id uuid.UUID) error {
	m.MarkSentCalled++
	return m.MarkSentFn(tx, id)
}

func (m *MockOutboxRepository) MarkFailed(tx *gorm.DB, id uuid.UUID, reason string) error {
	m.MarkFailedCalled++
	return m.MarkFailedFn(tx, id, reason)
}

func (m *MockOutboxRepository) MarkProcessing(tx *gorm.DB, id uuid.UUID) error {
	m.MarkProcessingCalled++
	return m.MarkProcessingFn(tx, id)
}