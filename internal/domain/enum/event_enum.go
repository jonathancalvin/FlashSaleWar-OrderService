package enum

type EventTopic string

const (
	OrderCreated   EventTopic = "order.v1.created"
	OrderPaid      EventTopic = "order.v1.paid"
	OrderCancelled EventTopic = "order.v1.cancelled"
	OrderExpired   EventTopic = "order.v1.expired"
)

type EventType string

const (
	EventTypeOrderCreated   EventType = "ORDER_CREATED"
	EventTypeOrderPaid      EventType = "ORDER_PAID"
	EventTypeOrderCancelled EventType = "ORDER_CANCELLED"
	EventTypeOrderExpired   EventType = "ORDER_EXPIRED"
)

var EventTypeToTopic = map[EventType]EventTopic{
	EventTypeOrderCreated:   OrderCreated,
	EventTypeOrderPaid:      OrderPaid,
	EventTypeOrderCancelled: OrderCancelled,
	EventTypeOrderExpired:   OrderExpired,
}