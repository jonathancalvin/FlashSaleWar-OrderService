package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/entity"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/enum"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/model"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/infrastructure/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OutboxWorker struct {
	DB       *gorm.DB
	Repo     repository.OutboxRepository
	Producer KafkaProducer
	Log      *logrus.Logger
}


func NewOutboxWorker(
	db *gorm.DB,
	repo repository.OutboxRepository,
	producer KafkaProducer,
	log *logrus.Logger,
) *OutboxWorker {
	return &OutboxWorker{
		DB:       db,
		Repo:     repo,
		Producer: producer,
		Log:      log,
	}
}

func (w *OutboxWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	w.Log.Info("outbox worker started")

	for {
		select {
		case <-ctx.Done():
			w.Log.Info("outbox worker stopped")
			return
		case <-ticker.C:
			w.process(ctx)
		}
	}
}

func (w *OutboxWorker) process(ctx context.Context) {
	err := w.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		events, err := w.Repo.FindPending(tx, 100)
		if err != nil {
			return err
		}

		for _, evt := range events {
			if err := w.handleEvent(tx, evt); err != nil {
				w.Log.WithFields(logrus.Fields{
					"event_id":   evt.ID,
					"event_type": evt.EventType,
				}).WithError(err).Error("failed to process outbox event")
			}
		}

		return nil
	})

	if err != nil {
		w.Log.WithError(err).Error("outbox worker transaction failed")
	}
}


func (w *OutboxWorker) handleEvent(
	tx *gorm.DB,
	evt entity.OutboxEvent,
) error {
	if err := w.Repo.MarkProcessing(tx, evt.ID); err != nil {
		return err
	}

	topic, ok := enum.EventTypeToTopic[enum.EventType(evt.EventType)]
	if !ok {
		err := fmt.Errorf("no topic mapping for event type: %s", evt.EventType)
		_ = w.Repo.MarkFailed(tx, evt.ID, err.Error())
		return err
	}

	// Setting up envelope
	var rawPayload json.RawMessage
    if err := json.Unmarshal(evt.Payload, &rawPayload); err != nil {
        _ = w.Repo.MarkFailed(tx, evt.ID, err.Error())
        return err
    }

    envelope := model.EventEnvelope{
        EventID:       evt.ID.String(),
        EventType:     evt.EventType,
        OccurredAt:    evt.CreatedAt,
        SchemaVersion: 1,
        Payload:       rawPayload,
    }

	if err := w.Producer.Publish(
		string(topic),
		evt.AggregateID.String(),
		envelope,
	); err != nil {
		_ = w.Repo.MarkFailed(tx, evt.ID, err.Error())
		return err
	}

	return w.Repo.MarkSent(tx, evt.ID)
}
