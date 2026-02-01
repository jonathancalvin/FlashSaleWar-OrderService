package enum

type OrderStatus string

const (
	StatusCreated        OrderStatus = "CREATED"
	StatusExpired        OrderStatus = "EXPIRED"
	StatusPaymentPending OrderStatus = "PAYMENT_PENDING"
	StatusPaid           OrderStatus = "PAID"
	StatusCancelled      OrderStatus = "CANCELLED"
)
