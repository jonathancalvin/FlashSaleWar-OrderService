package test_messaging

import (
	"context"
	"encoding/json"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	"github.com/google/uuid"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/entity"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/enum"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/model"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/infrastructure/messaging"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/test"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestOutboxWorker_Process_Success(t *testing.T) {
	log := logrus.New()
	db := test.SetupTestDB()

	envelope := model.EventEnvelope{
		EventID:       "evt-1",
		EventType:     string(enum.EventTypeOrderCreated),
		OccurredAt:    time.Now(),
		SchemaVersion: 1,
		Payload: model.OrderCreatedPayload{
			OrderID:   "00000000-0000-0000-0000-000000000002",
			UserID:    "user-123",
			TotalAmount:    20000,
			Currency:  "IDR",
			Items: []model.OrderItemPayload{
				{
					SkuID:    "sku-1",
					Quantity: 2,
					Price:    10000,
				},
			},
		},
	}

	payload, _ := json.Marshal(envelope)

	repo := &MockOutboxRepository{
		FindPendingFn: func(tx *gorm.DB, limit int) ([]entity.OutboxEvent, error) {
			return []entity.OutboxEvent{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AggregateID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					EventType:   string(enum.EventTypeOrderCreated),
					Payload:     payload,
					Status:      "PENDING",
				},
			}, nil
		},
		MarkSentFn: func(tx *gorm.DB, id uuid.UUID) error {
			return nil
		},
		MarkProcessingFn: func(tx *gorm.DB, id uuid.UUID) error { 
			return nil 
		},
		MarkFailedFn: func(tx *gorm.DB, id uuid.UUID, reason string) error {
			t.Fatal("should not be called")
			return nil
		},
	}

	producer := &MockProducer{
		PublishFn: func(topic string, key string, envelope model.EventEnvelope) error {
			assert.Equal(t, "order.v1.created", topic)
			assert.Equal(t, "00000000-0000-0000-0000-000000000002", key)
			assert.Equal(t, string(enum.EventTypeOrderCreated), envelope.EventType)
			return nil
		},
	}

	worker := messaging.NewOutboxWorker(db, repo, producer, log)
	
	ctx, cancel := context.WithCancel(context.Background())

	go worker.Run(ctx)

	assert.Eventually(t, func() bool {
        return atomic.LoadInt64((*int64)(unsafe.Pointer(&repo.MarkSentCalled))) == 1
    }, 2*time.Second, 100*time.Millisecond)

	cancel()

	assert.Equal(t, 1, producer.PublishCalled)
	assert.Equal(t, 1, repo.MarkSentCalled)
	assert.Equal(t, 1, repo.MarkProcessingCalled)
	assert.Equal(t, 0, repo.MarkFailedCalled)
}
