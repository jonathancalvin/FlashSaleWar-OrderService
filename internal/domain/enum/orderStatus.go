package enum

type OrderStatus string

const (
	StatusCreated        OrderStatus = "CREATED"
	StatusExpired        OrderStatus = "EXPIRED"
	StatusReserved	  	 OrderStatus = "RESERVED"
	StatusPaymentPending OrderStatus = "PAYMENT_PENDING"
	StatusPaid           OrderStatus = "PAID"
	StatusCancelled      OrderStatus = "CANCELLED"
)

var ExpirableOrderStatuses = []OrderStatus{
	StatusReserved,
	StatusPaymentPending,
}