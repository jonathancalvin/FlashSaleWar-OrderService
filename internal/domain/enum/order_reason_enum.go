package enum

type CancelReason string

const (
    UserCancel     CancelReason = "USER_CANCEL"
    StockExhausted CancelReason = "STOCK_EXHAUSTED"
)

type ExpireReason string

const (
    PaymentTTL     ExpireReason = "PAYMENT_TTL"
    ReservationTTL ExpireReason = "RESERVATION_TTL"
)