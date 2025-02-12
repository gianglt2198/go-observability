package models

// Coffee represents the coffee product in the database
type Coffee struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `gorm:"not null" json:"name"`
	Price       float64 `gorm:"not null" json:"price"`
	Description string  `gorm:"not null" json:"description"`
}

// TableName sets the table name for the Coffee model
func (Coffee) TableName() string {
	return "coffees"
}
