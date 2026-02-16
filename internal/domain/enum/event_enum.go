package enum

type EventTopic string

const (
	OrderTopic EventTopic = "order.v1.events"
	InventoryTopic EventTopic = "inventory.v1.events"
	PaymentTopic   EventTopic = "payment.v1.events"
)

type EventType string

const (
	EventTypeOrderCreated   EventType = "ORDER_CREATED"
	EventTypeOrderPaid      EventType = "ORDER_PAID"
	EventTypeOrderCancelled EventType = "ORDER_CANCELLED"
	EventTypeOrderExpired   EventType = "ORDER_EXPIRED"

	EventTypeInventoryReserved EventType = "INVENTORY_RESERVED"
	EventTypeInventoryReservationFailed EventType = "INVENTORY_RESERVATION_FAILED"
	EventTypeInventoryReleased EventType = "INVENTORY_RELEASED"
	
	EventTypePaymentIntentCreated EventType = "PAYMENT_INTENT_CREATED"
	EventTypePaymentFailed    EventType = "PAYMENT_FAILED"
	EventTypePaymentSucceeded EventType = "PAYMENT_SUCCEEDED"
)

var EventTypeToTopic = map[EventType]EventTopic{
	EventTypeOrderCreated:   OrderTopic,
    EventTypeOrderPaid:      OrderTopic,
    EventTypeOrderCancelled: OrderTopic,
    EventTypeOrderExpired:   OrderTopic,

	EventTypeInventoryReserved:          InventoryTopic,
	EventTypeInventoryReservationFailed: InventoryTopic,
	EventTypeInventoryReleased:          InventoryTopic,

	EventTypePaymentIntentCreated: PaymentTopic,
	EventTypePaymentFailed:        PaymentTopic,
}