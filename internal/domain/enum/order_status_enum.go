package enum

import "github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/domainerr"

type OrderStatus string

const (
	StatusCreated        OrderStatus = "CREATED"
	StatusReserved       OrderStatus = "RESERVED"
	StatusPaymentPending OrderStatus = "PAYMENT_PENDING"
	StatusPaid           OrderStatus = "PAID"
	StatusCancelled      OrderStatus = "CANCELLED"
	StatusExpired        OrderStatus = "EXPIRED"
)

var allowedOrderTransitions = map[OrderStatus]map[OrderStatus]bool{
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
	if next, ok := allowedOrderTransitions[from]; ok {
		return next[to]
	}
	return false
}

func ValidateTransition(from, to OrderStatus) error {
	if !CanTransition(from, to) {
		return domainerr.ErrInvalidOrderTransition
	}
	return nil
}