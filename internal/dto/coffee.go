package dto

type CreateCoffeeRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
}
type UpdateCoffeeRequest struct {
	ID          uint    `params:"id" validate:"required"`
	Name        string  `json:"name" validate:"omitempty,gte=1"`
	Description string  `json:"description" validate:"omitempty,gte=1"`
	Price       float64 `json:"price" validate:"omitempty,gt=0"`
}

type DeleteCoffeeRequest struct {
	ID uint `params:"id" validate:"required"`
}

type GetCoffeeRequest struct {
	ID uint `params:"id" validate:"required"`
}

type ListCoffeeRequest struct{}
