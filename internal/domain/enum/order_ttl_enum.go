package enum

import "time"

type OrderTTLType string

const (
	OrderTTLCreated        OrderTTLType = "CREATED"
	OrderTTLReserved       OrderTTLType = "RESERVED"
	OrderTTLPaymentPending OrderTTLType = "PAYMENT_PENDING"
)

var ttlConfig = map[OrderStatus]time.Duration{
    StatusCreated:        5 * time.Minute,
    StatusReserved:       2 * time.Minute,
    StatusPaymentPending: 10 * time.Minute,
}

func GetTTL(status OrderStatus) (time.Duration, bool) {
    duration, exists := ttlConfig[status]
    return duration, exists
}

func CalculateExpiryTime(status OrderStatus) *time.Time {
	if ttl, exists := GetTTL(status); exists {
		expiry := time.Now().Add(ttl)
		return &expiry
	}
	return nil
}