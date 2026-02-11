package enum

import "errors"

type CancelReason string

const (
    UserCancel     CancelReason = "USER_CANCEL"
    StockExhausted CancelReason = "STOCK_EXHAUSTED"
)

func (r CancelReason) IsValid() error {
    switch r {
    case UserCancel, StockExhausted:
        return nil
    default:
        return errors.New("invalid cancel reason value must be USER_CANCEL or STOCK_EXHAUSTED")
    }
}

type ExpireReason string

const (
    PaymentTTL     ExpireReason = "PAYMENT_TTL"
    ReservationTTL ExpireReason = "RESERVATION_TTL"
)