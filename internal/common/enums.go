package common

type (
	PaymentStatus string
)

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusPaid    PaymentStatus = "paid"
	PaymentStatusFailed  PaymentStatus = "failed"
)
