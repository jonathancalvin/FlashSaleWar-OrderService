package config

import (
	"github.com/IBM/sarama"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/infrastructure/messaging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewKafkaProducer(config *viper.Viper, log *logrus.Logger) (messaging.KafkaProducer, error) {
	brokers := config.GetStringSlice("kafka.brokers")
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll 
	saramaConfig.Producer.Retry.Max = 5
	saramaConfig.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, saramaConfig)
	if err != nil {
		log.Errorf("failed to create Kafka producer: %v", err)
		return nil, err
	}

	return &messaging.KafkaProducerImpl{
		Producer: producer,
		Log:      log,
	}, nil
}