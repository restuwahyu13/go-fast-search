package dto

type (
	CreateUsersDTO struct {
		Name        string `json:"name" validate:"required"`
		Email       string `json:"email" validate:"required,email"`
		Phone       string `json:"phone" validate:"required"`
		DateOfBirth string `json:"date_of_birth" validate:"required,time=2006-01-02"`
		Age         int    `json:"age" validate:"required"`
		Address     string `json:"address" validate:"required"`
		City        string `json:"city" validate:"required"`
		State       string `json:"state" validate:"required"`
		Direction   string `json:"direction" validate:"required"`
		Country     string `json:"country" validate:"required"`
		PostalCode  string `json:"postal_code" validate:"required"`
	}
)
