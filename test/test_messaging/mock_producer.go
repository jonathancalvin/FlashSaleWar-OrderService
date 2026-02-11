package test_messaging

import (
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/model"
)

type MockProducer struct {
	PublishFn   func(topic string, key string, envelope model.EventEnvelope) error
	PublishCalled int
	CloseFn     func() error
}

func (m *MockProducer) Publish(topic string, key string, envelope model.EventEnvelope) error {
	m.PublishCalled++
	return m.PublishFn(topic, key, envelope)
}

func (m *MockProducer) Close() error {
	return m.CloseFn()
}
