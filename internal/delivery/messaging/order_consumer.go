package messaging

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/application"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/enum"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/model"
	"github.com/sirupsen/logrus"
)

type EventHandler func(ctx context.Context, payload json.RawMessage) error

type OrderConsumer struct {
	Log *logrus.Logger
	Handlers map[enum.EventType]EventHandler
	OrderUC application.OrderUseCase
}

func NewOrderConsumer(log *logrus.Logger, orderUC application.OrderUseCase) *OrderConsumer {
	c := &OrderConsumer{
		Log: log,
		Handlers: make(map[enum.EventType]EventHandler),
		OrderUC: orderUC,
	}

	c.Handlers[enum.EventTypeInventoryReserved] = c.handleInventoryReserved
	c.Handlers[enum.EventTypeInventoryReservationFailed] = c.handleInventoryReservationFailed
	c.Handlers[enum.EventTypePaymentIntentCreated] = c.handlePaymentIntentCreated
	c.Handlers[enum.EventTypePaymentSucceeded] = c.handlePaymentSucceeded
	return c
}

func (c OrderConsumer) Consume(ctx context.Context, message *sarama.ConsumerMessage) error {
	var envelope model.EventEnvelope
	if err := json.Unmarshal(message.Value, &envelope); err != nil {
		return err
	}

	handler, ok := c.Handlers[envelope.EventType]
	if !ok {
		c.Log.Warnf("no handler for event type %s", envelope.EventType)
		return nil
	}

	c.Log.Infof("Received topic orders with event: %v from partition %d", envelope.EventType, message.Partition)

	payloadBytes, err := json.Marshal(envelope.Payload)
	if err != nil {
		c.Log.Warnf("failed to marshal payload for event id %s", envelope.EventID)
		return err
	}
	return handler(ctx, json.RawMessage(payloadBytes))
}

func (c OrderConsumer) handleInventoryReserved(ctx context.Context, payload json.RawMessage) error {
	var event model.InventoryReservedPayload
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}
	c.Log.Infof("Handling InventoryReserved for OrderID test %s : %s", enum.StatusReserved, event.OrderID)

	resp, err := c.OrderUC.UpdateOrderStatus(
		ctx,
		event.OrderID,
		enum.StatusReserved,
	)

	if err != nil {
		c.Log.WithError(err).Errorf("failed to update order status for OrderID: %s", event.OrderID)
		return err
	}

	c.Log.Infof("Updated order status for OrderID: %s with response: %v", event.OrderID, resp)

	return nil
}

func (c OrderConsumer) handleInventoryReservationFailed(ctx context.Context, payload json.RawMessage) error {
	var event model.InventoryReservationFailedPayload
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}
	c.Log.Infof("Handling InventoryReservationFailed for OrderID: %s", event.OrderID)

	model := model.CancelOrderRequest{
		OrderID: event.OrderID,
		Reason: enum.StockExhausted,
	}

	resp, err := c.OrderUC.CancelOrder(ctx, model)
	
	if err != nil {
		c.Log.WithError(err).Errorf("failed to cancel order for OrderID: %s", event.OrderID)
		return err
	}
	c.Log.Infof("Cancelled order for OrderID: %s with response: %v", event.OrderID, resp)
	return nil
}

func (c OrderConsumer) handlePaymentIntentCreated(ctx context.Context, payload json.RawMessage) error {
	var event model.PaymentIntentCreatedPayload
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}
	c.Log.Infof("Handling PaymentIntentCreated for OrderID: %s", event.OrderID)

	resp, err := c.OrderUC.UpdateOrderStatus(
		ctx,
		event.OrderID,
		enum.StatusPaymentPending,
	)

	if err != nil {
		c.Log.WithError(err).Errorf("failed to update order status for OrderID: %s", event.OrderID)
		return err
	}

	c.Log.Infof("Updated order status for OrderID: %s with response: %v", event.OrderID, resp)
	return nil
}

func (c OrderConsumer) handlePaymentSucceeded(ctx context.Context, payload json.RawMessage) error {
	var event model.PaymentSucceededPayload
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}
	c.Log.Infof("Handling PaymentSucceeded for OrderID: %s", event.OrderID)

	resp, err := c.OrderUC.MarkOrderPaid(
		ctx,
		event.OrderID,
	)
	
	if err != nil {
		c.Log.WithError(err).Errorf("failed to mark order as paid for OrderID: %s", event.OrderID)
		return err
	}

	c.Log.Infof("Marked order as paid for OrderID: %s with response: %v", event.OrderID, resp)
	return nil
}