package enum

type PaymentFailureReason string

const (
	Declined PaymentFailureReason = "DECLINED"
	Timeout  PaymentFailureReason = "TIMEOUT"
)

func (r PaymentFailureReason) IsValid() bool {
	switch r {
	case Declined, Timeout:
		return true
	default:
		return false
	}
}