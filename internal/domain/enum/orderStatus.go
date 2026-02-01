package enum

import "errors"

type OrderStatus string

const (
	StatusCreated        OrderStatus = "CREATED"
	StatusReserved       OrderStatus = "RESERVED"
	StatusPaymentPending OrderStatus = "PAYMENT_PENDING"
	StatusPaid           OrderStatus = "PAID"
	StatusCancelled      OrderStatus = "CANCELLED"
	StatusExpired        OrderStatus = "EXPIRED"
)

var ErrInvalidOrderTransition = errors.New("invalid order state transition")

var AllowedOrderTransitions = map[OrderStatus]map[OrderStatus]bool{
	StatusCreated: {
		StatusReserved:  true,
		StatusCancelled: true,
	},
	StatusReserved: {
		StatusPaymentPending: true,
		StatusCancelled:      true,
		StatusExpired:        true,
	},
	StatusPaymentPending: {
		StatusPaid:      true,
		StatusCancelled: true,
		StatusExpired:   true,
	},
}

var ExpirableOrderStatuses = []OrderStatus{
	StatusReserved,
	StatusPaymentPending,
}

func CanTransition(from, to OrderStatus) bool {
	if next, ok := AllowedOrderTransitions[from]; ok {
		return next[to]
	}
	return false
}

func ValidateTransition(from, to OrderStatus) error {
	if !CanTransition(from, to) {
		return ErrInvalidOrderTransition
	}
	return nil
}
