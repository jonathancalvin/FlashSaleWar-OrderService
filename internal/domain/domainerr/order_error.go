package domainerr

import "errors"

var (
	ErrInvalidOrderTransition = errors.New("invalid order state transition")
	ErrOrderNotFound          = errors.New("order not found")
	ErrOrderExpired           = errors.New("order already expired")
	ErrOrderUnauthorized      = errors.New("order does not belong to this user")
)
