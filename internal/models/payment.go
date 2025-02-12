package models

import "coffee-shop-api/internal/common"

type Payment struct {
	ID      uint                 `gorm:"primaryKey" json:"id"`
	OrderID uint                 `gorm:"not null" json:"order_id"`
	Amount  uint                 `gorm:"not null" json:"amount"`
	Status  common.PaymentStatus `gorm:"not null" json:"status"`
}

// TableName sets the table name for the Payment model
func (Payment) TableName() string {
	return "payments"
}
