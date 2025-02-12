package models

import (
	"time"
)

// Order represents a coffee order in the database
type Order struct {
	ID        uint      `gorm:"primaryKey"`
	OrderDate time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Total     float64   `gorm:"type:decimal(10,2)"`
	Status    string    `gorm:"type:varchar(20);default:'pending'"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Coffees   []Coffee `gorm:"many2many:order_coffees"`
}

// TableName sets the table name for the Coffee model
func (Order) TableName() string {
	return "orders"
}

// OrderCoffee represents the junction table between Order and Coffee
type OrderCoffee struct {
	OrderID   uint    `gorm:"primaryKey"`
	CoffeeID  uint    `gorm:"primaryKey"`
	Quantity  int     `gorm:"not null;default:1"`
	Subtotal  float64 `gorm:"type:decimal(10,2)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName sets the table name for the OrderCoffee model
func (OrderCoffee) TableName() string {
	return "order_coffees"
}
