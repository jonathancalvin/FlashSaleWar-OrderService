package domain

type Status string

const (
	StatusCreated        Status = "CREATED"
	StatusReserved       Status = "RESERVED"
	StatusPaymentPending Status = "PAYMENT_PENDING"
	StatusPaid           Status = "PAID"
	StatusCancelled      Status = "CANCELLED"
)
