package messaging

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type ConsumerHandler func(ctx context.Context, message *sarama.ConsumerMessage) error

type ConsumerGroupHandler struct {
	Handler ConsumerHandler
	Log     *logrus.Logger
}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			err := h.Handler(session.Context(), message)
			if err != nil {
				h.Log.WithError(err).Error("Failed to process message")
			} else {
				session.MarkMessage(message, "")
			}

		case <-session.Context().Done():
			return nil
		}
	}
}

func ConsumeTopics(ctx context.Context, consumerGroup sarama.ConsumerGroup, topics []string, log *logrus.Logger, handler ConsumerHandler) {
	consumerHandler := &ConsumerGroupHandler{
		Handler: handler,
		Log:     log,
	}

	go func() {
		for {
			if err := consumerGroup.Consume(ctx, topics, consumerHandler); err != nil {
				log.WithError(err).Error("Error from consumer")
			}

			if ctx.Err() != nil {
				log.Info("Context cancelled, stopping consumer")
				return
			}
		}
	}()

	go func() {
		for err := range consumerGroup.Errors() {
			log.WithError(err).Error("Consumer group error")
		}
	}()

	<-ctx.Done()
	log.Infof("Closing consumer group for topics: %v", topics)
	if err := consumerGroup.Close(); err != nil {
		log.WithError(err).Error("Error closing consumer group")
	}
}