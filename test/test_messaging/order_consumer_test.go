package test_messaging

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/delivery/messaging"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/enum"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/model"
	mocks "github.com/jonathancalvin/FlashSaleWar-OrderService/test/mock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrderConsumer_UnknownEvent_NoError(t *testing.T) {
	mockUC := new(mocks.MockOrderUseCase)
	logger := logrus.New()

	consumer := messaging.NewOrderConsumer(logger, mockUC)

	msg := buildMessage(
		enum.EventType("UNKNOWN_EVENT"),
		map[string]string{"foo": "bar"},
	)

	err := consumer.Consume(context.Background(), msg)

	assert.NoError(t, err)
}

func TestOrderConsumer_InvalidJSON_ReturnsError(t *testing.T) {
	mockUC := new(mocks.MockOrderUseCase)
	logger := logrus.New()

	consumer := messaging.NewOrderConsumer(logger, mockUC)

	msg := &sarama.ConsumerMessage{
		Value: []byte("invalid-json"),
	}

	err := consumer.Consume(context.Background(), msg)

	assert.Error(t, err)
}

func TestOrderConsumer_InventoryReserved_Success(t *testing.T) {
	mockUC := new(mocks.MockOrderUseCase)
	logger := logrus.New()

	orderID := "order-123"

	mockUC.
		On("UpdateOrderStatus", mock.Anything, orderID, enum.StatusReserved).
		Return(&model.OrderResponse{
			ID:     orderID,
			Status: enum.StatusReserved,
		}, nil).
		Once()

	consumer := messaging.NewOrderConsumer(logger, mockUC)

	msg := buildMessage(
		enum.EventTypeInventoryReserved,
		model.InventoryReservedPayload{
			OrderID: orderID,
		},
	)

	err := consumer.Consume(context.Background(), msg)

	assert.NoError(t, err)
	mockUC.AssertExpectations(t)
}


func TestOrderConsumer_InventoryReservationFailed_Success(t *testing.T) {
	mockUC := new(mocks.MockOrderUseCase)
	logger := logrus.New()
	
	orderID := "order-456"
	
	expectedReq := model.CancelOrderRequest{
		OrderID: orderID,
		Reason:  enum.StockExhausted,
	}
	
	mockUC.
	On("CancelOrder", mock.Anything, expectedReq).
	Return(&model.OrderResponse{
		ID: orderID,
		Status:  enum.StatusCancelled,
		}, nil).
		Once()
		
		consumer := messaging.NewOrderConsumer(logger, mockUC)
		
		msg := buildMessage(
			enum.EventTypeInventoryReservationFailed,
			model.InventoryReservationFailedPayload{
				OrderID: orderID,
			},
	)

	err := consumer.Consume(context.Background(), msg)

	assert.NoError(t, err)
	mockUC.AssertExpectations(t)
}

func TestOrderConsumer_PaymentIntentCreated_Success(t *testing.T) {
	mockUC := new(mocks.MockOrderUseCase)
	logger := logrus.New()

	orderID := "order-789"

	mockUC.
		On("UpdateOrderStatus", mock.Anything, orderID, enum.StatusPaymentPending).
		Return(&model.OrderResponse{
			ID:     orderID,
			Status: enum.StatusPaymentPending,
		}, nil).
		Once()

	consumer := messaging.NewOrderConsumer(logger, mockUC)

	msg := buildMessage(
		enum.EventTypePaymentIntentCreated,
		model.PaymentIntentCreatedPayload{
			OrderID:         orderID,
			PaymentIntentID: "pi_456",
		},
	)

	err := consumer.Consume(context.Background(), msg)

	assert.NoError(t, err)
	mockUC.AssertExpectations(t)
	
}

func TestOrderConsumer_PaymentSucceeded_Success(t *testing.T) {
	mockUC := new(mocks.MockOrderUseCase)
	logger := logrus.New()

	orderID := "order-123"

	mockUC.
		On("MarkOrderPaid", mock.Anything, orderID).
		Return(&model.OrderResponse{
			ID: orderID,
			Status:  enum.StatusPaid,
		}, nil).
		Once()

	consumer := messaging.NewOrderConsumer(logger, mockUC)

	msg := buildMessage(
		enum.EventTypePaymentSucceeded,
		model.PaymentSucceededPayload{
			PaymentIntentID: "pi_123",
			OrderID: orderID,
			GatewayEventID: "evt_123",
			PaidAt: time.Now(),
		},
	)

	err := consumer.Consume(context.Background(), msg)

	assert.NoError(t, err)
	mockUC.AssertExpectations(t)
}


func buildMessage(eventType enum.EventType, payload interface{}) *sarama.ConsumerMessage {
	envelope := model.EventEnvelope{
		EventID:   "evt-123",
		EventType: eventType,
		Payload:   payload,
	}

	bytes, _ := json.Marshal(envelope)

	return &sarama.ConsumerMessage{
		Value:     bytes,
		Partition: 0,
	}
}
