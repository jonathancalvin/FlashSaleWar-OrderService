package messaging

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/model"
)

type KafkaProducer interface {
	Publish(topic string, key string, payload model.EventEnvelope) error
	Close() error
}

type KafkaProducerImpl struct {
	Producer sarama.SyncProducer
	Log      *logrus.Logger
	done     chan struct{}
}

func (p *KafkaProducerImpl) Publish(topic string, key string, envelope model.EventEnvelope) error {
	payload, err := json.Marshal(envelope)
	if err != nil {
		p.Log.WithFields(logrus.Fields{
			"topic":    topic,
			"event_id": envelope.EventID,
		}).Errorf("failed to marshal event envelope: %v", err)
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(payload),
		Timestamp: envelope.OccurredAt,
	}

	partition, offset, err := p.Producer.SendMessage(msg)
	if err != nil {
		p.Log.WithFields(logrus.Fields{
			"topic":    topic,
			"event_id": envelope.EventID,
		}).Errorf("failed to send message to Kafka: %v", err)
		return err
	}

	p.Log.Infof("Message sent to partition %d at offset %d", partition, offset)
	return nil
}

func (p *KafkaProducerImpl) Close() error {
	return p.Producer.Close()
}